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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	helmconfig "tkestack.io/tke/pkg/application/helm/config"
)

// NewHelmClient return a new client used to run helm cmd
func NewHelmClient(ctx context.Context,
	platformClient platformversionedclient.PlatformV1Interface,
	clusterID string,
	namespace string) (*helmaction.Client, error) {
	cluster, err := platformClient.Clusters().Get(ctx, clusterID, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("can not found cluster by name %s", cluster))
	}
	fieldSelector := fields.OneTermEqualSelector("clusterName", clusterID).String()
	list, err := platformClient.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return nil, fmt.Errorf("get cluster's credential error: %w", err)
	} else if len(list.Items) == 0 {
		return nil, fmt.Errorf("get cluster's credential error, no cluster credential")
	}
	credential := list.Items[0]

	restConfig, err := credential.RESTConfig(cluster)
	if err != nil {
		return nil, fmt.Errorf("get cluster's externalRestConfig error: %w", err)
	}
	restClientGetter := &helmconfig.RESTClientGetter{RestConfig: restConfig}
	// we should set namespace here. If not, release will be installed in target namespace, but resources will not be installed in target namespace
	restClientGetter.Namespace = &namespace
	client := helmaction.NewClient("", restClientGetter)
	return client, nil
}

// NewHelmClientWithoutRESTClient return a new client used to run helm cmd
func NewHelmClientWithoutRESTClient() *helmaction.Client {
	client := helmaction.NewClient("", nil)
	return client
}
