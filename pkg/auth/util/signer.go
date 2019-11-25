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

package util

import (
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"time"
	"tkestack.io/tke/pkg/auth/types"

	"github.com/dgrijalva/jwt-go"
	"tkestack.io/tke/pkg/util/log"
)

var (
	minExpire = 1 * time.Second
	maxExpire = 100 * 365 * 24 * time.Hour
)

// APIClaims is the claims section of jwt token.
type APIClaims struct {
	*jwt.StandardClaims
	*KeyData
}

// KeyData contains the necessary info of api key validated.
type KeyData struct {
	UserName string `json:"usr,omitempty"`
	TenantID string `json:"ted,omitempty"`
}

// KeySigner is a interface used to generate api key for a user
type KeySigner interface {
	Generate(userName string, tenantID string, expire time.Duration) (*types.APIKeyData, error)
	Verify(apiKey string) (*APIClaims, error)
}

type genericKeySigner struct {
	keyGen KeyGen
	method string
}

// NewGenericKeySigner creates a generic key signer instance.
func NewGenericKeySigner(method string, store Storage) (KeySigner, error) {
	signMethod, ok := signMap[method]
	if !ok {
		return nil, fmt.Errorf("unknown sign method, only support RSA and HMAC")
	}

	keyGetter, err := newKeyGen(signMethod, store)
	if err != nil {
		return nil, fmt.Errorf("create sign key generator failed %v", err)
	}

	return &genericKeySigner{keyGen: keyGetter, method: signMethod}, nil
}

// Generate use generate api key from username and tenantID.
func (g *genericKeySigner) Generate(userName string, tenantID string, expire time.Duration) (*types.APIKeyData, error) {
	now := time.Now()
	if expire < minExpire || expire > maxExpire {
		return nil, fmt.Errorf("expire %v must not shorter than %v or longer than %v", expire, minExpire, maxExpire)
	}

	// jti ensures api key unique in case generating in the same time.
	bytes, err := uuid.GenerateRandomBytes(2)
	if err != nil {
		return nil, err
	}
	jti := hex.EncodeToString(bytes)
	claims := jwt.NewWithClaims(jwt.GetSigningMethod(g.method), &APIClaims{
		&jwt.StandardClaims{
			// set the issue time
			// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.6
			IssuedAt: now.Unix(),
			// set the expire time
			// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
			ExpiresAt: now.Add(expire).Unix(),
			// set the jwt id
			// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.7
			Id: jti,
		},
		&KeyData{
			UserName: userName,
			TenantID: tenantID,
		},
	})

	keys, err := g.keyGen.Get()
	if err != nil {
		return nil, err
	}

	apiKey, err := claims.SignedString(keys.SigningKey.Key)
	if err != nil {
		return nil, err
	}

	disabled := false
	deleted := false
	return &types.APIKeyData{
		APIkey:   apiKey,
		IssueAt:  now,
		ExpireAt: now.Add(expire),
		Disabled: &disabled,
		Deleted:  &deleted,
	}, nil
}

// Verify verifies api key and returns key info.
func (g *genericKeySigner) Verify(apiKey string) (*APIClaims, error) {
	keys, err := g.keyGen.Get()
	if err != nil {
		return nil, err
	}

	result, err := jwt.ParseWithClaims(apiKey, &APIClaims{}, func(token *jwt.Token) (interface{}, error) {
		return keys.SigningKeyPub.Key, nil
	})

	var (
		claims *APIClaims
		ok     bool
	)

	if result != nil {
		claims, ok = result.Claims.(*APIClaims)
		if !ok {
			return nil, fmt.Errorf("invalid format api key info: %+v", result.Claims)
		}

		if result.Valid {
			return claims, nil
		}
	}

	if err != nil {
		log.Warn("Verify api key failed", log.String("api key", substring(apiKey)), log.Any("result", result), log.Err(err))
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return claims, fmt.Errorf("not valid jwt token format")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return claims, fmt.Errorf("token is either expired or not active yet")
			} else {
				return claims, fmt.Errorf("token is invalid")
			}
		} else {
			return claims, fmt.Errorf("token is invalid")
		}
	}

	return claims, fmt.Errorf("token is invalid")
}

func substring(str string) string {
	return str[:minimum(20, len(str))]
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}
