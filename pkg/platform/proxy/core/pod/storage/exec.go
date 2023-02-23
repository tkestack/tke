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
	"net/url"

	corev1api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
	"tkestack.io/tke/pkg/util/log"
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

	location, transport, _, err := util.APIServerLocation(ctx, r.platformClient)
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
	if execOpts.Stdin {
		params.Add("stdin", "true")
	} else {
		params.Add("stdin", "false")
	}
	if execOpts.Stdout {
		params.Add("stdout", "true")
	} else {
		params.Add("stdout", "false")
	}
	if execOpts.Stderr {
		params.Add("stderr", "true")
	} else {
		params.Add("stderr", "false")
	}
	if execOpts.TTY {
		params.Add("tty", "true")
	} else {
		params.Add("tty", "false")
	}

	location.RawQuery = params.Encode()

	return &execHandler{
		upgradeAwareHandler: newThrottledUpgradeAwareProxyHandler(location, transport, false, true, responder),
	}, nil
}

type execHandler struct {
	upgradeAwareHandler *proxy.UpgradeAwareHandler
}

func (h *execHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqClone := utilnet.CloneRequest(req)
	reqClone.URL.Scheme = h.upgradeAwareHandler.Location.Scheme
	reqClone.URL.Host = h.upgradeAwareHandler.Location.Host
	reqClone.Header = nil
	resp, err := h.upgradeAwareHandler.Transport.RoundTrip(reqClone)
	if err != nil {
		log.Warnf("err %v", err)
	}
	outReq := resp.Request
	for k, vs := range req.Header {
		for _, v := range vs {
			outReq.Header.Add(k, v)
		}
	}
	log.Errorf("header: %v", outReq)
	h.upgradeAwareHandler.ServeHTTP(w, outReq)
}

func newThrottledUpgradeAwareProxyHandler(location *url.URL, transport http.RoundTripper, wrapTransport, upgradeRequired bool, responder rest.Responder) *proxy.UpgradeAwareHandler {
	handler := proxy.NewUpgradeAwareHandler(location, transport, wrapTransport, upgradeRequired, proxy.NewErrorResponder(responder))
	handler.MaxBytesPerSec = 0
	return handler
}
