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
	notify "tkestack.io/tke/api/notify"
)

// ReceiverLister helps list Receivers.
// All objects returned here must be treated as read-only.
type ReceiverLister interface {
	// List lists all Receivers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*notify.Receiver, err error)
	// Get retrieves the Receiver from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*notify.Receiver, error)
	ReceiverListerExpansion
}

// receiverLister implements the ReceiverLister interface.
type receiverLister struct {
	indexer cache.Indexer
}

// NewReceiverLister returns a new ReceiverLister.
func NewReceiverLister(indexer cache.Indexer) ReceiverLister {
	return &receiverLister{indexer: indexer}
}

// List lists all Receivers in the indexer.
func (s *receiverLister) List(selector labels.Selector) (ret []*notify.Receiver, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*notify.Receiver))
	})
	return ret, err
}

// Get retrieves the Receiver from the index for a given name.
func (s *receiverLister) Get(name string) (*notify.Receiver, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(notify.Resource("receiver"), name)
	}
	return obj.(*notify.Receiver), nil
}
