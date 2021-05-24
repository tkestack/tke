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

package internalversion

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	auth "tkestack.io/tke/api/auth"
)

// CustomPolicyBindingLister helps list CustomPolicyBindings.
// All objects returned here must be treated as read-only.
type CustomPolicyBindingLister interface {
	// List lists all CustomPolicyBindings in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*auth.CustomPolicyBinding, err error)
	// CustomPolicyBindings returns an object that can list and get CustomPolicyBindings.
	CustomPolicyBindings(namespace string) CustomPolicyBindingNamespaceLister
	CustomPolicyBindingListerExpansion
}

// customPolicyBindingLister implements the CustomPolicyBindingLister interface.
type customPolicyBindingLister struct {
	indexer cache.Indexer
}

// NewCustomPolicyBindingLister returns a new CustomPolicyBindingLister.
func NewCustomPolicyBindingLister(indexer cache.Indexer) CustomPolicyBindingLister {
	return &customPolicyBindingLister{indexer: indexer}
}

// List lists all CustomPolicyBindings in the indexer.
func (s *customPolicyBindingLister) List(selector labels.Selector) (ret []*auth.CustomPolicyBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*auth.CustomPolicyBinding))
	})
	return ret, err
}

// CustomPolicyBindings returns an object that can list and get CustomPolicyBindings.
func (s *customPolicyBindingLister) CustomPolicyBindings(namespace string) CustomPolicyBindingNamespaceLister {
	return customPolicyBindingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CustomPolicyBindingNamespaceLister helps list and get CustomPolicyBindings.
// All objects returned here must be treated as read-only.
type CustomPolicyBindingNamespaceLister interface {
	// List lists all CustomPolicyBindings in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*auth.CustomPolicyBinding, err error)
	// Get retrieves the CustomPolicyBinding from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*auth.CustomPolicyBinding, error)
	CustomPolicyBindingNamespaceListerExpansion
}

// customPolicyBindingNamespaceLister implements the CustomPolicyBindingNamespaceLister
// interface.
type customPolicyBindingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CustomPolicyBindings in the indexer for a given namespace.
func (s customPolicyBindingNamespaceLister) List(selector labels.Selector) (ret []*auth.CustomPolicyBinding, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*auth.CustomPolicyBinding))
	})
	return ret, err
}

// Get retrieves the CustomPolicyBinding from the indexer for a given namespace and name.
func (s customPolicyBindingNamespaceLister) Get(name string) (*auth.CustomPolicyBinding, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(auth.Resource("custompolicybinding"), name)
	}
	return obj.(*auth.CustomPolicyBinding), nil
}
