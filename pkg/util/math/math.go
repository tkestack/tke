/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package math

// Max is to get the maximum number of the float64 numbers.
func Max(a []float64) (*int, *float64) {
	if a == nil {
		return nil, nil
	}
	max := a[0]
	index := 0
	for i, value := range a {
		if value > max {
			index = i
			max = value
		}
	}

	return &index, &max
}

// Min is to get the minimum number of the float64 numbers.
func Min(a []float64) (*int, *float64) {
	if a == nil {
		return nil, nil
	}
	min := a[0]
	index := 0
	for i, value := range a {
		if value < min {
			index = i
			min = value
		}
	}

	return &index, &min
}

// Range is the difference between maximum and minimum.
func Range(a []float64) float64 {
	if a == nil {
		return 0
	}
	min := a[0]
	max := a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	return max - min
}
