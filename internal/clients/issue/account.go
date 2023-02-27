package issue

import (
	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/account/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"
)

func AccountPath(mount string, operator string, account string) string {
	return mount + "/issue/operator/" + operator + "/account/" + account
}

func ReadAccount(c *vault.Client, operator string, account string) (*v1alpha1.AccountParameters, error) {
	path := AccountPath(c.Mount, operator, account)

	resp, err := vault.Read[natsbackend.IssueAccountParameters](c, path)
	if err != nil {
		return nil, err
	}

	response := &v1alpha1.AccountParameters{
		Operator: resp.Operator,
		Claims:   resp.Claims,
	}
	return response, nil
}

func WriteAccount(c *vault.Client, operator string, account string, params *v1alpha1.AccountParameters) error {
	path := AccountPath(c.Mount, operator, account)
	request := &natsbackend.IssueAccountParameters{
		Operator: operator,
		Account:  account,
		Claims:   params.Claims,
	}
	return vault.Write(c, path, request)
}

func DeleteAccount(c *vault.Client, operator string, account string) error {
	path := AccountPath(c.Mount, operator, account)
	return vault.Delete(c, path)
}
