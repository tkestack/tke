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

const (
	// NameTag is used for route path
	NameTag = "name"

	// IDTag is used for route path
	IDTag = "id"

	// owner is used for route param
	UserTag = "username"

	// KeywordTag is used for route param
	KeywordTag = "keyword"

	// TypeTag is used for route param
	TypeTag = "type"

	// ScopeTag is used for route param
	ScopeTag = "scope"

	// ListScopeAll is used for listing all policies or roles
	ListScopeAll = "all"

	// ListScopeLocal is used for listing only policies or roles created by normal user
	ListScopeLocal = "local"

	// ListScopeSystem is used for listing only policies or roles created by system default
	ListScopeSystem = "system"

	// RoleIDPrefix to identify role in casbin
	RoleIDPrefix = "role-"

	// PolicyIDPrefix to identify policy in casbin
	PolicyIDPrefix = "policy-"

	// UserPrefix to identify user in casbin
	UserPrefix = "user-"

	// IssuerName is the name of issuer location.
	IssuerName = "oidc"
)
