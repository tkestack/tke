/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/discovery"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/platform/proxy"
	"tkestack.io/tke/pkg/platform/util"
)

// APIResourcesREST implement bucket call interface for cluster.
type APIResourcesREST struct {
	rest.Storage
	store          *registry.Store
	platformClient platforminternalclient.PlatformInterface
}

func (r *APIResourcesREST) New() runtime.Object {
	return &platform.ClusterGroupAPIResourceItems{}
}

func (r *APIResourcesREST) NewGetOptions() (runtime.Object, bool, string) {
	return &platform.ClusterGroupAPIResourceOptions{}, false, ""
}

// Get finds a resource in the storage by name and returns it.
func (r *APIResourcesREST) Get(ctx context.Context, clusterName string, options runtime.Object) (runtime.Object, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	cluster := clusterObject.(*platform.Cluster)
	if err := util.FilterCluster(ctx, cluster); err != nil {
		return nil, err
	}

	config, err := proxy.GetConfig(ctx, r.platformClient)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	discoveryclient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	lists, err := discoveryclient.ServerPreferredResources()
	failedGroup := ""
	if err != nil {
		failedGroup = err.Error()
	}
	items := make([]platform.ClusterGroupAPIResourceItems, 0)
	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}
		apiResources := make([]platform.ClusterGroupAPIResourceItem, 0)
		for _, groupResource := range list.APIResources {
			apiResources = append(apiResources, platform.ClusterGroupAPIResourceItem{
				Name:         groupResource.Name,
				SingularName: groupResource.SingularName,
				Namespaced:   groupResource.Namespaced,
				Group:        groupResource.Group,
				Version:      groupResource.Version,
				Kind:         groupResource.Kind,
				Verbs:        groupResource.Verbs,
				ShortNames:   groupResource.ShortNames,
				Categories:   groupResource.Categories,
			})
		}
		items = append(items, platform.ClusterGroupAPIResourceItems{
			GroupVersion: list.GroupVersion,
			APIResources: apiResources,
		})
	}
	return &platform.ClusterGroupAPIResourceItemsList{
		Items:            items,
		FailedGroupError: failedGroup,
	}, nil
}
