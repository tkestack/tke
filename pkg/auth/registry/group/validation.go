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

package group

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

// ValidateGroupName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateGroupName = apiMachineryValidation.NameIsDNSLabel

// ValidateGroup tests if required fields in the group are set.
func ValidateGroup(group *auth.Group, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&group.ObjectMeta, false, ValidateGroupName, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if err := validation.IsDisplayName(group.Spec.DisplayName); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("displayName"), group.Spec.DisplayName, err.Error()))
	}

	fldStatPath := field.NewPath("status")
	for i, subj := range group.Status.Users {
		if subj.ID == "" && subj.Name == "" {
			allErrs = append(allErrs, field.Required(fldStatPath.Child("users"), "must specify id or name"))
			continue
		}

		// if specify id, ensure name
		if subj.ID != "" {
			val, err := authClient.LocalIdentities().Get(subj.ID, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					allErrs = append(allErrs, field.NotFound(fldStatPath.Child("users"), subj.ID))
				} else {
					allErrs = append(allErrs, field.InternalError(fldStatPath.Child("users"), err))
				}
			} else {
				if val.Spec.TenantID != group.Spec.TenantID {
					allErrs = append(allErrs, field.Invalid(fldStatPath.Child("users"), subj.ID, "must in the same tenant with the group"))
				} else {
					group.Status.Users[i].Name = val.Spec.Username
				}
			}
		} else {
			localIdentity, err := util.GetLocalIdentity(authClient, group.Spec.TenantID, subj.Name)
			if err != nil && apierrors.IsNotFound(err) {
				continue
			}
			if err != nil {
				allErrs = append(allErrs, field.InternalError(fldStatPath.Child("subjects"), err))
			} else {
				group.Status.Users[i].ID = localIdentity.Name
			}
		}
	}

	return allErrs
}

// ValidateGroupUpdate tests if required fields in the group are set during
// an update.
func ValidateGroupUpdate(group *auth.Group, old *auth.Group, authClient authinternalclient.AuthInterface) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&group.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateGroup(group, authClient)...)

	fldSpecPath := field.NewPath("spec")
	if group.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), group.Spec.TenantID, "disallowed change the tenant"))
	}

	return allErrs
}
