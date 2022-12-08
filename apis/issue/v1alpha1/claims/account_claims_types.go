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

import "time"

// AccountClaims defines the body of an account JWT
type AccountClaims struct {
	ClaimsData `json:",inline"`        // mapstructure:",inline"`
	Account    `json:"nats,omitempty"` // mapstructure:"nats,omitempty"`
}

// Account holds account specific claims data
type Account struct {
	Imports Imports        `json:"imports,omitempty"` // mapstructure:"imports,omitempty"`
	Exports Exports        `json:"exports,omitempty"` // mapstructure:"exports,omitempty"`
	Limits  OperatorLimits `json:"limits,omitempty"`  // mapstructure:"limits,omitempty"`
	// TODO
	// SigningKeys        SigningKeys    `json:"signing_keys,omitempty"`
	Revocations        RevocationList   `json:"revocations,omitempty"`         // mapstructure:"revocations,omitempty"`
	DefaultPermissions Permissions      `json:"default_permissions,omitempty"` // mapstructure:"default_permissions,omitempty"`
	Mappings           Mapping          `json:"mappings,omitempty"`            // mapstructure:"mappings,omitempty"`
	Info               `json:",inline"` // mapstructure:",inline"`
	GenericFields      `json:",inline"` // mapstructure:",inline"`
}

// Imports is a list of import structs
type Imports []*Import

// Exports is a slice of exports
type Exports []*Export

// OperatorLimits are used to limit access by an account
type OperatorLimits struct {
	NatsLimits            `json:",inline"`                 // mapstructure:",inline"`
	AccountLimits         `json:",inline"`                 // mapstructure:",inline"`
	JetStreamLimits       `json:",inline"`                 // mapstructure:",inline"`
	JetStreamTieredLimits `json:"tiered_limits,omitempty"` // mapstructure:"tiered_limits,omitempty"`
}

// Import describes a mapping from another account into this one
type Import struct {
	Name string `json:"name,omitempty"` // mapstructure:"name,omitempty"`
	// Subject field in an import is always from the perspective of the
	// initial publisher - in the case of a stream it is the account owning
	// the stream (the exporter), and in the case of a service it is the
	// account making the request (the importer).
	Subject Subject `json:"subject,omitempty"` // mapstructure:"subject,omitempty"`
	Account string  `json:"account,omitempty"` // mapstructure:"account,omitempty"`
	Token   string  `json:"token,omitempty"`   // mapstructure:"token,omitempty"`
	// Deprecated: use LocalSubject instead
	// To field in an import is always from the perspective of the subscriber
	// in the case of a stream it is the client of the stream (the importer),
	// from the perspective of a service, it is the subscription waiting for
	// requests (the exporter). If the field is empty, it will default to the
	// value in the Subject field.
	To Subject `json:"to,omitempty"` // mapstructure:"to,omitempty"`
	// Local subject used to subscribe (for streams) and publish (for services) to.
	// This value only needs setting if you want to change the value of Subject.
	// If the value of Subject ends in > then LocalSubject needs to end in > as well.
	// LocalSubject can contain $<number> wildcard references where number references the nth wildcard in Subject.
	// The sum of wildcard reference and * tokens needs to match the number of * token in Subject.
	LocalSubject RenamingSubject `json:"local_subject,omitempty"` // mapstructure:"local_subject,omitempty"`
	Type         ExportType      `json:"type,omitempty"`          // mapstructure:"type,omitempty"`
	Share        bool            `json:"share,omitempty"`         // mapstructure:"share,omitempty"`
}

// Export represents a single export
type Export struct {
	Name                 string           `json:"name,omitempty"`                   // mapstructure:"name,omitempty"`
	Subject              Subject          `json:"subject,omitempty"`                // mapstructure:"subject,omitempty"`
	Type                 ExportType       `json:"type,omitempty"`                   // mapstructure:"type,omitempty"`
	TokenReq             bool             `json:"token_req,omitempty"`              // mapstructure:"token_req,omitempty"`
	Revocations          RevocationList   `json:"revocations,omitempty"`            // mapstructure:"revocations,omitempty"`
	ResponseType         ResponseType     `json:"response_type,omitempty"`          // mapstructure:"response_type,omitempty"`
	ResponseThreshold    time.Duration    `json:"response_threshold,omitempty"`     // mapstructure:"response_threshold,omitempty"`
	Latency              *ServiceLatency  `json:"service_latency,omitempty"`        // mapstructure:"service_latency,omitempty"`
	AccountTokenPosition uint             `json:"account_token_position,omitempty"` // mapstructure:"account_token_position,omitempty"`
	Advertise            bool             `json:"advertise,omitempty"`              // mapstructure:"advertise,omitempty"`
	Info                 `json:",inline"` // mapstructure:",inline"`
}

type NatsLimits struct {
	Subs    int64 `json:"subs,omitempty"`    // mapstructure:"subs,omitempty"`       // Max number of subscriptions
	Data    int64 `json:"data,omitempty"`    // mapstructure:"data,omitempty"`       // Max number of bytes
	Payload int64 `json:"payload,omitempty"` // mapstructure:"payload,omitempty"` // Max message payload
}

type AccountLimits struct {
	Imports         int64 `json:"imports,omitempty"`         // mapstructure:"imports,omitempty"`                 // Max number of imports
	Exports         int64 `json:"exports,omitempty"`         // mapstructure:"exports,omitempty"`                 // Max number of exports
	WildcardExports bool  `json:"wildcards,omitempty"`       // mapstructure:"wildcards,omitempty"`             // Are wildcards allowed in exports
	DisallowBearer  bool  `json:"disallow_bearer,omitempty"` // mapstructure:"disallow_bearer,omitempty"` // User JWT can't be bearer token
	Conn            int64 `json:"conn,omitempty"`            // mapstructure:"conn,omitempty"`                       // Max number of active connections
	LeafNodeConn    int64 `json:"leaf,omitempty"`            // mapstructure:"leaf,omitempty"`                       // Max number of active leaf node connections
}

// Subject is a string that represents a NATS subject
type Subject string

type RenamingSubject Subject

// ExportType defines the type of import/export.
type ExportType int

type JetStreamLimits struct {
	MemoryStorage        int64 `json:"mem_storage,omitempty"`           // mapstructure:"mem_storage,omitempty"`                     // Max number of bytes stored in memory across all streams. (0 means disabled)
	DiskStorage          int64 `json:"disk_storage,omitempty"`          // mapstructure:"disk_storage,omitempty"`                   // Max number of bytes stored on disk across all streams. (0 means disabled)
	Streams              int64 `json:"streams,omitempty"`               // mapstructure:"streams,omitempty"`                             // Max number of streams
	Consumer             int64 `json:"consumer,omitempty"`              // mapstructure:"consumer,omitempty"`                           // Max number of consumers
	MaxAckPending        int64 `json:"max_ack_pending,omitempty"`       // mapstructure:"max_ack_pending,omitempty"`             // Max ack pending of a Stream
	MemoryMaxStreamBytes int64 `json:"mem_max_stream_bytes,omitempty"`  // mapstructure:"mem_max_stream_bytes,omitempty"`   // Max bytes a memory backed stream can have. (0 means disabled/unlimited)
	DiskMaxStreamBytes   int64 `json:"disk_max_stream_bytes,omitempty"` // mapstructure:"disk_max_stream_bytes,omitempty"` // Max bytes a disk backed stream can have. (0 means disabled/unlimited)
	MaxBytesRequired     bool  `json:"max_bytes_required,omitempty"`    // mapstructure:"max_bytes_required,omitempty"`       // Max bytes required by all Streams
}

// RevocationList is used to store a mapping of public keys to unix timestamps
type RevocationList map[string]int64

type JetStreamTieredLimits map[string]JetStreamLimits

// Permissions are used to restrict subject access, either on a user or for everyone on a server by default
type Permissions struct {
	Pub  Permission          `json:"pub,omitempty"`  // mapstructure:"pub,omitempty"`
	Sub  Permission          `json:"sub,omitempty"`  // mapstructure:"sub,omitempty"`
	Resp *ResponsePermission `json:"resp,omitempty"` // mapstructure:"resp,omitempty"`
}

// Permission defines allow/deny subjects
type Permission struct {
	Allow StringList `json:"allow,omitempty"` // mapstructure:"allow,omitempty"`
	Deny  StringList `json:"deny,omitempty"`  // mapstructure:"deny,omitempty"`
}

// ResponsePermission can be used to allow responses to any reply subject
// that is received on a valid subscription.
type ResponsePermission struct {
	MaxMsgs int           `json:"max"` // mapstructure:"max"`
	Expires time.Duration `json:"ttl"` // mapstructure:"ttl"`
}

// ResponseType is used to store an export response type
type ResponseType string

// ServiceLatency is used when observing and exported service for
// latency measurements.
// Sampling 1-100, represents sampling rate, defaults to 100.
// Results is the subject where the latency metrics are published.
// A metric will be defined by the nats-server's ServiceLatency. Time durations
// are in nanoseconds.
// see https://github.com/nats-io/nats-server/blob/main/server/accounts.go#L524
// e.g.
//
//	{
//	 "app": "dlc22",
//	 "start": "2019-09-16T21:46:23.636869585-07:00",
//	 "svc": 219732,
//	 "nats": {
//	   "req": 320415,
//	   "resp": 228268,
//	   "sys": 0
//	 },
//	 "total": 768415
//	}
type ServiceLatency struct {
	Sampling SamplingRate `json:"sampling"` // mapstructure:"sampling"`
	Results  Subject      `json:"results"`  // mapstructure:"results"`
}

type SamplingRate int

type Info struct {
	Description string `json:"description,omitempty"` // mapstructure:"description,omitempty"`
	InfoURL     string `json:"info_url,omitempty"`    // mapstructure:"info_url,omitempty"`
}

type Mapping map[Subject][]WeightedMapping

// Mapping for publishes
type WeightedMapping struct {
	Subject Subject `json:"subject"`           // mapstructure:"subject"`
	Weight  uint8   `json:"weight,omitempty"`  // mapstructure:"weight,omitempty"`
	Cluster string  `json:"cluster,omitempty"` // mapstructure:"cluster,omitempty"`
}
