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

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	machineprovider "tkestack.io/tke/pkg/platform/provider/machine"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

const MaxTimeOffset = 5 * 300

// ValidateMachine validates a given machine.
func ValidateMachine(ctx context.Context, machine *platform.Machine, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	fldPath := field.NewPath("spec")
	allErrs := apimachineryvalidation.ValidateObjectMeta(&machine.ObjectMeta, false, apimachineryvalidation.NameIsDNSLabel, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateMachineSpecType(machine.Spec.Type, fldPath.Child("type"))...)
	allErrs = append(allErrs, ValidateClusterName(ctx, machine.Spec.ClusterName, fldPath.Child("clusterName"), platformClient)...)
	p, err := machineprovider.GetProvider(machine.Spec.Type)
	if err != nil {
		allErrs = append(allErrs, field.NotFound(fldPath, machine.Spec.Type))
	}

	allErrs = append(allErrs, p.Validate(ctx, machine)...)

	return allErrs
}

// ValidateMachineUpdate tests if an update to a machine is valid.
func ValidateMachineUpdate(ctx context.Context, machine *platform.Machine, oldMachine *platform.Machine) field.ErrorList {
	fldPath := field.NewPath("spec")
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&machine.ObjectMeta, &oldMachine.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.Type, oldMachine.Spec.Type, fldPath.Child("type"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.ClusterName, oldMachine.Spec.ClusterName, fldPath.Child("clusterName"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.IP, oldMachine.Spec.IP, fldPath.Child("ip"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.Labels, oldMachine.Spec.Labels, fldPath.Child("labels"))...)
	allErrs = append(allErrs, apimachineryvalidation.ValidateImmutableField(machine.Spec.Taints, oldMachine.Spec.Taints, fldPath.Child("taints"))...)
	allErrs = append(allErrs, ValidateMachineSpecType(machine.Spec.Type, fldPath.Child("type"))...)
	p, err := machineprovider.GetProvider(machine.Spec.Type)
	if err != nil {
		allErrs = append(allErrs, field.NotFound(fldPath, machine.Spec.Type))
	}

	allErrs = append(allErrs, p.ValidateUpdate(ctx, machine, oldMachine)...)

	return allErrs
}

// ValidateMachineSpecType validates a given type and call provider.Validate.
func ValidateMachineSpecType(machineType string, fldPath *field.Path) field.ErrorList {
	return utilvalidation.ValidateEnum(machineType, fldPath, machineprovider.Providers())
}

// ValidateClusterName validates a given clusterName and return cluster if exists.
func ValidateClusterName(ctx context.Context, clusterName string, fldPath *field.Path, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := field.ErrorList{}

	if clusterName == "" {
		allErrs = append(allErrs, field.Required(fldPath, "must specify cluster name"))
	}

	return allErrs
}
