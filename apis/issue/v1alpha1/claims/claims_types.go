package claims

// ClaimsData is the base struct for all claims
type ClaimsData struct {
	Audience  string `json:"aud,omitempty"`  // mapstructure:"aud,omitempty"`
	Expires   int64  `json:"exp,omitempty"`  // mapstructure:"exp,omitempty"`
	ID        string `json:"jti,omitempty"`  // mapstructure:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`  // mapstructure:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`  // mapstructure:"iss,omitempty"`
	Name      string `json:"name,omitempty"` // mapstructure:"name,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`  // mapstructure:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`  // mapstructure:"sub,omitempty"`
}
