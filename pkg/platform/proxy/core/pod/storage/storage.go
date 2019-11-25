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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// Storage includes storage for resources.
type Storage struct {
	Pod     *REST
	Status  *StatusREST
	Binding *BindingREST
	Events  *EventREST
	Log     *LogREST
}

// REST implements pkg/api/rest.StandardStorage.
type REST struct {
	*util.Store
}

// NewStorage returns a Storage object that will work against resources.
func NewStorage(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	podStore := &util.Store{
		NewFunc:        func() runtime.Object { return &corev1.Pod{} },
		NewListFunc:    func() runtime.Object { return &corev1.PodList{} },
		Namespaced:     true,
		PlatformClient: platformClient,
	}

	statusStore := *podStore

	return &Storage{
		Pod: &REST{podStore},
		Status: &StatusREST{
			store: &statusStore,
		},
		Binding: &BindingREST{
			platformClient: platformClient,
		},
		Events: &EventREST{
			platformClient: platformClient,
		},
		Log: &LogREST{
			platformClient: platformClient,
		},
	}
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"po"}
}

// Implement CategoriesProvider
var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return []string{"all"}
}

// BindingREST implements the REST endpoint for binding pods to nodes when etcd is in use.
type BindingREST struct {
	platformClient platforminternalclient.PlatformInterface
}

var _ = rest.Creater(&BindingREST{})

// NamespaceScoped fulfill rest.Scoper
func (r *BindingREST) NamespaceScoped() bool {
	return true
}

// New creates a new binding resource
func (r *BindingREST) New() runtime.Object {
	return &corev1.Binding{}
}

// Create ensures a pod is bound to a specific host.
func (r *BindingREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (out runtime.Object, err error) {
	client, requestInfo, err := util.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	result := &corev1.Binding{}
	if err := client.
		Post().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Body(obj).
		Do().
		Into(result); err != nil {
		return nil, err
	}
	return result, nil
}

// StatusREST implements the REST endpoint for changing the status of a pod.
type StatusREST struct {
	store *util.Store
}

// New creates a new pod resource
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
