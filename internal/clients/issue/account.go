package issue

import (
	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/issue/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func accountPath(mount string, operator string, account string) string {
	return mount + "/issue/operator/" + operator + "/account/" + account
}

func ReadAccount(c *vault.Client, operator string, account string) (*v1alpha1.AccountParameters, error) {
	path := accountPath(c.Mount, operator, account)
	return vault.Read[v1alpha1.AccountParameters](c, path)
}

func WriteAccount(c *vault.Client, operator string, account string, params *v1alpha1.AccountParameters) error {
	path := accountPath(c.Mount, operator, account)
	return vault.Write(c, path, params)
}

func DeleteAccount(c *vault.Client, operator string, account string) error {
	path := accountPath(c.Mount, operator, account)
	return vault.Delete(c, path)
}
