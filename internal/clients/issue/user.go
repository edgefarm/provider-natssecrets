package issue

import (
	"fmt"

	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/user/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"
	vaultv1alpha1 "github.com/edgefarm/vault-plugin-secrets-nats/pkg/claims/user/v1alpha1"
)

func UserPath(mount string, operator string, account string, user string) string {
	return mount + "/issue/operator/" + operator + "/account/" + account + "/user/" + user
}

func fixEmptySlices(params *vaultv1alpha1.UserClaims) {
	if params == nil {
		return
	}
	if params.Permissions.Pub.Allow == nil {
		params.Permissions.Pub.Allow = []string{}
	}
	if params.Permissions.Pub.Deny == nil {
		params.Permissions.Pub.Deny = []string{}
	}
	if params.Permissions.Sub.Allow == nil {
		params.Permissions.Sub.Allow = []string{}
	}
	if params.Permissions.Sub.Deny == nil {
		params.Permissions.Sub.Deny = []string{}
	}
}

func ReadUser(c *vault.Client, operator string, account string, user string) (*v1alpha1.UserParameters, *natsbackend.IssueUserStatus, error) {
	path := UserPath(c.Mount, operator, account, user)

	resp, err := vault.Read[natsbackend.IssueUserData](c, path)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		fixEmptySlices(&resp.Claims)
		return &v1alpha1.UserParameters{
			Operator:      resp.Operator,
			Account:       resp.Account,
			Claims:        resp.Claims,
			UseSigningKey: resp.UseSigningKey,
		}, &resp.Status, nil
	}
	return nil, nil, fmt.Errorf("user %s in account %s in operator %s not found", user, account, operator)
}

func WriteUser(c *vault.Client, operator string, account string, user string, params *v1alpha1.UserParameters) error {
	path := UserPath(c.Mount, operator, account, user)
	request := &natsbackend.IssueUserParameters{
		Operator:      operator,
		Account:       account,
		User:          user,
		Claims:        params.Claims,
		UseSigningKey: params.UseSigningKey,
	}
	return vault.Write(c, path, request)
}

func DeleteUser(c *vault.Client, operator string, account string, user string) error {
	path := UserPath(c.Mount, operator, account, user)
	return vault.Delete(c, path)
}
