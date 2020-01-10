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

package util

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-uuid"
	"gopkg.in/square/go-jose.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/util/log"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

var (
	DefaultAPISigningKey = "default-api-signing-key"
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
	Generate(username string, tenantID string, expire time.Duration) (*auth.APIKey, error)
	Verify(apiKey string) (*APIClaims, error)
}

// NewGenericKeySigner creates a generic key signer instance.
func NewGenericKeySigner(authclient authinternalclient.AuthInterface) KeySigner {
	return genericKeySigner{authClient: authclient}
}

type genericKeySigner struct {
	authClient authinternalclient.AuthInterface
}

// Generate use generate api key from username and tenantID.
func (g genericKeySigner) Generate(username string, tenantID string, expire time.Duration) (*auth.APIKey, error) {
	now := time.Now()

	// jti ensures api key unique in case generating in the same time.
	bytes, err := uuid.GenerateRandomBytes(2)
	if err != nil {
		return nil, err
	}
	jti := hex.EncodeToString(bytes)
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, &APIClaims{
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
			UserName: username,
			TenantID: tenantID,
		},
	})

	keys, err := g.authClient.APISigningKeys().Get(DefaultAPISigningKey, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get signing keys", log.Err(err))
		return nil, err
	}

	var privKey jose.JSONWebKey
	if err := privKey.UnmarshalJSON(keys.SigningKey); err != nil {
		log.Error("Failed to unmarshal signing priv key", log.ByteString("signingKey", keys.SigningKey), log.Err(err))
		return nil, err
	}

	apiKeyStr, err := claims.SignedString(privKey.Key)
	if err != nil {
		log.Error("Failed to sign claims", log.Err(err))
		return nil, err
	}

	apiKey := &auth.APIKey{
		Spec: auth.APIKeySpec{
			APIkey:   apiKeyStr,
			TenantID: tenantID,
			IssueAt:  metav1.NewTime(now),
			Username: username,
			ExpireAt: metav1.NewTime(now.Add(expire)),
		},
	}

	return apiKey, nil
}

func (g genericKeySigner) Verify(apiKey string) (*APIClaims, error) {
	keys, err := g.authClient.APISigningKeys().Get(DefaultAPISigningKey, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get signing keys", log.Err(err))
		return nil, err
	}

	var pubKey jose.JSONWebKey
	if err := pubKey.UnmarshalJSON(keys.SigningKeyPub); err != nil {
		log.Error("Failed to unmarshal signing pub key", log.ByteString("signingPubKey", keys.SigningKeyPub), log.Err(err))
		return nil, err
	}

	result, err := jwt.ParseWithClaims(apiKey, &APIClaims{}, func(token *jwt.Token) (interface{}, error) {
		return pubKey.Key, nil
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
		log.Warn("Verify api key failed", log.String("api key", apiKey), log.Any("result", result), log.Err(err))
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
