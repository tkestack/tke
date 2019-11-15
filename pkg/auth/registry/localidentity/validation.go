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

package localidentity

import (
	"fmt"
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/util/validation"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

var (
	reservedNames = map[string]bool{"system": true}
)

// ValidateLocalIdentity tests if required fields in the identity are set.
func ValidateLocalIdentity(authClient *authinternalclient.AuthClient, localIdentity *auth.LocalIdentity, validateCheck bool) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&localIdentity.ObjectMeta, false, apiMachineryValidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if localIdentity.Spec.UserName == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("userName"), "must specify userName"))
	}

	if err := validation.IsDNS1123Name(localIdentity.Spec.UserName); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("userName"), localIdentity.Spec.UserName, err.Error()))
	}

	if _, ok := reservedNames[localIdentity.Spec.UserName]; ok {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("userName"), localIdentity.Spec.UserName, `userName is reserved`))
	}

	if !validateCheck {
		if exists, err := localIdentityExists(authClient, localIdentity.Spec.TenantID, localIdentity.Spec.UserName); err != nil {
			allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("userName"), err))
		} else if exists {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("userName"), localIdentity.Spec.UserName,
				fmt.Sprintf("%s %q already exists", auth.Resource("localidentities").String(), localIdentity.Spec.UserName)))
		}

		if len(localIdentity.Spec.HashedPassword) == 0 {
			allErrs = append(allErrs, field.Required(fldSpecPath.Child("hashedPassword"), "password is empty or invalid"))
		}
	}

	if err := validation.IsDisplayName(localIdentity.Spec.DisplayName); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("displayName"), localIdentity.Spec.DisplayName, err.Error()))
	}

	if localIdentity.Spec.Email != "" {
		if err := validation.IsEmail(localIdentity.Spec.Email); err != nil {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("email"), localIdentity.Spec.Email, err.Error()))
		}
	}

	if localIdentity.Spec.PhoneNumber != "" {
		if err := validation.IsPhoneNumber(localIdentity.Spec.PhoneNumber); err != nil {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("phoneNumber"), localIdentity.Spec.PhoneNumber, err.Error()))
		}
	}

	return allErrs
}

// ValidateLocalIdentityUpdate tests if required fields in the localIdentity are set
// during an update.
func ValidateLocalIdentityUpdate(authClient *authinternalclient.AuthClient, localIdentity *auth.LocalIdentity, oldLocalIdentity *auth.LocalIdentity) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateLocalIdentity(authClient, localIdentity, true)...)
	allErrs = append(allErrs, apiMachineryValidation.ValidateObjectMetaUpdate(&localIdentity.ObjectMeta, &oldLocalIdentity.ObjectMeta, field.NewPath("metadata"))...)

	fldSpecPath := field.NewPath("spec")
	if localIdentity.Spec.TenantID != oldLocalIdentity.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), localIdentity.Spec.TenantID, "disallowed change the tenant"))
	}

	if localIdentity.Spec.UserName != oldLocalIdentity.Spec.UserName {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("userName"), localIdentity.Spec.UserName, "disallowed change the userName"))
	}

	return allErrs
}

func localIdentityExists(authClient *authinternalclient.AuthClient, tenantID, userName string) (bool, error) {
	tenantUserSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.tenantID", tenantID),
		fields.OneTermEqualSelector("spec.userName", userName))

	localIdentityList, err := authClient.LocalIdentities().List(v1.ListOptions{FieldSelector: tenantUserSelector.String()})
	if err != nil {
		return false, err
	}

	if len(localIdentityList.Items) > 0 {
		return true, nil
	}

	return false, nil
}
