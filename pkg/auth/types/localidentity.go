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

package types

import (
	"time"
)

// LocalIdentity is an object that contains the metadata about identify used to login
// to TKE.
//+k8s:openapi-gen=true
type LocalIdentity struct {
	Name     string    `json:"name,omitempty"`
	UID      string    `json:"uid,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
	UpdateAt time.Time `json:"updateAt,omitempty"`

	// Spec defines the desired identities of identity in this set.
	Spec   *LocalIdentitySpec
	Status *LocalIdentityStatus
}

// LocalIdentityList is the whole list of all identities.
//+k8s:openapi-gen=true
type LocalIdentityList struct {
	// List of identities.
	Items []*LocalIdentity `json:"items,omitempty"`
}

// LocalIdentitySpec is a description of an identity.
//+k8s:openapi-gen=true
type LocalIdentitySpec struct {
	HashedPassword   string            `json:"hashedPassword,omitempty"`
	OriginalPassword string            `json:"originalPassword,omitempty"`
	Groups           []string          `json:"groups,omitempty"`
	TenantID         string            `json:"tenantID,omitempty"`
	Extra            map[string]string `json:"extra,omitempty"`
}

// LocalIdentityStatus is a description of an identity status.
//+k8s:openapi-gen=true
type LocalIdentityStatus struct {
	Locked bool `json:"locked,omitempty"`
}
