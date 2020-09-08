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
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/addons/cniplugins"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/docker"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/gpu"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeconfig"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubelet"
	"tkestack.io/tke/pkg/platform/provider/baremetal/preflight"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util"
	typesv1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/cmdstring"
	"tkestack.io/tke/pkg/util/hosts"
)

const (
	sysctlFile       = "/etc/sysctl.conf"
	sysctlCustomFile = "/etc/sysctl.d/99-tke.conf"
	moduleFile       = "/etc/modules-load.d/tke.conf"
)

func (p *Provider) EnsureCopyFiles(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	for _, file := range cluster.Spec.Features.Files {
		err = machineSSH.CopyFile(file.Src, file.Dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) EnsurePreInstallHook(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	hook := cluster.Spec.Features.Hooks[platformv1.HookPreInstall]
	if hook == "" {
		return nil
	}

	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	cmd := strings.Split(hook, " ")[0]

	machineSSH.Execf("chmod +x %s", cmd)
	_, stderr, exit, err := machineSSH.Exec(hook)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", hook, exit, stderr, err)
	}
	return nil
}

func (p *Provider) EnsurePostInstallHook(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	hook := cluster.Spec.Features.Hooks[platformv1.HookPostInstall]
	if hook == "" {
		return nil
	}

	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	cmd := strings.Split(hook, " ")[0]

	machineSSH.Execf("chmod +x %s", cmd)
	_, stderr, exit, err := machineSSH.Exec(hook)
	if err != nil || exit != 0 {
		return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", hook, exit, stderr, err)
	}
	return nil
}

func (p *Provider) EnsureClean(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	_, err = machineSSH.CombinedOutput(fmt.Sprintf("rm -rf %s", constants.KubernetesDir))
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsurePreflight(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	err = preflight.RunNodeChecks(machineSSH)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureRegistryHosts(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	if !p.config.Registry.NeedSetHosts() {
		return nil
	}

	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	domains := []string{
		p.config.Registry.Domain,
	}
	if machine.Spec.TenantID != "" {
		domains = append(domains, machine.Spec.TenantID+"."+p.config.Registry.Domain)
	}
	for _, one := range domains {
		remoteHosts := hosts.RemoteHosts{Host: one, SSH: machineSSH}
		err := remoteHosts.Set(p.config.Registry.IP)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) EnsureKernelModule(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	s, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	modules := []string{"iptable_nat", "ip_vs", "ip_vs_rr", "ip_vs_wrr", "ip_vs_sh"}
	if _, err := s.CombinedOutput("modinfo br_netfilter"); err == nil {
		modules = append(modules, "br_netfilter")
	}
	var data bytes.Buffer
	for _, m := range modules {
		_, err := s.CombinedOutput(fmt.Sprintf("modprobe %s", m))
		if err != nil {
			return err
		}
		data.WriteString(m + "\n")
	}
	err = s.WriteFile(strings.NewReader(data.String()), moduleFile)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureSysctl(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	_, err = machineSSH.CombinedOutput(cmdstring.SetFileContent(sysctlFile, "^net.ipv4.ip_forward.*", "net.ipv4.ip_forward = 1"))
	if err != nil {
		return err
	}

	_, err = machineSSH.CombinedOutput(cmdstring.SetFileContent(sysctlFile, "^net.bridge.bridge-nf-call-iptables.*", "net.bridge.bridge-nf-call-iptables = 1"))
	if err != nil {
		return err
	}

	f, err := os.Open(path.Join(constants.ConfDir, "sysctl.conf"))
	if err == nil {
		err = machineSSH.WriteFile(f, sysctlCustomFile)
		if err != nil {
			return err
		}
	}

	_, err = machineSSH.CombinedOutput("sysctl --system")
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureDisableSwap(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	_, err = machineSSH.CombinedOutput(`swapoff -a && sed -i "s/^[^#]*swap/#&/" /etc/fstab`)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureManifestDir(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}
	_, err = machineSSH.CombinedOutput("mkdir -p /etc/kubernetes/manifests")
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureKubeconfig(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	masterEndpoint, err := util.GetMasterEndpoint(cluster.Status.Addresses)
	if err != nil {
		return err
	}

	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	option := &kubeconfig.Option{
		MasterEndpoint: masterEndpoint,
		ClusterName:    cluster.Name,
		CACert:         cluster.ClusterCredential.CACert,
		Token:          *cluster.ClusterCredential.Token,
	}
	err = kubeconfig.Install(machineSSH, option)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureNvidiaDriver(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	if !gpu.IsEnable(machine.Spec.Labels) {
		return nil
	}

	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	return gpu.InstallNvidiaDriver(machineSSH, &gpu.NvidiaDriverOption{})
}

func (p *Provider) EnsureNvidiaContainerRuntime(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	if !gpu.IsEnable(machine.Spec.Labels) {
		return nil
	}

	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	return gpu.InstallNvidiaContainerRuntime(machineSSH, &gpu.NvidiaContainerRuntimeOption{})
}

func (p *Provider) EnsureDocker(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	insecureRegistries := fmt.Sprintf(`"%s"`, p.config.Registry.Domain)
	if p.config.Registry.NeedSetHosts() && machine.Spec.TenantID != "" {
		insecureRegistries = fmt.Sprintf(`%s,"%s"`, insecureRegistries, machine.Spec.TenantID+"."+p.config.Registry.Domain)
	}

	extraArgs := cluster.Spec.DockerExtraArgs
	utilruntime.Must(mergo.Merge(&extraArgs, p.config.Docker.ExtraArgs))
	option := &docker.Option{
		InsecureRegistries: insecureRegistries,
		RegistryDomain:     p.config.Registry.Domain,
		IsGPU:              gpu.IsEnable(machine.Spec.Labels),
		ExtraArgs:          extraArgs,
	}
	err = docker.Install(machineSSH, option)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureKubelet(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	err = kubelet.Install(machineSSH, cluster.Spec.Version)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureCNIPlugins(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	option := &cniplugins.Option{}
	err = cniplugins.Install(machineSSH, option)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureConntrackTools(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	err = res.ConntrackTools.InstallWithDefault(machineSSH)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureKubeadm(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	err = kubeadm.Install(machineSSH, cluster.Spec.Version)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureJoinPhasePreflight(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(cluster, machine.Spec.IP), "preflight")
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseKubeletStart(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	machineSSH, err := machine.Spec.SSH()
	if err != nil {
		return err
	}

	err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(cluster, machine.Spec.IP), "kubelet-start")
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureMarkNode(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	clientset, err := cluster.Clientset()
	if err != nil {
		return err
	}

	node, err := apiclient.GetNodeByMachineIP(ctx, clientset, machine.Spec.IP)
	if err != nil {
		return err
	}
	err = apiclient.MarkNode(ctx, clientset, node.Name, machine.Spec.Labels, machine.Spec.Taints)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureNodeReady(ctx context.Context, machine *platformv1.Machine, cluster *typesv1.Cluster) error {
	clientset, err := cluster.Clientset()
	if err != nil {
		return err
	}

	return wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		node, err := apiclient.GetNodeByMachineIP(ctx, clientset, machine.Spec.IP)
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
