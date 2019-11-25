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
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"io"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/pkiutil"
)

var (
	signMap = map[string]string{
		"RSA":  string(jose.RS256),
		"HMAC": string(jose.HS256),
	}
)

// Storage defines interface to operate sign keys.
type Storage interface {
	CreateSignKey(method string, signKeys *types.SignKeys) error
	GetSignKey(method string) (*types.SignKeys, error)
	UpdateSignKey(method string, updater func(old types.SignKeys) (types.SignKeys, error)) error
}

// KeyGen contains basic method to operate key for signing api token.
type KeyGen interface {
	Get() (*types.SignKeys, error)
	generateKeys() error
}

func newKeyGen(method string, storage Storage) (KeyGen, error) {
	var keyGen = &genericKeygen{store: storage, method: method}

	keys, err := storage.GetSignKey(method)
	if err != nil && err != etcd.ErrNotFound {
		log.Error("Get api token key failed", log.Err(err))
		return nil, err
	}

	// if keys not exits, generate and store it
	if err != nil && err == etcd.ErrNotFound {
		log.Info("Api token sign key not exist, generate it")
		err := keyGen.generateKeys()
		if err != nil {
			log.Error("generate keys failed: failed", log.Err(err))
			return nil, fmt.Errorf("generate keys failed: %v", err)
		}
	}

	keyGen.keys = keys
	return keyGen, nil
}

type genericKeygen struct {
	store  Storage
	keys   *types.SignKeys
	method string
}

func (r *genericKeygen) Get() (*types.SignKeys, error) {
	if r.keys != nil {
		return r.keys, nil
	}

	if err := r.generateKeys(); err != nil {
		return nil, err
	}

	return r.keys, nil
}

func (r *genericKeygen) generateKeys() error {
	var (
		privKey interface{}
		pubKey  interface{}
		err     error
	)

	b := make([]byte, 20)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return err
	}

	keyID := hex.EncodeToString(b)
	switch r.method {
	case string(jose.RS256):
		privKey, err = pkiutil.NewPrivateKey()
		if err != nil {
			return err
		}
		pubKey = privKey.(*rsa.PrivateKey).Public()
	case string(jose.HS256):
		privKey = b
		pubKey = b
	default:
		return fmt.Errorf("unknown key sign method %s", r.method)
	}

	priv := &jose.JSONWebKey{
		Key:       privKey,
		KeyID:     keyID,
		Algorithm: "RS256",
		Use:       "sig",
	}
	pub := &jose.JSONWebKey{
		Key:       pubKey,
		KeyID:     keyID,
		Algorithm: "RS256",
		Use:       "sig",
	}

	keyData := types.SignKeys{
		SigningKey:    priv,
		SigningKeyPub: pub,
	}

	err = r.store.CreateSignKey(r.method, &keyData)
	if err != nil && err != etcd.ErrAlreadyExists {
		return err
	}

	r.keys = &types.SignKeys{}
	r.keys, err = r.store.GetSignKey(r.method)
	if err != nil {
		return err
	}

	return nil
}
