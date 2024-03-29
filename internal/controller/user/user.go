/*
Copyright 2022 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

	"github.com/edgefarm/provider-natssecrets/apis/user/v1alpha1"
	apisv1alpha1 "github.com/edgefarm/provider-natssecrets/apis/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
	"github.com/edgefarm/provider-natssecrets/internal/clients/creds"
	"github.com/edgefarm/provider-natssecrets/internal/clients/issue"
	"github.com/edgefarm/provider-natssecrets/internal/clients/jwt"
	"github.com/edgefarm/provider-natssecrets/internal/clients/nkey"
	"github.com/edgefarm/provider-natssecrets/internal/controller/features"

	deep "github.com/go-test/deep"
)

const (
	errNotUser      = "managed resource is not a User custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

// Setup adds a controller that reconciles User managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.UserGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.UserGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: vault.NewRootClient,
			logger:       o.Logger}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.User{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte) (*vault.Client, error)
	logger       logging.Logger
}

type ProviderConfigSecretSpec struct {
	// JWT is the NATS users JWT token to use for authentication.
	JWT string `json:"jwt"`
	// SeedKey is the NATS users seed key to use for authentication.
	SeedKey string `json:"seed_key"`
	// Address is the NATS address to use for authentication.
	Address string `json:"address"`
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return nil, errors.New(errNotUser)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	client, err := c.newServiceFn(data)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &external{
		client: client,
		log:    c.logger,
	}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	client *vault.Client
	log    logging.Logger
}

const (
	annotationExternalName = "crossplane.io/external-name"
)

func getExternalName(r *v1alpha1.User) (string, error) {
	annotations := r.GetAnnotations()
	if annotations != nil {
		if val, ok := annotations[annotationExternalName]; ok {
			return val, nil
		}
	}
	return "", fmt.Errorf("External name annotation not found for %s", r.GetName())
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotUser)
	}
	v1alpha1.FixEmptySlices(&cr.Spec.ForProvider.Claims)
	user, err := getExternalName(cr)
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	operator := cr.Spec.ForProvider.Operator
	account := cr.Spec.ForProvider.Account
	data, status, err := issue.ReadUser(c.client, operator, account, user)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	if data == nil {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	diff := deep.Equal(*data, cr.Spec.ForProvider)
	if diff != nil {
		c.log.Debug("Observe", "user", user)
		c.log.Debug("Compare failed", "diff", diff)
	}
	if !reflect.DeepEqual(data, &cr.Spec.ForProvider) {
		return managed.ExternalObservation{
			ResourceExists:    true,
			ResourceUpToDate:  false,
			ConnectionDetails: managed.ConnectionDetails{},
		}, nil
	}

	// receive nkey informations from vault
	details := managed.ConnectionDetails{}
	nk, err := nkey.ReadUser(c.client, operator, account, user)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}
	if nk == nil {
		cr.SetConditions(xpv1.Creating().WithMessage("Waiting for user nkey to be created"))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	accountNK, err := nkey.ReadAccount(c.client, operator, account)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}
	if accountNK == nil {
		cr.SetConditions(xpv1.Creating().WithMessage("Waiting for account nkey to be created"))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	// receive jwt informations from vault
	j, err := jwt.ReadUser(c.client, operator, account, user)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}
	if j == nil {
		cr.SetConditions(xpv1.Creating().WithMessage("Waiting for user JWT to be created"))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	// receive creds informations from vault
	userCreds, err := creds.Read(c.client, operator, account, user)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}
	if userCreds == nil {
		cr.SetConditions(xpv1.Creating().WithMessage("Waiting for user creds to be created"))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	operatorIssue, _, err := issue.ReadOperator(c.client, operator)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}
	if operatorIssue == nil {
		cr.SetConditions(xpv1.Creating().WithMessage("Waiting for operator issue to be created"))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	operatorJWT, err := jwt.ReadOperator(c.client, operator)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}

	sysaccountName := operatorIssue.Claims.SystemAccount
	if sysaccountName == "" {
		cr.SetConditions(xpv1.Unavailable().WithMessage("System account not found"))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}

	sysAccountJWT, err := jwt.ReadAccount(c.client, operator, sysaccountName)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}

	sysAccountNK, err := nkey.ReadAccount(c.client, operator, sysaccountName)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}

	// set connection details
	details["creds"] = []byte(userCreds.Creds)
	details["account-public-key"] = []byte(accountNK.PublicKey)
	details["address"] = []byte(operatorIssue.Claims.OperatorServiceURLs[0])

	if account != sysaccountName {
		details["seed"] = []byte(nk.Seed)
		details["jwt"] = []byte(j.JWT)
		details["dapr"] = []byte(fmt.Sprintf(`apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: streams
spec:
  type: pubsub.jetstream
  version: v1
  metadata:
  - name: natsURL
    value: "nats://localhost:4222"
  - name: jwt
    value: "%s"
  - name: seedKey
    value: "%s"`, j.JWT, nk.Seed))

		// Use the first operator service url as the address for now.
		// TODO: Support multiple operator service urls.
		if operatorIssue.Claims.OperatorServiceURLs != nil {
			ret, err := json.Marshal(ProviderConfigSecretSpec{
				Address: operatorIssue.Claims.OperatorServiceURLs[0],
				JWT:     j.JWT,
				SeedKey: nk.Seed,
			})
			if err != nil {
				return managed.ExternalObservation{}, errors.Wrap(err, "failed to marshal provider config secret")
			}
			details["provider-config-secret"] = ret
		}
	}

	if account == sysaccountName {
		details["operator-jwt"] = []byte(operatorJWT.JWT)
		details["sys-account-jwt"] = []byte(sysAccountJWT.JWT)
		details["sys-account-public-key"] = []byte(sysAccountNK.PublicKey)
	}

	cr.Status.AtProvider.Operator = operator
	cr.Status.AtProvider.Account = account
	cr.Status.AtProvider.User = user
	cr.Status.AtProvider.Issue = issue.UserPath(c.client.Mount, operator, account, user)
	cr.Status.AtProvider.NKeyPath = nkey.UserPath(c.client.Mount, operator, account, user)
	cr.Status.AtProvider.JWTPath = jwt.UserPath(c.client.Mount, operator, account, user)
	cr.Status.AtProvider.JWT = func() string {
		if status.User.JWT {
			return "true"
		}
		return "false"
	}()
	cr.Status.AtProvider.NKey = func() string {
		if status.User.Nkey {
			return "true"
		}
		return "false"
	}()
	cr.Status.AtProvider.Creds = creds.Path(c.client.Mount, operator, account, user)

	cr.SetConditions(xpv1.Available())
	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  true,
		ConnectionDetails: details,
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotUser)
	}

	operator := cr.Spec.ForProvider.Operator
	account := cr.Spec.ForProvider.Account
	user, err := getExternalName(cr)
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	err = issue.WriteUser(c.client, operator, account, user, &cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotUser)
	}

	operator := cr.Spec.ForProvider.Operator
	account := cr.Spec.ForProvider.Account
	user, err := getExternalName(cr)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}
	err = issue.WriteUser(c.client, operator, account, user, &cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return errors.New(errNotUser)
	}

	operator := cr.Spec.ForProvider.Operator
	account := cr.Spec.ForProvider.Account
	user, err := getExternalName(cr)
	if err != nil {
		return err
	}

	return issue.DeleteUser(c.client, operator, account, user)
}
