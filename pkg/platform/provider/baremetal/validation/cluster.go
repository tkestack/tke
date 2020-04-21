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
	"strconv"

	"github.com/thoas/go-funk"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/gpu"
	"tkestack.io/tke/pkg/spec"
)

var (
	nodePodNumAvails        = []int32{16, 32, 64, 128, 256}
	clusterServiceNumAvails = []int32{32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768}
)

// ValidateCluster validates a given Cluster.
func ValidateCluster(obj *platform.Cluster) field.ErrorList {
	allErrs := ValidatClusterSpec(&obj.Spec, field.NewPath("spec"), obj.Status.Phase)

	return allErrs
}

// ValidateCluster validates a given ClusterSpec.
func ValidatClusterSpec(spec *platform.ClusterSpec, fldPath *field.Path, phase platform.ClusterPhase) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateClusterSpecVersion(spec.Version, fldPath.Child("version"), phase)...)
	allErrs = append(allErrs, ValidateCIDR(spec.ClusterCIDR, fldPath.Child("clusterCIDR"))...)
	allErrs = append(allErrs, ValidateClusterProperty(spec.Properties, fldPath.Child("properties"))...)
	allErrs = append(allErrs, ValidateClusterMachines(spec.Machines, fldPath.Child("machines"))...)

	return allErrs
}

// ValidateClusterSpecVersion validates a given version.
func ValidateClusterSpecVersion(version string, fldPath *field.Path, phase platform.ClusterPhase) field.ErrorList {
	allErrs := field.ErrorList{}

	if phase == platform.ClusterInitializing && !funk.Contains(spec.K8sVersions, version) {
		allErrs = append(allErrs, field.NotSupported(fldPath, version, spec.K8sVersions))
	}

	return allErrs
}

// ValidateCIDR validates a given cidr.
func ValidateCIDR(cidr string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if cidr == "" {
		allErrs = append(allErrs, field.Required(fldPath, ""))
	} else {
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, cidr, err.Error()))
		}
	}

	return allErrs
}

// ValidateClusterProperty validates a given ClusterProperty.
func ValidateClusterProperty(property platform.ClusterProperty, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if property.MaxNodePodNum == nil {
		allErrs = append(allErrs, field.Required(fldPath, fmt.Sprintf("validate values are %v", nodePodNumAvails)))
	} else {
		if !funk.Contains(nodePodNumAvails, *property.MaxNodePodNum) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("maxNodePodNum"), *property.MaxNodePodNum,
				funk.Map(nodePodNumAvails, func(x int) string {
					return strconv.Itoa(x)
				}).([]string)))
		}
	}
	if property.MaxClusterServiceNum == nil {
		allErrs = append(allErrs, field.Required(fldPath, fmt.Sprintf("validate values are %v", clusterServiceNumAvails)))
	} else {
		if !funk.Contains(clusterServiceNumAvails, *property.MaxClusterServiceNum) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("maxClusterServiceNum"), *property.MaxNodePodNum,
				funk.Map(clusterServiceNumAvails, func(x int) string {
					return strconv.Itoa(x)
				}).([]string)))
		}
		if *property.MaxClusterServiceNum < 10 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("maxClusterServiceNum"), *property.MaxClusterServiceNum,
				"must be greater than or equal to 10 because kubeadm need the 10th ip"))
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
