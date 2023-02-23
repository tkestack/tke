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
	"tkestack.io/tke/api/authz"
	authzv1 "tkestack.io/tke/api/authz/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	configmapstorage "tkestack.io/tke/pkg/authz/registry/configmap/storage"
	mcrbstorage "tkestack.io/tke/pkg/authz/registry/multiclusterrolebinding/storage"
	policystorage "tkestack.io/tke/pkg/authz/registry/policy/storage"
	rolestorage "tkestack.io/tke/pkg/authz/registry/role/storage"
)

// StorageProvider is a REST type for core resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
	Authorizer           authorizer.Authorizer
	PlatformClient       platformversionedclient.PlatformV1Interface
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the APIGroupInfo
func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(authz.GroupName, authz.Scheme, authz.ParameterCodec, authz.Codecs)
	apiGroupInfo.VersionedResourcesStorageMap[authzv1.SchemeGroupVersion.Version] =
		s.v1Storage(apiResourceConfigSource, restOptionsGetter, s.PlatformClient)
	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return authz.GroupName
}

func (s *StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, platformClient platformversionedclient.PlatformV1Interface) map[string]rest.Storage {
	storageMap := make(map[string]rest.Storage)
	{
		configmapREST := configmapstorage.NewStorage(restOptionsGetter)
		policyREST := policystorage.NewStorage(restOptionsGetter, platformClient)
		rolestorageREST := rolestorage.NewStorage(restOptionsGetter, policyREST.Policy, platformClient)
		mcrbREST := mcrbstorage.NewStorage(restOptionsGetter, rolestorageREST.Role, platformClient)

		storageMap["policies"] = policyREST.Policy
		storageMap["roles"] = rolestorageREST.Role
		storageMap["roles/finalize"] = rolestorageREST.Finalize
		storageMap["multiclusterrolebindings"] = mcrbREST.MultiClusterRoleBinding
		storageMap["multiclusterrolebindings/status"] = mcrbREST.Status
		storageMap["multiclusterrolebindings/finalize"] = mcrbREST.Finalize
		storageMap["configmaps"] = configmapREST.ConfigMap
	}
	return storageMap
}
