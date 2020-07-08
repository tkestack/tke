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

package controller

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/apimachinery/pkg/runtime"
	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	apiserver "k8s.io/apiserver/pkg/server"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/apiserver/pkg/server/routes"
	componentconfig "k8s.io/component-base/config"
	"net/http"
	goruntime "runtime"
)

// BuildHandlerChain builds a handler chain with a base handler and CompletedConfig.
func BuildHandlerChain(apiHandler http.Handler, authorizationInfo *apiserver.AuthorizationInfo, authenticationInfo *apiserver.AuthenticationInfo, codec runtime.NegotiatedSerializer) http.Handler {
	requestInfoResolver := &apirequest.RequestInfoFactory{}
	failedHandler := genericapifilters.Unauthorized(codec, false)

	handler := apiHandler
	if authorizationInfo != nil {
		handler = genericapifilters.WithAuthorization(apiHandler, authorizationInfo.Authorizer, codec)
	}
	if authenticationInfo != nil {
		handler = genericapifilters.WithAuthentication(handler, authenticationInfo.Authenticator, failedHandler, nil)
	}
	handler = genericapifilters.WithRequestInfo(handler, requestInfoResolver)
	handler = genericfilters.WithPanicRecovery(handler)

	return handler
}

// NewBaseHandler takes in CompletedConfig and returns a handler.
func NewBaseHandler(c *componentconfig.DebuggingConfiguration, checks ...healthz.HealthChecker) *mux.PathRecorderMux {
	m := mux.NewPathRecorderMux("controller-manager")
	healthz.InstallHandler(m, checks...)
	if c.EnableProfiling {
		routes.Profiling{}.Install(m)
		if c.EnableContentionProfiling {
			goruntime.SetBlockProfileRate(1)
		}
	}
	m.Handle("/metrics", promhttp.Handler())

	return m
}
