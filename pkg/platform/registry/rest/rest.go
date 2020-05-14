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
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/storage"
	clusterstorage "tkestack.io/tke/pkg/platform/registry/cluster/storage"
	clusteraddontypestorage "tkestack.io/tke/pkg/platform/registry/clusteraddontype/storage"
	clustercredentialstorage "tkestack.io/tke/pkg/platform/registry/clustercredential/storage"
	configmapstorage "tkestack.io/tke/pkg/platform/registry/configmap/storage"
	cronhpastorage "tkestack.io/tke/pkg/platform/registry/cronhpa/storage"
	csioperatorstorage "tkestack.io/tke/pkg/platform/registry/csioperator/storage"
	helmstorage "tkestack.io/tke/pkg/platform/registry/helm/storage"
	ipamstorage "tkestack.io/tke/pkg/platform/registry/ipam/storage"
	lbcfstorage "tkestack.io/tke/pkg/platform/registry/lbcf/storage"
	logcollectorstorage "tkestack.io/tke/pkg/platform/registry/logcollector/storage"
	machinestorage "tkestack.io/tke/pkg/platform/registry/machine/storage"
	persistenteventstorage "tkestack.io/tke/pkg/platform/registry/persistentevent/storage"
	promstorage "tkestack.io/tke/pkg/platform/registry/prometheus/storage"
	registrystorage "tkestack.io/tke/pkg/platform/registry/registry/storage"
	tappcontrollertorage "tkestack.io/tke/pkg/platform/registry/tappcontroller/storage"
	volumedecoratorstorage "tkestack.io/tke/pkg/platform/registry/volumedecorator/storage"
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
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(platform.GroupName, platform.Scheme, platform.ParameterCodec, platform.Codecs)

	if apiResourceConfigSource.VersionEnabled(v1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[v1.SchemeGroupVersion.Version] = s.v1Storage(apiResourceConfigSource, restOptionsGetter, s.LoopbackClientConfig)
	}

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return platform.GroupName
}

func (s *StorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)
	storageMap := make(map[string]rest.Storage)

	{
		clusterREST := clusterstorage.NewStorage(restOptionsGetter, platformClient, loopbackClientConfig.Host, s.PrivilegedUsername)
		storageMap["clusters"] = clusterREST.Cluster
		storageMap["clusters/status"] = clusterREST.Status
		storageMap["clusters/finalize"] = clusterREST.Finalize
		storageMap["clusters/drain"] = clusterREST.Drain
		storageMap["clusters/proxy"] = clusterREST.Proxy
		storageMap["clusters/apply"] = clusterREST.Apply
		storageMap["clusters/helm"] = clusterREST.Helm
		storageMap["clusters/tapps"] = clusterREST.TappController
		storageMap["clusters/csis"] = clusterREST.CSI
		storageMap["clusters/pvcrs"] = clusterREST.PVCR
		storageMap["clusters/logcollector"] = clusterREST.LogCollector
		storageMap["clusters/cronhpas"] = clusterREST.CronHPA
		storageMap["clusters/addons"] = clusterREST.Addon
		storageMap["clusters/addontypes"] = clusterREST.AddonType
		storageMap["clusters/lbcflbdrivers"] = clusterREST.LBCFDriver
		storageMap["clusters/lbcflbs"] = clusterREST.LBCFLoadBalancer
		storageMap["clusters/lbcfbackendgroups"] = clusterREST.LBCFBackendGroup
		storageMap["clusters/lbcfbackendrecords"] = clusterREST.LBCFBackendRecord

		machineREST := machinestorage.NewStorage(restOptionsGetter, platformClient, s.PrivilegedUsername)
		storageMap["machines"] = machineREST.Machine
		storageMap["machines/status"] = machineREST.Status
		storageMap["machines/finalize"] = machineREST.Finalize

		clusterCredentialREST := clustercredentialstorage.NewStorage(restOptionsGetter, platformClient, s.PrivilegedUsername)
		storageMap["clustercredentials"] = clusterCredentialREST.ClusterCredential

		clusterAddonTypeREST := clusteraddontypestorage.NewStorage(restOptionsGetter)
		storageMap["clusteraddontypes"] = clusterAddonTypeREST.ClusterAddonType

		persistentEventREST := persistenteventstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["persistentevents"] = persistentEventREST.PersistentEvent
		storageMap["persistentevents/status"] = persistentEventREST.Status

		helmREST := helmstorage.NewStorage(restOptionsGetter, platformClient, s.PrivilegedUsername)
		storageMap["helms"] = helmREST.Helm
		storageMap["helms/status"] = helmREST.Status

		ipamREST := ipamstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["ipams"] = ipamREST.IPAM
		storageMap["ipams/status"] = ipamREST.Status

		configmapREST := configmapstorage.NewStorage(restOptionsGetter)
		storageMap["configmaps"] = configmapREST.ConfigMap

		registryREST := registrystorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["registries"] = registryREST.Registry

		tappControllerREST := tappcontrollertorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["tappcontrollers"] = tappControllerREST.TappController
		storageMap["tappcontrollers/status"] = tappControllerREST.Status

		csiOperatorREST := csioperatorstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["csioperators"] = csiOperatorREST.CSIOperator
		storageMap["csioperators/status"] = csiOperatorREST.Status

		volumeDecoratorREST := volumedecoratorstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["volumedecorators"] = volumeDecoratorREST.VolumeDecorator
		storageMap["volumedecorators/status"] = volumeDecoratorREST.Status

		logCollectorREST := logcollectorstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["logcollectors"] = logCollectorREST.LogCollector
		storageMap["logcollectors/status"] = logCollectorREST.Status

		cronHPAREST := cronhpastorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["cronhpas"] = cronHPAREST.CronHPA
		storageMap["cronhpas/status"] = cronHPAREST.Status

		promREST := promstorage.NewStorage(restOptionsGetter, s.PrivilegedUsername)
		storageMap["prometheuses"] = promREST.Prometheus
		storageMap["prometheuses/status"] = promREST.Status

		lbcfREST := lbcfstorage.NewStorage(restOptionsGetter, platformClient, s.PrivilegedUsername)
		storageMap["lbcfs"] = lbcfREST.LBCF
		storageMap["lbcfs/status"] = lbcfREST.Status
	}

	return storageMap
}
