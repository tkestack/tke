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

package auth

import (
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"net/url"
)

const RedirectURIKey = "redirect_uri"

// RedirectLogin to redirect the http request to login page.
func RedirectLogin(w http.ResponseWriter, r *http.Request, oauthConfig *oauth2.Config, disableOIDCProxy bool) {
	oauthURL := oauthConfig.AuthCodeURL(r.URL.String(), oauth2.AccessTypeOffline)
	if !disableOIDCProxy {
		originOAuthURL, err := url.Parse(oauthURL)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
			return
		}
		queries := originOAuthURL.Query()
		if redirectURI := queries.Get(RedirectURIKey); redirectURI != "" {
			redirectURL, err := url.Parse(redirectURI)
			if err != nil {
				responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), w)
				return
			}
			if r.TLS == nil {
				redirectURL.Scheme = "http"
			} else {
				redirectURL.Scheme = "https"
			}
			redirectURL.Host = r.Host
			queries.Set(RedirectURIKey, redirectURL.String())
		}

		newURL := &url.URL{
			Path:       originOAuthURL.Path,
			RawPath:    originOAuthURL.RawPath,
			ForceQuery: originOAuthURL.ForceQuery,
			RawQuery:   queries.Encode(),
			Fragment:   originOAuthURL.Fragment,
		}
		oauthURL = newURL.String()
	}
	http.Redirect(w, r, oauthURL, http.StatusFound)
}
