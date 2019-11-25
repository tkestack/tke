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

package role

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"tkestack.io/tke/pkg/auth/types"

	"github.com/coreos/etcd/clientv3"
	"tkestack.io/tke/pkg/util/etcd"
)

var (
	prefix = "/roles/"

	// defaultStorageTimeout will be applied to all storage's operations.
	defaultStorageTimeout = 5 * time.Second
)

// Storage is responsible for performing role crud actions onto etcd.
type Storage struct {
	db *clientv3.Client
}

// NewRoleStorage create the role storage instance.
func NewRoleStorage(db *clientv3.Client) *Storage {
	return &Storage{db}
}

// Create to create a new role into the etcd.
func (c *Storage) Create(role *types.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnCreate(ctx, c.db, path.Join(prefix, role.TenantID, role.ID), role)
}

// Get to get a role by by a given id.
func (c *Storage) Get(tenantID, id string) (*types.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	var req types.Role
	if err := etcd.GetKey(ctx, c.db, path.Join(prefix, tenantID, id), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// List to get all existing roles for a specify tenant.
func (c *Storage) List(tenantID string) (*types.RoleList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	res, err := c.db.Get(ctx, path.Join(prefix, tenantID)+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	roles := types.RoleList{}
	for _, v := range res.Kvs {
		var role types.Role
		if err := json.Unmarshal(v.Value, &role); err != nil {
			return nil, err
		}
		roles.Items = append(roles.Items, &role)
	}

	return &roles, nil
}

// ListAll to get all existing roles.
func (c *Storage) ListAll() (*types.RoleList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	res, err := c.db.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	roles := types.RoleList{}
	for _, v := range res.Kvs {
		var role types.Role
		if err := json.Unmarshal(v.Value, &role); err != nil {
			return nil, err
		}
		roles.Items = append(roles.Items, &role)
	}

	return &roles, nil
}

// Update to modify a existing role by given id.
func (c *Storage) Update(tenantID, id string, updater func(old types.Role) (types.Role, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnUpdate(ctx, c.db, path.Join(prefix, tenantID, id), func(currentValue []byte) ([]byte, error) {
		var current types.Role
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

// Delete to delete a existing role by given id.
func (c *Storage) Delete(tenantID, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	return etcd.DeleteKey(ctx, c.db, path.Join(prefix, tenantID, id))
}
