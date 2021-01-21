/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
 *
 */

package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformutil "tkestack.io/tke/pkg/platform/util"
)

// TKEClusterProvider get kubeconfig from tkestack platform api
type TKEClusterProvider struct {
	platformClient platformversionedclient.PlatformV1Interface
}

// NewTKEClusterProvider
func NewTKEClusterProvider(platformClient platformversionedclient.PlatformV1Interface) *TKEClusterProvider {
	return &TKEClusterProvider{
		platformClient: platformClient,
	}
}

func (t *TKEClusterProvider) RestConfig(clusterName string) (*rest.Config, error) {
	ctx := context.TODO()
	cls, err := t.platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cred, err := t.platformClient.
		ClusterCredentials().
		Get(ctx, cls.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return platformutil.GetExternalRestConfig(cls, cred)
}

func (t *TKEClusterProvider) Client(clusterName string, scheme *runtime.Scheme) (ctrlclient.Client, error) {
	config, err := t.RestConfig(clusterName)
	if err != nil {
		return nil, err
	}

	return ctrlclient.New(config, ctrlclient.Options{Scheme: scheme})
}
