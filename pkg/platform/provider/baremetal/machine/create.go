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

package machine

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/addons/cniplugins"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/docker"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeconfig"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubelet"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/marknode"
	"tkestack.io/tke/pkg/platform/provider/baremetal/preflight"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util/hosts"
)

const (
	sysctlFile       = "/etc/sysctl.conf"
	sysctlCustomFile = "/etc/sysctl.d/99-tke.conf"
	moduleFile       = "/etc/modules-load.d/tke.conf"
)

func (p *Provider) EnsureClean(m *Machine) error {
	_, err := m.CombinedOutput(fmt.Sprintf("rm -rf %s", constants.KubernetesDir))
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsurePreflight(m *Machine) error {
	err := preflight.RunNodeChecks(m)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureRegistryHosts(m *Machine) error {
	if !m.Registry.UseTKE() {
		return nil
	}

	domains := []string{
		m.Registry.Domain,
		m.Spec.TenantID + "." + m.Registry.Domain,
	}
	for _, one := range domains {
		remoteHosts := hosts.RemoteHosts{Host: one, SSH: m}
		err := remoteHosts.Set(m.Registry.IP)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) EnsureKernelModule(m *Machine) error {
	modules := []string{"iptable_nat"}
	var data bytes.Buffer
	s := m

	for _, m := range modules {
		_, err := s.CombinedOutput(fmt.Sprintf("modprobe %s", m))
		if err != nil {
			return err
		}
		data.WriteString(m + "\n")
	}
	err := s.WriteFile(strings.NewReader(data.String()), moduleFile)
	if err != nil {
		return err
	}

	return nil
}

func setFileContent(file, pattern, content string) string {
	return fmt.Sprintf("grep -Pq '%s' %s && sed -i 's;%s;%s;g' %s|| echo '%s' >> %s",
		pattern, file,
		pattern, content, file,
		content, file)
}

func (p *Provider) EnsureSysctl(m *Machine) error {
	_, err := m.CombinedOutput(setFileContent(sysctlFile, "net.ipv4.ip_forward.*", "net.ipv4.ip_forward = 1"))
	if err != nil {
		return err
	}

	_, err = m.CombinedOutput(setFileContent(sysctlFile, "net.bridge.bridge-nf-call-iptables.*", "net.bridge.bridge-nf-call-iptables = 1"))
	if err != nil {
		return err
	}

	f, err := os.Open(path.Join(constants.ConfDir, "sysctl.conf"))
	if err == nil {
		err = m.WriteFile(f, sysctlCustomFile)
		if err != nil {
			return err
		}
	}

	_, err = m.CombinedOutput("sysctl --system")
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureDisableSwap(m *Machine) error {
	_, err := m.CombinedOutput("swapoff -a && sed -i 's/^[^#]*swap/#&/' /etc/fstab")
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureKubeconfig(m *Machine) error {
	masterEndpoint, err := util.GetMasterEndpoint(m.Cluster.Status.Addresses)
	if err != nil {
		return err
	}

	option := &kubeconfig.Option{
		MasterEndpoint: masterEndpoint,
		ClusterName:    m.Cluster.Name,
		CACert:         m.ClusterCredential.CACert,
		Token:          *m.ClusterCredential.Token,
	}
	err = kubeconfig.Install(m, option)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureDocker(m *Machine) error {
	insecureRegistries := fmt.Sprintf(`"%s"`, m.Registry.Domain)
	if m.Config.Registry.UseTKE() {
		insecureRegistries = fmt.Sprintf(`%s,"%s"`, insecureRegistries, m.Spec.TenantID+"."+m.Registry.Domain)
	}

	option := &docker.Option{
		Version:            m.Docker.DefaultVersion,
		InsecureRegistries: insecureRegistries,
		RegistryDomain:     m.Registry.Domain,
		IsGPU:              IsGPU(m.Spec.Labels),
		ExtraArgs:          m.Cluster.Spec.DockerExtraArgs,
	}
	err := docker.Install(m, option)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureKubelet(m *Machine) error {
	option := &kubelet.Option{
		Version:   m.Cluster.Spec.Version,
		ExtraArgs: m.Cluster.Spec.KubeletExtraArgs,
	}
	err := kubelet.Install(m, option)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureCNIPlugins(m *Machine) error {
	option := &cniplugins.Option{
		Version: m.CNIPlugins.DefaultVersion,
	}
	err := cniplugins.Install(m, option)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureKubeadm(m *Machine) error {
	err := kubeadm.Install(m)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureJoinNode(m *Machine) error {
	masterEndpoint, err := util.GetMasterEndpoint(m.Cluster.Status.Addresses)
	if err != nil {
		return err
	}

	option := &kubeadm.JoinNodeOption{
		NodeName:             m.Spec.IP,
		BootstrapToken:       *m.ClusterCredential.BootstrapToken,
		ControlPlaneEndpoint: strings.TrimPrefix(masterEndpoint, "https://"),
	}
	err = kubeadm.JoinNode(m, option)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureMarkNode(m *Machine) error {
	if len(m.Spec.Labels) == 0 {
		return nil
	}

	option := &marknode.Option{
		NodeName: m.Spec.IP,
		Labels:   m.Spec.Labels,
	}
	err := marknode.Install(m.ClientSet, option)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureNodeReady(m *Machine) error {
	return wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		node, err := m.ClientSet.CoreV1().Nodes().Get(m.Spec.IP, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		for _, one := range node.Status.Conditions {
			if one.Type == corev1.NodeReady && one.Status == corev1.ConditionTrue {
				return true, nil
			}
		}

		return false, nil
	})
}
