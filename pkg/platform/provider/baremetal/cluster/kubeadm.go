/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package cluster

import (
	"fmt"
	"net"
	"path"

	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	utilsnet "k8s.io/utils/net"
	kubeadmv1beta2 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
	kubeletv1beta1 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubelet/config/v1beta1"
	kubeproxyv1alpha1 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeproxy/config/v1alpha1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/json"
	"tkestack.io/tke/pkg/util/version"
)

func (p *Provider) getKubeadmInitConfig(c *v1.Cluster) *kubeadm.InitConfig {
	config := new(kubeadm.InitConfig)
	config.InitConfiguration = p.getInitConfiguration(c)
	config.ClusterConfiguration = p.getClusterConfiguration(c)
	config.KubeProxyConfiguration = p.getKubeProxyConfiguration(c)
	config.KubeletConfiguration = p.getKubeletConfiguration(c)

	return config
}

func (p *Provider) getKubeadmJoinConfig(c *v1.Cluster, machineIP string) *kubeadmv1beta2.JoinConfiguration {
	apiServerEndpoint, err := c.HostForBootstrap()
	if err != nil {
		panic(err)
	}

	nodeRegistration := kubeadmv1beta2.NodeRegistrationOptions{}
	kubeletExtraArgs := p.getKubeletExtraArgs(c)
	if !utilsnet.IsIPv6String(c.Spec.Machines[0].IP) {
		kubeletExtraArgs["node-labels"] = fmt.Sprintf("%s=%s", apiclient.LabelMachineIPV4, machineIP)
	} else {
		kubeletExtraArgs["node-labels"] = apiclient.GetNodeIPV6Label(machineIP)
	}
	if c.Cluster.Spec.Features.EnableCilium && c.Cluster.Spec.NetworkArgs["networkMode"] == "underlay" {
		if asn, ok := c.Cluster.Spec.NetworkArgs["asn"]; ok {
			kubeletExtraArgs["node-labels"] = fmt.Sprintf("%s,%s=%s", kubeletExtraArgs["node-labels"], apiclient.LabelASNCilium, asn)
		}
		if switchIP, ok := c.Cluster.Spec.NetworkArgs["switch-ip"]; ok {
			kubeletExtraArgs["node-labels"] = fmt.Sprintf("%s,%s=%s", kubeletExtraArgs["node-labels"], apiclient.LabelSwitchIPCilium, switchIP)
		}
	}
	kubeletExtraArgs["node-labels"] = fmt.Sprintf("%s,%s=%s", kubeletExtraArgs["node-lables"], apiclient.LabelTopologyZone, "default")

	if _, ok := kubeletExtraArgs["hostname-override"]; !ok {
		if !c.Spec.HostnameAsNodename {
			nodeRegistration.Name = machineIP
		}
	}
	nodeRegistration.KubeletExtraArgs = kubeletExtraArgs
	// Specify cri runtime type
	if c.Cluster.Spec.Features.ContainerRuntime == "docker" {
		nodeRegistration.CRISocket = "/var/run/dockershim.sock"
	} else {
		nodeRegistration.CRISocket = "/var/run/containerd/containerd.sock"
	}

	return &kubeadmv1beta2.JoinConfiguration{
		NodeRegistration: nodeRegistration,
		Discovery: kubeadmv1beta2.Discovery{
			BootstrapToken: &kubeadmv1beta2.BootstrapTokenDiscovery{
				Token:                    *c.ClusterCredential.BootstrapToken,
				APIServerEndpoint:        apiServerEndpoint,
				UnsafeSkipCAVerification: true,
			},
			TLSBootstrapToken: *c.ClusterCredential.BootstrapToken,
		},
		ControlPlane: &kubeadmv1beta2.JoinControlPlane{
			CertificateKey: *c.ClusterCredential.CertificateKey,
		},
	}
}

func (p *Provider) getInitConfiguration(c *v1.Cluster) *kubeadmv1beta2.InitConfiguration {
	token, _ := kubeadmv1beta2.NewBootstrapTokenString(*c.ClusterCredential.BootstrapToken)

	nodeRegistration := kubeadmv1beta2.NodeRegistrationOptions{}
	kubeletExtraArgs := p.getKubeletExtraArgs(c)
	machineIP := c.Spec.Machines[0].IP
	if !utilsnet.IsIPv6String(c.Spec.Machines[0].IP) {
		kubeletExtraArgs["node-labels"] = fmt.Sprintf("%s=%s", apiclient.LabelMachineIPV4, machineIP)
	} else {
		kubeletExtraArgs["node-labels"] = apiclient.GetNodeIPV6Label(machineIP)
	}
	if c.Cluster.Spec.Features.EnableCilium && c.Cluster.Spec.NetworkArgs["networkMode"] == "underlay" {
		if asn, ok := c.Cluster.Spec.NetworkArgs["asn"]; ok {
			kubeletExtraArgs["node-labels"] = fmt.Sprintf("%s,%s=%s", kubeletExtraArgs["node-labels"], apiclient.LabelASNCilium, asn)
		}
		if switchIP, ok := c.Cluster.Spec.NetworkArgs["switch-ip"]; ok {
			kubeletExtraArgs["node-labels"] = fmt.Sprintf("%s,%s=%s", kubeletExtraArgs["node-labels"], apiclient.LabelSwitchIPCilium, switchIP)
		}
	}
	kubeletExtraArgs["node-labels"] = fmt.Sprintf("%s,%s=%s", kubeletExtraArgs["node-lables"], apiclient.LabelTopologyZone, "default")
	// add node ip for single stack ipv6 clusters.
	if _, ok := kubeletExtraArgs["node-ip"]; !ok {
		kubeletExtraArgs["node-ip"] = machineIP
	}
	if _, ok := kubeletExtraArgs["hostname-override"]; !ok {
		if !c.Spec.HostnameAsNodename {
			nodeRegistration.Name = machineIP
		}
	}
	nodeRegistration.KubeletExtraArgs = kubeletExtraArgs
	// Specify cri runtime type
	if c.Cluster.Spec.Features.ContainerRuntime == "docker" {
		nodeRegistration.CRISocket = "/var/run/dockershim.sock"
	} else {
		nodeRegistration.CRISocket = "/var/run/containerd/containerd.sock"
	}
	return &kubeadmv1beta2.InitConfiguration{
		BootstrapTokens: []kubeadmv1beta2.BootstrapToken{
			{
				Token:       token,
				Description: "TKE kubeadm bootstrap token",
				TTL:         &metav1.Duration{Duration: 0},
			},
		},
		NodeRegistration: nodeRegistration,
		LocalAPIEndpoint: kubeadmv1beta2.APIEndpoint{
			AdvertiseAddress: machineIP,
		},
		CertificateKey: *c.ClusterCredential.CertificateKey,
	}
}

func (p *Provider) getClusterConfiguration(c *v1.Cluster) *kubeadmv1beta2.ClusterConfiguration {
	controlPlaneEndpoint := net.JoinHostPort(constants.APIServerHostName, "6443")

	kubernetesVolume := kubeadmv1beta2.HostPathMount{
		Name:      "vol-dir-0",
		HostPath:  "/etc/kubernetes",
		MountPath: "/etc/kubernetes",
	}

	config := &kubeadmv1beta2.ClusterConfiguration{
		Networking: kubeadmv1beta2.Networking{
			DNSDomain:     c.Spec.DNSDomain,
			ServiceSubnet: c.Status.ServiceCIDR,
		},
		KubernetesVersion:    c.Spec.Version,
		ControlPlaneEndpoint: controlPlaneEndpoint,
		APIServer: kubeadmv1beta2.APIServer{
			ControlPlaneComponent: kubeadmv1beta2.ControlPlaneComponent{
				ExtraArgs:    p.getAPIServerExtraArgs(c),
				ExtraVolumes: []kubeadmv1beta2.HostPathMount{kubernetesVolume},
			},
			CertSANs: GetAPIServerCertSANs(c.Cluster),
		},
		ControllerManager: kubeadmv1beta2.ControlPlaneComponent{
			ExtraArgs:    p.getControllerManagerExtraArgs(c),
			ExtraVolumes: []kubeadmv1beta2.HostPathMount{kubernetesVolume},
		},
		Scheduler: kubeadmv1beta2.ControlPlaneComponent{
			ExtraArgs:    p.getSchedulerExtraArgs(c),
			ExtraVolumes: []kubeadmv1beta2.HostPathMount{kubernetesVolume},
		},
		DNS: kubeadmv1beta2.DNS{
			Type: kubeadmv1beta2.CoreDNS,
		},
		ImageRepository: p.getImagePrefix(c),
		ClusterName:     c.Name,
		FeatureGates: map[string]bool{
			"IPv6DualStack": c.Cluster.Spec.Features.IPv6DualStack},
	}

	if p.needSetCoreDNS(c.Spec.Version) {
		config.DNS.ImageTag = images.Get().CoreDNS.Tag
	}

	utilruntime.Must(json.Merge(&config.Etcd, &c.Spec.Etcd))
	if config.Etcd.Local != nil {
		config.Etcd.Local.ImageTag = images.Get().ETCD.Tag
		config.Etcd.Local.ExtraArgs = map[string]string{
			"quota-backend-bytes": "6442450944",
		}
	}

	return config
}

func (Provider) needSetCoreDNS(k8sVersion string) bool {
	return version.Compare(k8sVersion, constants.NeedUpgradeCoreDNSLowerK8sVersion) < 0 ||
		version.Compare(k8sVersion, constants.NeedUpgradeCoreDNSUpperK8sVersion) >= 0
}

func (p *Provider) getKubeProxyConfiguration(c *v1.Cluster) *kubeproxyv1alpha1.KubeProxyConfiguration {
	config := &kubeproxyv1alpha1.KubeProxyConfiguration{}
	config.Mode = "iptables"
	if c.Spec.Features.IPVS != nil && *c.Spec.Features.IPVS {
		config.Mode = "ipvs"
		config.ClusterCIDR = c.Spec.ClusterCIDR
		if c.Spec.Features.HA != nil {
			if c.Spec.Features.HA.TKEHA != nil {
				config.IPVS.ExcludeCIDRs = []string{fmt.Sprintf("%s/32", c.Spec.Features.HA.TKEHA.VIP)}
			}
			if c.Spec.Features.HA.ThirdPartyHA != nil {
				config.IPVS.ExcludeCIDRs = []string{fmt.Sprintf("%s/32", c.Spec.Features.HA.ThirdPartyHA.VIP)}
			}
		}
	}
	if utilsnet.IsIPv6CIDRString(c.Spec.ClusterCIDR) {
		config.BindAddress = "::"
	}

	return config
}

func (p *Provider) getKubeletConfiguration(c *v1.Cluster) *kubeletv1beta1.KubeletConfiguration {
	return &kubeletv1beta1.KubeletConfiguration{
		KubeReserved: map[string]string{
			"cpu":    "100m",
			"memory": "500Mi",
		},
		SystemReserved: map[string]string{
			"cpu":    "100m",
			"memory": "500Mi",
		},
		MaxPods: *c.Spec.Properties.MaxNodePodNum,
	}
}

func (p *Provider) getAPIServerExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"token-auth-file": constants.TokenFile,
	}
	if p.Config.AuditEnabled() {
		args["audit-policy-file"] = constants.KubernetesAuditPolicyConfigFile
		args["audit-webhook-config-file"] = constants.KubernetesAuditWebhookConfigFile
	}
	if c.AuthzWebhookEnabled() {
		args["authorization-webhook-config-file"] = constants.KubernetesAuthzWebhookConfigFile
		args["authorization-mode"] = "Node,RBAC,Webhook"
	}
	for k, v := range c.Spec.APIServerExtraArgs {
		args[k] = v
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.APIServerExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.Config.APIServer.ExtraArgs))

	return args
}

func (p *Provider) getControllerManagerExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"allocate-node-cidrs": "true",
		"cluster-cidr":        c.Spec.ClusterCIDR,
		"bind-address":        "0.0.0.0",
	}
	if c.Spec.Features.IPv6DualStack {
		args["node-cidr-mask-size-ipv4"] = fmt.Sprintf("%v", c.Status.NodeCIDRMaskSizeIPv4)
		args["node-cidr-mask-size-ipv6"] = fmt.Sprintf("%v", c.Status.NodeCIDRMaskSizeIPv6)
		args["service-cluster-ip-range"] = *c.Spec.ServiceCIDR
	} else {
		args["node-cidr-mask-size"] = fmt.Sprintf("%v", c.Status.NodeCIDRMaskSize)
		args["service-cluster-ip-range"] = c.Status.ServiceCIDR
	}
	if c.Spec.Features.EnableCilium && c.Spec.NetworkArgs["networkMode"] == "overlay" {
		args["configure-cloud-routes"] = "false"
		args["allocate-node-cidrs"] = "false"
	}
	for k, v := range c.Spec.ControllerManagerExtraArgs {
		args[k] = v
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.ControllerManagerExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.Config.ControllerManager.ExtraArgs))

	return args
}

func (p *Provider) getSchedulerExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"use-legacy-policy-config": "true",
		"policy-config-file":       constants.KubernetesSchedulerPolicyConfigFile,
		"bind-address":             "0.0.0.0",
	}
	for k, v := range c.Spec.SchedulerExtraArgs {
		args[k] = v
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.SchedulerExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.Config.Scheduler.ExtraArgs))

	return args
}

func (p *Provider) getKubeletExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"pod-infra-container-image": path.Join(p.getImagePrefix(c), images.Get().Pause.BaseName()),
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.KubeletExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.Config.Kubelet.ExtraArgs))

	return args
}
