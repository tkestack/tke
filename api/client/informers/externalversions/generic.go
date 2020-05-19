/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
	v1 "tkestack.io/tke/api/auth/v1"
	businessv1 "tkestack.io/tke/api/business/v1"
	logagentv1 "tkestack.io/tke/api/logagent/v1"
	monitorv1 "tkestack.io/tke/api/monitor/v1"
	notifyv1 "tkestack.io/tke/api/notify/v1"
	platformv1 "tkestack.io/tke/api/platform/v1"
	registryv1 "tkestack.io/tke/api/registry/v1"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=auth.tkestack.io, Version=v1
	case v1.SchemeGroupVersion.WithResource("apikeys"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().APIKeys().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("apisigningkeys"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().APISigningKeys().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("categories"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().Categories().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clients"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().Clients().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().ConfigMaps().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("groups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().Groups().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("identityproviders"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().IdentityProviders().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("localgroups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().LocalGroups().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("localidentities"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().LocalIdentities().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("policies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().Policies().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("projectpolicybindings"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().ProjectPolicyBindings().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("roles"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().Roles().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("rules"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().Rules().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("users"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().V1().Users().Informer()}, nil

		// Group=business.tkestack.io, Version=v1
	case businessv1.SchemeGroupVersion.WithResource("chartgroups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().V1().ChartGroups().Informer()}, nil
	case businessv1.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().V1().ConfigMaps().Informer()}, nil
	case businessv1.SchemeGroupVersion.WithResource("imagenamespaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().V1().ImageNamespaces().Informer()}, nil
	case businessv1.SchemeGroupVersion.WithResource("namespaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().V1().Namespaces().Informer()}, nil
	case businessv1.SchemeGroupVersion.WithResource("nsemigrations"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().V1().NsEmigrations().Informer()}, nil
	case businessv1.SchemeGroupVersion.WithResource("platforms"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().V1().Platforms().Informer()}, nil
	case businessv1.SchemeGroupVersion.WithResource("projects"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().V1().Projects().Informer()}, nil

		// Group=logagent.tkestack.io, Version=v1
	case logagentv1.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Logagent().V1().ConfigMaps().Informer()}, nil
	case logagentv1.SchemeGroupVersion.WithResource("logagents"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Logagent().V1().LogAgents().Informer()}, nil

		// Group=monitor.tkestack.io, Version=v1
	case monitorv1.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Monitor().V1().ConfigMaps().Informer()}, nil

		// Group=notify.tkestack.io, Version=v1
	case notifyv1.SchemeGroupVersion.WithResource("channels"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().V1().Channels().Informer()}, nil
	case notifyv1.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().V1().ConfigMaps().Informer()}, nil
	case notifyv1.SchemeGroupVersion.WithResource("messages"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().V1().Messages().Informer()}, nil
	case notifyv1.SchemeGroupVersion.WithResource("messagerequests"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().V1().MessageRequests().Informer()}, nil
	case notifyv1.SchemeGroupVersion.WithResource("receivers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().V1().Receivers().Informer()}, nil
	case notifyv1.SchemeGroupVersion.WithResource("receivergroups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().V1().ReceiverGroups().Informer()}, nil
	case notifyv1.SchemeGroupVersion.WithResource("templates"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().V1().Templates().Informer()}, nil

		// Group=platform.tkestack.io, Version=v1
	case platformv1.SchemeGroupVersion.WithResource("csioperators"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().CSIOperators().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("clusters"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().Clusters().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("clustercredentials"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().ClusterCredentials().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().ConfigMaps().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("cronhpas"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().CronHPAs().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("helms"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().Helms().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("ipams"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().IPAMs().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("lbcfs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().LBCFs().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("logcollectors"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().LogCollectors().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("machines"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().Machines().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("persistentevents"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().PersistentEvents().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("prometheuses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().Prometheuses().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("registries"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().Registries().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("tappcontrollers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().TappControllers().Informer()}, nil
	case platformv1.SchemeGroupVersion.WithResource("volumedecorators"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().V1().VolumeDecorators().Informer()}, nil

		// Group=registry.tkestack.io, Version=v1
	case registryv1.SchemeGroupVersion.WithResource("charts"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().V1().Charts().Informer()}, nil
	case registryv1.SchemeGroupVersion.WithResource("chartgroups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().V1().ChartGroups().Informer()}, nil
	case registryv1.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().V1().ConfigMaps().Informer()}, nil
	case registryv1.SchemeGroupVersion.WithResource("namespaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().V1().Namespaces().Informer()}, nil
	case registryv1.SchemeGroupVersion.WithResource("repositories"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().V1().Repositories().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
