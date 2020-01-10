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

package identityprovider

import (
	"fmt"

	"github.com/dexidp/dex/server"
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/util/validation"
)

// ValidateIdentityProvider tests if required fields in the signing key are set.
func ValidateIdentityProvider(idp *auth.IdentityProvider) field.ErrorList {

	allErrs := apiMachineryValidation.ValidateObjectMeta(&idp.ObjectMeta, false, apiMachineryValidation.NameIsDNSSubdomain, field.NewPath("metadata"))

	fldSpecPath := field.NewPath("spec")
	if err := validation.IsDNS1123Name(idp.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("categoryName"), idp.Name, err.Error()))
	}

	if idp.Spec.Name == "" {
		allErrs = append(allErrs, field.Required(fldSpecPath.Child("name"), "must specify name"))
	} else {
		if err := validation.IsDisplayName(idp.Spec.Name); err != nil {
			allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("displayName"), idp.Spec.Name, err.Error()))
		}
	}

	if _, ok := server.ConnectorsConfig[idp.Spec.Type]; !ok && idp.Spec.Type != server.LocalConnector {
		allErrs = append(allErrs, field.Invalid(field.NewPath("type"), idp.Spec.Type, fmt.Sprintf("only support %v", suppportIDPTypes())))
	}

	return allErrs
}

// ValidateIdentityProviderUpdate tests if required fields in the session are set during
// an update.
func ValidateIdentityProviderUpdate(idp *auth.IdentityProvider, oldIdentityProvider *auth.IdentityProvider) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateIdentityProvider(idp)...)
	allErrs = append(allErrs, apiMachineryValidation.ValidateObjectMetaUpdate(&idp.ObjectMeta, &oldIdentityProvider.ObjectMeta, field.NewPath("metadata"))...)

	fldSpecPath := field.NewPath("metadata")

	if idp.Name != oldIdentityProvider.Name {
		allErrs = append(allErrs, field.Invalid(fldSpecPath.Child("name"), idp.Name, "disallowed change the name"))
	}

	return allErrs
}

func suppportIDPTypes() []string {
	var supportTypes []string
	for key := range server.ConnectorsConfig {
		supportTypes = append(supportTypes, key)
	}

	supportTypes = append(supportTypes, server.LocalConnector)

	return supportTypes
}
