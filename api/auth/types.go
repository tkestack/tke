/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package auth

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentity is an object that contains the metadata about identify used to login
// to TKE.
type LocalIdentity struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	// Spec defines the desired identities of identity in this set.
	Spec   LocalIdentitySpec
	Status LocalIdentityStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentityList is the whole list of all identities.
type LocalIdentityList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of identities.
	Items []LocalIdentity
}

// LocalIdentitySpec is a description of an identity.
type LocalIdentitySpec struct {
	Username         string
	DisplayName      string
	Email            string
	PhoneNumber      string
	HashedPassword   string
	OriginalPassword string
	TenantID         string
	Groups           []string
	Extra            map[string]string
}

// LocalIdentityStatus is a description of an identity status.
type LocalIdentityStatus struct {
	Locked bool

	// The last time the local identity was updated.
	// +optional
	LastUpdateTime metav1.Time
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKey contains expiration time used to apply the api key.
type APIKey struct {
	metav1.TypeMeta

	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of APIkey in this set.
	Spec   APIKeySpec
	Status APIKeyStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentityList is the whole list of all identities.
type APIKeyList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of api keys.
	Items []APIKey
}

// APIKeySpec is a description of an apiKey.
type APIKeySpec struct {
	// APIkey is the jwt token used to authenticate user, and contains user info and sign.
	APIkey string `json:"apiKey,omitempty"`

	TenantID string `json:"tenantID,omitempty"`

	// Creator
	// +optional
	Username string `json:"username,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description"`

	// IssueAt is the created time for api key
	IssueAt metav1.Time `json:"issue_at,omitempty"`

	// ExpireAt is the expire time for api key
	ExpireAt metav1.Time `json:"expire_at,omitempty"`
}

// APIKeyStatus is a description of an api key status.
type APIKeyStatus struct {
	// Disabled represents whether the apikey has been disabled.
	Disabled bool `json:"disabled"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReq contains expiration time used to apply the api key.
type APIKeyReq struct {
	metav1.TypeMeta

	//Exipre is required, holds the duration of the api key become invalid. By default, 168h(= seven days)
	Expire metav1.Duration `json:"expire,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReqPassword contains userinfo and expiration time used to apply the api key.
type APIKeyReqPassword struct {
	metav1.TypeMeta

	// TenantID for user
	TenantID string `json:"tenantID,omitempty"`

	// Username
	Username string `json:"username,omitempty"`

	// Password (encoded by base64)
	Password string `json:"password,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description"`

	// Expire holds the duration of the api key become invalid. By default, 168h(= seven days)
	Expire metav1.Duration `json:"expire,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APISigningKey hold encryption and signing key for api key.
type APISigningKey struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	SigningKey    []byte
	SigningKeyPub []byte
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APISigningKeyList is the whole list of all signing keys.
type APISigningKeyList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []APISigningKey
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMap holds configuration data for tke to consume.
type ConfigMap struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string

	// BinaryData contains the binary data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// BinaryData can contain byte sequences that are not in the UTF-8 range.
	// The keys stored in BinaryData must not overlap with the ones in
	// the Data field, this is enforced during validation process.
	// +optional
	BinaryData map[string][]byte
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapList is a resource containing a list of ConfigMap objects.
type ConfigMapList struct {
	metav1.TypeMeta

	// +optional
	metav1.ListMeta

	// Items is the list of ConfigMaps.
	Items []ConfigMap
}
