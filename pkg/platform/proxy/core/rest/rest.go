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
func (c LegacyRESTStorageProvider) NewLegacyRESTStorage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) (*genericserver.APIGroupInfo, error) {
	apiGroupInfo := &genericserver.APIGroupInfo{
		PrioritizedVersions:          platform.Scheme.PrioritizedVersionsForGroup(""),
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		Scheme:                       platform.Scheme,
		ParameterCodec:               platform.ParameterCodec,
		NegotiatedSerializer:         platform.Codecs,
	}

	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	podStore := podstorage.NewStorage(restOptionsGetter, platformClient)
	podTemplateStore := podtemplatestorage.NewStorage(restOptionsGetter, platformClient)
	endpointsStore := endpointsstorage.NewStorage(restOptionsGetter, platformClient)
	replicationControllerStore := replicationcontrollerstorage.NewStorage(restOptionsGetter, platformClient)
	namespaceStore := namespacestorage.NewStorage(restOptionsGetter, platformClient)
	nodeStore := nodestorage.NewStorage(restOptionsGetter, platformClient)
	serviceStore := servicestorage.NewStorage(restOptionsGetter, platformClient)
	resourceQuotaStore := resourcequotastorage.NewStorage(restOptionsGetter, platformClient)
	componentStatusStore := componentstatusstorage.NewStorage(restOptionsGetter, platformClient)
	limitRangeStore := limitrangestorage.NewStorage(restOptionsGetter, platformClient)
	configMapStore := configmapstorage.NewStorage(restOptionsGetter, platformClient)
	persistentVolumeStore := persistentvolumestorage.NewStorage(restOptionsGetter, platformClient)
	persistentVolumeClaimStore := persistentvolumeclaimstorage.NewStorage(restOptionsGetter, platformClient)
	serviceAccountStore := serviceaccountstorage.NewStorage(restOptionsGetter, platformClient)
	secretStore := secretstorage.NewStorage(restOptionsGetter, platformClient)
	eventStore := eventstorage.NewStorage(restOptionsGetter, platformClient)

	restStorageMap := map[string]rest.Storage{
		"pods":                          podStore.Pod,
		"pods/status":                   podStore.Status,
		"pods/binding":                  podStore.Binding,
		"pods/events":                   podStore.Events,
		"pods/log":                      podStore.Log,
		"bindings":                      podStore.Binding,
		"podTemplates":                  podTemplateStore.PodTemplate,
		"replicationControllers":        replicationControllerStore.ReplicationController,
		"replicationControllers/pods":   replicationControllerStore.Pods,
		"replicationControllers/status": replicationControllerStore.Status,
		"replicationControllers/scale":  replicationControllerStore.Scale,
		"replicationControllers/events": replicationControllerStore.Events,
		"services":                      serviceStore.Service,
		"services/status":               serviceStore.Status,
		"services/events":               serviceStore.Events,
		"endpoints":                     endpointsStore.Endpoint,
		"limitRanges":                   limitRangeStore.LimitRange,
		"resourceQuotas":                resourceQuotaStore.ResourceQuota,
		"resourceQuotas/status":         resourceQuotaStore.Status,
		"namespaces":                    namespaceStore.Namespace,
		"namespaces/status":             namespaceStore.Status,
		"namespaces/finalize":           namespaceStore.Finalize,
		"nodes":                         nodeStore.Node,
		"nodes/status":                  nodeStore.Status,
		"events":                        eventStore.Event,
		"secrets":                       secretStore.Secret,
		"serviceAccounts":               serviceAccountStore.ServiceAccount,
		"persistentVolumes":             persistentVolumeStore.PersistentVolume,
		"persistentVolumes/status":      persistentVolumeStore.Status,
		"persistentVolumes/events":      persistentVolumeStore.Events,
		"persistentVolumeClaims":        persistentVolumeClaimStore.PersistentVolumeClaim,
		"persistentVolumeClaims/status": persistentVolumeClaimStore.Status,
		"persistentVolumeClaims/events": persistentVolumeClaimStore.Events,
		"configMaps":                    configMapStore.ConfigMap,
		"componentStatuses":             componentStatusStore.ComponentStatus,
	}

	apiGroupInfo.VersionedResourcesStorageMap["v1"] = restStorageMap
	return apiGroupInfo, nil
}
