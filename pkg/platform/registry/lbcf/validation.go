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

package lbcf

import (
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/util/validation"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apiMachineryValidation.ValidateNamespaceName

// ValidateLBCF tests if required fields in the cluster are set.
func ValidateLBCF(obj *platform.LBCF, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&obj.ObjectMeta, false, ValidateName, field.NewPath("metadata"))
	allErrs = append(allErrs, validation.ValidateCluster(platformClient, obj.Spec.ClusterName)...)

	return allErrs
}

// ValidateLBCFUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateLBCFUpdate(lbcf *platform.LBCF, old *platform.LBCF, platformClient platforminternalclient.PlatformInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&lbcf.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateLBCF(lbcf, platformClient)...)

	if lbcf.Spec.ClusterName != old.Spec.ClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), lbcf.Spec.ClusterName, "disallowed change the cluster name"))
	}

	if lbcf.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"), lbcf.Spec.TenantID, "disallowed change the tenant"))
	}

	if lbcf.Status.Phase == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "phase"), string(lbcf.Status.Phase)))
	}

	return allErrs
}
