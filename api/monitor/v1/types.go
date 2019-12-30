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

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Collector is a monitor component.
type Collector struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec CollectorSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status CollectorStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CollectorList is the whole list of all collectors which owned by a tenant.
type CollectorList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of Collector
	Items []Collector `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// CollectorSpec describes the attributes on a Collector.
type CollectorSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	// +optional
	Type CollectorType `json:"type,omitempty" protobuf:"bytes,3,opt,name=type"`
	// Version is the components version.
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,4,opt,name=version"`
	// Storage is the remote address for collector when writing/reading outside of cluster.
	// +optional
	Storage CollectorStorage `json:"storage,omitempty" protobuf:"bytes,5,opt,name=storage"`
	// NotifyWebhook is the address that alert messages send to, optional. If not set, a default webhook address "https://[notify-api-address]/webhook" will be used.
	// +optional
	NotifyWebhook string `json:"notifyWebhook,omitempty" protobuf:"bytes,6,opt,name=notifyWebhook"`
}

// CollectorStatus is information about the current status of a Collector.
type CollectorStatus struct {
	// Version is the version of collector.
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// Components is the components version such as node-exporter and alert manager.
	// +optional
	Components map[string]string `json:"components,omitempty" protobuf:"bytes,2,rep,name=components"`
	// Phase is the current lifecycle phase of the helm of cluster.
	// +optional
	Phase CollectorPhase `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase"`
	// Reason is a brief CamelCase string that describes any failure.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// RetryCount is a int between 0 and 5 that describes the time of retrying initializing.
	// +optional
	RetryCount int32 `json:"retryCount" protobuf:"varint,5,name=retryCount"`
	// LastReInitializingTimestamp is a timestamp that describes the last time of retrying initializing.
	// +optional
	LastReInitializingTimestamp metav1.Time `json:"lastReInitializingTimestamp" protobuf:"bytes,6,name=lastReInitializingTimestamp"`
}

// CollectorStorage is the remote write/read address for collector.
type CollectorStorage struct {
	WriteAddr []string `json:"writeAddr,omitempty" protobuf:"bytes,1,opt,name=writeAddr"`
	ReadAddr  []string `json:"readAddr,omitempty" protobuf:"bytes,2,opt,name=readAddr"`
}

// CollectorPhase defines the phase of collector constructor.
type CollectorPhase string

const (
	// CollectorPhaseInitializing means is wait initializing.
	CollectorPhaseInitializing CollectorPhase = "Initializing"
	// CollectorPhaseReinitializing means is reinitializing.
	CollectorPhaseReinitializing CollectorPhase = "Reinitializing"
	// CollectorPhaseChecking means is wait checking.
	CollectorPhaseChecking CollectorPhase = "Checking"
	// CollectorPhaseRunning means is running.
	CollectorPhaseRunning CollectorPhase = "Running"
	// CollectorPhaseUpgrading means is upgrading.
	CollectorPhaseUpgrading CollectorPhase = "Upgrading"
	// CollectorPhaseFailed means has been failed.
	CollectorPhaseFailed CollectorPhase = "Failed"
)

// CollectorType defines the type of collector.
type CollectorType string

const (
	// CollectorManaged means the collector managed by TKE.
	CollectorManaged CollectorType = "Managed"
	// CollectorImported means the prometheus installed by other.
	CollectorImportedPrometheus CollectorType = "ImportedPrometheus"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlarmPolicy is a policy of alarm system.
type AlarmPolicy struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of alarm policies in this set.
	// +optional
	Spec AlarmPolicySpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status AlarmPolicyStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlarmPolicyList is the whole list of all alarm policies which owned by a tenant.
type AlarmPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of alarm policies.
	Items []AlarmPolicy `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// AlarmPolicySpec describes the attributes on an alarm policy.
type AlarmPolicySpec struct {
	TenantID    string          `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string          `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Type        AlarmPolicyType `json:"type" protobuf:"bytes,3,opt,name=type"`
	// +patchMergeKey=metricName
	// +patchStrategy=merge
	Metrics          []AlarmMetric    `json:"metrics" protobuf:"bytes,4,rep,name=metrics" patchStrategy:"merge" patchMergeKey:"metricName"`
	Objects          string           `json:"objects" protobuf:"bytes,5,opt,name=objects"`
	ObjectsType      AlarmObjectsType `json:"objectsType" protobuf:"bytes,6,opt,name=objectsType"`
	StatisticsPeriod int64            `json:"statisticsPeriod,omitempty" protobuf:"varint,7,opt,name=statisticsPeriod"`
	// +optional
	Namespace *string `json:"namespace,omitempty" protobuf:"bytes,8,opt,name=namespace"`
	// +optional
	WorkloadType *WorkloadType `json:"workloadType,omitempty" protobuf:"bytes,9,opt,name=workloadType"`
	// +optional
	// +patchStrategy=merge
	ReceiverGroups []string `json:"receiverGroups,omitempty" protobuf:"bytes,10,opt,name=receiverGroups" patchStrategy:"merge"`
	// +optional
	// +patchStrategy=merge
	Receivers []string `json:"receivers,omitempty" protobuf:"bytes,11,opt,name=receivers" patchStrategy:"merge"`
	// +optional
	// +patchMergeKey=templateName
	// +patchStrategy=merge
	NotifyWays []AlarmNotifyWay `json:"notifyWays,omitempty" protobuf:"bytes,12,rep,name=notifyWays" patchStrategy:"merge" patchMergeKey:"templateName"`
}

// AlarmPolicyStatus is information about the current status of a AlarmPolicy.
type AlarmPolicyStatus struct {
	// +optional
	Phase AlarmPolicyPhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=AlarmPolicyPhase"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,2,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
}

// AlarmPolicyType defines the type of alarm policy.
type AlarmPolicyType string

const (
	// AlarmPolicyCluster indicates a cluster-wide alarm policy.
	AlarmPolicyCluster AlarmPolicyType = "Cluster"
	// AlarmPolicyNode indicates a node-wide alarm policy.
	AlarmPolicyNode AlarmPolicyType = "Node"
	// AlarmPolicyPod indicates a pod-wide alarm policy.
	AlarmPolicyPod AlarmPolicyType = "Pod"
)

type AlarmObjectsType string

const (
	AlarmObjectsAll  AlarmObjectsType = "All"
	AlarmObjectsPart AlarmObjectsType = "Part"
)

type AlarmMetric struct {
	Measurement    string `json:"measurement" protobuf:"bytes,1,opt,name=measurement"`
	MetricName     string `json:"metricName" protobuf:"bytes,2,opt,name=metricName"`
	ContinuePeriod int64  `json:"continuePeriod" protobuf:"varint,3,opt,name=continuePeriod"`
	// +optional
	DisplayName string `json:"displayName,omitempty" protobuf:"bytes,4,opt,name=displayName"`
	// +optional
	Evaluator *AlarmEvaluator `json:"evaluator,omitempty" protobuf:"bytes,5,opt,name=evaluator"`
	// +optional
	Unit string `json:"unit,omitempty" protobuf:"bytes,6,opt,name=unit"`
}

type AlarmEvaluator struct {
	Type  string `json:"type" protobuf:"bytes,1,opt,name=type"`
	Value string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

type WorkloadType string

const (
	WorkloadDeployment  WorkloadType = "Deployment"
	WorkloadDaemonset   WorkloadType = "Daemonset"
	WorkloadStatefulset WorkloadType = "Statefulset"
)

type AlarmNotifyWay struct {
	ChannelName  string `json:"channelName" protobuf:"bytes,1,opt,name=channelName"`
	TemplateName string `json:"templateName" protobuf:"bytes,2,opt,name=templateName"`
}

// AlarmPolicyPhase indicates the status of policy alarm in cluster.
type AlarmPolicyPhase string

// These are valid alarm policy status.
const (
	// AlarmPolicyPending indicates that the alarm policy has been declared, when
	// the alarm policy has not actually been created in the cluster.
	AlarmPolicyPending AlarmPolicyPhase = "Pending"
	// AlarmPolicyAvailable indicates the alarm policy is available.
	AlarmPolicyAvailable AlarmPolicyPhase = "Available"
	// AlarmPolicyFailed indicates that the alarm policy failed to be created in the
	// cluster or deleted in the cluster after it has been created.
	AlarmPolicyFailed AlarmPolicyPhase = "Failed"
	// AlarmPolicyTerminating means the alarm policy is undergoing graceful
	// termination.
	AlarmPolicyTerminating AlarmPolicyPhase = "Terminating"
)
