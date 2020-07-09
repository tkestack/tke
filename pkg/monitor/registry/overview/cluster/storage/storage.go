/*
 * Tencent is pleased to support the open source community by making TKEStack available.
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

package storage

import (
	"context"
	"fmt"

	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/monitor"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/monitor/util/cache"
	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericregistry "k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
)

// Storage includes storage for metrics and all sub resources.
type Storage struct {
	ClusterOverview *REST
}

// NewStorage returns a Storage object that will work against metrics.
func NewStorage(_ genericregistry.RESTOptionsGetter, platformClient platformversionedclient.PlatformV1Interface, cacher cache.Cacher) *Storage {
	log.Info("ClusterOverview NewStorage")
	return &Storage{
		ClusterOverview: &REST{
			platformClient: platformClient,
			cacher:         cacher,
		},
	}
}

// REST implements a RESTStorage for metrics against etcd.
type REST struct {
	rest.Storage
	platformClient platformversionedclient.PlatformV1Interface
	cacher         cache.Cacher
}

var _ rest.Creater = &REST{}
var _ rest.Scoper = &REST{}

// NamespaceScoped returns true if the storage is namespaced
func (r *REST) NamespaceScoped() bool {
	return false
}

// New returns an empty object that can be used with Create and Update after request data has been put into it.
func (r *REST) New() runtime.Object {
	return &monitor.ClusterOverview{}
}

// Create creates a new version of a resource.
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	clusterOverview, ok := obj.(*monitor.ClusterOverview)
	if !ok {
		return nil, errors.NewBadRequest("failed to processed request body")
	}
	_, tenantID := authentication.UsernameAndTenantID(ctx)
	listOptions := metav1.ListOptions{}
	if tenantID != "" {
		listOptions.FieldSelector = fmt.Sprintf("spec.tenantID=%s", tenantID)
	}
	clusterList, _ := r.platformClient.Clusters().List(ctx, listOptions)
	clusterIDs := make([]string, 0)
	for _, cls := range clusterList.Items {
		clusterIDs = append(clusterIDs, cls.GetName())
	}
	clusterOverview.Result = r.cacher.GetClusterOverviewResult(clusterIDs)
	return clusterOverview, nil
}
