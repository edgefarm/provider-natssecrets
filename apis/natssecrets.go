/*
Copyright 2020 The Crossplane Authors.

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

// Package apis contains Kubernetes API for the NatsSecrets provider.
package apis

import (
	"k8s.io/apimachinery/pkg/runtime"

	accountv1alpha1 "github.com/edgefarm/provider-natssecrets/apis/account/v1alpha1"
	accountsigningkeyv1alpha1 "github.com/edgefarm/provider-natssecrets/apis/accountSigningKey/v1alpha1"
	operatorv1alpha1 "github.com/edgefarm/provider-natssecrets/apis/operator/v1alpha1"
	operatorsigningkeyv1alpha1 "github.com/edgefarm/provider-natssecrets/apis/operatorSigningKey/v1alpha1"
	userv1alpha1 "github.com/edgefarm/provider-natssecrets/apis/user/v1alpha1"
	natssecretsv1alpha1 "github.com/edgefarm/provider-natssecrets/apis/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes,
		natssecretsv1alpha1.SchemeBuilder.AddToScheme,
		operatorv1alpha1.SchemeBuilder.AddToScheme,
		accountv1alpha1.SchemeBuilder.AddToScheme,
		userv1alpha1.SchemeBuilder.AddToScheme,
		accountsigningkeyv1alpha1.SchemeBuilder.AddToScheme,
		operatorsigningkeyv1alpha1.SchemeBuilder.AddToScheme,
	)
}

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
