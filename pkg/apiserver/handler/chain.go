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

package handler

import (
	"github.com/rs/cors"
	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	"net/http"
	"regexp"
	apiserverfilter "tkestack.io/tke/pkg/apiserver/filter"
	authfilter "tkestack.io/tke/pkg/auth/filter"
)

type Chain func(apiHandler http.Handler, c *genericapiserver.Config) http.Handler

// BuildHandlerChain returns the chained http Handler.
func BuildHandlerChain(ignoreAuthPathPrefixes []string) Chain {
	return func(apiHandler http.Handler, c *genericapiserver.Config) http.Handler {
		handler := authfilter.WithTKEAuthorization(apiHandler, c.Authorization.Authorizer, c.Serializer, ignoreAuthPathPrefixes)
		handler = genericfilters.WithMaxInFlightLimit(handler, c.MaxRequestsInFlight, c.MaxMutatingRequestsInFlight, c.LongRunningFunc)
		handler = genericapifilters.WithImpersonation(handler, c.Authorization.Authorizer, c.Serializer)
		handler = genericapifilters.WithAudit(handler, c.AuditBackend, c.AuditPolicyChecker, c.LongRunningFunc)
		failedHandler := genericapifilters.Unauthorized(c.Serializer, c.Authentication.SupportsBasicAuth)
		failedHandler = genericapifilters.WithFailedAuthenticationAudit(failedHandler, c.AuditBackend, c.AuditPolicyChecker)
		handler = apiserverfilter.WithAuthentication(handler, c.Authentication.Authenticator, failedHandler, c.Authentication.APIAudiences, ignoreAuthPathPrefixes)

		corsHandler := cors.New(cors.Options{
			AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
			AllowedHeaders: []string{
				"Content-Type",
				"Content-Length",
				"Accept-Encoding",
				"X-CSRF-Token",
				"Authorization",
				"X-Requested-With",
				"If-Modified-Since",
				"X-TKE-ClusterName",
				"X-TKE-FuzzyResourceName",
				"X-CsrfCode",
				"X-Referer",
				"X-SeqId",
				"X-Remote-Extra-RequestID",
			},
			ExposedHeaders: []string{
				"Date",
				"X-Remote-Extra-RequestID",
			},
			AllowCredentials: true,
			MaxAge:           86400,
			AllowOriginFunc: func(origin string) bool {
				for _, regexpStr := range c.CorsAllowedOriginList {
					r, err := regexp.Compile(regexpStr)
					if err != nil {
						return false
					}
					if r.MatchString(origin) {
						return true
					}
				}
				return false
			},
		})
		handler = corsHandler.Handler(handler)
		handler = genericfilters.WithTimeoutForNonLongRunningRequests(handler, c.LongRunningFunc, c.RequestTimeout)
		handler = genericfilters.WithWaitGroup(handler, c.LongRunningFunc, c.HandlerChainWaitGroup)
		handler = genericapifilters.WithRequestInfo(handler, c.RequestInfoResolver)
		handler = apiserverfilter.WithLocal(handler)
		handler = apiserverfilter.WithRequestID(handler)
		handler = genericfilters.WithPanicRecovery(handler)
		return handler
	}
}
