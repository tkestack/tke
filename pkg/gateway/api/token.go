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

package api

import (
	"context"
	"fmt"
	gooidc "github.com/coreos/go-oidc"
	"github.com/emicklei/go-restful"
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/gateway/auth"
	"tkestack.io/tke/pkg/gateway/token"
)

// UserInfo defines a data structure containing user information.
type UserInfo struct {
	Name   string              `json:"name"`
	UID    string              `json:"uid"`
	Groups []string            `json:"groups"`
	Extra  map[string][]string `json:"extra"`
}

func registerTokenRoute(container *restful.Container, oauthConfig *oauth2.Config, oidcHTTPClient *http.Client, oidcAuthenticator *oidc.Authenticator, disableOIDCProxy bool) {
	ws := new(restful.WebService)
	ws.Path(fmt.Sprintf("/apis/%s/%s/tokens", GroupName, Version))

	ws.Route(ws.
		POST("/").
		Doc("generate token by username and password").
		Operation("createPasswordToken").
		Produces(restful.MIME_JSON).
		Returns(http.StatusCreated, "Created", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		To(handleTokenGenerateFunc(oauthConfig, oidcHTTPClient)))
	ws.Route(ws.
		GET("info").
		Doc("obtain the user information corresponding to the token").
		Operation("getInfo").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Returns(http.StatusOK, "Ok", UserInfo{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		To(handleTokenInfo(oidcAuthenticator)))
	ws.Route(ws.
		GET("redirect").
		Doc("redirect to OpenID Connect server for authentication").
		Operation("createRedirect").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Returns(http.StatusFound, "Found", v1.Status{}).
		To(handleTokenRedirectFunc(oauthConfig, disableOIDCProxy)))
	ws.Route(ws.
		POST("renew").
		Doc("renew a token by refresh token").
		Operation("createRenewToken").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Returns(http.StatusCreated, "Created", v1.Status{}).
		Returns(http.StatusNoContent, "NoContent", v1.Status{}).
		Returns(http.StatusInternalServerError, "InternalError", v1.Status{}).
		Returns(http.StatusUnauthorized, "Unauthorized", v1.Status{}).
		To(handleTokenRenewFunc(oauthConfig, oidcHTTPClient)))
	container.Add(ws)
}

func handleTokenGenerateFunc(oauthConfig *oauth2.Config, oidcHTTPClient *http.Client) func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		username, password, err := retrievePassword(request.Request)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized(err.Error()), response.ResponseWriter)
			return
		}

		ctx := gooidc.ClientContext(context.Background(), oidcHTTPClient)
		t, err := oauthConfig.PasswordCredentialsToken(ctx, username, password)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), response.ResponseWriter)
			return
		}

		if err := token.ResponseToken(t, response.ResponseWriter); err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), response.ResponseWriter)
			return
		}

		responsewriters.WriteRawJSON(http.StatusCreated, v1.Status{
			Status: v1.StatusSuccess,
			Code:   http.StatusCreated,
		}, response.ResponseWriter)
	}
}

func handleTokenInfo(oidcAuthenticator *oidc.Authenticator) func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		t, err := token.RetrieveToken(request.Request)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized(err.Error()), response.ResponseWriter)
			return
		}
		r, authenticated, err := oidcAuthenticator.AuthenticateToken(request.Request.Context(), t.ID)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), response.ResponseWriter)
			return
		}
		if !authenticated {
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized("invalid token"), response.ResponseWriter)
			return
		}
		userInfo := &UserInfo{
			Name:   r.User.GetName(),
			UID:    r.User.GetUID(),
			Groups: r.User.GetGroups(),
			Extra:  r.User.GetExtra(),
		}
		responsewriters.WriteRawJSON(http.StatusOK, userInfo, response.ResponseWriter)
	}
}

func handleTokenRedirectFunc(oauthConfig *oauth2.Config, disableOIDCProxy bool) func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		// delete cookie
		token.DeleteCookie(response.ResponseWriter)
		// redirect
		auth.RedirectLogin(response.ResponseWriter, request.Request, oauthConfig, disableOIDCProxy)
	}
}

func handleTokenRenewFunc(oauthConfig *oauth2.Config, oidcHTTPClient *http.Client) func(*restful.Request, *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		t, err := token.RetrieveToken(request.Request)
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusUnauthorized, errors.NewUnauthorized(err.Error()), response.ResponseWriter)
			return
		}
		if t.Refresh == "" {
			responsewriters.WriteRawJSON(http.StatusNoContent, v1.Status{
				Status: v1.StatusSuccess,
				Code:   http.StatusNoContent,
			}, response.ResponseWriter)
			return
		}
		ctx := gooidc.ClientContext(context.Background(), oidcHTTPClient)
		tokenSource := oauthConfig.TokenSource(ctx, &oauth2.Token{
			RefreshToken: t.Refresh,
		})
		oauth2Token, err := tokenSource.Token()
		if err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), response.ResponseWriter)
			return
		}
		if err := token.ResponseToken(oauth2Token, response.ResponseWriter); err != nil {
			responsewriters.WriteRawJSON(http.StatusInternalServerError, errors.NewInternalError(err), response.ResponseWriter)
			return
		}
		responsewriters.WriteRawJSON(http.StatusCreated, v1.Status{
			Status: v1.StatusSuccess,
			Code:   http.StatusCreated,
		}, response.ResponseWriter)
	}
}

func retrievePassword(request *http.Request) (string, string, error) {
	userName := request.PostFormValue("username")
	password := request.PostFormValue("password")

	if len(userName) == 0 || len(password) == 0 {
		return "", "", fmt.Errorf("username or password is empty")
	}

	return userName, password, nil
}
