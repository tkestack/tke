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
	"context"
	"fmt"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
)

// _validateChartGroupName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var _validateChartGroupName = apimachineryvalidation.NameIsDNSLabel

// ValidateChartGroupCreate tests if required fields in the ChartGroup are set correctly.
func ValidateChartGroupCreate(ctx context.Context, chartGroup *business.ChartGroup,
	businessClient *businessinternalclient.BusinessClient, registryClient registryversionedclient.RegistryV1Interface) field.ErrorList {
	allErrs := validateChartGroup(chartGroup, businessClient, registryClient)

	fldName := field.NewPath("spec", "name")
	chartGroupList, err := businessClient.ChartGroups(chartGroup.Namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", chartGroup.Spec.TenantID, chartGroup.Name),
	})
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldName, fmt.Errorf("failed to check chartGroup conflicting, for %s", err)))
	} else if len(chartGroupList.Items) != 0 {
		allErrs = append(allErrs, field.Duplicate(fldName, chartGroup.Name))
	}

	return allErrs
}

// ValidateChartGroupUpdate tests if required fields in the ChartGroup are set during
// an update.
func ValidateChartGroupUpdate(ctx context.Context, chartGroup *business.ChartGroup, old *business.ChartGroup,
	businessClient *businessinternalclient.BusinessClient, registryClient registryversionedclient.RegistryV1Interface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&chartGroup.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, validateChartGroup(chartGroup, businessClient, registryClient)...)

	if chartGroup.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "tenantID"),
				chartGroup.Spec.TenantID, "disallowed change the tenant"))
	}

	if chartGroup.Status.Phase == "" {
		allErrs = append(allErrs,
			field.Required(field.NewPath("status", "phase"), "must specify a phase"))
	}

	return allErrs
}

// validateChartGroup tests if required fields in the ChartGroup are set.
func validateChartGroup(chartGroup *business.ChartGroup,
	businessClient *businessinternalclient.BusinessClient, registryClient registryversionedclient.RegistryV1Interface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&chartGroup.ObjectMeta,
		true, _validateChartGroupName, field.NewPath("metadata"))

	return allErrs
}
