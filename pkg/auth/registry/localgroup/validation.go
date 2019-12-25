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

package localgroup

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/util/validation"
)

// ValidateGroupName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateGroupName = apiMachineryValidation.NameIsDNSLabel

// ValidateLocalGroup tests if required fields in the group are set.
func ValidateLocalGroup(group *auth.LocalGroup, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&group.ObjectMeta, false, ValidateGroupName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")

	if err := validation.IsDisplayName(group.Spec.DisplayName); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("displayName"), group.Spec.DisplayName, err.Error()))
	}

	fldUserPath := field.NewPath("status", "users")
	for i, subj := range group.Status.Users {
		if subj.ID == "" {
			allErrs = append(allErrs, field.Required(fldUserPath, "must specify id"))
			continue
		}

		if subj.Name == "" {
			val, err := authClient.LocalIdentities().Get(subj.ID, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					allErrs = append(allErrs, field.NotFound(fldUserPath, subj.ID))
				} else {
					allErrs = append(allErrs, field.InternalError(fldUserPath, err))
				}
			} else {
				if val.Spec.TenantID != group.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldUserPath, subj.ID, "must in the same tenant with the group"))
				} else {
					group.Status.Users[i].Name = val.Spec.Username
				}
			}
		}
	}

	return allErrs
}

// ValidateLocalGroupUpdate tests if required fields in the group are set during
// an update.
func ValidateLocalGroupUpdate(group *auth.LocalGroup, old *auth.LocalGroup, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&group.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateLocalGroup(group, authClient)...)

	fldSpecPath := field.NewPath("spec")
	if group.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), group.Spec.TenantID, "disallowed change the tenant"))
	}

	return allErrs
}
