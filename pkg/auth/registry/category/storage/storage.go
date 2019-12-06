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

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"tkestack.io/tke/api/auth"
	"tkestack.io/tke/pkg/auth/registry/category"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for signing keys and all sub resources.
type Storage struct {
	*REST
}

// NewStorage returns a Storage object that will work against signing key.
func NewStorage(optsGetter generic.RESTOptionsGetter) *Storage {
	strategy := category.NewStrategy()
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &auth.Category{} },
		NewListFunc:              func() runtime.Object { return &auth.CategoryList{} },
		DefaultQualifiedResource: auth.Resource("categories"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    category.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create category etcd rest storage", log.Err(err))
	}

	return &Storage{
		&REST{store},
	}
}

// REST implements a RESTStorage for signing keys against etcd.
type REST struct {
	*registry.Store
}
