/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package validation

import (
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
)

// ValidateConfiguration validates `mc` and returns an error if it is invalid
func ValidateConfiguration(mc *meshconfig.MeshConfiguration) error {
	var allErrors []error

	componentField := field.NewPath("component")
	meshField := componentField.Child("meshManager")

	if mc.Components.MeshManager != nil {
		address := mc.Components.MeshManager.Address
		if len(address) == 0 {
			allErrors = append(allErrors, field.Required(meshField.Child("address"), ""))
		}
	}

	return utilerrors.NewAggregate(allErrors)
}
