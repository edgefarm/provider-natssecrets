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

	config "github.com/edgefarm/provider-natssecrets/apis/accountSigningKey/v1alpha1/config"
)

// AccountSigningKeyParameters are the configurable fields of a account.
type AccountSigningKeyParameters struct {
	// The name of the operator
	Operator string `json:"operator"`
	// The name of the account
	Account string `json:"account"`
	// Specifies how the account signing key is generated
	Config config.AccountSigningKeyConfig `json:"config,omitempty"`
}

// AccountSigningKeyObservation are the observable fields of a AccountSigningKey.
type AccountSigningKeyObservation struct {
	Operator string `json:"operator,omitempty"`
	Account  string `json:"account,omitempty"`
	NKey     string `json:"nkey,omitempty"`
	NKeyPath string `json:"nkeyPath,omitempty"`
}

// A AccountSigningKeySpec defines the desired state of a AccountSigningKey.
type AccountSigningKeySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       AccountSigningKeyParameters `json:"forProvider"`
}

// A AccountSigningKeyStatus represents the observed state of a AccountSigningKey.
type AccountSigningKeyStatus struct {
	xpv1.ResourceStatus `json:",inline"`

	AtProvider AccountSigningKeyObservation `json:"atProvider,omitempty"`

	// Status of this instance.
	Status string `json:"status,omitempty"`
}

// An AccountSigningKey is an API type for account signing keys.
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="OPERATOR",type="string",priority=1,JSONPath=".status.atProvider.operator"
// +kubebuilder:printcolumn:name="ACCOUNT",type="string",priority=1,JSONPath=".status.atProvider.account"
// +kubebuilder:printcolumn:name="NKEY",type="string",priority=1,JSONPath=".status.atProvider.nkey"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,natssecrets}
type AccountSigningKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccountSigningKeySpec   `json:"spec"`
	Status AccountSigningKeyStatus `json:"status,omitempty"`
}

// AccountSigningKeyList contains a list of AccountSigningKey
// +kubebuilder:object:root=true
type AccountSigningKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AccountSigningKey `json:"items"`
}

// AccountSigningKey type metadata.
var (
	AccountSigningKeyKind             = reflect.TypeOf(AccountSigningKey{}).Name()
	AccountSigningKeyGroupKind        = schema.GroupKind{Group: Group, Kind: AccountSigningKeyKind}.String()
	AccountSigningKeyKindAPIVersion   = AccountSigningKeyKind + "." + SchemeGroupVersion.String()
	AccountSigningKeyGroupVersionKind = SchemeGroupVersion.WithKind(AccountSigningKeyKind)
)

func init() {
	SchemeBuilder.Register(&AccountSigningKey{}, &AccountSigningKeyList{})
}
