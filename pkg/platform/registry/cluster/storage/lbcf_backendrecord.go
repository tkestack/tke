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
	"sort"
	"strings"
	"time"
	"tkestack.io/tke/pkg/util/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// LBCFAction defines the type of action
type LBCFAction string

const (
	// LBCFEvents defines the action of lbcf events
	LBCFEvents LBCFAction = "events"
)

// LBCFBackendRecordREST implements proxy LBCF LoadBalancer request to cluster of user.
type LBCFBackendRecordREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *LBCFBackendRecordREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *LBCFBackendRecordREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.LBCFProxyOptions{}, false, ""
}

// Connect returns a handler for the kube-apiserver proxy
func (r *LBCFBackendRecordREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
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
	return &lbcfBackendRecordHandler{
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

// New creates a new LBCF proxy options object
func (r *LBCFBackendRecordREST) New() runtime.Object {
	return &platform.LBCFProxyOptions{}
}

type lbcfBackendRecordHandler struct {
	transport         http.RoundTripper
	cluster           *platform.Cluster
	clusterCredential *platform.ClusterCredential
	location          *url.URL
	token             string
	namespace         string
	name              string
	action            string
}

func (h *lbcfBackendRecordHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
		loc.Path = fmt.Sprintf("%s/namespaces/%s/backendrecords", prefix, h.namespace)
	} else {
		loc.Path = fmt.Sprintf("%s/namespaces/%s/backendrecords/%s", prefix, h.namespace, h.name)
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

func (h *lbcfBackendRecordHandler) serveAction(w http.ResponseWriter, req *http.Request) {
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
	default:
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("unsupported action"), w)
	}
}

func (h *lbcfBackendRecordHandler) getEventList() (*corev1.EventList, error) {
	return getLBCFEvents(h.cluster, h.clusterCredential, recordResource, "BackendRecord", h.namespace, h.name)
}

func getLBCFEvents(cluster *platform.Cluster, credential *platform.ClusterCredential, resource schema.GroupVersionResource, kind, namespace, name string) (*corev1.EventList, error) {
	var clusterv1 platformv1.Cluster
	if err := apiplatformv1.Convert_platform_Cluster_To_v1_Cluster(cluster, &clusterv1, nil); err != nil {
		return nil, err
	}
	var clusterCredential platformv1.ClusterCredential
	if err := apiplatformv1.Convert_platform_ClusterCredential_To_v1_ClusterCredential(credential, &clusterCredential, nil); err != nil {
		return nil, err
	}
	dynamicClient, err := util.BuildExternalDynamicClientSet(&clusterv1, &clusterCredential)
	if err != nil {
		return nil, err
	}
	obj, err := dynamicClient.Resource(resource).Namespace(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	kubeclient, err := util.BuildClientSet(cluster, credential)
	if err != nil {
		return nil, err
	}

	eventList, err := util.GetEvents(kubeclient, string(obj.GetUID()), obj.GetNamespace(), obj.GetName(), kind)
	if err != nil {
		return nil, err
	}

	var events util.EventSlice
	for _, event := range eventList.Items {
		events = append(events, event)
	}

	sort.Sort(events)

	return &corev1.EventList{
		Items: events,
	}, nil
}
