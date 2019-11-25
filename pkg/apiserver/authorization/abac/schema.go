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

// Portions Copyright 2014 The Kubernetes Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package abac

import (
	"fmt"
)

// policy contains a single ABAC policy rule
type policy struct {

	// Spec describes the policy rule
	Spec policySpec
}

// policySpec contains the attributes for a policy rule
type policySpec struct {

	// User is the username this rule applies to.
	// Either user or group is required to match the request.
	// "*" matches all users and support regex match.
	User string

	// Group is the group this rule applies to.
	// Either user or group is required to match the request.
	// "*" matches all groups and support regex match.
	Group string

	// Readonly matches readonly requests when true, and all requests when false
	Readonly bool

	// Verb is the username this rule applies to.
	// support regex match.
	Verb string

	// APIGroup is the name of an API group. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all API groups
	APIGroup string

	// Resource is the name of a resource. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all resources and support regex match,
	Resource string

	// Namespace is the name of a namespace. APIGroup, Resource, and Namespace are required to match resource requests.
	// "*" matches all namespaces (including unnamespaced requests)
	Namespace string

	// NonResourcePath matches non-resource request paths.
	// "*" matches all paths
	// "/foo/*" matches all subpaths of foo
	NonResourcePath string
}

type policyLoadError struct {
	path string
	line int
	data []byte
	err  error
}

// Error describes parse policy failed info.
func (p policyLoadError) Error() string {
	if p.line >= 0 {
		return fmt.Sprintf("error reading policy file %s, line %d: %s: %v", p.path, p.line, string(p.data), p.err)
	}
	return fmt.Sprintf("error reading policy file %s: %v", p.path, p.err)
}
