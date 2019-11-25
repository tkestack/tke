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

package business

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Project is a project in TKE.
type Project struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec ProjectSpec
	// +optional
	Status ProjectStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProjectList is the whole list of all projects which owned by a tenant.
type ProjectList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of projects
	Items []Project
}

// ProjectSpec is a description of a project.
type ProjectSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers  []FinalizerName
	TenantID    string
	DisplayName string

	// Members represents the user list of project.
	Members []string
	// +optional
	ParentProjectName string
	// Clusters represents clusters that can be used and the resource limits of each cluster.
	// +optional
	Clusters ClusterHard
}

// ProjectStatus represents information about the status of a project.
type ProjectStatus struct {
	// +optional
	Locked *bool
	// +optional
	Phase ProjectPhase
	// Clusters represents clusters that have been used and the resource usage of each cluster.
	// +optional
	Clusters ClusterUsed
	// +optional
	CalculatedChildProjects []string
	// +optional
	CalculatedNamespaces []string
}

// ProjectPhase defines the phase of project constructor.
type ProjectPhase string

const (
	// ProjectActive indicates the project is active.
	ProjectActive ProjectPhase = "Active"
	// ProjectTerminating means the project is undergoing graceful termination.
	ProjectTerminating ProjectPhase = "Terminating"
)

// FinalizerName is the name identifying a finalizer during project lifecycle.
type FinalizerName string

const (
	// ProjectFinalize is an internal finalizer values to Project.
	ProjectFinalize FinalizerName = "project"
	// NamespaceFinalize is an internal finalizer values to Namespace.
	NamespaceFinalize FinalizerName = "namespace"
)

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[string]resource.Quantity

// HardQuantity is a straightforward wrapper of ResourceList.
type HardQuantity struct {
	Hard ResourceList
}

// UsedQuantity is a straightforward wrapper of ResourceList.
type UsedQuantity struct {
	Used ResourceList
}

// ClusterHard is a set of (cluster name, ResourceQuantity) pairs.
type ClusterHard map[string]HardQuantity

// ClusterUsed is a set of (cluster name, ResourceQuantity) pairs.
type ClusterUsed map[string]UsedQuantity

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Namespace is a namespace in cluster.
type Namespace struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of namespaces in this set.
	// +optional
	Spec NamespaceSpec
	// +optional
	Status NamespaceStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceList is the whole list of all namespaces which owned by a tenant.
type NamespaceList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of namespaces
	Items []Namespace
}

// NamespaceSpec represents a namespace in cluster of a project.
type NamespaceSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers  []FinalizerName
	TenantID    string
	ClusterName string
	Namespace   string
	// Hard represents the total resources of a namespace.
	// +optional
	Hard ResourceList
}

// NamespaceStatus represents information about the status of a namespace in project.
type NamespaceStatus struct {
	// +optional
	Phase NamespacePhase
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// The reason for the condition's last transition.
	// +optional
	Reason string
	// A human readable message indicating details about the transition.
	// +optional
	Message string
	// +optional
	ResourceQuotaName string
	// Used represents the resources of a namespace that are used.
	// +optional
	Used ResourceList
}

// NamespacePhase indicates the status of namespace in project.
type NamespacePhase string

// These are valid namespace status of project.
const (
	// NamespacePending indicates that the namespace has been declared, when the namespace
	// has not actually been created in the cluster.
	NamespacePending NamespacePhase = "Pending"
	// NamespaceAvailable indicates the namespace of the project is available.
	NamespaceAvailable NamespacePhase = "Available"
	// Namespace indicates that the namespace failed to be created in the cluster or
	// deleted in the cluster after it has been created.
	NamespaceFailed NamespacePhase = "Failed"
	// NamespaceTerminating means the namespace is undergoing graceful termination.
	NamespaceTerminating NamespacePhase = "Terminating"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Platform is a platform in TKE.
type Platform struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of platforms in this set.
	// +optional
	Spec PlatformSpec
}

// PlatformSpec is a description of a platform.
type PlatformSpec struct {
	TenantID       string
	Administrators []string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlatformList is the whole list of all platforms which owned by a tenant.
type PlatformList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of platform.
	Items []Platform
}

// +genclient
// +genclient:nonNamespaced
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Portal is a user in TKE.
type Portal struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Administrator indicates whether the user is a platform administrator
	Administrator bool
	// Projects represents the list of projects to which the user belongs, where the key represents
	// project name and the value represents the project display name.
	Projects map[string]string
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMap holds configuration data for tke to consume.
type ConfigMap struct {
	metav1.TypeMeta `json:",inline"`
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
