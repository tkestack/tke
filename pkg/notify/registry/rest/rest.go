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
	notifyinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/notify/internalversion"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	channelstorage "tkestack.io/tke/pkg/notify/registry/channel/storage"
	configmapstorage "tkestack.io/tke/pkg/notify/registry/configmap/storage"
	messagestorage "tkestack.io/tke/pkg/notify/registry/message/storage"
	messagerequeststorage "tkestack.io/tke/pkg/notify/registry/messagerequest/storage"
	receiverstorage "tkestack.io/tke/pkg/notify/registry/receiver/storage"
	receivergroupstorage "tkestack.io/tke/pkg/notify/registry/receivergroup/storage"
	templatestorage "tkestack.io/tke/pkg/notify/registry/template/storage"
)

// StorageProvider is a REST type for core resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
	PrivilegedUsername   string
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the APIGroupInfo
func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(notify.GroupName, notify.Scheme, notify.ParameterCodec, notify.Codecs)

	if apiResourceConfigSource.VersionEnabled(v1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = s.v1Storage(apiResourceConfigSource, restOptionsGetter, s.LoopbackClientConfig)
	}

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return notify.GroupName
}

func (s *StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	notifyClient := notifyinternalclient.NewForConfigOrDie(loopbackClientConfig)

	storageMap := make(map[string]rest.Storage)
	{
		channelREST := channelstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["channels"] = channelREST.Channel
		storageMap["channels/status"] = channelREST.Status
		storageMap["channels/finalize"] = channelREST.Finalize

		templateREST := templatestorage.NewStorage(restOptionsGetter, notifyClient, s.PrivilegedUsername)
		storageMap["templates"] = templateREST.Template

		configMapREST := configmapstorage.NewStorage(restOptionsGetter)
		storageMap["configmaps"] = configMapREST.ConfigMap

		messageREST := messagestorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["messages"] = messageREST.Message
		storageMap["messages/status"] = messageREST.Status

		messageRequestREST := messagerequeststorage.NewStorage(restOptionsGetter, notifyClient, s.PrivilegedUsername)
		storageMap["messagerequests"] = messageRequestREST.MessageRequest
		storageMap["messagerequests/status"] = messageRequestREST.Status

		receiverREST := receiverstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["receivers"] = receiverREST.Receiver

		receiverGroupREST := receivergroupstorage.NewStorage(restOptionsGetter, notifyClient, s.PrivilegedUsername)
		storageMap["receivergroups"] = receiverGroupREST.ReceiverGroup
	}

	return storageMap
}
