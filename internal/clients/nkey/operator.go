package nkey

import (
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"

	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func OperatorPath(mount string, operator string) string {
	return mount + "/nkey/operator/" + operator
}

func ReadOperator(c *vault.Client, operator string) (*natsbackend.NkeyData, error) {
	path := OperatorPath(c.Mount, operator)
	return vault.Read[natsbackend.NkeyData](c, path)
}
