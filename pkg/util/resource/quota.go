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

package resource

import (
	"fmt"

	apimachineryresource "k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/business"
	"tkestack.io/tke/pkg/util/log"
)

const (
	QuotaLimitErrorInfo  = "must set quota"
	AllocatableErrorInfo = "should NOT be more than remaining quantity"
	UpdateQuotaErrorInfo = "should NOT be less than used quantity"
)

// ValidateUpdateResource is used to verify that the adjusted quota when
// updating the resource quota cannot be less than the used one.
func ValidateUpdateResource(hard business.ResourceList, used business.ResourceList, fldPath *field.Path) field.ErrorList {
	if hard == nil {
		hard = business.ResourceList{}
	}
	if used == nil {
		used = business.ResourceList{}
	}

	allErrs := field.ErrorList{}
	for hardResourceName, hardResourceQuantity := range hard {
		usedResourceQuantity, usedResourceExist := used[hardResourceName]
		if !usedResourceExist {
			usedResourceQuantity = *apimachineryresource.NewQuantity(0, hardResourceQuantity.Format)
		}
		if hardResourceQuantity.Cmp(usedResourceQuantity) < 0 {
			log.Error("INSUFFICIENT", log.String("resourceName", hardResourceName), log.Int("cmp", hardResourceQuantity.Cmp(usedResourceQuantity)))
			allErrs = append(allErrs, field.Invalid(fldPath.Key(hardResourceName), hardResourceQuantity.String(),
				fmt.Sprintf("%s(%s)", UpdateQuotaErrorInfo, usedResourceQuantity.String())))
		}
	}
	return allErrs
}

// ValidateAllocatableResources is used to verify whether the resource quota to
// be allocated meets the total resource quota minus the condition of the used
// resource.
func ValidateAllocatableResources(wantAllocate business.ResourceList, hasAllocated business.ResourceList,
	capacity business.ResourceList, used business.ResourceList, fldPath *field.Path) field.ErrorList {
	if wantAllocate == nil {
		wantAllocate = business.ResourceList{}
	}
	if hasAllocated == nil {
		hasAllocated = business.ResourceList{}
	}
	if capacity == nil {
		capacity = business.ResourceList{}
	}
	if used == nil {
		used = business.ResourceList{}
	}

	allErrs := field.ErrorList{}

	for capacityResourceName := range capacity {
		_, wantAllocateExist := wantAllocate[capacityResourceName]
		if !wantAllocateExist {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(capacityResourceName), 0,
				fmt.Sprintf("%s '%s', for the parent project has set this kind of quota", QuotaLimitErrorInfo, capacityResourceName)))
		}
	}

	for wantResourceName, wantResourceQuantity := range wantAllocate {
		capacityResourceQuantity, capacityResourceExist := capacity[wantResourceName]
		if !capacityResourceExist {
			continue
		}
		usedResourceQuantity, usedResourceExist := used[wantResourceName]
		if !usedResourceExist {
			usedResourceQuantity = *apimachineryresource.NewQuantity(0, wantResourceQuantity.Format)
		}
		hasResourceQuantity, hasResourceExist := hasAllocated[wantResourceName]
		if hasResourceExist {
			usedResourceQuantity.Sub(hasResourceQuantity)
		}
		capacityResourceQuantity.Sub(usedResourceQuantity)
		if wantResourceQuantity.Cmp(capacityResourceQuantity) > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(wantResourceName), wantResourceQuantity.String(),
				fmt.Sprintf("%s(%s)", AllocatableErrorInfo, capacityResourceQuantity.String())))
		}
	}

	return allErrs
}
