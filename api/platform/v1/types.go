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

// Cluster is a Kubernetes cluster in
type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec ClusterSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ClusterStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList is the whole list of all clusters which owned by a tenant.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of clusters
	Items []Cluster `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ClusterMachine is the master machine definition of cluster.
type ClusterMachine struct {
	IP       string `json:"ip" protobuf:"bytes,1,opt,name=ip"`
	Port     int32  `json:"port" protobuf:"varint,2,opt,name=port"`
	Username string `json:"username" protobuf:"bytes,3,opt,name=username"`
	// +optional
	Password []byte `json:"password,omitempty" protobuf:"bytes,4,opt,name=password"`
	// +optional
	PrivateKey []byte `json:"privateKey,omitempty" protobuf:"bytes,5,opt,name=privateKey"`
	// +optional
	PassPhrase []byte `json:"passPhrase,omitempty" protobuf:"bytes,6,opt,name=passPhrase"`
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,7,opt,name=labels"`
}

// ClusterSpec is a description of a cluster.
type ClusterSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=FinalizerName"`
	TenantID   string          `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	// +optional
	DisplayName string `json:"displayName" protobuf:"bytes,3,opt,name=displayName"`
	Type        string `json:"type" protobuf:"bytes,4,opt,name=type"`
	Version     string `json:"version" protobuf:"bytes,5,opt,name=version"`
	// +optional
	NetworkType NetworkType `json:"networkType,omitempty" protobuf:"bytes,6,opt,name=networkType,casttype=NetworkType"`
	// +optional
	NetworkDevice string `json:"networkDevice,omitempty" protobuf:"bytes,7,opt,name=networkDevice"`
	// +optional
	ClusterCIDR string `json:"clusterCIDR,omitempty" protobuf:"bytes,8,opt,name=clusterCIDR"`
	// DNSDomain is the dns domain used by k8s services. Defaults to "cluster.local".
	DNSDomain string `json:"dnsDomain,omitempty" protobuf:"bytes,9,opt,name=dnsDomain"`
	// +optional
	PublicAlternativeNames []string `json:"publicAlternativeNames,omitempty" protobuf:"bytes,10,opt,name=publicAlternativeNames"`
	// +optional
	Features ClusterFeature `json:"features,omitempty" protobuf:"bytes,11,opt,name=features,casttype=ClusterFeature"`
	// +optional
	Properties ClusterProperty `json:"properties,omitempty" protobuf:"bytes,12,opt,name=properties,casttype=ClusterProperty"`
	// +optional
	Machines []ClusterMachine `json:"machines,omitempty" protobuf:"bytes,13,rep,name=addresses"`

	// +optional
	DockerExtraArgs map[string]string `json:"dockerExtraArgs,omitempty" protobuf:"bytes,14,name=dockerExtraArgs"`
	// +optional
	KubeletExtraArgs map[string]string `json:"kubeletExtraArgs,omitempty" protobuf:"bytes,15,name=kubeletExtraArgs"`
	// +optional
	APIServerExtraArgs map[string]string `json:"apiServerExtraArgs,omitempty" protobuf:"bytes,16,name=apiServerExtraArgs"`
	// +optional
	ControllerManagerExtraArgs map[string]string `json:"controllerManagerExtraArgs,omitempty" protobuf:"bytes,17,name=controllerManagerExtraArgs"`
	// +optional
	SchedulerExtraArgs map[string]string `json:"schedulerExtraArgs,omitempty" protobuf:"bytes,18,name=schedulerExtraArgs"`
}

// ClusterStatus represents information about the status of a cluster.
type ClusterStatus struct {
	// +optional
	Locked *bool `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`
	// +optional
	Version string `json:"version" protobuf:"bytes,2,opt,name=version"`
	// +optional
	Phase ClusterPhase `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase,casttype=ClusterPhase"`
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ClusterCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,4,rep,name=conditions"`
	// A human readable message indicating details about why the cluster is in this condition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
	// A brief CamelCase message indicating details about why the cluster is in this state.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,6,opt,name=reason"`
	// List of addresses reachable to the cluster.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Addresses []ClusterAddress `json:"addresses,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,7,rep,name=addresses"`
	// +optional
	Resource ClusterResource `json:"resource,omitempty" protobuf:"bytes,9,opt,name=resource,casttype=ClusterResource"`
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Components []ClusterComponent `json:"components,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,10,rep,name=components"`
	// +optional
	ServiceCIDR string `json:"serviceCIDR,omitempty" protobuf:"bytes,11,opt,name=serviceCIDR"`
	// +optional
	NodeCIDRMaskSize int32 `json:"nodeCIDRMaskSize,omitempty" protobuf:"varint,12,opt,name=nodeCIDRMaskSize"`
	// +optional
	DNSIP string `json:"dnsIP,omitempty" protobuf:"bytes,13,opt,name=dnsIP"`
	// +optional
	RegistryIPs []string `json:"registryIPs,omitempty" protobuf:"bytes,14,opt,name=registryIPs"`
}

// FinalizerName is the name identifying a finalizer during cluster lifecycle.
type FinalizerName string

const (
	// ClusterFinalize is an internal finalizer values to Cluster.
	ClusterFinalize FinalizerName = "cluster"

	// MachineFinalize is an internal finalizer values to Machine.
	MachineFinalize FinalizerName = "machine"
)

// NetworkType defines the network type of cluster.
type NetworkType string

const (
	// NetworkPhysics indicates the communication network using the physics network to establish the pod between nodes.
	NetworkPhysics NetworkType = "Physics"
	// NetworkVPC indicates the communication network using the VPC to establish the pod between nodes.
	NetworkVPC NetworkType = "VPC"
	// NetworkFlannel indicates the communication network using the flannel to establish the pod between nodes.
	NetworkFlannel NetworkType = "Flannel"
	// NetworkCalico indicates the communication network using the calico to establish the pod between nodes.
	NetworkCalico NetworkType = "Calico"
	// NetworkIPIP indicates the communication network using the IPIP to establish the pod between nodes.
	NetworkIPIP NetworkType = "IPIP"
)

// GPUType defines the gpu type of cluster.
type GPUType string

const (
	// GPUPhysical indicates the gpu type of cluster is physical.
	GPUPhysical GPUType = "Physical"
	// GPUVirtual indicates the gpu type of cluster is virtual.
	GPUVirtual GPUType = "Virtual"
)

// ClusterPhase defines the phase of cluster constructor.
type ClusterPhase string

const (
	// ClusterRunning is the normal running phase.
	ClusterRunning ClusterPhase = "Running"
	// ClusterInitializing is the initialize phase.
	ClusterInitializing ClusterPhase = "Initializing"
	// ClusterFailed is the failed phase.
	ClusterFailed ClusterPhase = "Failed"
	// ClusterTerminating means the cluster is undergoing graceful termination.
	ClusterTerminating ClusterPhase = "Terminating"
)

// ClusterCondition contains details for the current condition of this cluster.
type ClusterCondition struct {
	// Type is the type of the condition.
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// Last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// AddressType indicates the type of cluster apiserver access address.
type AddressType string

// These are valid address type of cluster.
const (
	// AddressPublic indicates the address of the apiserver accessed from the external network.(such as public lb)
	AddressPublic AddressType = "Public"
	// AddressAdvertise indicates the address of the apiserver accessed from the worker node.(such as internal lb)
	AddressAdvertise AddressType = "Advertise"
	// AddressReal indicates the real address of one apiserver
	AddressReal AddressType = "Real"
	// AddressInternal indicates the address of the apiserver accessed from TKE control plane.
	AddressInternal AddressType = "Internal"
	// AddressSupport used for vpc lb which bind to JNS gateway as known AddressInternal
	AddressSupport AddressType = "Support"
)

// ClusterAddress contains information for the cluster's address.
type ClusterAddress struct {
	// Cluster address type, one of Public, ExternalIP or InternalIP.
	Type AddressType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=AddressType"`
	// The cluster address.
	Host string `json:"host" protobuf:"bytes,2,opt,name=host"`
	Port int32  `json:"port" protobuf:"varint,3,name=port"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCredential records the credential information needed to access the cluster.
type ClusterCredential struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	TenantID    string `json:"tenantID" protobuf:"bytes,2,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,3,opt,name=clusterName"`

	// For TKE in global reuse
	// +optional
	ETCDCACert []byte `json:"etcdCACert,omitempty" protobuf:"bytes,4,opt,name=etcdCACert"`
	// +optional
	ETCDCAKey []byte `json:"etcdCAKey,omitempty" protobuf:"bytes,5,opt,name=etcdCAKey"`
	// +optional
	ETCDAPIClientCert []byte `json:"etcdAPIClientCert,omitempty" protobuf:"bytes,6,opt,name=etcdAPIClientCert"`
	// +optional
	ETCDAPIClientKey []byte `json:"etcdAPIClientKey,omitempty" protobuf:"bytes,7,opt,name=etcdAPIClientKey"`

	// For connect the cluster
	// +optional
	CACert []byte `json:"caCert,omitempty" protobuf:"bytes,8,opt,name=caCert"`
	// +optional
	CAKey []byte `json:"caKey,omitempty" protobuf:"bytes,9,opt,name=caKey"`
	// For kube-apiserver X509 auth
	// +optional
	ClientCert []byte `json:"clientCert,omitempty" protobuf:"bytes,10,opt,name=clientCert"`
	// For kube-apiserver X509 auth
	// +optional
	ClientKey []byte `json:"clientKey,omitempty" protobuf:"bytes,11,opt,name=clientKey"`
	// For kube-apiserver token auth
	// +optional
	Token *string `json:"token,omitempty" protobuf:"bytes,12,opt,name=token"`
	// For kubeadm init or join
	// +optional
	BootstrapToken *string `json:"bootstrapToken,omitempty" protobuf:"bytes,13,opt,name=bootstrapToken"`
	// For kubeadm init or join
	// +optional
	CertificateKey *string `json:"certificateKey,omitempty" protobuf:"bytes,14,opt,name=certificateKey"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCredentialList is the whole list of all ClusterCredential which owned by a tenant.
type ClusterCredentialList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of clusters
	Items []ClusterCredential `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ClusterFeature records the features that are enabled by the cluster.
type ClusterFeature struct {
	// +optional
	IPVS *bool `json:"ipvs,omitempty" protobuf:"varint,1,opt,name=ipvs"`
	// +optional
	PublicLB *bool `json:"publicLB,omitempty" protobuf:"varint,2,opt,name=publicLB"`
	// +optional
	InternalLB *bool `json:"internalLB,omitempty" protobuf:"varint,3,opt,name=internalLB"`
	// +optional
	GPUType *GPUType `json:"gpuType,omitempty" protobuf:"bytes,4,opt,name=gpuType"`
	// +optional
	EnableMasterSchedule bool `json:"enableMasterSchedule,omitempty" protobuf:"bytes,5,opt,name=enableMasterSchedule"`
	// +optional
	HA *HA `json:"ha,omitempty" protobuf:"bytes,6,opt,name=ha"`
	// +optional
	SkipConditions []string `json:"skipConditions,omitempty" protobuf:"bytes,7,opt,name=skipConditions"`
	// +optional
	Files []File `json:"files,omitempty" protobuf:"bytes,8,opt,name=files"`
	// +optional
	Hooks map[HookType]string `json:"hooks,omitempty" protobuf:"bytes,9,opt,name=hooks"`
}

type HA struct {
	TKEHA        *TKEHA        `json:"tke,omitempty" protobuf:"bytes,1,opt,name=tke"`
	ThirdPartyHA *ThirdPartyHA `json:"thirdParty,omitempty" protobuf:"bytes,2,opt,name=thirdParty"`
}

type TKEHA struct {
	VIP string `json:"vip" protobuf:"bytes,1,name=vip"`
}

type ThirdPartyHA struct {
	VIP   string `json:"vip" protobuf:"bytes,1,name=vip"`
	VPort int32  `json:"vport" protobuf:"bytes,2,name=vport"`
}

type File struct {
	Src string `json:"src" protobuf:"bytes,1,name=src"` // Only support regular file
	Dst string `json:"dst" protobuf:"bytes,2,name=dst"`
}

type HookType string

const (
	HookPreInstall  HookType = "PreInstall"
	HookPostInstall HookType = "PostInstall"
)

// ClusterProperty records the attribute information of the cluster.
type ClusterProperty struct {
	// +optional
	MaxClusterServiceNum *int32 `json:"maxClusterServiceNum,omitempty" protobuf:"bytes,1,opt,name=maxClusterServiceNum"`
	// +optional
	MaxNodePodNum *int32 `json:"maxNodePodNum,omitempty" protobuf:"bytes,2,opt,name=maxNodePodNum"`
	// +optional
	OversoldRatio map[string]string `json:"oversoldRatio,omitempty" protobuf:"bytes,3,opt,name=oversoldRatio"`
}

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[string]resource.Quantity

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList"`
	Requests ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList"`
}

// ClusterResource records the current available and maximum resource quota
// information for the cluster.
type ClusterResource struct {
	// Capacity represents the total resources of a cluster.
	// +optional
	Capacity ResourceList `json:"capacity,omitempty" protobuf:"bytes,1,rep,name=capacity,casttype=ResourceList"`
	// Allocatable represents the resources of a cluster that are available for scheduling.
	// Defaults to Capacity.
	// +optional
	Allocatable ResourceList `json:"allocatable,omitempty" protobuf:"bytes,2,rep,name=allocatable,casttype=ResourceList"`
	// +optional
	Allocated ResourceList `json:"allocated,omitempty" protobuf:"bytes,3,rep,name=allocated,casttype=ResourceList"`
}

// ClusterComponent records the number of copies of each component of the
// cluster master.
type ClusterComponent struct {
	Type     string                   `json:"type" protobuf:"bytes,1,opt,name=type"`
	Replicas ClusterComponentReplicas `json:"replicas" protobuf:"bytes,2,opt,name=replicas,casttype=ClusterComponentReplicas"`
}

// ClusterComponentReplicas records the number of copies of each state of each
// component of the cluster master.
type ClusterComponentReplicas struct {
	Desired   int32 `json:"desired" protobuf:"varint,1,name=desired"`
	Current   int32 `json:"current" protobuf:"varint,2,name=current"`
	Available int32 `json:"available" protobuf:"varint,3,name=available"`
	Updated   int32 `json:"updated" protobuf:"varint,4,name=updated"`
}

// AddonLevel indicates the level of cluster addon.
type AddonLevel string

// These are valid level of addon.
const (
	// LevelBasic is level for basic of cluster.
	LevelBasic AddonLevel = "Basic"
	// LevelEnhance is level for enhance of cluster.
	LevelEnhance AddonLevel = "Enhance"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=list,get
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAddon contains the Addon component for the current kubernetes cluster
type ClusterAddon struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec defines the desired identities of addons in this set.
	// +optional
	Spec ClusterAddonSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ClusterAddonStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAddonList is the whole list of all ClusterAddon.
type ClusterAddonList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of ClusterAddon
	Items []ClusterAddon `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ClusterAddonSpec indicates the specifications of the ClusterAddon.
type ClusterAddonSpec struct {
	// Addon type, one of Helm, PersistentEvent or LogCollector etc.
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`
	// AddonLevel is level of cluster addon.
	Level AddonLevel `json:"level" protobuf:"bytes,2,opt,name=level,casttype=AddonLevel"`
	// Version
	Version string `json:"version" protobuf:"bytes,3,opt,name=version"`
}

// ClusterAddonStatus is information about the current status of a ClusterAddon.
type ClusterAddonStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// Phase is the current lifecycle phase of the addon of cluster.
	// +optional
	Phase string `json:"phase,omitempty" protobuf:"bytes,2,opt,name=phase"`
	// Reason is a brief CamelCase string that describes any failure.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=list
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAddonType records the all addons of cluster available.
type ClusterAddonType struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Addon type, one of Helm, PersistentEvent or LogCollector etc.
	Type string `json:"type" protobuf:"bytes,2,opt,name=type"`
	// AddonLevel is level of cluster addon.
	Level AddonLevel `json:"level" protobuf:"bytes,3,opt,name=level,casttype=AddonLevel"`
	// LatestVersion is latest version of the addon.
	LatestVersion string `json:"latestVersion" protobuf:"bytes,4,opt,name=latestVersion"`
	// Description is desc of the addon.
	Description           string   `json:"description,omitempty" protobuf:"bytes,5,opt,name=description"`
	CompatibleClusterType []string `json:"compatibleClusterType,omitempty" protobuf:"bytes,6,rep,name=compatibleClusterType"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAddonTypeList is a resource containing a list of ClusterAddonType objects.
type ClusterAddonTypeList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// +optional
	Items []ClusterAddonType `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterApplyOptions is the query options to a kube-apiserver proxy call for cluster object.
type ClusterApplyOptions struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	NotUpdate bool `json:"notUpdate,omitempty" protobuf:"varint,1,opt,name=notUpdate"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Registry records the third-party image repository information stored by the
// user.
type Registry struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// +optional
	Spec RegistrySpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RegistryList is a resource containing a list of Registry objects.
type RegistryList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// +optional
	Items []Registry `json:"items,omitempty" protobuf:"bytes,2,opt,name=items"`
}

// RegistrySpec indicates the specifications of the third-party image repository.
type RegistrySpec struct {
	// +optional
	TenantID string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	// +optional
	DisplayName string `json:"displayName,omitempty" protobuf:"bytes,2,opt,name=displayName"`
	// +optional
	ClusterName string `json:"clusterName,omitempty" protobuf:"bytes,3,opt,name=clusterName"`
	// +optional
	URL string `json:"url,omitempty" protobuf:"bytes,4,opt,name=url"`
	// +optional
	UserName *string `json:"userName,omitempty" protobuf:"bytes,5,opt,name=userName"`
	// +optional
	Password *string `json:"password,omitempty" protobuf:"bytes,6,opt,name=password"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PersistentEvent is a recorder of kubernetes event.
type PersistentEvent struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec PersistentEventSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status PersistentEventStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PersistentEventList is the whole list of all clusters which owned by a tenant.
type PersistentEventList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of PersistentEvents
	Items []PersistentEvent `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// PersistentEventSpec describes the attributes on a PersistentEvent.
type PersistentEventSpec struct {
	TenantID          string            `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName       string            `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	PersistentBackEnd PersistentBackEnd `json:"persistentBackEnd,omitempty" protobuf:"bytes,3,opt,name=persistentBackEnd"`
	Version           string            `json:"version,omitempty" protobuf:"bytes,4,opt,name=version"`
}

// PersistentEventStatus is information about the current status of a
// PersistentEvent.
type PersistentEventStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// Phase is the current lifecycle phase of the persistent event of cluster.
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

// PersistentBackEnd indicates the backend type and attributes of the persistent
// log store.
type PersistentBackEnd struct {
	CLS *StorageBackEndCLS `json:"cls,omitempty" protobuf:"bytes,1,opt,name=cls"`
	ES  *StorageBackEndES  `json:"es,omitempty" protobuf:"bytes,2,opt,name=es"`
}

// StorageBackEndCLS records the attributes required when the backend storage
// type is CLS.
type StorageBackEndCLS struct {
	LogSetID string `json:"logSetID,omitempty" protobuf:"bytes,1,opt,name=logSetID"`
	TopicID  string `json:"topicID,omitempty" protobuf:"bytes,2,opt,name=topicID"`
}

// StorageBackEndES records the attributes required when the backend storage
// type is ElasticSearch.
type StorageBackEndES struct {
	IP        string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=ip"`
	Port      int32  `json:"port,omitempty" protobuf:"varint,2,opt,name=port"`
	Scheme    string `json:"scheme,omitempty" protobuf:"bytes,3,opt,name=scheme"`
	IndexName string `json:"indexName,omitempty" protobuf:"bytes,4,opt,name=indexName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmProxyOptions is the query options to a Helm-api proxy call.
type HelmProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Path is the URL path to use for the current proxy request to helm-api.
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Helm is a kubernetes package manager.
type Helm struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec HelmSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status HelmStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmList is the whole list of all helms which owned by a tenant.
type HelmList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of Helms
	Items []Helm `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// HelmSpec describes the attributes on a Helm.
type HelmSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// HelmStatus is information about the current status of a Helm.
type HelmStatus struct {
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPAMProxyOptions is the query options to a ipam-api proxy call.
type IPAMProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Path is the URL path to use for the current proxy request to ipam-api.
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPAM is a scheduler plugin for assigning IP.
type IPAM struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec IPAMSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status IPAMStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPAMList is the whole list of all IPAMs which owned by a tenant.
type IPAMList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of IPAMs
	Items []IPAM `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// IPAMSpec describes the attributes on a IPAM.
type IPAMSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// IPAMStatus is information about the current status of a IPAM.
type IPAMStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// Phase is the current lifecycle phase of the addon of cluster.
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

// ConditionStatus defines the status of Condition.
type ConditionStatus string

// These are valid condition statuses.
// "ConditionTrue" means a resource is in the condition.
// "ConditionFalse" means a resource is not in the condition.
// "ConditionUnknown" means server can't decide if a resource is in the condition
// or not.
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// AddonSpec describes the attributes on a Addon.
type AddonSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TappControllerProxyOptions is the query options to a kube-apiserver proxy call.
type TappControllerProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
	Action    string `json:"action,omitempty" protobuf:"bytes,3,opt,name=action"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TappController is a new kubernetes workload.
type TappController struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of tapp controller.
	// +optional
	Spec TappControllerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status TappControllerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TappControllerList is the whole list of all tapp controllers which owned by a tenant.
type TappControllerList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of tapp controllers
	Items []TappController `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// TappControllerSpec describes the attributes on a tapp controller.
type TappControllerSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// TappControllerStatus is information about the current status of a tapp controller.
type TappControllerStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// Phase is the current lifecycle phase of the tapp controller of cluster.
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CSIProxyOptions is the query options to a kube-apiserver proxy call for CSI crd object.
type CSIProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CSIOperator is a operator to manages CSI external components.
type CSIOperator struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of storage operator.
	// +optional
	Spec CSIOperatorSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status CSIOperatorStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CSIOperatorList is the whole list of all storage operators which owned by a tenant.
type CSIOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of storage operators.
	Items []CSIOperator `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// CSIOperatorSpec describes the attributes of a storage operator.
type CSIOperatorSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	// Version of the CSI operator.
	Version string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// CSIOperatorStatus is information about the current status of a storage operator.
type CSIOperatorStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// StorageVendorVersion will be set to the config version of the storage vendor.
	// +optional
	StorageVendorVersion string `json:"storageVendorVersion,omitempty" protobuf:"bytes,2,opt,name=storageVendorVersion"`
	// Phase is the current lifecycle phase of the tapp controller of cluster.
	// +optional
	Phase AddonPhase `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase"`
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PVCRProxyOptions is the query options to a kube-apiserver proxy call for PVCR crd object.
type PVCRProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeDecorator is a controller to manage PVC information.
type VolumeDecorator struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of volume decorator.
	// +optional
	Spec VolumeDecoratorSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status VolumeDecoratorStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeDecoratorList is the whole list of all VolumeDecorator which owned by a tenant.
type VolumeDecoratorList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of volume decorators.
	Items []VolumeDecorator `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// VolumeDecoratorSpec describes the attributes of a VolumeDecorator.
type VolumeDecoratorSpec struct {
	TenantID          string   `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName       string   `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version           string   `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
	VolumeTypes       []string `json:"volumeTypes,omitempty" protobuf:"bytes,4,opt,name=volumeTypes"`
	WorkloadAdmission bool     `json:"workloadAdmission,omitempty" protobuf:"bytes,5,opt,name=workloadAdmission"`
}

// VolumeDecoratorStatus is information about the current status of a VolumeDecorator.
type VolumeDecoratorStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// VolumeTypes is the supported volume types in this cluster.
	// +optional
	VolumeTypes []string `json:"volumeTypes,omitempty" protobuf:"bytes,2,opt,name=volumeTypes"`
	// WorkloadAdmission will be true to enable the workload admission webhook.
	// +optional
	WorkloadAdmission bool `json:"workloadAdmission,omitempty" protobuf:"bytes,3,opt,name=workloadAdmission"`
	// StorageVendorVersion will be set to the config version of the storage vendor.
	// +optional
	StorageVendorVersion string `json:"storageVendorVersion,omitempty" protobuf:"bytes,4,opt,name=storageVendorVersion"`
	// Phase is the current lifecycle phase of the volume decorator of cluster.
	// +optional
	Phase AddonPhase `json:"phase,omitempty" protobuf:"bytes,5,opt,name=phase"`
	// Reason is a brief CamelCase string that describes any failure.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,6,opt,name=reason"`
	// RetryCount is a int between 0 and 5 that describes the time of retrying initializing.
	// +optional
	RetryCount int32 `json:"retryCount" protobuf:"varint,7,name=retryCount"`
	// LastReInitializingTimestamp is a timestamp that describes the last time of retrying initializing.
	// +optional
	LastReInitializingTimestamp metav1.Time `json:"lastReInitializingTimestamp" protobuf:"bytes,8,name=lastReInitializingTimestamp"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogCollectorProxyOptions is the query options to a kube-apiserver proxy call for LogCollector crd object.
type LogCollectorProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogCollector is a manager to collect logs of workload.
type LogCollector struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of LogCollector.
	// +optional
	Spec LogCollectorSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status LogCollectorStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogCollectorList is the whole list of all LogCollector which owned by a tenant.
type LogCollectorList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of volume decorators.
	Items []LogCollector `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// LogCollectorSpec describes the attributes of a LogCollector.
type LogCollectorSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// LogCollectorStatus is information about the current status of a LogCollector.
type LogCollectorStatus struct {
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

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Machine instance in Kubernetes cluster
type Machine struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec defines the desired identities of the Machine.
	// +optional
	Spec MachineSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status MachineStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MachineSpec is a description of machine.
type MachineSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers  []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=FinalizerName"`
	TenantID    string          `json:"tenantID,omitempty" protobuf:"bytes,2,opt,name=tenantID"`
	ClusterName string          `json:"clusterName" protobuf:"bytes,3,opt,name=clusterName"`
	Type        string          `json:"type" protobuf:"bytes,4,opt,name=type"`
	IP          string          `json:"ip" protobuf:"bytes,5,opt,name=ip"`
	Port        int32           `json:"port" protobuf:"varint,6,opt,name=port"`
	Username    string          `json:"username" protobuf:"bytes,7,opt,name=username"`
	// +optional
	Password []byte `json:"password,omitempty" protobuf:"bytes,8,opt,name=password"`
	// +optional
	PrivateKey []byte `json:"privateKey,omitempty" protobuf:"bytes,9,opt,name=privateKey"`
	// +optional
	PassPhrase []byte `json:"passPhrase,omitempty" protobuf:"bytes,10,opt,name=passPhrase"`
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,11,opt,name=labels"`
}

// MachineStatus represents information about the status of an machine.
type MachineStatus struct {
	// +optional
	Locked *bool `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`
	// +optional
	Phase MachinePhase `json:"phase,omitempty" protobuf:"bytes,2,opt,name=phase,casttype=MachinePhase"`
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []MachineCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,3,rep,name=conditions"`
	// A human readable message indicating details about why the machine is in this condition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// A brief CamelCase message indicating details about why the machine is in this state.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// List of addresses reachable to the machine.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Addresses []MachineAddress `json:"addresses,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=addresses"`
	// Set of ids/uuids to uniquely identify the node.
	// +optional
	MachineInfo MachineSystemInfo `json:"machineInfo,omitempty" protobuf:"bytes,7,opt,name=machineInfo"`
}

// MachineSystemInfo is a set of ids/uuids to uniquely identify the node.
type MachineSystemInfo struct {
	// MachineID reported by the node. For unique machine identification
	// in the cluster this field is preferred. Learn more from man(5)
	// machine-id: http://man7.org/linux/man-pages/man5/machine-id.5.html
	MachineID string `json:"machineID,omitempty" protobuf:"bytes,1,opt,name=machineID"`
	// SystemUUID reported by the node. For unique machine identification
	// MachineID is preferred. This field is specific to Red Hat hosts
	// https://access.redhat.com/documentation/en-US/Red_Hat_Subscription_Management/1/html/RHSM/getting-system-uuid.html
	SystemUUID string `json:"systemUUID,omitempty" protobuf:"bytes,2,opt,name=systemUUID"`
	// Boot ID reported by the node.
	BootID string `json:"bootID,omitempty" protobuf:"bytes,3,opt,name=bootID"`
	// Kernel Version reported by the node.
	KernelVersion string `json:"kernelVersion,omitempty" protobuf:"bytes,4,opt,name=kernelVersion"`
	// OS Image reported by the node.
	OSImage string `json:"osImage,omitempty" protobuf:"bytes,5,opt,name=osImage"`
	// ContainerRuntime Version reported by the node.
	ContainerRuntimeVersion string `json:"containerRuntimeVersion,omitempty" protobuf:"bytes,6,opt,name=containerRuntimeVersion"`
	// Kubelet Version reported by the node.
	KubeletVersion string `json:"kubeletVersion,omitempty" protobuf:"bytes,7,opt,name=kubeletVersion"`
	// KubeProxy Version reported by the node.
	KubeProxyVersion string `json:"kubeProxyVersion,omitempty" protobuf:"bytes,8,opt,name=kubeProxyVersion"`
	// The Operating System reported by the node
	OperatingSystem string `json:"operatingSystem,omitempty" protobuf:"bytes,9,opt,name=operatingSystem"`
	// The Architecture reported by the node
	Architecture string `json:"architecture,omitempty" protobuf:"bytes,10,opt,name=architecture"`
}

// MachineAddress contains information for the machine's address.
type MachineAddress struct {
	// Machine address type, one of Public, ExternalIP or InternalIP.
	Type MachineAddressType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=MachineAddressType"`
	// The machine address.
	Address string `json:"address" protobuf:"bytes,2,opt,name=address"`
}

// MachineAddressType represents the type of machine address.
type MachineAddressType string

// These are valid address type of machine.
const (
	MachineHostName    MachineAddressType = "Hostname"
	MachineExternalIP  MachineAddressType = "ExternalIP"
	MachineInternalIP  MachineAddressType = "InternalIP"
	MachineExternalDNS MachineAddressType = "ExternalDNS"
	MachineInternalDNS MachineAddressType = "InternalDNS"
)

// MachineCondition contains details for the current condition of this Machine.
type MachineCondition struct {
	// Type is the type of the condition.
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// Last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// MachinePhase defines the phase of machine constructor
type MachinePhase string

const (
	// MachineRunning is the normal running phase
	MachineRunning MachinePhase = "Running"
	// MachineInitializing is the initialize phase
	MachineInitializing MachinePhase = "Initializing"
	// MachineFailed is the failed phase
	MachineFailed MachinePhase = "Failed"
	// MachineTerminating is the terminating phase
	MachineTerminating MachinePhase = "Terminating"
)

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineList is the whole list of all machine in an cluster.
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of clusters
	Items []Machine `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronHPAProxyOptions is the query options to a kube-apiserver proxy call.
type CronHPAProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronHPA is a new kubernetes workload.
type CronHPA struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of CronHPA.
	// +optional
	Spec CronHPASpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status CronHPAStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronHPAList is the whole list of all CronHPAs which owned by a tenant.
type CronHPAList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of CronHPAs
	Items []CronHPA `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// CronHPASpec describes the attributes on a CronHPA.
type CronHPASpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// CronHPAStatus is information about the current status of a CronHPA.
type CronHPAStatus struct {
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,1,opt,name=version"`
	// Phase is the current lifecycle phase of the CronHPA of cluster.
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LBCFProxyOptions is the query options to a kube-apiserver proxy call.
type LBCFProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	Name      string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
	Action    string `json:"action,omitempty" protobuf:"bytes,3,opt,name=action"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LBCF is a kubernetes load balancer manager.
type LBCF struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec LBCFSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status LBCFStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LBCFList is the whole list of all helms which owned by a tenant.
type LBCFList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of LBCFs
	Items []LBCF `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// LBCFSpec describes the attributes on a Helm.
type LBCFSpec struct {
	TenantID    string `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	ClusterName string `json:"clusterName" protobuf:"bytes,2,opt,name=clusterName"`
	Version     string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
}

// LBCFStatus is information about the current status of a Helm.
type LBCFStatus struct {
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
}
