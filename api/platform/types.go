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
	"fmt"
	"math/rand"
	"net"
	"os"
	"path"
	"strings"

	pkgerrors "github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	applicationv1 "tkestack.io/tke/api/application/v1"
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
	Taints     []corev1.Taint
	Proxy      ClusterMachineProxy
}

// ClusterMachine is the proxy definition of ClusterMachine.
type ClusterMachineProxy struct {
	Type       ProxyType
	IP         string
	Port       int32
	Username   string
	Password   []byte
	PrivateKey []byte
	PassPhrase []byte
}

// ProxyType describes diffirent type of proxy
type ProxyType string

const (
	// SSH jumper server proxy
	SSHJumpServer ProxyType = "SSHJumpServer"
	// SOCKS5 proxy
	SOCKS5 ProxyType = "SOCKS5"
)

const (
	// RegistrationCommandAnno contains base64 registration command of cluster net
	RegistrationCommandAnno = "tkestack.io/registration-command"
	// AnywhereEdtionLabel describe which anywhere edition will be deployed
	AnywhereEdtionLabel = "tkestack.io/anywhere-edtion"
	// AnywhereSubscriptionNameAnno describe sub name
	AnywhereSubscriptionNameAnno = "tkestack.io/anywhere-subscription-name"
	// AnywhereSubscriptionNameAnno describe sub namespace
	AnywhereSubscriptionNamespaceAnno = "tkestack.io/anywhere-subscription-namespace"
	// AnywhereLocalizationsAnno contains base64 localizations json data
	AnywhereLocalizationsAnno = "tkestack.io/anywhere-localizations"
	// AnywhereMachinesAnno contains base64 machines json data
	AnywhereMachinesAnno = "tkestack.io/anywhere-machines"
	// AnywhereUpgradeRetryComponentAnno describe curent retry component when upgrade failed
	AnywhereUpgradeRetryComponentAnno = "tkestack.io/anywhere-upgrade-retry-component"
	// AnywhereUpgradeRetryComponentAnno describe anywhere upgrade stats
	AnywhereUpgradeStatsAnno = "tkestack.io/anywhere-upgrade-stats"
	// ClusterNameLable contains related cluster's name for no-cluster resources
	ClusterNameLable = "tkestack.io/cluster-name"
	// HubAPIServerAnno describe hub cluster api server url
	HubAPIServerAnno = "tkestack.io/hub-api-server"
	// cluster credential token
	CredentialTokenAnno = "tkestack.io/credential-token"
	// AnywhereApplicationAnno contains base64 application json data
	AnywhereApplicationAnno = "tkestack.io/anywhere-application"
	// AnywhereValidateAnno is exist, the cluster will always return validate result
	AnywhereValidateAnno = "tkestack.io/anywhere-validate"
	// LocationBasedImagePrefixAnno is exist, the cluster will use it as k8s images prefix
	LocationBaseImagePrefixAnno = "tkestack.io/location-based-image-prefix"
)

// KubeVendorType describe the kubernetes provider of the cluster
// ref https://github.com/open-cluster-management/multicloud-operators-foundation/blob/e94b719de6d5f3541e948dd70ad8f1ff748aa452/pkg/apis/internal.open-cluster-management.io/v1beta1/clusterinfo_types.go#L137
type KubeVendorType string

const (
	// KubeVendorTKE TKE
	KubeVendorTKE KubeVendorType = "TKE"
	// KubeVendorOpenShift OpenShift
	KubeVendorOpenShift KubeVendorType = "OpenShift"
	// KubeVendorAKS Azure Kuberentes Service
	KubeVendorAKS KubeVendorType = "AKS"
	// KubeVendorEKS Elastic Kubernetes Service
	KubeVendorEKS KubeVendorType = "EKS"
	// KubeVendorGKE Google Kubernetes Engine
	KubeVendorGKE KubeVendorType = "GKE"
	// KubeVendorICP IBM Cloud Private
	KubeVendorICP KubeVendorType = "ICP"
	// KubeVendorIKS IBM Kubernetes Service
	KubeVendorIKS KubeVendorType = "IKS"
	// KubeVendorOSD OpenShiftDedicated
	KubeVendorOSD KubeVendorType = "OpenShiftDedicated"
	// KubeVendorOther other (unable to auto detect)
	KubeVendorOther KubeVendorType = "Other"
)

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
	// ServiceCIDR is used to set a separated CIDR for k8s service, it's exclusive with MaxClusterServiceNum.
	// +optional
	ServiceCIDR *string
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
	ScalingMachines []ClusterMachine
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

	// ClusterCredentialRef for isolate sensitive information.
	// If not specified, cluster controller will create one;
	// If specified, provider must make sure is valid.
	// +optional
	ClusterCredentialRef *corev1.LocalObjectReference

	// Etcd holds configuration for etcd.
	// +optional
	Etcd *Etcd
	// If true will use hostname as nodename, if false will use machine IP as nodename.
	// +optional
	HostnameAsNodename bool
	// +optional
	NetworkArgs map[string]string
	// BootstrapApps will install apps during creating cluster
	// +optional
	BootstrapApps BootstrapApps
	// AppVersion is the overall version of system components
	// +optional
	AppVersion string
	// ClusterLevel is the expect level of cluster
	// +optional
	ClusterLevel *string
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
	// +optional
	ClusterCIDR string
	// +optional
	SecondaryServiceCIDR string
	// +optional
	SecondaryClusterCIDR string
	// +optional
	NodeCIDRMaskSizeIPv4 int32
	// +optional
	NodeCIDRMaskSizeIPv6 int32
	// +optional
	KubeVendor KubeVendorType
	// AppVersion is the overall version of system components
	// +optional
	AppVersion string
	// ComponentPhase is the status of components, contains "deployed", "pending-upgrade", "failed" status
	// +optional
	ComponentPhase ComponentPhase
	// ClusterLevel is the real level of cluster
	// +optional
	ClusterLevel *string
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

// GPUType defines the gpu type of cluster.
type GPUType string

const (
	// GPUPhysical indicates the gpu type of cluster is physical.
	GPUPhysical GPUType = "Physical"
	// GPUVirtual indicates the gpu type of cluster is virtual.
	GPUVirtual GPUType = "Virtual"
)

type ContainerRuntimeType = string

const (
	Containerd ContainerRuntimeType = "containerd"
	Docker     ContainerRuntimeType = "docker"
)

// ClusterPhase defines the phase of cluster constructor.
type ClusterPhase string

const (
	// ClusterInitializing is the initialize phase.
	ClusterInitializing ClusterPhase = "Initializing"
	// ClusterWaiting indicates that the cluster is waiting for registration.
	ClusterWaiting ClusterPhase = "Waiting"
	// ClusterRunning is the normal running phase.
	ClusterRunning ClusterPhase = "Running"
	// ClusterFailed is the failed phase.
	ClusterFailed ClusterPhase = "Failed"
	// ClusterUpgrading means that the cluster is in upgrading process.
	ClusterUpgrading ClusterPhase = "Upgrading"
	// ClusterTerminating means the cluster is undergoing graceful termination.
	ClusterTerminating ClusterPhase = "Terminating"
	// ClusterUpscaling means the cluster is undergoing graceful up scaling.
	ClusterUpscaling ClusterPhase = "Upscaling"
	// ClusterDownscaling means the cluster is undergoing graceful down scaling.
	ClusterDownscaling ClusterPhase = "Downscaling"
	// ClusterRecovering means the cluster is recovering form confined.
	ClusterRecovering ClusterPhase = "Recovering"
)

// ComponentPhase defines the phase of anywhere cluster component
type ComponentPhase string

const (
	// ComponentDeployed is the normal phase of anywhere cluster component
	ComponentDeployed ComponentPhase = "deployed"
	// ComponentPendingUpgrade means the anywhere cluster component is upgrading
	ComponentPendingUpgrade ComponentPhase = "pending-upgrade"
	// ComponentFailed means the anywhere cluster component upgrade failed
	ComponentFailed ComponentPhase = "failed"
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
	Path string
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
	// Username is the username for basic authentication to the kubernetes cluster.
	// +optional
	Username string
	// Impersonate is the username to act-as.
	// +optional
	Impersonate string
	// ImpersonateGroups is the groups to imperonate.
	// +optional
	ImpersonateGroups []string
	// ImpersonateUserExtra contains additional information for impersonated user.
	// +optional
	ImpersonateUserExtra ImpersonateUserExtra
}

type ImpersonateUserExtra map[string]string

func (i ImpersonateUserExtra) ExtraToHeaders() map[string][]string {
	res := map[string][]string{}
	for k, v := range i {
		res[k] = strings.Split(v, ",")
	}
	return res
}

func (cc ClusterCredential) RESTConfig(cls *Cluster) *rest.Config {
	config := &rest.Config{}
	if cls != nil {
		host := clusterHost(cls)
		if len(host) != 0 {
			config.Host = fmt.Sprintf("https://%s", host)
		}
	}
	// If api-server does not sign the ip in address, set ca then request, it will report x509 certificate error, need to ignore the certificate
	if os.Getenv("TKE_IGNORE_CA") != "true" && cc.CACert != nil {
		config.TLSClientConfig.CAData = cc.CACert
	} else {
		config.TLSClientConfig.Insecure = true
	}
	if cc.ClientCert != nil && cc.ClientKey != nil {
		config.TLSClientConfig.CertData = cc.ClientCert
		config.TLSClientConfig.KeyData = cc.ClientKey
	}
	if cc.Token != nil {
		config.BearerToken = *cc.Token
	}

	config.Impersonate.UserName = cc.Impersonate
	config.Impersonate.Groups = cc.ImpersonateGroups
	config.Impersonate.Extra = cc.ImpersonateUserExtra.ExtraToHeaders()

	return config
}

func clusterHost(cluster *Cluster) string {
	address, err := clusterAddress(cluster)
	if err != nil {
		return ""
	}

	result := net.JoinHostPort(address.Host, fmt.Sprintf("%d", address.Port))
	if address.Path != "" {
		result = path.Join(result, address.Path)
	}

	return result
}

func clusterAddress(cluster *Cluster) (*ClusterAddress, error) {
	addrs := make(map[AddressType][]ClusterAddress)
	for _, one := range cluster.Status.Addresses {
		addrs[one.Type] = append(addrs[one.Type], one)
	}

	var address *ClusterAddress
	if len(addrs[AddressInternal]) != 0 {
		address = &addrs[AddressInternal][rand.Intn(len(addrs[AddressInternal]))]
	} else if len(addrs[AddressAdvertise]) != 0 {
		address = &addrs[AddressAdvertise][rand.Intn(len(addrs[AddressAdvertise]))]
	} else {
		if len(addrs[AddressReal]) != 0 {
			address = &addrs[AddressReal][rand.Intn(len(addrs[AddressReal]))]
		}
	}
	if address == nil {
		return nil, pkgerrors.New("no valid address for the cluster")
	}

	return address, nil
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
	// +optional
	CSIOperator *CSIOperatorFeature
	// For kube-apiserver authorization webhook
	// +optional
	AuthzWebhookAddr *AuthzWebhookAddr
	// +optional
	EnableMetricsServer bool
	// +optional
	EnableCilium bool

	ContainerRuntime ContainerRuntimeType
	// +optional
	IPv6DualStack bool
	// Upgrade control upgrade process.
	// +optional
	Upgrade Upgrade
}

type BootstrapApps []BootstapApp

type BootstapApp struct {
	App App
}

type App struct {
	// +optional
	metav1.ObjectMeta
	// +optional
	Spec applicationv1.AppSpec
}

type HA struct {
	TKEHA        *TKEHA
	ThirdPartyHA *ThirdPartyHA
}

type TKEHA struct {
	VIP  string
	VRID *int32
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

type CSIOperatorFeature struct {
	Version string
}

type AuthzWebhookAddr struct {
	// +optional
	Builtin *BuiltinAuthzWebhookAddr
	// +optional
	External *ExternalAuthzWebhookAddr
}

type BuiltinAuthzWebhookAddr struct{}

type ExternalAuthzWebhookAddr struct {
	IP   string `json:"ip" protobuf:"bytes,1,name=ip"`
	Port int32  `json:"port" protobuf:"varint,2,name=port"`
}

const (
	// node lifecycle hook
	HookPreInstall  HookType = "PreInstall"
	HookPostInstall HookType = "PostInstall"
	HookPreUpgrade  HookType = "PreUpgrade"
	HookPostUpgrade HookType = "PostUpgrade"

	// cluster lifecycle hook
	HookPreClusterInstall  HookType = "PreClusterInstall"
	HookPostClusterInstall HookType = "PostClusterInstall"
	HookPreClusterUpgrade  HookType = "PreClusterUpgrade"
	HookPostClusterUpgrade HookType = "PostClusterUpgrade"
	HookPreClusterDelete   HookType = "PreClusterDelete"
	HookPostClusterDelete  HookType = "PostClusterDelete"
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

// Etcd contains elements describing Etcd configuration.
type Etcd struct {

	// Local provides configuration knobs for configuring the local etcd instance
	// Local and External are mutually exclusive
	Local *LocalEtcd

	// External describes how to connect to an external etcd cluster
	// Local and External are mutually exclusive
	External *ExternalEtcd
}

// LocalEtcd describes that kubeadm should run an etcd cluster locally
type LocalEtcd struct {
	// DataDir is the directory etcd will place its data.
	// Defaults to "/var/lib/etcd".
	DataDir string

	// ExtraArgs are extra arguments provided to the etcd binary
	// when run inside a static pod.
	ExtraArgs map[string]string

	// ServerCertSANs sets extra Subject Alternative Names for the etcd server signing cert.
	ServerCertSANs []string
	// PeerCertSANs sets extra Subject Alternative Names for the etcd peer signing cert.
	PeerCertSANs []string
}

// ExternalEtcd describes an external etcd cluster
type ExternalEtcd struct {

	// Endpoints of etcd members. Useful for using external etcd.
	// If not provided, kubeadm will run etcd in a static pod.
	Endpoints []string
	// CAFile is an SSL Certificate Authority file used to secure etcd communication.
	CAFile string
	// CertFile is an SSL certification file used to secure etcd communication.
	CertFile string
	// KeyFile is an SSL key file used to secure etcd communication.
	KeyFile string
}

type Upgrade struct {
	// Upgrade mode, default value is Auto.
	Mode UpgradeMode
	// Upgrade strategy config.
	Strategy UpgradeStrategy
}

type UpgradeMode string

const (
	// Upgrade nodes automatically.
	UpgradeModeAuto = UpgradeMode("Auto")
	// Manual upgrade nodes which means user need label node with `platform.tkestack.io/need-upgrade`.
	UpgradeModeManual = UpgradeMode("Manual")
)

// UpgradeStrategy used to control the upgrade process.
type UpgradeStrategy struct {
	// The maximum number of pods that can be unready during the upgrade.
	// 0% means all pods need to be ready after evition.
	// 100% means ignore any pods unready which may be used in one worker node, use this carefully!
	// default value is 0%.
	MaxUnready intstr.IntOrString
	// Whether drain node before upgrade.
	// Draining node before upgrade is recommended.
	// But not all pod running as cows, a few running as pets.
	// If your pod can not accept be expelled from current node, this value should be false.
	DrainNodeBeforeUpgrade bool
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
	// Addon type, one of PersistentEvent or LogCollector etc.
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
	// Addon type, one of PersistentEvent or LogCollector etc.
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

// +k8s:conversion-gen:explicit-from=net/url.Values
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
	IP          string
	Port        int32
	Scheme      string
	IndexName   string
	User        string
	Password    string
	ReserveDays int32
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProxyOptions is the query options to a proxy call.
type ProxyOptions struct {
	metav1.TypeMeta

	// Path is the URL path to use for the current proxy request to helm-api.
	// +optional
	Path string
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

// +k8s:conversion-gen:explicit-from=net/url.Values
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

// +k8s:conversion-gen:explicit-from=net/url.Values
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
	Taints      []corev1.Taint
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
	// MachineInitializing is the initialize phase
	MachineInitializing MachinePhase = "Initializing"
	// MachineRunning is the normal running phase
	MachineRunning MachinePhase = "Running"
	// MachineFailed is the failed phase
	MachineFailed MachinePhase = "Failed"
	// MachineUpgrading means that the machine is in upgrading process.
	MachineUpgrading MachinePhase = "Upgrading"
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

// +k8s:conversion-gen:explicit-from=net/url.Values
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

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterGroupAPIResourceItemsList is the whole list of all ClusterAPIResource.
type ClusterGroupAPIResourceItemsList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	// List of ClusterAPIResource
	Items []ClusterGroupAPIResourceItems
	// Failed Group Error
	FailedGroupError string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=list,get
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterGroupAPIResourceItems contains the GKV for the current kubernetes cluster
type ClusterGroupAPIResourceItems struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta
	// groupVersion is the group and version this APIResourceList is for.
	GroupVersion string
	// resources contains the name of the resources and if they are namespaced.
	APIResources []ClusterGroupAPIResourceItem
}

// ClusterGroupAPIResourceItem specifies the name of a resource and whether it is namespaced.
type ClusterGroupAPIResourceItem struct {
	// name is the plural name of the resource.
	Name string
	// singularName is the singular name of the resource.  This allows clients to handle plural and singular opaquely.
	// The singularName is more correct for reporting status on a single item and both singular and plural are allowed
	// from the kubectl CLI interface.
	SingularName string
	// namespaced indicates if a resource is namespaced or not.
	Namespaced bool
	// group is the preferred group of the resource.  Empty implies the group of the containing resource list.
	// For subresources, this may have a different value, for example: Scale".
	Group string
	// version is the preferred version of the resource.  Empty implies the version of the containing resource list
	// For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource's group)".
	Version string
	// kind is the kind for the resource (e.g. 'Foo' is the kind for a resource 'foo')
	Kind string
	// verbs is a list of supported kube verbs (this includes get, list, watch, create,
	// update, patch, delete, deletecollection, and proxy)
	Verbs []string
	// shortNames is a list of suggested short names of the resource.
	ShortNames []string
	// categories is a list of the grouped resources this resource belongs to (e.g. 'all')
	Categories []string
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterGroupAPIResourceOptions is the query options.
type ClusterGroupAPIResourceOptions struct {
	metav1.TypeMeta
}
