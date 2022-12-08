package claims

// UserClaims defines a user JWT
type UserClaims struct {
	ClaimsData `json:",inline"`        // mapstructure:",inline"`
	User       `json:"nats,omitempty"` // mapstructure:"nats,omitempty"`
}

// User defines the user specific data in a user JWT
type User struct {
	UserPermissionLimits `json:",inline"` // mapstructure:",inline"`
	// IssuerAccount stores the public key for the account the issuer represents.
	// When set, the claim was issued by a signing key.
	IssuerAccount string           `json:"issuer_account,omitempty"` // mapstructure:"issuer_account,omitempty"`
	GenericFields `json:",inline"` // mapstructure:",inline"`
}

type UserPermissionLimits struct {
	Permissions            `json:",inline"` // mapstructure:",inline"`
	Limits                 `json:",inline"` // mapstructure:",inline"`
	BearerToken            bool             `json:"bearer_token,omitempty"`             // mapstructure:"bearer_token,omitempty"`
	AllowedConnectionTypes StringList       `json:"allowed_connection_types,omitempty"` // mapstructure:"allowed_connection_types,omitempty"`
}

// Limits are used to control acccess for users and importing accounts
type Limits struct {
	UserLimits `json:",inline"` // mapstructure:",inline"`
	NatsLimits `json:",inline"` // mapstructure:",inline"`
}

// Src is a comma separated list of CIDR specifications
type UserLimits struct {
	Src    CIDRList    `json:"src,omitempty"`            // mapstructure:"src,omitempty"`
	Times  []TimeRange `json:"times,omitempty"`          // mapstructure:"times,omitempty"`
	Locale string      `json:"times_location,omitempty"` // mapstructure:"times_location,omitempty"`
}

// TimeRange is used to represent a start and end time
type TimeRange struct {
	Start string `json:"start,omitempty"` // mapstructure:"start,omitempty"`
	End   string `json:"end,omitempty"`   // mapstructure:"end,omitempty"`
}

type CIDRList TagList
