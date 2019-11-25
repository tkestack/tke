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

// PolicyType represents policy or role created by user or system default
//+k8s:openapi-gen=true
type PolicyType int

const (
	// UserDefine stats policy or role created by normal user
	UserDefine PolicyType = iota

	// PreDefine stats policy or role created by system default
	PreDefine
)

// Policy defines a data structure containing a authorization strategy.
//+k8s:openapi-gen=true
type Policy struct {
	Name        string    `json:"name"`
	ID          string    `json:"id"`
	TenantID    string    `json:"tenantID"`
	Service     string    `json:"service"`
	Statement   Statement `json:"statement"`
	UserName    string    `json:"userName"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"createAt"`
	UpdateAt    time.Time `json:"updateAt"`

	//Type defines policy is created by system(1), or created by user(0)
	Type PolicyType `json:"type"`
}

// PolicyCreate defines the policy create request.
//+k8s:openapi-gen=true
type PolicyCreate struct {
	Name      string    `json:"name"`
	TenantID  string    `json:"tenantID"`
	Service   string    `json:"service"`
	Statement Statement `json:"statement"`
	// UserName  claims users attached to the policy created and split by ','. e.g: user1,user2.
	UserName    string `json:"userName"`
	Description string `json:"description"`
}

// Statement defines a series of action on resource can be done or not.
//+k8s:openapi-gen=true
type Statement struct {
	Actions  []string `json:"action"`
	Resource string   `json:"resource"`
	// Effect indicates action on the resource is allowed or not, can be "allow" or "deny"
	Effect string `json:"effect"`
}

// Permission defines a series of action on resource can be done or not.
//+k8s:openapi-gen=true
type Permission struct {
	AllowPerms map[string][]string `json:"allow"`
	DenyPerms  map[string][]string `json:"deny"`
}

// PolicyList is the whole list of policy.
//+k8s:openapi-gen=true
type PolicyList struct {
	// List of policies.
	Items []*Policy `json:"items"`
}

// PolicyMeta contains metadata of Policy used for in roles.
//+k8s:openapi-gen=true
type PolicyMeta struct {
	Name        string     `json:"name"`
	ID          string     `json:"id"`
	TenantID    string     `json:"tenantID"`
	Service     string     `json:"service"`
	Type        PolicyType `json:"type"`
	Description string     `json:"description"`
}

// PolicyOption is option for listing polices.
//+k8s:openapi-gen=true
type PolicyOption struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
	Name     string `json:"name"`
	Keyword  string `json:"keyword"`
	TenantID string `json:"tenantID"`

	// PolicyStorage list scope: local, system, all
	Scope string `json:"scope"`
}

// AttachInfo contains info to attach/detach users to/from policy or role.
//+k8s:openapi-gen=true
type AttachInfo struct {
	// role or policy id bond
	ID string `json:"id"`
	// name of users
	UserNames []string `json:"userNames"`
	// id of policies
	PolicyIDs []string `json:"policyIDs"`
	// id of tenant
	TenantID string `json:"tenantID"`
}
