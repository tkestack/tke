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
	"tkestack.io/tke/pkg/platform/util"

	autoscalingAPIV1 "k8s.io/api/autoscaling/v1"
	autoscalingAPIV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingAPIV2Beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
)

// Storage includes storage for resources.
type Storage struct {
	HorizontalPodAutoscaler *REST
	Status                  *StatusREST
	Events                  *EventREST
}

// REST implements pkg/api/rest.StandardStorage
type REST struct {
	*util.Store
}

// NewStorageV1 returns a Storage object that will work against resources.
func NewStorageV1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	scalerStore := &util.Store{
		NewFunc:        func() runtime.Object { return &autoscalingAPIV1.HorizontalPodAutoscaler{} },
		NewListFunc:    func() runtime.Object { return &autoscalingAPIV1.HorizontalPodAutoscalerList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *scalerStore

	return &Storage{
		HorizontalPodAutoscaler: &REST{scalerStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
	}
}

// NewStorageV2Beta1 returns a Storage object that will work against resources.
func NewStorageV2Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	scalerStore := &util.Store{
		NewFunc:        func() runtime.Object { return &autoscalingAPIV2Beta1.HorizontalPodAutoscaler{} },
		NewListFunc:    func() runtime.Object { return &autoscalingAPIV2Beta1.HorizontalPodAutoscalerList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *scalerStore

	return &Storage{
		HorizontalPodAutoscaler: &REST{scalerStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
	}
}

// NewStorageV2Beta2 returns a Storage object that will work against resources.
func NewStorageV2Beta2(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	scalerStore := &util.Store{
		NewFunc:        func() runtime.Object { return &autoscalingAPIV2Beta2.HorizontalPodAutoscaler{} },
		NewListFunc:    func() runtime.Object { return &autoscalingAPIV2Beta2.HorizontalPodAutoscalerList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *scalerStore

	return &Storage{
		HorizontalPodAutoscaler: &REST{scalerStore},
		Status: &StatusREST{
			store: &statusStore,
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
	return []string{"hpa"}
}

// Implement CategoriesProvider
var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return []string{"all"}
}

// StatusREST implements the REST endpoint for changing the status of a resource
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
