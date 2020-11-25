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

	// +kubebuilder:validation:Required

	Term Term `json:"term,omitempty"`

	// +kubebuilder:validation:Required

	Encryption Encryption `json:"encryption,omitempty"`
}

type Term struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string

	ExpiresAfter string `json:"expiresAfter,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time

	ExpiresDate string `json:"expiresDate,omitempty"`
}

type Encryption struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Format:=string

	// Use the github public key to encrypt the authentication file.
	Github Github `json:"github,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Format:=string

	// Use this public key to encrypt the authentication file.
	// key format is OpenSSH public key format.
	PublicKey string `json:"publicKey,omitempty"`
}

type Github struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format:=string

	// Github UserID
	User string `json:"user,omitempty"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional

	// KeyNumber is line number of https://github.com/<user>.keys
	// Default value: 1
	KeyNumber *int32 `json:"keyNumber,omitempty"`
}

// CinderellaStatus defines the observed state of Cinderella
type CinderellaStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time

	// ExpiredAt is expired at binding account
	// ExpiredAt is RFC 3339 date and time at which this resource will be deleted.
	ExpiredAt metav1.Time `json:"expiredAt,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=cin;prole
// +kubebuilder:printcolumn:name="Expired",type="date",JSONPath=".status.expiredAt"

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
