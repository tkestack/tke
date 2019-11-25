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

package localidentity

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/etcd"

	"github.com/coreos/etcd/clientv3"
)

var (
	prefix = "/localidentities/"

	// defaultStorageTimeout will be applied to all storage's operations.
	defaultStorageTimeout = 5 * time.Second
)

// Storage is responsible for performing local identity crud actions onto etcd directly.
type Storage struct {
	db *clientv3.Client
}

// NewLocalIdentity create the localidentity storage instance.
func NewLocalIdentity(db *clientv3.Client) *Storage {
	return &Storage{db}
}

// Create to create a new user in the etcd.
func (c *Storage) Create(identity *types.LocalIdentity) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnCreate(ctx, c.db, path.Join(prefix, identity.Spec.TenantID, identity.Name), identity)
}

// Get to get a user by by a given name.
func (c *Storage) Get(tenantID string, name string) (*types.LocalIdentity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	var req types.LocalIdentity
	if err := etcd.GetKey(ctx, c.db, path.Join(prefix, tenantID, name), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// List to get all existing users.
func (c *Storage) List(tenantID string) (*types.LocalIdentityList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	res, err := c.db.Get(ctx, path.Join(prefix, tenantID)+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	identities := types.LocalIdentityList{}
	for _, v := range res.Kvs {
		var localIdentity types.LocalIdentity
		if err := json.Unmarshal(v.Value, &localIdentity); err != nil {
			return nil, err
		}
		identities.Items = append(identities.Items, &localIdentity)
	}

	return &identities, nil
}

// Update to modify a existing user by given name.
func (c *Storage) Update(tenantID string, name string, updater func(old types.LocalIdentity) (types.LocalIdentity, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnUpdate(ctx, c.db, path.Join(prefix, tenantID, name), func(currentValue []byte) ([]byte, error) {
		var current types.LocalIdentity
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

// Delete to delete a existing user by given name.
func (c *Storage) Delete(tenantID string, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	return etcd.DeleteKey(ctx, c.db, path.Join(prefix, tenantID, name))
}
