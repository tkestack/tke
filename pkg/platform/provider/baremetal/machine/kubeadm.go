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
	"github.com/imdario/mergo"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeadmv1beta2 "tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
)

func (p *Provider) getKubeadmJoinConfig(c *v1.Cluster, machineIP string) *kubeadmv1beta2.JoinConfiguration {
	apiServerEndpoint, err := c.Host()
	if err != nil {
		panic(err)
	}

	kubeletExtraArgs := p.getKubeletExtraArgs(c)
	if _, ok := kubeletExtraArgs["node-ip"]; !ok {
		kubeletExtraArgs["node-ip"] = machineIP
	}

	return &kubeadmv1beta2.JoinConfiguration{
		NodeRegistration: kubeadmv1beta2.NodeRegistrationOptions{
			Name:             machineIP,
			KubeletExtraArgs: kubeletExtraArgs,
		},
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
	args := map[string]string{
		"pod-infra-container-image": images.Get().Pause.FullName(),
	}

	utilruntime.Must(mergo.Merge(&args, c.Spec.KubeletExtraArgs))
	utilruntime.Must(mergo.Merge(&args, p.config.Kubelet.ExtraArgs))

	return args
}
