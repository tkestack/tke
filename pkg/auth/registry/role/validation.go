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
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/validation"
)

// ValidateRoleName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateRoleName = apiMachineryValidation.NameIsDNSLabel

// ValidateRole tests if required fields in the role are set.
func ValidateRole(ctx context.Context, role *auth.Role, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&role.ObjectMeta, false, ValidateRoleName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if err := validation.IsDisplayName(role.Spec.DisplayName); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("displayName"), role.Spec.DisplayName, err.Error()))
	}

	var validPolicies []string
	for _, pid := range role.Spec.Policies {
		pol, err := authClient.Policies().Get(ctx, pid, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				log.Warn("policy of the role is not found, will removed it", log.String("role", role.Name), log.String("policy", pol.Name))
			} else {
				allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("policies"), err))
			}
		} else {
			if pol.Spec.TenantID != role.Spec.TenantID {
				allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("policies"), pid, "related policy must be in the same tenant"))
			} else {
				validPolicies = append(validPolicies, pid)
			}
		}
	}

	var validUsers []auth.Subject
	fldUserPath := field.NewPath("status", "users")
	for i, subj := range role.Status.Users {
		if subj.ID == "" {
			allErrs = append(allErrs, field.Required(fldUserPath, "must specify id"))
			continue
		}

		if subj.Name == "" {
			val, err := authClient.Users().Get(ctx, util.CombineTenantAndName(role.Spec.TenantID, subj.ID), metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Warn("user of the role is not found, will removed it", log.String("role", role.Name), log.String("user", subj.Name))
				} else {
					allErrs = append(allErrs, field.InternalError(fldUserPath, err))
				}
			} else {
				if val.Spec.TenantID != role.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldUserPath, subj.ID, "must in the same tenant with the role"))
				} else {
					role.Status.Users[i].Name = val.Spec.Name
					validUsers = append(validUsers, role.Status.Users[i])
				}
			}
		} else {
			validUsers = append(validUsers, role.Status.Users[i])
		}
	}

	var validGroups []auth.Subject
	fldGroupPath := field.NewPath("status", "groups")
	for i, subj := range role.Status.Groups {
		if subj.ID == "" {
			allErrs = append(allErrs, field.Required(fldGroupPath, "must specify id or name"))
			continue
		}

		if subj.Name == "" {
			val, err := authClient.Groups().Get(ctx, util.CombineTenantAndName(role.Spec.TenantID, subj.ID), metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					log.Warn("group of the role is not found, will removed it", log.String("role", role.Name), log.String("group", subj.Name))
				} else {
					allErrs = append(allErrs, field.InternalError(fldGroupPath, err))
				}
			} else {
				if val.Spec.TenantID != role.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldGroupPath, subj.ID, "must in the same tenant with the role"))
				} else {
					role.Status.Groups[i].Name = val.Spec.DisplayName
					validGroups = append(validGroups, role.Status.Groups[i])
				}
			}
		} else {
			validGroups = append(validGroups, role.Status.Groups[i])
		}
	}

	if len(allErrs) == 0 {
		role.Spec.Policies = validPolicies
		role.Status.Users = validUsers
		role.Status.Groups = validGroups
	}

	return allErrs
}

// ValidateRoleUpdate tests if required fields in the role are set during
// an update.
func ValidateRoleUpdate(ctx context.Context, role *auth.Role, old *auth.Role, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&role.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateRole(ctx, role, authClient)...)

	fldSpecPath := field.NewPath("spec")
	if role.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), role.Spec.TenantID, "disallowed change the tenant"))
	}

	return allErrs
}
