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

package util

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/pkg/auth/types"
)

// ValidateSubjectAccessReviewSpec validates SubjectAccessReviewSpec resource attributes and user
func ValidateSubjectAccessReviewSpec(spec types.SubjectAccessReviewSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if spec.ResourceAttributes == nil && spec.NonResourceAttributes == nil && len(spec.ResourceAttributesList) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("resourceAttributes"), spec.NonResourceAttributes, `exactly one of nonResourceAttributes or resourceAttributes or resourceAttributesList must be specified`))
	}
	if len(spec.User) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("user"), spec.User, `user must be specified`))
	}

	return allErrs
}

// ValidateSubjectAccessReview validates SubjectAccessReview
func ValidateSubjectAccessReview(sar *types.SubjectAccessReview) field.ErrorList {
	allErrs := ValidateSubjectAccessReviewSpec(sar.Spec, field.NewPath("spec"))
	return allErrs
}
