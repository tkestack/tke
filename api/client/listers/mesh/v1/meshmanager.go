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
	v1 "tkestack.io/tke/api/mesh/v1"
)

// MeshManagerLister helps list MeshManagers.
// All objects returned here must be treated as read-only.
type MeshManagerLister interface {
	// List lists all MeshManagers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.MeshManager, err error)
	// Get retrieves the MeshManager from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.MeshManager, error)
	MeshManagerListerExpansion
}

// meshManagerLister implements the MeshManagerLister interface.
type meshManagerLister struct {
	indexer cache.Indexer
}

// NewMeshManagerLister returns a new MeshManagerLister.
func NewMeshManagerLister(indexer cache.Indexer) MeshManagerLister {
	return &meshManagerLister{indexer: indexer}
}

// List lists all MeshManagers in the indexer.
func (s *meshManagerLister) List(selector labels.Selector) (ret []*v1.MeshManager, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.MeshManager))
	})
	return ret, err
}

// Get retrieves the MeshManager from the index for a given name.
func (s *meshManagerLister) Get(name string) (*v1.MeshManager, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("meshmanager"), name)
	}
	return obj.(*v1.MeshManager), nil
}
