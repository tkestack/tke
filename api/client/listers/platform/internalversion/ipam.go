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

// IPAMLister helps list IPAMs.
type IPAMLister interface {
	// List lists all IPAMs in the indexer.
	List(selector labels.Selector) (ret []*platform.IPAM, err error)
	// Get retrieves the IPAM from the index for a given name.
	Get(name string) (*platform.IPAM, error)
	IPAMListerExpansion
}

// iPAMLister implements the IPAMLister interface.
type iPAMLister struct {
	indexer cache.Indexer
}

// NewIPAMLister returns a new IPAMLister.
func NewIPAMLister(indexer cache.Indexer) IPAMLister {
	return &iPAMLister{indexer: indexer}
}

// List lists all IPAMs in the indexer.
func (s *iPAMLister) List(selector labels.Selector) (ret []*platform.IPAM, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*platform.IPAM))
	})
	return ret, err
}

// Get retrieves the IPAM from the index for a given name.
func (s *iPAMLister) Get(name string) (*platform.IPAM, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(platform.Resource("ipam"), name)
	}
	return obj.(*platform.IPAM), nil
}
