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
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/storage"
	controllerrevisionstorage "tkestack.io/tke/pkg/platform/proxy/apps/controllerrevision/storage"
	daemonsetstorage "tkestack.io/tke/pkg/platform/proxy/apps/daemonset/storage"
	deploymentstorage "tkestack.io/tke/pkg/platform/proxy/apps/deployment/storage"
	replicasetstorage "tkestack.io/tke/pkg/platform/proxy/apps/replicaset/storage"
	statefulsetstorage "tkestack.io/tke/pkg/platform/proxy/apps/statefulset/storage"
)

// StorageProvider is a REST type for oauth resources storage that implement
// RestStorageProvider interface
type StorageProvider struct {
	LoopbackClientConfig *restclient.Config
}

// Implement RESTStorageProvider
var _ storage.RESTStorageProvider = &StorageProvider{}

// NewRESTStorage is a factory constructor to creates and returns the APIGroupInfo
func (s *StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericserver.APIGroupInfo, bool) {
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(appsv1.GroupName, platform.Scheme, metav1.ParameterCodec, platform.Codecs)

	if apiResourceConfigSource.VersionEnabled(appsv1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[appsv1.SchemeGroupVersion.Version] = s.v1Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	if apiResourceConfigSource.VersionEnabled(appsv1beta1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[appsv1beta1.SchemeGroupVersion.Version] = s.v1Beta1Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	if apiResourceConfigSource.VersionEnabled(appsv1beta2.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[appsv1beta2.SchemeGroupVersion.Version] = s.v1Beta2Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return appsv1.GroupName
}

func (s *StorageProvider) v1Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	controllerRevisionStore := controllerrevisionstorage.NewStorageV1(restOptionsGetter, platformClient)
	daemonSetStore := daemonsetstorage.NewStorageV1(restOptionsGetter, platformClient)
	deploymentStore := deploymentstorage.NewStorageV1(restOptionsGetter, platformClient)
	statefulSetStore := statefulsetstorage.NewStorageV1(restOptionsGetter, platformClient)
	replicaSetStore := replicasetstorage.NewStorageV1(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"controllerrevisions":                  controllerRevisionStore.ControllerRevision,
		"daemonsets":                           daemonSetStore.DaemonSet,
		"daemonsets/pods":                      daemonSetStore.Pods,
		"daemonsets/status":                    daemonSetStore.Status,
		"daemonsets/events":                    daemonSetStore.Events,
		"deployments":                          deploymentStore.Deployment,
		"deployments/status":                   deploymentStore.Status,
		"deployments/scale":                    deploymentStore.Scale,
		"deployments/rollback":                 deploymentStore.RolloutUndo,
		"deployments/pods":                     deploymentStore.Pods,
		"deployments/events":                   deploymentStore.Events,
		"deployments/horizontalpodautoscalers": deploymentStore.HPAs,
		"statefulsets":                         statefulSetStore.StatefulSet,
		"statefulsets/status":                  statefulSetStore.Status,
		"statefulsets/pods":                    statefulSetStore.Pods,
		"statefulsets/scale":                   statefulSetStore.Scale,
		"statefulsets/events":                  statefulSetStore.Events,
		"replicasets":                          replicaSetStore.ReplicaSet,
		"replicasets/pods":                     replicaSetStore.Pods,
		"replicasets/status":                   replicaSetStore.Status,
		"replicasets/scale":                    replicaSetStore.Scale,
		"replicasets/events":                   replicaSetStore.Events,
	}

	return storageMap
}

func (s *StorageProvider) v1Beta1Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	controllerRevisionStore := controllerrevisionstorage.NewStorageV1Beta1(restOptionsGetter, platformClient)
	deploymentStore := deploymentstorage.NewStorageV1Beta1(restOptionsGetter, platformClient)
	statefulSetStore := statefulsetstorage.NewStorageV1Beta1(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"controllerrevisions":                  controllerRevisionStore.ControllerRevision,
		"deployments":                          deploymentStore.Deployment,
		"deployments/status":                   deploymentStore.Status,
		"deployments/scale":                    deploymentStore.Scale,
		"deployments/rollback":                 deploymentStore.Rollback,
		"deployments/pods":                     deploymentStore.Pods,
		"deployments/events":                   deploymentStore.Events,
		"deployments/horizontalpodautoscalers": deploymentStore.HPAs,
		"statefulsets":                         statefulSetStore.StatefulSet,
		"statefulsets/status":                  statefulSetStore.Status,
		"statefulsets/pods":                    statefulSetStore.Pods,
		"statefulsets/scale":                   statefulSetStore.Scale,
		"statefulsets/events":                  statefulSetStore.Events,
	}

	return storageMap
}

func (s *StorageProvider) v1Beta2Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	controllerRevisionStore := controllerrevisionstorage.NewStorageV1Beta2(restOptionsGetter, platformClient)
	daemonSetStore := daemonsetstorage.NewStorageV1Beta2(restOptionsGetter, platformClient)
	deploymentStore := deploymentstorage.NewStorageV1Beta2(restOptionsGetter, platformClient)
	statefulSetStore := statefulsetstorage.NewStorageV1Beta2(restOptionsGetter, platformClient)
	replicaSetStore := replicasetstorage.NewStorageV1Beta2(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"controllerrevisions":                  controllerRevisionStore.ControllerRevision,
		"daemonsets":                           daemonSetStore.DaemonSet,
		"daemonsets/pods":                      daemonSetStore.Pods,
		"daemonsets/status":                    daemonSetStore.Status,
		"daemonsets/events":                    daemonSetStore.Events,
		"deployments":                          deploymentStore.Deployment,
		"deployments/status":                   deploymentStore.Status,
		"deployments/scale":                    deploymentStore.Scale,
		"deployments/pods":                     deploymentStore.Pods,
		"deployments/events":                   deploymentStore.Events,
		"deployments/horizontalpodautoscalers": deploymentStore.HPAs,
		"statefulsets":                         statefulSetStore.StatefulSet,
		"statefulsets/status":                  statefulSetStore.Status,
		"statefulsets/pods":                    statefulSetStore.Pods,
		"statefulsets/scale":                   statefulSetStore.Scale,
		"statefulsets/events":                  statefulSetStore.Events,
		"replicasets":                          replicaSetStore.ReplicaSet,
		"replicasets/pods":                     replicaSetStore.Pods,
		"replicasets/status":                   replicaSetStore.Status,
		"replicasets/scale":                    replicaSetStore.Scale,
		"replicasets/events":                   replicaSetStore.Events,
	}

	return storageMap
}
