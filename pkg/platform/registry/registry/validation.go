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

package registry

import (
	"fmt"
	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"regexp"
	"strings"
	"tkestack.io/tke/api/platform"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apiMachineryValidation.NameIsDNSSubdomain

// ValidateRegistryConfig tests if required fields in the cluster are set.
func ValidateRegistryConfig(registry *platform.Registry) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&registry.ObjectMeta, false, ValidateName, field.NewPath("metadata"))
	if registry.Spec.DisplayName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "registry", "displayName"), "displayName should not be empty"))
	}

	if registry.Spec.URL == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "registry", "url"), "url should not be empty"))
	}

	for _, msg := range validateRegistryURL(registry.Spec.URL, true) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "registry", "url"), "url is invalid", msg))
	}

	return allErrs
}

// ValidateRegistryConfigUpdate tests if required fields in the namespace set are
// set during an update.
func ValidateRegistryConfigUpdate(registryConfig *platform.Registry, old *platform.Registry) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&registryConfig.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateRegistryConfig(registryConfig)...)

	return allErrs
}

// validateRegistryURL verify the url match a domain
// allow domain like foo.com or bar.com:8080
func validateRegistryURL(url string, prefix bool) []string {
	var DNS1123SubdomainMaxLength = 253
	var validateRegistryReg = `^(https?:\/\/)?(([a-zA-Z0-9_-])+(\.)?)*(:\d+)?(\/((\.)?(\?)?=?&?[a-zA-Z0-9_-](\?)?)*)*$`
	var validateRegistryURLRegexp = regexp.MustCompile(validateRegistryReg)

	if prefix {
		url = maskTrailingDash(url)
	}

	var errs []string
	if len(url) > DNS1123SubdomainMaxLength {
		errs = append(errs, fmt.Sprintf("must be no more than %d characters", len(url)))
	}
	if !validateRegistryURLRegexp.MatchString(url) {
		errs = append(errs, validation.RegexError("a valid registry url should be a domain like foo.com or bar.com:8080", validateRegistryReg, "example.com"))
	}
	return errs
}

// maskTrailingDash replaces the final character of a string with a subdomain safe
// value if is a dash.
func maskTrailingDash(name string) string {
	if strings.HasSuffix(name, "-") {
		return name[:len(name)-2] + "a"
	}
	return name
}
