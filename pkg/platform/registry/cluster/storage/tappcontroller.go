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
	"encoding/json"
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
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/util"
)

// Action is the action that specified in URL
type Action string

const (
	// Pods is an action that lists pods
	Pods Action = "pods"
	// Events is an action that lists events
	Events Action = "events"

	tappGroupName = "apps.tkestack.io"
)

// TappControllerREST implements proxy tapp controller request to cluster of user.
type TappControllerREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *TappControllerREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *TappControllerREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &platform.TappControllerProxyOptions{}, false, ""
}

// Connect returns a handler for the kube-apiserver proxy
func (r *TappControllerREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}
	proxyOpts := opts.(*platform.TappControllerProxyOptions)

	location, transport, token, err := util.APIServerLocationByCluster(ctx, cluster, r.platformClient)
	if err != nil {
		return nil, err
	}
	credential, err := util.ClusterCredential(r.platformClient, cluster.Name)
	if err != nil {
		return nil, err
	}
	return &tappControllerProxyHandler{
		location:          location,
		transport:         transport,
		cluster:           cluster,
		clusterCredential: credential,
		token:             token,
		namespace:         proxyOpts.Namespace,
		name:              proxyOpts.Name,
		action:            proxyOpts.Action,
	}, nil
}

// New creates a new tapp proxy options object
func (r *TappControllerREST) New() runtime.Object {
	return &platform.TappControllerProxyOptions{}
}

type tappControllerProxyHandler struct {
	transport         http.RoundTripper
	cluster           *platform.Cluster
	clusterCredential *platform.ClusterCredential
	location          *url.URL
	token             string
	namespace         string
	name              string
	action            string
}

func (h *tappControllerProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery

	prefix := fmt.Sprintf("/apis/%s/v1", tappGroupName)

	if len(h.action) != 0 {
		h.serveAction(w, req)
		return
	}

	if len(h.namespace) == 0 && len(h.name) == 0 {
		loc.Path = fmt.Sprintf("%s/tapps", prefix)
	} else if len(h.name) == 0 {
		loc.Path = fmt.Sprintf("%s/namespaces/%s/tapps", prefix, h.namespace)
	} else {
		loc.Path = fmt.Sprintf("%s/namespaces/%s/tapps/%s", prefix, h.namespace, h.name)
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

func (h *tappControllerProxyHandler) serveAction(w http.ResponseWriter, req *http.Request) {
	if len(h.namespace) == 0 || len(h.name) == 0 {
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("namespace and name must be specified"), w)
		return
	}
	switch h.action {
	case string(Pods):
		if podList, err := h.getPodList(); err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		} else {
			responsewriters.WriteRawJSON(http.StatusOK, podList, w)
		}
	case string(Events):
		if eventList, err := h.getEventList(); err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
		} else {
			responsewriters.WriteRawJSON(http.StatusOK, eventList, w)
		}
	default:
		responsewriters.WriteRawJSON(http.StatusBadRequest, errors.NewBadRequest("unsupported action"), w)
	}
}

func (h *tappControllerProxyHandler) getPodList() (*corev1.PodList, error) {
	tapp, err := getTapp(h.cluster, h.clusterCredential, h.namespace, h.name)
	if err != nil {
		return nil, err
	}

	selector, err := metav1.LabelSelectorAsSelector(tapp.Spec.Selector)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	kubeclient, err := util.BuildClientSet(h.cluster, h.clusterCredential)
	if err != nil {
		return nil, err
	}
	pods, err := kubeclient.CoreV1().Pods(h.namespace).List(metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	podList := &corev1.PodList{
		Items: make([]corev1.Pod, 0),
	}
	for _, pod := range pods.Items {
		for _, podReferences := range pod.ObjectMeta.OwnerReferences {
			if (podReferences.Kind == "TApp") && (podReferences.Name == tapp.Name) {
				podList.Items = append(podList.Items, pod)
			}
		}
	}
	return podList, nil
}

// Get retrieves the object from the storage. It is required to support Patch.
func (h *tappControllerProxyHandler) getEventList() (*corev1.EventList, error) {
	tapp, err := getTapp(h.cluster, h.clusterCredential, h.namespace, h.name)
	if err != nil {
		return nil, err
	}

	kubeclient, err := util.BuildClientSet(h.cluster, h.clusterCredential)
	if err != nil {
		return nil, err
	}
	// Get tapp events
	tappEvents, err := util.GetEvents(kubeclient, string(tapp.UID), tapp.Namespace, tapp.Name, "TApp")
	if err != nil {
		return nil, err
	}

	var events util.EventSlice
	for _, event := range tappEvents.Items {
		events = append(events, event)
	}

	podList, err := h.getPodList()
	if err != nil {
		return nil, err
	}

	// Get pod events
	for _, pod := range podList.Items {
		podEvents, err := util.GetEvents(kubeclient, string(pod.UID), pod.Namespace, pod.Name, "Pod")
		if err != nil {
			return nil, err
		}

		for _, podEvent := range podEvents.Items {
			events = append(events, podEvent)
		}
	}

	sort.Sort(events)

	return &corev1.EventList{
		Items: events,
	}, nil
}

func getTapp(cluster *platform.Cluster, credential *platform.ClusterCredential, namespace, name string) (*util.CustomResource, error) {
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
	tappResource := schema.GroupVersionResource{Group: tappGroupName, Version: "v1", Resource: "tapps"}
	content, err := dynamicclient.Resource(tappResource).Namespace(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, err := content.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var tapp util.CustomResource
	if err := json.Unmarshal(data, &tapp); err != nil {
		return nil, err
	}
	return &tapp, nil
}
