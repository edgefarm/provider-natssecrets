package vault

import (
	"errors"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/mitchellh/mapstructure"
)

var (
	ErrSecretNotFound = errors.New("no secret not found at given path")
	ErrVaultConfig    = errors.New("secret does not contain a valid vault configuration")
)

type Config struct {
	// Token is the Vault token to use for authentication.
	Token string `json:"token"`
	// Address is the Vault address to use for authentication.
	Address string `json:"address"`
	// Path is the mount root of the nats secrets engine.
	Path string `json:"path"`
}

type Logical interface {
	List(path string) (*api.Secret, error)
	Read(path string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
	Delete(path string) (*api.Secret, error)
}

type Client struct {
	Logical Logical
	Mount   string
}

func NewRootClient(creds []byte) (*Client, error) {

	var config Config
	if err := jsonutil.DecodeJSON(creds, &config); err != nil {
		return nil, err
	}

	if config.Token == "" || config.Address == "" {
		return nil, ErrVaultConfig
	}

	api, err := api.NewClient(&api.Config{Address: config.Address})
	if err != nil {
		return &Client{}, err
	}

	api.SetToken(config.Token)

	return &Client{
		Logical: api.Logical(),
		Mount:   config.Path,
	}, nil
}

func Write[T any](c *Client, path string, params *T) error {
	data := map[string]interface{}{}
	err := mapstructure.Decode(params, &data)
	if err != nil {
		return err
	}
	_, err = c.Logical.Write(path, data)
	return err
}

func Read[T any](c *Client, path string) (*T, error) {
	secret, err := c.Logical.Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	var params T
	err = mapstructure.Decode(secret.Data, &params)
	if err != nil {
		return nil, err
	}

	return &params, nil
}

func Delete(c *Client, path string) error {
	_, err := c.Logical.Delete(path)
	return err
}
