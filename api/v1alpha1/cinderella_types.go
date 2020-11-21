/*


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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CinderellaSpec defines the desired state of Cinderella
type CinderellaSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Cinderella. Edit Cinderella_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// CinderellaStatus defines the observed state of Cinderella
type CinderellaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Cinderella is the Schema for the cinderellas API
type Cinderella struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CinderellaSpec   `json:"spec,omitempty"`
	Status CinderellaStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CinderellaList contains a list of Cinderella
type CinderellaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cinderella `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cinderella{}, &CinderellaList{})
}
