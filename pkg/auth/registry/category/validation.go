/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package category

import (
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/util/validation"
)

// ValidateCategory tests if required fields in the signing key are set.
func ValidateCategory(category *auth.Category) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&category.ObjectMeta, false, apiMachineryValidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if err := validation.IsDNS1123Name(category.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("categoryName"), category.Name, err.Error()))
	}

	if category.Spec.DisplayName == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("displayName"), "must specify displayName"))
	} else {
		if err := validation.IsDisplayName(category.Spec.DisplayName); err != nil {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("displayName"), category.Spec.DisplayName, err.Error()))
		}
	}

	for _, act := range category.Spec.Actions {
		if act.Name == "" {
			allErrs = append(allErrs, field.Required(fldSpecPath.Child("actions", "name"), "act must specify name"))
		}
	}

	return allErrs
}

// ValidateCategoryUpdate tests if required fields in the session are set during
// an update.
func ValidateCategoryUpdate(category *auth.Category, oldCategory *auth.Category) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateCategory(category)...)
	allErrs = append(allErrs, apiMachineryValidation.ValidateObjectMetaUpdate(&category.ObjectMeta, &oldCategory.ObjectMeta, field.NewPath("metadata"))...)

	fldSpecPath := field.NewPath("spec")

	if category.Name != oldCategory.Name {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("categoryName"), category.Name, "disallowed change the categoryName"))
	}

	return allErrs
}
