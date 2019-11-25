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
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	restclient "k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/storage"
	daemonsetstorage "tkestack.io/tke/pkg/platform/proxy/apps/daemonset/storage"
	deploymentstorage "tkestack.io/tke/pkg/platform/proxy/apps/deployment/storage"
	replicasetstorage "tkestack.io/tke/pkg/platform/proxy/apps/replicaset/storage"
	controllerstorage "tkestack.io/tke/pkg/platform/proxy/extensions/controller/storage"
	ingressstorage "tkestack.io/tke/pkg/platform/proxy/extensions/ingress/storage"
	networkpolicystorage "tkestack.io/tke/pkg/platform/proxy/networking/networkpolicy/storage"
	podsecuritypolicystorage "tkestack.io/tke/pkg/platform/proxy/policy/podsecuritypolicy/storage"
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
	apiGroupInfo := genericserver.NewDefaultAPIGroupInfo(v1beta1.GroupName, platform.Scheme, metav1.ParameterCodec, platform.Codecs)

	if apiResourceConfigSource.VersionEnabled(v1beta1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[v1beta1.SchemeGroupVersion.Version] = s.v1Beta1Storage(restOptionsGetter, s.LoopbackClientConfig)
	}

	return apiGroupInfo, true
}

// GroupName return the api group name
func (*StorageProvider) GroupName() string {
	return v1beta1.GroupName
}

func (s *StorageProvider) v1Beta1Storage(restOptionsGetter generic.RESTOptionsGetter, loopbackClientConfig *restclient.Config) map[string]rest.Storage {
	platformClient := platforminternalclient.NewForConfigOrDie(loopbackClientConfig)

	controllerStore := controllerstorage.NewStorageV1Beta1(restOptionsGetter, platformClient)
	daemonSetStore := daemonsetstorage.NewStorageExtensionsV1Beta1(restOptionsGetter, platformClient)
	deploymentStore := deploymentstorage.NewStorageExtensionsV1Beta1(restOptionsGetter, platformClient)
	ingressStore := ingressstorage.NewStorageV1Beta1(restOptionsGetter, platformClient)
	replicaSetStore := replicasetstorage.NewStorageExtensionsV1Beta1(restOptionsGetter, platformClient)
	networkPolicyStore := networkpolicystorage.NewStorageExtensionsV1Beta1(restOptionsGetter, platformClient)
	podSecurityPolicyStore := podsecuritypolicystorage.NewStorageExtensionsV1Beta1(restOptionsGetter, platformClient)

	storageMap := map[string]rest.Storage{
		"replicationcontrollers":               controllerStore.ReplicationController,
		"replicationcontrollers/pods":          controllerStore.Pods,
		"replicationcontrollers/scale":         controllerStore.Scale,
		"daemonsets":                           daemonSetStore.DaemonSet,
		"daemonsets/pods":                      daemonSetStore.Pods,
		"daemonsets/status":                    daemonSetStore.Status,
		"daemonsets/events":                    daemonSetStore.Events,
		"deployments":                          deploymentStore.Deployment,
		"deployments/status":                   deploymentStore.Status,
		"deployments/scale":                    deploymentStore.Scale,
		"deployments/rollback":                 deploymentStore.Rollback,
		"deployments/pods":                     deploymentStore.Pods,
		"deployments/events":                   deploymentStore.Events,
		"deployments/horizontalpodautoscalers": deploymentStore.HPAs,
		"ingresses":                            ingressStore.Ingress,
		"ingresses/status":                     ingressStore.Status,
		"ingresses/events":                     ingressStore.Events,
		"replicasets":                          replicaSetStore.ReplicaSet,
		"replicasets/pods":                     replicaSetStore.Pods,
		"replicasets/status":                   replicaSetStore.Status,
		"replicasets/scale":                    replicaSetStore.Scale,
		"replicasets/events":                   replicaSetStore.Events,
		"networkpolicies":                      networkPolicyStore.NetworkPolicy,
		"podSecurityPolicies":                  podSecurityPolicyStore.PodSecurityPolicy,
	}

	return storageMap
}
