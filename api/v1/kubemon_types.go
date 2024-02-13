/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KubeMonSpec defines the desired state of KubeMon
type KubeMonSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Species string `json:"species"`
	// +kubebuilder:validation:Minimum:1
	// +kubebuilder:validation:Minimum:99
	Level *int32 `json:"level"`
	Owner string `json:"owner,omitempty"`
}

// KubeMonStatus defines the observed state of KubeMon
type KubeMonStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	HP *int32 `json:"hp,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KubeMon is the Schema for the kubemons API
type KubeMon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeMonSpec   `json:"spec,omitempty"`
	Status KubeMonStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KubeMonList contains a list of KubeMon
type KubeMonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeMon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeMon{}, &KubeMonList{})
}
