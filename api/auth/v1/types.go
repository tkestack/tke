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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentity is an object that contains the metadata about identify used to
// login to TKE.
type LocalIdentity struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of identity in this set.
	// +optional
	Spec LocalIdentitySpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status LocalIdentityStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentityList is the whole list of all identities.
type LocalIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of identities.
	Items []LocalIdentity `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// LocalIdentitySpec is a description of an identity.
type LocalIdentitySpec struct {
	UserName         string `json:"userName" protobuf:"bytes,7,opt,name=name"`
	DisplayName      string `json:"displayName" protobuf:"bytes,8,opt,name=displayName"`
	Email            string `json:"email" protobuf:"bytes,9,opt,name=email"`
	PhoneNumber      string `json:"phoneNumber" protobuf:"bytes,10,opt,name=phone"`
	HashedPassword   string `json:"hashedPassword,omitempty" protobuf:"bytes,4,opt,name=hashedPassword"`
	OriginalPassword string `json:"originalPassword,omitempty" protobuf:"bytes,5,opt,name=originalPassword"`
	// +optional
	TenantID string `json:"tenantID,omitempty" protobuf:"bytes,2,opt,name=tenantID"`

	Groups []string `json:"groups,omitempty" protobuf:"bytes,6,rep,name=groups"`
	// +optional
	Extra map[string]string `json:"extra,omitempty" protobuf:"bytes,3,rep,name=extra"`
}

// LocalIdentityStatus is a description of an identity status.
type LocalIdentityStatus struct {
	// +optional
	Locked bool `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`

	// The last time the local identity was updated.
	// +optional
	LastUpdateTime metav1.Time `protobuf:"bytes,2,opt,name=lastUpdateTime"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKey contains expiration time used to apply the api key.
type APIKey struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=objectMeta"`

	// Spec defines the desired identities of APIkey in this set.
	Spec   APIKeySpec   `protobuf:"bytes,2,opt,name=spec"`
	Status APIKeyStatus `protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentityList is the whole list of all identities.
type APIKeyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=listMeta"`
	// List of api keys.
	Items []APIKey `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}

// APIKeySpec is a description of an apiKey.
type APIKeySpec struct {
	// APIkey is the jwt token used to authenticate user, and contains user info and sign.
	APIkey string `json:"apiKey,omitempty" protobuf:"bytes,1,opt,name=apiKey"`

	TenantID string `json:"tenantID,omitempty" protobuf:"bytes,5,opt,name=tenantID"`

	// Description describes api keys usage.
	// +optional
	Description string `json:"description,omitempty" protobuf:"bytes,2,opt,name=description"`

	// IssueAt is the created time for api key
	IssueAt metav1.Time `json:"issue_at,omitempty" protobuf:"bytes,3,opt,name=issue_at,json=issueAt"`

	// ExpireAt is the expire time for api key
	ExpireAt metav1.Time `json:"expire_at,omitempty" protobuf:"bytes,4,opt,name=expire_at,json=expireAt"`
}

// APIKeyStatus is a description of an api key status.
type APIKeyStatus struct {
	// Disabled represents whether the apikey has been disabled.
	// +optional
	Disabled *bool `json:"disabled,omitempty" protobuf:"varint,1,opt,name=disabled"`

	// Deleted represents whether the apikey has been deleted.
	// +optional
	Deleted *bool `json:"deleted,omitempty" protobuf:"varint,2,opt,name=deleted"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReq contains expiration time used to apply the api key.
type APIKeyReq struct {
	metav1.TypeMeta `json:",inline"`

	//Exipre is required, holds the duration of the api key become invalid. By default, 168h(= seven days)
	// +optional
	Expire metav1.Duration `json:"expire,omitempty" protobuf:"bytes,2,opt,name=expire"`

	// Description describes api keys usage.
	Description string `json:"description" protobuf:"bytes,3,opt,name=description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReqPassword contains userinfo and expiration time used to apply the api key.
type APIKeyReqPassword struct {
	metav1.TypeMeta `json:",inline"`

	// TenantID for user
	TenantID string `json:"tenantID,omitempty" protobuf:"bytes,2,opt,name=tenantID"`

	// UserName
	UserName string `json:"username,omitempty" protobuf:"bytes,3,opt,name=username"`

	// Password (encoded by base64)
	Password string `json:"password,omitempty" protobuf:"bytes,4,opt,name=password"`

	// Description describes api keys usage.
	// +optional
	Description string `json:"description" protobuf:"bytes,5,opt,name=description"`

	// Expire holds the duration of the api key become invalid. By default, 168h(= seven days)
	// +optional
	Expire metav1.Duration `json:"expire,omitempty" protobuf:"bytes,6,opt,name=expire"`
}
