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

package authenticator

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"tkestack.io/tke/pkg/auth/util"

	genericapiserver "k8s.io/apiserver/pkg/server"

	"gopkg.in/square/go-jose.v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/server"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/pkiutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
)

type apiSigningKeysHookHandler struct {
	authClient authinternalclient.AuthInterface
}

// NewAPISigningKeyHookHandler creates a new apiSigningKeysHookHandler object.
func NewAPISigningKeyHookHandler(authClient authinternalclient.AuthInterface) genericapiserver.PostStartHookProvider {
	return &apiSigningKeysHookHandler{
		authClient: authClient,
	}
}

func (d *apiSigningKeysHookHandler) PostStartHook() (string, server.PostStartHookFunc, error) {
	return "generate-default-api-signing-keys", func(context server.PostStartHookContext) error {
		_, err := d.authClient.APISigningKeys().Get(util.DefaultAPISigningKey, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			privateKey, pubKey, err := generateKey()
			if err != nil {
				return err
			}

			signingKey := &auth.APISigningKey{
				ObjectMeta: metav1.ObjectMeta{
					Name: util.DefaultAPISigningKey,
				},
				SigningKey:    privateKey,
				SigningKeyPub: pubKey,
			}

			if _, err := d.authClient.APISigningKeys().Create(signingKey); err != nil {
				log.Error("Failed to create the api signing key", log.Err(err))
				return err
			}

			return nil
		}
		return err
	}, nil
}

func generateKey() ([]byte, []byte, error) {
	key, err := pkiutil.NewPrivateKey()
	if err != nil {
		log.Error("Failed generate signing key", log.Err(err))
		return nil, nil, err
	}

	b := make([]byte, 20)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	keyID := hex.EncodeToString(b)
	privateKey := &jose.JSONWebKey{
		Key:       key,
		KeyID:     keyID,
		Algorithm: "RS256",
		Use:       "sig",
	}
	pubKey := &jose.JSONWebKey{
		Key:       key.Public(),
		KeyID:     keyID,
		Algorithm: "RS256",
		Use:       "sig",
	}

	privateKeyBytes, err := privateKey.MarshalJSON()
	if err != nil {
		return nil, nil, err
	}

	pubKeyBytes, err := pubKey.MarshalJSON()
	if err != nil {
		return nil, nil, err
	}

	return privateKeyBytes, pubKeyBytes, nil
}
