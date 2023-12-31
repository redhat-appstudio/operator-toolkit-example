/*
Copyright 2023.

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

// BarSpec defines the desired state of Bar
type BarSpec struct {
	// Foo is the name of the Foo resource associated with this resource
	// +optional
	Foo string `json:"foo,omitempty"`
}

//+kubebuilder:object:root=true

// Bar is the Schema for the bars API
type Bar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BarSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// BarList contains a list of Bar
type BarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bar `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bar{}, &BarList{})
}
