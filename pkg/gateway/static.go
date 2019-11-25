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

package gateway

import (
	"golang.org/x/oauth2"
	"io"
	"k8s.io/apiserver/pkg/server/mux"
	"net/http"
	"strings"
	"tkestack.io/tke/pkg/gateway/assets"
	"tkestack.io/tke/pkg/gateway/auth"
	"tkestack.io/tke/pkg/gateway/token"
	"tkestack.io/tke/pkg/util/log"
)

func registerStaticRoute(m *mux.PathRecorderMux, oauthConfig *oauth2.Config, disableOIDCProxy bool) {
	m.HandlePrefix("/", withRewrite(assets.Server, oauthConfig, disableOIDCProxy))
}

func withRewrite(handler http.Handler, oauthConfig *oauth2.Config, disableOIDCProxy bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestPath := strings.ToLower(strings.TrimLeft(request.URL.Path, "/"))
		if requestPath == "" || requestPath == "index.html" {
			serveIndex(writer, request, oauthConfig, disableOIDCProxy)
			return
		}
		nfRW := &notFoundResponseWriter{ResponseWriter: writer}
		handler.ServeHTTP(nfRW, request)
		if nfRW.status == http.StatusNotFound {
			serveIndex(writer, request, oauthConfig, disableOIDCProxy)
		}
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request, oauthConfig *oauth2.Config, disableOIDCProxy bool) {
	if oauthConfig != nil {
		_, err := token.RetrieveToken(r)
		if err != nil {
			auth.RedirectLogin(w, r, oauthConfig, disableOIDCProxy)
			return
		}
	}

	rc, err := assets.Open("index.html")
	if err != nil {
		log.Error("Failed to get index.html file", log.Err(err))
		http.NotFound(w, r)
		return
	}
	defer func() {
		_ = rc.Close()
	}()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, errCopy := io.Copy(w, rc)
	if errCopy != nil {
		http.Error(w, "failed to rewrite file", http.StatusInternalServerError)
		return
	}
}

type notFoundResponseWriter struct {
	http.ResponseWriter // embed http.ResponseWriter
	status              int
}

func (w *notFoundResponseWriter) WriteHeader(status int) {
	// Store the status for our own use
	w.status = status
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *notFoundResponseWriter) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}
	// Lie that we successfully written it
	return len(p), nil
}
