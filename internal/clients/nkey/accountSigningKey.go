package nkey

import (
	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/accountSigningKey/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func accountSigningKeyPath(mount string, operator string, account string, key string) string {
	return mount + "/nkey/operator/" + operator + "/account/" + account + "/signing/" + key
}

func ReadAccountSigningKey(c *vault.Client, operator string, account string, key string) (*v1alpha1.AccountSigningKeyParameters, error) {
	path := accountSigningKeyPath(c.Mount, operator, account, key)
	return vault.Read[v1alpha1.AccountSigningKeyParameters](c, path)
}

func WriteAccountSigningKey(c *vault.Client, operator string, account string, key string, params *v1alpha1.AccountSigningKeyParameters) error {
	path := accountSigningKeyPath(c.Mount, operator, account, key)
	return vault.Write(c, path, params)
}

func DeleteAccountSigningKey(c *vault.Client, operator string, account string, key string) error {
	path := accountSigningKeyPath(c.Mount, operator, account, key)
	return vault.Delete(c, path)
}
