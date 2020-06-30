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

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[string]resource.Quantity

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList"`
	Requests ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Prometheus is a kubernetes package manager.
type Prometheus struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec PrometheusSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status PrometheusStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PrometheusList is the whole list of all prometheus which owned by a tenant.
type PrometheusList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of Prometheuss
	Items []Prometheus `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// PrometheusSpec describes the attributes on a Prometheus.
type PrometheusSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
	// SubVersion is the components version such as node-exporter.
	SubVersion map[string]string `json:"subVersion,omitempty" protobuf:"bytes,4,opt,name=subVersion"`
	// RemoteAddress is the remote address for prometheus when writing/reading outside of cluster.
	RemoteAddress PrometheusRemoteAddr `json:"remoteAddress,omitempty" protobuf:"bytes,5,opt,name=remoteAddress"`
	// +optional
	// NotifyWebhook is the address that alert messages send to, optional. If not set, a default webhook address "https://[notify-api-address]/webhook" will be used.
	NotifyWebhook string `json:"notifyWebhook,omitempty" protobuf:"bytes,6,opt,name=notifyWebhook"`
	// +optional
	// Resources is the resource request and limit for prometheus
	Resources ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,7,opt,name=resources"`
	// +optional
	// RunOnMaster indicates whether to add master Affinity for all monitor components or not
	RunOnMaster bool `json:"runOnMaster,omitempty" protobuf:"bytes,8,opt,name=runOnMaster"`
	// +optional
	// AlertRepeatInterval indicates repeat interval of alerts
	AlertRepeatInterval string `json:"alertRepeatInterval,omitempty" protobuf:"bytes,9,opt,name=alertRepeatInterval"`
}

// PrometheusStatus is information about the current status of a Prometheus.
type PrometheusStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// Phase is the current lifecycle phase of the helm of cluster.
	// +optional
	Phase AddonPhase `json:"phase,omitempty" protobuf:"bytes,2,opt,name=phase"`
	// Reason is a brief CamelCase string that describes any failure.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	// RetryCount is a int between 0 and 5 that describes the time of retrying initializing.
	// +optional
	RetryCount int32 `json:"retryCount" protobuf:"varint,4,name=retryCount"`
	// LastReInitializingTimestamp is a timestamp that describes the last time of retrying initializing.
	// +optional
	LastReInitializingTimestamp metav1.Time `json:"lastReInitializingTimestamp" protobuf:"bytes,5,name=lastReInitializingTimestamp"`
	// SubVersion is the components version such as node-exporter.
	SubVersion map[string]string `json:"subVersion,omitempty" protobuf:"bytes,6,opt,name=subVersion"`
}

// PrometheusRemoteAddr is the remote write/read address for prometheus
type PrometheusRemoteAddr struct {
	WriteAddr []string `json:"writeAddr,omitempty" protobuf:"bytes,1,opt,name=writeAddr"`
	ReadAddr  []string `json:"readAddr,omitempty" protobuf:"bytes,2,opt,name=readAddr"`
}

// AddonPhase defines the phase of helm constructor.
type AddonPhase string

const (
	// AddonPhaseInitializing means is wait initializing.
	AddonPhaseInitializing AddonPhase = "Initializing"
	// AddonPhaseReinitializing means is reinitializing.
	AddonPhaseReinitializing AddonPhase = "Reinitializing"
	// AddonPhaseChecking means is wait checking.
	AddonPhaseChecking AddonPhase = "Checking"
	// AddonPhaseRunning means is running.
	AddonPhaseRunning AddonPhase = "Running"
	// AddonPhaseUpgrading means is upgrading.
	AddonPhaseUpgrading AddonPhase = "Upgrading"
	// AddonPhaseFailed means has been failed.
	AddonPhaseFailed AddonPhase = "Failed"
	// AddonPhasePending means the controller is proceeding deploying
	AddonPhasePending AddonPhase = "Pending"
	// AddonPhaseUnhealthy means some pods of GPUManager is partial running
	AddonPhaseUnhealthy AddonPhase = "Unhealthy"
	// AddonPhaseTerminating means addon terminating
	AddonPhaseTerminating AddonPhase = "Terminating"
	// AddonPhaseUnknown means addon unknown
	AddonPhaseUnknown AddonPhase = "Unknown"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=create
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Metric defines the structure for querying monitoring data requests and results.
type Metric struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +optional
	Query MetricQuery `json:"query,omitempty" protobuf:"bytes,2,opt,name=query"`
	// +optional
	JSONResult string `json:"jsonResult,omitempty" protobuf:"bytes,3,opt,name=jsonResult"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MetricList is the whole list of all metrics.
type MetricList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of metrics
	Items []Metric `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type MetricQuery struct {
	Table string `json:"table" protobuf:"bytes,1,opt,name=table"`
	// +optional
	StartTime *int64 `json:"startTime,omitempty" protobuf:"varint,2,opt,name=startTime"`
	// +optional
	EndTime    *int64                 `json:"endTime,omitempty" protobuf:"varint,3,opt,name=endTime"`
	Fields     []string               `json:"fields" protobuf:"bytes,4,rep,name=fields"`
	Conditions []MetricQueryCondition `json:"conditions" protobuf:"bytes,5,rep,name=conditions"`
	// +optional
	OrderBy string `json:"orderBy,omitempty" protobuf:"bytes,6,opt,name=orderBy"`
	// +optional
	Order   string   `json:"order,omitempty" protobuf:"bytes,7,opt,name=order"`
	GroupBy []string `json:"groupBy" protobuf:"bytes,8,rep,name=groupBy"`
	Limit   int32    `json:"limit" protobuf:"varint,9,opt,name=limit"`
	Offset  int32    `json:"offset" protobuf:"varint,10,opt,name=offset"`
}

type MetricQueryCondition struct {
	Key   string `json:"key" protobuf:"bytes,1,opt,name=key"`
	Expr  string `json:"expr" protobuf:"bytes,2,opt,name=expr"`
	Value string `json:"value" protobuf:"bytes,3,opt,name=value"`
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
