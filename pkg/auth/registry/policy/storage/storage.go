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
	"context"

	"tkestack.io/tke/pkg/auth/registry/policy"
	"tkestack.io/tke/pkg/util/log"

	metaInternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"

	"tkestack.io/tke/api/auth"
)

// Storage includes storage for policies and all sub resources.
type Storage struct {
	Policy *REST
}

// NewStorage returns a Storage object that will work against policies.
func NewStorage(optsGetter generic.RESTOptionsGetter) *Storage {
	strategy := policy.NewStrategy()
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &auth.Policy{} },
		NewListFunc:              func() runtime.Object { return &auth.PolicyList{} },
		DefaultQualifiedResource: auth.Resource("policies"),
		PredicateFunc:            policy.MatchPolicy,

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    policy.GetAttrs,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create policy etcd rest storage", log.Err(err))
	}

	return &Storage{
		Policy: &REST{store},
	}
}

// REST implements a RESTStorage for clusters against etcd.
type REST struct {
	*registry.Store
}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"pol"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metaInternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	return r.Store.List(ctx, wrappedOptions)
}
