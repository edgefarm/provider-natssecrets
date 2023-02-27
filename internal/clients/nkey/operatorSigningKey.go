package nkey

import (
	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/operatorSigningKey/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func operatorSigningKeyPath(mount string, operator string, key string) string {
	return mount + "/nkey/operator/" + operator + "/signing/" + key
}

func ReadOperatorSigningKey(c *vault.Client, operator string, key string) (*v1alpha1.OperatorSigningKeyParameters, error) {
	path := operatorSigningKeyPath(c.Mount, operator, key)
	return vault.Read[v1alpha1.OperatorSigningKeyParameters](c, path)
}

func WriteOperatorSigningKey(c *vault.Client, operator string, key string, params *v1alpha1.OperatorSigningKeyParameters) error {
	path := operatorSigningKeyPath(c.Mount, operator, key)
	return vault.Write(c, path, params)
}

func DeleteOperatorSigningKey(c *vault.Client, operator string, key string) error {
	path := operatorSigningKeyPath(c.Mount, operator, key)
	return vault.Delete(c, path)
}
