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

package projectrole

import (
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

// ValidateProjectPolicy tests if required fields in the policy are set.
func ValidateProjectPolicy(binding *auth.ProjectPolicyBinding, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&binding.ObjectMeta, false, ValidateBindingName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")

	if binding.Spec.PolicyID == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("policyID"), "must specify policyID"))
	} else {
		pol, err := authClient.Policies().Get(binding.Spec.PolicyID, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				allErrs = append(allErrs, field.NotFound(fldSpecPath.Child("policyID"), binding.Spec.PolicyID))
			} else {
				allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("policyID"), err))
			}
		}

		if pol.Spec.Scope != auth.PolicyProject {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("scope"), pol.Spec.Scope, "policy must be project scope"))
		}
	}

	// TODO maybe need to check projectID exists?
	if binding.Spec.ProjectID == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("projectID"), "must specify projectID"))
	}

	var validUsers []auth.Subject
	fldUserPath := field.NewPath("spec", "users")
	for i, subj := range binding.Spec.Users {
		if subj.ID == "" {
			allErrs = append(allErrs, field.Required(fldUserPath, "must specify subject id"))
			continue
		}

		if subj.Name == "" {
			val, err := authClient.Users().Get(util.CombineTenantAndName(binding.Spec.TenantID, subj.ID), metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Warn("user of the policy is not found, will removed it", log.String("policy", binding.Name), log.String("user", subj.Name))
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
		} else {
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
			val, err := authClient.Groups().Get(util.CombineTenantAndName(binding.Spec.TenantID, subj.ID), metav1.GetOptions{})
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

// ValidateProjectPolicyUpdate tests if required fields in the policy are set during
// an update.
func ValidateProjectPolicyUpdate(new *auth.ProjectPolicyBinding, old *auth.ProjectPolicyBinding, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&new.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateProjectPolicy(new, authClient)...)

	fldSpecPath := field.NewPath("spec")
	if new.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), new.Spec.TenantID, "disallowed change the tenant"))
	}

	if new.Spec.ProjectID != old.Spec.ProjectID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("projectID"), new.Spec.ProjectID, "disallowed change the projectID"))
	}

	if new.Spec.PolicyID != old.Spec.PolicyID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("policyID"), new.Spec.PolicyID, "disallowed change the policyID"))
	}

	return allErrs
}
