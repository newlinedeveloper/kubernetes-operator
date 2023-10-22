/*
Copyright 2023 Veera.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CloudnloudSpec defines the desired state of Cloudnloud
type CloudnloudSpec struct {
	Start       int              `json:"start"`
	End         int              `json:"end"`
	Replicas    int              `json:"replicas"`
	Deployments []NamespacedName `json:"deployments"`
}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// CloudnloudStatus defines the observed state of Cloudnloud
type CloudnloudStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cloudnloud is the Schema for the cloudnlouds API
type Cloudnloud struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudnloudSpec   `json:"spec,omitempty"`
	Status CloudnloudStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CloudnloudList contains a list of Cloudnloud
type CloudnloudList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cloudnloud `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cloudnloud{}, &CloudnloudList{})
}
