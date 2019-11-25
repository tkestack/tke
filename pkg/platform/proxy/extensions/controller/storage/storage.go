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

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/util"
)

// Storage includes storage for resources.
type Storage struct {
	ReplicationController *REST
	Pods                  *PodREST
	Scale                 *ScaleREST
}

// REST implements pkg/api/rest.StandardStorage
type REST struct {
	rest.Storage
}

var _ rest.Storage = &REST{}

// NewStorageV1Beta1 returns a Storage object that will work against resources.
func NewStorageV1Beta1(_ genericregistry.RESTOptionsGetter, platformClient platforminternalclient.PlatformInterface) *Storage {
	return &Storage{
		ReplicationController: &REST{},
		Scale: &ScaleREST{
			platformClient: platformClient,
		},
		Pods: &PodREST{
			platformClient: platformClient,
		},
	}
}

// New creates a new Scale object
func (r *REST) New() runtime.Object {
	return &extensionsv1beta1.ReplicationControllerDummy{}
}

// NamespaceScoped returns if the object must be in a namespace.
func (r *REST) NamespaceScoped() bool {
	return true
}

// ScaleREST implements the REST endpoint for scale the extension.
type ScaleREST struct {
	rest.Storage
	platformClient platforminternalclient.PlatformInterface
}

// ScaleREST implements Patcher
var _ = rest.Patcher(&ScaleREST{})

// New creates a new Scale object
func (r *ScaleREST) New() runtime.Object {
	return &extensionsv1beta1.Scale{}
}

// Get finds a resource in the storage by name and returns it.
func (r *ScaleREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	client, requestInfo, err := util.RESTClient(ctx, r.platformClient)
	if err != nil {
		return nil, err
	}

	result := &extensionsv1beta1.Scale{}
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

	result := &extensionsv1beta1.Scale{}
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
