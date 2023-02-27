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

package accountSigningKey

// +kubebuilder:object:generate=true
// Specifies how the account signing key is created.
type AccountSigningKeyConfig struct {
	// Import specifies an existing Kubernetes secret to import as the account signing key.
	// +kubebuilder:validation:Optional
	Import ImportAccountSigningKey `json:"import,omitempty"`
}

// ImportAccountSigningKey will import a signing key from a Kubernetes secret.
type ImportAccountSigningKey struct {
	// References a Kubernetes secret that contains the account signing key.
	// +kubebuilder:validation:Required
	SecretRef *LocalObjectReference `json:"secretRef"`
}

// LocalObjectReference references a Kubernetes object
type LocalObjectReference struct {
	// The name of the Kubernetes secret that contains the account signing key.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}
