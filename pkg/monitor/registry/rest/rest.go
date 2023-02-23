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

package rest

import (
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/monitor"
	v1 "tkestack.io/tke/api/monitor/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	configmapstorage "tkestack.io/tke/pkg/monitor/registry/configmap/storage"
	metricstorage "tkestack.io/tke/pkg/monitor/registry/metric/storage"
	clusteroverview "tkestack.io/tke/pkg/monitor/registry/overview/cluster/storage"
	promstorage "tkestack.io/tke/pkg/monitor/registry/prometheus/storage"
	monitorstorage "tkestack.io/tke/pkg/monitor/storage"
	"tkestack.io/tke/pkg/monitor/util/cache"
)

// StorageProvider is a REST type for core resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
	PrivilegedUsername   string
	MetricStorage        monitorstorage.MetricStorage
	Cacher               cache.Cacher
	PlatformClient       platformversionedclient.PlatformV1Interface
	BusinessClient       businessversionedclient.BusinessV1Interface
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the APIGroupInfo
func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(monitor.GroupName, monitor.Scheme, monitor.ParameterCodec, monitor.Codecs)

	apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = s.v1Storage(apiResourceConfigSource, restOptionsGetter, s.LoopbackClientConfig)

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return monitor.GroupName
}

func (s *StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	storageMap := make(map[string]rest.Storage)
	{

		configMapREST := configmapstorage.NewStorage(restOptionsGetter)
		storageMap["configmaps"] = configMapREST.ConfigMap

		metricREST := metricstorage.NewStorage(restOptionsGetter, s.MetricStorage)
		storageMap["metrics"] = metricREST.Metric

		clusterOverviewREST := clusteroverview.NewStorage(restOptionsGetter, s.PlatformClient, s.BusinessClient, s.Cacher)
		storageMap["clusteroverviews"] = clusterOverviewREST.ClusterOverview

		promREST := promstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["prometheuses"] = promREST.Prometheus
		storageMap["prometheuses/status"] = promREST.Status
	}

	return storageMap
}
