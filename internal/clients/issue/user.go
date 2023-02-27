package issue

import (
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"

	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func userPath(mount string, operator string, account string, user string) string {
	return mount + "/issue/operator/" + operator + "/account/" + account + "/user/" + user
}

func ReadUser(c *vault.Client, operator string, account string, user string) (*natsbackend.IssueUserParameters, error) {
	path := userPath(c.Mount, operator, account, user)
	return vault.Read[natsbackend.IssueUserParameters](c, path)
}

func WriteUser(c *vault.Client, operator string, account string, user string, params *natsbackend.IssueUserParameters) error {
	path := userPath(c.Mount, operator, account, user)
	return vault.Write(c, path, params)
}

func DeleteUser(c *vault.Client, operator string, account string, user string) error {
	path := userPath(c.Mount, operator, account, user)
	return vault.Delete(c, path)
}
