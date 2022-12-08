package jwt

import (
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func OperatorPath(mount string, operator string) string {
	return mount + "/jwt/operator/" + operator
}

func ReadOperator(c *vault.Client, operator string) (*JWTData, error) {
	path := OperatorPath(c.Mount, operator)
	return vault.Read[JWTData](c, path)
}
