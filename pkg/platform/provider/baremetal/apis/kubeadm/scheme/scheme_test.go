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

package scheme

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"tkestack.io/tke/pkg/platform/provider/baremetal/apis/kubeadm/v1beta2"
)

func Test_Decode(t *testing.T) {
	clusterConfig := &v1beta2.ClusterConfiguration{
		KubernetesVersion: "1.18.2",
		ClusterName:       "test",
	}
	const mediaType = runtime.ContentTypeYAML
	info, ok := runtime.SerializerInfoForMediaType(Codecs.SupportedMediaTypes(), mediaType)
	assert.True(t, ok)
	encoder := Codecs.EncoderForVersion(info.Serializer, v1beta2.SchemeGroupVersion)
	data, err := runtime.Encode(encoder, clusterConfig)
	assert.Nil(t, err)
	fmt.Println(string(data))
}
