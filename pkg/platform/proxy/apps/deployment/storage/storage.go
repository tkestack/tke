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
	appsV1Beta1 "k8s.io/api/apps/v1beta1"
	appsV1Beta2 "k8s.io/api/apps/v1beta2"
	autoscalingV1API "k8s.io/api/autoscaling/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// Storage includes storage for resources.
type Storage struct {
	Deployment *REST
	Status     *StatusREST
	Scale      *ScaleREST
	Rollback   *RollbackREST
	Pods       *PodREST
	HPAs       *HPARest
	Events     *EventREST
}

// REST implements pkg/api/rest.StandardStorage
type REST struct {
	*util.Store
}

// NewStorageV1 returns a Storage object that will work against resources.
func NewStorageV1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	deploymentStore := &util.Store{
		NewFunc:        func() runtime.Object { return &appsV1.Deployment{} },
		NewListFunc:    func() runtime.Object { return &appsV1.DeploymentList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *deploymentStore

	return &Storage{
		Deployment: &REST{deploymentStore},
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
	}
}

// NewStorageV1Beta1 returns a Storage object that will work against resources.
func NewStorageV1Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	deploymentStore := &util.Store{
		NewFunc:     func() runtime.Object { return &appsV1Beta1.Deployment{} },
		NewListFunc: func() runtime.Object { return &appsV1Beta1.DeploymentList{} },
		Namespaced:  true,

		PlatformClient: platformClient,
	}

	statusStore := *deploymentStore

	return &Storage{
		Deployment: &REST{deploymentStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Scale: &ScaleREST{

			platformClient: platformClient,
		},
		Rollback: &RollbackREST{

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
	}
}

// NewStorageExtensionsV1Beta1 returns a Storage object that will work against resources.
func NewStorageExtensionsV1Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	deploymentStore := &util.Store{
		NewFunc:     func() runtime.Object { return &extensionsv1beta1.Deployment{} },
		NewListFunc: func() runtime.Object { return &extensionsv1beta1.DeploymentList{} },
		Namespaced:  true,

		PlatformClient: platformClient,
	}

	statusStore := *deploymentStore

	return &Storage{
		Deployment: &REST{deploymentStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Scale: &ScaleREST{
			platformClient: platformClient,
		},
		Rollback: &RollbackREST{
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
	}
}

// NewStorageV1Beta2 returns a Storage object that will work against resources.
func NewStorageV1Beta2(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	deploymentStore := &util.Store{
		NewFunc:        func() runtime.Object { return &appsV1Beta2.Deployment{} },
		NewListFunc:    func() runtime.Object { return &appsV1Beta2.DeploymentList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *deploymentStore

	return &Storage{
		Deployment: &REST{deploymentStore},
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
	}
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"deploy"}
}

// Implement CategoriesProvider
var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return []string{"all"}
}

// StatusREST implements the REST endpoint for scale the deployment.
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

// ScaleREST implements the REST endpoint for scale the deployment.
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
	case extensionsv1beta1.SchemeGroupVersion:
		return extensionsv1beta1.SchemeGroupVersion.WithKind("Scale")
	case appsV1Beta1.SchemeGroupVersion:
		return appsV1Beta1.SchemeGroupVersion.WithKind("Scale")
	case appsV1Beta2.SchemeGroupVersion:
		return appsV1Beta2.SchemeGroupVersion.WithKind("Scale")
	default:
		return autoscalingV1API.SchemeGroupVersion.WithKind("Scale")
	}
}

// New creates a new Scale object.
func (r *ScaleREST) New() runtime.Object {
	return &autoscalingV1API.Scale{}
}

// Get finds a resource in the storage by name and returns it.
func (r *ScaleREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	client, requestInfo, err := util.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	result := &autoscalingV1API.Scale{}
	if err := client.
		Get().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(name).
		VersionedParams(options, metav1.ParameterCodec).
		Do().
		Into(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Update finds a resource in the storage and updates it.
func (r *ScaleREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	client, requestInfo, err := util.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, false, err
	}

	obj, err := objInfo.UpdatedObject(ctx, nil)
	if err != nil {
		return nil, false, errors.NewInternalError(err)
	}

	result := &autoscalingV1API.Scale{}
	if err := client.
		Put().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(name).
		Body(obj).
		Do().
		Into(result); err != nil {
		return nil, false, err
	}

	return result, true, nil
}

// RollbackREST implements Creater
type RollbackREST struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

var _ = rest.Creater(&RollbackREST{})
var _ = rest.GroupVersionKindProvider(&RollbackREST{})

// GroupVersionKind is used to specify a particular GroupVersionKind to discovery.
func (r *RollbackREST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return extensionsv1beta1.SchemeGroupVersion.WithKind("DeploymentRollback")
}

// New creates a new Rollback object
func (r *RollbackREST) New() runtime.Object {
	return &extensionsv1beta1.DeploymentRollback{}
}

// Create inserts a new item according to the unique key from the object.
func (r *RollbackREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	client, requestInfo, err := util.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	result := &metav1.Status{}
	if err := client.
		Post().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(requestInfo.Name).
		Body(obj).
		Do().
		Into(result); err != nil {
		return nil, err
	}
	return result, nil
}
