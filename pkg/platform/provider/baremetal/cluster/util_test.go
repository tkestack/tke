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

package cluster

import "testing"

func TestGetServiceCIDRAndNodeCIDRMaskSize(t *testing.T) {
	type args struct {
		clusterCIDR          string
		maxClusterServiceNum int32
		maxNodePodNum        int32
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int32
		wantErr bool
	}{
		{
			name: "maxClusterServiceNum == 0",
			args: args{
				clusterCIDR:          "192.168.0.0/24",
				maxClusterServiceNum: 0,
			},
			wantErr: true,
		},
		{
			name: "maxNodePodNum == 0",
			args: args{
				clusterCIDR:          "192.168.0.0/24",
				maxClusterServiceNum: 0,
			},
			wantErr: true,
		},
		{
			name: "maxClusterServiceNum maxNodePodNum < clusterCIDR size",
			args: args{
				clusterCIDR:          "192.168.0.0/24",
				maxClusterServiceNum: 32,
				maxNodePodNum:        64,
			},
			want:    "192.168.0.224/27",
			want1:   26,
			wantErr: false,
		},
		{
			name: "maxClusterServiceNum == clusterCIDR size",
			args: args{
				clusterCIDR:          "192.168.0.0/24",
				maxClusterServiceNum: 256,
				maxNodePodNum:        64,
			},
			want:    "192.168.0.0/24",
			want1:   26,
			wantErr: false,
		},
		{
			name: "maxClusterServiceNum > clusterCIDR size",
			args: args{
				clusterCIDR:          "192.168.0.0/24",
				maxClusterServiceNum: 257,
				maxNodePodNum:        64,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetServiceCIDRAndNodeCIDRMaskSize(tt.args.clusterCIDR, tt.args.maxClusterServiceNum, tt.args.maxNodePodNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServiceCIDRAndNodeCIDRMaskSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetServiceCIDRAndNodeCIDRMaskSize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetServiceCIDRAndNodeCIDRMaskSize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
