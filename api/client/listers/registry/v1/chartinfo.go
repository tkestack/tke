/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1 "tkestack.io/tke/api/registry/v1"
)

// ChartInfoLister helps list ChartInfos.
// All objects returned here must be treated as read-only.
type ChartInfoLister interface {
	// List lists all ChartInfos in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ChartInfo, err error)
	// ChartInfos returns an object that can list and get ChartInfos.
	ChartInfos(namespace string) ChartInfoNamespaceLister
	ChartInfoListerExpansion
}

// chartInfoLister implements the ChartInfoLister interface.
type chartInfoLister struct {
	indexer cache.Indexer
}

// NewChartInfoLister returns a new ChartInfoLister.
func NewChartInfoLister(indexer cache.Indexer) ChartInfoLister {
	return &chartInfoLister{indexer: indexer}
}

// List lists all ChartInfos in the indexer.
func (s *chartInfoLister) List(selector labels.Selector) (ret []*v1.ChartInfo, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ChartInfo))
	})
	return ret, err
}

// ChartInfos returns an object that can list and get ChartInfos.
func (s *chartInfoLister) ChartInfos(namespace string) ChartInfoNamespaceLister {
	return chartInfoNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ChartInfoNamespaceLister helps list and get ChartInfos.
// All objects returned here must be treated as read-only.
type ChartInfoNamespaceLister interface {
	// List lists all ChartInfos in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ChartInfo, err error)
	// Get retrieves the ChartInfo from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.ChartInfo, error)
	ChartInfoNamespaceListerExpansion
}

// chartInfoNamespaceLister implements the ChartInfoNamespaceLister
// interface.
type chartInfoNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ChartInfos in the indexer for a given namespace.
func (s chartInfoNamespaceLister) List(selector labels.Selector) (ret []*v1.ChartInfo, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ChartInfo))
	})
	return ret, err
}

// Get retrieves the ChartInfo from the indexer for a given namespace and name.
func (s chartInfoNamespaceLister) Get(name string) (*v1.ChartInfo, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("chartinfo"), name)
	}
	return obj.(*v1.ChartInfo), nil
}
