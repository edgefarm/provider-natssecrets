package nkey

type NkeyData struct {
	Seed       string `mapstructure:"seed,omitempty"`
	PublicKey  string `mapstructure:"public_key,omitempty"`
	PrivateKey string `mapstructure:"private_key,omitempty"`
}
