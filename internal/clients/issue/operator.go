package issue

import (
	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/operator/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func OperatorPath(mount string, operator string) string {
	return mount + "/issue/operator/" + operator
}

func ReadOperator(c *vault.Client, operator string) (*v1alpha1.OperatorParameters, error) {
	path := OperatorPath(c.Mount, operator)
	return vault.Read[v1alpha1.OperatorParameters](c, path)
}

func WriteOperator(c *vault.Client, operator string, params *v1alpha1.OperatorParameters) error {
	path := OperatorPath(c.Mount, operator)
	return vault.Write(c, path, params)
}

func DeleteOperator(c *vault.Client, operator string) error {
	path := OperatorPath(c.Mount, operator)
	return vault.Delete(c, path)
}
