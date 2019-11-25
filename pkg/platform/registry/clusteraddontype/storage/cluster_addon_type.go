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
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/registry/clusteraddontype"
)

// Storage includes storage for clusters and all sub resources.
type Storage struct {
	ClusterAddonType *REST
}

// NewStorage returns a Storage object that will work against clusters.
func NewStorage(_ genericregistry.RESTOptionsGetter) *Storage {
	return &Storage{
		ClusterAddonType: &REST{},
	}
}

// REST implements a RESTStorage for clusters against etcd.
type REST struct {
	rest.Storage
}

var _ rest.ShortNamesProvider = &REST{}
var _ rest.Lister = &REST{}
var _ rest.Scoper = &REST{}

// ShortNames implements the ShortNamesProvider interface. Returns a list of short names for a resource.
func (r *REST) ShortNames() []string {
	return []string{"cat"}
}

// NamespaceScoped returns true if the storage is namespaced
func (r *REST) NamespaceScoped() bool {
	return false
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *REST) New() runtime.Object {
	return &platform.ClusterAddonType{}
}

// Get finds a resource in the storage by name and returns it.
func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	var ct *platform.ClusterAddonType
	for k, v := range clusteraddontype.Types {
		if strings.EqualFold(string(k), name) {
			ct = &platform.ClusterAddonType{
				ObjectMeta: metav1.ObjectMeta{
					Name: strings.ToLower(string(k)),
				},
				Type:          string(k),
				Level:         v.Level,
				LatestVersion: v.LatestVersion,
			}
		}
	}
	if ct == nil {
		return nil, errors.NewNotFound(platform.Resource("clusteraddontype"), name)
	}
	return ct, nil
}

// NewList returns an empty object that can be used with the List call.
func (r *REST) NewList() runtime.Object {
	return &platform.ClusterAddonTypeList{}
}

// List selects resources in the storage which match to the selector. 'options' can be nil.
func (r *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	l := &platform.ClusterAddonTypeList{
		Items: make([]platform.ClusterAddonType, len(clusteraddontype.Types)),
	}
	i := 0
	for k, v := range clusteraddontype.Types {
		l.Items[i] = platform.ClusterAddonType{
			ObjectMeta: metav1.ObjectMeta{
				Name: strings.ToLower(string(k)),
			},
			Type:          string(k),
			Level:         v.Level,
			LatestVersion: v.LatestVersion,
			Description:   v.Description,
		}
		i++
	}
	return l, nil
}
