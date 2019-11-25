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

import (
	"testing"
)
import corev1 "k8s.io/api/core/v1"

func TestParseTemplateTo(t *testing.T) {
	type args struct {
		strtmpl string
		obj     interface{}
		dst     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"a",
			args{
				`
apiVersion: v1
kind: ServiceAccount
metadata:
  name: flannel
  namespace: kube-system
`,
				nil,
				&corev1.ServiceAccount{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseTemplateTo(tt.args.strtmpl, tt.args.obj, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("ParseTemplateTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
