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
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/logagent"
	"tkestack.io/tke/pkg/logagent/util"
	"tkestack.io/tke/pkg/util/log"
)
// TokenREST implements the REST endpoint.
type FileDownloadREST struct {
	//rest.Storage
	store *registry.Store
	PlatformClient platformversionedclient.PlatformV1Interface
}

type FileDownloadRequest struct {
	PodName 	string `json:"pod"`
	Namespace 	string `json:"namespace"`
	Container 	string `json:"container"`
	Path      	string  `json:"path"`
}


var _ = rest.Connecter(&FileDownloadREST{})

func (r *FileDownloadREST)  New() runtime.Object {
	return &logagent.LogFileDownload{}
}

func (r *FileDownloadREST)  NewConnectOptions() (runtime.Object, bool, string) {
	return &logagent.LogFileDownload{}, false, ""
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *FileDownloadREST) ConnectMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
}

// Connect returns a handler for the kube-apiserver proxy
func (r *FileDownloadREST) Connect(ctx context.Context, loagentName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	logagentObject, err := r.store.Get(ctx, loagentName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	logagent := logagentObject.(*logagent.LogAgent)
	return &logCollectorProxyHandler{
		clusterId: logagent.Spec.ClusterName,
		platformClient: r.PlatformClient,
		location: &url.URL{Scheme:"http",},
	}, nil
}

type logCollectorProxyHandler struct {
	clusterId string
	platformClient platformversionedclient.PlatformV1Interface
	location  *url.URL
}

func (h *logCollectorProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.WriteResponseError(w, util.ErrorInvalidParameter, "unable to read request body")
		log.Infof("unable to ready body %v", err)
		return
	}
	reqConfig := &FileDownloadRequest{}
	if err := json.Unmarshal(body, reqConfig); err != nil {
		util.WriteResponseError(w, util.ErrorInvalidParameter, "unable to ummarshal request")
		log.Errorf("unable to unmarshal body %v", err)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	hostIp, err := util.GetClusterPodIp(h.clusterId, reqConfig.Namespace,reqConfig.PodName, h.platformClient)
	if err != nil {
		util.WriteResponseError(w, util.ErrorInternalError, "unable to find host for this request")
		log.Errorf("unable to get hostip %v", err)
		return
	}
	log.Infof("get host ip is %v body is %v", hostIp, req.Body)
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery
	loc.Path = "/v1/logfile/download"
	loc.Host = hostIp+":8090"
	newReq := req.WithContext(context.Background())
	newReq.Header = netutil.CloneHeader(req.Header)
	newReq.URL = &loc
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: loc.Scheme, Host: loc.Host})
	reverseProxy.FlushInterval = 100 * time.Millisecond
	reverseProxy.ErrorLog = log.StdErrLogger()
	reverseProxy.ServeHTTP(w, newReq)
}
