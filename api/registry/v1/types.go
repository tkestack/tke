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

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartGroup is a chart container in chartmuseum registry.
type ChartGroup struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of chart group in this set.
	// +optional
	Spec ChartGroupSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ChartGroupStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartGroupList is the whole list of all chart groups which owned by a tenant.
type ChartGroupList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of chart groups
	Items []ChartGroup `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ChartGroupSpec is a description of a chart group.
type ChartGroupSpec struct {
	Name     string `json:"name" protobuf:"bytes,1,opt,name=name"`
	TenantID string `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	// +optional
	DisplayName string `json:"displayName,omitempty" protobuf:"bytes,3,opt,name=displayName"`
	// +optional
	Visibility Visibility `json:"visibility,omitempty" protobuf:"bytes,4,opt,name=visibility,casttype=Visibility"`
	// +optional
	Type RepoType `json:"type,omitempty" protobuf:"bytes,5,opt,name=type"`
	// +optional
	Description string `json:"description,omitempty" protobuf:"bytes,6,opt,name=description"`
	// +optional
	Projects []string `json:"projects,omitempty" protobuf:"bytes,7,opt,name=projects"`
	// +optional
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,8,rep,name=finalizers,casttype=FinalizerName"`
	// +optional
	Users []string `json:"users,omitempty" protobuf:"bytes,9,opt,name=users"`
	// +optional
	ImportedInfo ChartGroupImport `json:"importedInfo,omitempty" protobuf:"bytes,10,opt,name=importedInfo"`
	// +optional
	Creator string `json:"creator,omitempty" protobuf:"bytes,11,opt,name=creator"`
}

// ChartGroupImport is a description of an import chart group.
type ChartGroupImport struct {
	Addr string `json:"addr" protobuf:"bytes,1,opt,name=addr"`
	// +optional
	Username string `json:"username,omitempty" protobuf:"bytes,2,opt,name=username"`
	// +optional
	Password string `json:"password,omitempty" protobuf:"bytes,3,opt,name=password"`
}

// ChartGroupStatus represents information about the status of a chart group.
type ChartGroupStatus struct {
	// +optional
	Locked     *bool `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`
	ChartCount int32 `json:"chartCount" protobuf:"varint,2,opt,name=chartCount"`
	// +optional
	Phase ChartGroupPhase `json:"phase" protobuf:"bytes,3,opt,name=phase,casttype=ChartGroupPhase"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
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
	// ChartGroupFailed indicates that the chart group failed to be created or deleted
	// after it has been created.
	ChartGroupFailed ChartGroupPhase = "Failed"
	// ChartGroupTerminating means the chart group is undergoing graceful termination.
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

// Chart is a chart in chart group of chartmuseum registry.
type Chart struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of chart in this set.
	// +optional
	Spec ChartSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ChartStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartList is the whole list of all charts which owned by a chart group.
type ChartList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of charts
	Items []Chart `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ChartSpec struct {
	Name           string `json:"name" protobuf:"bytes,1,opt,name=name"`
	TenantID       string `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	ChartGroupName string `json:"chartGroupName" protobuf:"bytes,3,opt,name=chartGroupName"`
	// +optional
	DisplayName string `json:"displayName,omitempty" protobuf:"bytes,4,opt,name=displayName"`
	// +optional
	Visibility Visibility `json:"visibility,omitempty" protobuf:"bytes,5,opt,name=visibility,casttype=Visibility"`
	// +optional
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,6,rep,name=finalizers,casttype=FinalizerName"`
}

type ChartStatus struct {
	// +optional
	Locked    *bool          `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`
	PullCount int32          `json:"pullCount" protobuf:"varint,2,opt,name=pullCount"`
	Versions  []ChartVersion `json:"versions" protobuf:"bytes,3,rep,name=versions"`
	// +optional
	Phase ChartPhase `json:"phase" protobuf:"bytes,4,opt,name=phase,casttype=ChartPhase"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,5,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,6,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,7,opt,name=message"`
}

type ChartVersion struct {
	Version     string      `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	ChartSize   int64       `json:"chartSize,omitempty" protobuf:"varint,2,opt,name=chartSize"`
	TimeCreated metav1.Time `json:"timeCreated,omitempty" protobuf:"bytes,3,opt,name=timeCreated"`
	Description string      `json:"description,omitempty" protobuf:"bytes,4,opt,name=description"`
	AppVersion  string      `json:"appVersion,omitempty" protobuf:"bytes,5,opt,name=appVersion"`
	Icon        string      `json:"icon,omitempty" protobuf:"bytes,6,opt,name=icon"`
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
	metav1.TypeMeta `json:",inline"`

	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// +optional
	Cluster string `json:"cluster,omitempty" protobuf:"bytes,2,opt,name=cluster"`
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

// +genclient
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChartInfo describes detail of a chart version.
type ChartInfo struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of a chart.
	// +optional
	Spec ChartInfoSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// ChartInfoSpec is a description of a ChartInfo.
type ChartInfoSpec struct {
	// +optional
	Readme map[string]string `json:"readme,omitempty" protobuf:"bytes,1,opt,name=readme"`
	// +optional
	Values map[string]string `json:"values,omitempty" protobuf:"bytes,2,opt,name=values"`
	// +optional
	RawFiles map[string]string `json:"rawFiles,omitempty" protobuf:"bytes,3,opt,name=rawFiles"`
	// +optional
	ChartSpec `json:",inline" protobuf:"bytes,4,opt,name=chartSpec"`
	// +optional
	ChartVersion `json:",inline" protobuf:"bytes,5,opt,name=chartVersion"`
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
