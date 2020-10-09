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
	"fmt"
	"net"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"

	netutils "k8s.io/utils/net"
	csioperatorimage "tkestack.io/tke/pkg/platform/provider/baremetal/phases/csioperator/images"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/gpu"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/ipallocator"
	"tkestack.io/tke/pkg/util/validation"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

var (
	nodePodNumAvails        = []int32{16, 32, 64, 128, 256}
	clusterServiceNumAvails = []int32{32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768}
)

// ValidateCluster validates a given Cluster.
func ValidateCluster(obj *types.Cluster) field.ErrorList {
	allErrs := ValidatClusterSpec(&obj.Spec, field.NewPath("spec"), obj.Status.Phase)

	return allErrs
}

// ValidatClusterSpec validates a given ClusterSpec.
func ValidatClusterSpec(spec *platform.ClusterSpec, fldPath *field.Path, phase platform.ClusterPhase) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateClusterSpecVersion(spec.Version, fldPath.Child("version"), phase)...)
	allErrs = append(allErrs, ValidateCIDRs(spec, fldPath)...)
	allErrs = append(allErrs, ValidateClusterProperty(spec, fldPath.Child("properties"))...)
	allErrs = append(allErrs, ValidateClusterMachines(spec.Machines, fldPath.Child("machines"))...)
	allErrs = append(allErrs, ValidateClusterFeature(&spec.Features, fldPath.Child("features"))...)

	return allErrs
}

// ValidateClusterSpecVersion validates a given version.
func ValidateClusterSpecVersion(version string, fldPath *field.Path, phase platform.ClusterPhase) field.ErrorList {
	allErrs := field.ErrorList{}
	if phase == platform.ClusterInitializing {
		allErrs = utilvalidation.ValidateEnum(version, fldPath, spec.K8sVersions)
	}

	return allErrs
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
	cidr = *spec.ServiceCIDR
	if spec.ServiceCIDR != nil {
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

func ValidateClusterFeature(feature *platform.ClusterFeature, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if feature.CSIOperator != nil {
		allErrs = append(allErrs, ValidateCSIOperator(feature.CSIOperator, fldPath.Child("csiOperator"))...)
	}

	return allErrs
}

func ValidateCSIOperator(csioperator *platform.CSIOperatorFeature, fldPath *field.Path) field.ErrorList {
	return utilvalidation.ValidateEnum(csioperator.Version, fldPath.Child("version"), csioperatorimage.Versions())
}
