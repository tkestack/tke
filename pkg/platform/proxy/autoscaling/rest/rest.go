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
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/storage"
	hpastorage "tkestack.io/tke/pkg/platform/proxy/autoscaling/horizontalpodautoscaler/storage"
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
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(autoscalingv1.GroupName, platform.Scheme, metav1.ParameterCodec, platform.Codecs)

	if apiResourceConfigSource.VersionEnabled(autoscalingv1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[autoscalingv1.SchemeGroupVersion.Version] = s.v1Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	if apiResourceConfigSource.VersionEnabled(autoscalingv2beta1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[autoscalingv2beta1.SchemeGroupVersion.Version] = s.v2Beta1Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	if apiResourceConfigSource.VersionEnabled(autoscalingv2beta2.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[autoscalingv2beta2.SchemeGroupVersion.Version] = s.v2beta2Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return autoscalingv1.GroupName
}

func (s *StorageProvider) v1Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	hpaStore := hpastorage.NewStorageV1(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"horizontalpodautoscalers":        hpaStore.HorizontalPodAutoscaler,
		"horizontalpodautoscalers/status": hpaStore.Status,
		"horizontalpodautoscalers/events": hpaStore.Events,
	}

	return storageMap
}

func (s *StorageProvider) v2Beta1Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	hpaStore := hpastorage.NewStorageV2Beta1(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"horizontalpodautoscalers":        hpaStore.HorizontalPodAutoscaler,
		"horizontalpodautoscalers/status": hpaStore.Status,
		"horizontalpodautoscalers/events": hpaStore.Events,
	}

	return storageMap
}

func (s *StorageProvider) v2beta2Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	hpaStore := hpastorage.NewStorageV2Beta2(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"horizontalpodautoscalers":        hpaStore.HorizontalPodAutoscaler,
		"horizontalpodautoscalers/status": hpaStore.Status,
		"horizontalpodautoscalers/events": hpaStore.Events,
	}

	return storageMap
}
