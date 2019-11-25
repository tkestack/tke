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
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"k8s.io/apimachinery/pkg/util/wait"
	genericauthenticator "k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"

	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/auth/registry/apikey"
	"tkestack.io/tke/pkg/auth/registry/localidentity"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
)

var (
	apiKeyRotateInterval = 30 * time.Minute

	defaultAPIKeyTimeout = 7 * 24 * time.Hour
)

// APIKeyAuthenticator provides a function to generate and verify jwt-format api key.
type APIKeyAuthenticator struct {
	store         apikey.Storage
	identityStore *localidentity.Storage
	keySigner     util.KeySigner
}

// NewAPIKeyAuthenticator creates new APIKeyAuthenticator object.
func NewAPIKeyAuthenticator(method string, store apikey.Storage, identityStore *localidentity.Storage) (*APIKeyAuthenticator, error) {
	apiKeyAuth := &APIKeyAuthenticator{store: store, identityStore: identityStore}
	var err error
	apiKeyAuth.keySigner, err = util.NewGenericKeySigner(method, store)
	if err != nil {
		return nil, err
	}

	apiKeyAuth.StartRotateAPIkey()
	return apiKeyAuth, nil
}

// AuthenticateToken verifies jwt-format api key and returns user info.
func (h *APIKeyAuthenticator) AuthenticateToken(ctx context.Context, token string) (*genericauthenticator.Response, bool, error) {
	startTime := time.Now()
	defer func() {
		log.Info("Finish verifying api key", log.String("api key", token), log.Duration("processTime", time.Since(startTime)))
	}()

	tokenInfo, err := h.keySigner.Verify(token)
	if err != nil {
		return nil, false, err
	}

	keyData, err := h.store.GetAPIKey(tokenInfo.TenantID, tokenInfo.UserName, token)
	if err != nil {
		if err != etcd.ErrNotFound {
			log.Error("Get key data failed", log.Err(err))
			return nil, false, err
		}

		log.Warn("Orphaned api key")
	}

	if keyData != nil && (*keyData.Disabled || *keyData.Deleted) {
		log.Info("Api key has been disabled or deleted")
		return nil, false, fmt.Errorf("api key has been disabled or deleted")
	}

	info := &user.DefaultInfo{Name: tokenInfo.UserName}
	info.Extra = map[string][]string{}
	info.Extra[genericoidc.TenantIDKey] = []string{tokenInfo.TenantID}
	info.Extra["expireAt"] = []string{time.Unix(tokenInfo.ExpiresAt, 0).String()}
	info.Extra["issueAt"] = []string{time.Unix(tokenInfo.IssuedAt, 0).String()}
	info.Extra["description"] = []string{keyData.Description}

	return &genericauthenticator.Response{User: info}, true, nil
}

// CreateToken creates a new api key with user info and expiration.
func (h *APIKeyAuthenticator) CreateToken(tenantID, userName, description string, expire time.Duration) (*types.APIKeyData, error) {
	startTime := time.Now()
	defer func() {
		log.Info("Finish generating api key", log.String("userName", userName), log.String("tenantID", tenantID),
			log.Duration("expire", expire), log.Duration("processTime", time.Since(startTime)))
	}()

	if expire == 0 {
		expire = defaultAPIKeyTimeout
	}
	keyData, err := h.keySigner.Generate(userName, tenantID, expire)
	if err != nil {
		log.Error("Generate api key failed", log.Err(err))
		return nil, err
	}

	newKeyData := *keyData
	newKeyData.Description = description
	newKeyData.APIkey = ""
	if err := h.store.CreateAPIKey(tenantID, userName, keyData.APIkey, &newKeyData); err != nil {
		log.Error("Insert APIKey data failed", log.Err(err))
		return nil, err
	}

	return keyData, nil
}

// CreateTokenWithPassword creates a new api key with user password and expiration.
func (h *APIKeyAuthenticator) CreateTokenWithPassword(tenantID, userName, password, description string, expire time.Duration) (*types.APIKeyData, error) {
	startTime := time.Now()
	defer func() {
		log.Info("Finish generating api key", log.String("userName", userName), log.String("tenantID", tenantID),
			log.Duration("expire", expire), log.Duration("processTime", time.Since(startTime)))
	}()

	if expire == 0 {
		expire = defaultAPIKeyTimeout
	}

	localIdentity, err := h.identityStore.Get(tenantID, userName)
	if err != nil {
		log.Error("Get user failed", log.String("user", userName), log.Err(err))
		return nil, err
	}

	bytes, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		log.Error("Decode password failed", log.Err(err))
		return nil, err
	}

	hashBytes, err := base64.StdEncoding.DecodeString(localIdentity.Spec.HashedPassword)
	if err != nil {
		log.Error("Parse hash password failed", log.String("hashedPassword", localIdentity.Spec.HashedPassword), log.Err(err))
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(hashBytes, bytes); err != nil {
		log.Error("Invalid password", log.ByteString("input password", []byte(password)), log.ByteString("store password", hashBytes))
		return nil, fmt.Errorf("password is wrong")
	}

	keyData, err := h.keySigner.Generate(userName, tenantID, expire)
	if err != nil {
		log.Error("Generate api key failed", log.Err(err))
		return nil, err
	}

	newKeyData := *keyData
	newKeyData.Description = description
	newKeyData.APIkey = ""
	if err := h.store.CreateAPIKey(tenantID, userName, keyData.APIkey, &newKeyData); err != nil {
		log.Error("Insert APIKey data failed", log.Err(err))
		return nil, err
	}

	return keyData, nil
}

// UpdateToken used to disable or delete api key.
func (h *APIKeyAuthenticator) UpdateToken(keyReq *types.APIKeyData, tenantID, userName string) error {
	startTime := time.Now()
	defer func() {
		log.Info("Finish update api key for user", log.String("userName", userName), log.String("tenantID", tenantID),
			log.Duration("processTime", time.Since(startTime)))
	}()

	tokenInfo, err := h.keySigner.Verify(keyReq.APIkey)
	if err != nil && tokenInfo == nil {
		return err
	}

	if tokenInfo.TenantID != tenantID || tokenInfo.UserName != userName {
		log.Warnf("User (%s,%s) try to operate api key %s belong to (%s, %s)", tenantID, userName, keyReq.APIkey, tokenInfo.TenantID, tokenInfo.UserName)
		return fmt.Errorf("API key only can be disabled or deleted by owner")
	}

	updater := func(current types.APIKeyData) (types.APIKeyData, error) {
		if keyReq.Deleted != nil {
			// if delete key, disable key
			if *keyReq.Deleted {
				*current.Disabled = true
			} else {
				*current.Disabled = false
			}

			current.Deleted = keyReq.Deleted
		}

		if !*current.Deleted && keyReq.Disabled != nil {
			*current.Disabled = *keyReq.Disabled
		}

		current.Description = keyReq.Description
		return current, nil
	}

	return h.store.UpdateAPIKey(tokenInfo.TenantID, tokenInfo.UserName, keyReq.APIkey, updater)
}

// ListAPIKeys will list a whole list for a given user.
func (h *APIKeyAuthenticator) ListAPIKeys(tenantID, userName string) (*types.APIKeyList, error) {
	startTime := time.Now()
	defer func() {
		log.Info("Finish list api keys for user", log.String("userName", userName), log.String("tenantID", tenantID),
			log.Duration("processTime", time.Since(startTime)))
	}()

	keyList, err := h.store.ListAPIKeys(tenantID, userName)
	if err != nil {
		log.Error("List all api keys for user failed", log.Err(err))
		return nil, err
	}

	return keyList, nil
}

// StartRotateAPIkey clears all expired api token which has been disabled before.
func (h *APIKeyAuthenticator) StartRotateAPIkey() {
	go wait.Until(h.store.RotateAPIKeys, apiKeyRotateInterval, wait.NeverStop)
}
