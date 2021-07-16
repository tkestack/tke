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

package cluster

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	utilmath "tkestack.io/tke/pkg/util/math"
	"tkestack.io/tke/pkg/util/ssh"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

const MaxTimeOffset = 5 * 300

// ValidateCluster validates a given ClusterSpec.
func ValidatClusterSpec(spec *platform.ClusterSpec, fldPath *field.Path, validateMachine bool) field.ErrorList {
	allErrs := field.ErrorList{}
	if validateMachine {
		allErrs = append(allErrs, ValidateClusterMachines(spec.Machines, fldPath.Child("machines"))...)
	}
	allErrs = append(allErrs, ValidateClusterFeature(&spec.Features, fldPath.Child("features"))...)

	return allErrs
}

// ValidateClusterMachines validates a given CluterMachines.
func ValidateClusterMachines(machines []platform.ClusterMachine, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if machines == nil {
		return allErrs
	}

	var masters []*ssh.SSH
	for i, one := range machines {
		sshErrors := ValidateSSH(fldPath.Index(i), one.IP, int(one.Port), one.Username, one.Password, one.PrivateKey, one.PassPhrase)
		if sshErrors != nil {
			allErrs = append(allErrs, sshErrors...)
		} else {
			master, _ := one.SSH()
			masters = append(masters, master)
		}
	}

	if len(masters) == len(machines) {
		allErrs = append(allErrs, ValidateMasterTimeOffset(fldPath, masters)...)

	}

	return allErrs
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
			allErrs = append(allErrs, field.Invalid(fldPath, file.Src, "must be a regular file or directory"))
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

// ValidateSSH validates a given ssh config.
func ValidateSSH(fldPath *field.Path, ip string, port int, user string, password []byte, privateKey []byte, passPhrase []byte) field.ErrorList {
	allErrs := field.ErrorList{}

	for _, msg := range validation.IsValidIP(ip) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ip"), ip, msg))

	}
	for _, msg := range validation.IsValidPortNum(port) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("port"), port, msg))
	}
	if password == nil && privateKey == nil {
		allErrs = append(allErrs, field.Required(fldPath, "must specify password or privateKey"))
	}

	if len(allErrs) != 0 {
		return allErrs
	}

	sshConfig := &ssh.Config{
		User:        user,
		Host:        ip,
		Port:        port,
		Password:    string(password),
		PrivateKey:  privateKey,
		PassPhrase:  passPhrase,
		DialTimeOut: time.Second,
		Retry:       0,
	}
	s, err := ssh.New(sshConfig)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, "", err.Error()))
	} else {
		output, err := s.CombinedOutput("whoami")
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, "", err.Error()))
		}
		if strings.TrimSpace(string(output)) != "root" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("user"), user, `must be root or set sudo without password`))
		}
	}

	return allErrs
}
