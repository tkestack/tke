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
	"fmt"
	"sync"

	istioconfig "istio.io/client-go/pkg/apis/config/v1alpha2"
	istioscheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"tkestack.io/tke/pkg/util/log"
)

// Client facade interface
type Client interface {
	// RestConfig given cluster id to get Kubernetes cluster rest config
	RestConfig(cluster string) (*rest.Config, error)
	// Cluster given cluster id to get Kubernetes cluster client
	Cluster(cluster string) (ctrlclient.Client, error)
	// Istio given cluster id to get Istio client
	Istio(cluster string) (ctrlclient.Client, error)
}

// ClusterProvider retrieve cluster kubeconfig provider interface
type ClusterProvider interface {
	// RestConfig Kubernetes cluster rest config
	RestConfig(clusterName string) (*rest.Config, error)
	// Client Kubernetes cluster client with scheme
	Client(clusterName string, scheme *runtime.Scheme) (ctrlclient.Client, error)
}

var (
	ClusterNameToClient      sync.Map
	ClusterNameToIstioClient sync.Map
)

type client struct {
	clusterProvider ClusterProvider
}

// New new client
func New(provider ClusterProvider) Client {
	return &client{
		clusterProvider: provider,
	}
}

// RestConfig
func (c *client) RestConfig(clusterName string) (*rest.Config, error) {
	restConfig, err := c.clusterProvider.RestConfig(clusterName)
	return restConfig, err
}

// Cluster
func (c *client) Cluster(clusterName string) (ctrlclient.Client, error) {
	// First check from cache
	if item, ok := ClusterNameToClient.Load(clusterName); ok {
		// Check if is available
		clusterClient := item.(ctrlclient.Client)
		ns := &corev1.Namespace{}
		err := clusterClient.Get(context.TODO(), types.NamespacedName{Name: metav1.NamespaceSystem}, ns)
		if err == nil {
			return clusterClient, nil
		}
		ClusterNameToClient.Delete(clusterName)
	}

	clusterClient, err := c.clusterProvider.Client(clusterName, kubescheme.Scheme)
	if err != nil {
		log.Errorf("create kubernetes client failed: %v", err)
		return nil, fmt.Errorf("create kubernetes client failed")
	}
	ClusterNameToClient.Store(clusterName, clusterClient)

	return clusterClient, nil
}

// Istio
func (c *client) Istio(clusterName string) (ctrlclient.Client, error) {
	// First check from cache
	if item, ok := ClusterNameToIstioClient.Load(clusterName); ok {
		// Check if is available
		istioClient := item.(ctrlclient.Client)
		err := istioClient.List(context.TODO(), &istioconfig.RuleList{})
		if err == nil {
			return istioClient, nil
		}
		log.Warnf("load istio client from cache: %v", err)
		ClusterNameToIstioClient.Delete(clusterName)
	}

	clusterClient, err := c.clusterProvider.Client(clusterName, istioscheme.Scheme)
	if err != nil {
		log.Errorf("create istio client failed: %v", err)
		return nil, fmt.Errorf("create istio client failed")
	}
	ClusterNameToIstioClient.Store(clusterName, clusterClient)

	return clusterClient, nil
}
