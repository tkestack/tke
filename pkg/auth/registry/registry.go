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

package registry

import (
	"tkestack.io/tke/pkg/auth/registry/apikey"
	"tkestack.io/tke/pkg/auth/registry/category"
	"tkestack.io/tke/pkg/auth/registry/localidentity"
	"tkestack.io/tke/pkg/auth/registry/policy"
	"tkestack.io/tke/pkg/auth/registry/role"

	"github.com/coreos/etcd/clientv3"
	"github.com/dexidp/dex/storage"
)

// Registry represents the module of all object operations on backend storage.
type Registry struct {
	dexStorage           storage.Storage
	localIdentityStorage *localidentity.Storage
	roleStorage          *role.Storage
	policyStorage        *policy.Storage
	categoryStorage      *category.Storage
	apikeyStorage        apikey.Storage
}

// NewRegistry to create the new Registry object instance.
func NewRegistry(db *clientv3.Client, store storage.Storage) (*Registry, error) {
	return &Registry{
		dexStorage:           store,
		localIdentityStorage: localidentity.NewLocalIdentity(db),
		roleStorage:          role.NewRoleStorage(db),
		policyStorage:        policy.NewPolicyStorage(db),
		categoryStorage:      category.NewCategoryStorage(db),
		apikeyStorage:        apikey.NewAPIKeyStorage(db),
	}, nil
}

// LocalIdentityStorage returns user storage instance.
func (r *Registry) LocalIdentityStorage() *localidentity.Storage {
	return r.localIdentityStorage
}

// RoleStorage returns role storage instance.
func (r *Registry) RoleStorage() *role.Storage {
	return r.roleStorage
}

// PolicyStorage returns policy storage instance.
func (r *Registry) PolicyStorage() *policy.Storage {
	return r.policyStorage
}

// DexStorage returns dex storage instance
func (r *Registry) DexStorage() storage.Storage {
	return r.dexStorage
}

// CategoryStorage returns category storage instance.
func (r *Registry) CategoryStorage() *category.Storage {
	return r.categoryStorage
}

// APIKeyStorage returns apikey storage instance.
func (r *Registry) APIKeyStorage() apikey.Storage {
	return r.apikeyStorage
}
