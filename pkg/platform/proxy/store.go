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

package proxy

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	clientrest "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	apiserverutil "tkestack.io/tke/pkg/apiserver/util"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
)

var proxyConvert = conversion.NewConverter(conversion.DefaultNameFunc)

func init() {
	_ = proxyConvert.RegisterUntypedConversionFunc((*metainternalversion.ListOptions)(nil), (*v1.ListOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return metainternalversion.Convert_internalversion_ListOptions_To_v1_ListOptions(a.(*metainternalversion.ListOptions), b.(*v1.ListOptions), scope)
	})

}

// Store implements pkg/api/rest.StandardStorage.
type Store struct {
	// NewFunc returns a new instance of the type this registry returns for a
	// GET of a single object
	NewFunc func() runtime.Object
	// NewListFunc returns a new list of the type this registry
	NewListFunc func() runtime.Object
	// DefaultQualifiedResource is the pluralized name of the resource.
	// This field is used if there is no request info present in the context.
	// See qualifiedResourceFromContext for details.
	DefaultQualifiedResource schema.GroupResource
	// TableConvertor is an optional interface for transforming items or lists
	// of items into tabular output. If unset, the default will be used.
	TableConvertor rest.TableConvertor

	Namespaced     bool
	PlatformClient platforminternalclient.PlatformInterface
}

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

// ConvertToTable converts objects to metav1.Table objects using default table
// convertor.
func (s *Store) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*v1.Table, error) {
	if s.TableConvertor != nil {
		return s.TableConvertor.ConvertToTable(ctx, object, tableOptions)
	}
	return rest.NewDefaultTableConvertor(s.qualifiedResourceFromContext(ctx)).ConvertToTable(ctx, object, tableOptions)
}

// List returns a list of items matching labels and field according to the
// backend kubernetes api server.
func (s *Store) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	client, requestInfo, err := RESTClient(ctx, s.PlatformClient)
	if err != nil {
		return nil, err
	}

	fuzzyResourceName := filter.FuzzyResourceFrom(ctx)
	options, fuzzyResourceName = apiserverutil.InterceptFuzzyResourceNameFromListOptions(options, fuzzyResourceName)
	v1options := &v1.ListOptions{}
	err = proxyConvert.Convert(options, v1options, &conversion.Meta{})
	if err != nil {
		return nil, fmt.Errorf("convert failed: %v", err)
	}

	result := s.NewListFunc()
	if err := client.
		Get().
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		SpecificallyVersionedParams(v1options, platform.ParameterCodec, v1.SchemeGroupVersion).
		Do(ctx).
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
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		Name(name).
		VersionedParams(options, platform.ParameterCodec).
		Do(ctx).
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
	v1options := &v1.ListOptions{}
	err = proxyConvert.Convert(options, v1options, &conversion.Meta{})
	if err != nil {
		return nil, fmt.Errorf("convert failed: %v", err)
	}

	return client.Get().
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		SpecificallyVersionedParams(v1options, platform.ParameterCodec, v1.SchemeGroupVersion).
		Watch(ctx)
}

// Create inserts a new item according to the unique key from the object.
func (s *Store) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *v1.CreateOptions) (runtime.Object, error) {
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
		SetHeader("Content-Type", requestBody.ContentType).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		VersionedParams(options, platform.ParameterCodec).
		Body(requestBody.Data).
		Do(ctx).
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
		SetHeader("Content-Type", requestBody.ContentType).
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		VersionedParams(options, platform.ParameterCodec).
		Name(name).
		Body(requestBody.Data).
		Do(ctx).
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
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SubResource(requestInfo.Subresource).
		VersionedParams(options, platform.ParameterCodec).
		Name(name).
		Body(options).
		Do(ctx)
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
	v1listoptions := &v1.ListOptions{}
	err = proxyConvert.Convert(listOptions, v1listoptions, &conversion.Meta{})
	if err != nil {
		return nil, fmt.Errorf("convert failed: %v", err)
	}

	result := client.
		Delete().
		NamespaceIfScoped(requestInfo.Namespace, requestInfo.Namespace != "" && requestInfo.Resource != "namespaces").
		Resource(requestInfo.Resource).
		SpecificallyVersionedParams(v1listoptions, platform.ParameterCodec, v1.SchemeGroupVersion).
		Body(options).
		Do(ctx)

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

// qualifiedResourceFromContext attempts to retrieve a GroupResource from the context's request info.
// If the context has no request info, DefaultQualifiedResource is used.
func (s *Store) qualifiedResourceFromContext(ctx context.Context) schema.GroupResource {
	if info, ok := genericapirequest.RequestInfoFrom(ctx); ok {
		return schema.GroupResource{Group: info.APIGroup, Resource: info.Resource}
	}
	// some implementations access storage directly and thus the context has no RequestInfo
	return s.DefaultQualifiedResource
}
