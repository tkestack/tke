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

package machine

import (
	"fmt"

	"github.com/imdario/mergo"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	utilsnet "k8s.io/utils/net"
	platformv1 "tkestack.io/tke/api/platform/v1"
	kubeadmv1beta2 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
)

func (p *Provider) getKubeadmJoinConfig(c *v1.Cluster, machineIP string) *kubeadmv1beta2.JoinConfiguration {
	apiServerEndpoint, err := c.Host()
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
	}
}

func (p *Provider) getKubeletExtraArgs(c *v1.Cluster) map[string]string {
	args := map[string]string{}
	if c.Cluster.Spec.Features.ContainerRuntime == platformv1.Docker {
		args = map[string]string{
			"pod-infra-container-image": images.Get().Pause.FullName(),
		}
	}
	// for containerd runtimes, no need to set pod-infra-container-image

	utilruntime.Must(mergo.Merge(&args, c.Spec.KubeletExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.config.Kubelet.ExtraArgs))

	return args
}
