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
	fmt "fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// App is a app bootstrap in TKE.
type App struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of bootstrap in this set.
	// +optional
	Spec AppSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status AppStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppList is the whole list of all bootstraps.
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of bootstraps
	Items []App `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// AppSpec is a description of a project.
type AppSpec struct {
	Type            AppType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=AppType"`
	TenantID        string  `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	Name            string  `json:"name" protobuf:"bytes,3,opt,name=name"`
	TargetCluster   string  `json:"targetCluster" protobuf:"bytes,4,opt,name=targetCluster"`
	TargetNamespace string  `json:"targetNamespace" protobuf:"bytes,5,opt,name=targetNamespace"`
	// +optional
	Chart Chart `json:"chart" protobuf:"bytes,6,opt,name=chart,casttype=Chart"`
	// Values holds the values for this app.
	// +optional
	Values AppValues `json:"values,omitempty" protobuf:"bytes,7,opt,name=values,casttype=AppValues"`
	// +optional
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,8,rep,name=finalizers,casttype=FinalizerName"`
	DryRun     bool            `json:"dryRun" protobuf:"bytes,9,opt,name=dryRun"`
}

// Chart is a description of a chart.
type Chart struct {
	TenantID       string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ChartGroupName string `json:"chartGroupName" protobuf:"bytes,2,opt,name=chartGroupName"`
	// ChartName is the name of the chart.
	ChartName string `json:"chartName" protobuf:"bytes,3,opt,name=chartName"`
	// ChartVersion is the version of the chart.
	ChartVersion string      `json:"chartVersion" protobuf:"bytes,4,opt,name=chartVersion"`
	RepoURL      string      `json:"repoURL" protobuf:"bytes,5,opt,name=repoURL"`
	RepoUsername string      `json:"repoUsername" protobuf:"bytes,6,opt,name=repoUsername"`
	RepoPassword string      `json:"repoPassword" protobuf:"bytes,7,opt,name=repoPassword"`
	ImportedRepo bool        `json:"importedRepo" protobuf:"bytes,8,opt,name=importedRepo"`
	InstallPara  InstallPara `json:"installPara" protobuf:"bytes,9,opt,name=installPara"`
	UpgradePara  UpgradePara `json:"upgradePara" protobuf:"bytes,10,opt,name=upgradePara"`
}

//parameters used to install a chart
type InstallPara struct {
	HelmPublicPara `json:",inline" protobuf:"bytes,1,opt,name=helmPublicPara"`
}

//parameters used to upgrade a chart
type UpgradePara struct {
	HelmPublicPara `json:",inline" protobuf:"bytes,1,opt,name=helmPublicPara"`
}

//public parameters used in helm install and helm upgrade command
type HelmPublicPara struct {
	//Client timeout when installiing or upgrading helm release, override default clientTimeOut
	Timeout time.Duration `json:"clientTimeout" protobuf:"bytes,1,opt,name=clientTimeout"`
	// CreateNamespace create namespace when install helm release
	CreateNamespace bool `json:"createNamespace" protobuf:"bytes,2,opt,name=createNamespace"`
	// Atomic, if true, for install case, will uninstall failed release, for upgrade case, will roll back on failure.
	Atomic bool `json:"atomic" protobuf:"bytes,3,opt,name=atomic"`
	// Wait, if true, will wait until all Pods, PVCs, Services, and minimum number of Pods of a Deployment,StatefulSet,
	//or ReplicaSet are in a ready state before marking the release as successful, or wait until client timeout
	Wait bool `json:"wait" protobuf:"bytes,4,opt,name=wait"`
	// WaitForJobs, if true, wait until all Jobs have been completed before marking the release as successful
	// or wait until client timeout
	WaitForJobs bool `json:"waitForJobs" protobuf:"bytes,5,opt,name=waitForJobs"`
}

// AppStatus represents information about the status of a bootstrap.
type AppStatus struct {
	// Phase the release is in, one of ('ChartFetched',
	// 'ChartFetchFailed', 'Installing', 'Upgrading', 'Succeeded',
	// 'RollingBack', 'RolledBack', 'RollbackFailed')
	// +optional
	Phase AppPhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=AppPhase"`
	// ObservedGeneration is the most recent generation observed by
	// the operator.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,2,opt,name=observedGeneration"`
	// ReleaseStatus is the status as given by Helm for the release
	// managed by this resource.
	// +optional
	ReleaseStatus string `json:"releaseStatus,omitempty" protobuf:"bytes,3,opt,name=releaseStatus"`
	// ReleaseLastUpdated is the last updated time for the release
	// +optional
	ReleaseLastUpdated metav1.Time `json:"releaseLastUpdated,omitempty" protobuf:"bytes,4,opt,name=releaseLastUpdated"`
	// Revision holds the Git hash or version of the chart currently
	// deployed.
	// +optional
	Revision int64 `json:"revision,omitempty" protobuf:"varint,5,opt,name=revision"`
	// RollbackRevision specify the target rollback version of the chart
	// +optional
	RollbackRevision int64 `json:"rollbackRevision,omitempty" protobuf:"varint,6,opt,name=rollbackRevision"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,7,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,8,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,9,opt,name=message"`
	// Dryrun result.
	// +optional
	Manifest string `json:"manifest" protobuf:"bytes,10,opt,name=manifest"`
}

// +genclient
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppHistory is a app history in TKE.
type AppHistory struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of bootstrap in this set.
	// +optional
	Spec AppHistorySpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// AppHistorySpec is a description of a AppHistory.
type AppHistorySpec struct {
	Type            AppType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=AppType"`
	TenantID        string  `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	Name            string  `json:"name" protobuf:"bytes,3,opt,name=name"`
	TargetCluster   string  `json:"targetCluster" protobuf:"bytes,4,opt,name=targetCluster"`
	TargetNamespace string  `json:"targetNamespace" protobuf:"bytes,5,opt,name=targetNamespace"`
	// +optional
	Histories []History `json:"histories,omitempty" protobuf:"bytes,6,opt,name=histories,casttype=History"`
}

// History is a history of a app.
type History struct {
	Revision    int64       `json:"revision,omitempty" protobuf:"varint,1,opt,name=revision"`
	Updated     metav1.Time `json:"updated,omitempty" protobuf:"bytes,2,opt,name=updated"`
	Status      string      `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	Chart       string      `json:"chart,omitempty" protobuf:"bytes,4,opt,name=chart"`
	AppVersion  string      `json:"appVersion,omitempty" protobuf:"bytes,5,opt,name=appVersion"`
	Description string      `json:"description,omitempty" protobuf:"bytes,6,opt,name=description"`
	Manifest    string      `json:"manifest,omitempty" protobuf:"bytes,7,opt,name=manifest"`
}

// +genclient
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppResource is a app resource in TKE.
type AppResource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of bootstrap in this set.
	// +optional
	Spec AppResourceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// Resources is a map info of different resources.
type Resources map[string]ResourceValues

// ResourceValues masks the value so protobuf can generate
// You can view related issues : https://github.com/kubernetes/kubernetes/issues/46024
// +protobuf.nullable=true
// +protobuf.options.(gogoproto.goproto_stringer)=false
type ResourceValues []string

func (t ResourceValues) String() string {
	return fmt.Sprintf("%v", []string(t))
}

// AppResourceSpec is a description of a AppResource.
type AppResourceSpec struct {
	Type            AppType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=AppType"`
	TenantID        string  `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	Name            string  `json:"name" protobuf:"bytes,3,opt,name=name"`
	TargetCluster   string  `json:"targetCluster" protobuf:"bytes,4,opt,name=targetCluster"`
	TargetNamespace string  `json:"targetNamespace" protobuf:"bytes,5,opt,name=targetNamespace"`
	// +optional
	Resources Resources `json:"resources,omitempty" protobuf:"bytes,6,opt,name=resources,casttype=Resources"`
}

// AppValues string the values for this app.
type AppValues struct {
	RawValuesType RawValuesType `json:"rawValuesType,omitempty" protobuf:"bytes,1,opt,name=rawValuesType,casttype=RawValuesType"`
	RawValues     string        `json:"rawValues,omitempty" protobuf:"bytes,2,opt,name=rawValues"`
	Values        []string      `json:"values,omitempty" protobuf:"bytes,3,opt,name=values"`
}

// RawValuesType indicates the type of rawValues.
type RawValuesType string

const (
	// RawValuesTypeJson means the type of rawValues is json
	RawValuesTypeJson RawValuesType = "json"
	// RawValuesTypeYaml means the type of rawValues is yaml
	RawValuesTypeYaml RawValuesType = "yaml"
)

// FinalizerName is the name identifying a finalizer during resource lifecycle.
type FinalizerName string

const (
	// AppFinalize is an internal finalizer values to App.
	AppFinalize FinalizerName = "app"
)

// AppPhase indicates the phase of app.
type AppPhase string

const (
	// Installing means the installation for the App is running.
	AppPhaseInstalling AppPhase = "Installing"
	// InstallFailed means the installation for the App failed.
	AppPhaseInstallFailed AppPhase = "InstallFailed"
	// Upgrading means the upgrade for the App is running.
	AppPhaseUpgrading AppPhase = "Upgrading"
	// Succeeded means the dry-run, installation, or upgrade for the
	// App succeeded.
	AppPhaseSucceeded AppPhase = "Succeeded"
	// Failed means the upgrade for the App
	// failed.
	AppPhaseUpgradFailed AppPhase = "UpgradFailed"
	// RollingBack means a rollback for the App is running.
	AppPhaseRollingBack AppPhase = "RollingBack"
	// RolledBack means the App has been rolled back.
	AppPhaseRolledBack AppPhase = "RolledBack"
	// RolledBackFailed means the rollback for the App failed.
	AppPhaseRollbackFailed AppPhase = "RollbackFailed"

	// AppPhaseTerminating means the app is undergoing graceful termination.
	AppPhaseTerminating AppPhase = "Terminating"
	// SyncFailed means the synchrony for the App
	// failed.
	AppPhaseSyncFailed AppPhase = "SyncFailed"
)

// AppType indicates the type of app.
type AppType string

const (
	// AppTypeHelmV3 means the app is a Helm3 release
	AppTypeHelmV3 AppType = "HelmV3"
)

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RollbackProxyOptions is the query options to an app rollback proxy call.
type RollbackProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	Revision int64 `json:"revision,omitempty" protobuf:"varint,1,opt,name=revision"`
	// +optional
	Cluster string `json:"cluster,omitempty" protobuf:"bytes,2,opt,name=cluster"`
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
