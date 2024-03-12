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

package upgradejob

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/application"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apimachineryvalidation.ValidateNamespaceName

// ValidateUpgradeJob tests if required fields in the cluster are set.
func ValidateUpgradeJob(job *application.UpgradeJob) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&job.ObjectMeta, true, ValidateName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if job.Spec.TenantID == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("tenantID"), "must specify tenantID"))
	}
	/*
		if job.Spec.BatchNum == nil || *job.Spec.BatchNum < MinBatchNum {
			allErrs = append(allErrs, field.Required(fldSpecPath.Child("batchNum"), "must set batchNum correctly"))
		}
		if job.Spec.Target == "" {
			allErrs = append(allErrs, field.Required(fldSpecPath.Child("target"), "must set target correctly"))
		}
		if job.Spec.MaxFailed != nil && *job.Spec.MaxFailed < 0 {
			allErrs = append(allErrs, field.Required(fldSpecPath.Child("maxFailed"), "must set maxFailed correctly"))
		}
		if job.Spec.MaxSurge != nil && *job.Spec.MaxSurge <= 0 {
			allErrs = append(allErrs, field.Required(fldSpecPath.Child("maxSurge"), "must set maxSurge correctly"))
		}
	*/
	return allErrs
}

// ValidateUpgradeJobUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateUpgradeJobUpdate(job *application.UpgradeJob, old *application.UpgradeJob) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&job.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateUpgradeJob(job)...)

	return allErrs
}
