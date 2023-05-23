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
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	appsv1alpha1 "github.com/clusternet/apis/apps/v1alpha1"
	clustersv1beta1 "github.com/clusternet/apis/clusters/v1beta1"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/thoas/go-funk"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
	kubeaggregatorclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	utilsnet "k8s.io/utils/net"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	application "tkestack.io/tke/api/application/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/addons/cniplugins"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/authzwebhook"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/containerd"
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
	"tkestack.io/tke/pkg/platform/provider/baremetal/util"
	"tkestack.io/tke/pkg/platform/provider/util/mark"
	v1 "tkestack.io/tke/pkg/platform/types/v1"
	"tkestack.io/tke/pkg/util/apiclient"
	"tkestack.io/tke/pkg/util/cmdstring"
	containerregistryutil "tkestack.io/tke/pkg/util/containerregistry"
	"tkestack.io/tke/pkg/util/extenderapi"
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, file := range c.Spec.Features.Files {
		for _, machine := range machines {
			machineSSH, err := machine.SSH()
			if err != nil {
				return err
			}
			s, err := os.Stat(file.Src)
			if err != nil {
				return err
			}
			if s.Mode().IsDir() {
				if err != nil {
					return err
				}
				err = machineSSH.CopyDir(file.Src, file.Dst)
				if err != nil {
					return err
				}
			} else {
				err = machineSSH.CopyFile(file.Src, file.Dst)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *Provider) EnsurePreClusterInstallHook(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	return util.ExcuteCustomizedHook(ctx, c, platformv1.HookPreClusterInstall, c.Spec.Machines[:1])
}

func (p *Provider) EnsurePreInstallHook(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	return util.ExcuteCustomizedHook(ctx, c, platformv1.HookPreInstall, machines)
}

func (p *Provider) EnsurePostInstallHook(ctx context.Context, c *v1.Cluster) error {

	return util.ExcuteCustomizedHook(ctx, c, platformv1.HookPostInstall, c.Spec.Machines)
}

func (p *Provider) EnsurePostClusterInstallHook(ctx context.Context, c *v1.Cluster) error {

	return util.ExcuteCustomizedHook(ctx, c, platformv1.HookPostClusterInstall, c.Spec.Machines[:1])
}

func (p *Provider) EnsurePreflight(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = preflight.RunMasterChecks(c, machineSSH)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureRegistryHosts(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	domains := []string{
		p.Config.Registry.Domain,
	}
	if c.Spec.TenantID != "" {
		domains = append(domains, c.Spec.TenantID+"."+p.Config.Registry.Domain)
	}
	domains = append(domains, constants.MirrorsRegistryHostName)
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		for _, one := range domains {
			remoteHosts := hosts.RemoteHosts{Host: one, SSH: machineSSH}
			ip := p.Config.Registry.IP
			if len(p.Config.Registry.IP) == 0 {
				ip = c.GetMainIP()
			}
			err := remoteHosts.Set(ip)
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
		}
	}

	return nil
}

func (p *Provider) EnsureKernelModule(ctx context.Context, c *v1.Cluster) error {
	var data bytes.Buffer
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
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
	if cluster.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
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
		clusterCIDR      = cluster.Spec.ClusterCIDR
		serviceCIDR      string
		nodeCIDRMaskSize int32
		err              error
	)

	// dual stack case
	if cluster.Spec.Features.IPv6DualStack {
		clusterCidrs := strings.Split(serviceCIDR, ",")
		serviceCidrs := strings.Split(clusterCIDR, ",")
		for _, cidr := range clusterCidrs {
			if maskSize, isIPv6 := CalcNodeCidrSize(cidr); isIPv6 {
				cluster.Status.NodeCIDRMaskSizeIPv6 = maskSize
				cluster.Status.SecondaryClusterCIDR = cidr
			} else {
				cluster.Status.NodeCIDRMaskSizeIPv4 = maskSize
				cluster.Status.ClusterCIDR = cidr
			}
		}
		for _, cidr := range serviceCidrs {
			if utilsnet.IsIPv6CIDRString(cidr) {
				cluster.Status.SecondaryServiceCIDR = cidr
			} else {
				cluster.Status.ServiceCIDR = cidr
			}
		}
		return nil
	}
	// single stack case incldue ipv4 and ipv6
	if cluster.Spec.ServiceCIDR != nil {
		serviceCIDR = *cluster.Spec.ServiceCIDR
		if utilsnet.IsIPv6CIDRString(clusterCIDR) {
			nodeCIDRMaskSize, _ = CalcNodeCidrSize(clusterCIDR)
		} else {
			nodeCIDRMaskSize, err = GetNodeCIDRMaskSize(clusterCIDR, *cluster.Spec.Properties.MaxNodePodNum)
			if err != nil {
				return errors.Wrap(err, "GetNodeCIDRMaskSize error")
			}
		}
	} else {
		serviceCIDR, nodeCIDRMaskSize, err = GetServiceCIDRAndNodeCIDRMaskSize(clusterCIDR, *cluster.Spec.Properties.MaxClusterServiceNum, *cluster.Spec.Properties.MaxNodePodNum)
		if err != nil {
			return errors.Wrap(err, "GetServiceCIDRAndNodeCIDRMaskSize error")
		}
	}

	cluster.Status.ServiceCIDR = serviceCIDR
	cluster.Status.ClusterCIDR = clusterCIDR
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
	return completePlatformClusterAddresses(cluster.Cluster)
}

func completePlatformClusterAddresses(cluster *platformv1.Cluster) error {
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
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

func (p *Provider) EnsureContainerRuntime(ctx context.Context, c *v1.Cluster) error {
	if c.Cluster.Spec.Features.ContainerRuntime == platformv1.Docker {
		return p.EnsureDocker(ctx, c)
	}
	return p.EnsureContainerd(ctx, c)
}

func (p *Provider) getImagePrefix(c *v1.Cluster) string {
	if anno, ok := c.Annotations[platformv1.LocationBasedImagePrefixAnno]; ok {
		return anno
	}
	return containerregistryutil.GetPrefix()
}

func (p *Provider) EnsureContainerd(ctx context.Context, c *v1.Cluster) error {
	insecureRegistries := []string{p.Config.Registry.Domain}
	if c.Spec.TenantID != "" {
		insecureRegistries = append(insecureRegistries, c.Spec.TenantID+"."+p.Config.Registry.Domain)
	}
	prefix := p.getImagePrefix(c)
	option := &containerd.Option{
		InsecureRegistries: insecureRegistries,
		SandboxImage:       path.Join(prefix, images.Get().Pause.BaseName()),
		// for mirror, we just need domain in prefix
		RegistryMirror: strings.Split(prefix, "/")[0],
	}
	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		option.IsGPU = gpu.IsEnable(machine.Labels)
		err = containerd.Install(machineSSH, option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureDocker(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	insecureRegistries := fmt.Sprintf(`"%s"`, p.Config.Registry.Domain)
	if c.Spec.TenantID != "" {
		insecureRegistries = fmt.Sprintf(`%s,"%s"`, insecureRegistries, c.Spec.TenantID+"."+p.Config.Registry.Domain)
	}
	extraArgs := c.Spec.DockerExtraArgs
	utilruntime.Must(mergo.Merge(&extraArgs, p.Config.Docker.ExtraArgs))
	option := &docker.Option{
		InsecureRegistries: insecureRegistries,
		RegistryDomain:     p.Config.Registry.Domain,
		ExtraArgs:          extraArgs,
	}
	for _, machine := range machines {
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	option := &image.Option{Version: c.Spec.Version, RegistryDomain: p.Config.Registry.Domain, KubeImages: images.KubecomponetNames}
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		err = image.PullKubernetesImages(c, machineSSH, option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureConntrackTools(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		option := &kubeadm.Option{
			RuntimeType: c.Spec.Features.ContainerRuntime,
			Version:     c.Spec.Version,
		}
		err = kubeadm.Install(machineSSH, option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

// EnsureKeepalivedInit make sure all master node has cleaning iptable table so in kubeadm join time apiserver may not join it self.
// keepalived only installs in master node 0 before kubeadm init phase to prevet from vip failover in kubeadm join(etcd phase)
func (p *Provider) EnsureKeepalivedInit(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	if c.Spec.Features.HA == nil || c.Spec.Features.HA.TKEHA == nil {
		return nil
	}

	kubernetesSvcIP, err := kubernetesSvcIP(c)
	if err != nil {
		return err
	}

	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		keepalived.ClearLoadBalance(machineSSH, c.Spec.Features.HA.TKEHA.VIP, kubernetesSvcIP)
	}
	if c.Status.Phase == platformv1.ClusterRunning {
		return nil
	}
	option := &keepalived.Option{
		IP:              c.Spec.Machines[0].IP,
		VIP:             c.Spec.Features.HA.TKEHA.VIP,
		VRID:            c.Spec.Features.HA.TKEHA.VRID,
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	if c.Spec.Features.HA == nil || c.Spec.Features.HA.ThirdPartyHA == nil {
		return nil
	}

	for _, machine := range machines {
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
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	if !c.AuthzWebhookEnabled() {
		return nil
	}
	isGlobalCluster := (c.Cluster.Name == "global")
	isClusterUpscaling := (c.Status.Phase == platformv1.ClusterUpscaling)
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		authzEndpoint, ok := c.AuthzWebhookExternEndpoint()
		if !ok {
			if isGlobalCluster {
				authzEndpoint, _ = c.AuthzWebhookBuiltinEndpoint()
			} else {
				authzEndpoint = p.Config.AuthzWebhook.Endpoint
			}
		}
		option := authzwebhook.Option{
			AuthzWebhookEndpoint: authzEndpoint,
			IsGlobalCluster:      isGlobalCluster,
			IsClusterUpscaling:   isClusterUpscaling,
		}
		err = authzwebhook.Install(machineSSH, &option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureAuditConfig(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	auditPolicyData, _ := ioutil.ReadFile(constants.AuditPolicyConfigFile)
	auditWebhookConfig, err := template.ParseString(auditWebhookConfig, map[string]interface{}{
		"AuditBackendAddress": p.Config.Audit.Address,
		"ClusterName":         c.Name,
	})
	if err != nil {
		return errors.Wrap(err, "parse auditWebhookConfig error")
	}
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}
		if p.Config.AuditEnabled() {
			if len(auditPolicyData) != 0 {
				err = machineSSH.WriteFile(bytes.NewReader(auditPolicyData), constants.KubernetesAuditPolicyConfigFile)
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}
				err = machineSSH.WriteFile(bytes.NewReader(auditWebhookConfig), constants.KubernetesAuditWebhookConfigFile)
				if err != nil {
					return errors.Wrap(err, machine.IP)
				}
			}
		}
	}

	return nil
}

func (p *Provider) EnsurePrepareForControlplane(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	oidcCa, _ := ioutil.ReadFile(constants.OIDCConfigFile)
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
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		tokenData := fmt.Sprintf(tokenFileTemplate, *c.ClusterCredential.Token)
		err = machineSSH.WriteFile(strings.NewReader(tokenData), constants.TokenFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		err = machineSSH.WriteFile(bytes.NewReader(schedulerPolicyConfig), constants.KubernetesSchedulerPolicyConfigFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		if len(oidcCa) != 0 {
			err = machineSSH.WriteFile(bytes.NewReader(oidcCa), constants.OIDCCACertFile)
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
		}
	}

	return nil
}

func (p *Provider) EnsureKubeadmInitPhaseKubeletStart(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	phase := "kubelet-start"
	kubeletExtraArgs := p.getKubeletExtraArgs(c)
	if _, ok := kubeletExtraArgs["hostname-override"]; !ok {
		if !c.Spec.HostnameAsNodename {
			phase += fmt.Sprintf(" --node-name=%s", c.Spec.Machines[0].IP)
		}
	}
	err = kubeadm.WriteInitConfig(machineSSH, p.getKubeadmInitConfig(c))
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, phase)
}

func (p *Provider) EnsureKubeadmInitPhaseCerts(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, "certs all")
}

func (p *Provider) EnsureKubeadmInitPhaseKubeConfig(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, "kubeconfig all")
}

func (p *Provider) EnsureKubeadmInitPhaseControlPlane(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, "control-plane all")
}

func (p *Provider) EnsureKubeadmInitPhaseETCD(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, "etcd local")
}

func (p *Provider) EnsureKubeadmInitPhaseUploadConfig(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, "upload-config all ")
}

func (p *Provider) EnsureKubeadmInitPhaseUploadCerts(ctx context.Context, c *v1.Cluster) error {
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, "upload-certs --upload-certs")
}

func (p *Provider) EnsureKubeadmInitPhaseBootstrapToken(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, "bootstrap-token")
}

func (p *Provider) EnsureKubeadmInitPhaseAddon(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	return kubeadm.Init(machineSSH, "addon all")
}

func (p *Provider) EnsureGalaxy(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	clientset, err := c.ClientsetForBootstrap()
	if err != nil {
		return err
	}
	backendType := "vxlan"
	clusterSpec := c.Cluster.Spec
	if clusterSpec.NetworkArgs != nil {
		backendTypeArg, ok := clusterSpec.NetworkArgs["backendType"]
		if ok {
			backendType = backendTypeArg
		}
	}
	return galaxy.Install(ctx, clientset, &galaxy.Option{
		Version:     galaxyimages.LatestVersion,
		NodeCIDR:    clusterSpec.ClusterCIDR,
		NetDevice:   clusterSpec.NetworkDevice,
		BackendType: backendType,
	})
}

func (p *Provider) clusterMachineIPs(c *v1.Cluster) []string {
	ips := []string{}
	for _, mc := range c.Spec.Machines {
		ips = append(ips, mc.IP)
	}
	return ips
}

func (p *Provider) EnsureJoinPhasePreflight(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines[1:]}[len(c.Spec.ScalingMachines) > 0]

	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "preflight", p.clusterMachineIPs(c))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseControlPlanePrepare(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines[1:]}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "control-plane-prepare all", p.clusterMachineIPs(c))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseKubeletStart(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines[1:]}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "kubelet-start", p.clusterMachineIPs(c))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseControlPlaneJoinETCD(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines[1:]}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "control-plane-join etcd", p.clusterMachineIPs(c))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureJoinPhaseControlPlaneJoinUpdateStatus(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines[1:]}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		err = kubeadm.Join(machineSSH, p.getKubeadmJoinConfig(c, machine.IP), "control-plane-join update-status", p.clusterMachineIPs(c))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureStoreCredential(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}

	data, err := machineSSH.ReadFile(constants.CACertName)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(c.ClusterCredential.CACert, data) {
		c.ClusterCredential.CACert = data
		c.IsCredentialChanged = true
	}

	data, err = machineSSH.ReadFile(constants.CAKeyName)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(c.ClusterCredential.CAKey, data) {
		c.ClusterCredential.CAKey = data
		c.IsCredentialChanged = true
	}

	data, err = machineSSH.ReadFile(constants.EtcdCACertName)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(c.ClusterCredential.ETCDCACert, data) {
		c.ClusterCredential.ETCDCACert = data
		c.IsCredentialChanged = true
	}

	data, err = machineSSH.ReadFile(constants.EtcdCAKeyName)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(c.ClusterCredential.ETCDCAKey, data) {
		c.ClusterCredential.ETCDCAKey = data
		c.IsCredentialChanged = true
	}

	data, err = machineSSH.ReadFile(constants.APIServerEtcdClientCertName)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(c.ClusterCredential.ETCDAPIClientCert, data) {
		c.ClusterCredential.ETCDAPIClientCert = data
		c.IsCredentialChanged = true
	}

	data, err = machineSSH.ReadFile(constants.APIServerEtcdClientKeyName)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(c.ClusterCredential.ETCDAPIClientKey, data) {
		c.ClusterCredential.ETCDAPIClientKey = data
		c.IsCredentialChanged = true
	}

	if c.IsCredentialChanged {
		c.RegisterRestConfig(c.ClusterCredential.RESTConfig(c.Cluster))
	}

	return nil
}

func (p *Provider) EnsurePatchAnnotation(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]

	// from k8s 1.18, kubeadm will add built-in annotations to etcd and kube-apiserver
	// we should handle such case when add tkestack annotations according to different case
	prefix := `  annotations:\n`
	cmdTpl := `
		idx=3
		yaml='%s'
		annotations='%s'
		if line=$(grep "annotations" -n ${yaml});then
			if grep -q "tke.prometheus.io/scrape" ${yaml};then
				exit
			else
				idx=$(echo $line | cut -d":" -f1)
				annotations='%s'
			fi
		fi
		sed -i "${idx}a\\${annotations}" ${yaml}`
	fileData := map[string]string{
		constants.EtcdPodManifestFile:                  `    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "https"\n    prometheus.io/port: "2379"`,
		constants.KubeAPIServerPodManifestFile:         `    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "https"\n    prometheus.io/port: "6443"`,
		constants.KubeControllerManagerPodManifestFile: `    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "https"\n    prometheus.io/port: "10257"`,
		constants.KubeSchedulerPodManifestFile:         `    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "https"\n    prometheus.io/port: "10259"`,
	}
	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		for file, data := range fileData {
			cmd := fmt.Sprintf(cmdTpl, file, prefix+data, data)
			_, stderr, exit, err := machineSSH.Exec(cmd)
			if err != nil || exit != 0 {
				return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
			}
		}
	}

	return nil
}

func (p *Provider) EnsureKubelet(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	for _, machine := range machines {
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]
	option := &cniplugins.Option{}
	for _, machine := range machines {
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
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	var exit int
	var stderr string
	_ = wait.PollImmediate(5*time.Second, 3*time.Minute, func() (bool, error) {
		cmd := "kubectl cluster-info"
		_, stderr, exit, err = machineSSH.Exec(cmd)
		if err != nil {
			err = fmt.Errorf("check apiserver failed: exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
			return false, nil
		}

		return true, nil
	})
	return err
}

func (p *Provider) EnsureMarkControlPlane(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]

	clientset, err := c.ClientsetForBootstrap()
	if err != nil {
		return err
	}

	for _, machine := range machines {
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
		node, err := apiclient.GetNodeByMachineIP(ctx, clientset, machine.IP)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
		err = apiclient.MarkNode(ctx, clientset, node.Name, machine.Labels, machine.Taints)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureNvidiaDevicePlugin(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}

	if c.Cluster.Spec.Features.GPUType == nil || *c.Cluster.Spec.Features.GPUType == platformv1.GPUVirtual {
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
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}

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
		"ContainerRuntime":       c.Spec.Features.ContainerRuntime,
	}

	err = apiclient.CreateResourceWithFile(ctx, client, constants.GPUManagerManifest, option)
	if err != nil {
		return errors.Wrap(err, "install gpu manager error")
	}

	return nil
}

func (p *Provider) EnsureMetricsServer(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	if !c.Cluster.Spec.Features.EnableMetricsServer {
		return nil
	}
	client, err := c.Clientset()
	if err != nil {
		return err
	}
	config, err := c.RESTConfig()
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

func (p *Provider) EnsureCilium(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}
	if !c.Cluster.Spec.Features.EnableCilium {
		return nil
	}
	// old cilium interface should be deleted
	if err := util.CleanFlannelInterfaces("cilium"); err != nil {
		return err
	}
	client, err := c.Clientset()
	if err != nil {
		return err
	}
	// default networkMode is overlay
	networkMode := "overlay"
	clusterSpec := c.Cluster.Spec
	if clusterSpec.NetworkArgs != nil {
		if networkTypeArg, ok := clusterSpec.NetworkArgs["networkMode"]; ok {
			networkMode = networkTypeArg
		}
	}
	option := map[string]interface{}{
		"CiliumImage":         images.Get().Cilium.FullName(),
		"CiliumOperatorImage": images.Get().CiliumOperator.FullName(),
		"IpamdImage":          images.Get().Ipamd.FullName(),
		"MasqImage":           images.Get().Masq.FullName(),
		"CiliumRouterImage":   images.Get().CiliumRouter.FullName(),
		"NetworkMode":         networkMode,
		"ClusterCIDR":         c.Cluster.Spec.ClusterCIDR,
		"MaskSize":            c.Cluster.Status.NodeCIDRMaskSize,
		"MaxNodePodNum":       c.Cluster.Spec.Properties.MaxNodePodNum,
	}

	err = apiclient.CreateResourceWithDir(ctx, client, constants.CiliumManifest, option)
	if err != nil {
		return errors.Wrap(err, "install Cilium error")
	}

	return nil
}

func (p *Provider) EnsureCSIOperator(ctx context.Context, c *v1.Cluster) error {
	if c.Status.Phase == platformv1.ClusterUpscaling {
		return nil
	}

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

func (p *Provider) EnsureKeepalivedWithLBOption(ctx context.Context, c *v1.Cluster) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]

	if c.Spec.Features.HA == nil || c.Spec.Features.HA.TKEHA == nil {
		return nil
	}

	ipvs := c.Spec.Features.IPVS != nil && *c.Spec.Features.IPVS
	kubernetesSvcIP, err := kubernetesSvcIP(c)
	if err != nil {
		return err
	}

	for _, machine := range machines {
		s, err := machine.SSH()
		if err != nil {
			return err
		}

		option := &keepalived.Option{
			IP:              machine.IP,
			VIP:             c.Spec.Features.HA.TKEHA.VIP,
			VRID:            c.Spec.Features.HA.TKEHA.VRID,
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]

	if c.Spec.Features.HA == nil || c.Spec.Features.HA.ThirdPartyHA == nil {
		return nil
	}

	for _, machine := range machines {
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
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]

	for _, machine := range machines {
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

func (p *Provider) setAPIServerHost(ctx context.Context, c *v1.Cluster, ip string) error {
	machines := map[bool][]platformv1.ClusterMachine{
		true:  c.Spec.ScalingMachines,
		false: c.Spec.Machines}[len(c.Spec.ScalingMachines) > 0]

	for _, machine := range machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		remoteHosts := hosts.RemoteHosts{Host: constants.APIServerHostName, SSH: machineSSH}
		err = remoteHosts.Set(ip)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}
	return nil
}

func (p *Provider) EnsureInitAPIServerHost(ctx context.Context, c *v1.Cluster) error {
	// Set host to master0 IP firstly, when local apiserver is at work, modify the host to 127.0.0.1.
	return p.setAPIServerHost(ctx, c, c.Spec.Machines[0].IP)
}

func (p *Provider) EnsureModifyAPIServerHost(ctx context.Context, c *v1.Cluster) error {
	return p.setAPIServerHost(ctx, c, "127.0.0.1")
}

func (p *Provider) EnsureClusternetRegistration(ctx context.Context, c *v1.Cluster) error {
	if c.Annotations[platformv1.RegistrationCommandAnno] == "" {
		log.FromContext(ctx).Info("registration command is empty, skip EnsureClusternetRegistration")
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(c.Annotations[platformv1.RegistrationCommandAnno])
	if err != nil {
		return fmt.Errorf("decode registration command failed: %v", err)
	}
	machineSSH, err := c.Spec.Machines[0].SSH()
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("echo \"%s\" | kubectl apply -f -", string(data))
	_, err = machineSSH.CombinedOutput(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) EnsureAnywhereEdtion(ctx context.Context, c *v1.Cluster) error {
	if c.Labels[platformv1.AnywhereEdtionLabel] == "" {
		log.FromContext(ctx).Info("anywhere edtion is empty, skip EnsureAnywhereEdtion")
		return nil
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	hubClient, err := extenderapi.GetExtenderClient(config)
	if err != nil {
		return err
	}
	current, err := extenderapi.GetManagedCluster(hubClient, c.Name)
	if err != nil {
		return err
	}

	if c.Annotations[platformv1.AnywhereLocalizationsAnno] != "" {
		localizationsJSON, err := base64.StdEncoding.DecodeString(c.Annotations[platformv1.AnywhereLocalizationsAnno])
		if err != nil {
			return fmt.Errorf("decode localizations failed: %v", err)

		}
		localizations := new(appsv1alpha1.LocalizationList)
		err = json.Unmarshal(localizationsJSON, localizations)
		if err != nil {
			return fmt.Errorf("unmarshal localization failed %v", err)
		}

		for _, l := range localizations.Items {
			l.Namespace = current.Namespace
			err := hubClient.Create(ctx, &l)
			if err != nil && !apierrors.IsAlreadyExists(err) {
				return fmt.Errorf("create localization %+v failed: %v", l, err)
			}
		}
	}

	desired := current.DeepCopy()
	desired.Labels[platformv1.AnywhereEdtionLabel] = c.Labels[platformv1.AnywhereEdtionLabel]
	err = hubClient.Patch(ctx, desired, runtimeclient.MergeFrom(current))
	if err != nil {
		return fmt.Errorf("patch managed cls failed %v", err)
	}
	return nil
}

func (p *Provider) EnsureCheckAnywhereSubscription(ctx context.Context, c *v1.Cluster) error {
	if c.Annotations[platformv1.AnywhereSubscriptionNameAnno] == "" {
		log.FromContext(ctx).Info("anywhere subscription name is empty, skip subscription")
		return nil
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	hubClient, err := extenderapi.GetExtenderClient(config)
	if err != nil {
		return err
	}
	mcls, err := extenderapi.GetManagedCluster(hubClient, c.Name)
	if err != nil {
		return err
	}
	sub, err := extenderapi.GetSubscription(hubClient, c.Annotations[platformv1.AnywhereSubscriptionNameAnno], c.Annotations[platformv1.AnywhereSubscriptionNamespaceAnno])
	if err != nil {
		return err
	}
	_ = wait.PollImmediate(15*time.Second, 10*time.Minute, func() (bool, error) {
		for i, feed := range sub.Spec.Feeds {
			var helmrelease *appsv1alpha1.HelmRelease
			helmrelease, err = extenderapi.GetHelmRelease(hubClient, extenderapi.GenerateHelmReleaseName(sub.Namespace, sub.Name, feed.Namespace, feed.Name), mcls.Namespace)
			if err != nil {
				if apierrors.IsNotFound(err) {
					return false, nil
				}
				err = fmt.Errorf("get helmrelease %s failed: %v", feed.Name, err)
				return false, err
			}
			if helmrelease != nil && helmrelease.Status.Phase != release.StatusDeployed {
				err = fmt.Errorf("%d/%d charts are deployed, %s is not deployed, phase: %s, description: %s, notes: %s",
					i,
					len(sub.Spec.Feeds),
					feed.Name,
					helmrelease.Status.Phase,
					helmrelease.Status.Description,
					helmrelease.Status.Notes,
				)
				if helmrelease.Status.Phase == release.StatusFailed {
					log.FromContext(ctx).Errorf("cluster %s install chart %s failed, phase: %s, description: %s, notes: %s", c.Name, feed.Name, helmrelease.Status.Phase, helmrelease.Status.Description, helmrelease.Status.Notes)
				}
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		return err
	}
	// Update appVersion after all system components deployed
	c.Status.AppVersion = c.Spec.AppVersion
	c.Status.ComponentPhase = platformv1.ComponentDeployed
	return nil
}

// Ensure anywhere addon applications
func (p *Provider) EnsureAnywhereAddons(ctx context.Context, c *v1.Cluster) error {
	config, err := c.RESTConfig()
	if err != nil {
		return err
	}
	extenderClient, err := extenderapi.GetExtenderClient(config)
	if err != nil {
		return err
	}
	if c.Annotations[platformv1.AnywhereApplicationAnno] != "" {
		applicationJSON, err := base64.StdEncoding.DecodeString(c.Annotations[platformv1.AnywhereApplicationAnno])
		if err != nil {
			return fmt.Errorf("decode application JSON failed: %v", err)

		}
		applications := &application.AppList{}
		err = json.Unmarshal(applicationJSON, applications)
		if err != nil {
			return fmt.Errorf("unmarshal application failed %v", err)
		}

		for _, app := range applications.Items {
			err := extenderClient.Create(ctx, &app)
			if err != nil && !apierrors.IsAlreadyExists(err) {
				return fmt.Errorf("create application %+v failed: %v", app, err)
			}
		}
	}
	return nil
}

// update cluster to connect remote cluster apiserver
func (p *Provider) EnsureClusterAddressReal(ctx context.Context, c *v1.Cluster) error {
	var hubAPIServerURL *url.URL
	var err error
	if urlValue, ok := c.Annotations[platformv1.HubAPIServerAnno]; ok {
		hubAPIServerURL, err = url.Parse(urlValue)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("cluster %s annotation %s dont exist", c.Name, platformv1.HubAPIServerAnno)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	hubClient, err := extenderapi.GetExtenderClient(config)
	if err != nil {
		return err
	}
	var currentManagerCluster *clustersv1beta1.ManagedCluster
	_ = wait.PollImmediate(5*time.Second, 3*time.Minute, func() (bool, error) {
		currentManagerCluster, err = extenderapi.GetManagedCluster(hubClient, c.Name)
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("get mcls failed: %v", err)
	}

	hubAPIServerPort, err := strconv.ParseInt(hubAPIServerURL.Port(), 10, 32)
	if err != nil {
		return err
	}
	address := platformv1.ClusterAddress{
		Type: platformv1.AddressReal,
		Host: hubAPIServerURL.Hostname(),
		Port: int32(hubAPIServerPort),
		Path: fmt.Sprintf("/apis/proxies.clusternet.io/v1alpha1/sockets/%s/proxy/direct", currentManagerCluster.Spec.ClusterID),
	}
	c.Status.Addresses = make([]platformv1.ClusterAddress, 0)
	c.Status.Addresses = append(c.Status.Addresses, address)
	return nil
}

// update cluster credential to connect remote cluster apiserver
func (p *Provider) EnsureModifyClusterCredential(ctx context.Context, c *v1.Cluster) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	hubClient, err := extenderapi.GetExtenderClient(config)
	if err != nil {
		return err
	}
	var currentManagerCluster *clustersv1beta1.ManagedCluster
	_ = wait.PollImmediate(5*time.Second, 3*time.Minute, func() (bool, error) {
		currentManagerCluster, err = extenderapi.GetManagedCluster(hubClient, c.Name)
		if err != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("get mcls failed: %v", err)
	}

	inClusterClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	mclsSecret, err := inClusterClient.CoreV1().Secrets(currentManagerCluster.Namespace).Get(ctx, "child-cluster-deployer", metav1.GetOptions{})
	if err != nil {
		return err
	}
	// token is decode data
	token, ok := mclsSecret.Data["token"]
	if !ok {
		return fmt.Errorf("mcls %s dont have token data in child-cluster-deployer secret in %s namespace", currentManagerCluster.Name, currentManagerCluster.Namespace)
	}

	cc := c.ClusterCredential
	if cc.Annotations == nil {
		cc.Annotations = make(map[string]string)
	}
	if _, ok := cc.Annotations[platformv1.CredentialTokenAnno]; !ok {
		credentialToken := *cc.Token
		cc.Annotations[platformv1.CredentialTokenAnno] = base64.StdEncoding.EncodeToString(([]byte)(credentialToken))
	}
	cc.Token = nil
	cc.Impersonate = "clusternet"
	cc.Username = "system:anonymous"
	cc.ImpersonateUserExtra = platformv1.ImpersonateUserExtra{
		"clusternet-token": string(token),
	}
	c.RegisterRestConfig(c.ClusterCredential.RESTConfig(c.Cluster))
	c.IsCredentialChanged = true
	return nil
}

func (p *Provider) EnsureKubeAPIServerRestart(ctx context.Context, c *v1.Cluster) error {
	if c.Spec.Machines == nil || len(c.Spec.Machines) == 0 {
		return fmt.Errorf("cluster %s dont have machine info", c.Name)
	}

	for _, machine := range c.Spec.Machines {
		machineSSH, err := machine.SSH()
		if err != nil {
			return err
		}

		podName := fmt.Sprintf("kube-apiserver-%s", machine.IP)
		cmd := fmt.Sprintf("kubectl delete pod %s -n kube-system", podName)
		_, err = machineSSH.CombinedOutput(cmd)
		if err != nil {
			return err
		}
	}

	clientSet, err := c.ClientsetForBootstrap()
	if err != nil {
		return err
	}

	ok := false
	_ = wait.PollImmediate(5*time.Second, 3*time.Minute, func() (bool, error) {
		ok, err = apiclient.CheckPodReadyWithLabel(ctx, clientSet, "kube-system", "component=kube-apiserver")
		if err != nil {
			return false, nil
		}
		return ok, nil
	})
	if err != nil {
		return fmt.Errorf("check kube-apiserver pod failed: %v", err)
	}
	if !ok {
		return fmt.Errorf("kube-apiserver is not ready yet")
	}
	return nil
}

func (p *Provider) EnsureRegisterGlobalCluster(ctx context.Context, c *v1.Cluster) error {
	var err error
	platformClient, err := c.PlatformClientsetForBootstrap()
	if err != nil {
		return fmt.Errorf("get platfomr client failed: %v", err)
	}

	// ensure api ready
	_ = wait.PollImmediate(5*time.Second, 3*time.Minute, func() (bool, error) {
		_, err = platformClient.Clusters().List(ctx, metav1.ListOptions{})
		if err != nil {
			err = fmt.Errorf("check cluster resources failed %v", err)
			return false, nil
		}
		_, err = platformClient.ClusterCredentials().List(ctx, metav1.ListOptions{})
		if err != nil {
			err = fmt.Errorf("check cluster credential resources failed %v", err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return err
	}

	globalCluster := c.DeepCopy()
	globalClusterCredential := c.ClusterCredential.DeepCopy()
	globalClusterName := "global"
	globalClusterCredentialName := fmt.Sprintf("cc-%s", globalClusterName)

	globalCluster.Name = globalClusterName
	globalCluster.ResourceVersion = ""
	globalCluster.UID = ""
	globalCluster.Status.Phase = platformv1.ClusterRunning
	if globalCluster.Spec.ClusterCredentialRef == nil {
		return fmt.Errorf("cluster %s dont have credential reference", globalCluster.Name)
	}
	globalCluster.Spec.ClusterCredentialRef.Name = globalClusterCredentialName
	globalCluster.Spec.Type = "Baremetal"
	globalCluster.Spec.DisplayName = "TKE"
	globalCluster.Status.Addresses = make([]platformv1.ClusterAddress, 0)
	if err = completePlatformClusterAddresses(globalCluster); err != nil {
		return fmt.Errorf("complete platfor cluster addr failed: %v", err)
	}

	for i := range globalCluster.Spec.Machines {
		globalCluster.Spec.Machines[i].Proxy = platformv1.ClusterMachineProxy{}
	}

	globalClusterCredential.Name = globalClusterCredentialName
	globalClusterCredential.ResourceVersion = ""
	globalClusterCredential.UID = ""
	globalClusterCredential.OwnerReferences = nil
	globalClusterCredential.ClusterName = globalClusterName
	globalClusterCredential.Username = ""
	globalClusterCredential.Impersonate = ""
	globalClusterCredential.ImpersonateUserExtra = nil
	delete(globalClusterCredential.Labels, platformv1.ClusterNameLable)
	if token, ok := globalClusterCredential.Annotations[platformv1.CredentialTokenAnno]; ok {
		tokenBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return fmt.Errorf("decode annotaions platformv1.CredentialTokenAnno %s failed: %v", token, err)
		}
		tokenStr := string(tokenBytes)
		globalClusterCredential.Token = &tokenStr
		delete(globalClusterCredential.Annotations, platformv1.CredentialTokenAnno)
	} else {
		return fmt.Errorf("cluster %s credential %s dont have token annotation", c.Name, c.ClusterCredential.Name)
	}

	globalCluster.SetCondition(platformv1.ClusterCondition{
		Type:    "EnsureGlobalClusterRegistration",
		Status:  platformv1.ConditionTrue,
		Message: "",
		Reason:  "",
	}, false)

	_, err = platformClient.ClusterCredentials().Get(ctx, globalClusterCredential.Name, metav1.GetOptions{})
	if err == nil {
		err := platformClient.ClusterCredentials().Delete(ctx, globalClusterCredential.Name, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("clean cluster credential failed: %v", err)
		}
	}
	_, err = platformClient.ClusterCredentials().Create(ctx, globalClusterCredential, metav1.CreateOptions{})
	if err != nil {
		if err != nil {
			return fmt.Errorf("create cluster credential failed: %v", err)
		}
		return err
	}

	_, err = platformClient.Clusters().Get(ctx, globalCluster.Name, metav1.GetOptions{})
	if err == nil {
		err := platformClient.Clusters().Delete(ctx, globalCluster.Name, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("clean cluster failed: %v", err)
		}
	}
	_, err = platformClient.Clusters().Create(ctx, globalCluster, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create cluster failed: %v", err)
	}

	return nil
}
