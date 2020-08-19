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

package monitor

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[string]resource.Quantity

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	Limits   ResourceList
	Requests ResourceList
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Prometheus is a systems monitoring and alerting toolkit.
type Prometheus struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec PrometheusSpec
	// +optional
	Status PrometheusStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PrometheusList is the whole list of all prometheus which owned by a tenant.
type PrometheusList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of Prometheuss
	Items []Prometheus
}

// PrometheusSpec describes the attributes on a Prometheus.
type PrometheusSpec struct {
	TenantID      string
	ClusterName   string
	Version       string
	SubVersion    map[string]string
	RemoteAddress PrometheusRemoteAddr
	// +optional
	NotifyWebhook string
	// +optional
	Resources ResourceRequirements
	// +optional
	RunOnMaster bool
	// +optional
	AlertRepeatInterval string
	// +optional
	WithNPD bool
}

// PrometheusStatus is information about the current status of a Prometheus.
type PrometheusStatus struct {
	// +optional
	Version string
	// Phase is the current lifecycle phase of the helm of cluster.
	// +optional
	Phase AddonPhase
	// Reason is a brief CamelCase string that describes any failure.
	// +optional
	Reason string
	// RetryCount is a int between 0 and 5 that describes the time of retrying initializing.
	// +optional
	RetryCount int32
	// LastReInitializingTimestamp is a timestamp that describes the last time of retrying initializing.
	// +optional
	LastReInitializingTimestamp metav1.Time
	// SubVersion is the components version such as node-exporter.
	SubVersion map[string]string
}

// PrometheusRemoteAddr is the remote write/read address for prometheus
type PrometheusRemoteAddr struct {
	WriteAddr []string
	ReadAddr  []string
}

// AddonPhase defines the phase of addon
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
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// +optional
	Query MetricQuery
	// +optional
	JSONResult string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MetricList is the whole list of all metrics.
type MetricList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of metrics
	Items []Metric
}

type MetricQuery struct {
	Table string
	// +optional
	StartTime *int64
	// +optional
	EndTime    *int64
	Fields     []string
	Conditions []MetricQueryCondition
	// +optional
	OrderBy string
	// +optional
	Order   string
	GroupBy []string
	Limit   int32
	Offset  int32
}

type MetricQueryCondition struct {
	Key   string
	Expr  string
	Value string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterOverview defines the structure for clusters' overview data request and result.
type ClusterOverview struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// +optional
	Result *ClusterOverviewResult
}

type ClusterOverviewResult struct {
	ClusterCount     int32
	ClusterAbnormal  int32
	ProjectCount     int32
	ProjectAbnormal  int32
	NodeCount        int32
	NodeAbnormal     int32
	WorkloadCount    int32
	WorkloadAbnormal int32
	Clusters         []*ClusterStatistic
}

type ClusterStatistic struct {
	ClusterID                string
	ClusterPhase             string
	NodeCount                int32
	NodeAbnormal             int32
	WorkloadCount            int32
	WorkloadAbnormal         int32
	HasMetricServer          bool
	CPUUsed                  float64
	CPURequest               float64
	CPULimit                 float64
	CPUCapacity              float64
	CPUAllocatable           float64
	CPURequestRate           string
	CPUAllocatableRate       string
	CPUUsage                 string
	MemUsed                  int64
	MemRequest               int64
	MemLimit                 int64
	MemCapacity              int64
	MemAllocatable           int64
	MemRequestRate           string
	MemAllocatableRate       string
	MemUsage                 string
	SchedulerHealthy         bool
	ControllerManagerHealthy bool
	EtcdHealthy              bool
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
