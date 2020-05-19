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

package chart

import (
	"context"
	"fmt"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	"tkestack.io/tke/api/registry"
)

// ValidateChartName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateChartName = apimachineryvalidation.NameIsDNSLabel

// ValidateChart tests if required fields in the message are set.
func ValidateChart(ctx context.Context, chart *registry.Chart, registryClient *registryinternalclient.RegistryClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&chart.ObjectMeta, true, ValidateChartName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if chart.Spec.Name == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("name"), "must specify name"))
	}
	if chart.Spec.ChartGroupName == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("chartGroupName"), "must specify chartGroupName"))
	}

	if chart.Spec.ChartGroupName != "" && chart.Spec.Name != "" {
		chartGroupList, err := registryClient.ChartGroups().List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", chart.Spec.TenantID, chart.Spec.ChartGroupName),
		})
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("chartGroupName"), err))
		} else if len(chartGroupList.Items) == 0 {
			allErrs = append(allErrs, field.NotFound(fldSpecPath.Child("chartGroupName"), chart.Spec.ChartGroupName))
		} else {
			chartGroup := chartGroupList.Items[0]
			if chart.ObjectMeta.Namespace != chartGroup.ObjectMeta.Name {
				allErrs = append(allErrs, field.NotFound(field.NewPath("metadata", "namespace"), chart.ObjectMeta.Namespace))
			}

			chartList, err := registryClient.Charts(chartGroup.ObjectMeta.Name).List(ctx, metav1.ListOptions{
				FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s,spec.chartGroupName=%s", chart.Spec.TenantID, chart.Spec.Name, chart.Spec.ChartGroupName),
			})
			if err != nil {
				allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("name"), err))
			} else if len(chartList.Items) > 0 {
				allErrs = append(allErrs, field.Duplicate(fldSpecPath.Child("name"), chart.Spec.TenantID))
			}
		}
	}

	visibilities := sets.NewString(string(registry.VisibilityPrivate), string(registry.VisibilityPublic))
	if !visibilities.Has(string(chart.Spec.Visibility)) {
		allErrs = append(allErrs, field.NotSupported(fldSpecPath.Child("visibility"), chart.Spec.Visibility, visibilities.List()))
	}

	return allErrs
}

// ValidateChartUpdate tests if required fields in the chart are set
// during an update.
func ValidateChartUpdate(ctx context.Context, chart *registry.Chart, old *registry.Chart) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&chart.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))

	if chart.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}

	if chart.Spec.Name != old.Spec.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "name"), chart.Spec.Name, "disallowed change the name"))
	}

	if chart.Spec.ChartGroupName != old.Spec.ChartGroupName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "chartGroupName"), chart.Spec.Name, "disallowed change the namespace"))
	}

	return allErrs
}
