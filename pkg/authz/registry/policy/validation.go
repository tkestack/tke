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

package policy

import (
	"context"
	"fmt"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/authz"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
)

var ValidatePolicyName = apimachineryvalidation.NameIsDNSLabel

// ValidatePolicy tests if required fields in the cluster are set.
func ValidatePolicy(policy *authz.Policy, platformClient platformversionedclient.PlatformV1Interface) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&policy.ObjectMeta, true, ValidatePolicyName, field.NewPath("metadata"))
	for i, rule := range policy.Rules {
		fldPath := field.NewPath("rules").Index(i)
		if len(rule.Verbs) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("verbs"), "verbs must contain at least one value"))
		}
		if len(rule.NonResourceURLs) > 0 {
			if len(rule.APIGroups) > 0 || len(rule.Resources) > 0 || len(rule.ResourceNames) > 0 {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("nonResourceURLs"), rule.NonResourceURLs, "rules cannot apply to both regular resources and non-resource URLs"))
			}
			return allErrs
		}
		if len(rule.APIGroups) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("apiGroups"), "resource rules must supply at least one api group"))
		}
		if len(rule.Resources) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("resources"), "resource rules must supply at least one resource"))
		}
	}
	return allErrs
}

// ValidatePolicyUpdate tests if required fields in the namespace set are
// set during an update.
func ValidatePolicyUpdate(ctx context.Context, policy *authz.Policy, old *authz.Policy, platformClient platformversionedclient.PlatformV1Interface) field.ErrorList {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID == "" {
		tenantID = "default"
	}
	if tenantID != "default" && tenantID != policy.Namespace {
		return append(field.ErrorList{}, field.Required(field.NewPath("metadata", "namespace"), fmt.Sprintf("tenant '%s' can't update policy '%s/%s'", tenantID, policy.Namespace, policy.Name)))
	}
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&policy.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidatePolicy(policy, platformClient)...)
	return allErrs
}
