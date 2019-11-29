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
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
)

// ValidatePolicyName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidatePolicyName = apiMachineryValidation.NameIsDNSLabel

// ValidatePolicy tests if required fields in the policy are set.
func ValidatePolicy(policy *auth.Policy) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&policy.ObjectMeta, false, ValidatePolicyName, field.NewPath("metadata"))

	fldStatPath := field.NewPath("spec", "statement")

	if len(policy.Spec.Statement.Actions) == 0 {
		allErrs = append(allErrs, field.Required(fldStatPath.Child("actions"), "must specify actions"))
	}

	if len(policy.Spec.Statement.Resources) == 0 {
		allErrs = append(allErrs, field.Required(fldStatPath.Child("resources"), "must specify resources"))
	}

	if policy.Spec.Statement.Effect == "" {
		allErrs = append(allErrs, field.Required(fldStatPath.Child( "effect"), "must specify resources"))
	} else if policy.Spec.Statement.Effect != auth.Allow && policy.Spec.Statement.Effect != auth.Deny {
		allErrs = append(allErrs, field.Invalid(fldStatPath.Child( "effect"), policy.Spec.Statement.Effect, "must specify one of: `allow` or `deny`"))
	}

	return allErrs
}

// ValidatePolicyUpdate tests if required fields in the policy are set during
// an update.
func ValidatePolicyUpdate(policy *auth.Policy, old *auth.Policy) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&policy.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidatePolicy(policy)...)
	return allErrs
}
