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

package mark

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/tke/pkg/util/apiclient"
)

const (
	Name      = "tke"
	Namespace = metav1.NamespaceSystem
)

func Get(ctx context.Context, clientset kubernetes.Interface) (*corev1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps(Namespace).Get(ctx, Name, metav1.GetOptions{})
}

func Create(ctx context.Context, clientset kubernetes.Interface) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: Namespace,
		},
	}

	return apiclient.CreateOrUpdateConfigMap(ctx, clientset, cm)
}
