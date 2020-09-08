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

package custompolicybinding

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

// ValidateBindingName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateBindingName = apiMachineryValidation.NameIsDNSLabel

// ValidateProjectPolicyBinding tests if required fields in the projectpolicybinding are set.
func ValidateProjectPolicyBinding(ctx context.Context, binding *auth.CustomPolicyBinding, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&binding.ObjectMeta, true, ValidateBindingName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")

	if binding.Spec.PolicyID == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("policyID"), "must specify policyID"))
	} else {
		_, err := authClient.Policies().Get(ctx, binding.Spec.PolicyID, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				allErrs = append(allErrs, field.NotFound(fldSpecPath.Child("policyID"), binding.Spec.PolicyID))
			} else {
				allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("policyID"), err))
			}
		}
	}

	if binding.Spec.Domain == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("domain"), "must specify domain"))
	}

	if binding.Spec.RulePrefix == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("rulePrefix"), "must specify rulePrefix"))
	}

	if len(binding.Spec.Resources) == 0 {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("resources"), "must specify resources"))
	}

	var validUsers []auth.Subject
	fldUserPath := field.NewPath("spec", "users")
	for i, subj := range binding.Spec.Users {
		if subj.ID == "" && subj.Name == "" {
			allErrs = append(allErrs, field.Required(fldUserPath, "must specify subject id or name"))
			continue
		}

		switch {
		case subj.ID != "" && subj.Name == "":
			val, err := authClient.Users().Get(ctx, util.CombineTenantAndName(binding.Spec.TenantID, subj.ID), metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Warn("user is not found, will removed it", log.String("policy", binding.Name), log.String("user", subj.ID))
				} else {
					allErrs = append(allErrs, field.InternalError(fldUserPath, err))
				}
			} else {
				if val.Spec.TenantID != binding.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldUserPath, subj.ID, "must in the same tenant with the policy"))
				} else {
					binding.Spec.Users[i].Name = val.Spec.Name
					validUsers = append(validUsers, binding.Spec.Users[i])
				}
			}
		case subj.ID == "" && subj.Name != "":
			user, err := util.GetUserByName(ctx, authClient, binding.Spec.TenantID, subj.Name)
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Warn("user is not found in tenant, will removed it", log.String("policy", binding.Name), log.String("user", subj.Name))
				} else {
					allErrs = append(allErrs, field.InternalError(fldUserPath, err))
				}
			} else {
				if user.Spec.TenantID != binding.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldUserPath, subj.ID, "must in the same tenant with the project"))
				} else {
					binding.Spec.Users[i].ID = user.Spec.ID
					validUsers = append(validUsers, binding.Spec.Users[i])
				}
			}
		default:
			validUsers = append(validUsers, binding.Spec.Users[i])
		}
	}

	var validGroups []auth.Subject
	fldGroupPath := field.NewPath("spec", "groups")
	for i, subj := range binding.Spec.Groups {
		if subj.ID == "" {
			allErrs = append(allErrs, field.Required(fldGroupPath, "must specify id"))
			continue
		}

		if subj.Name == "" {
			val, err := authClient.Groups().Get(ctx, util.CombineTenantAndName(binding.Spec.TenantID, subj.ID), metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Warn("group of the policy is not found, will removed it", log.String("policy", binding.Name), log.String("group", subj.Name))
				} else {
					allErrs = append(allErrs, field.InternalError(fldGroupPath, err))
				}
			} else {
				if val.Spec.TenantID != binding.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldGroupPath, subj.ID, "must in the same tenant with the policy"))
				} else {
					binding.Spec.Groups[i].Name = val.Spec.DisplayName
					validGroups = append(validGroups, binding.Spec.Groups[i])
				}
			}
		} else {
			validGroups = append(validGroups, binding.Spec.Groups[i])
		}
	}
	if len(allErrs) == 0 {
		binding.Spec.Users = validUsers
		binding.Spec.Groups = validGroups
	}

	log.Debug("binding spec", log.Any("spec", binding.Spec))

	return allErrs
}

// ValidateProjectPolicyBindingUpdate tests if required fields in the policy are set during
// an update.
func ValidateProjectPolicyBindingUpdate(ctx context.Context, new *auth.CustomPolicyBinding, old *auth.CustomPolicyBinding, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&new.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateProjectPolicyBinding(ctx, new, authClient)...)

	fldSpecPath := field.NewPath("spec")
	if new.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), new.Spec.TenantID, "disallowed change the tenant"))
	}

	if new.Spec.PolicyID != old.Spec.PolicyID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("policyID"), new.Spec.PolicyID, "disallowed change the policyID"))
	}

	if new.Spec.RulePrefix != old.Spec.RulePrefix {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("rulePrefix"), "disallowed change the rulePrefix"))
	}
	if new.Spec.RulePrefix == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("rulePrefix"), "must specify rulePrefix"))
	}

	if len(new.Spec.Resources) == 0 {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("resources"), "must specify resources"))
	}

	return allErrs
}
