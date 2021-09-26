/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package cluster

import (
	"testing"
)

func TestProvider_coreDNSNeedUpgrade(t *testing.T) {
	type args struct {
		k8sVersion string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Test k8s 1.18",
			args{
				k8sVersion: "1.18.4",
			},
			false,
		},
		{
			"Test k8s 1.19",
			args{
				k8sVersion: "1.19.7",
			},
			true,
		},
		{
			"Test k8s 1.20",
			args{
				k8sVersion: "1.20.4-tke.1",
			},
			true,
		},
		{
			"Test k8s 1.21",
			args{
				k8sVersion: "1.21.4-tke.1",
			},
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Provider{}
			if got := p.coreDNSNeedUpgrade(tt.args.k8sVersion); got != tt.want {
				t.Errorf("Provider.coreDNSNeedUpgrade() = %v, want %v", got, tt.want)
			}
		})
	}
}
