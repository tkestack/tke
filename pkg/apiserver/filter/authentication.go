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

package filter

import (
	"errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/component-base/metrics"
	"k8s.io/component-base/metrics/legacyregistry"
	"net/http"
	"strings"
	"tkestack.io/tke/pkg/util/log"
)

const (
	successLabel = "success"
	failureLabel = "failure"
	errorLabel   = "error"
)

var (
	authenticatedUserCounter = metrics.NewCounterVec(
		&metrics.CounterOpts{
			Name:           "tke_authenticated_user_requests",
			Help:           "Counter of authenticated requests broken out by username.",
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"username"},
	)

	authenticatedAttemptsCounter = metrics.NewCounterVec(
		&metrics.CounterOpts{
			Name: "tke_authentication_attempts",
			Help: "Counter of authenticated attempts.",
		},
		[]string{"result"},
	)

	authenticatedIgnoreCounter = metrics.NewCounterVec(
		&metrics.CounterOpts{
			Name:           "tke_authenticated_ignore_requests",
			Help:           "Counter of ignored authentication requests by path prefix and request method.",
			StabilityLevel: metrics.ALPHA,
		},
		[]string{"path_prefix", "method"},
	)

	defaultIgnoreAuthPathPrefixes = []string{
		"/metrics",
		"/debug",
		"/openapi",
		"/version",
		"/swagger",
		"/favicon.ico",
		"/healthz",
	}
)

func init() {
	legacyregistry.MustRegister(authenticatedUserCounter)
	legacyregistry.MustRegister(authenticatedAttemptsCounter)
}

// WithAuthentication creates an http handler that tries to authenticate the given request as a user, and then
// stores any such user found onto the provided context for the request. If authentication fails or returns an error
// the failed handler is used. On success, "Authorization" header is removed from the request and handler
// is invoked to serve the request.
func WithAuthentication(handler http.Handler, auth authenticator.Request, failed http.Handler, apiAuds authenticator.Audiences, ignorePathPrefixes []string) http.Handler {
	if auth == nil {
		log.Warnf("Authentication is disabled")
		return handler
	}
	allIgnorePathPrefixes := MakeAllIgnoreAuthPathPrefixes(ignorePathPrefixes)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodOptions {
			authenticatedIgnoreCounter.WithLabelValues("ALL", http.MethodOptions).Inc()
			handler.ServeHTTP(w, req)
			return
		}

		ignorePathPrefix := ""
		reqPath := strings.ToLower(req.URL.Path)
		for _, pathPrefix := range allIgnorePathPrefixes {
			if strings.HasPrefix(reqPath, strings.ToLower(pathPrefix)) {
				ignorePathPrefix = pathPrefix
				break
			}
		}
		if ignorePathPrefix != "" {
			authenticatedIgnoreCounter.WithLabelValues(ignorePathPrefix, "ALL").Inc()
			handler.ServeHTTP(w, req)
			return
		}

		if len(apiAuds) > 0 {
			req = req.WithContext(authenticator.WithAudiences(req.Context(), apiAuds))
		}
		resp, ok, err := auth.AuthenticateRequest(req)
		if err != nil || !ok {
			if err != nil {
				log.Errorf("Unable to authenticate the request due to an error: %v", err)
				authenticatedAttemptsCounter.WithLabelValues(errorLabel).Inc()
			} else if !ok {
				authenticatedAttemptsCounter.WithLabelValues(failureLabel).Inc()
			}

			failed.ServeHTTP(w, req)
			return
		}

		// authorization header is not required anymore in case of a successful authentication.
		req.Header.Del("Authorization")

		req = req.WithContext(genericapirequest.WithUser(req.Context(), resp.User))

		authenticatedUserCounter.WithLabelValues(compressUsername(resp.User.GetName())).Inc()
		authenticatedAttemptsCounter.WithLabelValues(successLabel).Inc()

		handler.ServeHTTP(w, req)
	})
}

func MakeAllIgnoreAuthPathPrefixes(pathPrefixes []string) []string {
	if len(pathPrefixes) > 0 {
		return append(defaultIgnoreAuthPathPrefixes, pathPrefixes...)
	}
	return defaultIgnoreAuthPathPrefixes
}

func Unauthorized(s runtime.NegotiatedSerializer, supportsBasicAuth bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if supportsBasicAuth {
			w.Header().Set("WWW-Authenticate", `Basic realm="tke"`)
		}
		ctx := req.Context()
		requestInfo, found := genericapirequest.RequestInfoFrom(ctx)
		if !found {
			responsewriters.InternalError(w, req, errors.New("no RequestInfo found in the context"))
			return
		}

		gv := schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}
		responsewriters.ErrorNegotiated(apierrors.NewUnauthorized("Unauthorized"), s, gv, w, req)
	})
}

// compressUsername maps all possible usernames onto a small set of categories
// of usernames. This is done both to limit the cardinality of the
// authorized_user_requests metric, and to avoid pushing actual usernames in the
// metric.
func compressUsername(username string) string {
	switch {
	// Known internal identities.
	case username == "admin" ||
		username == "client":
		return username
	// Probably an email address.
	case strings.Contains(username, "@"):
		return "email_id"
	// Anything else (custom service accounts, custom external identities, etc.)
	default:
		return "other"
	}
}
