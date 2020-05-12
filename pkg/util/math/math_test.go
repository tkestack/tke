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

import "testing"

func TestRange(t *testing.T) {
	type args struct {
		a []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"nil",
			args{
				nil,
			},
			0,
		},
		{
			"size=1",
			args{
				[]float64{0},
			},
			0,
		},
		{
			"size=2",
			args{
				[]float64{0, 1},
			},
			1,
		},
		{
			"size=3 sort",
			args{
				[]float64{0, 1, 2},
			},
			2,
		},
		{
			"size=3 random",
			args{
				[]float64{2, 0, 1},
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Range(tt.args.a); got != tt.want {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}
