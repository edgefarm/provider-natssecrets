package jwt

import (
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"

	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func UserPath(mount string, operator string, account string, user string) string {
	return mount + "/jwt/operator/" + operator + "/account/" + account + "/user/" + user
}

func ReadUser(c *vault.Client, operator string, account string, user string) (*natsbackend.JWTData, error) {
	path := UserPath(c.Mount, operator, account, user)
	return vault.Read[natsbackend.JWTData](c, path)
}
