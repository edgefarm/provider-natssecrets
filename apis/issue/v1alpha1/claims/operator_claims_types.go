/*
Copyright 2018-2019 The NATS Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/////////////////////////////////////////////////////////////////////////////////
// IMPORTANT NOTE: This file is a copy of the types from the NATS go repo.     //
// Modifications are necessary to be able to generate the CRDs:                //
// * Some Data needs json tags `json:",inline"`                                //
// Some Types are not directly supported and therefore have been commented out://
// * SigningKeys are not supported                                             //
/////////////////////////////////////////////////////////////////////////////////

package claims

// OperatorClaims define the data for an operator JWT
type OperatorClaims struct {
	ClaimsData `json:",inline"`        // mapstructure:",inline"`
	Operator   `json:"nats,omitempty"` // mapstructure:"nats,omitempty"`
}

// Operator specific claims
type Operator struct {
	// TODO:
	// // Slice of other operator NKeys that can be used to sign on behalf of the main
	// // operator identity.
	// SigningKeys StringList `json:"signing_keys,omitempty"`
	// AccountServerURL is a partial URL like "https://host.domain.org:<port>/jwt/v1"
	// tools will use the prefix and build queries by appending /accounts/<account_id>
	// or /operator to the path provided. Note this assumes that the account server
	// can handle requests in a nats-account-server compatible way. See
	// https://github.com/nats-io/nats-account-server.
	AccountServerURL string `json:"account_server_url,omitempty"` // mapstructure:"account_server_url,omitempty"`
	// A list of NATS urls (tls://host:port) where tools can connect to the server
	// using proper credentials.
	OperatorServiceURLs StringList `json:"operator_service_urls,omitempty"` // mapstructure:"operator_service_urls,omitempty"`
	// Identity of the system account
	SystemAccount string `json:"system_account,omitempty"` // mapstructure:"system_account,omitempty"`
	// Min Server version
	AssertServerVersion string `json:"assert_server_version,omitempty"` // mapstructure:"assert_server_version,omitempty"`
	// Signing of subordinate objects will require signing keys
	StrictSigningKeyUsage bool             `json:"strict_signing_key_usage,omitempty"` // mapstructure:"strict_signing_key_usage,omitempty"`
	GenericFields         `json:",inline"` // mapstructure:",inline"`
}

// StringList is a wrapper for an array of strings
type StringList []string

type GenericFields struct {
	Tags    TagList   `json:"tags,omitempty"`    // mapstructure:"tags,omitempty"`
	Type    ClaimType `json:"type,omitempty"`    // mapstructure:"type,omitempty"`
	Version int       `json:"version,omitempty"` // mapstructure:"version,omitempty"`
}

// TagList is a unique array of lower case strings
// All tag list methods lower case the strings in the arguments
type TagList []string

// ClaimType is used to indicate the type of JWT being stored in a Claim
type ClaimType string
