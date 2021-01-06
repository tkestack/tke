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

package application

import (
	"context"
	"fmt"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/application"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	"yunion.io/x/pkg/util/sets"
)

// ValidateApplicationName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateApplicationName = apimachineryvalidation.NameIsDNSLabel

// ValidateApplication tests if required fields in the message are set.
func ValidateApplication(ctx context.Context, app *application.App, applicationClient applicationinternalclient.ApplicationInterface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&app.ObjectMeta, true, ValidateApplicationName, field.NewPath("metadata"))

	fldMetadataPath := field.NewPath("metadata")
	if app.ObjectMeta.Namespace == "" {
		allErrs = append(allErrs, field.Required(fldMetadataPath.Child("namespace"), "must specify namespace"))
	}

	fldSpecPath := field.NewPath("spec")
	if app.Spec.Name == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("name"), "must specify name"))
	}
	if app.Spec.TargetCluster == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("targetCluster"), "must specify targetCluster"))
	}
	types := sets.NewString(string(application.AppTypeHelmV3))
	if !types.Has(string(app.Spec.Type)) {
		allErrs = append(allErrs, field.NotSupported(fldSpecPath.Child("type"), app.Spec.Type, types.List()))
	}
	if app.Spec.Type == application.AppTypeHelmV3 {
		fldChartPath := field.NewPath("spec", "chart")
		if app.Spec.Chart.ChartGroupName == "" {
			allErrs = append(allErrs, field.Required(fldChartPath.Child("chartGroupName"), "must specify chartGroupName"))
		}
		if app.Spec.Chart.ChartName == "" {
			allErrs = append(allErrs, field.Required(fldChartPath.Child("chartName"), "must specify chartName"))
		}
		if app.Spec.Chart.TenantID == "" {
			allErrs = append(allErrs, field.Required(fldChartPath.Child("tenantID"), "must specify chart tenantID"))
		}
	}

	applicationList, err := applicationClient.Apps(app.ObjectMeta.Namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s,spec.targetCluster=%s", app.Spec.TenantID, app.Spec.Name, app.Spec.TargetCluster),
	})
	if err != nil {
		allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("name"), err))
	} else if len(applicationList.Items) > 0 {
		allErrs = append(allErrs, field.Duplicate(fldSpecPath.Child("name"), app.Spec.Name))
	}

	return allErrs
}

// ValidateApplicationUpdate tests if required fields in the app are set
// during an update.
func ValidateApplicationUpdate(ctx context.Context, app *application.App, old *application.App) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&app.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))

	if app.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}

	if app.Spec.Type != old.Spec.Type {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "type"), app.Spec.Type, "disallowed change the type"))
	}

	if app.Spec.Name != old.Spec.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "name"), app.Spec.Name, "disallowed change the name"))
	}

	if app.Spec.TargetCluster != old.Spec.TargetCluster {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "targetCluster"), app.Spec.TargetCluster, "disallowed change the targetCluster"))
	}

	return allErrs
}
