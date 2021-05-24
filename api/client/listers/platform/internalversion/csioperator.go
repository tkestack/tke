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
	platform "tkestack.io/tke/api/platform"
)

// CSIOperatorLister helps list CSIOperators.
// All objects returned here must be treated as read-only.
type CSIOperatorLister interface {
	// List lists all CSIOperators in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*platform.CSIOperator, err error)
	// Get retrieves the CSIOperator from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*platform.CSIOperator, error)
	CSIOperatorListerExpansion
}

// cSIOperatorLister implements the CSIOperatorLister interface.
type cSIOperatorLister struct {
	indexer cache.Indexer
}

// NewCSIOperatorLister returns a new CSIOperatorLister.
func NewCSIOperatorLister(indexer cache.Indexer) CSIOperatorLister {
	return &cSIOperatorLister{indexer: indexer}
}

// List lists all CSIOperators in the indexer.
func (s *cSIOperatorLister) List(selector labels.Selector) (ret []*platform.CSIOperator, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*platform.CSIOperator))
	})
	return ret, err
}

// Get retrieves the CSIOperator from the index for a given name.
func (s *cSIOperatorLister) Get(name string) (*platform.CSIOperator, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(platform.Resource("csioperator"), name)
	}
	return obj.(*platform.CSIOperator), nil
}
