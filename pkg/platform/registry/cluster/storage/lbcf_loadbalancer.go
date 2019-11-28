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
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
	"tkestack.io/tke/pkg/util/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	apiplatformv1 "tkestack.io/tke/api/platform/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/util"
)

const (
	// LBCFGetBackendGroup defines the action of backend groups
	LBCFGetBackendGroup LBCFAction = "backendgroups"
)

// LBCFLoadBalancerREST implements proxy LBCF LoadBalancer request to cluster of user.
type LBCFLoadBalancerREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *LBCFLoadBalancerREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *LBCFLoadBalancerREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.LBCFProxyOptions{}, false, ""
}

// Connect returns a handler for the kube-apiserver proxy
func (r *LBCFLoadBalancerREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.LBCFProxyOptions)

	location, transport, token, err := util.APIServerLocationByCluster(ctx, cluster, r.platformClient)
	if err != nil {
		return nil, err
	}
	credential, err := util.ClusterCredential(r.platformClient, cluster.Name)
	if err != nil {
		return nil, err
	}
	return &lbcfLoadBalancerProxyHandler{
		location:          location,
		cluster:           cluster,
		clusterCredential: credential,
		transport:         transport,
		token:             token,
		namespace:         proxyOpts.Namespace,
		name:              proxyOpts.Name,
		action:            proxyOpts.Action,
	}, nil
}

// New creates a new LBCF LoadBalancer proxy options object
func (r *LBCFLoadBalancerREST) New() runtime.Object {
	return &platform.LBCFProxyOptions{}
}

type lbcfLoadBalancerProxyHandler struct {
	transport         http.RoundTripper
	cluster           *platform.Cluster
	clusterCredential *platform.ClusterCredential
	location          *url.URL
	token             string
	namespace         string
	name              string
	action            string
}

func (h *lbcfLoadBalancerProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery

	// todo: Change the apigroup here once the integration pipeline configuration is complete using the tapp in the tkestack group
	prefix := "/apis/lbcf.tkestack.io/v1beta1"

	if len(h.action) > 0 {
		h.serveAction(w, req)
		return
	}

	if len(h.namespace) == 0 && len(h.name) == 0 {
		loc.Path = fmt.Sprintf("%s/loadbalancers", prefix)
	} else if len(h.name) == 0 {
		loc.Path = fmt.Sprintf("%s/namespaces/%s/loadbalancers", prefix, h.namespace)
	} else {
		loc.Path = fmt.Sprintf("%s/namespaces/%s/loadbalancers/%s", prefix, h.namespace, h.name)
	}

	// WithContext creates a shallow clone of the request with the new context.
	newReq := req.WithContext(context.Background())
	newReq.Header = netutil.CloneHeader(req.Header)
	newReq.URL = &loc
	if h.token != "" {
		newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(h.token)))
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: h.location.Scheme, Host: h.location.Host})
	reverseProxy.Transport = h.transport
	reverseProxy.FlushInterval = 100 * time.Millisecond
	reverseProxy.ErrorLog = log.StdErrLogger()
	reverseProxy.ServeHTTP(w, newReq)
}

func (h *lbcfLoadBalancerProxyHandler) serveAction(w http.ResponseWriter, req *http.Request) {
	if len(h.namespace) == 0 || len(h.name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("namespace and name must be specified"), w)
		return
	}
	switch h.action {
	case string(LBCFEvents):
		if eventList, err := h.getEventList(); err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		} else {
			responsewriters.WriteRawJSON(http.StatusOK, eventList, w)
		}
	case string(LBCFGetBackendGroup):
		if groupList, err := h.getBackendGroups(); err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		} else {
			responsewriters.WriteRawJSON(http.StatusOK, groupList, w)
		}
	default:
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("unsupported action"), w)
	}
}

var (
	lbResource     = schema.GroupVersionResource{Group: "lbcf.tkestack.io", Version: "v1beta1", Resource: "loadbalancers"}
	bgResource     = schema.GroupVersionResource{Group: "lbcf.tkestack.io", Version: "v1beta1", Resource: "backendgroups"}
	recordResource = schema.GroupVersionResource{Group: "lbcf.tkestack.io", Version: "v1beta1", Resource: "backendrecords"}
)

func (h *lbcfLoadBalancerProxyHandler) getEventList() (*corev1.EventList, error) {
	return getLBCFEvents(h.cluster, h.clusterCredential, lbResource, "LoadBalancer", h.namespace, h.name)
}

func (h *lbcfLoadBalancerProxyHandler) getBackendGroups() ([]unstructured.Unstructured, error) {
	var clusterv1 platformv1.Cluster
	if err := apiplatformv1.Convert_platform_Cluster_To_v1_Cluster(h.cluster, &clusterv1, nil); err != nil {
		return nil, err
	}
	var clusterCredential platformv1.ClusterCredential
	if err := apiplatformv1.Convert_platform_ClusterCredential_To_v1_ClusterCredential(h.clusterCredential, &clusterCredential, nil); err != nil {
		return nil, err
	}
	dynamicClient, err := util.BuildExternalDynamicClientSet(&clusterv1, &clusterCredential)
	if err != nil {
		return nil, err
	}

	selector := labels.SelectorFromValidatedSet(labels.Set(map[string]string{
		"lbcf.tkestack.io/lb-name": h.name,
	}))
	groupList, err := dynamicClient.Resource(bgResource).Namespace(h.namespace).List(metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}

	return groupList.Items, nil
}
