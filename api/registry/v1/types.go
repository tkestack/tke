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
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Namespace is an image container in registry.
type Namespace struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of namespace in this set.
	// +optional
	Spec NamespaceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status NamespaceStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceList is the whole list of all namespaces which owned by a tenant.
type NamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of namespaces
	Items []Namespace `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// NamespaceSpec is a description of a namespace.
type NamespaceSpec struct {
	Name     string `json:"name" protobuf:"bytes,1,opt,name=name"`
	TenantID string `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	// +optional
	DisplayName string `json:"displayName,omitempty" protobuf:"bytes,3,opt,name=displayName"`
	// +optional
	Visibility Visibility `json:"visibility,omitempty" protobuf:"bytes,4,opt,name=visibility,casttype=Visibility"`
}

// NamespaceStatus represents information about the status of a namespace.
type NamespaceStatus struct {
	// +optional
	Locked    *bool `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`
	RepoCount int32 `json:"repoCount" protobuf:"varint,2,opt,name=repoCount"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Repository is a repo in namespace of registry.
type Repository struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of repository in this set.
	// +optional
	Spec RepositorySpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status RepositoryStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RepositoryList is the whole list of all repositories which owned by a namespace.
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of repositories
	Items []Repository `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type RepositorySpec struct {
	Name          string `json:"name" protobuf:"bytes,1,opt,name=name"`
	TenantID      string `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	NamespaceName string `json:"namespaceName" protobuf:"bytes,3,opt,name=namespaceName"`
	// +optional
	DisplayName string `json:"displayName,omitempty" protobuf:"bytes,4,opt,name=displayName"`
	// +optional
	Visibility Visibility `json:"visibility,omitempty" protobuf:"bytes,5,opt,name=visibility,casttype=Visibility"`
}

type RepositoryStatus struct {
	// +optional
	Locked    *bool           `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`
	PullCount int32           `json:"pullCount" protobuf:"varint,2,opt,name=pullCount"`
	Tags      []RepositoryTag `json:"tags" protobuf:"bytes,3,rep,name=tags"`
}

type RepositoryTag struct {
	Name        string      `json:"name" protobuf:"bytes,1,opt,name=name"`
	Digest      string      `json:"digest" protobuf:"bytes,2,opt,name=digest"`
	TimeCreated metav1.Time `json:"timeCreated,omitempty" protobuf:"bytes,3,opt,name=timeCreated"`
}

// Visibility defines the visible properties of the repo or namespace.
type Visibility string

const (
	// VisibilityPublic indicates the namespace or repo is public.
	VisibilityPublic Visibility = "Public"
	// VisibilityPrivate indicates the namespace or repo is private.
	VisibilityPrivate Visibility = "Private"
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
