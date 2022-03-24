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

package role

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/authz"
)

var ValidateRoleName = apimachineryvalidation.NameIsDNSLabel

func ValidateRole(role *authz.Role) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&role.ObjectMeta, true, ValidateRoleName, field.NewPath("metadata"))

	return allErrs
}

// ValidateRoleUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateRoleUpdate(role *authz.Role, old *authz.Role) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&role.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateRole(role)...)

	return allErrs
}
