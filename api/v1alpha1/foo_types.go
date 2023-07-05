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
	"github.com/redhat-appstudio/operator-toolkit/conditions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FooSpec defines the desired state of Foo
type FooSpec struct {
	// DesiredReplicas is the number of Bar replicas that should exist at any given moment
	DesiredReplicas int `json:"desiredReplicas"`
}

// FooStatus defines the observed state of Foo
type FooStatus struct {
	// Conditions represent the latest available observations for the Foo resource
	// +optional
	Conditions []metav1.Condition `json:"conditions"`

	// Replicas is a slice containing the list of replica names for this resource
	Replicas []string `json:"replicas,omitempty"`
}

// MarkHealthy marks the Foo resource as healthy using the reason passed as a parameter
func (f *Foo) MarkHealthy(reason conditions.ConditionReason) {
	conditions.SetCondition(&f.Status.Conditions, healthConditionType, metav1.ConditionTrue, reason)
}

// MarkUnhealthy marks the Foo resource as unhealthy
func (f *Foo) MarkUnhealthy() {
	conditions.SetCondition(&f.Status.Conditions, healthConditionType, metav1.ConditionFalse, NotEnoughReplicasReason)
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Health",type=string,JSONPath=`.status.conditions[?(@.type=="Health")].reason`

// Foo is the Schema for the foos API
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooSpec   `json:"spec,omitempty"`
	Status FooStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FooList contains a list of Foo
type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Foo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Foo{}, &FooList{})
}
