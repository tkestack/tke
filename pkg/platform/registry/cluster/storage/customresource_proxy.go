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
 */

package storage

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/platform/proxy"
)

type CustomResourceHandler struct {
	LoopbackClientConfig *rest.Config
}

var (
	// Scheme is the default instance of runtime.Scheme to which types in the TKE API are already registered.
	Scheme = runtime.NewScheme()
	// Codecs provides access to encoding and decoding for the scheme
	Codecs = serializer.NewCodecFactory(Scheme)

	unversionedVersion = schema.GroupVersion{Group: "", Version: "v1"}
	unversionedTypes   = []runtime.Object{
		&metav1.Status{},
	}
)

func init() {
	Scheme.AddUnversionedTypes(unversionedVersion, unversionedTypes...)
}

// ServeHTTP is a proxy for unregister custom resource
func (n *CustomResourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestInfo, ok := request.RequestInfoFrom(r.Context())
	if !ok {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
		)
		return
	}

	// crd resource start with /apis
	if requestInfo.APIPrefix != "apis" {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("request crd is validate")),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
		)
		return
	}

	platformClient := platforminternalclient.NewForConfigOrDie(n.LoopbackClientConfig)
	config, err := proxy.GetConfig(r.Context(), platformClient)

	if err != nil {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(err),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
		)
		return
	}

	TLSClientConfig := &tls.Config{}
	TLSClientConfig.InsecureSkipVerify = true
	if config.TLSClientConfig.CertData != nil && config.TLSClientConfig.KeyData != nil {
		cert, err := tls.X509KeyPair(config.TLSClientConfig.CertData, config.TLSClientConfig.KeyData)
		if err != nil {
			responsewriters.ErrorNegotiated(
				apierrors.NewGenericServerResponse(http.StatusUnauthorized, requestInfo.Verb, schema.GroupResource{
					Group:    requestInfo.APIGroup,
					Resource: requestInfo.Resource,
				}, requestInfo.Name,
					err.Error(), 0, true),
				Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
			)
			return
		}
		TLSClientConfig.Certificates = []tls.Certificate{cert}
	} else if config.BearerToken == "" {
		responsewriters.ErrorNegotiated(
			apierrors.NewGenericServerResponse(http.StatusUnauthorized, requestInfo.Verb, schema.GroupResource{
				Group:    requestInfo.APIGroup,
				Resource: requestInfo.Resource,
			}, requestInfo.Name,
				fmt.Sprintf("%s has NO BearerToken", filter.ClusterFrom(r.Context())),
				0, true),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
		)
		return
	}

	reserveProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		// TODO: support apiserver with path
		Host: strings.TrimPrefix(config.Host, "https://"),
	})
	reserveProxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       TLSClientConfig,
	}

	reserveProxy.Director = buildDirector(r, config)
	reserveProxy.FlushInterval = 100 * time.Millisecond
	reserveProxy.ServeHTTP(w, r)
}

func buildDirector(r *http.Request, config *rest.Config) func(req *http.Request) {
	clusterName := filter.ClusterFrom(r.Context())
	userName, tenantID := authentication.UsernameAndTenantID(r.Context())
	return func(req *http.Request) {
		req.Header.Set(filter.ClusterNameHeaderKey, clusterName)
		req.Header.Set("X-Remote-User", userName)
		req.Header.Set("X-Remote-Extra-TenantID", tenantID)
		if config.BearerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.BearerToken))
		}
		req.URL = &url.URL{
			Scheme: "https",
			// TODO: support apiserver with path
			Host: strings.TrimPrefix(config.Host, "https://"),
			Path: r.RequestURI,
		}
	}
}
