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

package platform

import (
	"fmt"

	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/business"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	"tkestack.io/tke/cmd/tke-business-api/app/options"
)

// ValidatePlatformName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidatePlatformName = apimachineryvalidation.NameIsDNSLabel

// ValidatePlatform tests if required fields in the platform are set.
func ValidatePlatform(platform *business.Platform, businessClient *businessinternalclient.BusinessClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&platform.ObjectMeta, false, ValidatePlatformName, field.NewPath("metadata"))

	if platform.Spec.TenantID != "" && platform.Name != options.DefaultPlatform {
		platformList, err := businessClient.Platforms().List(metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.tenantID=%s", platform.Spec.TenantID),
		})
		if err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath("spec", "tenantID"), err))
		} else if len(platformList.Items) > 0 {
			allErrs = append(allErrs, field.Duplicate(field.NewPath("spec", "tenantID"), platform.Spec.TenantID))
		}
	}

	return allErrs
}

// ValidatePlatformUpdate tests if required fields in the platform are set during
// an update.
func ValidatePlatformUpdate(platform *business.Platform, old *business.Platform) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&platform.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))

	return allErrs
}
