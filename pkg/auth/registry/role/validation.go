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

package role

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

// ValidateRoleName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateRoleName = apiMachineryValidation.NameIsDNSLabel

// ValidateRole tests if required fields in the role are set.
func ValidateRole(role *auth.Role, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&role.ObjectMeta, false, ValidateRoleName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if err := validation.IsDisplayName(role.Spec.DisplayName); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("displayName"), role.Spec.DisplayName, err.Error()))
	}

	for _, pid := range role.Spec.Policies {
		pol, err := authClient.Policies().Get(pid, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				allErrs = append(allErrs, field.NotFound(fldSpecPath.Child("policies"), pid))
			} else {
				allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("policies"), err))
			}
		} else {
			if pol.Spec.TenantID != role.Spec.TenantID {
				allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("policies"), pid, "related policy must be in the same tenant"))
			}
		}
	}

	fldStatPath := field.NewPath("status")
	for i, subj := range role.Status.Subjects {
		if subj.ID == "" && subj.Name == "" {
			allErrs = append(allErrs, field.Required(fldStatPath.Child("subjects"), "must specify id or name"))
			continue
		}

		// if specify id, ensure name
		if subj.ID != "" {
			val, err := authClient.LocalIdentities().Get(subj.ID, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					allErrs = append(allErrs, field.NotFound(fldStatPath.Child("subjects"), subj.ID))
				} else {
					allErrs = append(allErrs, field.InternalError(fldStatPath.Child("subjects"), err))
				}
			} else {
				if val.Spec.TenantID != val.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldStatPath.Child("subjects"), subj.ID, "must in the same tenant with the role"))
				} else {
					role.Status.Subjects[i].Name = val.Spec.Username
				}
			}
		} else {
			localIdentity, err := util.GetLocalIdentity(authClient, role.Spec.TenantID, subj.Name)
			if err != nil && apierrors.IsNotFound(err) {
				continue
			}
			if err != nil {
				allErrs = append(allErrs, field.InternalError(fldStatPath.Child("subjects"), err))
			} else {
				role.Status.Subjects[i].ID = localIdentity.Name
			}
		}
	}

	return allErrs
}

// ValidateRoleUpdate tests if required fields in the role are set during
// an update.
func ValidateRoleUpdate(role *auth.Role, old *auth.Role, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&role.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateRole(role, authClient)...)

	fldSpecPath := field.NewPath("spec")
	if role.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), role.Spec.TenantID, "disallowed change the tenant"))
	}

	return allErrs
}
