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

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/thoas/go-funk"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
	kubeaggregatorclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/addons/cniplugins"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/authzwebhook"
	csioperatorimage "tkestack.io/tke/pkg/platform/provider/baremetal/phases/csioperator/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/docker"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy"
	galaxyimages "tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/gpu"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/image"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/keepalived"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeconfig"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubelet"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/thirdpartyha"
	"tkestack.io/tke/pkg/platform/provider/baremetal/preflight"
	"tkestack.io/tke/pkg/platform/provider/baremetal/res"
	"tkestack.io/tke/pkg/platform/provider/util/mark"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/cmdstring"
	containerregistryutil "tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/hosts"
	"tkestack.io/tke/pkg/util/log"
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
	var data bytes.Buffer
	for _, machine := range c.Spec.Machines {
		modules := []string{"iptable_nat", "ip_vs", "ip_vs_rr", "ip_vs_wrr", "ip_vs_sh"}

		s, err := machine.SSH()
		if err != nil {
			return err
		}
		if _, err := s.CombinedOutput("modinfo br_netfilter"); err == nil {
			modules = append(modules, "br_netfilter")
		}

		for _, m := range modules {
			_, err := s.CombinedOutput(fmt.Sprintf("modprobe %s", m))
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
			data.WriteString(m + "\n")
		}
		err = s.WriteFile(strings.NewReader(data.String()), moduleFile)
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

		_, err = machineSSH.CombinedOutput(`swapoff -a && sed -i "s/^[^#]*swap/#&/" /etc/fstab`)
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
		completeServiceIP,
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

func kubernetesSvcIP(cluster *v1.Cluster) (string, error) {
	ip, err := GetIndexedIP(cluster.Status.ServiceCIDR, constants.KUBERNETES)
	if err != nil {
		return "", errors.Wrap(err, "get kubernetesSvcIP error")
	}

	return ip.String(), nil
}

func completeDNS(cluster *v1.Cluster) error {
	ip, err := GetIndexedIP(cluster.Status.ServiceCIDR, constants.DNSIPIndex)
	if err != nil {
		return errors.Wrap(err, "get DNS IP error")
	}
	cluster.Status.DNSIP = ip.String()

	return nil
}

func completeServiceIP(cluster *v1.Cluster) error {
	if cluster.Annotations == nil {
		cluster.Annotations = make(map[string]string)
	}
	for index, name := range map[int]string{
		constants.GPUQuotaAdmissionIPIndex: constants.GPUQuotaAdmissionIPAnnotaion,
		constants.GalaxyIPAMIPIndex:        constants.GalaxyIPAMIPIndexAnnotaion,
	} {
		ip, err := GetIndexedIP(cluster.Status.ServiceCIDR, index)
		if err != nil {
			return errors.Wrap(err, "get service IP error")
		}
		cluster.Annotations[name] = ip.String()
	}

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
	extraArgs := c.Spec.DockerExtraArgs
	utilruntime.Must(mergo.Merge(&extraArgs, p.config.Docker.ExtraArgs))
	option := &docker.Option{
		InsecureRegistries: insecureRegistries,
		RegistryDomain:     p.config.Registry.Domain,
		ExtraArgs:          extraArgs,
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

func (p *Provider) EnsureKubernetesImages(ctx context.Context, c *v1.Cluster) error {
	option := &image.Option{Version: c.Spec.Version, RegistryDomain: p.config.Registry.Domain}
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		err = image.PullKubernetesImages(machineSSH, option)
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

// EnsureKeepalivedInit make sure all master node has cleaning iptable table so in kubeadm join time apiserver may not join it self.
// keepalived only installs in master node 0 before kubeadm init phase to prevet from vip failover in kubeadm join(etcd phase)
func (p *Provider) EnsureKeepalivedInit(ctx context.Context, c *v1.Cluster) error {
	if c.Spec.Features.HA == nil || c.Spec.Features.HA.TKEHA == nil {
		return nil
	}

	kubernetesSvcIP, err := kubernetesSvcIP(c)
	if err != nil {
		return err
	}

	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		keepalived.ClearLoadBalance(machineSSH, c.Spec.Features.HA.TKEHA.VIP, kubernetesSvcIP)
	}

	option := &keepalived.Option{
		IP:              c.Spec.Machines[0].IP,
		VIP:             c.Spec.Features.HA.TKEHA.VIP,
		LoadBalance:     false,
		IPVS:            false,
		KubernetesSvcIP: kubernetesSvcIP,
	}

	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}

	err = keepalived.Install(machineSSH, option)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) EnsureThirdPartyHAInit(ctx context.Context, c *v1.Cluster) error {
	if c.Spec.Features.HA == nil || c.Spec.Features.HA.ThirdPartyHA == nil {
		return nil
	}

	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		option := thirdpartyha.Option{
			IP:    machine.IP,
			VIP:   c.Spec.Features.HA.ThirdPartyHA.VIP,
			VPort: c.Spec.Features.HA.ThirdPartyHA.VPort,
		}

		thirdpartyha.Clear(machineSSH, &option)
	}

	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}

	option := thirdpartyha.Option{
		IP:    c.Spec.Machines[0].IP,
		VIP:   c.Spec.Features.HA.ThirdPartyHA.VIP,
		VPort: c.Spec.Features.HA.ThirdPartyHA.VPort,
	}

	err = thirdpartyha.Install(machineSSH, &option)
	if err != nil {
		return err
	}

	return nil
}
func (p *Provider) EnsureAuthzWebhook(ctx context.Context, c *v1.Cluster) error {
	if !c.AuthzWebhookEnabled() {
		return nil
	}

	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		authzEndpoint, ok := c.AuthzWebhookExternEndpoint()
		if !ok {
			authzEndpoint = p.config.AuthzWebhook.Endpoint
		}

		option := authzwebhook.Option{AuthzWebhookEndpoint: authzEndpoint}
		err = authzwebhook.Install(machineSSH, &option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsurePrepareForControlplane(ctx context.Context, c *v1.Cluster) error {
	oidcCa, _ := ioutil.ReadFile(constants.OIDCConfigFile)
	auditPolicyData, _ := ioutil.ReadFile(constants.AuditPolicyConfigFile)
	GPUQuotaAdmissionHost := c.Annotations[constants.GPUQuotaAdmissionIPAnnotaion]
	if GPUQuotaAdmissionHost == "" {
		GPUQuotaAdmissionHost = "gpu-quota-admission"
	}
	GalaxyIPAMHost := c.Annotations[constants.GalaxyIPAMIPIndexAnnotaion]
	if GalaxyIPAMHost == "" {
		GalaxyIPAMHost = "galaxy-ipam"
	}
	schedulerPolicyConfig, err := template.ParseString(schedulerPolicyConfig, map[string]interface{}{
		"GPUQuotaAdmissionHost": GPUQuotaAdmissionHost,
		"GalaxyIPAMHost":        GalaxyIPAMHost,
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
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		tokenData := fmt.Sprintf(tokenFileTemplate, *c.ClusterCredential.Token)
		err = machineSSH.WriteFile(strings.NewReader(tokenData), constants.TokenFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		err = machineSSH.WriteFile(bytes.NewReader(schedulerPolicyConfig), constants.KuberentesSchedulerPolicyConfigFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		if len(oidcCa) != 0 {
			err = machineSSH.WriteFile(bytes.NewReader(oidcCa), constants.OIDCCACertFile)
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
		}

		if p.config.AuditEnabled() {
			if len(auditPolicyData) != 0 {
				err = machineSSH.WriteFile(bytes.NewReader(auditPolicyData), constants.KubernetesAuditPolicyConfigFile)
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}
				err = machineSSH.WriteFile(bytes.NewReader(auditWebhookConfig), constants.KuberentesAuditWebhookConfigFile)
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
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubelet.Install(machineSSH, c.Spec.Version)
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
	return wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		clientset, err := c.ClientsetForBootstrap()
		if err != nil {
			log.FromContext(ctx).Error(err, "Create clientset error")
			return false, nil
		}
		result := clientset.Discovery().RESTClient().Get().AbsPath("/healthz").Do(ctx)
		statusCode := 0
		result.StatusCode(&statusCode)
		if statusCode != http.StatusOK {
			log.FromContext(ctx).Error(result.Error(), "check healthz error", "statusCode", statusCode)
			return false, nil
		}

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
		"GPUQuotaAdmissionHost":  c.Annotations[constants.GPUQuotaAdmissionIPAnnotaion],
	}

	err = apiclient.CreateResourceWithFile(ctx, client, constants.GPUManagerManifest, option)
	if err != nil {
		return errors.Wrap(err, "install gpu manager error")
	}

	return nil
}

func (p *Provider) EnsureMetricsServer(ctx context.Context, c *v1.Cluster) error {
	client, err := c.Clientset()
	if err != nil {
		return err
	}
	config, err := c.RESTConfig(&rest.Config{})
	if err != nil {
		return err
	}
	kaClient, err := kubeaggregatorclientset.NewForConfig(config)
	if err != nil {
		return err
	}
	option := map[string]interface{}{
		"MetricsServerImage": images.Get().MetricsServer.FullName(),
		"AddonResizerImage":  images.Get().AddonResizer.FullName(),
	}

	err = apiclient.CreateKAResourceWithFile(ctx, client, kaClient, constants.MetricsServerManifest, option)
	if err != nil {
		return errors.Wrap(err, "install metrics server error")
	}

	return nil
}

func (p *Provider) EnsureCSIOperator(ctx context.Context, c *v1.Cluster) error {
	if c.Cluster.Spec.Features.CSIOperator == nil {
		return nil
	}

	log.FromContext(ctx).Info("csi-perator will be created")

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
		return errors.Wrap(err, "install csi-operator error")
	}

	log.FromContext(ctx).Info("csi-perator already created")

	return nil
}

func (p *Provider) EnsureKeepalivedWithLB(ctx context.Context, c *v1.Cluster) error {
	if c.Spec.Features.HA == nil || c.Spec.Features.HA.TKEHA == nil {
		return nil
	}

	ipvs := c.Spec.Features.IPVS != nil && *c.Spec.Features.IPVS
	kubernetesSvcIP, err := kubernetesSvcIP(c)
	if err != nil {
		return err
	}

	for _, machine := range c.Spec.Machines {
		s, err := machine.SSH()
		if err != nil {
			return err
		}

		option := &keepalived.Option{
			IP:              machine.IP,
			VIP:             c.Spec.Features.HA.TKEHA.VIP,
			LoadBalance:     true,
			IPVS:            ipvs,
			KubernetesSvcIP: kubernetesSvcIP,
		}

		err = keepalived.Install(s, option)
		if err != nil {
			return err
		}

		log.FromContext(ctx).Info("keepalived created success.", "node", machine.IP)
	}

	return nil
}

func (p *Provider) EnsureThirdPartyHA(ctx context.Context, c *v1.Cluster) error {
	if c.Spec.Features.HA == nil || c.Spec.Features.HA.ThirdPartyHA == nil {
		return nil
	}

	for _, machine := range c.Spec.Machines {
		s, err := machine.SSH()
		if err != nil {
			return err
		}

		option := thirdpartyha.Option{
			IP:    machine.IP,
			VIP:   c.Spec.Features.HA.ThirdPartyHA.VIP,
			VPort: c.Spec.Features.HA.ThirdPartyHA.VPort,
		}

		err = thirdpartyha.Install(s, &option)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) EnsureCleanup(ctx context.Context, c *v1.Cluster) error {
	for _, machine := range c.Spec.Machines {
		_, err := machine.SSH()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) EnsureCreateClusterMark(ctx context.Context, c *v1.Cluster) error {
	clientset, err := c.Clientset()
	if err != nil {
		return err
	}

	return mark.Create(ctx, clientset)
}
