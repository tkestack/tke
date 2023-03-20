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
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	"tkestack.io/tke/api/application"
	applicationv1 "tkestack.io/tke/api/application/v1"
	applicationinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/application/internalversion"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	registryversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/registry/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	appconfig "tkestack.io/tke/pkg/application/config"
	applicationstorage "tkestack.io/tke/pkg/application/registry/application/storage"
	configmapstorage "tkestack.io/tke/pkg/application/registry/configmap/storage"
)

// StorageProvider is a REST type for core resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
	PlatformClient       platformversionedclient.PlatformV1Interface
	RegistryClient       registryversionedclient.RegistryV1Interface
	Authorizer           authorizer.Authorizer
	RepoConfiguration    appconfig.RepoConfiguration
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the APIGroupInfo
func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(application.GroupName, application.Scheme, application.ParameterCodec, application.Codecs)

	apiGroupInfo.VersionedResourcesStorageMap[applicationv1.SchemeGroupVersion.Version] =
		s.v1Storage(apiResourceConfigSource, restOptionsGetter, s.LoopbackClientConfig)

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return application.GroupName
}

func (s *StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	applicationClient := applicationinternalclient.NewForConfigOrDie(loopbackClientConfig)
	storageMap := make(map[string]rest.Storage)
	{
		configMapREST := configmapstorage.NewStorage(restOptionsGetter)
		storageMap["configmaps"] = configMapREST.ConfigMap

		appRESTStorage := applicationstorage.NewStorage(restOptionsGetter, applicationClient)
		appREST := applicationstorage.NewREST(appRESTStorage.App,
			applicationClient,
			s.PlatformClient,
			s.RegistryClient,
			s.Authorizer,
			s.RepoConfiguration)
		appHistoryREST := applicationstorage.NewHistoryREST(appRESTStorage.App, applicationClient, s.PlatformClient)
		appResourceREST := applicationstorage.NewResourceREST(appRESTStorage.App, applicationClient, s.PlatformClient)
		appRollbackREST := applicationstorage.NewRollbackREST(appRESTStorage.App, applicationClient, s.PlatformClient)
		storageMap["apps"] = appREST
		storageMap["apps/histories"] = appHistoryREST
		storageMap["apps/resources"] = appResourceREST
		storageMap["apps/rollback"] = appRollbackREST
		storageMap["apps/status"] = appRESTStorage.Status
		storageMap["apps/finalize"] = appRESTStorage.Finalize
	}

	return storageMap
}
