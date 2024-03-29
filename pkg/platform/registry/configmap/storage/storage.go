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

package storage

import (
	"context"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/platform"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/platform/registry/configmap"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for namespace set and all sub resources.
type Storage struct {
	ConfigMap *REST
}

// NewStorage returns a Storage object that will work against namespace sets.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter) *Storage {
	strategy := configmap.NewStrategy()
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &platform.ConfigMap{} },
		NewListFunc:              func() runtime.Object { return &platform.ConfigMapList{} },
		DefaultQualifiedResource: platform.Resource("configmaps"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create configmap etcd rest storage", log.Err(err))
	}

	return &Storage{
		ConfigMap: &REST{store},
	}
}

// REST implements a RESTStorage for namespace sets against etcd.
type REST struct {
	*registry.Store
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"cm"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	return r.Store.List(ctx, wrappedOptions)
}
