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
	"tkestack.io/tke/api/business"
	businessv1 "tkestack.io/tke/api/business/v1"
	businessinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/business/internalversion"
	authversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/auth/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	"tkestack.io/tke/cmd/tke-business-api/app/options"
	"tkestack.io/tke/pkg/apiserver/storage"
	chartgroupstorage "tkestack.io/tke/pkg/business/registry/chartgroup/storage"
	configmapstorage "tkestack.io/tke/pkg/business/registry/configmap/storage"
	imagenamespacestorage "tkestack.io/tke/pkg/business/registry/imagenamespace/storage"
	namespacestorage "tkestack.io/tke/pkg/business/registry/namespace/storage"
	platformstorage "tkestack.io/tke/pkg/business/registry/platform/storage"
	portalstorage "tkestack.io/tke/pkg/business/registry/portal/storage"
	projectstorage "tkestack.io/tke/pkg/business/registry/project/storage"
)

// StorageProvider is a REST type for core resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
	PlatformClient       platformversionedclient.PlatformV1Interface
	RegistryClient       registryversionedclient.RegistryV1Interface
	AuthClient           authversionedclient.AuthV1Interface
	PrivilegedUsername   string
	Features             *options.FeatureOptions
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the APIGroupInfo
func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(business.GroupName, business.Scheme, business.ParameterCodec, business.Codecs)

	if apiResourceConfigSource.VersionEnabled(businessv1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[businessv1.SchemeGroupVersion.Version] =
			s.v1Storage(apiResourceConfigSource, restOptionsGetter, s.LoopbackClientConfig, s.Features)
	}

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return business.GroupName
}

func (s *StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config,
	features *options.FeatureOptions) map[string]rest.Storage {
	businessClient := businessinternalclient.NewForConfigOrDie(loopbackClientConfig)

	storageMap := make(map[string]rest.Storage)
	{
		projectREST := projectstorage.NewStorage(restOptionsGetter, businessClient, s.PlatformClient, s.PrivilegedUsername, features)
		storageMap["projects"] = projectREST.Project
		storageMap["projects/status"] = projectREST.Status
		storageMap["projects/finalize"] = projectREST.Finalize

		namespaceREST := namespacestorage.NewStorage(restOptionsGetter, businessClient, s.PlatformClient, s.PrivilegedUsername)
		storageMap["namespaces"] = namespaceREST.Namespace
		storageMap["namespaces/status"] = namespaceREST.Status

		platformREST := platformstorage.NewStorage(restOptionsGetter, businessClient, s.PrivilegedUsername)
		storageMap["platforms"] = platformREST.Platform

		portalREST := portalstorage.NewStorage(restOptionsGetter, businessClient, s.AuthClient)
		storageMap["portal"] = portalREST.Portal

		configMapREST := configmapstorage.NewStorage(restOptionsGetter)
		storageMap["configmaps"] = configMapREST.ConfigMap

		if s.RegistryClient != nil {
			imageNamespaceREST := imagenamespacestorage.NewStorage(restOptionsGetter, businessClient, s.RegistryClient, s.PrivilegedUsername)
			storageMap["imagenamespaces"] = imageNamespaceREST.ImageNamespace
			storageMap["imagenamespaces/status"] = imageNamespaceREST.Status
			storageMap["imagenamespaces/finalize"] = imageNamespaceREST.Finalize

			chartGroupREST := chartgroupstorage.NewStorage(restOptionsGetter, businessClient, s.RegistryClient, s.PrivilegedUsername)
			storageMap["chartgroups"] = chartGroupREST.ChartGroup
			storageMap["chartgroups/status"] = chartGroupREST.Status
			storageMap["chartgroups/finalize"] = chartGroupREST.Finalize
		}
	}

	return storageMap
}
