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

package apikey

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"path"
	"strings"
	"time"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/etcd"
	"tkestack.io/tke/pkg/util/log"
)

var (
	apiKeyPrefix = "/apikeys/"

	signKeyPrefix = "/signkeys/"

	defaultStorageTimeout = 5 * time.Second

	apiKeyRotateInterval = 30 * time.Second
)

// Storage is responsible for performing category crud actions onto etcd directly.
type Storage interface {
	CreateSignKey(method string, signKeys *types.SignKeys) error
	GetSignKey(method string) (*types.SignKeys, error)
	UpdateSignKey(method string, updater func(old types.SignKeys) (types.SignKeys, error)) error

	CreateAPIKey(tenantID, userName, apiKey string, data *types.APIKeyData) error
	GetAPIKey(tenantID, userName, apiKey string) (*types.APIKeyData, error)
	ListAPIKeys(tenantID, userName string) (*types.APIKeyList, error)
	ListAllAPIKeys() (*types.APIKeyList, error)
	UpdateAPIKey(tenantID, userName, apiKey string, updater func(old types.APIKeyData) (types.APIKeyData, error)) error
	DeleteAPIKey(tenantID, userName, apiKey string) error
	RotateAPIKeys()
}

type etcdStorage struct {
	db *clientv3.Client
}

// NewAPIKeyStorage creates the apiKey storage instance.
func NewAPIKeyStorage(db *clientv3.Client) Storage {
	return &etcdStorage{db}
}

// CreateSignKey to create a signKeys for the method in the etcd.
func (e *etcdStorage) CreateSignKey(method string, signKeys *types.SignKeys) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnCreate(ctx, e.db, path.Join(signKeyPrefix, method), signKeys)
}

// CreateSignKey to create a signKeys for the method in the etcd.
func (e *etcdStorage) GetSignKey(method string) (*types.SignKeys, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	var req types.SignKeys
	if err := etcd.GetKey(ctx, e.db, path.Join(signKeyPrefix, method), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// UpdateSignKey to modify a existing signkeys by given method.
func (e *etcdStorage) UpdateSignKey(method string, updater func(old types.SignKeys) (types.SignKeys, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnUpdate(ctx, e.db, path.Join(signKeyPrefix, method), func(currentValue []byte) ([]byte, error) {
		var current types.SignKeys
		if len(currentValue) > 0 {
			if err := json.Unmarshal(currentValue, &current); err != nil {
				return nil, err
			}
		}
		updated, err := updater(current)
		if err != nil {
			return nil, err
		}
		return json.Marshal(updated)
	})
}

// CreateAPIKey to create a apikey for the method in the etcd.
func (e *etcdStorage) CreateAPIKey(tenantID, userName, apiKey string, signKeys *types.APIKeyData) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnCreate(ctx, e.db, path.Join(apiKeyPrefix, tenantID, userName, apiKey), signKeys)
}

func (e *etcdStorage) GetAPIKey(tenantID, userName, apiKey string) (*types.APIKeyData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	var req types.APIKeyData
	if err := etcd.GetKey(ctx, e.db, path.Join(apiKeyPrefix, tenantID, userName, apiKey), &req); err != nil {
		return nil, err
	}

	req.APIkey = apiKey
	return &req, nil
}

// ListAPIKeys to get all existing api key list for given user.
func (e *etcdStorage) ListAPIKeys(tenantID, userName string) (*types.APIKeyList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	keyDir := path.Join(apiKeyPrefix, tenantID, userName) + "/"
	res, err := e.db.Get(ctx, keyDir, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	keyList := types.APIKeyList{}
	now := time.Now()
	for _, v := range res.Kvs {
		var keyData types.APIKeyData
		if err := json.Unmarshal(v.Value, &keyData); err != nil {
			return nil, err
		}

		if keyData.Deleted == nil || !*keyData.Deleted {
			keyData.APIkey = strings.TrimPrefix(string(v.Key), keyDir)

			if keyData.ExpireAt.Before(now) {
				keyData.Expired = true
			}
			keyList.Items = append(keyList.Items, &keyData)
		}
	}

	return &keyList, nil
}

// ListAPIKeys to get all existing api key list for given user.
func (e *etcdStorage) ListAllAPIKeys() (*types.APIKeyList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	res, err := e.db.Get(ctx, apiKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	keyList := types.APIKeyList{}
	for _, v := range res.Kvs {
		var keyData types.APIKeyData
		if err := json.Unmarshal(v.Value, &keyData); err != nil {
			return nil, err
		}
		keyData.APIkey = string(v.Key)
		keyList.Items = append(keyList.Items, &keyData)
	}

	return &keyList, nil
}

// UpdateAPIKey to modify a existing apikey info.
func (e *etcdStorage) UpdateAPIKey(tenantID, userName, apiKey string, updater func(old types.APIKeyData) (types.APIKeyData, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnUpdate(ctx, e.db, path.Join(apiKeyPrefix, tenantID, userName, apiKey), func(currentValue []byte) ([]byte, error) {
		var current types.APIKeyData
		if len(currentValue) > 0 {
			if err := json.Unmarshal(currentValue, &current); err != nil {
				return nil, err
			}
		}
		updated, err := updater(current)
		if err != nil {
			return nil, err
		}
		return json.Marshal(updated)
	})
}

// DeleteAPIKey to delete a existing api key.
func (e *etcdStorage) DeleteAPIKey(tenantID, userName, apiKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	return etcd.DeleteKey(ctx, e.db, path.Join(apiKeyPrefix, tenantID, userName, apiKey))
}

func (e *etcdStorage) RotateAPIKeys() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	res, err := e.db.Get(ctx, apiKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		log.Error("List all api keys failed", log.Err(err))
		return
	}

	now := time.Now()
	for _, v := range res.Kvs {
		var keyData types.APIKeyData
		if err := json.Unmarshal(v.Value, &keyData); err != nil {
			log.Warn("Unmarshal api key data failed", log.Err(err))
			continue
		}

		if keyData.Deleted != nil && *keyData.Deleted && keyData.ExpireAt.Before(now) {
			if err := etcd.DeleteKey(context.Background(), e.db, string(v.Key)); err != nil {
				log.Warn("Remove deleted  expired api key failed", log.Err(err))
				continue
			}
		}
	}
}
