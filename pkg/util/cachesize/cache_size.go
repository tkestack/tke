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

package cachesize

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NewHeuristicWatchCacheSizes returns a map of suggested watch cache sizes
// based on total memory.
func NewHeuristicWatchCacheSizes(expectedRAMCapacityMB int) map[schema.GroupResource]int {
	size := expectedRAMCapacityMB / 10

	// We should specify cache size for a given resource only if it
	// is supposed to have non-default value.
	watchCacheSizes := make(map[schema.GroupResource]int)
	watchCacheSizes[schema.GroupResource{Resource: "cluster", Group: "platform.tkestack.io"}] = maxInt(5*size, 1000)
	watchCacheSizes[schema.GroupResource{Resource: "machine"}] = maxInt(10*size, 1000)
	watchCacheSizes[schema.GroupResource{Resource: "project"}] = maxInt(5*size, 1000)
	return watchCacheSizes
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
