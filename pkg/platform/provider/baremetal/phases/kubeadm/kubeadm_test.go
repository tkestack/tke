/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package kubeadm

import (
	"testing"
)

func Test_sameVersion(t *testing.T) {
	type args struct {
		ver1               string
		ver2               string
		ignorePatchVersion bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"same version",
			args{
				ver1: "1.20.4",
				ver2: "1.20.4",
			},
			true,
			false,
		},
		{
			"same version with v",
			args{
				ver1: "1.20.4",
				ver2: "v1.20.4",
			},
			true,
			false,
		},
		{
			"diffrent version",
			args{
				ver1: "1.20.4",
				ver2: "1.20.4-tke.1",
			},
			false,
			false,
		},
		{
			"diffrent minor version",
			args{
				ver1:               "1.19.4",
				ver2:               "1.20.4",
				ignorePatchVersion: true,
			},
			false,
			false,
		},
		{
			"ignore patch version",
			args{
				ver1:               "1.20.4",
				ver2:               "1.20.4-tke.1",
				ignorePatchVersion: true,
			},
			true,
			false,
		},
		{
			"same inc version with v",
			args{
				ver1:               "v1.20.4-tke.1",
				ver2:               "1.20.4-tke.1",
				ignorePatchVersion: false,
			},
			true,
			false,
		},
		{
			"empty",
			args{
				ver1:               "",
				ver2:               "1.20.4-tke.1",
				ignorePatchVersion: false,
			},
			false,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sameVersion(tt.args.ver1, tt.args.ver2, tt.args.ignorePatchVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("sameVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sameVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
