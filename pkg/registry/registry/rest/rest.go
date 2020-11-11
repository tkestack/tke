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
	"encoding/base64"
	"fmt"
	"net/http"

	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	registryinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/registry/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	businessversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/registry"
	v1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	harbor "tkestack.io/tke/pkg/registry/harbor/client"
	chartstorage "tkestack.io/tke/pkg/registry/registry/chart/storage"
	chartgroupstorage "tkestack.io/tke/pkg/registry/registry/chartgroup/storage"
	configmapstorage "tkestack.io/tke/pkg/registry/registry/configmap/storage"
	namespacestorage "tkestack.io/tke/pkg/registry/registry/namespace/storage"
	repositorystorage "tkestack.io/tke/pkg/registry/registry/repository/storage"
	"tkestack.io/tke/pkg/util/transport"
)

// StorageProvider is a REST type for core resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
	ExternalScheme       string
	ExternalHost         string
	ExternalPort         int
	ExternalCAFile       string
	PrivilegedUsername   string
	AuthClient           authversionedclient.AuthV1Interface
	BusinessClient       businessversionedclient.BusinessV1Interface
	PlatformClient       platformversionedclient.PlatformV1Interface
	RegistryConfig       *registryconfig.RegistryConfiguration
	Authorizer           authorizer.Authorizer
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the APIGroupInfo
func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(registry.GroupName, registry.Scheme, registry.ParameterCodec, registry.Codecs)

	if apiResourceConfigSource.VersionEnabled(v1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = s.v1Storage(apiResourceConfigSource, restOptionsGetter, s.LoopbackClientConfig)
	}

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return registry.GroupName
}

func (s *StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	registryClient := registryinternalclient.NewForConfigOrDie(loopbackClientConfig)

	var harborClient *harbor.APIClient = nil

	if s.RegistryConfig.HarborEnabled {
		tr, _ := transport.NewOneWayTLSTransport(s.RegistryConfig.HarborCAFile, true)
		headers := make(map[string]string)
		headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(
			s.RegistryConfig.Security.AdminUsername+":"+s.RegistryConfig.Security.AdminPassword),
		)
		cfg := &harbor.Configuration{
			BasePath:      fmt.Sprintf("https://%s/api/v2.0", s.RegistryConfig.DomainSuffix),
			DefaultHeader: headers,
			UserAgent:     "Swagger-Codegen/1.0.0/go",
			HTTPClient: &http.Client{
				Transport: tr,
			},
		}
		harborClient = harbor.NewAPIClient(cfg)
	}

	storageMap := make(map[string]rest.Storage)
	{

		configMapREST := configmapstorage.NewStorage(restOptionsGetter)
		storageMap["configmaps"] = configMapREST.ConfigMap

		namespaceREST := namespacestorage.NewStorage(restOptionsGetter, registryClient, s.PrivilegedUsername, harborClient)
		storageMap["namespaces"] = namespaceREST.Namespace
		storageMap["namespaces/status"] = namespaceREST.Status

		repositoryREST := repositorystorage.NewStorage(restOptionsGetter, registryClient, s.PrivilegedUsername, harborClient)
		storageMap["repositories"] = repositoryREST.Repository
		storageMap["repositories/status"] = repositoryREST.Status

		chartGroupRESTStorage := chartgroupstorage.NewStorage(restOptionsGetter, registryClient, s.AuthClient, s.BusinessClient, s.PrivilegedUsername)
		chartGroupREST := chartgroupstorage.NewREST(chartGroupRESTStorage.ChartGroup, registryClient, s.AuthClient)
		storageMap["chartgroups"] = chartGroupREST
		storageMap["chartgroups/status"] = chartGroupRESTStorage.Status
		storageMap["chartgroups/finalize"] = chartGroupRESTStorage.Finalize

		chartREST := chartstorage.NewStorage(restOptionsGetter, registryClient, s.AuthClient, s.BusinessClient, s.PrivilegedUsername)
		chartVersionREST := chartstorage.NewVersionREST(chartREST.Chart, s.PlatformClient, registryClient, s.RegistryConfig,
			s.ExternalScheme,
			s.ExternalHost,
			s.ExternalPort,
			s.ExternalCAFile,
			s.Authorizer)
		storageMap["charts"] = chartREST.Chart
		storageMap["charts/status"] = chartREST.Status
		storageMap["charts/finalize"] = chartREST.Finalize
		storageMap["charts/version"] = chartVersionREST
	}

	return storageMap
}
