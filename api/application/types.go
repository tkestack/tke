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

package application

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// App is a app bootstrap in TKE.
type App struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of bootstrap in this set.
	// +optional
	Spec AppSpec
	// +optional
	Status AppStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppList is the whole list of all bootstraps.
type AppList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of bootstraps
	Items []App
}

// AppSpec is a description of a bootstrap.
type AppSpec struct {
	Type            AppType
	TenantID        string
	Name            string
	TargetCluster   string
	TargetNamespace string
	// +optional
	Chart Chart
	// Values holds the values for this app.
	// +optional
	Values     AppValues
	Finalizers []FinalizerName
	// +optional
	DryRun bool
}

// Chart is a description of a chart.
type Chart struct {
	TenantID       string
	ChartGroupName string
	ChartName      string
	ChartVersion   string
	RepoURL        string
	RepoUsername   string
	RepoPassword   string
	ImportedRepo   bool

	InstallPara InstallPara //parameters used to install a chart
	UpgradePara UpgradePara //parameters used to upgrade a chart

}

//parameters used to install a chart
type InstallPara struct {
	HelmPublicPara
}

//parameters used to upgrade a chart
type UpgradePara struct {
	HelmPublicPara
}

//public parameters used in helm install and helm upgrade command
type HelmPublicPara struct {
	//Client timeout when installiing or upgrading helm release, override default clientTimeOut
	Timeout time.Duration
	// CreateNamespace create namespace when install helm release
	CreateNamespace bool
	// Atomic, if true, for install case, will uninstall failed release, for upgrade case, will roll back on failure.
	Atomic bool
	// Wait, if true, will wait until all Pods, PVCs, Services, and minimum number of Pods of a Deployment,StatefulSet,
	//or ReplicaSet are in a ready state before marking the release as successful, or wait until client timeout
	Wait bool
	// WaitForJobs, if true, wait until all Jobs have been completed before marking the release as successful
	// or wait until client timeout
	WaitForJobs bool
}

// AppStatus represents information about the status of a bootstrap.
type AppStatus struct {
	// Phase the release is in, one of ('ChartFetched',
	// 'ChartFetchFailed', 'Installing', 'Upgrading', 'Succeeded',
	// 'RollingBack', 'RolledBack', 'RollbackFailed')
	// +optional
	Phase AppPhase
	// ObservedGeneration is the most recent generation observed by
	// the operator.
	// +optional
	ObservedGeneration int64
	// ReleaseStatus is the status as given by Helm for the release
	// managed by this resource.
	// +optional
	ReleaseStatus string
	// ReleaseLastUpdated is the last updated time for the release
	// +optional
	ReleaseLastUpdated metav1.Time
	// Revision holds the Git hash or version of the chart currently
	// deployed.
	// +optional
	Revision int64
	// RollbackRevision specify the target rollback version of the chart
	// +optional
	RollbackRevision int64
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// The reason for the condition's last transition.
	// +optional
	Reason string
	// A human readable message indicating details about the transition.
	// +optional
	Message string
	// Dryrun result.
	// +optional
	Manifest string
}

// +genclient
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppHistory is a app history in TKE.
type AppHistory struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of bootstrap in this set.
	// +optional
	Spec AppHistorySpec
}

// AppHistorySpec is a description of a AppHistory.
type AppHistorySpec struct {
	Type            AppType
	TenantID        string
	Name            string
	TargetCluster   string
	TargetNamespace string
	// +optional
	Histories []History
}

// History is a history of a app.
type History struct {
	Revision    int64
	Updated     metav1.Time
	Status      string
	Chart       string
	AppVersion  string
	Description string
	Manifest    string
}

// +genclient
// +genclient:noVerbs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppResource is a app resource in TKE.
type AppResource struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of bootstrap in this set.
	// +optional
	Spec AppResourceSpec
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
	Type            AppType
	TenantID        string
	Name            string
	TargetCluster   string
	TargetNamespace string
	// +optional
	Resources Resources
}

// AppValues string the values for this app.
type AppValues struct {
	RawValuesType RawValuesType
	RawValues     string
	Values        []string
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
	metav1.TypeMeta

	Revision int64
	Cluster  string
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
