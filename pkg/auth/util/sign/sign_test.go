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

package sign

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"
	"time"

	"gopkg.in/square/go-jose.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"tkestack.io/tke/api/auth"
	fakeauth "tkestack.io/tke/api/client/clientset/internalversion/fake"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/pkiutil"
)

var (
	userName = "test"
	tenantID = "default"
	expire   = time.Hour
)

func TestGenericKeySigner(t *testing.T) {
	authClient := fakeauth.NewSimpleClientset()

	priv, pub, _ := generateKey()

	_, err := authClient.Auth().APISigningKeys().Create(&auth.APISigningKey{
		ObjectMeta: metav1.ObjectMeta{
			Name: DefaultAPISigningKey},
		SigningKey:    priv,
		SigningKeyPub: pub})

	signer := NewGenericKeySigner(authClient.Auth())
	apiKey, err := signer.Generate(userName, tenantID, expire)
	if err != nil {
		t.Fatal(err)
	}

	_, err = signer.Verify(apiKey.Spec.APIkey)
	if err != nil {
		t.Fatal(err)
	}
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
