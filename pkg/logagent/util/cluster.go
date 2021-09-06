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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes"

	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	v1platform "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/util"
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

	kubeClient, err := BuildExternalClientSetWithName(ctx, platformClient, clusterName)
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
	credential, err := GetClusterCredentialV1(ctx, platformClient, cluster)
	if err != nil {
		log.Errorf("unable to get credential %v", err)
		return nil, nil, "", err
	}

	transport, err := BuildTransportV1(credential)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}
	host, err := util.ClusterV1Host(cluster)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}

	token := ""
	if credential.Token != nil {
		token = *credential.Token
	}
	return &url.URL{
		Scheme: "https",
		Host:   host,
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

// BuildTransport create the http transport for communicate to backend
// kubernetes api server.
func BuildTransportV1(credential *platformv1.ClusterCredential) (http.RoundTripper, error) {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	if len(credential.CACert) > 0 {
		transport.TLSClientConfig = &tls.Config{
			RootCAs: rootCertPool(credential.CACert),
		}
	} else {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if credential.ClientKey != nil && credential.ClientCert != nil {
		cert, err := tls.X509KeyPair(credential.ClientCert, credential.ClientKey)
		if err != nil {
			return nil, err
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}

	return transport, nil
}

// rootCertPool returns nil if caData is empty.  When passed along, this will mean "use system CAs".
// When caData is not empty, it will be the ONLY information used in the CertPool.
func rootCertPool(caData []byte) *x509.CertPool {
	// What we really want is a copy of x509.systemRootsPool, but that isn't exposed.  It's difficult to build (see the go
	// code for a look at the platform specific insanity), so we'll use the fact that RootCAs == nil gives us the system values
	// It doesn't allow trusting either/or, but hopefully that won't be an issue
	if len(caData) == 0 {
		return nil
	}

	// if we have caData, use it
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)
	return certPool
}

// BuildExternalClientSetWithName creates the clientset of kubernetes by given cluster
// name and returns it.
func BuildExternalClientSetWithName(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, name string) (*kubernetes.Clientset, error) {
	cluster, err := platformClient.Clusters().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	clientset, err := BuildExternalClientSet(ctx, cluster, platformClient)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// BuildExternalClientSet creates the clientset of kubernetes by given cluster object and returns it.
func BuildExternalClientSet(ctx context.Context, cluster *platformv1.Cluster, client platformversionedclient.PlatformV1Interface) (*kubernetes.Clientset, error) {
	credential, err := GetClusterCredentialV1(ctx, client, cluster)
	if err != nil {
		return nil, err
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	return util.BuildVersionedClientSet(cluster, credential)
}

// GetClusterCredentialV1 returns the versioned cluster's credential
func GetClusterCredentialV1(ctx context.Context, client platformversionedclient.PlatformV1Interface, cluster *platformv1.Cluster) (*platformv1.ClusterCredential, error) {
	var (
		credential *platformv1.ClusterCredential
		err        error
	)

	if cluster.Spec.ClusterCredentialRef != nil {
		credential, err = client.ClusterCredentials().Get(ctx, cluster.Spec.ClusterCredentialRef.Name, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
	} else {
		clusterName := cluster.Name
		fieldSelector := fields.OneTermEqualSelector("clusterName", clusterName).String()
		clusterCredentials, err := client.ClusterCredentials().List(ctx, metav1.ListOptions{FieldSelector: fieldSelector})
		if err != nil {
			return nil, err
		}
		if len(clusterCredentials.Items) == 0 {
			return nil, errors.NewNotFound(platform.Resource("ClusterCredential"), clusterName)
		}
		credential = &clusterCredentials.Items[0]
	}

	return credential, nil
}
