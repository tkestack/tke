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

package apisigningkey

import (
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
)

// ValidateSigningKey tests if required fields in the signing key are set.
func ValidateSigningKey(signingKey *auth.APISigningKey) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&signingKey.ObjectMeta, false, apiMachineryValidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	return allErrs
}

// ValidateSigningKeyUpdate tests if required fields in the session are set during
// an update.
func ValidateSigningKeyUpdate(signingKey *auth.APISigningKey, oldSigningKey *auth.APISigningKey) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateSigningKey(signingKey)...)
	allErrs = append(allErrs, apiMachineryValidation.ValidateObjectMetaUpdate(&signingKey.ObjectMeta, &oldSigningKey.ObjectMeta, field.NewPath("metadata"))...)

	return allErrs
}
