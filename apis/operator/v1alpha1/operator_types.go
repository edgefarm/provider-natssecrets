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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	operatorv1 "github.com/edgefarm/vault-plugin-secrets-nats/pkg/claims/operator/v1alpha1"
)

// OperatorParameters are the configurable fields of a Operator.
type OperatorParameters struct {
	Operator      string                    `json:"operator"`
	SystemAccount string                    `json:"system_account,omitempty"`
	SigningKeys   []string                  `json:"signing_keys,omitempty"`
	Claims        operatorv1.OperatorClaims `json:"operator_claims,omitempty"`
}

// OperatorObservation are the observable fields of a Operator.
type OperatorObservation struct {
	Operator string `json:"operator,omitempty"`
	Issue    string `json:"issue,omitempty"`
	NKey     string `json:"nkey,omitempty"`
	JWT      string `json:"jwt,omitempty"`
}

// A OperatorSpec defines the desired state of a Operator.
type OperatorSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       OperatorParameters `json:"forProvider"`
}

// A OperatorStatus represents the observed state of a Operator.
type OperatorStatus struct {
	xpv1.ResourceStatus `json:",inline"`

	AtProvider OperatorObservation `json:"atProvider,omitempty"`

	// Status of this instance.
	Status string `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// A Operator is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,natssecrets}
type Operator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OperatorSpec   `json:"spec"`
	Status OperatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OperatorList contains a list of Operator
type OperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Operator `json:"items"`
}

// Operator type metadata.
var (
	OperatorKind             = reflect.TypeOf(Operator{}).Name()
	OperatorGroupKind        = schema.GroupKind{Group: Group, Kind: OperatorKind}.String()
	OperatorKindAPIVersion   = OperatorKind + "." + SchemeGroupVersion.String()
	OperatorGroupVersionKind = SchemeGroupVersion.WithKind(OperatorKind)
)

func init() {
	SchemeBuilder.Register(&Operator{}, &OperatorList{})
}
