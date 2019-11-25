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
	"gopkg.in/square/go-jose.v2"
	"testing"
	"time"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/etcd"

	"gotest.tools/assert"
)

var (
	store    Storage
	userName = "test"
	tenantID = "default"
)

// memoryKeyStorage stores SignKeys into memory, just for test.
type memoryKeyStorage struct {
	cache map[string]*types.SignKeys
}

func (e *memoryKeyStorage) CreateSignKey(method string, signKeys *types.SignKeys) error {
	if _, ok := e.cache[method]; ok {
		return etcd.ErrAlreadyExists
	}

	e.cache[method] = signKeys
	return nil
}

func (e *memoryKeyStorage) GetSignKey(method string) (*types.SignKeys, error) {
	v, ok := e.cache[method]
	if !ok {
		return nil, etcd.ErrNotFound
	}

	return v, nil
}

func (e *memoryKeyStorage) UpdateSignKey(method string, updater func(old types.SignKeys) (types.SignKeys, error)) error {
	panic("implement me")
}

func setup() {
	store = &memoryKeyStorage{cache: map[string]*types.SignKeys{}}
}

func TestRSASignAndVerify(t *testing.T) {
	setup()
	keyGen, _ := newKeyGen(string(jose.RS256), store)
	keySigner := genericKeySigner{
		keyGen: keyGen,
		method: string(jose.RS256),
	}

	apiKey, err := keySigner.Generate(userName, tenantID, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(apiKey.APIkey)
	info, err := keySigner.Verify(apiKey.APIkey)
	if err != nil {
		t.Fatal(err)
	}

	assert.Assert(t, info.TenantID == tenantID)
	assert.Assert(t, info.UserName == userName)
}

func TestRSASignAndVerifyFailed(t *testing.T) {
	setup()
	keyGen, _ := newKeyGen(string(jose.RS256), store)
	keySigner := genericKeySigner{
		keyGen: keyGen,
		method: string(jose.RS256),
	}

	str := ""
	_, err := keySigner.Verify(str)
	if err == nil {
		t.Fatal("expected verify failed, but success")
	}

	str = "123"
	_, err = keySigner.Verify(str)
	if err == nil {
		t.Fatal("expected verify failed, but success")
	}
}

func TestHMACSignAndVerify(t *testing.T) {
	setup()
	keyGen, _ := newKeyGen(string(jose.HS256), store)
	keySigner := genericKeySigner{
		keyGen: keyGen,
		method: string(jose.HS256),
	}

	apiKey, err := keySigner.Generate(userName, tenantID, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	info, err := keySigner.Verify(apiKey.APIkey)
	if err != nil {
		t.Fatal(err)
	}

	assert.Assert(t, info.TenantID == tenantID)
	assert.Assert(t, info.UserName == userName)
}

func TestHMACSignAndVerifyFailed(t *testing.T) {
	setup()
	keyGen, _ := newKeyGen(string(jose.HS256), store)
	keySigner := genericKeySigner{
		keyGen: keyGen,
		method: string(jose.HS256),
	}

	str := ""
	_, err := keySigner.Verify(str)
	if err == nil {
		t.Fatal("expected verify failed, but success")
	}

	str = "123456789"
	_, err = keySigner.Verify(str)
	if err == nil {
		t.Fatal("expected verify failed, but success")
	}
}
