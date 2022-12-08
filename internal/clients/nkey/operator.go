package nkey

import (
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func OperatorPath(mount string, operator string) string {
	return mount + "/nkey/operator/" + operator
}

func ReadOperator(c *vault.Client, operator string) (*NkeyData, error) {
	path := OperatorPath(c.Mount, operator)
	return vault.Read[NkeyData](c, path)
}
