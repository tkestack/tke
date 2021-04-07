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
	"path"
	"sort"
	"strings"
	"time"
	"encoding/json"
	platformv1 "tkestack.io/tke/api/platform/v1"
	corev1 "k8s.io/api/core/v1"
	"tkestack.io/tke/pkg/util/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"tkestack.io/tke/pkg/platform/util"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
)

const (

	// Events is an action that lists events
	Events Action = "events"
)

// HpcREST implements proxy HPC request to cluster of user.    //hpc rest实现了将 hpc的request转发至 用户集群
type HpcREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied    //能够代理的request method
func (r *HpcREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *HpcREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.HpcProxyOptions{}, false, ""
}

// Connect returns a handler for the kube-apiserver proxy
func (r *HpcREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.HpcProxyOptions)

	location, transport, token, err := util.APIServerLocationByCluster(ctx, cluster, r.platformClient)
	credential, err := util.GetClusterCredential(ctx, r.platformClient, cluster)
	if err != nil {
		return nil, err
	}
	return &HpcProxyHandler{
		location:  location,
		transport: transport,
		token:     token,
		namespace: proxyOpts.Namespace,
		name:      proxyOpts.Name,
		action:    proxyOpts.Action,
		cluster:           cluster,
		clusterCredential: credential,
	}, nil
}

// New creates a new HpcCollector proxy options object
func (r *HpcREST) New() runtime.Object {
	return &platform.HpcProxyOptions{}
}
type HpcProxyHandler struct {
	transport http.RoundTripper
	location  *url.URL
	token     string
	namespace string
	name      string
	action    string
	cluster           *platform.Cluster
	clusterCredential *platform.ClusterCredential
}

func (h *HpcProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery

	// todo: Change the apigroup here once the integration pipeline configuration is complete using the hpc in the tkestack group
	prefix := "/apis/autoscaling.cloud.tencent.com/v1"

	if len(h.action) > 0{
		h.serveAction(w,req)
		return
	}

	if len(h.namespace) == 0 && len(h.name) == 0 {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/horizontalpodcronscalers", prefix))
	} else if len(h.name) == 0 {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/namespaces/%s/horizontalpodcronscalers", prefix, h.namespace))
	} else {
		loc.Path = path.Join(loc.Path, fmt.Sprintf("%s/namespaces/%s/horizontalpodcronscalers/%s", prefix, h.namespace, h.name))
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


func (h *HpcProxyHandler) serveAction(w http.ResponseWriter, req *http.Request) {
	if len(h.namespace) == 0 || len(h.name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("namespace and name must be specified"), w)
		return
	}
	switch h.action {
	case string(Events):
		if eventList, err := h.getEventList(req.Context()); err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		} else {
			responsewriters.WriteRawJSON(http.StatusOK, eventList, w)
		}
	default:
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("unsupported action"), w)
	}
}

// Get retrieves the object from the storage. It is required to support Patch.
func (h *HpcProxyHandler) getEventList(ctx context.Context) (*corev1.EventList, error) {
	hpc, err := getHpc(ctx, h.cluster, h.clusterCredential, h.namespace, h.name)
	if err != nil {
		return nil, err
	}

	kubeclient, err := util.BuildClientSet(ctx, h.cluster, h.clusterCredential)
	if err != nil {
		return nil, err
	}
	// Get hpc events
	hpcEvents, err := util.GetEvents(ctx, kubeclient, string(hpc.UID), hpc.Namespace, hpc.Name, "HorizontalPodCronscaler")
	if err != nil {
		return nil, err
	}

	var events util.EventSlice
	for _, event := range hpcEvents.Items {
		events = append(events, event)
	}
	sort.Sort(events)

	return &corev1.EventList{
		Items: events,
	}, nil
}


func getHpc(ctx context.Context, cluster *platform.Cluster, credential *platform.ClusterCredential, namespace, name string) (*util.CustomResource, error) {
	var clusterv1 platformv1.Cluster
	if err := platformv1.Convert_platform_Cluster_To_v1_Cluster(cluster, &clusterv1, nil); err != nil {
		return nil, err
	}
	var clusterCredential platformv1.ClusterCredential
	if err := platformv1.Convert_platform_ClusterCredential_To_v1_ClusterCredential(credential, &clusterCredential, nil); err != nil {
		return nil, err
	}

	dynamicclient, err := util.BuildExternalDynamicClientSet(&clusterv1, &clusterCredential)
	if err != nil {
		return nil, err
	}
	hpcResource := schema.GroupVersionResource{Group: "autoscaling.cloud.tencent.com", Version: "v1", Resource: "horizontalpodcronscalers"}
	content, err := dynamicclient.Resource(hpcResource).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, err := content.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var hpc util.CustomResource
	if err := json.Unmarshal(data, &hpc); err != nil {
		return nil, err
	}
	return &hpc, nil
}
