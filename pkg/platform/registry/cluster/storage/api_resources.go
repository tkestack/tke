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
	"net/http"
	"strings"
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
	return &platform.ClusterAPIResourcesOptions{}, false, ""
}

// Get finds a resource in the storage by name and returns it.
func (r *APIResourcesREST) Get(ctx context.Context, clusterName string, options runtime.Object) (runtime.Object, error) {
	opts := options.(*platform.ClusterAPIResourcesOptions)
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
	// userName, tenantID := authentication.UsernameAndTenantID(ctx)
	// config.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
	// 	return &headerAdder{
	// 		headers: map[string]string{"X-Remote-Extra-TenantID":tenantID,"X-Remote-User":userName},
	// 		rt:      rt,
	// 	}
	// }
	discoveryclient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	lists, err := discoveryclient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	items := make([]platform.ClusterGroupAPIResourceItems, 0)
	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}
		apiResources := make([]platform.ClusterGroupAPIResourceItem, 0)
		for _, groupResource := range list.APIResources {
			if allowGroupVersion(list.GroupVersion, opts.OnlySecure) {
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
		}
		items = append(items, platform.ClusterGroupAPIResourceItems{
			GroupVersion: list.GroupVersion,
			APIResources: apiResources,
		})
	}
	return &platform.ClusterGroupAPIResourceItemsList{
		Items: items,
	}, nil
}

type headerAdder struct {
	headers map[string]string
	rt      http.RoundTripper
}

func (h *headerAdder) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	return h.rt.RoundTrip(req)
}

func allowGroupVersion(groupVersion string, onlySecure bool) bool {
	if !onlySecure {
		return true
	}
	if groupVersion == "v1" {
		return true
	}
	gvs := strings.Split(groupVersion, "/")
	if len(gvs) < 2 {
		return false
	}
	group := gvs[0]
	if group == "apps" ||
		group == "autoscaling" ||
		group == "batch" ||
		group == "extensions" ||
		group == "networking.k8s.io" ||
		group == "policy" ||
		group == "scheduling.k8s.io" ||
		group == "settings.k8s.io" ||
		group == "storage.k8s.io" {
		return true
	}
	if strings.HasSuffix(group, "istio.io") ||
		strings.HasSuffix(group, "tke.cloud.tencent.com") ||
		strings.HasSuffix(group, "tkestack.io") {
		return true
	}
	return false
}
