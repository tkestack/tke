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

import "time"

// Category defines a category of actions for policy.
//+k8s:openapi-gen=true
type Category struct {
	// Name identifies policy category
	Name string `json:"name,omitempty"`

	// DisplayName used to display category name
	DisplayName string `json:"displayName,omitempty"`

	TenantID    string `json:"tenantID,omitempty"`
	Description string `json:"description,omitempty"`
	// Actions represents a series of actions work on the policy category
	Actions map[string]Action `json:"actions"`

	CreateAt time.Time `json:"createAt,omitempty"`
	UpdateAt time.Time `json:"updateAt,omitempty"`
}

// Action defines a action verb for authorization.
//+k8s:openapi-gen=true
type Action struct {
	// Name represents user access review request verb.
	Name string `json:"name,omitempty"`
	// Description describes the action.
	Description string `json:"description,omitempty"`
}

// CategoryList is the whole list of policy Category.
//+k8s:openapi-gen=true
type CategoryList struct {
	// List of category.
	Items []*Category `json:"items,omitempty"`
}
