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
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/gpu"

	"tkestack.io/tke/pkg/platform/provider/baremetal/images"
	galaxyimages "tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy/images"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/addons/cniplugins"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/docker"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/galaxy"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeadm"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubeconfig"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/kubelet"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/markcontrolplane"
	"tkestack.io/tke/pkg/platform/provider/baremetal/preflight"
	"tkestack.io/tke/pkg/platform/provider/baremetal/util/hosts"
	"tkestack.io/tke/pkg/util/log"
)

const (
	sysctlFile       = "/etc/sysctl.conf"
	sysctlCustomFile = "/etc/sysctl.d/99-tke.conf"
	moduleFile       = "/etc/modules-load.d/tke.conf"
)

func (p *Provider) EnsurePreflight(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		s := c.SSH[machine.IP]

		err := preflight.RunMasterChecks(s)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureRegistryHosts(c *Cluster) error {
	if !c.Registry.UseTKE() {
		return nil
	}

	domains := []string{
		c.Registry.Domain,
		c.Spec.TenantID + "." + c.Registry.Domain,
	}
	for _, machine := range c.Spec.Machines {
		s := c.SSH[machine.IP]

		for _, one := range domains {
			remoteHosts := hosts.RemoteHosts{Host: one, SSH: s}
			err := remoteHosts.Set(c.Registry.IP)
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
		}
	}

	return nil
}

func (p *Provider) EnsureKernelModule(c *Cluster) error {
	modules := []string{"iptable_nat"}
	var data bytes.Buffer
	for _, machine := range c.Spec.Machines {
		s := c.SSH[machine.IP]

		for _, m := range modules {
			_, err := s.CombinedOutput(fmt.Sprintf("modprobe %s", m))
			if err != nil {
				return errors.Wrap(err, machine.IP)
			}
			data.WriteString(m + "\n")
		}
		err := s.WriteFile(strings.NewReader(data.String()), moduleFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func setFileContent(file, pattern, content string) string {
	return fmt.Sprintf("grep -Pq '%s' %s && sed -i 's;%s;%s;g' %s|| echo '%s' >> %s",
		pattern, file,
		pattern, content, file,
		content, file)
}

func (p *Provider) EnsureSysctl(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		s := c.SSH[machine.IP]

		_, err := s.CombinedOutput(setFileContent(sysctlFile, "net.ipv4.ip_forward.*", "net.ipv4.ip_forward = 1"))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		_, err = s.CombinedOutput(setFileContent(sysctlFile, "net.bridge.bridge-nf-call-iptables.*", "net.bridge.bridge-nf-call-iptables = 1"))
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		f, err := os.Open(path.Join(constants.ConfDir, "sysctl.conf"))
		if err == nil {
			err = s.WriteFile(f, sysctlCustomFile)
			if err != nil {
				return err
			}
		}

		_, err = s.CombinedOutput("sysctl --system")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureDisableSwap(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		s := c.SSH[machine.IP]

		_, err := s.CombinedOutput("swapoff -a && sed 's/^[^#]*swap/#&/' /etc/fstab")
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

// 因为validate那里没法更新对象（不能存储）
// PreCrete，在api中错误只能panic，响应不会有报错提示，所以只能挪到这里处理
func (p *Provider) EnsureClusterComplete(cluster *Cluster) error {
	serviceCIDR, nodeCIDRMaskSize, err := GetServiceCIDRAndNodeCIDRMaskSize(cluster.Spec.ClusterCIDR, *cluster.Spec.Properties.MaxClusterServiceNum, *cluster.Spec.Properties.MaxNodePodNum)
	if err != nil {
		return errors.Wrap(err, "GetServiceCIDRAndNodeCIDRMaskSize error")
	}
	cluster.Status.ServiceCIDR = serviceCIDR
	cluster.Status.NodeCIDRMaskSize = nodeCIDRMaskSize

	dnsIP, err := GetDNSIP(cluster.Status.ServiceCIDR)
	if err != nil {
		return errors.Wrap(err, "get DNS IP error")
	}
	cluster.Status.DNSIP = dnsIP.String()
	for _, m := range cluster.Spec.Machines {
		cluster.AddAddress(platformv1.AddressReal, m.IP, 6443)
	}

	token := ksuid.New().String()
	cluster.ClusterCredential.Token = &token

	bootstrapToken, err := bootstraputil.GenerateBootstrapToken()
	if err != nil {
		return err
	}
	cluster.ClusterCredential.BootstrapToken = &bootstrapToken

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	certificateKey := hex.EncodeToString(bytes)
	cluster.ClusterCredential.CertificateKey = &certificateKey

	return nil
}

func (p *Provider) EnsureKubeconfig(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		option := &kubeconfig.Option{
			MasterEndpoint: "https://127.0.0.1:6443",
			ClusterName:    c.Name,
			CACert:         c.ClusterCredential.CACert,
			Token:          *c.ClusterCredential.Token,
		}
		err := kubeconfig.Install(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureNvidiaDriver(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		if !gpu.IsEnable(machine.Labels) {
			continue
		}
		err := gpu.InstallNvidiaDriver(c.SSH[machine.IP], &gpu.NvidiaDriverOption{Version: p.config.NvidiaDriver.DefaultVersion})
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureNvidiaContainerRuntime(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		if !gpu.IsEnable(machine.Labels) {
			continue
		}
		err := gpu.InstallNvidiaContainerRuntime(c.SSH[machine.IP], &gpu.NvidiaContainerRuntimeOption{Version: p.config.NvidiaContainerRuntime.DefaultVersion})
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureDocker(c *Cluster) error {
	insecureRegistries := fmt.Sprintf(`"%s"`, c.Registry.Domain)
	if c.Config.Registry.UseTKE() {
		insecureRegistries = fmt.Sprintf(`%s,"%s"`, insecureRegistries, c.Spec.TenantID+"."+c.Registry.Domain)
	}
	option := &docker.Option{
		Version:            c.Docker.DefaultVersion,
		InsecureRegistries: insecureRegistries,
		RegistryDomain:     c.Registry.Domain,
		ExtraArgs:          c.Spec.DockerExtraArgs,
	}
	for _, machine := range c.Spec.Machines {
		option.IsGPU = gpu.IsEnable(machine.Labels)
		err := docker.Install(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureKubeadm(c *Cluster) error {
	for _, machine := range c.Spec.Machines {
		err := kubeadm.Install(c.SSH[machine.IP])
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsurePrepareForControlplane(c *Cluster) error {
	oidcCa, _ := ioutil.ReadFile(path.Join(constants.ConfDir, constants.OIDCCACertName))

	for _, machine := range c.Spec.Machines {
		tokenData := fmt.Sprintf(tokenFileTemplate, *c.ClusterCredential.Token)
		err := c.SSH[machine.IP].WriteFile(strings.NewReader(tokenData), constants.TokenFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		err = c.SSH[machine.IP].WriteFile(strings.NewReader(schedulerPolicyConfig), constants.SchedulerPolicyConfigFile)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}

		if len(oidcCa) != 0 {
			err = c.SSH[machine.IP].WriteFile(bytes.NewReader(oidcCa), constants.OIDCCACertFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getKubeadmInitOption(c *Cluster) *kubeadm.InitOption {
	controlPlaneEndpoint := fmt.Sprintf("%s:6443", c.Spec.Machines[0].IP)
	addr := c.Address(platformv1.AddressAdvertise)
	if addr != nil {
		controlPlaneEndpoint = fmt.Sprintf("%s:%d", addr.Host, addr.Port)
	}
	return &kubeadm.InitOption{
		KubeadmConfigFileName: constants.KubeadmConfigFileName,
		NodeName:              c.Spec.Machines[0].IP,
		BootstrapToken:        *c.ClusterCredential.BootstrapToken,
		CertificateKey:        *c.ClusterCredential.CertificateKey,

		ETCDImageTag:         images.Get().ETCD.Tag,
		CoreDNSImageTag:      images.Get().CoreDNS.Tag,
		KubernetesVersion:    c.Spec.Version,
		ControlPlaneEndpoint: controlPlaneEndpoint,

		DNSDomain:             c.Spec.DNSDomain,
		ServiceSubnet:         c.Status.ServiceCIDR,
		NodeCIDRMaskSize:      c.Status.NodeCIDRMaskSize,
		ClusterCIDR:           c.Spec.ClusterCIDR,
		ServiceClusterIPRange: c.Status.ServiceCIDR,

		APIServerExtraArgs:         c.Spec.APIServerExtraArgs,
		ControllerManagerExtraArgs: c.Spec.ControllerManagerExtraArgs,
		SchedulerExtraArgs:         c.Spec.SchedulerExtraArgs,

		ImageRepository: c.Registry.Prefix,
		ClusterName:     c.Name,
	}
}

func (p *Provider) EnsureKubeadmInitKubeletStartPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c),
		fmt.Sprintf("kubelet-start --node-name=%s", c.Spec.Machines[0].IP))
}

func (p *Provider) EnsureKubeadmInitCertsPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "certs all")
}

func (p *Provider) EnsureKubeadmInitKubeConfigPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "kubeconfig all")
}

func (p *Provider) EnsureKubeadmInitControlPlanePhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "control-plane all")
}

func (p *Provider) EnsureKubeadmInitEtcdPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "etcd local")
}

func (p *Provider) EnsureKubeadmInitUploadConfigPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "upload-config all ")
}

func (p *Provider) EnsureKubeadmInitUploadCertsPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "upload-certs --upload-certs")
}

func (p *Provider) EnsureKubeadmInitBootstrapTokenPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "bootstrap-token")
}

func (p *Provider) EnsureKubeadmInitAddonPhase(c *Cluster) error {
	return kubeadm.Init(c.SSH[c.Spec.Machines[0].IP], getKubeadmInitOption(c), "addon all")
}

func (p *Provider) EnsureGalaxy(c *Cluster) error {
	clientset, err := c.Clientset()
	if err != nil {
		return err
	}
	return galaxy.Install(clientset, &galaxy.Option{
		Version:   galaxyimages.LatestVersion,
		NodeCIDR:  c.Cluster.Spec.ClusterCIDR,
		NetDevice: c.Cluster.Spec.NetworkDevice,
	})
}

func (p *Provider) EnsureJoinControlePlane(c *Cluster) error {
	oidcCa, _ := ioutil.ReadFile(path.Join(constants.ConfDir, constants.OIDCCACertName))
	option := &kubeadm.JoinControlePlaneOption{
		BootstrapToken:       *c.ClusterCredential.BootstrapToken,
		CertificateKey:       *c.ClusterCredential.CertificateKey,
		ControlPlaneEndpoint: fmt.Sprintf("%s:6443", c.Spec.Machines[0].IP),
		OIDCCA:               oidcCa,
	}
	for _, machine := range c.Spec.Machines[1:] {
		option.NodeName = machine.IP
		err := kubeadm.JoinControlePlane(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureStoreCredential(c *Cluster) error {
	data, err := c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.CACertName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.CACertName)
	}
	c.ClusterCredential.CACert = data

	data, err = c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.EtcdCACertName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.EtcdCACertName)
	}
	c.ClusterCredential.ETCDCACert = data

	data, err = c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.APIServerEtcdClientCertName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.APIServerEtcdClientCertName)
	}
	c.ClusterCredential.ETCDAPIClientCert = data

	data, err = c.SSH[c.Spec.Machines[0].IP].ReadFile(constants.APIServerEtcdClientKeyName)
	if err != nil {
		return errors.Wrapf(err, "read %s error", constants.APIServerEtcdClientKeyName)
	}
	c.ClusterCredential.ETCDAPIClientKey = data

	return nil
}

func (p *Provider) EnsurePatchAnnotation(c *Cluster) error {
	fileData := map[string]string{
		constants.EtcdPodManifestFile:                  `  annotations:\n    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "https"\n    prometheus.io/port: "2379"`,
		constants.KubeAPIServerPodManifestFile:         `  annotations:\n    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "https"\n    prometheus.io/port: "6443"`,
		constants.KubeControllerManagerPodManifestFile: `  annotations:\n    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "http"\n    prometheus.io/port: "10252"`,
		constants.KubeSchedulerPodManifestFile:         `  annotations:\n    scheduler.alpha.kubernetes.io/critical-pod: ""\n    tke.prometheus.io/scrape: "true"\n    prometheus.io/scheme: "http"\n    prometheus.io/port: "10251"`,
	}
	for _, machine := range c.Spec.Machines {
		for file, data := range fileData {
			cmd := fmt.Sprintf(`grep 'prometheus.io/port' %s || sed -i '3a\%s' %s`, file, data, file)
			_, stderr, exit, err := c.SSH[machine.IP].Exec(cmd)
			if err != nil || exit != 0 {
				return fmt.Errorf("exec %q failed:exit %d:stderr %s:error %s", cmd, exit, stderr, err)
			}
		}
	}

	return nil
}

func (p *Provider) EnsureKubelet(c *Cluster) error {
	option := &kubelet.Option{
		Version:   c.Spec.Version,
		ExtraArgs: c.Spec.KubeletExtraArgs,
	}
	for _, machine := range c.Spec.Machines {
		err := kubelet.Install(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureCNIPlugins(c *Cluster) error {
	option := &cniplugins.Option{
		Version: c.CNIPlugins.DefaultVersion,
	}
	for _, machine := range c.Spec.Machines {
		err := cniplugins.Install(c.SSH[machine.IP], option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureKubeadmInitWaitControlPlanePhase(c *Cluster) error {
	start := time.Now()

	return wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
		healthStatus := 0
		clientset, err := c.Clientset()
		if err != nil {
			log.Warn(err.Error())
			return false, nil
		}
		clientset.Discovery().RESTClient().Get().AbsPath("/healthz").Do().StatusCode(&healthStatus)
		if healthStatus != http.StatusOK {
			return false, nil
		}

		log.Infof("All control plane components are healthy after %f seconds\n", time.Since(start).Seconds())
		return true, nil
	})
}

func (p *Provider) EnsureMarkControlPlane(c *Cluster) error {
	clientset, err := c.Clientset()
	if err != nil {
		return err
	}

	option := &markcontrolplane.Option{}
	if c.Name == "global" {
		log.Info("don't set taints because this special cluster is resvered for TKE bootstrap.")
	} else {
		option.Taints = []corev1.Taint{
			{
				Key:    "node-role.kubernetes.io/master",
				Effect: corev1.TaintEffectNoSchedule,
			},
		}
	}

	for _, machine := range c.Spec.Machines {
		option.NodeName = machine.IP
		err := markcontrolplane.Install(clientset, option)
		if err != nil {
			return errors.Wrap(err, machine.IP)
		}
	}

	return nil
}

func (p *Provider) EnsureNvidiaDevicePlugin(c *Cluster) error {
	if c.Cluster.Spec.Features.GPUType == nil {
		return nil
	}

	client, err := c.Clientset()
	if err != nil {
		return err
	}
	option := &gpu.NvidiaDevicePluginOption{
		Image: images.Get().NvidiaDevicePlugin.FullName(),
	}
	err = gpu.InstallNvidiaDevicePlugin(client, option)
	if err != nil {
		return err
	}

	return nil
}
