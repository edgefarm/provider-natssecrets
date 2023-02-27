package nkey

import (
	"fmt"

	natsbackend "github.com/edgefarm/vault-plugin-secrets-nats"

	v1alpha1 "github.com/edgefarm/provider-natssecrets/apis/accountSigningKey/v1alpha1"
	vault "github.com/edgefarm/provider-natssecrets/internal/clients"
)

func AccountSigningKeyPath(mount string, operator string, account string, key string) string {
	return mount + "/nkey/operator/" + operator + "/account/" + account + "/signing/" + key
}

func ReadAccountSigningKey(c *vault.Client, operator string, account string, key string) (*natsbackend.NkeyParameters, error) {
	path := AccountSigningKeyPath(c.Mount, operator, account, key)
	return vault.Read[natsbackend.NkeyParameters](c, path)
}

func WriteAccountSigningKey(c *vault.Client, operator string, account string, key string, params *v1alpha1.AccountSigningKeyParameters) error {
	path := AccountSigningKeyPath(c.Mount, operator, account, key)

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

func DeleteAccountSigningKey(c *vault.Client, operator string, account string, key string) error {
	path := AccountSigningKeyPath(c.Mount, operator, account, key)
	return vault.Delete(c, path)
}
