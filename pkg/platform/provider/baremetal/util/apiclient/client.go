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

package apiclient

import (
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	tkeclientset "tkestack.io/tke/api/client/clientset/versioned"
)

// GetClientset return clientset
func GetClientset(masterEndpoint string, token string, caCert []byte) (*kubernetes.Clientset, error) {
	restConfig := &rest.Config{
		Host:        masterEndpoint,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caCert,
		},
		Timeout: 5 * time.Second,
	}

	return kubernetes.NewForConfig(restConfig)
}

// GetPlatformClientset return clientset
func GetPlatformClientset(masterEndpoint string, token string, caCert []byte) (tkeclientset.Interface, error) {
	restConfig := &rest.Config{
		Host:        masterEndpoint,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caCert,
		},
		Timeout: 5 * time.Second,
	}

	return tkeclientset.NewForConfig(restConfig)
}
