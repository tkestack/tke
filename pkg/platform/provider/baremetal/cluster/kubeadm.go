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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	platformv1 "tkestack.io/tke/api/platform/v1"
	kubeadmv1beta2 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
	kubeletv1beta1 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubelet/config/v1beta1"
	kubeproxyv1alpha1 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeproxy/config/v1alpha1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
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

func (p *Provider) getKubeadmJoinConfig(c *v1.Cluster, nodeName string) *kubeadmv1beta2.JoinConfiguration {
	apiServerEndpoint, err := c.HostForBootstrap()
	if err != nil {
		panic(err)
	}

	return &kubeadmv1beta2.JoinConfiguration{
		NodeRegistration: kubeadmv1beta2.NodeRegistrationOptions{
			Name: nodeName,
			KubeletExtraArgs: map[string]string{
				"pod-infra-container-image": images.Get().Pause.FullName(),
			},
		},
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

	return &kubeadmv1beta2.InitConfiguration{
		BootstrapTokens: []kubeadmv1beta2.BootstrapToken{
			{
				Token:       token,
				Description: "TKE kubeadm bootstrap token",
				TTL:         &metav1.Duration{Duration: 0},
			},
		},
		NodeRegistration: kubeadmv1beta2.NodeRegistrationOptions{
			Name: c.Spec.Machines[0].IP,
			KubeletExtraArgs: map[string]string{
				"pod-infra-container-image": images.Get().Pause.FullName(),
			},
		},
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
	}

	return config
}

func (p *Provider) getKubeProxyConfiguration(c *v1.Cluster) *kubeproxyv1alpha1.KubeProxyConfiguration {
	kubeProxyMode := "iptables"
	if c.Spec.Features.IPVS != nil && *c.Spec.Features.IPVS {
		kubeProxyMode = "ipvs"
	}

	return &kubeproxyv1alpha1.KubeProxyConfiguration{
		Mode: kubeproxyv1alpha1.ProxyMode(kubeProxyMode),
	}
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
	}
}

func (p *Provider) getAPIServerExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"token-auth-file": constants.TokenFile,
	}
	if p.config.Audit.Address != "" {
		args["audit-policy-file"] = constants.AuditPolicyConfigFile
		args["audit-webhook-config-file"] = constants.AuditWebhookConfigFile
	}
	for k, v := range c.Spec.APIServerExtraArgs {
		args[k] = v
	}

	return args
}

func (p *Provider) getControllerManagerExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"allocate-node-cidrs":      "true",
		"node-cidr-mask-size":      fmt.Sprintf("%v", c.Status.NodeCIDRMaskSize),
		"cluster-cidr":             c.Spec.ClusterCIDR,
		"service-cluster-ip-range": c.Status.ServiceCIDR,
	}
	for k, v := range c.Spec.ControllerManagerExtraArgs {
		args[k] = v
	}

	return args
}

func (p *Provider) getSchedulerExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{
		"use-legacy-policy-config": "true",
		"policy-config-file":       constants.SchedulerPolicyConfigFile,
	}
	for k, v := range c.Spec.SchedulerExtraArgs {
		args[k] = v
	}

	return args
}
