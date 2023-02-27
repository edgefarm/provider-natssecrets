package issue

import (
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"

	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func accountPath(mount string, operator string, account string) string {
	return mount + "/issue/operator/" + operator + "/account/" + account
}

func ReadAccount(c *vault.Client, operator string, account string) (*natsbackend.IssueAccountParameters, error) {
	path := accountPath(c.Mount, operator, account)
	return vault.Read[natsbackend.IssueAccountParameters](c, path)
}

func WriteAccount(c *vault.Client, operator string, account string, params *natsbackend.IssueAccountParameters) error {
	path := accountPath(c.Mount, operator, account)
	return vault.Write(c, path, params)
}

func DeleteAccount(c *vault.Client, operator string, account string) error {
	path := accountPath(c.Mount, operator, account)
	return vault.Delete(c, path)
}
