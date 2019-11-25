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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Project is a project in TKE.
type Project struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of project in this set.
	// +optional
	Spec ProjectSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ProjectStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProjectList is the whole list of all projects which owned by a tenant.
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of projects
	Items []Project `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ProjectSpec is a description of a project.
type ProjectSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers  []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=FinalizerName"`
	TenantID    string          `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	DisplayName string          `json:"displayName" protobuf:"bytes,3,opt,name=displayName"`

	// Members represents the user list of project.
	Members []string `json:"members" protobuf:"bytes,4,rep,name=members"`

	// ParentProjectName indicates the superior project name of this service.
	// +optional
	ParentProjectName string `json:"parentProjectName,omitempty" protobuf:"bytes,5,opt,name=parentProjectName"`
	// Clusters represents clusters that can be used and the resource limits of each cluster.
	// +optional
	Clusters ClusterHard `json:"clusters,omitempty" protobuf:"bytes,6,rep,name=clusters,casttype=ClusterHard"`
}

// ProjectStatus represents information about the status of a project.
type ProjectStatus struct {
	// +optional
	Locked *bool `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`
	// +optional
	Phase ProjectPhase `json:"phase,omitempty" protobuf:"bytes,2,opt,name=phase,casttype=ProjectPhase"`
	// Clusters represents clusters that have been used and the resource usage of each cluster.
	// +optional
	Clusters ClusterUsed `json:"clusters,omitempty" protobuf:"bytes,3,rep,name=clusters,casttype=ClusterUsed"`
	// +optional
	CalculatedChildProjects []string `json:"calculatedChildProjects,omitempty" protobuf:"bytes,4,rep,name=calculatedChildProjects"`
	// +optional
	CalculatedNamespaces []string `json:"calculatedNamespaces,omitempty" protobuf:"bytes,5,rep,name=calculatedNamespaces"`
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
	Hard ResourceList `json:"hard,omitempty" protobuf:"bytes,1,rep,name=hard,casttype=ResourceList"`
}

// ClusterHard is a set of (cluster name, HardQuantity) pairs.
type ClusterHard map[string]HardQuantity

// UsedQuantity is a straightforward wrapper of ResourceList.
type UsedQuantity struct {
	Used ResourceList `json:"used,omitempty" protobuf:"bytes,1,rep,name=used,casttype=ResourceList"`
}

// ClusterUsed is a set of (cluster name, UsedQuantity) pairs.
type ClusterUsed map[string]UsedQuantity

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Namespace is a namespace in cluster.
type Namespace struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of namespaces in this set.
	// +optional
	Spec NamespaceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status NamespaceStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceList is the whole list of all namespaces which owned by a tenant.
type NamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of namespaces
	Items []Namespace `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// NamespaceSpec represents a namespace in cluster of a project.
type NamespaceSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers  []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=FinalizerName"`
	TenantID    string          `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	ClusterName string          `json:"clusterName" protobuf:"bytes,3,opt,name=clusterName"`
	Namespace   string          `json:"namespace" protobuf:"bytes,4,opt,name=namespace"`
	// Hard represents the total resources of a namespace.
	// +optional
	Hard ResourceList `json:"hard,omitempty" protobuf:"bytes,5,rep,name=hard,casttype=ResourceList"`
}

// NamespaceStatus represents information about the status of a namespace in project.
type NamespaceStatus struct {
	// +optional
	Phase NamespacePhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=NamespacePhase"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,2,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// +optional
	ResourceQuotaName string `json:"resourceQuotaName,omitempty" protobuf:"bytes,5,opt,name=resourceQuotaName"`
	// Used represents the resources of a namespace that are used.
	// +optional
	Used ResourceList `json:"used,omitempty" protobuf:"bytes,6,rep,name=used,casttype=ResourceList"`
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
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of platforms in this set.
	// +optional
	Spec PlatformSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// PlatformSpec is a description of a platform.
type PlatformSpec struct {
	TenantID       string   `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	Administrators []string `json:"administrators" protobuf:"bytes,2,rep,name=administrators"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlatformList is the whole list of all platforms which owned by a tenant.
type PlatformList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of platform.
	Items []Platform `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Portal is a user in TKE.
type Portal struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Administrator indicates whether the user is a platform administrator
	Administrator bool `json:"administrator" protobuf:"varint,2,opt,name=administrator"`
	// Projects represents the list of projects to which the user belongs, where the key represents
	// project name and the value represents the project display name.
	Projects map[string]string `json:"projects" protobuf:"bytes,3,rep,name=projects"`
}

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
