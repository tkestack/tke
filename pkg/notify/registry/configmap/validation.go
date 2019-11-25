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

package configmap

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/notify"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apimachineryvalidation.ValidateNamespaceName

// ValidateConfigMap tests if required fields in the cluster are set.
func ValidateConfigMap(configmap *notify.ConfigMap) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&configmap.ObjectMeta, false, ValidateName, field.NewPath("metadata"))

	return allErrs
}

// ValidateConfigMapUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateConfigMapUpdate(configmap *notify.ConfigMap, old *notify.ConfigMap) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&configmap.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateConfigMap(configmap)...)

	return allErrs
}
