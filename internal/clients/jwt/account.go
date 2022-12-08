package jwt

import (
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func AccountPath(mount string, operator string, account string) string {
	return mount + "/jwt/operator/" + operator + "/account/" + account
}

func ReadAccount(c *vault.Client, operator string, account string) (*JWTData, error) {
	path := AccountPath(c.Mount, operator, account)
	return vault.Read[JWTData](c, path)
}
