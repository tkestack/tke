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
	v1 "tkestack.io/tke/api/platform/v1"
)

// TappControllerLister helps list TappControllers.
type TappControllerLister interface {
	// List lists all TappControllers in the indexer.
	List(selector labels.Selector) (ret []*v1.TappController, err error)
	// Get retrieves the TappController from the index for a given name.
	Get(name string) (*v1.TappController, error)
	TappControllerListerExpansion
}

// tappControllerLister implements the TappControllerLister interface.
type tappControllerLister struct {
	indexer cache.Indexer
}

// NewTappControllerLister returns a new TappControllerLister.
func NewTappControllerLister(indexer cache.Indexer) TappControllerLister {
	return &tappControllerLister{indexer: indexer}
}

// List lists all TappControllers in the indexer.
func (s *tappControllerLister) List(selector labels.Selector) (ret []*v1.TappController, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.TappController))
	})
	return ret, err
}

// Get retrieves the TappController from the index for a given name.
func (s *tappControllerLister) Get(name string) (*v1.TappController, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("tappcontroller"), name)
	}
	return obj.(*v1.TappController), nil
}
