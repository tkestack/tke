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

package registry

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FinalizerName is the name identifying a finalizer during resource lifecycle.
type FinalizerName string

const (
	// ChartGroupFinalize is an internal finalizer values to ChartGroup.
	ChartGroupFinalize FinalizerName = "chartgroup"
	// ChartFinalize is an internal finalizer values to Chart.
	ChartFinalize FinalizerName = "chart"
	// RegistryClientUserAgent is the user agent for tke registry client
	RegistryClientUserAgent = "tke-registry-client"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Namespace is an image container in registry.
type Namespace struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of namespace in this set.
	// +optional
	Spec NamespaceSpec
	// +optional
	Status NamespaceStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceList is the whole list of all namespaces which owned by a tenant.
type NamespaceList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of namespaces
	Items []Namespace
}

// NamespaceSpec is a description of a namespace.
type NamespaceSpec struct {
	Name     string
	TenantID string
	// +optional
	DisplayName string
	// +optional
	Visibility Visibility
}

// NamespaceStatus represents information about the status of a namespace.
type NamespaceStatus struct {
	// +optional
	Locked    *bool
	RepoCount int32
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Repository is a repo in namespace of registry.
type Repository struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of repository in this set.
	// +optional
	Spec RepositorySpec
	// +optional
	Status RepositoryStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RepositoryList is the whole list of all repositories which owned by a namespace.
type RepositoryList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of repositories
	Items []Repository
}

type RepositorySpec struct {
	Name          string
	TenantID      string
	NamespaceName string
	// +optional
	DisplayName string
	// +optional
	Visibility Visibility
}

type RepositoryStatus struct {
	// +optional
	Locked    *bool
	PullCount int32
	Tags      []RepositoryTag
}

type RepositoryTag struct {
	Name        string
	Digest      string
	TimeCreated metav1.Time
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartGroup is a chart container in chartmuseum registry.
type ChartGroup struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of chartgroup in this set.
	// +optional
	Spec ChartGroupSpec
	// +optional
	Status ChartGroupStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartGroupList is the whole list of all chartgroups which owned by a tenant.
type ChartGroupList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of chartgroups
	Items []ChartGroup
}

// ChartGroupSpec is a description of a chartgroup.
type ChartGroupSpec struct {
	Name     string
	TenantID string
	// +optional
	DisplayName string
	// +optional
	Visibility Visibility
	// +optional
	Type RepoType
	// +optional
	Description string
	// +optional
	Projects []string
	// +optional
	Finalizers []FinalizerName
	// +optional
	Users []string
	// +optional
	ImportedInfo ChartGroupImport
	// +optional
	Creator string
}

// ChartGroupImport is a description of an import chart group.
type ChartGroupImport struct {
	Addr string
	// +optional
	Username string
	// +optional
	Password string
}

// ChartGroupStatus represents information about the status of a chartgroup.
type ChartGroupStatus struct {
	// +optional
	Locked     *bool
	ChartCount int32
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

// ChartGroupPhase indicates the phase of chartgroups.
type ChartGroupPhase string

// These are valid phases of chartgroups.
const (
	// ChartGroupPending indicates that the chartgroup has been declared,
	// when the chartgroup has not actually been created.
	ChartGroupPending ChartGroupPhase = "Pending"
	// ChartGroupAvailable indicates the chartgroup of the project is available.
	ChartGroupAvailable ChartGroupPhase = "Available"
	// ChartGroupFailed indicates that the chartgroup failed to be created or deleted
	// after it has been created.
	ChartGroupFailed ChartGroupPhase = "Failed"
	// ChartGroupTerminating means the chartgroup is undergoing graceful termination.
	ChartGroupTerminating ChartGroupPhase = "Terminating"
)

// ChartPhase indicates the phase of chart.
type ChartPhase string

// These are valid phases of charts.
const (
	// ChartPending indicates that the chart has been declared,
	// when the chart has not actually been created.
	ChartPending ChartPhase = "Pending"
	// ChartAvailable indicates the chart of the project is available.
	ChartAvailable ChartPhase = "Available"
	// ChartFailed indicates that the chart failed to be created or deleted
	// after it has been created.
	ChartFailed ChartPhase = "Failed"
	// ChartTerminating means the chart is undergoing graceful termination.
	ChartTerminating ChartPhase = "Terminating"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Chart is a chart in chartgroup of chartmuseum registry.
type Chart struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of chart in this set.
	// +optional
	Spec ChartSpec
	// +optional
	Status ChartStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartList is the whole list of all charts which owned by a chartgroup.
type ChartList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of charts
	Items []Chart
}

type ChartSpec struct {
	Name           string
	TenantID       string
	ChartGroupName string
	// +optional
	DisplayName string
	// +optional
	Visibility Visibility
	// +optional
	Finalizers []FinalizerName
}

type ChartStatus struct {
	// +optional
	Locked    *bool
	PullCount int32
	Versions  []ChartVersion
	// +optional
	Phase ChartPhase
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

type ChartVersion struct {
	Version     string
	ChartSize   int64
	TimeCreated metav1.Time
	Description string
	AppVersion  string
	Icon        string
}

// Visibility defines the visible properties of the repo or namespace.
type Visibility string

// RepoType defines the type properties of the repo or namespace.
type RepoType string

const (
	// VisibilityPublic indicates the namespace or repo is public.
	VisibilityPublic Visibility = "Public"
	// VisibilityUser indicates the namespace or repo is user.
	VisibilityUser Visibility = "User"
	// VisibilityProject indicates the namespace or repo is project.
	VisibilityProject Visibility = "Project"

	// VisibilityPrivate indicates the namespace or repo is private.
	// Deprecated!
	VisibilityPrivate Visibility = "Private"

	// RepoTypeSelfBuilt indicates the type of namespace or repo is selfbuilt.
	RepoTypeSelfBuilt RepoType = "SelfBuilt"
	// RepoTypeImported indicates the type of namespace or repo is imported.
	RepoTypeImported RepoType = "Imported"
	// RepoTypeSystem indicates the type of namespace or repo is system.
	RepoTypeSystem RepoType = "System"

	// RepoTypeProject indicates the type of namespace or repo is project.
	// Deprecated!
	RepoTypeProject RepoType = "project"
	// RepoTypePersonal indicates the type of namespace or repo is personal.
	// Deprecated!
	RepoTypePersonal RepoType = "personal"

	// ScopeTypeAll indicates all of namespace or repo is all.
	ScopeTypeAll string = "all"
	// ScopeTypePublic indicates all of namespace or repo is public.
	ScopeTypePublic string = "public"
	// ScopeTypeUser indicates all of namespace or repo is user.
	ScopeTypeUser string = "user"
	// ScopeTypeProject indicates all of namespace or repo is project.
	ScopeTypeProject string = "project"
)

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartProxyOptions is the query options to a ChartInfo proxy call.
type ChartProxyOptions struct {
	metav1.TypeMeta

	Version   string
	Cluster   string
	Namespace string
}

// +genclient
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartInfo describes detail of a chart version.
type ChartInfo struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of a chart.
	// +optional
	Spec ChartInfoSpec
}

// ChartInfoSpec is a description of a ChartInfo.
type ChartInfoSpec struct {
	// +optional
	Readme map[string]string
	// +optional
	Values map[string]string
	// +optional
	RawFiles map[string]string
	// +optional
	ChartSpec
	// +optional
	ChartVersion
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
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
