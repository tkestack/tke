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
	"fmt"
	"math/rand"
	"net"
	"os"
	"path"
	strings "strings"

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
	// If specified, the node's taints.
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty" protobuf:"bytes,8,opt,name=taints"`
	// +optional
	Proxy ClusterMachineProxy `json:"proxy,omitempty" protobuf:"bytes,9,opt,name=proxy"`
}

// ClusterMachine is the proxy definition of ClusterMachine.
type ClusterMachineProxy struct {
	Type ProxyType `json:"type" protobuf:"bytes,1,opt,name=type"`
	IP   string    `json:"ip" protobuf:"bytes,2,opt,name=ip"`
	Port int32     `json:"port" protobuf:"varint,3,opt,name=port"`
	// +optional
	Username string `json:"username,omitempty" protobuf:"bytes,4,opt,name=username"`
	// +optional
	Password []byte `json:"password,omitempty" protobuf:"bytes,5,opt,name=password"`
	// +optional
	PrivateKey []byte `json:"privateKey,omitempty" protobuf:"bytes,6,opt,name=privateKey"`
	// +optional
	PassPhrase []byte `json:"passPhrase,omitempty" protobuf:"bytes,7,opt,name=passPhrase"`
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
	LocationBasedImagePrefixAnno = "tkestack.io/location-based-image-prefix"
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
	// ServiceCIDR is used to set a separated CIDR for k8s service, it's exclusive with MaxClusterServiceNum.
	// +optional
	ServiceCIDR *string `json:"serviceCIDR,omitempty" protobuf:"bytes,19,opt,name=serviceCIDR"`
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

	// ClusterCredentialRef for isolate sensitive information.
	// If not specified, cluster controller will create one;
	// If specified, provider must make sure is valid.
	// +optional
	ClusterCredentialRef *corev1.LocalObjectReference `json:"clusterCredentialRef,omitempty" probobuf:"bytes,20,opt,name=clusterCredentialRef" protobuf:"bytes,20,opt,name=clusterCredentialRef"`

	// Etcd holds configuration for etcd.
	// +optional
	Etcd *Etcd `json:"etcd,omitempty" protobuf:"bytes,21,opt,name=etcd"`
	// If true will use hostname as nodename, if false will use machine IP as nodename.
	// +optional
	HostnameAsNodename bool `json:"hostnameAsNodename,omitempty" protobuf:"bytes,23,opt,name=hostnameAsNodename"`
	// +optional
	NetworkArgs map[string]string `json:"networkArgs,omitempty" protobuf:"bytes,24,name=networkArgs"`
	// +optional
	ScalingMachines []ClusterMachine `json:"scalingMachines,omitempty" protobuf:"bytes,25,opt,name=scalingMachines"`
	// BootstrapApps will install apps during creating cluster
	// +optional
	BootstrapApps BootstrapApps `json:"bootstrapApps,omitempty" protobuf:"bytes,26,opt,name=bootstrapApps"`
	// AppVersion is the overall version of system components
	// +optional
	AppVersion string `json:"appVersion,omitempty" protobuf:"bytes,27,opt,name=appVersion"`
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
	// +optional
	SecondaryServiceCIDR string `json:"secondaryServiceCIDR,omitempty" protobuf:"bytes,15,opt,name=secondaryServiceCIDR"`
	// +optional
	ClusterCIDR string `json:"clusterCIDR,omitempty" protobuf:"bytes,16,opt,name=clusterCIDR"`
	// +optional
	SecondaryClusterCIDR string `json:"secondaryClusterCIDR,omitempty" protobuf:"bytes,17,opt,name=secondaryClusterCIDR"`
	// +optional
	NodeCIDRMaskSizeIPv4 int32 `json:"nodeCIDRMaskSizeIPv4,omitempty" protobuf:"varint,18,opt,name=nodeCIDRMaskSizeIPv4"`
	// +optional
	NodeCIDRMaskSizeIPv6 int32 `json:"nodeCIDRMaskSizeIPv6,omitempty" protobuf:"varint,19,opt,name=nodeCIDRMaskSizeIPv6"`
	// +optional
	KubeVendor KubeVendorType `json:"kubeVendor" protobuf:"bytes,20,opt,name=kubeVendor"`
	// AppVersion is the overall version of system components
	// +optional
	AppVersion string `json:"appVersion,omitempty" protobuf:"bytes,21,opt,name=appVersion"`
	// ComponentPhase is the status of components, contains "deployed", "pending-upgrade", "failed" status
	// +optional
	ComponentPhase ComponentPhase `json:"componentPhase,omitempty" protobuf:"bytes,22,opt,name=componentPhase"`
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
	// ClusterConfined is the Confined phase.
	ClusterConfined ClusterPhase = "Confined"
	// ClusterIdling is the Idling phase.
	ClusterIdling ClusterPhase = "Idling"
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
	Path string `json:"path" protobuf:"bytes,4,opt,name=path"`
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
	// Username is the username for basic authentication to the kubernetes cluster.
	// +optional
	Username string `json:"username,omitempty" protobuf:"bytes,15,opt,name=username"`
	// Impersonate is the username to act-as.
	// +optional
	Impersonate string `json:"as,omitempty" protobuf:"bytes,16,opt,name=as"`
	// ImpersonateGroups is the groups to imperonate.
	// +optional
	ImpersonateGroups []string `json:"as-groups,omitempty" protobuf:"bytes,17,opt,name=asGroups"`
	// ImpersonateUserExtra contains additional information for impersonated user.
	// +optional
	ImpersonateUserExtra ImpersonateUserExtra `json:"as-user-extra,omitempty" protobuf:"bytes,18,opt,name=asUserExtra"`
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
	// +optional
	CSIOperator *CSIOperatorFeature `json:"csiOperator,omitempty" protobuf:"bytes,10,opt,name=csiOperator"`
	// For kube-apiserver authorization webhook
	// +optional
	AuthzWebhookAddr *AuthzWebhookAddr `json:"authzWebhookAddr,omitempty" protobuf:"bytes,11,opt,name=authzWebhookAddr"`
	// +optional
	EnableMetricsServer bool `json:"enableMetricsServer,omitempty" protobuf:"bytes,12,opt,name=enableMetricsServer"`
	// +optional
	IPv6DualStack bool `json:"ipv6DualStack,omitempty" protobuf:"bytes,13,opt,name=ipv6DualStack"`
	// +optional
	EnableCilium bool `json:"enableCilium,omitempty" protobuf:"bytes,14,opt,name=enableCilium"`

	ContainerRuntime ContainerRuntimeType `json:"containerRuntime" protobuf:"bytes,15,opt,name=containerRuntime"`
	// Upgrade control upgrade process.
	// +optional
	Upgrade Upgrade `json:"upgrade,omitempty" protobuf:"bytes,22,opt,name=upgrade"`
}

type BootstrapApps []BootstrapApp

type BootstrapApp struct {
	App App `json:"app,omitempty" protobuf:"bytes,1,opt,name=app"`
}

type App struct {
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// +optional
	Spec applicationv1.AppSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type HA struct {
	TKEHA        *TKEHA        `json:"tke,omitempty" protobuf:"bytes,1,opt,name=tke"`
	ThirdPartyHA *ThirdPartyHA `json:"thirdParty,omitempty" protobuf:"bytes,2,opt,name=thirdParty"`
}

type TKEHA struct {
	VIP  string `json:"vip" protobuf:"bytes,1,name=vip"`
	VRID *int32 `json:"vrid,omitempty" protobuf:"bytes,2,name=vrid"`
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

type CSIOperatorFeature struct {
	Version string `json:"version" protobuf:"bytes,1,name=version"`
}

type AuthzWebhookAddr struct {
	// +optional
	Builtin *BuiltinAuthzWebhookAddr `json:"builtin,omitempty" protobuf:"bytes,1,opt,name=builtin"`
	// +optional
	External *ExternalAuthzWebhookAddr `json:"external,omitempty" protobuf:"bytes,2,opt,name=external"`
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
	MaxClusterServiceNum *int32 `json:"maxClusterServiceNum,omitempty" protobuf:"bytes,1,opt,name=maxClusterServiceNum"`
	// +optional
	MaxNodePodNum *int32 `json:"maxNodePodNum,omitempty" protobuf:"bytes,2,opt,name=maxNodePodNum"`
	// +optional
	OversoldRatio map[string]string `json:"oversoldRatio,omitempty" protobuf:"bytes,3,opt,name=oversoldRatio"`
}

// Etcd contains elements describing Etcd configuration.
type Etcd struct {

	// Local provides configuration knobs for configuring the local etcd instance
	// Local and External are mutually exclusive
	Local *LocalEtcd `json:"local,omitempty" protobuf:"bytes,1,opt,name=local"`

	// External describes how to connect to an external etcd cluster
	// Local and External are mutually exclusive
	External *ExternalEtcd `json:"external,omitempty" protobuf:"bytes,2,opt,name=external"`
}

// LocalEtcd describes that kubeadm should run an etcd cluster locally
type LocalEtcd struct {
	// DataDir is the directory etcd will place its data.
	// Defaults to "/var/lib/etcd".
	DataDir string `json:"dataDir" protobuf:"bytes,1,opt,name=dataDir"`

	// ExtraArgs are extra arguments provided to the etcd binary
	// when run inside a static pod.
	ExtraArgs map[string]string `json:"extraArgs,omitempty" protobuf:"bytes,2,rep,name=extraArgs"`

	// ServerCertSANs sets extra Subject Alternative Names for the etcd server signing cert.
	ServerCertSANs []string `json:"serverCertSANs,omitempty" protobuf:"bytes,3,rep,name=serverCertSANs"`
	// PeerCertSANs sets extra Subject Alternative Names for the etcd peer signing cert.
	PeerCertSANs []string `json:"peerCertSANs,omitempty" protobuf:"bytes,4,rep,name=peerCertSANs"`
}

// ExternalEtcd describes an external etcd cluster.
// Kubeadm has no knowledge of where certificate files live and they must be supplied.
type ExternalEtcd struct {
	// Endpoints of etcd members. Required for ExternalEtcd.
	Endpoints []string `json:"endpoints" protobuf:"bytes,1,rep,name=endpoints"`

	// CAFile is an SSL Certificate Authority file used to secure etcd communication.
	// Required if using a TLS connection.
	CAFile string `json:"caFile" protobuf:"bytes,2,opt,name=caFile"`

	// CertFile is an SSL certification file used to secure etcd communication.
	// Required if using a TLS connection.
	CertFile string `json:"certFile" protobuf:"bytes,3,opt,name=certFile"`

	// KeyFile is an SSL key file used to secure etcd communication.
	// Required if using a TLS connection.
	KeyFile string `json:"keyFile" protobuf:"bytes,4,opt,name=keyFile"`
}

type Upgrade struct {
	// Upgrade mode, default value is Auto.
	// +optional
	Mode UpgradeMode `json:"mode,omitempty" protobuf:"bytes,1,opt,name=mode"`
	// Upgrade strategy config.
	// +optional
	Strategy UpgradeStrategy `json:"strategy,omitempty" protobuf:"bytes,2,opt,name=strategy"`
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
	// +optional
	MaxUnready *intstr.IntOrString `json:"maxUnready,omitempty" protobuf:"bytes,1,opt,name=maxUnready"`
	// Whether drain node before upgrade.
	// Draining node before upgrade is recommended.
	// But not all pod running as cows, a few running as pets.
	// If your pod can not accept be expelled from current node, this value should be false.
	// +optional
	DrainNodeBeforeUpgrade *bool `json:"drainNodeBeforeUpgrade,omitempty" protobuf:"varint,2,opt,name=drainNodeBeforeUpgrade"`
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
	// Addon type, one of PersistentEvent or LogCollector etc.
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

// +k8s:conversion-gen:explicit-from=net/url.Values
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
	IP          string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=ip"`
	Port        int32  `json:"port,omitempty" protobuf:"varint,2,opt,name=port"`
	Scheme      string `json:"scheme,omitempty" protobuf:"bytes,3,opt,name=scheme"`
	IndexName   string `json:"indexName,omitempty" protobuf:"bytes,4,opt,name=indexName"`
	User        string `json:"user,omitempty" protobuf:"bytes,5,opt,name=user"`
	Password    string `json:"password,omitempty" protobuf:"bytes,6,opt,name=password"`
	ReserveDays int32  `json:"reserveDays,omitempty" protobuf:"varint,7,opt,name=reserveDays"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProxyOptions is the query options to a proxy call.
type ProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Path is the URL path to use for the current proxy request.
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
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

// +k8s:conversion-gen:explicit-from=net/url.Values
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

// +k8s:conversion-gen:explicit-from=net/url.Values
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
	// If specified, the node's taints.
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty" protobuf:"bytes,12,opt,name=taints"`
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
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of clusters
	Items []Machine `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
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

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterGroupAPIResourceItemsList is the whole list of all ClusterAPIResource.
type ClusterGroupAPIResourceItemsList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,2,opt,name=metadata"`
	// List of ClusterGroupAPIResourceItems
	Items []ClusterGroupAPIResourceItems `json:"items" protobuf:"bytes,3,rep,name=items"`
	// Failed Group Error
	FailedGroupError string `json:"failedGroupError" protobuf:"bytes,4,rep,name=failedGroupError"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=list,get
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterGroupAPIResourceItems contains the GKV for the current kubernetes cluster
type ClusterGroupAPIResourceItems struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// groupVersion is the group and version this APIResourceList is for.
	GroupVersion string `json:"groupVersion" protobuf:"bytes,2,opt,name=groupVersion"`
	// resources contains the name of the resources and if they are namespaced.
	APIResources []ClusterGroupAPIResourceItem `json:"resources" protobuf:"bytes,3,rep,name=resources"`
}

// ClusterGroupAPIResourceItem specifies the name of a resource and whether it is namespaced.
type ClusterGroupAPIResourceItem struct {
	// name is the plural name of the resource.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// singularName is the singular name of the resource.  This allows clients to handle plural and singular opaquely.
	// The singularName is more correct for reporting status on a single item and both singular and plural are allowed
	// from the kubectl CLI interface.
	SingularName string `json:"singularName" protobuf:"bytes,2,opt,name=singularName"`
	// namespaced indicates if a resource is namespaced or not.
	Namespaced bool `json:"namespaced" protobuf:"varint,3,opt,name=namespaced"`
	// group is the preferred group of the resource.  Empty implies the group of the containing resource list.
	// For subresources, this may have a different value, for example: Scale".
	Group string `json:"group" protobuf:"bytes,4,opt,name=group"`
	// version is the preferred version of the resource.  Empty implies the version of the containing resource list
	// For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource's group)".
	Version string `json:"version" protobuf:"bytes,5,opt,name=version"`
	// kind is the kind for the resource (e.g. 'Foo' is the kind for a resource 'foo')
	Kind string `json:"kind" protobuf:"bytes,6,opt,name=kind"`
	// verbs is a list of supported kube verbs (this includes get, list, watch, create,
	// update, patch, delete, deletecollection, and proxy)
	Verbs []string `json:"verbs" protobuf:"bytes,7,rep,name=verbs"`
	// shortNames is a list of suggested short names of the resource.
	ShortNames []string `json:"shortNames" protobuf:"bytes,8,rep,name=shortNames"`
	// categories is a list of the grouped resources this resource belongs to (e.g. 'all')
	Categories []string `json:"categories" protobuf:"bytes,9,rep,name=categories"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterGroupAPIResourceOptions is the query options.
type ClusterGroupAPIResourceOptions struct {
	metav1.TypeMeta `json:",inline"`
}
