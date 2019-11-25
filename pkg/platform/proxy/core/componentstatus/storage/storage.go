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
	metaInternalVersion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// Storage includes storage for resources.
type Storage struct {
	ComponentStatus *REST
}

// REST implements pkg/api/rest.StandardStorage.
type REST struct {
	rest.Storage

	platformClient platforminternalclient.PlatformInterface
}

// NewStorage returns a Storage object that will work against resources.
func NewStorage(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	return &Storage{
		ComponentStatus: &REST{
			platformClient: platformClient,
		},
	}
}

// Implement ShortNamesProvider
var _ rest.ShortNamesProvider = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of
// short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"cs"}
}

// NamespaceScoped returns if the object must be in a namespace.
func (r *REST) NamespaceScoped() bool {
	return false
}

// New returns an empty object that can be used with Create and Update after
// request data has been put into it.
func (r *REST) New() runtime.Object {
	return &corev1.ComponentStatus{}
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &corev1.ComponentStatusList{}
}

// List selects resources in the storage which match to the selector. 'options'
// can be nil.
func (r *REST) List(ctx context.Context, options *metaInternalVersion.ListOptions) (runtime.Object, error) {
	client, requestInfo, err := util.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	result := r.NewList()
	if err := client.
		Get().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		SpecificallyVersionedParams(options, metaInternalVersion.ParameterCodec, metaInternalVersion.SchemeGroupVersion).
		Do().
		Into(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	client, requestInfo, err := util.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	result := r.New()
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
