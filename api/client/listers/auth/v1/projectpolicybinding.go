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
	v1 "tkestack.io/tke/api/auth/v1"
)

// ProjectPolicyBindingLister helps list ProjectPolicyBindings.
// All objects returned here must be treated as read-only.
type ProjectPolicyBindingLister interface {
	// List lists all ProjectPolicyBindings in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ProjectPolicyBinding, err error)
	// Get retrieves the ProjectPolicyBinding from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.ProjectPolicyBinding, error)
	ProjectPolicyBindingListerExpansion
}

// projectPolicyBindingLister implements the ProjectPolicyBindingLister interface.
type projectPolicyBindingLister struct {
	indexer cache.Indexer
}

// NewProjectPolicyBindingLister returns a new ProjectPolicyBindingLister.
func NewProjectPolicyBindingLister(indexer cache.Indexer) ProjectPolicyBindingLister {
	return &projectPolicyBindingLister{indexer: indexer}
}

// List lists all ProjectPolicyBindings in the indexer.
func (s *projectPolicyBindingLister) List(selector labels.Selector) (ret []*v1.ProjectPolicyBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ProjectPolicyBinding))
	})
	return ret, err
}

// Get retrieves the ProjectPolicyBinding from the index for a given name.
func (s *projectPolicyBindingLister) Get(name string) (*v1.ProjectPolicyBinding, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("projectpolicybinding"), name)
	}
	return obj.(*v1.ProjectPolicyBinding), nil
}
