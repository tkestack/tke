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

package upgradepolicy

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/application"
)

const MinBatchNum = 1
const MinBatchIntervalSeconds = int32(30)
const DefaultMaxSurge = int32(3)
const DefaultMaxFailed = int32(0)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apimachineryvalidation.ValidateNamespaceName

// ValidateUpgradePolicy tests if required fields in the cluster are set.
func ValidateUpgradePolicy(up *application.UpgradePolicy) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&up.ObjectMeta, false, ValidateName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if up.Spec.BatchNum == nil || *up.Spec.BatchNum < MinBatchNum {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("batchNum"), "must set batchNum correctly"))
	}
	if up.Spec.MaxFailed != nil && *up.Spec.MaxFailed < 0 {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("maxFailed"), "must set maxFailed correctly"))
	}
	if up.Spec.MaxSurge != nil && *up.Spec.MaxSurge <= 0 {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("maxSurge"), "must set maxSurge correctly"))
	}

	return allErrs
}

// ValidateUpgradePolicyUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateUpgradePolicyUpdate(up *application.UpgradePolicy, old *application.UpgradePolicy) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&up.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateUpgradePolicy(up)...)

	return allErrs
}
