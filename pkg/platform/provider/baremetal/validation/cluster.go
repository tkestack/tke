/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

package validation

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	k8serror "k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	netutils "k8s.io/utils/net"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	csioperatorimage "tkestack.io/tke/pkg/platform/provider/baremetal/phases/csioperator/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/gpu"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/platform/util"
	vendor "tkestack.io/tke/pkg/platform/util/kubevendor"
	"tkestack.io/tke/pkg/spec"
	pkgutil "tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/ipallocator"
	"tkestack.io/tke/pkg/util/log"
	utilmath "tkestack.io/tke/pkg/util/math"
	"tkestack.io/tke/pkg/util/ssh"
	"tkestack.io/tke/pkg/util/validation"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

var (
	nodePodNumAvails        = []int32{16, 32, 64, 128, 256}
	clusterServiceNumAvails = []int32{32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768}
	supportedOSList         = []string{}
	reservePorts            = []int{
		// kube-apiserver
		6443,
		// etcd
		2379, 2380,
		// kubelet
		10250, 10251, 10252, 10253, 10254, 10255,
		// ingress
		80, 443, 8181,
		// node exporter
		9100,
		// tke-gateway
		31180, 31443,
		// tke-auth
		31138,
	}
)

// ValidateCluster validates a given Cluster.
func ValidateCluster(platformClient platformv1client.PlatformV1Interface, obj *types.Cluster) field.ErrorList {
	allErrs := ValidatClusterSpec(platformClient, obj.Name, obj.Cluster, field.NewPath("spec"), obj.Status.Phase, true)
	return allErrs
}

// ValidateCluster validates a given Cluster.
func ValidateClusterUpdate(platformClient platformv1client.PlatformV1Interface, cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList {
	fldPath := field.NewPath("spec")
	allErrs := ValidatClusterSpec(platformClient, cluster.Name, cluster.Cluster, fldPath, cluster.Status.Phase, false)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.NetworkDevice, oldCluster.Spec.NetworkDevice, fldPath.Child("networkDevice"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.ClusterCIDR, oldCluster.Spec.ClusterCIDR, fldPath.Child("clusterCIDR"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.DNSDomain, oldCluster.Spec.DNSDomain, fldPath.Child("dnsDomain"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.DockerExtraArgs, oldCluster.Spec.DockerExtraArgs, fldPath.Child("dockerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.KubeletExtraArgs, oldCluster.Spec.KubeletExtraArgs, fldPath.Child("kubeletExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.APIServerExtraArgs, oldCluster.Spec.APIServerExtraArgs, fldPath.Child("apiServerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.ControllerManagerExtraArgs, oldCluster.Spec.ControllerManagerExtraArgs, fldPath.Child("controllerManagerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.SchedulerExtraArgs, oldCluster.Spec.SchedulerExtraArgs, fldPath.Child("schedulerExtraArgs"))...)
	allErrs = append(allErrs, ValidateClusterScale(cluster.Cluster, oldCluster.Cluster, fldPath.Child("machines"))...)

	return allErrs
}

// ValidateClusterScale tests if master scale up/down to a cluster is valid.
func ValidateClusterScale(cluster *platform.Cluster, oldCluster *platform.Cluster, fldPath *field.Path) field.ErrorList {

	allErrs := field.ErrorList{}
	if len(cluster.Spec.Machines) == len(oldCluster.Spec.Machines) {
		return allErrs
	}
	ha := cluster.Spec.Features.HA
	if ha == nil {
		allErrs = append(allErrs, field.Invalid(fldPath, cluster.Spec.Machines, "HA configuration should enabled for master scale"))
		return allErrs
	}
	if ha.TKEHA == nil && ha.ThirdPartyHA == nil {
		allErrs = append(allErrs, field.Invalid(fldPath, cluster.Spec.Machines, "tkestack HA or third party HA should enabled for master scale"))
		return allErrs
	}
	_, err := util.PrepareClusterScale(cluster, oldCluster)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, cluster.Spec.Machines, err.Error()))
	}
	return allErrs
}

// ValidatClusterSpec validates a given ClusterSpec.
func ValidatClusterSpec(platformClient platformv1client.PlatformV1Interface, clusterName string, cls *platform.Cluster, fldPath *field.Path, phase platform.ClusterPhase, validateMachine bool) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateClusterSpecVersion(platformClient, clusterName, cls.Spec.Version, fldPath.Child("version"), phase)...)
	allErrs = append(allErrs, ValidateCIDRs(cls, fldPath)...)
	allErrs = append(allErrs, ValidateClusterProperty(&cls.Spec, fldPath.Child("properties"))...)
	if validateMachine {
		allErrs = append(allErrs, ValidateClusterMachines(cls, fldPath.Child("machines"))...)
	}
	allErrs = append(allErrs, ValidateClusterGPUMachines(cls.Spec.Machines, fldPath.Child("machines"))...)
	allErrs = append(allErrs, ValidateClusterFeature(&cls.Spec, fldPath.Child("features"))...)

	return allErrs
}

// ValidateClusterSpecVersion validates a given version.
func ValidateClusterSpecVersion(platformClient platformv1client.PlatformV1Interface, clsName, version string, fldPath *field.Path, phase platform.ClusterPhase) field.ErrorList {
	allErrs := field.ErrorList{}

	k8sValidVersions, err := getK8sValidVersions(platformClient, clsName)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldPath, err))
		return allErrs
	}

	if phase == platform.ClusterInitializing {
		allErrs = utilvalidation.ValidateEnum(version, fldPath, k8sValidVersions)
	}
	if phase == platform.ClusterUpgrading {
		c, err := platformClient.Clusters().Get(context.Background(), clsName, metav1.GetOptions{})
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
			return allErrs
		}
		dstKubevendor := vendor.GetKubeVendor(version)
		if err := validateKubevendor(c.Status.KubeVendor, dstKubevendor); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath,
				err,
				"current kubevendor is not supported to upgrade to input version"))
		}
	}

	return allErrs
}

// ValidateClusterMachines validates a given CluterMachines.
func ValidateClusterMachines(cls *platform.Cluster, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	proxyErrs := field.ErrorList{}
	sshErrs := field.ErrorList{}
	timeErrs := field.ErrorList{}
	osErrs := field.ErrorList{}
	mcReErrs := field.ErrorList{}
	routeErrs := field.ErrorList{}
	portsErrs := field.ErrorList{}

	proxyResult := TKEValidateResult{}
	sshResult := TKEValidateResult{}
	timeResult := TKEValidateResult{}
	osResult := TKEValidateResult{}
	mcReResult := TKEValidateResult{}
	routeResult := TKEValidateResult{}
	portsResult := TKEValidateResult{}

	var masters []*ssh.SSH
	for i, one := range cls.Spec.Machines {
		var proxy ssh.Proxy
		switch one.Proxy.Type {
		case platform.SSHJumpServer:
			sshproxy := ssh.JumpServer{}
			sshproxy.Host = one.Proxy.IP
			sshproxy.Port = int(one.Proxy.Port)
			sshproxy.User = one.Proxy.Username
			sshproxy.Password = string(one.Proxy.Password)
			sshproxy.PrivateKey = one.Proxy.PrivateKey
			sshproxy.PassPhrase = one.Proxy.PassPhrase
			sshproxy.DialTimeOut = time.Second
			sshproxy.Retry = 0
			proxy = sshproxy
		}
		proxyErrs = append(proxyErrs, ValidateProxy(fldPath.Index(i), proxy)...)
		proxyResult.Checked = true
		// if proxy has err, no need to check ssh
		if len(proxyErrs) == 0 {
			sshErrs = append(sshErrs, ValidateSSH(fldPath.Index(i), one.IP, int(one.Port), one.Username, one.Password, one.PrivateKey, one.PassPhrase, proxy)...)
			// when get ssh err or last machine ssh is checked, ssh can be considered checked
			if len(sshErrs) != 0 || i == len(cls.Spec.Machines)-1 {
				sshResult.Checked = true
			}
		}
		if len(sshErrs) == 0 && len(proxyErrs) == 0 {
			master, _ := one.SSH()
			masters = append(masters, master)
		}
	}

	if len(masters) == len(cls.Spec.Machines) {
		timeErrs = ValidateMasterTimeOffset(fldPath, masters)
		timeResult.Checked = true

		if len(supportedOSList) != 0 {
			osErrs = ValidateOSVersion(fldPath, masters)
			osResult.Checked = true
		} else {
			log.Warn("skip validate OS since supported OS list is empty")
		}

		mcReErrs = ValidateMachineResource(fldPath, masters)
		mcReResult.Checked = true

		routeErrs = ValidateDefaultRoute(fldPath, masters, cls.Spec.NetworkDevice)
		routeResult.Checked = true

		portsErrs = ValidateReservePorts(fldPath, masters)
		portsResult.Checked = true
	}
	if _, ok := cls.Annotations[platform.AnywhereValidateAnno]; ok {
		proxyResult.Name = AnywhereValidateItemTunnelConnectivity
		proxyResult.Description = "Verify Proxy Tunnel Connectivity"
		proxyResult.ErrorList = proxyErrs

		sshResult.Name = AnywhereValidateItemSSH
		sshResult.Description = "Verify SSH is Available"
		sshResult.ErrorList = sshErrs

		timeResult.Name = AnywhereValidateItemTimeDiff
		timeResult.Description = fmt.Sprintf("Verify Clock Gap between Master nodes is not More than %d Second(s)", MaxTimeOffset)
		timeResult.ErrorList = timeErrs

		osResult.Name = AnywhereValidateItemOSVersion
		osResult.Description = "Verify Target Machine OS"
		osResult.ErrorList = osErrs

		mcReResult.Name = AnywhereValidateItemMachineResource
		mcReResult.Description = "Verify Target Machine Resource"
		mcReResult.ErrorList = mcReErrs

		routeResult.Name = AnywhereValidateItemDefaultRoute
		routeResult.Description = "Verify Default Route Network Interface"
		routeResult.ErrorList = routeErrs

		portsResult.Name = AnywhereValidateItemReservePorts
		portsResult.Description = "Verify ReservePorts Status"
		portsResult.ErrorList = portsErrs

		allErrs = append(allErrs,
			proxyResult.ToFieldError(),
			sshResult.ToFieldError(),
			timeResult.ToFieldError(),
			osResult.ToFieldError(),
			mcReResult.ToFieldError(),
			routeResult.ToFieldError(),
			portsResult.ToFieldError())
	} else {
		allErrs = append(allErrs, proxyErrs...)
		allErrs = append(allErrs, sshErrs...)
		allErrs = append(allErrs, timeErrs...)
		allErrs = append(allErrs, osErrs...)
		allErrs = append(allErrs, mcReErrs...)
		allErrs = append(allErrs, routeErrs...)
		allErrs = append(allErrs, portsErrs...)
	}

	return allErrs
}

func ValidateProxy(fldPath *field.Path, proxy ssh.Proxy) field.ErrorList {
	allErrs := field.ErrorList{}
	sshConfig := &ssh.Config{
		User:        "validate",
		Host:        "127.0.0.1",
		Port:        22,
		Password:    base64.StdEncoding.EncodeToString([]byte("validate")),
		DialTimeOut: time.Second,
		Retry:       0,
		Proxy:       proxy,
	}
	s, err := ssh.New(sshConfig)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("proxy"), "", err.Error()))
		return allErrs
	}
	err = s.CheckProxyTunnel()
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("proxy"), "", err.Error()))
	}
	return allErrs
}

func ValidateOSVersion(fldPath *field.Path, sshs []*ssh.SSH) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, one := range sshs {
		os, err := ssh.OSVersion(one)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host, err.Error()))
			continue
		}
		if !pkgutil.InStringSlice(supportedOSList, os) {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host,
				fmt.Sprintf("target os %s is not in expected os list %v", os, supportedOSList)))
		}
	}
	return allErrs
}

func ValidateReservePorts(fldPath *field.Path, sshs []*ssh.SSH) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, one := range sshs {
		err := ssh.ReservePorts(one, reservePorts)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host, err.Error()))
		}
	}
	return allErrs
}

func ValidateDefaultRoute(fldPath *field.Path, sshs []*ssh.SSH, expectedNetInterface string) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, one := range sshs {
		realNetInterface := ssh.GetNetworkInterface(one, one.Host)
		if realNetInterface != expectedNetInterface {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host,
				fmt.Sprintf("network interface of IP %s is %s not %s", one.Host, realNetInterface, expectedNetInterface)))
		}
		realNetInterface = ssh.GetDefaultRouteInterface(one)
		if realNetInterface != expectedNetInterface {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host,
				fmt.Sprintf("network interface of default route is %s not %s", realNetInterface, expectedNetInterface)))
		}
	}
	return allErrs
}

func ValidateMachineResource(fldPath *field.Path, sshs []*ssh.SSH) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, one := range sshs {
		size, err := ssh.DiskAvail(one, MachineResourceRequstDiskPath)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host, err.Error()))
			continue
		}
		if size < MachineResourceRequstDiskSpace {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host,
				fmt.Sprintf("%s disk space %d GiB is smaller than request size %d GiB", MachineResourceRequstDiskPath, size, MachineResourceRequstDiskSpace)))
		}
		size, err = ssh.DiskAvail(one, MachineResourceRequstLogDiskPath)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host, err.Error()))
			continue
		}
		if size < MachineResourceRequstLogDiskSpace {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host,
				fmt.Sprintf("%s disk space %d GiB is smaller than request size %d GiB", MachineResourceRequstLogDiskPath, size, MachineResourceRequstLogDiskSpace)))
		}
		cpuNum, err := ssh.NumCPU(one)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host, err.Error()))
			continue
		}
		if cpuNum < MachineResourceRequstCPU {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host,
				fmt.Sprintf("cpu number %d is smaller than request %d", cpuNum, MachineResourceRequstCPU)))
		}
		memInBytes, err := ssh.MemoryCapacity(one)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host, err.Error()))
			continue
		}
		memInGiB := math.Ceil(float64(memInBytes) / 1024 / 1024 / 1024)
		if memInGiB < MachineResourceRequstMemory {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), one.Host,
				fmt.Sprintf("memory %d GiB is smaller than request %d GiB", int(memInGiB), MachineResourceRequstMemory)))
		}
	}

	return allErrs
}

func RegisterSupportedOSList(list []string) {
	supportedOSList = list
}

// ValidateMasterTimeOffset validates a given master time offset.
func ValidateMasterTimeOffset(fldPath *field.Path, masters []*ssh.SSH) field.ErrorList {
	allErrs := field.ErrorList{}

	times := make([]float64, 0, len(masters))
	for _, one := range masters {
		t, err := ssh.Timestamp(one)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
			return allErrs
		}
		times = append(times, float64(t))
	}
	maxIndex, maxTime := utilmath.Max(times)
	minIndex, minTime := utilmath.Min(times)
	offset := int(*maxTime) - int(*minTime)
	if offset > MaxTimeOffset {
		allErrs = append(allErrs, field.Invalid(fldPath, "",
			fmt.Sprintf("the time offset(%v-%v=%v) between node(%v) with node(%v) exceeds %d seconds, please unify machine time between nodes by using ntp or manual", int(*maxTime), int(*minTime), offset, masters[*maxIndex].Host, masters[*minIndex].Host, MaxTimeOffset)))
	}

	return allErrs
}

func getK8sValidVersions(platformClient platformv1client.PlatformV1Interface, clsName string) (validVersions []string, err error) {
	if clsName == "global" || platformClient == nil {
		return spec.K8sVersions, nil
	}

	cluster, err := platformClient.Clusters().Get(context.Background(), "global", metav1.GetOptions{})
	if err != nil {
		if k8serror.IsNotFound(err) {
			log.Warnf("global cluster is not exist")

			return spec.K8sVersions, nil
		}
		return nil, err
	}

	client, err := util.BuildExternalClientSet(context.Background(), cluster, platformClient)
	if err != nil {
		return nil, err
	}

	_, k8sValidVersions, err := util.GetPlatformVersionsFromClusterInfo(context.Background(), client)

	return k8sValidVersions, err
}

func validateKubevendor(srcKubevendor, dstKubevendor platformv1.KubeVendorType) (err error) {
	notSupportUpgradeMessage := "not support upgrade from vendor %v to vendor %v"
	switch srcKubevendor {
	case platformv1.KubeVendorTKE:
		if dstKubevendor != platformv1.KubeVendorTKE {
			return fmt.Errorf(notSupportUpgradeMessage, srcKubevendor, dstKubevendor)
		}
	case platformv1.KubeVendorOther:
		if dstKubevendor != platformv1.KubeVendorOther && dstKubevendor != platformv1.KubeVendorTKE {
			return fmt.Errorf(notSupportUpgradeMessage, srcKubevendor, dstKubevendor)
		}
	default:
		return fmt.Errorf(notSupportUpgradeMessage, srcKubevendor, dstKubevendor)
	}
	return nil
}

// ValidateCIDRs validates clusterCIDR and serviceCIDR.
func ValidateCIDRs(cls *platform.Cluster, specPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	result := TKEValidateResult{}
	var clusterCIDR, serviceCIDR *net.IPNet

	checkFunc := func(path *field.Path, cidr string) {
		cidrs := strings.Split(cidr, ",")
		dualStackEnabled := cls.Spec.Features.IPv6DualStack
		switch {
		// if DualStack only valid one cidr or two cidrs with one of each IP family
		case dualStackEnabled && len(cidrs) > 2:
			allErrs = append(allErrs, field.Invalid(path, cidr, "only one CIDR allowed or a valid DualStack CIDR (e.g. 10.100.0.0/16,fde4:8dba:82e1::/48)"))
		// if DualStack and two cidrs validate if there is at least one of each IP family
		case dualStackEnabled && len(cidrs) == 2:
			isDual, err := netutils.IsDualStackCIDRStrings(cidrs)
			if err != nil || !isDual {
				allErrs = append(allErrs, field.Invalid(path, cidr, "must be a valid DualStack CIDR (e.g. 10.100.0.0/16,fde4:8dba:82e1::/48)"))
			}
		// if not DualStack only one CIDR allowed
		case !dualStackEnabled && len(cidrs) > 1:
			allErrs = append(allErrs, field.Invalid(path, cidr, "only one CIDR allowed (e.g. 10.100.0.0/16 or fde4:8dba:82e1::/48)"))
		// if we are here means that len(cidrs) == 1, we need to validate it
		default:
			_, cidrX, err := net.ParseCIDR(cidr)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(path, cidr, "must be a valid CIDR block (e.g. 10.100.0.0/16 or fde4:8dba:82e1::/48)"))
			}
			if path == specPath.Child("clusterCIDR") {
				clusterCIDR = cidrX
				for i, mc := range cls.Spec.Machines {
					if clusterCIDR.Contains(net.ParseIP(mc.IP)) {
						allErrs = append(allErrs, field.Invalid(path.Index(i), cidr,
							fmt.Sprintf("cannot use CIDR %s, since this CIDR contains node IP %s", cidr, mc.IP)))
					}

				}
			} else {
				serviceCIDR = cidrX
			}
		}
	}

	fldPath := specPath.Child("clusterCIDR")
	cidr := cls.Spec.ClusterCIDR
	if len(cidr) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath, cidr, "ClusterCIDR is empty string"))
	} else {
		checkFunc(fldPath, cidr)
	}

	fldPath = specPath.Child("serviceCIDR")
	if cls.Spec.ServiceCIDR != nil {
		cidr = *cls.Spec.ServiceCIDR
		if len(cidr) == 0 {
			allErrs = append(allErrs, field.Invalid(fldPath, cidr, "ServiceCIDR is empty string"))
		} else {
			checkFunc(fldPath, cidr)
			if clusterCIDR != nil && serviceCIDR != nil {
				if err := validation.IsSubNetOverlapped(clusterCIDR, serviceCIDR); err != nil {
					allErrs = append(allErrs, field.Invalid(fldPath, cidr, err.Error()))
				}
				if _, err := ipallocator.GetIndexedIP(serviceCIDR, 10); err != nil {
					allErrs = append(allErrs, field.Invalid(fldPath, cidr,
						"must contains at least 10 ips, because kubeadm need the 10th ip"))
				}
			}
		}
	}

	if _, ok := cls.Annotations[platform.AnywhereValidateAnno]; ok {
		result.Name = AnywhereValidateItemHostNetOverlapping
		result.Description = "Verify Node IP(s) and CIDR Config"
		result.ErrorList = allErrs
		result.Checked = true

		return field.ErrorList{result.ToFieldError()}

	}
	return allErrs
}

// ValidateClusterProperty validates a given ClusterProperty.
func ValidateClusterProperty(spec *platform.ClusterSpec, propPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	properties := spec.Properties

	fldPath := propPath.Child("maxNodePodNum")
	if properties.MaxNodePodNum == nil {
		allErrs = append(allErrs, field.Required(fldPath, fmt.Sprintf("validate values are %v", nodePodNumAvails)))
	} else {
		allErrs = utilvalidation.ValidateEnum(*properties.MaxNodePodNum, fldPath, nodePodNumAvails)
	}

	fldPath = propPath.Child("maxClusterServiceNum")
	if properties.MaxClusterServiceNum == nil {
		if spec.ServiceCIDR == nil { // not set serviceCIDR, need set maxClusterServiceNum
			allErrs = append(allErrs, field.Required(fldPath, fmt.Sprintf("validate values are %v", clusterServiceNumAvails)))
		}
	} else {
		if spec.ServiceCIDR != nil { // spec.serviceCIDR and properties.maxClusterServiceNum can't be used together
			allErrs = append(allErrs, field.Forbidden(fldPath, "can't be used together with spec.serviceCIDR"))
		} else {
			allErrs = utilvalidation.ValidateEnum(*properties.MaxClusterServiceNum, fldPath, clusterServiceNumAvails)
			if *properties.MaxClusterServiceNum < 10 {
				allErrs = append(allErrs, field.Invalid(fldPath, *properties.MaxClusterServiceNum,
					"must be greater than or equal to 10 because kubeadm need the 10th ip"))
			}
		}
	}

	return allErrs
}

// ValidateClusterGPUMachines validates a given GPUMachines.
func ValidateClusterGPUMachines(machines []platform.ClusterMachine, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if machines == nil {
		allErrs = append(allErrs, field.Required(fldPath, ""))
	} else {
		for i, machine := range machines {
			idxPath := fldPath.Index(i)
			if s, err := machine.SSH(); err == nil {
				if gpu.IsEnable(machine.Labels) {
					if !gpu.MachineIsSupport(s) {
						allErrs = append(allErrs, field.Invalid(idxPath.Child("labels"), machine.Labels, "don't has GPU card"))
					}
				}
			}
		}
	}

	return allErrs
}

func ValidateClusterFeature(spec *platform.ClusterSpec, fldPath *field.Path) field.ErrorList {
	features := spec.Features
	allErrs := field.ErrorList{}
	if features.CSIOperator != nil {
		allErrs = append(allErrs, ValidateCSIOperator(features.CSIOperator, fldPath.Child("csiOperator"))...)
	}
	if features.IPVS != nil {
		allErrs = append(allErrs, ValidateIPVS(spec, features.IPVS, fldPath.Child("ipvs"))...)
	}

	return allErrs
}

func ValidateCSIOperator(csioperator *platform.CSIOperatorFeature, fldPath *field.Path) field.ErrorList {
	return utilvalidation.ValidateEnum(csioperator.Version, fldPath.Child("version"), csioperatorimage.Versions())
}

func ValidateIPVS(spec *platform.ClusterSpec, ipvs *bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if *ipvs {
		if spec.ServiceCIDR == nil {
			allErrs = append(allErrs, field.Invalid(fldPath, ipvs, "ClusterCIDR is not allowed empty string when enable ipvs"))
		}
	}
	return allErrs
}
