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
	"os"
	"path"
	"path/filepath"
	"strings"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/types"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

// ValidateCluster validates a given Cluster.
func ValidateCluster(cluster *types.Cluster) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&cluster.ObjectMeta, false, apimachineryvalidation.NameIsDNSLabel, field.NewPath("metadata"))

	allErrs = append(allErrs, ValidatClusterSpec(&cluster.Spec, field.NewPath("spec"), true)...)
	p, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return append(allErrs, field.NotFound(field.NewPath("spec").Child("type"), cluster.Spec.Type))
	}
	allErrs = append(allErrs, p.Validate(cluster)...)

	return allErrs
}

// ValidateClusterUpdate tests if an update to a cluster is valid.
func ValidateClusterUpdate(cluster *types.Cluster, oldCluster *types.Cluster) field.ErrorList {
	fldPath := field.NewPath("spec")

	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&cluster.ObjectMeta, &oldCluster.ObjectMeta, field.NewPath("metadata"))

	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.Type, oldCluster.Spec.Type, fldPath.Child("type"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.NetworkDevice, oldCluster.Spec.NetworkDevice, fldPath.Child("networkDevice"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.ClusterCIDR, oldCluster.Spec.ClusterCIDR, fldPath.Child("clusterCIDR"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.DNSDomain, oldCluster.Spec.DNSDomain, fldPath.Child("dnsDomain"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.DockerExtraArgs, oldCluster.Spec.DockerExtraArgs, fldPath.Child("dockerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.KubeletExtraArgs, oldCluster.Spec.KubeletExtraArgs, fldPath.Child("kubeletExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.APIServerExtraArgs, oldCluster.Spec.APIServerExtraArgs, fldPath.Child("apiServerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.ControllerManagerExtraArgs, oldCluster.Spec.ControllerManagerExtraArgs, fldPath.Child("controllerManagerExtraArgs"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(cluster.Spec.SchedulerExtraArgs, oldCluster.Spec.SchedulerExtraArgs, fldPath.Child("schedulerExtraArgs"))...)

	allErrs = append(allErrs, ValidatClusterSpec(&cluster.Spec, field.NewPath("spec"), false)...)
	p, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return append(allErrs, field.NotFound(field.NewPath("spec").Child("type"), cluster.Spec.Type))
	}
	allErrs = append(allErrs, p.ValidateUpdate(cluster, oldCluster)...)

	return allErrs
}

// ValidateCluster validates a given ClusterSpec.
func ValidatClusterSpec(spec *platform.ClusterSpec, fldPath *field.Path, validateMachine bool) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateClusterType(spec.Type, fldPath.Child("type"))...)
	allErrs = append(allErrs, ValidateClusterFeature(&spec.Features, fldPath.Child("features"))...)

	return allErrs
}

// ValidateClusterType validates a given type.
func ValidateClusterType(clusterType string, fldPath *field.Path) field.ErrorList {
	return utilvalidation.ValidateEnum(clusterType, fldPath, clusterprovider.Providers())
}

// ValidateClusterFeature validates a given ClusterFeature.
func ValidateClusterFeature(feature *platform.ClusterFeature, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateGPUType(feature.GPUType, fldPath.Child("gpuType"))...)
	allErrs = append(allErrs, ValidateHA(feature.HA, fldPath.Child("ha"))...)
	allErrs = append(allErrs, ValidateFiles(feature.Files, fldPath.Child("files"))...)
	allErrs = append(allErrs, ValidateHooks(feature.Hooks, fldPath.Child("hooks"), feature.Files, fldPath.Child("files"))...)

	return allErrs
}

// ValidateGPUType validates a given GPUType.
func ValidateGPUType(gpuType *platform.GPUType, fldPath *field.Path) field.ErrorList {
	if gpuType == nil {
		return field.ErrorList{}
	}
	return utilvalidation.ValidateEnum(*gpuType, fldPath.Child("gpuType"),
		[]interface{}{
			platform.GPUPhysical,
			platform.GPUVirtual,
		})
}

// ValidateFiles validates a given Files.
func ValidateFiles(files []platform.File, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, file := range files {
		fldPath := fldPath.Index(i).Child("src")
		s, err := os.Stat(file.Src)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, file.Src, err.Error()))
			continue
		}
		if !s.Mode().IsRegular() && !s.Mode().IsDir() {
			allErrs = append(allErrs, field.Invalid(fldPath, file.Src, fmt.Sprintf("must be a regular file or directory")))
			continue
		}
	}

	return allErrs
}

// ValidateHooks validates a given Hooks.
func ValidateHooks(hooks map[platform.HookType]string, fldPath *field.Path, files []platform.File, filesFldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if hooks == nil {
		return allErrs
	}

	hookMap := make(map[string]bool)
	filesMap := make(map[string]string)

	for _, f := range files {
		s, err := os.Stat(f.Src)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(filesFldPath, f.Src, fmt.Sprintf("get %s status failed: %v", f.Src, err)))
			continue
		}
		if s.Mode().IsRegular() {
			hookMap[f.Dst] = true
		} else {
			filesMap[filepath.Clean(f.Dst)] = f.Src
		}
	}

	for k, v := range hooks {
		exist1, exist2 := true, true
		fldPath := fldPath.Key(string(k))
		allErrs = append(allErrs, utilvalidation.ValidateEnum(k, fldPath,
			[]interface{}{
				platform.HookPreInstall,
				platform.HookPostInstall,
				platform.HookPreUpgrade,
				platform.HookPostUpgrade,
				platform.HookPreClusterInstall,
				platform.HookPostClusterInstall,
				platform.HookPreClusterDelete,
				platform.HookPostClusterUpgrade,
			})...)

		cmd := strings.Split(v, " ")[0]
		if _, ok := hookMap[cmd]; !ok {
			exist1 = false
		}
		_, err := os.Stat(path.Join(filesMap[filepath.Dir(cmd)], filepath.Base(cmd)))
		if err != nil {
			exist2 = false
		}
		if !exist1 && !exist2 {
			allErrs = append(allErrs, field.Invalid(fldPath, cmd, fmt.Sprintf("hook file is not exists in %s", filesFldPath.String())))
		}
	}

	return allErrs
}

// ValidateHA validates a given HA.
func ValidateHA(ha *platform.HA, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ha == nil {
		return allErrs
	}

	if ha.TKEHA != nil {
		for _, msg := range validation.IsValidIP(ha.TKEHA.VIP) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("tke").Child("vip"), ha.TKEHA.VIP, msg))

		}

		if ha.TKEHA.VRID != nil {
			if *ha.TKEHA.VRID < 1 || *ha.TKEHA.VRID > 255 {
				msg := "must be a valid vrid, range [1, 255]"
				allErrs = append(allErrs, field.Invalid(fldPath.Child("tke").Child("vrid"), ha.TKEHA.VRID, msg))
			}
		}
	}
	if ha.ThirdPartyHA != nil {
		for _, msg := range validation.IsValidIP(ha.ThirdPartyHA.VIP) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("thirdParty").Child("vip"), ha.ThirdPartyHA.VIP, msg))

		}
		for _, msg := range validation.IsValidPortNum(int(ha.ThirdPartyHA.VPort)) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("thirdParty").Child("vport"), ha.ThirdPartyHA.VPort, msg))
		}
	}

	return allErrs
}
