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
	corev1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
)

// ConvertToCoreV1ResourceList convert the resource list in tke to kubernetes
// core resource list.
func ConvertToCoreV1ResourceList(resourceList map[string]apiresource.Quantity) corev1.ResourceList {
	capacity := make(corev1.ResourceList)
	if len(resourceList) > 0 {
		for k, v := range resourceList {
			capacity[corev1.ResourceName(k)] = v
		}
	}
	return capacity
}

// ConvertFromCoreV1ResourceList convert the kubernetes core resource list to
// common resource list.
func ConvertFromCoreV1ResourceList(resourceList corev1.ResourceList) map[string]apiresource.Quantity {
	capacity := make(map[string]apiresource.Quantity)
	if len(resourceList) > 0 {
		for k, v := range resourceList {
			capacity[string(k)] = v
		}
	}
	return capacity
}
