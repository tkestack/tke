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

package validation

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
	utilvalidation "tkestack.io/tke/pkg/util/validation"
)

// ValidateCluster validates a given Cluster.
func ValidateCluster(obj *platform.Cluster) field.ErrorList {
	allErrs := ValidatClusterAddresses(obj.Status.Addresses, field.NewPath("status", "addresses"))

	return allErrs
}

// ValidatClusterAddresses validates a given ClusterAddresses.
func ValidatClusterAddresses(addresses []platform.ClusterAddress, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(addresses) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "addresses"), "must specify at least one obj access address"))
	} else {
		for i, address := range addresses {
			fldPath := fldPath.Index(i)
			allErrs = utilvalidation.ValidateEnum(address.Type, fldPath.Child("type"), []interface{}{
				platform.AddressAdvertise,
				platform.AddressReal,
			})
			if address.Host == "" {
				allErrs = append(allErrs, field.Required(fldPath.Child("host"), "must specify host"))
			}
			for _, msg := range validation.IsValidPortNum(int(address.Port)) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("port"), address.Port, msg))
			}
			if address.Path != "" && !strings.HasPrefix(address.Path, "/") {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("path"), address.Path, "must start by `/`"))
			}

			url := fmt.Sprintf("https://%s:%d", address.Host, address.Port)
			if address.Path != "" {
				url = fmt.Sprintf("%s%s", url, address.Path)
			}
			err := utilvalidation.IsValiadURL(url, 5*time.Second)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath, address, err.Error()))
			}
		}
	}
	return allErrs
}
