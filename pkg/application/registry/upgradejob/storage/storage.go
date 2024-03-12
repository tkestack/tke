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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/application"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/application/registry/upgradejob"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for upgradejob and all sub resources.
type Storage struct {
	UpgradeJob *REST
}

// NewStorage returns a Storage object that will work against upgradejob.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter) *Storage {
	strategy := upgradejob.NewStrategy()
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &application.UpgradeJob{} },
		NewListFunc:              func() runtime.Object { return &application.UpgradeJobList{} },
		DefaultQualifiedResource: application.Resource("upgradejobs"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create upgradejob etcd rest storage", log.Err(err))
	}

	return &Storage{
		UpgradeJob: &REST{store},
	}
}

// REST implements a RESTStorage for upgradejob against etcd.
type REST struct {
	*registry.Store
}

var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"uj"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	return r.Store.List(ctx, wrappedOptions)
}

// Watch selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) Watch(ctx context.Context, options *metainternal.ListOptions) (watch.Interface, error) {
	wrappedOptions := apiserverutil.PredicateListOptions(ctx, options)
	return r.Store.Watch(ctx, wrappedOptions)
}

func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	return r.Store.Create(ctx, obj, createValidation, options)
}
