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

// +kubebuilder:validation:Enum=ClusterRole;Role
type RoleKind string

const (
	NamespaceRole RoleKind = "Role"
	ClusterRole   RoleKind = "ClusterRole"
)

type Role struct {
	// +kubebuilder:validation:Required

	// Types of roles to bind
	// Valid values are:
	// - "Role":
	// - "ClusterRole":
	Kind RoleKind `json:"kind,omitempty"`

	// +kubebuilder:validation:Required

	// `ClusterRole` or `Role` Name
	Name string `json:"name,omitempty"`
}

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CinderellaSpec defines the desired state of Cinderella
type CinderellaSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1

	// roles for temporary user
	Roles []Role `json:"roles,omitempty"`

	// +kubebuilder:validation:Required

	// expiration term of temporary user
	Term Term `json:"term,omitempty"`

	// +kubebuilder:validation:Required

	// Encryption by public key for passing files to temporary user
	Encryption Encryption `json:"encryption,omitempty"`
}

// Term is expiration of temporary user
// This is expressed as a date time, or a deadline such as 60min later.
type Term struct {
	// +kubebuilder:validation:Optional

	// Temporary user is will be invalidated after specified value of ExpiresAfter
	// The unit is minutes.
	ExpiresAfter *int32 `json:"expiresAfter,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time

	// RFC 3339
	// e.g. "2020-12-01T00:00:00+09:00"
	ExpiresDate string `json:"expiresDate,omitempty"`
}

// Note: Encryptionって名前は微妙な気がしてる...

// Encryption is specifies where to get the gpg public key.
type Encryption struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Format:=string

	// Use the github registered public key by user to encrypt the authentication file.
	Github Github `json:"github,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Format:=string

	// Use this public key to encrypt the authentication file.
	// key format must be OpenSSH public key format.
	PublicKey string `json:"publicKey,omitempty"`
}

type Github struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format:=string

	// github user ID, fetch from https://github.com/<user>.gpg
	User string `json:"user,omitempty"`

	// TODO: default 動かない
	// _+kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional

	// KeyNumber is line number of https://github.com/<user>.gpg
	// Default value: 1
	KeyNumber *int32 `json:"keyNumber,omitempty"`
}

// CinderellaStatus defines the observed state of Cinderella
type CinderellaStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Type=boolean

	// Expired this resource
	Expired *bool `json:"expired,omitempty"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time

	// ExpiredAt is expired at binding account
	// ExpiredAt is RFC 3339 format date and time at which this resource will be deleted.
	ExpiredAt metav1.Time `json:"expiredAt,omitempty"`
}

// NOTE:
// https://book.kubebuilder.io/reference/markers/crd.html
// https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#additional-printer-columns

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=cin,scope=Cluster
// +kubebuilder:printcolumn:name="Expired-At",type="string",JSONPath=".status.expiredAt",format="date-time"
// +kubebuilder:printcolumn:name="Expired",type="boolean",JSONPath=".status.expired"

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
