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

package sort

import (
	"strings"

	"tkestack.io/tke/api/registry"
)

// ChartSlice chart array
type ChartSlice []registry.Chart

// Len return length
func (o ChartSlice) Len() int { return len(o) }

// Swap will swap data by index
func (o ChartSlice) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

// ChartsByName sort charts by chartgroup name and name
type ChartsByName struct {
	ChartSlice
	Desc bool
}

// Less 根据target升序排序
func (o ChartsByName) Less(i, j int) bool {
	c1 := strings.Compare(o.ChartSlice[i].Spec.ChartGroupName, o.ChartSlice[j].Spec.ChartGroupName)
	if o.Desc {
		if c1 == 0 {
			return strings.Compare(o.ChartSlice[i].Spec.Name, o.ChartSlice[j].Spec.Name) > 0
		}
		return c1 > 0
	}
	if c1 == 0 {
		return strings.Compare(o.ChartSlice[i].Spec.Name, o.ChartSlice[j].Spec.Name) <= 0
	}
	return c1 <= 0
}
