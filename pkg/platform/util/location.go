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
	"path"

	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/apiserver/filter"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/request"
	restclient "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
)

// APIServerLocationByCluster returns a URL and transport which one can use to
// send traffic for the specified cluster api server.
func APIServerLocationByCluster(ctx context.Context, cluster *platformv1.Cluster) (*url.URL, http.RoundTripper, string, error) {
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, nil, "", errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}
	if cluster.Status.Phase != platformv1.ClusterRunning {
		return nil, nil, "", errors.NewServiceUnavailable(fmt.Sprintf("cluster %s status is abnormal", cluster.ObjectMeta.Name))
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, nil, "", errors.NewForbidden(platform.Resource("clusters"), cluster.ObjectMeta.Name, fmt.Errorf("cluster is been locked"))
	}

	provider, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}

	restconfig, err := provider.GetRestConfig(ctx, cluster, username)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}

	transport, err := restclient.TransportFor(restconfig)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}
	address, err := clusterAddress(cluster)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}

	token := restconfig.BearerToken

	// Otherwise, return the requested scheme and port, and the proxy transport
	return &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%v:%v", address.Host, address.Port),
		Path:   address.Path,
	}, transport, token, nil
}

// APIServerLocation returns a URL and transport which one can use to send
// traffic for the specified kube api server.
func APIServerLocation(ctx context.Context, platformClient platforminternalclient.PlatformInterface) (*url.URL, http.RoundTripper, string, error) {
	clusterName := filter.ClusterFrom(ctx)
	if clusterName == "" {
		return nil, nil, "", errors.NewBadRequest("clusterName is required")
	}

	requestInfo, ok := request.RequestInfoFrom(ctx)
	if !ok {
		return nil, nil, "", errors.NewBadRequest("unable to get request info from context")
	}
	cluster, err := platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})

	if err != nil {
		return nil, nil, "", err
	}

	clusterv1 := &platformv1.Cluster{}
	err = platformv1.Convert_platform_Cluster_To_v1_Cluster(cluster, clusterv1, nil)
	if err != nil {
		return nil, nil, "", errors.NewInternalError(err)
	}

	location, transport, token, err := APIServerLocationByCluster(ctx, clusterv1)
	if err != nil {
		return nil, nil, "", err
	}
	location.Path = path.Join(location.Path, requestInfo.Path)
	return location, transport, token, nil
}
