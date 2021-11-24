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
	"net/http"
	"net/url"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"

	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	v1platform "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/util/addon"
	"tkestack.io/tke/pkg/util/log"
)

// ClusterNameToClient mapping cluster to kubernetes client
// clusterName => kubernetes.Interface
var ClusterNameToClient sync.Map

// GetClusterClient get kubernetes client via cluster name
func GetClusterClient(ctx context.Context, clusterName string, platformClient platformversionedclient.PlatformV1Interface) (kubernetes.Interface, error) {
	// First check from cache
	if item, ok := ClusterNameToClient.Load(clusterName); ok {
		// Check if is available
		kubeClient := item.(kubernetes.Interface)
		_, err := kubeClient.CoreV1().Services(metav1.NamespaceSystem).List(ctx, metav1.ListOptions{})
		if err == nil {
			return kubeClient, nil
		}
		ClusterNameToClient.Delete(clusterName)
	}

	kubeClient, err := addon.BuildExternalClientSetWithName(ctx, platformClient, clusterName)
	if err != nil {
		return nil, err
	}

	ClusterNameToClient.Store(clusterName, kubeClient)

	return kubeClient, nil
}

//TODO: use api && controller instead of proxy
func APIServerLocationByCluster(ctx context.Context, clusterName string, platformClient platformversionedclient.PlatformV1Interface) (*url.URL, http.RoundTripper, string, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, nil, "", errors.NewBadRequest("unable to get request info from context")
	}
	cluster, err := platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})
	if err != nil {
		log.Errorf("unable to get cluster %v", err)
		return nil, nil, "", err
	}
	if cluster.Status.Phase != v1platform.ClusterRunning {
		return nil, nil, "", errors.NewServiceUnavailable(fmt.Sprintf("cluster %s status is abnormal", cluster.ObjectMeta.Name))
	}
	credential, err := addon.GetClusterCredentialV1(ctx, platformClient, cluster)
	if err != nil {
		log.Errorf("unable to get credential %v", err)
		return nil, nil, "", err
	}

	restConfig := credential.RESTConfig(cluster)
	transport, err := restclient.TransportFor(restConfig)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}

	token := ""
	if credential.Token != nil {
		token = *credential.Token
	}
	return &url.URL{
		Scheme: "https",
		Host:   restConfig.Host,
		Path:   requestInfo.Path,
	}, transport, token, nil
}

//use cache to optimize this function
func GetClusterPodIP(ctx context.Context, clusterName, namespace, podName string, platformClient platformversionedclient.PlatformV1Interface) (string, error) {
	client, err := GetClusterClient(ctx, clusterName, platformClient)
	if err != nil {
		log.Errorf("unable to get cluster client %v", err)
		return "", err
	}
	pod, err := client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		log.Errorf("unable to get pod in cluster %v err=%v", clusterName, err)
		return "", err
	}
	return pod.Status.HostIP, nil
}
