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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	componentstatusstorage "tkestack.io/tke/pkg/platform/proxy/core/componentstatus/storage"
	configmapstorage "tkestack.io/tke/pkg/platform/proxy/core/configmap/storage"
	endpointsstorage "tkestack.io/tke/pkg/platform/proxy/core/endpoint/storage"
	eventstorage "tkestack.io/tke/pkg/platform/proxy/core/event/storage"
	limitrangestorage "tkestack.io/tke/pkg/platform/proxy/core/limitrange/storage"
	namespacestorage "tkestack.io/tke/pkg/platform/proxy/core/namespace/storage"
	nodestorage "tkestack.io/tke/pkg/platform/proxy/core/node/storage"
	persistentvolumestorage "tkestack.io/tke/pkg/platform/proxy/core/persistentvolume/storage"
	persistentvolumeclaimstorage "tkestack.io/tke/pkg/platform/proxy/core/persistentvolumeclaim/storage"
	podstorage "tkestack.io/tke/pkg/platform/proxy/core/pod/storage"
	podtemplatestorage "tkestack.io/tke/pkg/platform/proxy/core/podtemplate/storage"
	replicationcontrollerstorage "tkestack.io/tke/pkg/platform/proxy/core/replicationcontroller/storage"
	resourcequotastorage "tkestack.io/tke/pkg/platform/proxy/core/resourcequota/storage"
	secretstorage "tkestack.io/tke/pkg/platform/proxy/core/secret/storage"
	servicestorage "tkestack.io/tke/pkg/platform/proxy/core/service/storage"
	serviceaccountstorage "tkestack.io/tke/pkg/platform/proxy/core/serviceaccount/storage"
)

// LegacyRESTStorageProvider provides information needed to build RESTStorage for core, but
// does NOT implement the "normal" RESTStorageProvider (yet!)
type LegacyRESTStorageProvider struct {
}

// LegacyRESTStorage returns stateful information about particular instances of
// REST storage for wiring controllers.
type LegacyRESTStorage struct {
}

// NewLegacyRESTStorage creates the APIGroupInfo by given rest options and
// returns it.
func (c LegacyRESTStorageProvider) NewLegacyRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) (*genericserver.APIGroupInfo, error) {
	apiGroupInfo := &genericserver.APIGroupInfo{
		PrioritizedVersions:          platform.Scheme.PrioritizedVersionsForGroup(""),
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		Scheme:                       platform.Scheme,
		ParameterCodec:               platform.ParameterCodec,
		NegotiatedSerializer:         platform.Codecs,
	}

	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	podStorage := podstorage.NewStorage(restOptionsGetter, platformClient)
	podTemplateStorage := podtemplatestorage.NewStorage(restOptionsGetter, platformClient)
	endpointsStorage := endpointsstorage.NewStorage(restOptionsGetter, platformClient)
	replicationControllerStorage := replicationcontrollerstorage.NewStorage(restOptionsGetter, platformClient)
	namespaceStorage := namespacestorage.NewStorage(restOptionsGetter, platformClient)
	nodeStorage := nodestorage.NewStorage(restOptionsGetter, platformClient)
	serviceStorage := servicestorage.NewStorage(restOptionsGetter, platformClient)
	resourceQuotaStorage := resourcequotastorage.NewStorage(restOptionsGetter, platformClient)
	componentStatusStorage := componentstatusstorage.NewStorage(restOptionsGetter, platformClient)
	limitRangeStorage := limitrangestorage.NewStorage(restOptionsGetter, platformClient)
	configMapStorage := configmapstorage.NewStorage(restOptionsGetter, platformClient)
	persistentVolumeStorage := persistentvolumestorage.NewStorage(restOptionsGetter, platformClient)
	persistentVolumeClaimStorage := persistentvolumeclaimstorage.NewStorage(restOptionsGetter, platformClient)
	serviceAccountStorage := serviceaccountstorage.NewStorage(restOptionsGetter, platformClient)
	secretStorage := secretstorage.NewStorage(restOptionsGetter, platformClient)
	eventStorage := eventstorage.NewStorage(restOptionsGetter, platformClient)

	storage := map[string]rest.Storage{}
	if resource := "pods"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = podStorage.Pod
		storage[resource+"/status"] = podStorage.Status
		storage[resource+"/log"] = podStorage.Log
		storage[resource+"/exec"] = podStorage.Exec
		storage[resource+"/binding"] = podStorage.Binding
		storage[resource+"/events"] = podStorage.Events

	}
	if resource := "bindings"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = podStorage.Binding
	}

	if resource := "podtemplates"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = podTemplateStorage.PodTemplate
	}

	if resource := "replicationcontrollers"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = replicationControllerStorage.ReplicationController
		storage[resource+"/pods"] = replicationControllerStorage.Pods
		storage[resource+"/status"] = replicationControllerStorage.Status
		storage[resource+"/scale"] = replicationControllerStorage.Scale
		storage[resource+"/events"] = replicationControllerStorage.Events
	}

	if resource := "services"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = serviceStorage.Service
		storage[resource+"/status"] = serviceStorage.Status
		storage[resource+"/events"] = serviceStorage.Events
	}

	if resource := "endpoints"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = endpointsStorage.Endpoint
	}

	if resource := "nodes"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = nodeStorage.Node
		storage[resource+"/pods"] = nodeStorage.Pods
		storage[resource+"/status"] = nodeStorage.Status
	}

	if resource := "events"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = eventStorage.Event
	}

	if resource := "limitranges"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = limitRangeStorage.LimitRange
	}

	if resource := "resourcequotas"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = resourceQuotaStorage.ResourceQuota
		storage[resource+"/status"] = resourceQuotaStorage.Status
	}

	if resource := "namespaces"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = namespaceStorage.Namespace
		storage[resource+"/status"] = namespaceStorage.Status
		storage[resource+"/finalize"] = namespaceStorage.Finalize
	}

	if resource := "secrets"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = secretStorage.Secret
	}

	if resource := "serviceaccounts"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = serviceAccountStorage.ServiceAccount
	}

	if resource := "persistentvolumes"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = persistentVolumeStorage.PersistentVolume
		storage[resource+"/status"] = persistentVolumeStorage.Status
		storage[resource+"/events"] = persistentVolumeStorage.Events
	}

	if resource := "persistentvolumeclaims"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = persistentVolumeClaimStorage.PersistentVolumeClaim
		storage[resource+"/status"] = persistentVolumeClaimStorage.Status
		storage[resource+"/events"] = persistentVolumeClaimStorage.Events
	}

	if resource := "configmaps"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = configMapStorage.ConfigMap
	}

	if resource := "componentstatuses"; apiResourceConfigSource.ResourceEnabled(corev1.SchemeGroupVersion.WithResource(resource)) {
		storage[resource] = componentStatusStorage.ComponentStatus
	}

	if len(storage) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap["v1"] = storage
	}

	return apiGroupInfo, nil
}
