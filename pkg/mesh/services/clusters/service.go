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
 *
 */

package clusters

import (
	"context"
	"net/http"
	"net/url"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	restclient "k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	clusterclient "tkestack.io/tke/pkg/mesh/external/kubernetes"
	"tkestack.io/tke/pkg/mesh/services"
	"tkestack.io/tke/pkg/util/log"
)

type clusterService struct {
	platformClient platformversionedclient.PlatformV1Interface
	clients        clusterclient.Client
}

var _ services.ClusterService = &clusterService{}

func New(
	platformClient platformversionedclient.PlatformV1Interface,
	clients clusterclient.Client,
) services.ClusterService {

	return &clusterService{
		platformClient: platformClient,
		clients:        clients,
	}
}

func (c *clusterService) Get(ctx context.Context, clusterName string) (*platformv1.Cluster, error) {
	return c.platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})
}

func (c *clusterService) List(ctx context.Context) ([]platformv1.Cluster, error) {
	list, err := c.platformClient.Clusters().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	clusters := make([]platformv1.Cluster, 0)
	clusters = append(clusters, list.Items...)
	return clusters, nil
}

func (c *clusterService) ListNamespaces(ctx context.Context, clusterName string) ([]corev1.Namespace, error) {
	client, err := c.clients.Cluster(clusterName)
	if err != nil {
		return nil, err
	}
	ret := &corev1.NamespaceList{}
	err = client.List(ctx, ret)
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}

func (c *clusterService) ListAll(ctx context.Context, clusterName string) ([]corev1.Namespace, error) {
	client, err := c.clients.Cluster(clusterName)
	if err != nil {
		return nil, err
	}
	ret := &corev1.NamespaceList{}
	err = client.List(ctx, ret)
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}

// ListService list k8s service by namespace and label selector
// if namespace is empty, return all namespace service
// if selector is not nil, return services which matches label selector
func (c *clusterService) ListServices(
	ctx context.Context, clusterName,
	namespace string, selector labels.Selector,
) ([]corev1.Service, error) {

	client, err := c.clients.Cluster(clusterName)
	if err != nil {
		return nil, err
	}
	ret := &corev1.ServiceList{}
	err = client.List(ctx, ret, &ctrlclient.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}
	return ret.Items, nil
}

func (c *clusterService) Proxy(
	ctx context.Context, clusterName string,
) (transport http.RoundTripper, host string, err error) {

	restConfig, err := c.clients.RestConfig(clusterName)
	if err != nil {
		return nil, "", err
	}
	transport, err = restclient.TransportFor(restConfig)
	if err != nil {
		log.Errorf("create kubernetes proxy transport error: %v", err)
		return nil, "", err
	}
	u, err := url.Parse(restConfig.Host)
	if err != nil {
		return nil, "", err
	}
	host = u.Host
	return transport, host, err
}
