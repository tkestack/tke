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

	"tkestack.io/tke/api/application"
)

// AppSlice chartgroup array
type AppSlice []application.App

// Len return length
func (o AppSlice) Len() int { return len(o) }

// Swap will swap data by index
func (o AppSlice) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

// AppsByName sort chartgroups by chartgroup name
type AppsByName struct {
	AppSlice
	Desc bool
}

// Less 根据target升序排序
func (o AppsByName) Less(i, j int) bool {
	if o.Desc {
		return strings.Compare(o.AppSlice[i].Spec.Name, o.AppSlice[j].Spec.Name) > 0
	}
	return strings.Compare(o.AppSlice[i].Spec.Name, o.AppSlice[j].Spec.Name) <= 0
}
