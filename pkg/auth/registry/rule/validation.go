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

package rule

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
)

// ValidateRule tests if required fields in the policy are set.
func ValidateRule(rule *auth.Rule) field.ErrorList {
	allErrs := field.ErrorList{}
	fldStatPath := field.NewPath("spec")

	if rule.Spec.PType == "" {
		allErrs = append(allErrs, field.Required(fldStatPath.Child("ptype"), "must specify ptype"))
	}

	if rule.Spec.V0 == "" {
		allErrs = append(allErrs, field.Required(fldStatPath.Child("v0"), "must specify v0"))
	}

	return allErrs
}

// ValidateRuleUpdate tests if required fields in the policy are set during
// an update.
func ValidateRuleUpdate(rule *auth.Rule, old *auth.Rule) field.ErrorList {
	//allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&policy.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateRule(rule)...)
	return nil
}
