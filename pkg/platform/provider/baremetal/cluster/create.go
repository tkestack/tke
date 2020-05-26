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

package cluster

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"tkestack.io/tke/pkg/platform/provider/baremetal/res"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/thoas/go-funk"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/addons/cniplugins"
	csioperatorimage "tkestack.io/tke/pkg/platform/provider/baremetal/phases/csioperator/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/docker"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy"
	galaxyimages "tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/gpu"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeconfig"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubelet"
	"tkestack.io/tke/pkg/platform/provider/baremetal/preflight"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/cmdstring"
	containerregistryutil "tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/hosts"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/template"
)

const (
	sysctlFile       = "/etc/sysctl.conf"
	sysctlCustomFile = "/etc/sysctl.d/99-tke.conf"
	moduleFile       = "/etc/modules-load.d/tke.conf"
)

func (p *Provider) EnsureCopyFiles(ctx context.Context, c *v1.Cluster) error {
	for _, file := range c.Spec.Features.Files {
		for _, machine := range c.Spec.Machines {
			machineSSH, err := machine.SSH()
			if err != nil {
				return err
			}

			err = machineSSH.CopyFile(file.Src, file.Dst)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Provider) EnsurePreInstallHook(ctx context.Context, c *v1.Cluster) error {
	hook := c.Spec.Features.Hooks[platformv1.HookPreInstall]
	if hook == "" {
		return nil
	}
	cmd := strings.Split(hook, " ")[0]

	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		machineSSH.Execf("chmod +x %s", cmd)
		_, stderr, exit, err := machineSSH.Exec(hook)
		if err != nil || exit != 0 {
			return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", hook, exit, stderr, err)
		}
	}
	return nil
}

func (p *Provider) EnsurePostInstallHook(ctx context.Context, c *v1.Cluster) error {
	hook := c.Spec.Features.Hooks[platformv1.HookPostInstall]
	if hook == "" {
		return nil
	}
	cmd := strings.Split(hook, " ")[0]

	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		machineSSH.Execf("chmod +x %s", cmd)
		_, stderr, exit, err := machineSSH.Exec(hook)
		if err != nil || exit != 0 {
			return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", hook, exit, stderr, err)
		}
	}
	return nil
}

func (p *Provider) EnsurePreflight(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = preflight.RunMasterChecks(machineSSH)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureRegistryHosts(ctx context.Context, c *v1.Cluster) error {
	if !p.config.Registry.NeedSetHosts() {
		return nil
	}

	domains := []string{
		p.config.Registry.Domain,
	}
	if c.Spec.TenantID != "" {
		domains = append(domains, c.Spec.TenantID+"."+p.config.Registry.Domain)
	}
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		for _, one := range domains {
			remoteHosts := hosts.RemoteHosts{Host: one, SSH: machineSSH}
			err := remoteHosts.Set(p.config.Registry.IP)
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
		}
	}

	return nil
}

func (p *Provider) EnsureKernelModule(ctx context.Context, c *v1.Cluster) error {
	modules := []string{"iptable_nat"}
	var data bytes.Buffer
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		for _, m := range modules {
			_, err := machineSSH.CombinedOutput(fmt.Sprintf("modprobe %s", m))
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
			data.WriteString(m + "\n")
		}
		err = machineSSH.WriteFile(strings.NewReader(data.String()), moduleFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureSysctl(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		_, err = machineSSH.CombinedOutput(cmdstring.SetFileContent(sysctlFile, "^net.ipv4.ip_forward.*", "net.ipv4.ip_forward = 1"))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		_, err = machineSSH.CombinedOutput(cmdstring.SetFileContent(sysctlFile, "^net.bridge.bridge-nf-call-iptables.*", "net.bridge.bridge-nf-call-iptables = 1"))
		if err != nil {
			return errors.Wrap(err, machine.IP)
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
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureDisableSwap(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		_, err = machineSSH.CombinedOutput("swapoff -a && sed -i 's/^[^#]*swap/#&/' /etc/fstab")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

// 因为validate那里没法更新对象（不能存储）
// PreCrete，在api中错误只能panic，响应不会有报错提示，所以只能挪到这里处理
func (p *Provider) EnsureClusterComplete(ctx context.Context, cluster *v1.Cluster) error {
	funcs := []func(cluster *v1.Cluster) error{
		completeNetworking,
		completeDNS,
		completeGPU,
		completeAddresses,
		completeCredential,
	}
	for _, f := range funcs {
		if err := f(cluster); err != nil {
			return err
		}
	}
	return nil
}

func completeNetworking(cluster *v1.Cluster) error {
	var (
		serviceCIDR      string
		nodeCIDRMaskSize int32
		err              error
	)

	if cluster.Spec.ServiceCIDR != nil {
		serviceCIDR = *cluster.Spec.ServiceCIDR
		nodeCIDRMaskSize, err = GetNodeCIDRMaskSize(cluster.Spec.ClusterCIDR, *cluster.Spec.Properties.MaxNodePodNum)
		if err != nil {
			return errors.Wrap(err, "GetNodeCIDRMaskSize error")
		}
	} else {
		serviceCIDR, nodeCIDRMaskSize, err = GetServiceCIDRAndNodeCIDRMaskSize(cluster.Spec.ClusterCIDR, *cluster.Spec.Properties.MaxClusterServiceNum, *cluster.Spec.Properties.MaxNodePodNum)
		if err != nil {
			return errors.Wrap(err, "GetServiceCIDRAndNodeCIDRMaskSize error")
		}
	}
	cluster.Status.ServiceCIDR = serviceCIDR
	cluster.Status.NodeCIDRMaskSize = nodeCIDRMaskSize

	return nil
}

func completeDNS(cluster *v1.Cluster) error {
	ip, err := GetIndexedIP(cluster.Status.ServiceCIDR, constants.DNSIPIndex)
	if err != nil {
		return errors.Wrap(err, "get DNS IP error")
	}
	cluster.Status.DNSIP = ip.String()

	return nil
}

func completeGPU(cluster *v1.Cluster) error {
	ip, err := GetIndexedIP(cluster.Status.ServiceCIDR, constants.GPUQuotaAdmissionIPIndex)
	if err != nil {
		return errors.Wrap(err, "get gpu quota admission IP error")
	}
	if cluster.Annotations == nil {
		cluster.Annotations = make(map[string]string)
	}
	cluster.Annotations[constants.GPUQuotaAdmissionIPAnnotaion] = ip.String()

	return nil
}

func completeAddresses(cluster *v1.Cluster) error {
	for _, m := range cluster.Spec.Machines {
		cluster.AddAddress(platformv1.AddressReal, m.IP, 6443)
	}

	if cluster.Spec.Features.HA != nil {
		if cluster.Spec.Features.HA.TKEHA != nil {
			cluster.AddAddress(platformv1.AddressAdvertise, cluster.Spec.Features.HA.TKEHA.VIP, 6443)
		}
		if cluster.Spec.Features.HA.ThirdPartyHA != nil {
			cluster.AddAddress(platformv1.AddressAdvertise, cluster.Spec.Features.HA.ThirdPartyHA.VIP, cluster.Spec.Features.HA.ThirdPartyHA.VPort)
		}
	}

	return nil
}

func completeCredential(cluster *v1.Cluster) error {
	token := ksuid.New().String()
	cluster.ClusterCredential.Token = &token

	bootstrapToken, err := bootstraputil.GenerateBootstrapToken()
	if err != nil {
		return err
	}
	cluster.ClusterCredential.BootstrapToken = &bootstrapToken

	certBytes := make([]byte, 32)
	if _, err := rand.Read(certBytes); err != nil {
		return err
	}
	certificateKey := hex.EncodeToString(certBytes)
	cluster.ClusterCredential.CertificateKey = &certificateKey

	return nil
}

func (p *Provider) EnsureKubeconfig(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		option := &kubeconfig.Option{
			MasterEndpoint: "https://127.0.0.1:6443",
			ClusterName:    c.Name,
			CACert:         c.ClusterCredential.CACert,
			Token:          *c.ClusterCredential.Token,
		}
		err = kubeconfig.Install(machineSSH, option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureNvidiaDriver(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		if !gpu.IsEnable(machine.Labels) {
			continue
		}
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = gpu.InstallNvidiaDriver(machineSSH, &gpu.NvidiaDriverOption{})
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureNvidiaContainerRuntime(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		if !gpu.IsEnable(machine.Labels) {
			continue
		}
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = gpu.InstallNvidiaContainerRuntime(machineSSH, &gpu.NvidiaContainerRuntimeOption{})
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureDocker(ctx context.Context, c *v1.Cluster) error {
	insecureRegistries := fmt.Sprintf(`"%s"`, p.config.Registry.Domain)
	if p.config.Registry.NeedSetHosts() && c.Spec.TenantID != "" {
		insecureRegistries = fmt.Sprintf(`%s,"%s"`, insecureRegistries, c.Spec.TenantID+"."+p.config.Registry.Domain)
	}
	option := &docker.Option{
		InsecureRegistries: insecureRegistries,
		RegistryDomain:     p.config.Registry.Domain,
		ExtraArgs:          c.Spec.DockerExtraArgs,
	}
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		option.IsGPU = gpu.IsEnable(machine.Labels)
		err = docker.Install(machineSSH, option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureConntrackTools(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = res.ConntrackTools.InstallWithDefault(machineSSH)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureKubeadm(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Install(machineSSH, c.Spec.Version)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsurePrepareForControlplane(ctx context.Context, c *v1.Cluster) error {
	oidcCa, _ := ioutil.ReadFile(path.Join(constants.ConfDir, constants.OIDCCACertName))
	auditPolicyData, _ := ioutil.ReadFile(path.Join(constants.ConfDir, constants.AuditPolicyConfigFile))
	GPUQuotaAdmissionHost := c.Annotations[constants.GPUQuotaAdmissionIPAnnotaion]
	if GPUQuotaAdmissionHost == "" {
		GPUQuotaAdmissionHost = "gpu-quota-admission"
	}
	schedulerPolicyConfig, err := template.ParseString(schedulerPolicyConfig, map[string]interface{}{
		"GPUQuotaAdmissionHost": GPUQuotaAdmissionHost,
	})
	if err != nil {
		return errors.Wrap(err, "parse schedulerPolicyConfig error")
	}
	auditWebhookConfig, err := template.ParseString(auditWebhookConfig, map[string]interface{}{
		"AuditBackendAddress": p.config.Audit.Address,
		"ClusterName":         c.Name,
	})
	if err != nil {
		return errors.Wrap(err, "parse auditWebhookConfig error")
	}
	for i, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		tokenData := fmt.Sprintf(tokenFileTemplate, *c.ClusterCredential.Token)
		err = machineSSH.WriteFile(strings.NewReader(tokenData), constants.TokenFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		err = machineSSH.WriteFile(bytes.NewReader(schedulerPolicyConfig), constants.SchedulerPolicyConfigFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		if len(oidcCa) != 0 {
			err = machineSSH.WriteFile(bytes.NewReader(oidcCa), constants.OIDCCACertFile)
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
		}

		if c.Spec.Features.HA != nil {
			if c.Spec.Features.HA.TKEHA != nil {
				networkInterface := ssh.GetNetworkInterface(machineSSH, machine.IP)
				if networkInterface == "" {
					return fmt.Errorf("can't get network interface by %s", machine.IP)
				}

				data, err := template.ParseFile(constants.ManifestsDir+"keepalived/keepalived.conf", map[string]interface{}{
					"Interface": networkInterface,
					"VIP":       c.Spec.Features.HA.TKEHA.VIP,
				})
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}
				err = machineSSH.WriteFile(bytes.NewReader(data), constants.KeepavliedConfigFile)
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}

				data, err = template.ParseFile(constants.ManifestsDir+"keepalived/keepalived.yaml", map[string]interface{}{
					"Image": images.Get().Keepalived.FullName(),
				})
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}
				err = machineSSH.WriteFile(bytes.NewReader(data), constants.KeepavlivedManifestFile)
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}
			}
			if c.Spec.Features.HA.ThirdPartyHA != nil && i > 0 { // forward rest control-plane to first master
				cmd := fmt.Sprintf("iptables -t nat -I PREROUTING 1 -p tcp --dport 6443 -j DNAT --to-destination %s:6443",
					c.Spec.Machines[0].IP)
				_, stderr, exit, err := machineSSH.Exec(cmd)
				if err != nil || exit != 0 {
					return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
				}

				cmd = fmt.Sprintf("iptables -t nat -I POSTROUTING 1 -p tcp -d %s --dport 6443 -j SNAT --to-source %s",
					c.Spec.Machines[0].IP, machine.IP)
				_, stderr, exit, err = machineSSH.Exec(cmd)
				if err != nil || exit != 0 {
					return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
				}
			}
		}

		if p.config.Audit.Address != "" {
			if len(auditPolicyData) != 0 {
				err = machineSSH.WriteFile(bytes.NewReader(auditPolicyData), path.Join(constants.KubernetesDir, constants.AuditPolicyConfigFile))
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}
				err = machineSSH.WriteFile(bytes.NewReader(auditWebhookConfig), path.Join(constants.KubernetesDir, constants.AuditWebhookConfigFile))
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}
			}
		}
	}

	return nil
}

func (p *Provider) EnsureKubeadmInitPhaseKubeletStart(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c),
		fmt.Sprintf("kubelet-start --node-name=%s", c.Spec.Machines[0].IP))
}

func (p *Provider) EnsureKubeadmInitPhaseCerts(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c), "certs all")
}

func (p *Provider) EnsureKubeadmInitPhaseKubeConfig(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c), "kubeconfig all")
}

func (p *Provider) EnsureKubeadmInitPhaseControlPlane(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c), "control-plane all")
}

func (p *Provider) EnsureKubeadmInitPhaseETCD(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c), "etcd local")
}

func (p *Provider) EnsureKubeadmInitPhaseUploadConfig(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c), "upload-config all ")
}

func (p *Provider) EnsureKubeadmInitPhaseUploadCerts(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c), "upload-certs --upload-certs")
}

func (p *Provider) EnsureKubeadmInitPhaseBootstrapToken(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c), "bootstrap-token")
}

func (p *Provider) EnsureKubeadmInitPhaseAddon(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, p.getKubeadmInitConfig(c), "addon all")
}

func (p *Provider) EnsureGalaxy(ctx context.Context, c *v1.Cluster) error {
	clientset, err := c.ClientsetForBootstrap()
	if err != nil {
		return err
	}
	return galaxy.Install(ctx, clientset, &galaxy.Option{
		Version:   galaxyimages.LatestVersion,
		NodeCIDR:  c.Cluster.Spec.ClusterCIDR,
		NetDevice: c.Cluster.Spec.NetworkDevice,
	})
}

func (p *Provider) EnsureJoinPhasePreflight(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines[1:] {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "preflight")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseControlPlanePrepare(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines[1:] {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "control-plane-prepare all")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseKubeletStart(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines[1:] {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "kubelet-start")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseControlPlaneJoinETCD(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines[1:] {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "control-plane-join etcd")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseControlPlaneJoinUpdateStatus(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines[1:] {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "control-plane-join update-status")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureStoreCredential(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}

	data, err := machineSSH.ReadFile(constants.CACertName)
	if err != nil {
		return err
	}
	c.ClusterCredential.CACert = data

	data, err = machineSSH.ReadFile(constants.CAKeyName)
	if err != nil {
		return err
	}
	c.ClusterCredential.CAKey = data

	data, err = machineSSH.ReadFile(constants.EtcdCACertName)
	if err != nil {
		return err
	}
	c.ClusterCredential.ETCDCACert = data

	data, err = machineSSH.ReadFile(constants.EtcdCAKeyName)
	if err != nil {
		return err
	}
	c.ClusterCredential.ETCDCAKey = data

	data, err = machineSSH.ReadFile(constants.APIServerEtcdClientCertName)
	if err != nil {
		return err
	}
	c.ClusterCredential.ETCDAPIClientCert = data

	data, err = machineSSH.ReadFile(constants.APIServerEtcdClientKeyName)
	if err != nil {
		return err
	}
	c.ClusterCredential.ETCDAPIClientKey = data

	return nil
}

func (p *Provider) EnsurePatchAnnotation(ctx context.Context, c *v1.Cluster) error {
	fileData := map[string]string{
		constants.EtcdPodManifestFile:                  `  annotations:\n    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "https"\n    prometheus.io/port: "2379"`,
		constants.KubeAPIServerPodManifestFile:         `  annotations:\n    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "https"\n    prometheus.io/port: "6443"`,
		constants.KubeControllerManagerPodManifestFile: `  annotations:\n    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "http"\n    prometheus.io/port: "10252"`,
		constants.KubeSchedulerPodManifestFile:         `  annotations:\n    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "http"\n    prometheus.io/port: "10251"`,
	}
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		for file, data := range fileData {
			cmd := fmt.Sprintf(`grep 'prometheus.io/port' %s || sed -i '3a\%s' %s`, file, data, file)
			_, stderr, exit, err := machineSSH.Exec(cmd)
			if err != nil || exit != 0 {
				return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
			}
		}
	}

	return nil
}

func (p *Provider) EnsureKubelet(ctx context.Context, c *v1.Cluster) error {
	option := &kubelet.Option{
		Version:   c.Spec.Version,
		ExtraArgs: c.Spec.KubeletExtraArgs,
	}
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubelet.Install(machineSSH, option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureCNIPlugins(ctx context.Context, c *v1.Cluster) error {
	option := &cniplugins.Option{}
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = cniplugins.Install(machineSSH, option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureKubeadmInitPhaseWaitControlPlane(ctx context.Context, c *v1.Cluster) error {
	start := time.Now()

	return wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		healthStatus := 0
		clientset, err := c.ClientsetForBootstrap()
		if err != nil {
			log.Warn(err.Error())
			return false, nil
		}
		clientset.Discovery().RESTClient().Get().AbsPath("/healthz").Do(ctx).StatusCode(&healthStatus)
		if healthStatus != http.StatusOK {
			return false, nil
		}

		log.Infof("All control plane components are healthy after %f seconds\n", time.Since(start).Seconds())
		return true, nil
	})
}

func (p *Provider) EnsureMarkControlPlane(ctx context.Context, c *v1.Cluster) error {
	clientset, err := c.ClientsetForBootstrap()
	if err != nil {
		return err
	}

	for _, machine := range c.Spec.Machines {
		if machine.Labels == nil {
			machine.Labels = make(map[string]string)
		}
		machine.Labels[constants.LabelNodeRoleMaster] = ""

		if !c.Spec.Features.EnableMasterSchedule {
			taint := corev1.Taint{
				Key:    constants.LabelNodeRoleMaster,
				Effect: corev1.TaintEffectNoSchedule,
			}
			if !funk.Contains(machine.Taints, taint) {
				machine.Taints = append(machine.Taints, taint)
			}
		}
		err := apiclient.MarkNode(ctx, clientset, machine.IP, machine.Labels, machine.Taints)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureNvidiaDevicePlugin(ctx context.Context, c *v1.Cluster) error {
	if c.Cluster.Spec.Features.GPUType == nil {
		return nil
	}

	client, err := c.ClientsetForBootstrap()
	if err != nil {
		return err
	}
	option := &gpu.NvidiaDevicePluginOption{
		Image: images.Get().NvidiaDevicePlugin.FullName(),
	}
	err = gpu.InstallNvidiaDevicePlugin(ctx, client, option)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureGPUManager(ctx context.Context, c *v1.Cluster) error {
	if c.Cluster.Spec.Features.GPUType == nil {
		return nil
	}

	if *c.Cluster.Spec.Features.GPUType != platformv1.GPUVirtual {
		return nil
	}

	client, err := c.Clientset()
	if err != nil {
		return err
	}

	option := map[string]interface{}{
		"GPUManagerImage":        images.Get().GPUManager.FullName(),
		"BusyboxImage":           images.Get().Busybox.FullName(),
		"GPUQuotaAdmissionImage": images.Get().GPUQuotaAdmission.FullName(),
	}

	err = apiclient.CreateResourceWithFile(ctx, client, constants.GPUManagerManifest, option)
	if err != nil {
		return errors.Wrap(err, "install gpu manager error")
	}

	return nil
}

func (p *Provider) EnsureCSIOperator(ctx context.Context, c *v1.Cluster) error {
	if c.Cluster.Spec.Features.CSIOperator == nil {
		return nil
	}

	log.Info("csi-perator in feature will be created.", log.String("name", c.Cluster.Name))

	client, err := c.Clientset()
	if err != nil {
		return err
	}

	option := map[string]interface{}{
		"CSIOperatorImage": csioperatorimage.Get(c.Cluster.Spec.Features.CSIOperator.Version).CSIOperator.FullName(),
		"RegistryDomain":   containerregistryutil.GetPrefix(),
	}

	err = apiclient.CreateResourceWithFile(ctx, client, constants.CSIOperatorManifest, option)
	if err != nil {
		return errors.Wrap(err, "install csi operator error")
	}

	log.Info("csi-perator in feature already created.", log.String("name", c.Cluster.Name))

	return nil
}

func (p *Provider) EnsureCleanup(ctx context.Context, c *v1.Cluster) error {
	for i, machine := range c.Spec.Machines {
		s, err := machine.SSH()
		if err != nil {
			return err
		}

		if c.Spec.Features.HA != nil {
			if c.Spec.Features.HA.ThirdPartyHA != nil && i > 0 {
				cmd := fmt.Sprintf("iptables -t nat -D PREROUTING -p tcp --dport 6443 -j DNAT --to-destination %s:6443",
					c.Spec.Machines[0].IP)
				_, stderr, exit, err := s.Exec(cmd)
				if err != nil || exit != 0 {
					return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
				}

				cmd = fmt.Sprintf("iptables -t nat -D POSTROUTING -p tcp -d %s --dport 6443 -j SNAT --to-source %s",
					c.Spec.Machines[0].IP, machine.IP)
				_, stderr, exit, err = s.Exec(cmd)
				if err != nil || exit != 0 {
					return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
				}
			}
		}

	}
	return nil
}
