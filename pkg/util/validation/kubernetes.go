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
	"reflect"

	"github.com/thoas/go-funk"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateEnum validates a given enum.
// nil or nil pointer is valid.
// zero value is invalid.
func ValidateEnum(value interface{}, fldPath *field.Path, values interface{}) field.ErrorList {
	allErrs := field.ErrorList{}

	if value == nil {
		return allErrs
	}

	validValues := funk.Map(values, func(i interface{}) string {
		return fmt.Sprintf("%v", i)
	}).([]string)

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return allErrs
		}
	} else {
		if v.IsZero() {
			allErrs = append(allErrs, field.Required(fldPath, fmt.Sprintf("valid values: %v", validValues)))
			return allErrs
		}
	}

	if !funk.Contains(validValues, value) {
		allErrs = append(allErrs, field.NotSupported(fldPath, value, validValues))
	}

	return allErrs
}
