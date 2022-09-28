/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2022 Tencent. All Rights Reserved.
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

package api

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/gateway/token"
	"tkestack.io/tke/pkg/util/log"
)

func registerPlatformInfoRoute(container *restful.Container, oidcAuthenticator *oidc.Authenticator, client kubernetes.Interface) {
	ws := new(restful.WebService)
	ws.Path(fmt.Sprintf("/apis/%s/%s/platforminfo", GroupName, Version))
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON, restful.MIME_OCTET)
	ws.Route(ws.
		GET("/").
		Doc("get platform info of TKE").
		Operation("getPlatformInfo").
		Returns(http.StatusOK, "Ok", v1.ConfigMap{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.ConfigMap{}).
		To(handlePlatformInfoFunc(oidcAuthenticator, client)))
	container.Add(ws)
}

func handlePlatformInfoFunc(oidcAuthenticator *oidc.Authenticator, client kubernetes.Interface) func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		t, err := token.RetrieveToken(request.Request)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized(err.Error()), response.ResponseWriter)
			return
		}
		_, authenticated, _ := oidcAuthenticator.AuthenticateToken(request.Request.Context(), t.ID)

		if !authenticated {
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized("invalid token"), response.ResponseWriter)
			return
		}

		cm, err := client.CoreV1().ConfigMaps("tke").Get(request.Request.Context(), "cluster-info", metav1.GetOptions{})
		if err != nil {
			log.Errorf("get cluster-info failed: %v", err)
			cm = &v1.ConfigMap{}
		}
		responsewriters.WriteRawJSON(http.StatusOK, *cm, response.ResponseWriter)
	}
}
