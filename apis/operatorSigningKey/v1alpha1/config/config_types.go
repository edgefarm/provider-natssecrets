/*
Copyright 2023 The Crossplane Authors.

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

package operatorSigningKey

// +kubebuilder:object:generate=true
// Specifies how the operator signing key is created.
type OperatorSigningKeyConfig struct {
	// Specifies an existing Kubernetes secret to import as the operator signing key.
	// +kubebuilder:validation:Optional
	Import ImportOperatorSigningKey `json:"import,omitempty"`
}

// ImportOperatorSigningKey will import a signing key from a Kubernetes secret.
type ImportOperatorSigningKey struct {
	// References a Kubernetes secret that contains the operator signing key.
	// +kubebuilder:validation:Required
	SecretRef *LocalObjectReference `json:"secretRef,omitempty"`
}

// LocalObjectReference references a Kubernetes object
type LocalObjectReference struct {
	// The name of the Kubernetes secret that contains the operator signing key.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// The namespace of the Kubernetes secret that contains the operator signing key.
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
	// The key of the Kubernetes secret that contains the operator signing key's seed.
	// +kubebuilder:validation:Required
	Key string `json:"key"`
}
