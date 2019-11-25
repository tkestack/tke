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

package apiserver

import (
	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	certv1beta1 "k8s.io/api/certificates/v1beta1"
	coordinationv1 "k8s.io/api/coordination/v1"
	coordinationv1beta1 "k8s.io/api/coordination/v1beta1"
	corev1 "k8s.io/api/core/v1"
	eventsv1beta1 "k8s.io/api/events/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	nodev1alpha1 "k8s.io/api/node/v1alpha1"
	nodev1beta "k8s.io/api/node/v1beta1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"
	schedulingv1beta "k8s.io/api/scheduling/v1beta1"
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1alpha1 "k8s.io/api/storage/v1alpha1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"sync"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/apiserver/storage"
	admissionrest "tkestack.io/tke/pkg/platform/proxy/admissionregistration/rest"
	appsrest "tkestack.io/tke/pkg/platform/proxy/apps/rest"
	autoscalingrest "tkestack.io/tke/pkg/platform/proxy/autoscaling/rest"
	batchrest "tkestack.io/tke/pkg/platform/proxy/batch/rest"
	certrest "tkestack.io/tke/pkg/platform/proxy/certificates/rest"
	coordinationrest "tkestack.io/tke/pkg/platform/proxy/coordination/rest"
	corerest "tkestack.io/tke/pkg/platform/proxy/core/rest"
	eventsrest "tkestack.io/tke/pkg/platform/proxy/events/rest"
	extensionsrest "tkestack.io/tke/pkg/platform/proxy/extensions/rest"
	networkingrest "tkestack.io/tke/pkg/platform/proxy/networking/rest"
	noderest "tkestack.io/tke/pkg/platform/proxy/node/rest"
	policyrest "tkestack.io/tke/pkg/platform/proxy/policy/rest"
	rbacrest "tkestack.io/tke/pkg/platform/proxy/rbac/rest"
	schedulingrest "tkestack.io/tke/pkg/platform/proxy/scheduling/rest"
	settingsrest "tkestack.io/tke/pkg/platform/proxy/settings/rest"
	storagerest "tkestack.io/tke/pkg/platform/proxy/storage/rest"
	platformrest "tkestack.io/tke/pkg/platform/registry/rest"
	"tkestack.io/tke/pkg/util/log"
)

// ExtraConfig contains the additional configuration of apiserver.
type ExtraConfig struct {
	ServerName              string
	APIResourceConfigSource serverstorage.APIResourceConfigSource
	StorageFactory          serverstorage.StorageFactory
	VersionedInformers      versionedinformers.SharedInformerFactory
	ClusterProviders        *sync.Map
	MachineProviders        *sync.Map
	PrivilegedUsername      string
}

// Config contains the core configuration instance of apiserver and
// additional configuration.
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embed a private pointer of Config.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// APIServer contains state for a tke api server.
type APIServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

// Complete fills in any fields not set that are required to have valid data.
// It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	return CompletedConfig{&c}
}

// New returns a new instance of APIServer from the given config.
func (c completedConfig) New(delegationTarget genericapiserver.DelegationTarget) (*APIServer, error) {
	s, err := c.GenericConfig.New(c.ExtraConfig.ServerName, delegationTarget)
	if err != nil {
		return nil, err
	}

	m := &APIServer{
		GenericAPIServer: s,
	}

	// install legacy rest storage
	if c.ExtraConfig.APIResourceConfigSource.VersionEnabled(corev1.SchemeGroupVersion) {
		legacyRESTStorageProvider := corerest.LegacyRESTStorageProvider{}
		m.InstallLegacyAPI(&c, c.GenericConfig.RESTOptionsGetter, legacyRESTStorageProvider)
	}

	// The order here is preserved in discovery.
	restStorageProviders := []storage.RESTStorageProvider{
		&admissionrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&appsrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&autoscalingrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&batchrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&extensionsrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&networkingrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&coordinationrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&certrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&eventsrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&noderest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&policyrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&rbacrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&schedulingrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&settingsrest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},
		&storagerest.StorageProvider{LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig},

		&platformrest.StorageProvider{
			ClusterProviders:     c.ExtraConfig.ClusterProviders,
			MachineProviders:     c.ExtraConfig.MachineProviders,
			LoopbackClientConfig: c.GenericConfig.LoopbackClientConfig,
			PrivilegedUsername:   c.ExtraConfig.PrivilegedUsername,
		},
	}
	m.InstallAPIs(c.ExtraConfig.APIResourceConfigSource, c.GenericConfig.RESTOptionsGetter, restStorageProviders...)

	return m, nil
}

// InstallLegacyAPI will install the legacy API.
func (m *APIServer) InstallLegacyAPI(c *completedConfig, restOptionsGetter generic.RESTOptionsGetter, legacyRESTStorageProvider corerest.LegacyRESTStorageProvider) {
	apiGroupInfo, err := legacyRESTStorageProvider.NewLegacyRESTStorage(restOptionsGetter, c.GenericConfig.LoopbackClientConfig)
	if err != nil {
		log.Fatalf("Error building core storage: %v", err)
	}

	if err := m.GenericAPIServer.InstallLegacyAPIGroup(genericapiserver.DefaultLegacyAPIPrefix, apiGroupInfo); err != nil {
		log.Fatalf("Error in registering group versions: %v", err)
	}
}

// InstallAPIs will install the APIs for the restStorageProviders if they are enabled.
func (m *APIServer) InstallAPIs(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter, restStorageProviders ...storage.RESTStorageProvider) {
	var apiGroupsInfo []genericapiserver.APIGroupInfo

	for _, restStorageBuilder := range restStorageProviders {
		groupName := restStorageBuilder.GroupName()
		if !apiResourceConfigSource.AnyVersionForGroupEnabled(groupName) {
			log.Infof("Skipping disabled API group %q.", groupName)
			continue
		}
		apiGroupInfo, enabled := restStorageBuilder.NewRESTStorage(apiResourceConfigSource, restOptionsGetter)
		if !enabled {
			log.Warnf("Problem initializing API group %q, skipping.", groupName)
			continue
		}
		log.Infof("Enabling API group %q.", groupName)

		if postHookProvider, ok := restStorageBuilder.(genericapiserver.PostStartHookProvider); ok {
			name, hook, err := postHookProvider.PostStartHook()
			if err != nil {
				log.Fatalf("Error building PostStartHook: %v", err)
			}
			m.GenericAPIServer.AddPostStartHookOrDie(name, hook)
		}

		apiGroupsInfo = append(apiGroupsInfo, apiGroupInfo)
	}

	for i := range apiGroupsInfo {
		if err := m.GenericAPIServer.InstallAPIGroup(&apiGroupsInfo[i]); err != nil {
			log.Fatalf("Error in registering group versions: %v", err)
		}
	}
}

// DefaultAPIResourceConfigSource returns which groupVersion enabled and its
// resources enabled/disabled.
func DefaultAPIResourceConfigSource() *serverstorage.ResourceConfig {
	ret := serverstorage.NewResourceConfig()
	// NOTE: GroupVersions listed here will be enabled by default. Don't put alpha versions in the list.
	ret.EnableVersions(
		admissionv1beta1.SchemeGroupVersion,

		autoscalingv1.SchemeGroupVersion,
		autoscalingv2beta1.SchemeGroupVersion,

		appsv1.SchemeGroupVersion,
		appsv1beta2.SchemeGroupVersion,
		appsv1beta1.SchemeGroupVersion,

		batchv1.SchemeGroupVersion,
		batchv1beta1.SchemeGroupVersion,
		batchv2alpha1.SchemeGroupVersion,

		corev1.SchemeGroupVersion,
		certv1beta1.SchemeGroupVersion,

		extensionsv1beta1.SchemeGroupVersion,
		eventsv1beta1.SchemeGroupVersion,

		networkingv1.SchemeGroupVersion,
		networkingv1beta1.SchemeGroupVersion,

		coordinationv1.SchemeGroupVersion,
		coordinationv1beta1.SchemeGroupVersion,

		policyv1beta1.SchemeGroupVersion,
		rbacv1alpha1.SchemeGroupVersion,
		rbacv1beta1.SchemeGroupVersion,
		rbacv1.SchemeGroupVersion,

		schedulingv1alpha1.SchemeGroupVersion,
		schedulingv1beta.SchemeGroupVersion,
		schedulingv1.SchemeGroupVersion,

		settingsv1alpha1.SchemeGroupVersion,

		nodev1alpha1.SchemeGroupVersion,
		nodev1beta.SchemeGroupVersion,

		storagev1.SchemeGroupVersion,
		storagev1alpha1.SchemeGroupVersion,
		storagev1beta1.SchemeGroupVersion,

		platformv1.SchemeGroupVersion,
	)
	return ret
}
