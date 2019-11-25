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
	"tkestack.io/tke/api/business/v1"
)

// SubClusterHardFromUsed is used to sub Hard from Used
func SubClusterHardFromUsed(used *v1.ClusterUsed, needed v1.ClusterHard) {
	for clusterName, clusterNeeded := range needed {
		clusterUsed, clusterUsedExist := (*used)[clusterName]
		if !clusterUsedExist {
			continue
		}
		for k, v := range clusterNeeded.Hard {
			usedValue, ok := clusterUsed.Used[k]
			if ok {
				usedValue.Sub(v)
				clusterUsed.Used[k] = usedValue
			}
		}
		(*used)[clusterName] = clusterUsed
	}
}

// AddClusterHardToUsed is used to add Hard to Used
func AddClusterHardToUsed(used *v1.ClusterUsed, delta v1.ClusterHard) {
	for clusterName, clusterHard := range delta {
		clusterUsed, clusterUsedExist := (*used)[clusterName]
		if !clusterUsedExist {
			clusterUsed = v1.UsedQuantity{
				Used: make(v1.ResourceList),
			}
		}
		for k, v := range clusterHard.Hard {
			usedValue, ok := clusterUsed.Used[k]
			if ok {
				v.Add(usedValue)
			}
			clusterUsed.Used[k] = v
		}
		(*used)[clusterName] = clusterUsed
	}
}
