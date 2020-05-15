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

package apiclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
)

func TestCreateOrUpdateNamespace(t *testing.T) {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/Users/chenglong/.kube/config")
	assert.Nil(t, err)
	client := kubernetes.NewForConfigOrDie(cfg)
	CreateOrUpdateNamespace(context.Background(), client, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "a",
		},
	})
}
