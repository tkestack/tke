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

const CertOptionValidDays = "validDays"

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
	Finalizers []FinalizerName
	TenantID   string
	// +optional
	DisplayName string

	// Users represents the user list of project.
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
	// +optional
	CachedSpecClusters ClusterHard
	// +optional
	CachedParent *string
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// The reason for the condition's last transition.
	// +optional
	Reason string
	// A human readable message indicating details about the transition.
	// +optional
	Message string
}

// ProjectPhase defines the phase of project constructor.
type ProjectPhase string

const (
	// ProjectActive indicates the project is active.
	ProjectActive ProjectPhase = "Active"
	// ProjectPending indicates that the project has been declared.
	ProjectPending ProjectPhase = "Pending"
	// ProjectTerminating means the project is undergoing graceful termination.
	ProjectTerminating ProjectPhase = "Terminating"
	// ProjectFailed indicates that the project has been failed.
	ProjectFailed ProjectPhase = "Failed"
)

// FinalizerName is the name identifying a finalizer during project lifecycle.
type FinalizerName string

const (
	// ProjectFinalize is an internal finalizer values to Project.
	ProjectFinalize FinalizerName = "project"
	// NamespaceFinalize is an internal finalizer values to Namespace.
	NamespaceFinalize FinalizerName = "namespace"
	// ImageNamespaceFinalize is an internal finalizer values to ImageNamespace.
	ImageNamespaceFinalize FinalizerName = "imagenamespace"
	// ChartGroupFinalize is an internal finalizer values to ChartGroup.
	ChartGroupFinalize FinalizerName = "chartgroup"
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// NamespaceCertOptions is query options of getting namespace with a x509 certificate.
type NamespaceCertOptions struct {
	metav1.TypeMeta

	ValidDays string
}

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
	Finalizers         []FinalizerName
	TenantID           string
	ClusterName        string
	ClusterType        string
	ClusterVersion     string
	ClusterDisplayName string
	Namespace          string
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
	// +optional
	CachedSpecHard ResourceList
	// +optional
	Certificate *NamespaceCert
}

// NamespaceCert represents a x509 certificate of a namespace in project.
type NamespaceCert struct {
	// +optional
	CertPem []byte
	// +optional
	KeyPem []byte
	// +optional
	CACertPem []byte
	// +optional
	APIServer string
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
	// NamespaceLocked indicates the namespace is locked.
	NamespaceLocked NamespacePhase = "Locked"
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

// PortalProject is a project extension info for portal.
type PortalProject struct {
	// Phases of projects.
	Phase string `json:"phase" protobuf:"bytes,1,opt,name=phase"`
	// Parents of projects.
	Parent string `json:"parent" protobuf:"bytes,2,opt,name=parent"`
}

// ProjectExtension is a map from project name to PortalProject.
type ProjectExtension map[string]PortalProject

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
	// Extension is extension info. for projects.
	Extension ProjectExtension
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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageNamespace is an image namespace.
type ImageNamespace struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of namespaces in this set.
	// +optional
	Spec ImageNamespaceSpec
	// +optional
	Status ImageNamespaceStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageNamespaceList is the whole list of all image namespaces which owned by a tenant.
type ImageNamespaceList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of namespaces
	Items []ImageNamespace
}

// ImageNamespaceSpec represents an image namespace.
type ImageNamespaceSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers []FinalizerName
	Name       string
	TenantID   string
	// +optional
	DisplayName string
}

// ImageNamespaceStatus represents information about the status of an image namespace.
type ImageNamespaceStatus struct {
	// +optional
	Phase ImageNamespacePhase
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// The reason for the condition's last transition.
	// +optional
	Reason string
	// A human readable message indicating details about the transition.
	// +optional
	Message string
}

// ImageNamespacePhase indicates the phase of image namespaces.
type ImageNamespacePhase string

// These are valid phases of image namespaces.
const (
	// ImageNamespacePending indicates that the image namespace has been declared,
	// when the image namespace has not actually been created.
	ImageNamespacePending ImageNamespacePhase = "Pending"
	// ImageNamespaceAvailable indicates the image namespace of the project is available.
	ImageNamespaceAvailable ImageNamespacePhase = "Available"
	// ImageNamespaceLocked indicates the image namespace of the project is locked.
	ImageNamespaceLocked ImageNamespacePhase = "Locked"
	// ImageNamespaceFailed indicates that the image namespace failed to be created or deleted
	// after it has been created.
	ImageNamespaceFailed ImageNamespacePhase = "Failed"
	// ImageNamespaceTerminating means the image namespace is undergoing graceful termination.
	ImageNamespaceTerminating ImageNamespacePhase = "Terminating"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartGroup is a chart group.
type ChartGroup struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of namespaces in this set.
	// +optional
	Spec ChartGroupSpec
	// +optional
	Status ChartGroupStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartGroupList is the whole list of all chart groups which owned by a tenant.
type ChartGroupList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of ChartGroups
	Items []ChartGroup
}

// ChartGroupSpec represents a chart group.
type ChartGroupSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers []FinalizerName
	Name       string
	TenantID   string
	// +optional
	DisplayName string
}

// ChartGroupStatus represents information about the status of a chart group.
type ChartGroupStatus struct {
	// +optional
	Phase ChartGroupPhase
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// The reason for the condition's last transition.
	// +optional
	Reason string
	// A human readable message indicating details about the transition.
	// +optional
	Message string
}

// ChartGroupPhase indicates the phase of chart groups.
type ChartGroupPhase string

// These are valid phases of chart groups.
const (
	// ChartGroupPending indicates that the chart group has been declared,
	// when the chart group has not actually been created.
	ChartGroupPending ChartGroupPhase = "Pending"
	// ChartGroupAvailable indicates the chart group of the project is available.
	ChartGroupAvailable ChartGroupPhase = "Available"
	// ChartGroupLocked indicates the chart group of the project is locked.
	ChartGroupLocked ChartGroupPhase = "Locked"
	// ChartGroupFailed indicates that the chart group failed to be created or deleted
	// after it has been created.
	ChartGroupFailed ChartGroupPhase = "Failed"
	// ChartGroupTerminating means the chart group is undergoing graceful termination.
	ChartGroupTerminating ChartGroupPhase = "Terminating"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NsEmigration is a namespace emigration.
type NsEmigration struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of emigrations in this set.
	// +optional
	Spec NsEmigrationSpec
	// +optional
	Status NsEmigrationStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NsEmigrationList is the whole list of all namespace emigrations which owned by a tenant.
type NsEmigrationList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of namespace emigrations
	Items []NsEmigration
}

// NsEmigrationSpec represents a namespace emigration.
type NsEmigrationSpec struct {
	TenantID    string
	Namespace   string
	NsShowName  string
	Destination string
}

// NsEmigrationStatus represents information about the status of a namespace emigration.
type NsEmigrationStatus struct {
	// +optional
	Phase NsEmigrationPhase
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// The reason for the condition's last transition.
	// +optional
	Reason string
	// A human readable message indicating details about the transition.
	// +optional
	Message string
}

// NsEmigrationPhase indicates the phase of namespace emigrations.
type NsEmigrationPhase string

// These are valid phases of namespace emigrations.
const (
	// NsEmigrationPending indicates that the emigration is waiting to be executed.
	NsEmigrationPending NsEmigrationPhase = "Pending"
	// NsEmigrationOldOneLocked indicates that old namespace has been locked.
	NsEmigrationOldOneLocked NsEmigrationPhase = "OldOneLocked"
	// NsEmigrationOldOneDetached indicates that old namespace has been detached from k8s cluster namespace.
	NsEmigrationOldOneDetached NsEmigrationPhase = "OldOneDetached"
	// NsEmigrationNewOneCreated indicates that new namespace has been created.
	NsEmigrationNewOneCreated NsEmigrationPhase = "NewOneCreated"
	// NsEmigrationOldOneTerminating indicates that old namespace is terminating.
	NsEmigrationOldOneTerminating NsEmigrationPhase = "OldOneTerminating"
	// NsEmigrationFinished indicates that the emigration finished.
	NsEmigrationFinished NsEmigrationPhase = "Finished"
	// NsEmigrationFailed indicates that the emigration failed.
	NsEmigrationFailed NsEmigrationPhase = "Failed"
)
