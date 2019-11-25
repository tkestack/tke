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

// InInt32Slice checks if all elements in the array are integers.
func InInt32Slice(slice []int32, num int32) bool {
	if len(slice) == 0 {
		return false
	}
	exist := false
	for _, n := range slice {
		if n == num {
			exist = true
			break
		}
	}
	return exist
}

// InStringSlice checks if all elements in the array are strings.
func InStringSlice(slice []string, str string) bool {
	if len(slice) == 0 {
		return false
	}
	exist := false
	for _, s := range slice {
		if s == str {
			exist = true
			break
		}
	}
	return exist
}

// DiffStringSlice returns the difference between two given string arrays,
// including deleted and added elements.
func DiffStringSlice(origins []string, updated []string) (added []string, removed []string) {
	if origins == nil {
		origins = []string{}
	}
	if updated == nil {
		updated = []string{}
	}
	for _, origin := range origins {
		if !InStringSlice(updated, origin) {
			removed = append(removed, origin)
		}
	}
	for _, update := range updated {
		if !InStringSlice(origins, update) {
			added = append(added, update)
		}
	}
	return
}
