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

// Policy is a rbac template in TKE.
type Policy struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	DisplayName string `json:"displayName" protobuf:"bytes,2,opt,name=displayName"`
	// +optional
	TenantID string `json:"tenantID" protobuf:"bytes,3,opt,name=tenantID"`
	// +optional
	Username string `json:"username" protobuf:"bytes,4,opt,name=username"`
	// +optional
	Description string              `json:"description" protobuf:"bytes,5,opt,name=description"`
	Scope       Scope               `json:"scope" protobuf:"bytes,6,opt,name=scope"`
	Rules       []rbacv1.PolicyRule `json:"rules" protobuf:"bytes,7,rep,name=rules"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyList is the whole list of all rbac templates.
type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of policies
	Items []Policy `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Role is a collection with multiple policies.
type Role struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	DisplayName string `json:"displayName" protobuf:"bytes,2,opt,name=displayName"`
	// +optional
	TenantID string `json:"tenantID" protobuf:"bytes,3,opt,name=tenantID"`
	// +optional
	Username string `json:"username" protobuf:"bytes,4,opt,name=username"`
	// +optional
	Description string   `json:"description" protobuf:"bytes,5,opt,name=description"`
	Scope       Scope    `json:"scope" protobuf:"bytes,6,opt,name=scope"`
	Policies    []string `json:"policies" protobuf:"bytes,7,rep,name=policies"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RoleList is the whole list of policy.
type RoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of rules.
	Items []Role `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MultiClusterRoleBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MultiClusterRoleBindingSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            MultiClusterRoleBindingStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type MultiClusterRoleBindingSpec struct {
	TenantID string   `json:"tenantID" protobuf:"bytes,1,name=tenantID"`
	Username string   `json:"username" protobuf:"bytes,2,name=username"`
	RoleName string   `json:"roleName" protobuf:"bytes,3,name=roleName"`
	Clusters []string `json:"clusters" protobuf:"bytes,4,rep,name=clusters"`
}

type MultiClusterRoleBindingStatus struct {
	// +optional
	Phase BindingPhase `json:"phase" protobuf:"bytes,1,opt,name=phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MultiClusterRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of rules.
	Items []MultiClusterRoleBinding `json:"items" protobuf:"bytes,2,rep,name=items"`
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
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMap holds configuration data for tke to consume.
type ConfigMap struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`

	// BinaryData contains the binary data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// BinaryData can contain byte sequences that are not in the UTF-8 range.
	// The keys stored in BinaryData must not overlap with the ones in
	// the Data field, this is enforced during validation process.
	// +optional
	BinaryData map[string][]byte `json:"binaryData,omitempty" protobuf:"bytes,3,rep,name=binaryData"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapList is a resource containing a list of ConfigMap objects.
type ConfigMapList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of ConfigMaps.
	Items []ConfigMap `json:"items" protobuf:"bytes,2,rep,name=items"`
}
