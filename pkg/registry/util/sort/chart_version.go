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
	"sort"

	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/version"
)

// ByVersion implements sort.Interface for []ChartVersion based on
// the version field.
type ByVersion []registryv1.ChartVersion

func (a ByVersion) Len() int      { return len(a) }
func (a ByVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByVersion) Less(i, j int) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Failed compare chart version %s and %s", a[i].Version, a[j].Version)
		}
	}()
	return version.Compare(a[i].Version, a[j].Version) > 0
}

// ByChartVersion sort chart versions from high to low
func ByChartVersion(versions []registryv1.ChartVersion) []registryv1.ChartVersion {
	sort.Sort(ByVersion(versions))
	return versions
}
