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

package category

import (
	"context"
	"encoding/json"
	"path"
	"reflect"
	"time"

	"tkestack.io/tke/pkg/util/log"

	"tkestack.io/tke/pkg/auth/types"

	"github.com/coreos/etcd/clientv3"
	"tkestack.io/tke/pkg/util/etcd"
)

var (
	prefix = "/categories/"

	// defaultStorageTimeout will be applied to all storage's operations.
	defaultStorageTimeout = 5 * time.Second
)

// Storage is responsible for performing category crud actions onto etcd directly.
type Storage struct {
	db *clientv3.Client
}

// NewCategoryStorage create the category storage instance.
func NewCategoryStorage(db *clientv3.Client) *Storage {
	return &Storage{db}
}

// Create to create a new category in the etcd.
func (c *Storage) Create(category *types.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnCreate(ctx, c.db, path.Join(prefix, category.TenantID, category.Name), category)
}

// Get to get a category by a given id.
func (c *Storage) Get(tenantID, id string) (*types.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	var req types.Category
	if err := etcd.GetKey(ctx, c.db, path.Join(prefix, tenantID, id), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// List to get all existing categoryList for specify tenant.
func (c *Storage) List(tenantID string) (*types.CategoryList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	res, err := c.db.Get(ctx, path.Join(prefix, tenantID)+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	categoryList := types.CategoryList{}
	for _, v := range res.Kvs {
		var category types.Category
		if err := json.Unmarshal(v.Value, &category); err != nil {
			return nil, err
		}
		categoryList.Items = append(categoryList.Items, &category)
	}

	return &categoryList, nil
}

// ListAll to get all existing categoryList.
func (c *Storage) ListAll() (*types.CategoryList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	res, err := c.db.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	categoryList := types.CategoryList{}
	for _, v := range res.Kvs {
		var category types.Category
		if err := json.Unmarshal(v.Value, &category); err != nil {
			return nil, err
		}
		categoryList.Items = append(categoryList.Items, &category)
	}

	return &categoryList, nil
}

// Update to modify a existing category by given id.
func (c *Storage) Update(tenantID, id string, updater func(old types.Category) (types.Category, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()
	return etcd.TxnUpdate(ctx, c.db, path.Join(prefix, tenantID, id), func(currentValue []byte) ([]byte, error) {
		var current types.Category
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

// Delete to delete a existing category by given id.
func (c *Storage) Delete(tenantID, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	return etcd.DeleteKey(ctx, c.db, path.Join(prefix, tenantID, id))
}

// Load to load predefine categories into storage.
func (c *Storage) Load(tenantID string, categoryList []*types.Category) error {
	for _, cate := range categoryList {
		cate.TenantID = tenantID
		cateExist, err := c.Get(tenantID, cate.Name)
		if err != nil && err == etcd.ErrNotFound {
			log.Info("Category not exists, create it", log.String("name", cate.Name))
			cate.CreateAt = time.Now()
			cate.UpdateAt = cate.CreateAt

			err := c.Create(cate)
			if err != nil {
				log.Error("Create predefine category failed", log.Any("category", cate), log.Err(err))
				return err
			}

			continue
		} else if err != nil {
			log.Error("Get predefine category failed", log.Any("category", cate), log.Err(err))
			continue
		}

		isChanged := false
		if cate.Name != cateExist.Name || cate.DisplayName != cateExist.DisplayName || cate.Description != cateExist.Description {
			isChanged = true
		}

		for act, desc := range cate.Actions {
			if descExist, ok := cateExist.Actions[act]; !ok {
				cateExist.Actions[act] = desc
				isChanged = true
			} else {
				if !reflect.DeepEqual(desc, descExist) {
					cateExist.Actions[act] = desc
					isChanged = true
				}
			}
		}

		if isChanged {
			log.Info("Category has changed, update it", log.String("name", cate.Name))
			updater := func(current types.Category) (types.Category, error) {
				current.Name = cate.Name
				current.TenantID = cate.TenantID
				current.DisplayName = cate.DisplayName
				current.Description = cate.Description
				cate.Actions = cateExist.Actions
				current.UpdateAt = time.Now()
				return current, nil
			}

			err = c.Update(cate.TenantID, cate.Name, updater)
			if err != nil {
				log.Info("Update predefine category failed", log.Any("category", cate))
			}
		}
	}

	return nil
}
