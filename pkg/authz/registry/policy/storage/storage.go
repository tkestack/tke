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
	"fmt"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/authz"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/authz/registry/policy"
	"tkestack.io/tke/pkg/util/log"
)

// Storage includes storage for configmap and all sub resources.
type Storage struct {
	Policy *REST
}

// NewStorage returns a Storage object that will work against configmap.
func NewStorage(optsGetter genericregistry.RESTOptionsGetter, platformClient platformversionedclient.PlatformV1Interface) *Storage {
	strategy := policy.NewStrategy(platformClient)
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &authz.Policy{} },
		NewListFunc:              func() runtime.Object { return &authz.PolicyList{} },
		DefaultQualifiedResource: authz.Resource("policies"),
		ShouldDeleteDuringUpdate: policy.ShouldDeleteDuringUpdate,
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
	}
	store.TableConvertor = rest.NewDefaultTableConvertor(store.DefaultQualifiedResource)
	options := &genericregistry.StoreOptions{
		RESTOptions: optsGetter,
	}

	if err := store.CompleteWithOptions(options); err != nil {
		log.Panic("Failed to create configmap etcd rest storage", log.Err(err))
	}
	return &Storage{
		Policy: &REST{store},
	}
}

// REST implements a RESTStorage for configmap against etcd.
type REST struct {
	*registry.Store
}

var _ rest.ShortNamesProvider = &REST{}
var _ rest.Getter = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"pol"}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	return r.Store.List(ctx, options)
}

// Delete enforces life-cycle rules for policy termination
func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID == "" {
		tenantID = "default"
	}
	object, err := r.Get(ctx, name, &metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	pol := object.(*authz.Policy)
	if tenantID != "default" && pol.Namespace != tenantID {
		return nil, false, fmt.Errorf("tenant '%s' can't delete policy '%s/%s'", tenantID, pol.Namespace, pol.Name)
	}
	return r.Store.Delete(ctx, name, deleteValidation, options)
}
