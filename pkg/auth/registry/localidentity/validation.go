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
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/validation"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

var (
	reservedNames = map[string]bool{"system": true}
)

// ValidateLocalIdentity tests if required fields in the identity are set.
func ValidateLocalIdentity(authClient authinternalclient.AuthInterface, localIdentity *auth.LocalIdentity, updateCheck bool) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&localIdentity.ObjectMeta, false, apiMachineryValidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if localIdentity.Spec.Username == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("username"), "must specify username"))
	} else {
		if err := validation.IsDNS1123Name(localIdentity.Spec.Username); err != nil {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("username"), localIdentity.Spec.Username, err.Error()))
		}

		if _, ok := reservedNames[localIdentity.Spec.Username]; ok {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("username"), localIdentity.Spec.Username, `username is reserved`))
		}

		if !updateCheck {
			if exists, err := localIdentityExists(authClient, localIdentity.Spec.TenantID, localIdentity.Spec.Username); err != nil {
				allErrs = append(allErrs, field.InternalError(fldSpecPath.Child("username"), err))
			} else if exists {
				allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("username"), localIdentity.Spec.Username,
					fmt.Sprintf("%s %q already exists", auth.Resource("localidentities").String(), localIdentity.Spec.Username)))
			}
		}
	}

	if !updateCheck {
		if localIdentity.Spec.HashedPassword == "" {
			allErrs = append(allErrs, field.Required(fldSpecPath.Child("hashedPassword"), "must specify hashedPassword"))
		} else if bcrypted, err := util.BcryptPassword(localIdentity.Spec.HashedPassword); err != nil {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("hashedPassword"), localIdentity.Spec.HashedPassword, err.Error()))
		} else {
			localIdentity.Spec.HashedPassword = bcrypted
		}
	} else {
		if localIdentity.Spec.HashedPassword != "" {
			if bcrypted, err := util.BcryptPassword(localIdentity.Spec.HashedPassword); err != nil {
				allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("hashedPassword"), localIdentity.Spec.HashedPassword, err.Error()))
			} else {
				localIdentity.Spec.HashedPassword = bcrypted
			}
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
func ValidateLocalIdentityUpdate(authClient authinternalclient.AuthInterface, localIdentity *auth.LocalIdentity, oldLocalIdentity *auth.LocalIdentity) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateLocalIdentity(authClient, localIdentity, true)...)
	allErrs = append(allErrs, apiMachineryValidation.ValidateObjectMetaUpdate(&localIdentity.ObjectMeta, &oldLocalIdentity.ObjectMeta, field.NewPath("metadata"))...)

	fldSpecPath := field.NewPath("spec")
	if localIdentity.Spec.TenantID != oldLocalIdentity.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("tenantID"), localIdentity.Spec.TenantID, "disallowed change the tenant"))
	}

	if localIdentity.Spec.Username != oldLocalIdentity.Spec.Username {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("username"), localIdentity.Spec.Username, "disallowed change the username"))
	}

	if localIdentity.Spec.HashedPassword == "" {
		localIdentity.Spec.HashedPassword = oldLocalIdentity.Spec.HashedPassword
	}

	return allErrs
}

// ValidateLocalIdentityPasswordUpdate tests if required fields in the passwordReq are set
// during an update.
func ValidateLocalIdentityPasswordUpdate(localIdentity *auth.LocalIdentity, passwordReq *auth.PasswordReq) error {
	err := util.VerifyDecodedPassword(passwordReq.OriginalPassword, localIdentity.Spec.HashedPassword)
	if err != nil {
		log.Error("Invalid original password", log.String("original password", passwordReq.OriginalPassword), log.Err(err))
		return fmt.Errorf("verify original password failed: %v", err)
	}

	log.Info("local", log.Any("password", localIdentity.Spec.HashedPassword))
	if passwordReq.HashedPassword == "" {
		return fmt.Errorf("must specify hashedPassword")
	}

	if bcrypted, err := util.BcryptPassword(passwordReq.HashedPassword); err != nil {
		return fmt.Errorf("bcrypt password failed: %v", err)
	} else {
		localIdentity.Spec.HashedPassword = bcrypted
	}

	return nil
}

func localIdentityExists(authClient authinternalclient.AuthInterface, tenantID, username string) (bool, error) {
	tenantUserSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.tenantID", tenantID),
		fields.OneTermEqualSelector("spec.username", username))

	localIdentityList, err := authClient.LocalIdentities().List(v1.ListOptions{FieldSelector: tenantUserSelector.String()})
	if err != nil {
		return false, err
	}

	if len(localIdentityList.Items) > 0 {
		return true, nil
	}

	return false, nil
}
