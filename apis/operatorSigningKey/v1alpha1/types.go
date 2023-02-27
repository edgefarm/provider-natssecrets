/*
Copyright 2022 The Crossplane Authors.

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

package v1alpha1

import (
	"reflect"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	config "github.com/edgefarm/provider-natssecrets/apis/operatorSigningKey/v1alpha1/config"
)

// OperatorSigningKeyParameters are the configurable fields of a Operator.
type OperatorSigningKeyParameters struct {
	// The name of the operator
	Operator string `json:"operator"`
	// Specifies how the operator signing key is generated
	Config config.OperatorSigningKeyConfig `json:"config,omitempty"`
}

// OperatorSigningKeyObservation are the observable fields of a OperatorSigningKey.
type OperatorSigningKeyObservation struct {
	Operator string `json:"operator,omitempty"`
	NKey     string `json:"nkey,omitempty"`
}

// A OperatorSigningKeySpec defines the desired state of a OperatorSigningKey.
type OperatorSigningKeySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       OperatorSigningKeyParameters `json:"forProvider"`
}

// A OperatorSigningKeyStatus represents the observed state of a OperatorSigningKey.
type OperatorSigningKeyStatus struct {
	xpv1.ResourceStatus `json:",inline"`

	AtProvider OperatorSigningKeyObservation `json:"atProvider,omitempty"`

	// Status of this instance.
	Status string `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// An OperatorSigningKey is an API type for operator signing keys.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="OPERATOR",type="string",priority=1,JSONPath=".status.atProvider.operator"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,natssecrets}
type OperatorSigningKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OperatorSigningKeySpec   `json:"spec"`
	Status OperatorSigningKeyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OperatorSigningKeyList contains a list of OperatorSigningKey
type OperatorSigningKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OperatorSigningKey `json:"items"`
}

// OperatorSigningKey type metadata.
var (
	OperatorSigningKeyKind             = reflect.TypeOf(OperatorSigningKey{}).Name()
	OperatorSigningKeyGroupKind        = schema.GroupKind{Group: Group, Kind: OperatorSigningKeyKind}.String()
	OperatorSigningKeyKindAPIVersion   = OperatorSigningKeyKind + "." + SchemeGroupVersion.String()
	OperatorSigningKeyGroupVersionKind = SchemeGroupVersion.WithKind(OperatorSigningKeyKind)
)

func init() {
	SchemeBuilder.Register(&OperatorSigningKey{}, &OperatorSigningKeyList{})
}
