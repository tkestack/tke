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

package collector

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/monitor"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apimachineryvalidation.ValidateNamespaceName

// ValidateCollector tests if required fields in the cluster are set.
func ValidateCollector(collector *monitor.Collector) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&collector.ObjectMeta, false, ValidateName, field.NewPath("metadata"))

	if len(collector.Spec.ClusterName) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusterName"), "must specify a cluster name"))
	}

	return allErrs
}

// ValidateCollectorUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateCollectorUpdate(collector *monitor.Collector, old *monitor.Collector) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&collector.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateCollector(collector)...)

	if collector.Spec.ClusterName != old.Spec.ClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), collector.Spec.ClusterName, "disallowed change the cluster name"))
	}

	if collector.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"), collector.Spec.TenantID, "disallowed change the tenant"))
	}

	if collector.Spec.Type != old.Spec.Type {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "type"), collector.Spec.Type, "type is immutable"))
	}

	if collector.Status.Phase == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "phase"), string(collector.Status.Phase)))
	}

	return allErrs
}
