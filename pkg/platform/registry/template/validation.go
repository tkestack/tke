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
	"fmt"
	"reflect"

	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
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

// Validate validates a given template.
func ValidateTemplate(template *platform.Template) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&template.ObjectMeta, false, ValidateName, field.NewPath("metadata"))

	v := reflect.ValueOf(template.Spec.Content)
	count := v.NumField()
	var contentCount int
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if !f.IsNil() {
			contentCount++
		}
	}
	if contentCount != 1 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "content"), "must specify one content"))
	}
	if template.Spec.Content.Deployment != nil {
		if template.Spec.Type != Deployment {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "content", "spec", Deployment), fmt.Sprintf("must specify spec.type is %s", Deployment)))
		}
	}
	if template.Spec.Content.StatefulSet != nil {
		if template.Spec.Type != StatefulSet {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "content", "spec", StatefulSet), fmt.Sprintf("must specify spec.type is %s", StatefulSet)))
		}
	}
	if template.Spec.Content.DaemonSet != nil {
		if template.Spec.Type != DaemonSet {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "content", "spec", DaemonSet), fmt.Sprintf("must specify spec.type is %s", DaemonSet)))
		}
	}
	if template.Spec.Content.Job != nil {
		if template.Spec.Type != Job {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "content", "spec", Job), fmt.Sprintf("must specify spec.type is %s", Job)))
		}
	}
	if template.Spec.Content.CronJob != nil {
		if template.Spec.Type != CronJob {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "content", "spec", CronJob), fmt.Sprintf("must specify spec.type is %s", CronJob)))
		}
	}
	if template.Spec.Content.Tapp != nil {
		if template.Spec.Type != TApp {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "content", "spec", TApp), fmt.Sprintf("must specify spec.type is %s", TApp)))
		}
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
