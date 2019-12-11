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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"

	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/auth/registry/category"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for signing keys and all sub resources.
type Storage struct {
	*REST
}

// NewStorage returns a Storage object that will work against signing key.
func NewStorage(optsGetter generic.RESTOptionsGetter, authClient authinternalclient.AuthInterface) *Storage {
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

	return &Storage{&REST{store, authClient}}
}

// REST implements a RESTStorage for signing keys against etcd.
type REST struct {
	*registry.Store

	authClient authinternalclient.AuthInterface
}

func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {

	cat := obj.(*auth.Category)

	categorySelector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.categoryName", cat.Spec.CategoryName))

	categoryList, err := r.authClient.Categories().List(metav1.ListOptions{FieldSelector: categorySelector.String()})
	if err != nil {
		return nil, err
	}

	if len(categoryList.Items) != 0 {
		return nil, apierrors.NewConflict(
			auth.Resource("categories"),
			cat.Spec.CategoryName,
			fmt.Errorf("categoryName must be different"),
		)
	}

	return r.Store.Create(ctx, obj, createValidation, options)
}
