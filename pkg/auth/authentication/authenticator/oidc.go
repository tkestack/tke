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
	"time"

	"github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	genericauthenticator "k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	oidcclaims "tkestack.io/tke/pkg/auth/authentication/oidc/claims"
	"tkestack.io/tke/pkg/util/log"
)

// TokenAuthenticator provides a function to verify token.
type TokenAuthenticator struct {
	IDTokenVerifier *oidc.IDTokenVerifier

	authClient authinternalclient.AuthInterface
}

// NewTokenAuthenticator creates new TokenAuthenticator object.
func NewTokenAuthenticator(authClient authinternalclient.AuthInterface) *TokenAuthenticator {
	return &TokenAuthenticator{authClient: authClient}
}

// AuthenticateToken verifies oidc token and returns user info.
func (h *TokenAuthenticator) AuthenticateToken(ctx context.Context, token string) (*genericauthenticator.Response, bool, error) {
	startTime := time.Now()
	defer func() {
		log.Debug("Finish verifying oidc bearer token", log.String("token", token), log.Duration("processTime", time.Since(startTime)))
	}()

	if h.IDTokenVerifier == nil {
		return nil, false, errors.New("Authenticator not initialized")
	}

	idToken, err := h.IDTokenVerifier.Verify(ctx, token)
	if err != nil {
		log.Error("Failed to verify the oidc bearer token", log.String("token", token), log.Err(err))
		return nil, false, err
	}
	var claims oidcclaims.IDTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		log.Error("Failed to unmarshal the id token", log.Any("idToken", idToken), log.Err(err))
		return nil, false, err
	}

	info := &user.DefaultInfo{Name: claims.Name}
	info.Groups = claims.Groups

	info.Extra = map[string][]string{}
	info.Extra[genericoidc.TenantIDKey] = []string{claims.FederatedIDClaims.ConnectorID}
	info.Extra["expireAt"] = []string{time.Unix(claims.Expiry, 0).String()}
	info.Extra["issueAt"] = []string{time.Unix(claims.IssuedAt, 0).String()}
	log.Debug("OIDC authenticateToken result", log.Any("user info", info))
	return &genericauthenticator.Response{User: info}, true, nil
}
