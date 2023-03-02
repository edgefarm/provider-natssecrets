package nkey

import (
	"fmt"

	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/operatorSigningKey/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"
)

func OperatorSigningKeyPath(mount string, operator string, key string) string {
	return mount + "/nkey/operator/" + operator + "/signing/" + key
}

func ReadOperatorSigningKey(c *vault.Client, operator string, key string) (*natsbackend.NkeyParameters, bool, error) {
	path := OperatorSigningKeyPath(c.Mount, operator, key)
	data, err := vault.Read[natsbackend.NkeyData](c, path)
	if err != nil {
		return nil, false, err
	}
	ret := &natsbackend.NkeyParameters{}
	status := false
	if data != nil {
		ret = &natsbackend.NkeyParameters{
			Seed: data.Seed,
		}
		status = data.Seed != ""
		return ret, status, nil
	}
	return nil, false, fmt.Errorf("operator signing key %s not found", key)
}

func WriteOperatorSigningKey(c *vault.Client, operator string, key string, params *v1alpha1.OperatorSigningKeyParameters) error {
	path := OperatorSigningKeyPath(c.Mount, operator, key)

	seed := ""
	if params.Config.Import.SecretRef != nil {
		if params.Config.Import.SecretRef.Key == "" {
			return fmt.Errorf("secret key is missing")
		}
		if params.Config.Import.SecretRef.Name == "" {
			return fmt.Errorf("secret name is missing")
		}
		if params.Config.Import.SecretRef.Namespace == "" {
			return fmt.Errorf("secret namespace is missing")
		}
		var err error
		seed, err = GetSeedFromSecret(params.Config.Import.SecretRef.Namespace, params.Config.Import.SecretRef.Name, params.Config.Import.SecretRef.Key)
		if err != nil {
			return err
		}
	}
	request := &natsbackend.NkeyParameters{
		Seed: seed,
	}

	return vault.Write(c, path, request)
}

func DeleteOperatorSigningKey(c *vault.Client, operator string, key string) error {
	path := OperatorSigningKeyPath(c.Mount, operator, key)
	return vault.Delete(c, path)
}
