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

package chartgroup

import (
	"fmt"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
)

// ValidateChartGroupName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateChartGroupName = apimachineryvalidation.NameIsDNSLabel

// ValidateChartGroup tests if required fields in the chart group are set.
func ValidateChartGroup(chartGroup *registry.ChartGroup, registryClient *registryinternalclient.RegistryClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&chartGroup.ObjectMeta, false, ValidateChartGroupName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if chartGroup.Spec.Name == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("name"), "must specify name"))
	} else {
		chartGroupList, err := registryClient.ChartGroups().List(metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", chartGroup.Spec.TenantID, chartGroup.Spec.Name),
		})
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("name"), err))
		} else if len(chartGroupList.Items) > 0 {
			allErrs = append(allErrs, field.Duplicate(fldSpecPath.Child("name"), chartGroup.Spec.Name))
		}
	}

	visibilities := sets.NewString(string(registry.VisibilityPrivate), string(registry.VisibilityPublic))
	if !visibilities.Has(string(chartGroup.Spec.Visibility)) {
		allErrs = append(allErrs, field.NotSupported(fldSpecPath.Child("visibility"), chartGroup.Spec.Visibility, visibilities.List()))
	}

	return allErrs
}

// ValidateChartGroupUpdate tests if required fields in the chart group are set during
// an update.
func ValidateChartGroupUpdate(chartGroup *registry.ChartGroup, old *registry.ChartGroup) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&chartGroup.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))

	if chartGroup.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"), chartGroup.Spec.TenantID, "disallowed change the tenant"))
	}

	if chartGroup.Spec.Name != old.Spec.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "name"), chartGroup.Spec.Name, "disallowed change the name"))
	}

	return allErrs
}
