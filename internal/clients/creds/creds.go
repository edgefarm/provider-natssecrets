package creds

import (
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"

	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func Path(mount string, operator string, account string, user string) string {
	return mount + "/creds/operator/" + operator + "/account/" + account + "/user/" + user
}

func Read(c *vault.Client, operator string, account string, user string) (*natsbackend.CredsData, error) {
	path := Path(c.Mount, operator, account, user)
	return vault.Read[natsbackend.CredsData](c, path)
}
