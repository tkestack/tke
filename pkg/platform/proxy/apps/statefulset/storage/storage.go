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

	"tkestack.io/tke/pkg/platform/proxy"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
)

// Storage includes storage for resources.
type Storage struct {
	StatefulSet *REST
	Pods        *PodREST
	ReplicaSets *ReplicaSetsREST
	Status      *StatusREST
	Scale       *ScaleREST
	Events      *EventREST
	HPAs        *HPARest
}

// REST implements pkg/api/rest.StandardStorage
type REST struct {
	*proxy.Store
}

// NewStorageV1 returns a Storage object that will work against resources.
func NewStorageV1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	statefulSetStore := &proxy.Store{
		NewFunc:        func() runtime.Object { return &appsv1.StatefulSet{} },
		NewListFunc:    func() runtime.Object { return &appsv1.StatefulSetList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *statefulSetStore

	return &Storage{
		StatefulSet: &REST{statefulSetStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Pods: &PodREST{
			platformClient: platformClient,
		},
		Scale: &ScaleREST{
			platformClient: platformClient,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
		HPAs: &HPARest{
			platformClient: platformClient,
		},
		ReplicaSets: &ReplicaSetsREST{
			platformClient: platformClient,
		},
	}
}

// NewStorageV1Beta1 returns a Storage object that will work against resources.
func NewStorageV1Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	statefulSetStore := &proxy.Store{
		NewFunc:        func() runtime.Object { return &appsv1beta1.StatefulSet{} },
		NewListFunc:    func() runtime.Object { return &appsv1beta1.StatefulSetList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *statefulSetStore

	return &Storage{
		StatefulSet: &REST{statefulSetStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Pods: &PodREST{
			platformClient: platformClient,
		},
		Scale: &ScaleREST{
			platformClient: platformClient,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
		HPAs: &HPARest{
			platformClient: platformClient,
		},
		ReplicaSets: &ReplicaSetsREST{
			platformClient: platformClient,
		},
	}
}

// NewStorageV1Beta2 returns a Storage object that will work against resources.
func NewStorageV1Beta2(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	statefulSetStore := &proxy.Store{
		NewFunc:        func() runtime.Object { return &appsv1beta2.StatefulSet{} },
		NewListFunc:    func() runtime.Object { return &appsv1beta2.StatefulSetList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *statefulSetStore

	return &Storage{
		StatefulSet: &REST{statefulSetStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Scale: &ScaleREST{
			platformClient: platformClient,
		},
		Pods: &PodREST{
			platformClient: platformClient,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
		HPAs: &HPARest{
			platformClient: platformClient,
		},
		ReplicaSets: &ReplicaSetsREST{
			platformClient: platformClient,
		},
	}
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"sts"}
}

// Implement CategoriesProvider
var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return []string{"all"}
}

// StatusREST implements the REST endpoint for changing the status of a replication controller
type StatusREST struct {
	rest.Storage
	store *proxy.Store
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

// ScaleREST implements the REST endpoint for scale the statefulset.
type ScaleREST struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

// ScaleREST implements Patcher
var _ = rest.Patcher(&ScaleREST{})
var _ = rest.GroupVersionKindProvider(&ScaleREST{})

// GroupVersionKind is used to specify a particular GroupVersionKind to discovery.
func (r *ScaleREST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	switch containingGV {
	case appsv1beta1.SchemeGroupVersion:
		return appsv1beta1.SchemeGroupVersion.WithKind("Scale")
	case appsv1beta2.SchemeGroupVersion:
		return appsv1beta2.SchemeGroupVersion.WithKind("Scale")
	default:
		return autoscalingv1.SchemeGroupVersion.WithKind("Scale")
	}
}

// New creates a new Scale object.
func (r *ScaleREST) New() runtime.Object {
	return &autoscalingv1.Scale{}
}

// Get finds a resource in the storage by name and returns it.
func (r *ScaleREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	client, requestInfo, err := proxy.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	result := &autoscalingv1.Scale{}
	if err := client.
		Get().
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(name).
		VersionedParams(options, metav1.ParameterCodec).
		Do(ctx).
		Into(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Update finds a resource in the storage and updates it.
func (r *ScaleREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	client, requestInfo, err := proxy.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, false, err
	}

	if requestInfo.Verb == "patch" {
		requestBody, ok := filter.RequestBodyFrom(ctx)
		if !ok {
			return nil, false, errors.NewBadRequest("request body is required")
		}
		result := &autoscalingv1.Scale{}
		if err := client.
			Patch(types.PatchType(requestBody.ContentType)).
			NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "").
			Resource(requestInfo.Resource).
			SubResource(requestInfo.Subresource).
			Name(name).
			Body(requestBody.Data).
			Do(ctx).
			Into(result); err != nil {
			return nil, false, err
		}
		return result, true, nil
	}

	obj, err := objInfo.UpdatedObject(ctx, nil)
	if err != nil {
		return nil, false, errors.NewInternalError(err)
	}

	result := &autoscalingv1.Scale{}
	if err := client.
		Put().
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(name).
		Body(obj).
		Do(ctx).
		Into(result); err != nil {
		return nil, false, err
	}

	return result, true, nil
}
