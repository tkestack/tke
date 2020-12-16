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

	"github.com/thoas/go-funk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/registry/clusteraddontype"
	"tkestack.io/tke/pkg/platform/util"
)

// AddonTypeREST implements the REST endpoint.
type AddonTypeREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

var _ = rest.Getter(&AddonTypeREST{})

// New returns an empty object that can be used with Create after request data
// has been put into it.
func (r *AddonTypeREST) New() runtime.Object {
	return &platform.ClusterAddonTypeList{}
}

// Get finds a resource in the storage by name and returns it.
func (r *AddonTypeREST) Get(ctx context.Context, clusterName string, options *metav1.GetOptions) (runtime.Object, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, options)
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}

	l := &platform.ClusterAddonTypeList{
		Items: make([]platform.ClusterAddonType, 0),
	}

	for _, v := range clusteraddontype.Types {
		if funk.ContainsString(v.CompatibleClusterType, cluster.Spec.Type) {
			l.Items = append(l.Items, v)
		}
	}
	return l, nil
}
