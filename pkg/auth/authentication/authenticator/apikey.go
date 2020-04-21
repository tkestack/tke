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
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	genericauthenticator "k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
)

// APIKeyAuthenticator provides a function to generate and verify jwt-format api key.
type APIKeyAuthenticator struct {
	authClient authinternalclient.AuthInterface
	keySigner  util.KeySigner
}

// NewAPIKeyAuthenticator creates new APIKeyAuthenticator object.
func NewAPIKeyAuthenticator(authClient authinternalclient.AuthInterface) (*APIKeyAuthenticator, error) {
	keySigner := util.NewGenericKeySigner(authClient)
	apiKeyAuth := &APIKeyAuthenticator{authClient, keySigner}
	return apiKeyAuth, nil
}

// AuthenticateToken verifies jwt-format api key and returns user info.
func (h *APIKeyAuthenticator) AuthenticateToken(ctx context.Context, token string) (*genericauthenticator.Response, bool, error) {
	startTime := time.Now()
	defer func() {
		log.Debug("Finish verifying api key", log.String("api key", token), log.Duration("processTime", time.Since(startTime)))
	}()

	tokenInfo, err := h.keySigner.Verify(token)
	if err != nil {
		return nil, false, err
	}

	selector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.tenantID", tokenInfo.TenantID),
		fields.OneTermEqualSelector("spec.apiKey", token))

	apiKeyList, err := h.authClient.APIKeys().List(metav1.ListOptions{FieldSelector: selector.String()})
	if err != nil {
		log.Error("List api keys failed", log.String("api key", token), log.Err(err))
		return nil, false, err
	}

	if len(apiKeyList.Items) == 0 {
		log.Error("Api key is verified, but not found in storage", log.String("api key", token))
		return nil, false, fmt.Errorf("api key has been deleted")
	}

	apiKey := apiKeyList.Items[0]
	if apiKey.Status.Disabled {
		log.Info("Api key has been disabled or deleted", log.String("api key", token))
		return nil, false, fmt.Errorf("api key has been disabled")
	}

	info := &user.DefaultInfo{Name: tokenInfo.UserName}

	localIdentity, err := util.GetLocalIdentity(h.authClient, tokenInfo.TenantID, info.Name)
	if err != nil {
		log.Error("Get localIdentity failed", log.String("localIdentity", info.Name), log.Err(err))
		return nil, false, err
	}

	info.UID = localIdentity.ObjectMeta.Name
	groups, err := util.GetGroupsForUser(h.authClient, localIdentity.ObjectMeta.Name)
	if err == nil {
		for _, g := range groups.Items {
			info.Groups = append(info.Groups, g.ObjectMeta.Name)
		}
	}

	info.Extra = map[string][]string{}
	info.Extra[genericoidc.TenantIDKey] = []string{tokenInfo.TenantID}
	info.Extra["expireAt"] = []string{time.Unix(tokenInfo.ExpiresAt, 0).String()}
	info.Extra["issueAt"] = []string{time.Unix(tokenInfo.IssuedAt, 0).String()}
	info.Extra["description"] = []string{apiKey.Spec.Description}

	log.Debug("APIkey authenticateToken result", log.Any("user info", info))
	return &genericauthenticator.Response{User: info}, true, nil
}
