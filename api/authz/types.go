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

package authz

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Scope string

const (
	PlatformScope     Scope = "Platform"
	MultiClusterScope Scope = "MultiCluster"
	BusinessScope     Scope = "Business"
)

// +genclient
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Policy struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	DisplayName string

	// +optional
	TenantID string

	// Username is Creator
	// +optional
	Username string

	// +optional
	Description string

	Scope Scope

	Rules []rbacv1.PolicyRule
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyList is the whole list of all policies.
type PolicyList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	// List of policies
	Items []Policy
}

// +genclient
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Role is a collection with multiple policies.
type Role struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	DisplayName string

	// +optional
	TenantID string

	// Username is Creator
	// +optional
	Username string

	// +optional
	Description string

	Scope Scope

	// policyNamespace/policyName
	Policies []string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RoleList is the whole list of policy.
type RoleList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of rules.
	Items []Role
}

// +genclient
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MultiClusterRoleBinding struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   MultiClusterRoleBindingSpec
	Status MultiClusterRoleBindingStatus
}

type MultiClusterRoleBindingSpec struct {
	Username string
	// roleNamespace/roleName
	RoleName string
	Clusters []string
}

type MultiClusterRoleBindingStatus struct {
	// +optional
	Phase BindingPhase
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MultiClusterRoleBindingList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of rules.
	Items []MultiClusterRoleBinding
}

type BindingPhase string

const (
	BindingActive      BindingPhase = "Active"
	BindingTerminating BindingPhase = "Terminating"
)

type FinalizerName string

const (
	PolicyFinalize                  FinalizerName = "policy"
	RoleFinalize                    FinalizerName = "role"
	MultiClusterRoleBindingFinalize FinalizerName = "rolebinding"
)

// +genclient
// +genclient:nonNamespaced
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
