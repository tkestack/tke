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

package platform

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is a Kubernetes cluster in TKE.
type Cluster struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec ClusterSpec
	// +optional
	Status ClusterStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList is the whole list of all clusters which owned by a tenant.
type ClusterList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of clusters
	Items []Cluster
}

// ClusterMachine is the master machine definition of cluster.
type ClusterMachine struct {
	IP         string
	Port       int32
	Username   string
	Password   []byte
	PrivateKey []byte
	PassPhrase []byte
	Labels     map[string]string
}

// ClusterSpec is a description of a cluster.
type ClusterSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers []FinalizerName
	TenantID   string
	// +optional
	DisplayName string
	Type        string
	Version     string
	// +optional
	NetworkType NetworkType
	// +optional
	NetworkDevice string
	// +optional
	ClusterCIDR string
	// +optional
	// DNSDomain is the dns domain used by k8s services. Defaults to "cluster.local".
	DNSDomain string
	// +optional
	PublicAlternativeNames []string
	// +optional
	Features ClusterFeature
	// +optional
	Properties ClusterProperty
	// +optional
	Machines []ClusterMachine

	// +optional
	DockerExtraArgs map[string]string
	// +optional
	KubeletExtraArgs map[string]string
	// +optional
	APIServerExtraArgs map[string]string
	// +optional
	ControllerManagerExtraArgs map[string]string
	// +optional
	SchedulerExtraArgs map[string]string
}

// ClusterStatus represents information about the status of a cluster.
type ClusterStatus struct {
	// +optional
	Locked *bool
	// +optional
	Version string
	// +optional
	Phase ClusterPhase
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ClusterCondition
	// A human readable message indicating details about why the cluster is in this condition.
	// +optional
	Message string
	// A brief CamelCase message indicating details about why the cluster is in this state.
	// +optional
	Reason string
	// List of addresses reachable to the cluster.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Addresses []ClusterAddress
	// +optional
	Resource ClusterResource
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Components []ClusterComponent
	// +optional
	ServiceCIDR string
	// +optional
	NodeCIDRMaskSize int32
	// +optional
	DNSIP string
	// +optional
	RegistryIPs []string
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
	Type string
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status ConditionStatus
	// Last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string
	// Human-readable message indicating details about last transition.
	// +optional
	Message string
}

// AddressType indicates the type of cluster apiserver access address.
type AddressType string

// These are valid address type of node.
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
	Type AddressType
	// The cluster address.
	Host string
	Port int32
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCredential records the credential information needed to access the cluster.
type ClusterCredential struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	TenantID    string
	ClusterName string

	// For TKE in global reuse
	// +optional
	ETCDCACert []byte
	// +optional
	ETCDCAKey []byte
	// +optional
	ETCDAPIClientCert []byte
	// +optional
	ETCDAPIClientKey []byte

	// For validate the server cert
	// +optional
	CACert []byte
	// +optional
	CAKey []byte
	// For kube-apiserver X509 auth
	// +optional
	ClientCert []byte
	// For kube-apiserver X509 auth
	// +optional
	ClientKey []byte
	// For kube-apiserver token auth
	// +optional
	Token *string
	// For kubeadm init or join
	// +optional
	BootstrapToken *string
	// For kubeadm init or join
	// +optional
	CertificateKey *string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterCredentialList is the whole list of all ClusterCredential which owned by a tenant.
type ClusterCredentialList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of clusters
	Items []ClusterCredential
}

// ClusterFeature records the features that are enabled by the cluster.
type ClusterFeature struct {
	// +optional
	IPVS *bool
	// +optional
	PublicLB *bool
	// +optional
	InternalLB *bool
	// +optional
	GPUType *GPUType
	// +optional
	EnableMasterSchedule bool
	// +optional
	HA *HA
	// +optional
	SkipConditions []string
	// +optional
	Files []File
	// +optional
	Hooks map[HookType]string
}

type HA struct {
	TKEHA        *TKEHA
	ThirdPartyHA *ThirdPartyHA
}

type TKEHA struct {
	VIP string
}

type ThirdPartyHA struct {
	VIP   string
	VPort int32
}

type File struct {
	Src string // Only support regular file
	Dst string
}

type HookType string

const (
	HookPreInstall  HookType = "PreInstall"
	HookPostInstall HookType = "PostInstall"
)

// ClusterProperty records the attribute information of the cluster.
type ClusterProperty struct {
	// +optional
	MaxClusterServiceNum *int32
	// +optional
	MaxNodePodNum *int32
	// +optional
	OversoldRatio map[string]string
}

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[string]resource.Quantity

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	Limits   ResourceList
	Requests ResourceList
}

// ClusterResource records the current available and maximum resource quota
// information for the cluster.
type ClusterResource struct {
	// Capacity represents the total resources of a cluster.
	// +optional
	Capacity ResourceList
	// Allocatable represents the resources of a cluster that are available for scheduling.
	// Defaults to Capacity.
	// +optional
	Allocatable ResourceList
	// +optional
	Allocated ResourceList
}

// ClusterComponent records the number of copies of each component of the
// cluster master.
type ClusterComponent struct {
	Type     string
	Replicas ClusterComponentReplicas
}

// ClusterComponentReplicas records the number of copies of each state of each
// component of the cluster master.
type ClusterComponentReplicas struct {
	Desired   int32
	Current   int32
	Available int32
	Updated   int32
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
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// Spec defines the desired identities of addons in this set.
	// +optional
	Spec ClusterAddonSpec
	// +optional
	Status ClusterAddonStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAddonList is the whole list of all ClusterAddon.
type ClusterAddonList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	// List of ClusterAddon
	Items []ClusterAddon
}

// ClusterAddonSpec indicates the specifications of the ClusterAddon.
type ClusterAddonSpec struct {
	// Addon type, one of Helm, PersistentEvent or LogCollector etc.
	Type string
	// AddonLevel is level of cluster addon.
	Level AddonLevel
	// Version
	Version string
}

// ClusterAddonStatus is information about the current status of a ClusterAddon.
type ClusterAddonStatus struct {
	// +optional
	Version string
	// Phase is the current lifecycle phase of the addon of cluster.
	// +optional
	Phase string
	// Reason is a brief CamelCase string that describes any failure.
	// +optional
	Reason string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=list
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAddonType records the all addons of cluster available.
type ClusterAddonType struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// Addon type, one of Helm, PersistentEvent or LogCollector etc.
	Type string
	// AddonLevel is level of cluster addon.
	Level AddonLevel
	// LatestVersion is latest version of the addon.
	LatestVersion string
	// Description is desc of the addon.
	Description           string
	CompatibleClusterType []string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAddonTypeList is a resource containing a list of ClusterAddonType objects.
type ClusterAddonTypeList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	// +optional
	Items []ClusterAddonType
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterApplyOptions is the query options to a kube-apiserver proxy call for cluster object.
type ClusterApplyOptions struct {
	metav1.TypeMeta
	// +optional
	NotUpdate bool
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Registry records the third-party image repository information stored by the
// user.
type Registry struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// +optional
	Spec RegistrySpec
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RegistryList is a resource containing a list of Registry objects.
type RegistryList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	// +optional
	Items []Registry
}

// RegistrySpec indicates the specifications of the third-party image repository.
type RegistrySpec struct {
	// +optional
	TenantID string
	// +optional
	DisplayName string
	// +optional
	ClusterName string
	// +optional
	URL string
	// +optional
	UserName *string
	// +optional
	Password *string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PersistentEvent is a recorder of kubernetes event.
type PersistentEvent struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec PersistentEventSpec
	// +optional
	Status PersistentEventStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PersistentEventList is the whole list of all clusters which owned by a tenant.
type PersistentEventList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of PersistentEvents
	Items []PersistentEvent
}

// PersistentEventSpec describes the attributes on a PersistentEvent.
type PersistentEventSpec struct {
	TenantID          string
	ClusterName       string
	PersistentBackEnd PersistentBackEnd
	Version           string
}

// PersistentEventStatus is information about the current status of a
// PersistentEvent.
type PersistentEventStatus struct {
	// +optional
	Version string
	// Phase is the current lifecycle phase of the persistent event of cluster.
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
}

// PersistentBackEnd indicates the backend type and attributes of the persistent
// log store.
type PersistentBackEnd struct {
	CLS *StorageBackEndCLS
	ES  *StorageBackEndES
}

// StorageBackEndCLS records the attributes required when the backend storage
// type is CLS.
type StorageBackEndCLS struct {
	LogSetID string
	TopicID  string
}

// StorageBackEndES records the attributes required when the backend storage
// type is ElasticSearch.
type StorageBackEndES struct {
	IP        string
	Port      int32
	Scheme    string
	IndexName string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmProxyOptions is the query options to a Helm-api proxy call.
type HelmProxyOptions struct {
	metav1.TypeMeta

	// Path is the URL path to use for the current proxy request to helm-api.
	// +optional
	Path string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Helm is a kubernetes package manager.
type Helm struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec HelmSpec
	// +optional
	Status HelmStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmList is the whole list of all helms which owned by a tenant.
type HelmList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of Helms
	Items []Helm
}

// HelmSpec describes the attributes on a Helm.
type HelmSpec struct {
	TenantID    string
	ClusterName string
	Version     string
}

// HelmStatus is information about the current status of a Helm.
type HelmStatus struct {
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
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Prometheus is a kubernetes package manager.
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPAMProxyOptions is the query options to a Helm-api proxy call.
type IPAMProxyOptions struct {
	metav1.TypeMeta

	// Path is the URL path to use for the current proxy request to helm-api.
	// +optional
	Path string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPAM is a scheduler plugin for assigning IP.
type IPAM struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of clusters in this set.
	// +optional
	Spec IPAMSpec
	// +optional
	Status IPAMStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPAMList is the whole list of all IPAMs which owned by a tenant.
type IPAMList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of IPAMs
	Items []IPAM
}

// IPAMSpec describes the attributes on a IPAM.
type IPAMSpec struct {
	TenantID    string
	ClusterName string
	Version     string
}

// IPAMStatus is information about the current status of a IPAM.
type IPAMStatus struct {
	// +optional
	Version string
	// Phase is the current lifecycle phase of the addon of cluster.
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
}

// +genclient
// +genclient:nonNamespaced
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
	TenantID    string
	ClusterName string
	Version     string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TappControllerProxyOptions is the query options to a kube-apiserver proxy call.
type TappControllerProxyOptions struct {
	metav1.TypeMeta

	Namespace string
	Name      string
	Action    string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TappController is a new kubernetes workload.
type TappController struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of tapp controller.
	// +optional
	Spec TappControllerSpec
	// +optional
	Status TappControllerStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TappControllerList is the whole list of all tapp controllers which owned by a tenant.
type TappControllerList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of tapp controllers
	Items []TappController
}

// TappControllerSpec describes the attributes on a tapp controller.
type TappControllerSpec struct {
	TenantID    string
	ClusterName string
	Version     string
}

// TappControllerStatus is information about the current status of a tapp controller.
type TappControllerStatus struct {
	// +optional
	Version string
	// Phase is the current lifecycle phase of the tapp controller of cluster.
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
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CSIProxyOptions is the query options to a kube-apiserver proxy call for CSI crd object.
type CSIProxyOptions struct {
	metav1.TypeMeta

	Namespace string
	Name      string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CSIOperator is a operator to manages CSI external components.
type CSIOperator struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of storage operator.
	// +optional
	Spec CSIOperatorSpec
	// +optional
	Status CSIOperatorStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CSIOperatorList is the whole list of all storage operators which owned by a tenant.
type CSIOperatorList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of storage operators.
	Items []CSIOperator
}

// CSIOperatorSpec describes the attributes of a storage operator.
type CSIOperatorSpec struct {
	TenantID    string
	ClusterName string
	// Version of the CSI operator.
	Version string
}

// CSIOperatorStatus is information about the current status of a storage operator.
type CSIOperatorStatus struct {
	// +optional
	Version string
	// StorageVendorVersion will be set to the config version of the storage vendor.
	// +optional
	StorageVendorVersion string
	// Phase is the current lifecycle phase of the csi operator of cluster.
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
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PVCRProxyOptions is the query options to a kube-apiserver proxy call for PVCR crd object.
type PVCRProxyOptions struct {
	metav1.TypeMeta

	Namespace string
	Name      string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeDecorator is a controller to manage PVC information.
type VolumeDecorator struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of volume decorator.
	// +optional
	Spec VolumeDecoratorSpec
	// +optional
	Status VolumeDecoratorStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeDecoratorList is the whole list of all VolumeDecorator which owned by a tenant.
type VolumeDecoratorList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of volume decorators.
	Items []VolumeDecorator
}

// VolumeDecoratorSpec describes the attributes of a VolumeDecorator.
type VolumeDecoratorSpec struct {
	TenantID          string
	ClusterName       string
	Version           string
	VolumeTypes       []string
	WorkloadAdmission bool
}

// VolumeDecoratorStatus is information about the current status of a VolumeDecorator.
type VolumeDecoratorStatus struct {
	// +optional
	Version string
	// VolumeTypes is the supported volume types in this cluster.
	// +optional
	VolumeTypes []string
	// WorkloadAdmission will be true to enable the workload admission webhook.
	// +optional
	WorkloadAdmission bool
	// StorageVendorVersion will be set to the config version of the storage vendor.
	// +optional
	StorageVendorVersion string
	// Phase is the current lifecycle phase of the volume-decorator of cluster.
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
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogCollectorProxyOptions is the query options to a kube-apiserver proxy call for LogCollector crd object.
type LogCollectorProxyOptions struct {
	metav1.TypeMeta

	Namespace string
	Name      string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogCollector is a manager to collect logs of workload.
type LogCollector struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of LogCollector.
	// +optional
	Spec LogCollectorSpec
	// +optional
	Status LogCollectorStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogCollectorList is the whole list of all LogCollector which owned by a tenant.
type LogCollectorList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of volume decorators.
	Items []LogCollector
}

// LogCollectorSpec describes the attributes of a LogCollector.
type LogCollectorSpec struct {
	TenantID    string
	ClusterName string
	Version     string
}

// LogCollectorStatus is information about the current status of a LogCollector.
type LogCollectorStatus struct {
	// +optional
	Version string
	// Phase is the current lifecycle phase of the LogCollector of cluster.
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
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Machine instance in Kubernetes cluster
type Machine struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of the Machine.
	// +optional
	Spec MachineSpec
	// +optional
	Status MachineStatus
}

// MachineSpec is a description of machine.
type MachineSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// +optional
	Finalizers  []FinalizerName
	TenantID    string
	ClusterName string
	Type        string
	IP          string
	Port        int32
	Username    string
	Password    []byte
	PrivateKey  []byte
	PassPhrase  []byte
	Labels      map[string]string
}

// MachineStatus represents information about the status of an machine.
type MachineStatus struct {
	// +optional
	Locked *bool

	// +optional
	Phase MachinePhase
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []MachineCondition
	// A human readable message indicating details about why the machine is in this condition.
	// +optional
	Message string
	// A brief CamelCase message indicating details about why the machine is in this state.
	// +optional
	Reason string
	// List of addresses reachable to the machine.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Addresses []MachineAddress
	// Set of ids/uuids to uniquely identify the node.
	// +optional
	MachineInfo MachineSystemInfo
}

// MachineSystemInfo is a set of ids/uuids to uniquely identify the node.
type MachineSystemInfo struct {
	// MachineID reported by the node. For unique machine identification
	// in the cluster this field is preferred. Learn more from man(5)
	// machine-id: http://man7.org/linux/man-pages/man5/machine-id.5.html
	MachineID string
	// SystemUUID reported by the node. For unique machine identification
	// MachineID is preferred. This field is specific to Red Hat hosts
	// https://access.redhat.com/documentation/en-US/Red_Hat_Subscription_Management/1/html/RHSM/getting-system-uuid.html
	SystemUUID string
	// Boot ID reported by the node.
	BootID string
	// Kernel Version reported by the node.
	KernelVersion string
	// OS Image reported by the node.
	OSImage string
	// ContainerRuntime Version reported by the node.
	ContainerRuntimeVersion string
	// Kubelet Version reported by the node.
	KubeletVersion string
	// KubeProxy Version reported by the node.
	KubeProxyVersion string
	// The Operating System reported by the node
	OperatingSystem string
	// The Architecture reported by the node
	Architecture string
}

// MachineAddress contains information for the machine's address.
type MachineAddress struct {
	// Machine address type, one of Public, ExternalIP or InternalIP.
	Type MachineAddressType
	// The machine address.
	Address string
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
	Type string
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status ConditionStatus
	// Last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string
	// Human-readable message indicating details about last transition.
	// +optional
	Message string
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
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	// List of clusters
	Items []Machine
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronHPAProxyOptions is the query options to a kube-apiserver proxy call.
type CronHPAProxyOptions struct {
	metav1.TypeMeta

	Namespace string
	Name      string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronHPA is a new kubernetes workload.
type CronHPA struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of CronHPA.
	// +optional
	Spec CronHPASpec
	// +optional
	Status CronHPAStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronHPAList is the whole list of all CronHPAs which owned by a tenant.
type CronHPAList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of CronHPAs
	Items []CronHPA
}

// CronHPASpec describes the attributes on a CronHPA.
type CronHPASpec struct {
	TenantID    string
	ClusterName string
	Version     string
}

// CronHPAStatus is information about the current status of a CronHPA.
type CronHPAStatus struct {
	// +optional
	Version string
	// Phase is the current lifecycle phase of the CronHPA of cluster.
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
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LBCFProxyOptions is the query options to a kube-apiserver proxy call.
type LBCFProxyOptions struct {
	metav1.TypeMeta

	Namespace string
	Name      string
	Action    string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LBCF is a kubernetes load balancer manager
type LBCF struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of LBCF.
	// +optional
	Spec LBCFSpec
	// +optional
	Status LBCFStatus
}

// LBCFSpec defines the desired identities of LBCF.
type LBCFSpec struct {
	TenantID    string
	ClusterName string
	Version     string
}

// LBCFStatus is information about the current status of a LBCF.
type LBCFStatus struct {
	// +optional
	Version string
	// Phase is the current lifecycle phase of the CronHPA of cluster.
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
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LBCFList is the whole list of all LBCF which owned by a tenant.
type LBCFList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	// List of CronHPAs
	Items []LBCF
}
