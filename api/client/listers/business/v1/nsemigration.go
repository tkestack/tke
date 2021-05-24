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
	v1 "tkestack.io/tke/api/business/v1"
)

// NsEmigrationLister helps list NsEmigrations.
// All objects returned here must be treated as read-only.
type NsEmigrationLister interface {
	// List lists all NsEmigrations in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.NsEmigration, err error)
	// NsEmigrations returns an object that can list and get NsEmigrations.
	NsEmigrations(namespace string) NsEmigrationNamespaceLister
	NsEmigrationListerExpansion
}

// nsEmigrationLister implements the NsEmigrationLister interface.
type nsEmigrationLister struct {
	indexer cache.Indexer
}

// NewNsEmigrationLister returns a new NsEmigrationLister.
func NewNsEmigrationLister(indexer cache.Indexer) NsEmigrationLister {
	return &nsEmigrationLister{indexer: indexer}
}

// List lists all NsEmigrations in the indexer.
func (s *nsEmigrationLister) List(selector labels.Selector) (ret []*v1.NsEmigration, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.NsEmigration))
	})
	return ret, err
}

// NsEmigrations returns an object that can list and get NsEmigrations.
func (s *nsEmigrationLister) NsEmigrations(namespace string) NsEmigrationNamespaceLister {
	return nsEmigrationNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// NsEmigrationNamespaceLister helps list and get NsEmigrations.
// All objects returned here must be treated as read-only.
type NsEmigrationNamespaceLister interface {
	// List lists all NsEmigrations in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.NsEmigration, err error)
	// Get retrieves the NsEmigration from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.NsEmigration, error)
	NsEmigrationNamespaceListerExpansion
}

// nsEmigrationNamespaceLister implements the NsEmigrationNamespaceLister
// interface.
type nsEmigrationNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all NsEmigrations in the indexer for a given namespace.
func (s nsEmigrationNamespaceLister) List(selector labels.Selector) (ret []*v1.NsEmigration, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.NsEmigration))
	})
	return ret, err
}

// Get retrieves the NsEmigration from the indexer for a given namespace and name.
func (s nsEmigrationNamespaceLister) Get(name string) (*v1.NsEmigration, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("nsemigration"), name)
	}
	return obj.(*v1.NsEmigration), nil
}
