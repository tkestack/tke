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

package template

import (
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apiMachineryValidation.ValidateNamespaceName

const (
	Deployment  = "Deployment"
	StatefulSet = "StatefulSet"
	DaemonSet   = "DaemonSet"
	Job         = "Job"
	CronJob     = "CronJob"
	TApp        = "TApp"
)

var templateTypes = sets.NewString(Deployment, StatefulSet, Deployment, DaemonSet, Job, CronJob, TApp)

// Validate validates a given template.
func ValidateTemplate(template *platform.Template) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&template.ObjectMeta, false, ValidateName, field.NewPath("metadata"))

	if !templateTypes.Has(template.Spec.Type) {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "type"), "must be one of these: Deployment, StatefulSet, DaemonSet, Job, CronJob, TApp"))
	}

	if len(template.Spec.Content) <= 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "content"), "must not be none."))
	}

	return allErrs
}

// ValidateTemplateSpecType validates a given type and call provider.Validate.
func ValidateTemplateSpecType(templateType string, fldPath *field.Path) field.ErrorList {
	return utilvalidation.ValidateEnum(templateType, fldPath, []string{"Tapp", "Deployment", "Job", "Statefulset", ""})
}

// ValidateTemplateUpdate tests if an update to a template is valid.
func ValidateTemplateUpdate(template *platform.Template, oldTemplate *platform.Template) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&template.ObjectMeta, &oldTemplate.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateTemplate(template)...)

	if template.Spec.TenantID != oldTemplate.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"), oldTemplate.Spec.TenantID, "disallowed change the tenant"))
	}

	return allErrs
}
