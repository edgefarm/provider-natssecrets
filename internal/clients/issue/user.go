package issue

import (
	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/user/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"
)

func UserPath(mount string, operator string, account string, user string) string {
	return mount + "/issue/operator/" + operator + "/account/" + account + "/user/" + user
}

func ReadUser(c *vault.Client, operator string, account string, user string) (*v1alpha1.UserParameters, error) {
	path := UserPath(c.Mount, operator, account, user)

	resp, err := vault.Read[natsbackend.IssueUserParameters](c, path)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.UserParameters{
		Operator: resp.Operator,
		Account:  resp.Account,
		Claims:   resp.Claims,
	}, nil
}

func WriteUser(c *vault.Client, operator string, account string, user string, params *v1alpha1.UserParameters) error {
	path := UserPath(c.Mount, operator, account, user)
	request := &natsbackend.IssueUserParameters{
		Operator: operator,
		Account:  account,
		User:     user,
		Claims:   params.Claims,
	}
	return vault.Write(c, path, request)
}

func DeleteUser(c *vault.Client, operator string, account string, user string) error {
	path := UserPath(c.Mount, operator, account, user)
	return vault.Delete(c, path)
}
