/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package assets

import (
	"compress/gzip"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"

	"golang.org/x/oauth2"
	"k8s.io/apiserver/pkg/server/mux"
	"tkestack.io/tke/pkg/gateway/auth"
	"tkestack.io/tke/pkg/gateway/token"
)

const (
	RootDir       = "assets/"
	IndexTmplName = "index.tmpl.html"
	IndexFileName = "index.html"
)

var (
	indexReg = regexp.MustCompile(`/tkestack.*`)
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func RegisterRoute(m *mux.PathRecorderMux, oauthConfig *oauth2.Config, disableOIDCProxy bool) {
	handler := func() http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			if oauthConfig != nil {
				_, err := token.RetrieveToken(request)
				if err != nil {
					auth.RedirectLogin(writer, request, oauthConfig, disableOIDCProxy)
					return
				}
			}

			if indexReg.MatchString(request.URL.Path) {
				request.URL.Path = "/"
			}
			http.ServeFile(
				writer,
				request,
				path.Join(RootDir, request.URL.Path))
		}
	}

	m.HandlePrefix("/", gzipHandler(handler()))
}
