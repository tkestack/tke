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

// AddonPhase defines the phase of logagent constructor.
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
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogAgent is a manager to collect logs of workload.
type LogAgent struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of LogCollector.
	// +optional
	Spec LogAgentSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status LogAgentStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogAgentList is the whole list of all logagents which owned by a tenant.
type LogAgentList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of volume decorators.
	Items []LogAgent `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// LogAgentSpec describes the attributes of a Logagent.
type LogAgentSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// LogAgentStatus is information about the current status of a Logagent.
type LogAgentStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// Phase is the current lifecycle phase of the LogCollector of cluster.
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
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogAgentProxyOptions is the query options to a kube-apiserver proxy call for LogAgent crd object.
type LogAgentProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogFileTree
type LogFileTree struct {
	metav1.TypeMeta `json:",inline"`
	Spec            LogFileTreeSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

type LogFileTreeSpec struct {
	ClusterId string `json:"clusterId,omitempty" protobuf:"bytes,1,opt,name=clusterId"`
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
	Container string `json:"container,omitempty" protobuf:"bytes,3,opt,name=container"`
	Pod       string `json:"pod,omitempty" protobuf:"bytes,4,opt,name=pod"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogFileContent
type LogFileContent struct {
	metav1.TypeMeta `json:",inline"`
	Spec            LogFileContentSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

type LogFileContentSpec struct {
	ClusterId string `json:"clusterId,omitempty" protobuf:"bytes,1,opt,name=clusterId"`
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
	Container string `json:"container,omitempty" protobuf:"bytes,3,opt,name=container"`
	Pod       string `json:"pod,omitempty" protobuf:"bytes,4,opt,name=pod"`
	Start     int32  `json:"start,omitempty" protobuf:"varint,5,opt,name=start"`
	Length    int32  `json:"length,omitempty" protobuf:"varint,6,opt,name=length"`
	Filepath  string `json:"filepath,omitempty" protobuf:"bytes,7,opt,name=filepath"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogFileProxyOptions is the query options to a kube-apiserver proxy call for LogFile requests.
type LogFileProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Pod       string `json:"pod,omitempty" protobuf:"bytes,2,opt,name=pod"`
	Container string `json:"container,omitempty" protobuf:"bytes,3,opt,name=container"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogEsDetection
type LogEsDetection struct {
	metav1.TypeMeta `json:",inline"`

	Scheme   string `json:"scheme,omitempty" protobuf:"bytes,1,opt,name=scheme"`
	IP       string `json:"ip,omitempty" protobuf:"bytes,2,opt,name=ip"`
	Port     string `json:"port,omitempty" protobuf:"bytes,3,opt,name=port"`
	User     string `json:"user,omitempty" protobuf:"bytes,4,opt,name=user"`
	Password string `json:"password,omitempty" protobuf:"bytes,5,opt,name=password"`
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
