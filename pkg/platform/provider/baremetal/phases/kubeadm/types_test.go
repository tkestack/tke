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

package kubeadm

import (
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
)

func TestConfig_Marshal(t *testing.T) {
	c := &Config{
		InitConfiguration: &v1beta2.InitConfiguration{
			TypeMeta:       metav1.TypeMeta{},
			CertificateKey: "a",
		},
	}
	data, err := c.Marshal()
	fmt.Println(string(data), err)
}
