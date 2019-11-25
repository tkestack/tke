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

package policy

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
	prefix = "/policies/"

	// defaultStorageTimeout will be applied to all storage's operations.
	defaultStorageTimeout = 5 * time.Second
)

// Storage is responsible for performing policy crud actions onto etcd directly.
type Storage struct {
	db *clientv3.Client
}

// NewPolicyStorage create the policy storage instance.
func NewPolicyStorage(db *clientv3.Client) *Storage {
	return &Storage{db}
}

// Create to create a new policy in the etcd.
func (c *Storage) Create(policy *types.Policy) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnCreate(ctx, c.db, path.Join(prefix, policy.TenantID, policy.ID), policy)
}

// Get to get a policy by by a given id.
func (c *Storage) Get(tenantID, id string) (*types.Policy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	var req types.Policy
	if err := etcd.GetKey(ctx, c.db, path.Join(prefix, tenantID, id), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// List to get all existing policies for specify tenant.
func (c *Storage) List(tenantID string) (*types.PolicyList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	res, err := c.db.Get(ctx, path.Join(prefix, tenantID)+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	policies := types.PolicyList{}
	for _, v := range res.Kvs {
		var policy types.Policy
		if err := json.Unmarshal(v.Value, &policy); err != nil {
			return nil, err
		}
		policies.Items = append(policies.Items, &policy)
	}

	return &policies, nil
}

// ListAll to get all existing policies.
func (c *Storage) ListAll() (*types.PolicyList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	res, err := c.db.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	policies := types.PolicyList{}
	for _, v := range res.Kvs {
		var policy types.Policy
		if err := json.Unmarshal(v.Value, &policy); err != nil {
			return nil, err
		}
		policies.Items = append(policies.Items, &policy)
	}

	return &policies, nil
}

// Update to modify a existing policy by given id.
func (c *Storage) Update(tenantID, id string, updater func(old types.Policy) (types.Policy, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnUpdate(ctx, c.db, path.Join(prefix, tenantID, id), func(currentValue []byte) ([]byte, error) {
		var current types.Policy
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

// Delete to delete a existing policy by given id.
func (c *Storage) Delete(tenantID, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	return etcd.DeleteKey(ctx, c.db, path.Join(prefix, tenantID, id))
}
