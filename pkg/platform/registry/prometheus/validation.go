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

package prometheus

import (
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apiMachineryValidation.ValidateNamespaceName

// ValidatePrometheus tests if required fields in the cluster are set.
func ValidatePrometheus(prom *platform.Prometheus) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&prom.ObjectMeta, false, ValidateName, field.NewPath("metadata"))

	if len(prom.Spec.ClusterName) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusterName"), "must specify a cluster name"))
	}

	return allErrs
}

// ValidatePrometheusUpdate tests if required fields in the namespace set are
// set during an update.
func ValidatePrometheusUpdate(prom *platform.Prometheus, old *platform.Prometheus) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&prom.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidatePrometheus(prom)...)

	if prom.Spec.ClusterName != old.Spec.ClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), prom.Spec.ClusterName, "disallowed change the cluster name"))
	}

	if prom.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"), prom.Spec.TenantID, "disallowed change the tenant"))
	}

	if prom.Status.Phase == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "phase"), string(prom.Status.Phase)))
	}

	return allErrs
}
