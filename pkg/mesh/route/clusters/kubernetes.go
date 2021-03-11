/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package clusters

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/pkg/mesh/services/rest"
	"tkestack.io/tke/pkg/mesh/util/errors"
	"tkestack.io/tke/pkg/mesh/util/web"
	"tkestack.io/tke/pkg/util/log"
)

// Get get cluster
func (h *clusterHandler) Get(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest

	clusterName := req.PathParameter("cluster")

	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	ctx := req.Request.Context()
	cl, err := h.clusterService.Get(ctx, clusterName)
	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleAPIError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = cl
}

// List list cluster
func (h *clusterHandler) List(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	ctx := req.Request.Context()
	cls, err := h.clusterService.List(ctx)
	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleAPIError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = cls
}

// ListNamespacs list namespaces
func (h *clusterHandler) ListNamespaces(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	clusterName := req.PathParameter("cluster")

	ctx := req.Request.Context()
	ret, err := h.clusterService.ListNamespaces(ctx, clusterName)
	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleAPIError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = ret
}

// ListServices list services
func (h clusterHandler) ListServices(req *restful.Request, resp *restful.Response) {
	result := rest.NewResult(false, "")
	status := http.StatusBadRequest
	defer func() {
		_ = resp.WriteHeaderAndEntity(status, result)
	}()

	cluster := req.PathParameter(web.ClusterKey)
	namespace := req.PathParameter(web.NamespaceKey)

	ctx := req.Request.Context()
	list, err := h.clusterService.ListServices(ctx, cluster, namespace, nil)

	if err != nil {
		var ok bool
		ok, status, result.Err = errors.HandleAPIError(err)
		if !ok {
			status = http.StatusInternalServerError
			result.Err = err.Error()
		}
		return
	}

	status = http.StatusOK
	result.Result = true
	result.Data = list
}

func (h clusterHandler) Proxy(req *restful.Request, resp *restful.Response) {
	clusterName := req.PathParameter(web.ClusterKey)
	clusterCrdResourceAPIPath := req.PathParameter("path")

	transport, host, err := h.clusterService.Proxy(context.TODO(), clusterName)
	if err != nil {
		log.Errorf("%v", err)
		/*{
			"kind": "Status",
			"apiVersion": "v1",
			"metadata": {},
			"status": "Failure",
			"message": "deployments.apps \"provider-1-0-000\" not found",
			"reason": "NotFound",
			"details": {
				"name": "provider-1-0-000",
				"group": "apps",
				"kind": "deployments"
			},
			"code": 404
		}*/
		result := metav1.Status{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Status",
				APIVersion: "v1",
			},
			Status:  "Failure",
			Message: err.Error(),
			Reason:  metav1.StatusReasonInternalError,
			Code:    http.StatusInternalServerError,
		}
		_ = resp.WriteHeaderAndEntity(http.StatusInternalServerError, result)
		return
	}

	location := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   clusterCrdResourceAPIPath,
	}
	reverseProxy := &nrdProxyHandler{
		location:  location,
		transport: transport,
	}
	reverseProxy.ServeHTTP(resp, req.Request)
}

type nrdCloseTransport struct {
	originalTransport http.RoundTripper
}

var _ http.RoundTripper = &nrdCloseTransport{}

func (t *nrdCloseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Close = true
	return t.originalTransport.RoundTrip(req)
}

type nrdProxyHandler struct {
	transport http.RoundTripper
	location  *url.URL
	token     string
}

func (h *nrdProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery

	// WithContext creates a shallow clone of the request with the new context.
	newReq := req.WithContext(context.Background())
	out := make(http.Header, len(req.Header))
	for key, values := range req.Header {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	newReq.Header = out
	newReq.URL = &loc
	if h.token != "" {
		newReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(h.token)))
	}

	reserveProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: h.location.Scheme,
		Host:   h.location.Host,
	})
	reserveProxy.Transport = h.transport
	reserveProxy.FlushInterval = 100 * time.Millisecond
	reserveProxy.ServeHTTP(w, newReq)
}
