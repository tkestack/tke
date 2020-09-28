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

	"github.com/imdario/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	platformv1 "tkestack.io/tke/api/platform/v1"
	kubeadmv1beta2 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
	kubeletv1beta1 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubelet/config/v1beta1"
	kubeproxyv1alpha1 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeproxy/config/v1alpha1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/json"
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
	// add label to get node by machine ip.
	kubeletExtraArgs["node-labels"] = fields.OneTermEqualSelector(string(apiclient.LabelMachineIP), machineIP).String()
	nodeRegistration.KubeletExtraArgs = kubeletExtraArgs

	if !c.Spec.HostnameAsNodename {
		nodeRegistration.Name = machineIP
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
	// add label to get node by machine ip.
	kubeletExtraArgs["node-labels"] = fields.OneTermEqualSelector(string(apiclient.LabelMachineIP), c.Spec.Machines[0].IP).String()
	nodeRegistration.KubeletExtraArgs = kubeletExtraArgs

	if !c.Spec.HostnameAsNodename {
		nodeRegistration.Name = c.Spec.Machines[0].IP
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
			AdvertiseAddress: c.Spec.Machines[0].IP,
		},
		CertificateKey: *c.ClusterCredential.CertificateKey,
	}
}

func (p *Provider) getClusterConfiguration(c *v1.Cluster) *kubeadmv1beta2.ClusterConfiguration {
	controlPlaneEndpoint := fmt.Sprintf("%s:6443", c.Spec.Machines[0].IP)
	addr := c.Address(platformv1.AddressAdvertise)
	if addr != nil {
		controlPlaneEndpoint = fmt.Sprintf("%s:%d", addr.Host, addr.Port)
	}

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
			ImageMeta: kubeadmv1beta2.ImageMeta{
				ImageTag: images.Get().CoreDNS.Tag,
			},
		},
		ImageRepository: p.config.Registry.Prefix,
		ClusterName:     c.Name,
	}

	utilruntime.Must(json.Merge(&config.Etcd, &c.Spec.Etcd))
	if config.Etcd.Local != nil {
		config.Etcd.Local.ImageTag = images.Get().ETCD.Tag

		if config.Etcd.Local.ExtraArgs != nil && p.config.Etcd.ExtraArgs != nil {
			utilruntime.Must(mergo.Merge(&config.Etcd.Local.ExtraArgs, p.config.Etcd.ExtraArgs))
		}
	}

	return config
}

func (p *Provider) getKubeProxyConfiguration(c *v1.Cluster) *kubeproxyv1alpha1.KubeProxyConfiguration {
	config := &kubeproxyv1alpha1.KubeProxyConfiguration{}
	config.Mode = "iptables"
	if c.Spec.Features.IPVS != nil && *c.Spec.Features.IPVS {
		config.Mode = "ipvs"
		config.ClusterCIDR = c.Spec.ClusterCIDR
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
	if p.config.AuditEnabled() {
		args["audit-policy-file"] = constants.KubernetesAuditPolicyConfigFile
		args["audit-webhook-config-file"] = constants.KuberentesAuditWebhookConfigFile
	}
	if c.AuthzWebhookEnabled() {
		args["authorization-webhook-config-file"] = constants.KubernetesAuthzWebhookConfigFile
		args["authorization-mode"] = "Node,RBAC,Webhook"
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.APIServerExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.config.APIServer.ExtraArgs))

	return args
}

func (p *Provider) getControllerManagerExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"allocate-node-cidrs":      "true",
		"node-cidr-mask-size":      fmt.Sprintf("%v", c.Status.NodeCIDRMaskSize),
		"cluster-cidr":             c.Spec.ClusterCIDR,
		"service-cluster-ip-range": c.Status.ServiceCIDR,
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.ControllerManagerExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.config.ControllerManager.ExtraArgs))

	return args
}

func (p *Provider) getSchedulerExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"use-legacy-policy-config": "true",
		"policy-config-file":       constants.KuberentesSchedulerPolicyConfigFile,
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.SchedulerExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.config.Scheduler.ExtraArgs))

	return args
}

func (p *Provider) getKubeletExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"pod-infra-container-image": images.Get().Pause.FullName(),
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.KubeletExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.config.Kubelet.ExtraArgs))

	return args
}
