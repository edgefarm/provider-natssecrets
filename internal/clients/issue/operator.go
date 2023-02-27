package issue

import (
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"

	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func OperatorPath(mount string, operator string) string {
	return mount + "/issue/operator/" + operator
}

func ReadOperator(c *vault.Client, operator string) (*natsbackend.IssueOperatorParameters, error) {
	path := OperatorPath(c.Mount, operator)
	return vault.Read[natsbackend.IssueOperatorParameters](c, path)
}

func WriteOperator(c *vault.Client, operator string, params *natsbackend.IssueOperatorParameters) error {
	path := OperatorPath(c.Mount, operator)
	return vault.Write(c, path, params)
}

func DeleteOperator(c *vault.Client, operator string) error {
	path := OperatorPath(c.Mount, operator)
	return vault.Delete(c, path)
}
