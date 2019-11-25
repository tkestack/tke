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

package storage

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"strings"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/apiserver/cluster"
	"tkestack.io/tke/pkg/platform/util"
)

// DrainREST implements list versions of cluster
type DrainREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *DrainREST) ConnectMethods() []string {
	return []string{"POST"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *DrainREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.HelmProxyOptions{}, true, "path"
}

// Connect returns a handler for the helm-api proxy
func (r *DrainREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	c := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, c); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.HelmProxyOptions)

	clientset, err := util.ClientSetByCluster(ctx, c, r.platformClient)
	if err != nil {
		return nil, err
	}

	return &drainHandler{
		requestPath: proxyOpts.Path,
		clientset:   clientset,
	}, nil
}

// New creates a new helm proxy options object
func (r *DrainREST) New() runtime.Object {
	return &platform.HelmProxyOptions{}
}

type drainHandler struct {
	requestPath string
	clientset   kubernetes.Interface
}

func (h *drainHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	nodeName := strings.Trim(h.requestPath, "/")
	node, err := h.clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			responsewriters.WriteRawJSON(http.StatusNotFound, errors.NewNotFound(corev1.Resource("Node"), nodeName), w)
			return
		}
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}

	err = cluster.DrainNode(h.clientset, node)
	if err != nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		return
	}
	responsewriters.WriteRawJSON(http.StatusCreated, node, w)
}
