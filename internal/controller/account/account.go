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

package account

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/edgefarm/provider-natssecrets/apis/account/v1alpha1"
	"github.com/edgefarm/provider-natssecrets/internal/clients/issue"
	"github.com/edgefarm/provider-natssecrets/internal/clients/jwt"
	"github.com/edgefarm/provider-natssecrets/internal/clients/nkey"
	"github.com/edgefarm/provider-natssecrets/internal/controller/features"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	apisv1alpha1 "github.com/edgefarm/provider-natssecrets/apis/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

const (
	errNotAccount   = "managed resource is not a Account custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

// Setup adds a controller that reconciles Account managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.AccountGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.AccountGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: vault.NewRootClient,
		}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.Account{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte) (*vault.Client, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Account)
	if !ok {
		return nil, errors.New(errNotAccount)
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

	return &external{client: client}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	client *vault.Client
}

const (
	annotationExternalName = "crossplane.io/external-name"
)

func getExternalName(r *v1alpha1.Account) (string, error) {
	annotations := r.GetAnnotations()
	if annotations != nil {
		if val, ok := annotations[annotationExternalName]; ok {
			return val, nil
		}
	}
	return "", fmt.Errorf("External name annotation not found for %s", r.GetName())
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Account)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotAccount)
	}

	operator := cr.Spec.ForProvider.Operator
	account, err := getExternalName(cr)
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	data, status, err := issue.ReadAccount(c.client, operator, account)
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

	if !reflect.DeepEqual(data, &cr.Spec.ForProvider) {
		return managed.ExternalObservation{
			ResourceExists:    true,
			ResourceUpToDate:  false,
			ConnectionDetails: managed.ConnectionDetails{},
		}, nil
	}

	// receive nkey informations from vault
	details := managed.ConnectionDetails{}
	nk, err := nkey.ReadAccount(c.client, operator, account)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}
	if nk == nil {
		cr.SetConditions(xpv1.Creating().WithMessage("Waiting for account nkey to be created"))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	// receive jwt informations from vault
	j, err := jwt.ReadAccount(c.client, operator, account)
	if err != nil {
		cr.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}
	if j == nil {
		cr.SetConditions(xpv1.Creating().WithMessage("Waiting for account JWT to be created"))
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
		}, nil
	}

	// set connection details
	details["pub"] = []byte(nk.PublicKey)
	details["jwt"] = []byte(j.JWT)

	cr.Status.AtProvider.Operator = operator
	cr.Status.AtProvider.Account = account
	cr.Status.AtProvider.Issue = issue.AccountPath(c.client.Mount, operator, account)
	cr.Status.AtProvider.NKeyPath = nkey.AccountPath(c.client.Mount, operator, account)
	cr.Status.AtProvider.JWTPath = jwt.AccountPath(c.client.Mount, operator, account)
	cr.Status.AtProvider.JWT = func() string {
		if status.Account.JWT {
			return "true"
		}
		return "false"
	}()
	cr.Status.AtProvider.NKey = func() string {
		if status.Account.Nkey {
			return "true"
		}
		return "false"
	}()
	cr.Status.AtProvider.Pushed = func() string {
		if status.AccountServer.Synced {
			return "true"
		}
		return "false"
	}()
	cr.Status.AtProvider.LastPushed = func() string {
		if status.AccountServer.LastSync > 0 {
			return time.Unix(status.AccountServer.LastSync, 0).Format(time.RFC3339)
		}
		return "never"
	}()
	cr.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  true,
		ConnectionDetails: details,
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Account)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotAccount)
	}

	cr.SetConditions(xpv1.Creating())

	operator := cr.Spec.ForProvider.Operator
	account, err := getExternalName(cr)
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	err = issue.WriteAccount(c.client, operator, account, &cr.Spec.ForProvider)
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
	cr, ok := mg.(*v1alpha1.Account)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotAccount)
	}

	operator := cr.Spec.ForProvider.Operator
	account, err := getExternalName(cr)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	// check if operator id has changed
	if cr.Status.AtProvider.Operator != "" && cr.Status.AtProvider.Account != "" {
		// when changed, delete old operator
		if cr.Status.AtProvider.Operator != operator {
			err := issue.DeleteOperator(c.client, cr.Status.AtProvider.Operator)
			if err != nil {
				return managed.ExternalUpdate{}, err
			}
		}
	}

	err = issue.WriteAccount(c.client, operator, account, &cr.Spec.ForProvider)
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
	cr, ok := mg.(*v1alpha1.Account)
	if !ok {
		return errors.New(errNotAccount)
	}

	operator := cr.Spec.ForProvider.Operator
	account, err := getExternalName(cr)
	if err != nil {
		return err
	}

	return issue.DeleteAccount(c.client, operator, account)
}
