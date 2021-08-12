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
	"fmt"
	"net"
	"strings"

	k8serror "k8s.io/apimachinery/pkg/api/errors"
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
	"tkestack.io/tke/pkg/util/ipallocator"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/validation"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

var (
	nodePodNumAvails        = []int32{16, 32, 64, 128, 256}
	clusterServiceNumAvails = []int32{32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768}
)

// ValidateCluster validates a given Cluster.
func ValidateCluster(platformClient platformv1client.PlatformV1Interface, obj *types.Cluster) field.ErrorList {
	allErrs := ValidatClusterSpec(platformClient, obj.Name, &obj.Spec, field.NewPath("spec"), obj.Status.Phase)

	return allErrs
}

// ValidatClusterSpec validates a given ClusterSpec.
func ValidatClusterSpec(platformClient platformv1client.PlatformV1Interface, clusterName string, spec *platform.ClusterSpec, fldPath *field.Path, phase platform.ClusterPhase) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateClusterSpecVersion(platformClient, clusterName, spec.Version, fldPath.Child("version"), phase)...)
	allErrs = append(allErrs, ValidateCIDRs(spec, fldPath)...)
	allErrs = append(allErrs, ValidateClusterProperty(spec, fldPath.Child("properties"))...)
	allErrs = append(allErrs, ValidateClusterMachines(spec.Machines, fldPath.Child("machines"))...)
	allErrs = append(allErrs, ValidateClusterFeature(spec, fldPath.Child("features"))...)

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
func ValidateCIDRs(spec *platform.ClusterSpec, specPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	var clusterCIDR, serviceCIDR *net.IPNet

	checkFunc := func(path *field.Path, cidr string) {
		cidrs := strings.Split(cidr, ",")
		dualStackEnabled := spec.Features.IPv6DualStack
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
			} else {
				serviceCIDR = cidrX
			}
		}
	}

	fldPath := specPath.Child("clusterCIDR")
	cidr := spec.ClusterCIDR
	if len(cidr) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath, cidr, "ClusterCIDR is empty string"))
	} else {
		checkFunc(fldPath, cidr)
	}

	fldPath = specPath.Child("serviceCIDR")
	if spec.ServiceCIDR != nil {
		cidr = *spec.ServiceCIDR
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

// ValidateClusterMachines validates a given CluterMachines.
func ValidateClusterMachines(machines []platform.ClusterMachine, fldPath *field.Path) field.ErrorList {
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
