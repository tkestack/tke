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

package util

import (
	"strings"
	businessv1 "tkestack.io/tke/api/business/v1"
)

func ResourceAdd(a businessv1.ResourceList, b businessv1.ResourceList) businessv1.ResourceList {
	result := businessv1.ResourceList{}
	for key, value := range a {
		quantity := value.DeepCopy()
		if other, found := b[key]; found {
			quantity.Add(other)
		}
		result[key] = quantity
	}
	for key, value := range b {
		if _, found := result[key]; !found {
			quantity := value.DeepCopy()
			result[key] = quantity
		}
	}
	return result
}

func ResourceNameTranslate(rName string) string {
	var res string
	if strings.Contains(rName, "limits.") {
		res = strings.TrimPrefix(rName, "limits.")
	} else if strings.Contains(rName, "requests.") {
		res = strings.TrimPrefix(rName, "requests.")
	} else {
		res = rName
	}
	return strings.ReplaceAll(res, ".", "_")
}
