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
	"k8s.io/api/storage/v1"
	"k8s.io/api/storage/v1alpha1"
	"k8s.io/api/storage/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/storage"
	storageclassstorage "tkestack.io/tke/pkg/platform/proxy/storage/storageclass/storage"
	volumeattachmentstorage "tkestack.io/tke/pkg/platform/proxy/storage/volumeattachment/storage"
)

// StorageProvider is a REST type for oauth resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the
// APIGroupInfo
func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(v1.GroupName, platform.Scheme, metav1.ParameterCodec, platform.Codecs)

	if apiResourceConfigSource.VersionEnabled(v1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = s.v1Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	if apiResourceConfigSource.VersionEnabled(v1alpha1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[v1alpha1.SchemeGroupVersion.Version] = s.v1Alpha1Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	if apiResourceConfigSource.VersionEnabled(v1beta1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[v1beta1.SchemeGroupVersion.Version] = s.v1Beta1Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return v1.GroupName
}

func (s *StorageProvider) v1Alpha1Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	volumeAttachmentStore := volumeattachmentstorage.NewStorageV1Alpha1(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"volumeattachments": volumeAttachmentStore.VolumeAttachment,
	}

	return storageMap
}

func (s *StorageProvider) v1Beta1Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	volumeAttachmentStore := volumeattachmentstorage.NewStorageV1Beta1(restOptionsGetter, platformClient)
	storageClassStore := storageclassstorage.NewStorageV1Beta1(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"volumeattachments":     volumeAttachmentStore.VolumeAttachment,
		"storageclasses":        storageClassStore.StorageClass,
		"storageclasses/events": storageClassStore.Events,
	}

	return storageMap
}

func (s *StorageProvider) v1Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	storageClassStore := storageclassstorage.NewStorageV1(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"storageclasses":        storageClassStore.StorageClass,
		"storageclasses/events": storageClassStore.Events,
	}

	return storageMap
}
