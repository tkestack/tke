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
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/logagent"
	v1 "tkestack.io/tke/api/logagent/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	configmapstorage "tkestack.io/tke/pkg/logagent/registry/configmap/storage"
	logagentstorage "tkestack.io/tke/pkg/logagent/registry/logagent/storage"
)

// StorageProvider is a REST type for core resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
	PrivilegedUsername   string
	PlatformClient       platformversionedclient.PlatformV1Interface //used by structs like logfile tree to get cluster client and then communicate with clusters
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(logagent.GroupName, logagent.Scheme, logagent.ParameterCodec, logagent.Codecs)

	apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = s.v1Storage(apiResourceConfigSource, restOptionsGetter, s.LoopbackClientConfig)

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return logagent.GroupName
}

func (s *StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	//do we need client??
	storageMap := make(map[string]rest.Storage)
	{
		logagentRest := logagentstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername, s.PlatformClient)
		storageMap["logagents"] = logagentRest.LogAgent
		storageMap["logagents/status"] = logagentRest.Status
		storageMap["logagents/filetree"] = logagentRest.LogFileTree
		storageMap["logagents/filecontent"] = logagentRest.LogFileContent
		storageMap["logagents/logcollector"] = logagentRest.LogagentProxy
		storageMap["logagents/filedownload"] = logagentRest.LogfileProxy
		storageMap["logagents/esdetection"] = logagentRest.LogESDetection
		configMapREST := configmapstorage.NewStorage(restOptionsGetter)
		storageMap["configmaps"] = configMapREST.ConfigMap
	}
	return storageMap
}
