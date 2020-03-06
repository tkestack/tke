/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

package identityprovider

import (
	"fmt"
	"net/http"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	genericrequest "k8s.io/apiserver/pkg/endpoints/request"
	"tkestack.io/tke/api/auth"
)

type DexHander struct {
	handler http.Handler
}

func (t *DexHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.handler == nil {
		responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(fmt.Errorf("dex oidc server init failed")).Status(), w)
		return
	}

	// Inject header and form value for identity provider login use.
	if strings.HasPrefix(r.URL.String(), fmt.Sprintf("/%s/auth", auth.IssuerName)) {
		for k, v := range r.Header {
			r = r.WithContext(genericrequest.WithValue(r.Context(), k, v))
		}
		err := r.ParseForm()
		if err == nil {
			for k, v := range r.Form {
				r = r.WithContext(genericrequest.WithValue(r.Context(), k, v))
			}
		}
	}

	t.handler.ServeHTTP(w, r)
}
