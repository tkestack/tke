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
	"k8s.io/apimachinery/pkg/fields"
	"net/http"
	"net/url"
	platformv1 "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	v1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/apiserver/filter"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/request"
	tkeinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
)

// APIServerLocationByCluster returns a URL and transport which one can use to
// send traffic for the specified cluster api server.
func APIServerLocationByCluster(ctx context.Context, cluster *platform.Cluster, platformClient tkeinternalclient.PlatformInterface) (*url.URL, http.RoundTripper, string, error) {
	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, nil, "", errors.NewBadRequest("unable to get request info from context")
	}

	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, nil, "", errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}
	if cluster.Status.Phase != platform.ClusterRunning {
		return nil, nil, "", errors.NewServiceUnavailable(fmt.Sprintf("cluster %s status is abnormal", cluster.ObjectMeta.Name))
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, nil, "", errors.NewForbidden(platform.Resource("clusters"), cluster.ObjectMeta.Name, fmt.Errorf("cluster is been locked"))
	}

	clusterCredential, err := ClusterCredential(platformClient, cluster.Name)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}

	transport, err := BuildTransport(clusterCredential)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}
	address, err := ClusterAddress(cluster)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}

	token := ""
	if clusterCredential.Token != nil {
		token = *clusterCredential.Token
	}

	// Otherwise, return the requested scheme and port, and the proxy transport
	return &url.URL{
		Scheme: "https",
		Host:   address,
		Path:   requestInfo.Path,
	}, transport, token, nil
}

// APIServerLocation returns a URL and transport which one can use to send
// traffic for the specified kube api server.
func APIServerLocation(ctx context.Context, platformClient tkeinternalclient.PlatformInterface) (*url.URL, http.RoundTripper, string, error) {
	clusterName := filter.ClusterFrom(ctx)
	if clusterName == "" {
		return nil, nil, "", errors.NewBadRequest("clusterName is required")
	}

	cluster, err := platformClient.Clusters().Get(clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, "", err
	}

	return APIServerLocationByCluster(ctx, cluster, platformClient)
}

// ClusterCredential returns the cluster's credential
func ClusterCredential(client tkeinternalclient.PlatformInterface, clusterName string) (*platform.ClusterCredential, error) {
	fieldSelector := fields.OneTermEqualSelector("clusterName", clusterName).String()
	clusterCredentials, err := client.ClusterCredentials().List(metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return nil, err
	}
	if len(clusterCredentials.Items) == 0 {
		return nil, errors.NewNotFound(platform.Resource("ClusterCredential"), clusterName)
	}

	return &clusterCredentials.Items[0], nil
}

// ClusterCredentialV1 returns the versioned cluster's credential
func ClusterCredentialV1(client platformv1.PlatformV1Interface, clusterName string) (*v1.ClusterCredential, error) {
	fieldSelector := fields.OneTermEqualSelector("clusterName", clusterName).String()
	clusterCredentials, err := client.ClusterCredentials().List(metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return nil, err
	}
	if len(clusterCredentials.Items) == 0 {
		return nil, errors.NewNotFound(platform.Resource("ClusterCredential"), clusterName)
	}

	return &clusterCredentials.Items[0], nil
}
