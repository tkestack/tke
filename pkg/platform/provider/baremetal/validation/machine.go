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
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/provider/baremetal/phases/gpu"
)

// ValidateMachine validates a given machine.
func ValidateMachine(machine *platform.Machine) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateMachineSpec(&machine.Spec, field.NewPath("spec"))...)

	return allErrs
}

// ValidateMachineSpec validates a given machine spec.
func ValidateMachineSpec(spec *platform.MachineSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	s, err := spec.SSH()
	if err == nil {
		if gpu.IsEnable(spec.Labels) {
			if !gpu.MachineIsSupport(s) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("labels"), spec.Labels, "must have GPU card if set GPU label"))
			}
		}
	}

	return allErrs
}
