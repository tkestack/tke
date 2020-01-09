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

package policy

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/validation"
)

// ValidatePolicyName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidatePolicyName = apiMachineryValidation.NameIsDNSLabel

// ValidatePolicy tests if required fields in the policy are set.
func ValidatePolicy(policy *auth.Policy, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&policy.ObjectMeta, false, ValidatePolicyName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if err := validation.IsDisplayName(policy.Spec.DisplayName); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("displayName"), policy.Spec.DisplayName, err.Error()))
	}

	if policy.Spec.Type == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("type"), "must specify type"))
	} else if policy.Spec.Type != auth.PolicyCustom && policy.Spec.Type != auth.PolicyDefault {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("type"), policy.Spec.Type, "must specify one of: `custom` or `default`"))
	}

	if policy.Spec.Category == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("category"), policy.Spec.Category))
	}

	fldStmtPath := field.NewPath("spec", "statement")
	if len(policy.Spec.Statement.Actions) == 0 {
		allErrs = append(allErrs, field.Required(fldStmtPath.Child("actions"), "must specify actions"))
	}

	if len(policy.Spec.Statement.Resources) == 0 {
		allErrs = append(allErrs, field.Required(fldStmtPath.Child("resources"), "must specify resources"))
	}

	if policy.Spec.Statement.Effect == "" {
		allErrs = append(allErrs, field.Required(fldStmtPath.Child("effect"), "must specify effect"))
	} else if policy.Spec.Statement.Effect != auth.Allow && policy.Spec.Statement.Effect != auth.Deny {
		allErrs = append(allErrs, field.Invalid(fldStmtPath.Child("effect"), policy.Spec.Statement.Effect, "must specify one of: `allow` or `deny`"))
	}

	fldUserPath := field.NewPath("status", "users")
	for i, subj := range policy.Status.Users {
		if subj.ID == "" {
			allErrs = append(allErrs, field.Required(fldUserPath, "must specify subject id"))
			continue
		}

		if subj.Name == "" {
			val, err := authClient.Users().Get(util.CombineTenantAndName(policy.Spec.TenantID, subj.ID), metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					allErrs = append(allErrs, field.NotFound(fldUserPath, subj.ID))
				} else {
					allErrs = append(allErrs, field.InternalError(fldUserPath, err))
				}
			} else {
				if val.Spec.TenantID != policy.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldUserPath, subj.ID, "must in the same tenant with the policy"))
				} else {
					policy.Status.Users[i].Name = val.Spec.Name
				}
			}
		}

	}

	fldGroupPath := field.NewPath("status", "groups")
	for i, subj := range policy.Status.Groups {
		if subj.ID == "" {
			allErrs = append(allErrs, field.Required(fldGroupPath, "must specify id"))
			continue
		}

		if subj.Name == "" {
			val, err := authClient.Groups().Get(util.CombineTenantAndName(policy.Spec.TenantID, subj.ID), metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					allErrs = append(allErrs, field.NotFound(fldGroupPath, subj.ID))
				} else {
					allErrs = append(allErrs, field.InternalError(fldGroupPath, err))
				}
			} else {
				if val.Spec.TenantID != policy.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldGroupPath, subj.ID, "must in the same tenant with the policy"))
				} else {
					policy.Status.Groups[i].Name = val.Spec.DisplayName
				}
			}
		}

	}

	return allErrs
}

// ValidatePolicyUpdate tests if required fields in the policy are set during
// an update.
func ValidatePolicyUpdate(policy *auth.Policy, old *auth.Policy, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&policy.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidatePolicy(policy, authClient)...)

	fldSpecPath := field.NewPath("spec")
	if policy.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), policy.Spec.TenantID, "disallowed change the tenant"))
	}

	return allErrs
}
