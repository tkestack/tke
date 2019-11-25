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
	corev1api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/registry/rest"
	"net/http"
	"net/url"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// Support both GET and POST methods. We must support GET for browsers that want
// to use WebSockets.
var upgradeableMethods = []string{"GET", "POST"}

// ExecREST implements the exec endpoint for a Pod
type ExecREST struct {
	platformClient platforminternalclient.PlatformInterface
}

// New creates a new Pod log options object
func (r *ExecREST) New() runtime.Object {
	return &corev1api.PodExecOptions{}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *ExecREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &corev1api.PodExecOptions{}, false, ""
}

// ConnectMethods returns the methods supported by exec
func (r *ExecREST) ConnectMethods() []string {
	return upgradeableMethods
}

// Connect returns a handler for the pod exec proxy
func (r *ExecREST) Connect(ctx context.Context, name string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	execOpts, ok := opts.(*corev1api.PodExecOptions)
	if !ok {
		return nil, fmt.Errorf("invalid options object: %#v", opts)
	}

	location, transport, token, err := util.APIServerLocation(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	if execOpts.Container != "" {
		params.Add("container", execOpts.Container)
	}
	if execOpts.Command != nil {
		params["command"] = execOpts.Command
	}
	params.Add("stdin", "true")
	params.Add("stdout", "true")
	params.Add("stderr", "true")
	params.Add("tty", "true")

	location.RawQuery = params.Encode()

	return &execHandler{
		upgradeAwareHandler: newThrottledUpgradeAwareProxyHandler(location, transport, false, true, responder),
		token:               token,
	}, nil
}

type execHandler struct {
	upgradeAwareHandler *proxy.UpgradeAwareHandler
	token               string
}

func (h *execHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newReq := req.WithContext(req.Context())
	newReq.Header = utilnet.CloneHeader(req.Header)
	if h.token != "" {
		newReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.token))
	}
	h.upgradeAwareHandler.ServeHTTP(w, newReq)
}

func newThrottledUpgradeAwareProxyHandler(location *url.URL, transport http.RoundTripper, wrapTransport, upgradeRequired bool, responder rest.Responder) *proxy.UpgradeAwareHandler {
	handler := proxy.NewUpgradeAwareHandler(location, transport, wrapTransport, upgradeRequired, proxy.NewErrorResponder(responder))
	handler.InterceptRedirects = true
	handler.RequireSameHostRedirects = true
	handler.MaxBytesPerSec = 0
	return handler
}
