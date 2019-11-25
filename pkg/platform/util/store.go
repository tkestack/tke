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

package util

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/rest"
	clientrest "k8s.io/client-go/rest"
	"reflect"
	"strings"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
)

// Store implements pkg/api/rest.StandardStorage.
type Store struct {
	// NewFunc returns a new instance of the type this registry returns for a
	// GET of a single object
	NewFunc func() runtime.Object
	// NewListFunc returns a new list of the type this registry
	NewListFunc    func() runtime.Object
	Namespaced     bool
	PlatformClient platforminternalclient.PlatformInterface
}

// var _ rest.Exporter = &Store{}
// var _ rest.TableConvertor = &Store{}

// New implements RESTStorage.New.
func (s *Store) New() runtime.Object {
	return s.NewFunc()
}

// NewList implements rest.Lister.
func (s *Store) NewList() runtime.Object {
	return s.NewListFunc()
}

// NamespaceScoped indicates whether the resource is namespaced.
func (s *Store) NamespaceScoped() bool {
	return s.Namespaced
}

// List returns a list of items matching labels and field according to the
// backend kubernetes api server.
func (s *Store) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	client, requestInfo, err := RESTClient(ctx, s.PlatformClient)
	if err != nil {
		return nil, err
	}

	fuzzyResourceName := filter.FuzzyResourceFrom(ctx)
	if options != nil && options.FieldSelector != nil {
		if name, ok := options.FieldSelector.RequiresExactMatch("metadata.name"); ok {
			options.FieldSelector, _ = options.FieldSelector.Transform(func(k, v string) (string, string, error) {
				if k == "metadata.name" {
					return "", "", nil
				}
				return k, v, nil
			})
			fuzzyResourceName = name
		}
	}

	result := s.NewListFunc()
	if err := client.
		Get().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		SpecificallyVersionedParams(options, metainternalversion.ParameterCodec, v1.SchemeGroupVersion).
		Do().
		Into(result); err != nil {
		return nil, err
	}

	if fuzzyResourceName != "" {
		if _, ok := result.(v1.ListInterface); ok {
			v := reflect.ValueOf(result).Elem()
			if items := v.FieldByName("Items"); items != (reflect.Value{}) {
				if items.Kind() == reflect.Slice || items.Kind() == reflect.Array {
					newResult := make([]reflect.Value, 0)
					for i := 0; i < items.Len(); i++ {
						item := items.Index(i)
						for j := 0; j < item.Type().NumField(); j++ {
							itemChild := item.Field(j)
							if metadata, ok := itemChild.Interface().(v1.ObjectMeta); ok {
								if strings.Contains(strings.ToLower(metadata.Name), strings.ToLower(fuzzyResourceName)) {
									newResult = append(newResult, item)
								}
								break
							}
						}
					}
					slice := reflect.MakeSlice(items.Type(), 0, 0)
					newResultValue := reflect.Append(slice, newResult...)
					v.FieldByName("Items").Set(newResultValue)
				}
			}
		}
	}

	return result, nil
}

// Get retrieves the item from storage.
func (s *Store) Get(ctx context.Context, name string, options *v1.GetOptions) (runtime.Object, error) {
	client, requestInfo, err := RESTClient(ctx, s.PlatformClient)
	if err != nil {
		return nil, err
	}

	result := s.New()
	if err := client.
		Get().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(name).
		VersionedParams(options, v1.ParameterCodec).
		Do().
		Into(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Watch makes a matcher for the given label and field, and calls
// WatchPredicate. If possible, you should customize PredicateFunc to produce
// a matcher that matches by key. SelectionPredicate does this for you
// automatically.
func (s *Store) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	client, requestInfo, err := RESTClient(ctx, s.PlatformClient)
	if err != nil {
		return nil, err
	}

	options.Watch = true

	return client.Get().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		SpecificallyVersionedParams(options, metainternalversion.ParameterCodec, metainternalversion.SchemeGroupVersion).
		Watch()
}

// Create inserts a new item according to the unique key from the object.
func (s *Store) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, _ *v1.CreateOptions) (runtime.Object, error) {
	client, requestInfo, err := RESTClient(ctx, s.PlatformClient)
	if err != nil {
		return nil, err
	}

	requestBody, ok := filter.RequestBodyFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("request body is required")
	}

	result := s.New()
	if err := client.
		Post().
		Context(ctx).
		SetHeader("Content-Type", requestBody.ContentType).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Body(requestBody.Data).
		Do().
		Into(result); err != nil {
		return nil, err
	}
	return result, nil
}

// Update finds a resource in the storage and updates it. Some implementations
// may allow updates creates the object - they should set the created boolean
// to true.
func (s *Store) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *v1.UpdateOptions) (runtime.Object, bool, error) {
	client, requestInfo, err := RESTClient(ctx, s.PlatformClient)
	if err != nil {
		return nil, false, err
	}
	requestBody, ok := filter.RequestBodyFrom(ctx)
	if !ok {
		return nil, false, errors.NewBadRequest("request body is required")
	}

	result := s.New()

	var req *clientrest.Request
	if requestInfo.Verb == "patch" {
		req = client.Patch(types.PatchType(requestBody.ContentType))
	} else if requestInfo.Verb == "update" || requestInfo.Verb == "put" {
		req = client.Put()
	} else {
		return nil, false, errors.NewBadRequest("unsupported request method")
	}
	if err := req.
		Context(ctx).
		SetHeader("Content-Type", requestBody.ContentType).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(name).
		Body(requestBody.Data).
		Do().
		Into(result); err != nil {
		return nil, false, err
	}

	return result, true, nil
}

// Delete finds a resource in the storage and deletes it.
func (s *Store) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *v1.DeleteOptions) (runtime.Object, bool, error) {
	client, requestInfo, err := RESTClient(ctx, s.PlatformClient)
	if err != nil {
		return nil, false, err
	}

	result := client.
		Delete().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(name).
		Body(options).
		Do()
	resultErr := result.Error()
	if resultErr != nil {
		return nil, false, resultErr
	}
	returnedObj, err := result.Get()
	if err != nil {
		return nil, false, err
	}
	return returnedObj, true, nil
}

// DeleteCollection selects all resources in the storage matching given 'listOptions'
// and deletes them. If 'options' are provided, the resource will attempt to honor
// them or return an invalid request error.
// DeleteCollection may not be atomic - i.e. it may delete some objects and still
// return an error after it. On success, returns a list of deleted objects.
func (s *Store) DeleteCollection(ctx context.Context, options *v1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	client, requestInfo, err := RESTClient(ctx, s.PlatformClient)
	if err != nil {
		return nil, err
	}

	result := client.
		Delete().
		Context(ctx).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SpecificallyVersionedParams(listOptions, metainternalversion.ParameterCodec, metainternalversion.SchemeGroupVersion).
		Body(options).
		Do()

	resultErr := result.Error()
	if resultErr != nil {
		return nil, resultErr
	}
	returnedObj, err := result.Get()
	if err != nil {
		return nil, err
	}
	return returnedObj, nil
}
