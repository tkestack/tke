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
	appsV1 "k8s.io/api/apps/v1"
	appsV1Beta2 "k8s.io/api/apps/v1beta2"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// Storage includes storage for resources.
type Storage struct {
	DaemonSet *REST
	Status    *StatusREST
	Events    *EventREST
	Pods      *PodREST
}

// REST implements pkg/api/rest.StandardStorage
type REST struct {
	*util.Store
}

// NewStorageV1 returns a Storage object that will work against resources.
func NewStorageV1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	daemonSetStore := &util.Store{
		NewFunc:        func() runtime.Object { return &appsV1.DaemonSet{} },
		NewListFunc:    func() runtime.Object { return &appsV1.DaemonSetList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *daemonSetStore

	return &Storage{
		DaemonSet: &REST{daemonSetStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Pods: &PodREST{
			platformClient: platformClient,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
	}
}

// NewStorageV1Beta2 returns a Storage object that will work against resources.
func NewStorageV1Beta2(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	daemonSetStore := &util.Store{
		NewFunc:        func() runtime.Object { return &appsV1Beta2.DaemonSet{} },
		NewListFunc:    func() runtime.Object { return &appsV1Beta2.DaemonSetList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *daemonSetStore

	return &Storage{
		DaemonSet: &REST{daemonSetStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Pods: &PodREST{
			platformClient: platformClient,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
	}
}

// NewStorageExtensionsV1Beta1 returns a Storage object that will work against resources.
func NewStorageExtensionsV1Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	daemonSetStore := &util.Store{
		NewFunc:        func() runtime.Object { return &extensionsv1beta1.DaemonSet{} },
		NewListFunc:    func() runtime.Object { return &extensionsv1beta1.DaemonSetList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *daemonSetStore

	return &Storage{
		DaemonSet: &REST{daemonSetStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Pods: &PodREST{
			platformClient: platformClient,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
	}
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"ds"}
}

// Implement CategoriesProvider
var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return []string{"all"}
}

// StatusREST implements the REST endpoint for changing the status of a resource.
type StatusREST struct {
	rest.Storage
	store *util.Store
}

// StatusREST implements Patcher
var _ = rest.Patcher(&StatusREST{})

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *StatusREST) New() runtime.Object {
	return r.store.New()
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}
