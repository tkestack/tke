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

package authenticator

import (
	"context"
	"crypto/tls"
	gooidc "github.com/coreos/go-oidc"
	"gopkg.in/square/go-jose.v2"
	"k8s.io/apiserver/pkg/server"
	"net/http"
	"tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/util/log"
)

// ProviderHookHandler is an token authentication handler to provide a post start hook for the api server.
type ProviderHookHandler struct {
	handler         *TokenAuthenticator
	publicAddress   string
	internalAddress string
	ctx             context.Context
}

// NewProviderHookHandler creates a new ProviderHookHandler object.
func NewProviderHookHandler(ctx context.Context, publicAddress string, internalAddress string, handler *TokenAuthenticator) *ProviderHookHandler {
	return &ProviderHookHandler{
		ctx:             ctx,
		publicAddress:   publicAddress,
		internalAddress: internalAddress,
		handler:         handler,
	}
}

// PostStartHook provides a function that is called after the server has started.
func (h *ProviderHookHandler) PostStartHook() (string, server.PostStartHookFunc, error) {
	return "create-authn-provider", func(_ server.PostStartHookContext) error {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Transport: tr,
		}
		ctx := gooidc.ClientContext(h.ctx, client)
		idTokenVerifier, err := oidc.NewIDTokenVerifier(ctx, h.internalAddress, h.publicAddress, &gooidc.Config{
			SkipClientIDCheck: true,
			SupportedSigningAlgs: []string{
				string(jose.RS256),
			},
		})
		if err != nil {
			log.Error("Failed to create the oidc verifier", log.Err(err))
			return err
		}

		h.handler.idTokenVerifier = idTokenVerifier
		return nil
	}, nil
}
