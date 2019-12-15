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
	jsoniter "github.com/json-iterator/go"
	"k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/pkg/monitor/storage"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Storage includes storage for metrics and all sub resources.
type Storage struct {
	Metric *REST
}

// NewStorage returns a Storage object that will work against metrics.
func NewStorage(_ genericregistry.RESTOptionsGetter, metricStorage storage.MetricStorage) *Storage {
	return &Storage{
		Metric: &REST{
			metricStorage: metricStorage,
		},
	}
}

// REST implements a RESTStorage for metrics against etcd.
type REST struct {
	rest.Storage
	metricStorage storage.MetricStorage
}

var _ rest.ShortNamesProvider = &REST{}
var _ rest.Creater = &REST{}
var _ rest.Scoper = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"metric"}
}

// NamespaceScoped returns true if the storage is namespaced
func (r *REST) NamespaceScoped() bool {
	return false
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *REST) New() runtime.Object {
	return &monitor.Metric{}
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &monitor.MetricList{}
}

// Create creates a new version of a resource.
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	metric, ok := obj.(*monitor.Metric)
	if !ok {
		return nil, errors.NewBadRequest("failed to processed request body")
	}
	result, err := r.metricStorage.Query(&metric.Query)
	if err != nil {
		return nil, err
	}
	jsonResult, err := json.MarshalToString(result)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	return &monitor.Metric{
		Query:      metric.Query,
		JSONResult: jsonResult,
	}, nil
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	return &monitor.MetricList{
		Items: make([]monitor.Metric, 0),
	}, nil
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return nil, errors.NewNotFound(monitor.Resource("metric"), name)
}
