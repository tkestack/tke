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

package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

var (
	// ErrNotFound is the error returned by storages if a resource cannot be found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is the error returned by storages if a resource ID is taken during a create.
	ErrAlreadyExists = errors.New("ID already exists")
)

// TxnCreate creates a new object into etcd with given key and value and ensure transaction consistency.
func TxnCreate(ctx context.Context, db *clientv3.Client, key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	txn := db.Txn(ctx)
	res, err := txn.
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, string(b))).
		Commit()
	if err != nil {
		return err
	}

	if !res.Succeeded {
		return ErrAlreadyExists
	}
	return nil
}

// TxnUpdate updates a existing object into etcd with given key and update function and ensure transaction consistency.
func TxnUpdate(ctx context.Context, db *clientv3.Client, key string, update func(current []byte) ([]byte, error)) error {
	getResp, err := db.Get(ctx, key)
	if err != nil {
		return err
	}
	var currentValue []byte
	var modRev int64
	if len(getResp.Kvs) > 0 {
		currentValue = getResp.Kvs[0].Value
		modRev = getResp.Kvs[0].ModRevision
	}

	updatedValue, err := update(currentValue)
	if err != nil {
		return err
	}

	txn := db.Txn(ctx)
	updateResp, err := txn.
		If(clientv3.Compare(clientv3.ModRevision(key), "=", modRev)).
		Then(clientv3.OpPut(key, string(updatedValue))).
		Commit()
	if err != nil {
		return err
	}
	if !updateResp.Succeeded {
		return fmt.Errorf("failed to update key=%q: concurrent conflicting update happened", key)
	}
	return nil
}

// GetKey get a existing object from etcd with given key
func GetKey(ctx context.Context, db *clientv3.Client, key string, value interface{}) error {
	r, err := db.Get(ctx, key)
	if err != nil {
		return err
	}
	if r.Count == 0 {
		return ErrNotFound
	}
	return json.Unmarshal(r.Kvs[0].Value, value)
}

// DeleteKey delete a existing object from etcd with given key
func DeleteKey(ctx context.Context, db *clientv3.Client, key string) error {
	res, err := db.Delete(ctx, key)
	if err != nil {
		return err
	}
	if res.Deleted == 0 {
		return ErrNotFound
	}

	return nil
}
