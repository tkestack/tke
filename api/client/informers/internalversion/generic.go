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

package internalversion

import (
	"fmt"

	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
	auth "tkestack.io/tke/api/auth"
	business "tkestack.io/tke/api/business"
	logagent "tkestack.io/tke/api/logagent"
	monitor "tkestack.io/tke/api/monitor"
	notify "tkestack.io/tke/api/notify"
	platform "tkestack.io/tke/api/platform"
	registry "tkestack.io/tke/api/registry"
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
	// Group=auth.tkestack.io, Version=internalVersion
	case auth.SchemeGroupVersion.WithResource("apikeys"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().APIKeys().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("apisigningkeys"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().APISigningKeys().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("categories"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().Categories().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("clients"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().Clients().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().ConfigMaps().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("groups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().Groups().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("identityproviders"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().IdentityProviders().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("localgroups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().LocalGroups().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("localidentities"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().LocalIdentities().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("policies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().Policies().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("projects"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().Projects().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("projectpolicybindings"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().ProjectPolicyBindings().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("roles"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().Roles().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("rules"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().Rules().Informer()}, nil
	case auth.SchemeGroupVersion.WithResource("users"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Auth().InternalVersion().Users().Informer()}, nil

		// Group=business.tkestack.io, Version=internalVersion
	case business.SchemeGroupVersion.WithResource("chartgroups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().InternalVersion().ChartGroups().Informer()}, nil
	case business.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().InternalVersion().ConfigMaps().Informer()}, nil
	case business.SchemeGroupVersion.WithResource("imagenamespaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().InternalVersion().ImageNamespaces().Informer()}, nil
	case business.SchemeGroupVersion.WithResource("namespaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().InternalVersion().Namespaces().Informer()}, nil
	case business.SchemeGroupVersion.WithResource("nsemigrations"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().InternalVersion().NsEmigrations().Informer()}, nil
	case business.SchemeGroupVersion.WithResource("platforms"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().InternalVersion().Platforms().Informer()}, nil
	case business.SchemeGroupVersion.WithResource("projects"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Business().InternalVersion().Projects().Informer()}, nil

		// Group=logagent.tkestack.io, Version=internalVersion
	case logagent.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Logagent().InternalVersion().ConfigMaps().Informer()}, nil
	case logagent.SchemeGroupVersion.WithResource("logagents"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Logagent().InternalVersion().LogAgents().Informer()}, nil

		// Group=monitor.tkestack.io, Version=internalVersion
	case monitor.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Monitor().InternalVersion().ConfigMaps().Informer()}, nil
	case monitor.SchemeGroupVersion.WithResource("prometheuses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Monitor().InternalVersion().Prometheuses().Informer()}, nil

		// Group=notify.tkestack.io, Version=internalVersion
	case notify.SchemeGroupVersion.WithResource("channels"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().InternalVersion().Channels().Informer()}, nil
	case notify.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().InternalVersion().ConfigMaps().Informer()}, nil
	case notify.SchemeGroupVersion.WithResource("messages"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().InternalVersion().Messages().Informer()}, nil
	case notify.SchemeGroupVersion.WithResource("messagerequests"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().InternalVersion().MessageRequests().Informer()}, nil
	case notify.SchemeGroupVersion.WithResource("receivers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().InternalVersion().Receivers().Informer()}, nil
	case notify.SchemeGroupVersion.WithResource("receivergroups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().InternalVersion().ReceiverGroups().Informer()}, nil
	case notify.SchemeGroupVersion.WithResource("templates"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notify().InternalVersion().Templates().Informer()}, nil

		// Group=platform.tkestack.io, Version=internalVersion
	case platform.SchemeGroupVersion.WithResource("csioperators"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().CSIOperators().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("clusters"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().Clusters().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("clustercredentials"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().ClusterCredentials().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().ConfigMaps().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("cronhpas"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().CronHPAs().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("helms"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().Helms().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("ipams"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().IPAMs().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("lbcfs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().LBCFs().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("logcollectors"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().LogCollectors().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("machines"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().Machines().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("persistentevents"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().PersistentEvents().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("prometheuses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().Prometheuses().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("registries"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().Registries().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("tappcontrollers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().TappControllers().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("templates"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().Templates().Informer()}, nil
	case platform.SchemeGroupVersion.WithResource("volumedecorators"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Platform().InternalVersion().VolumeDecorators().Informer()}, nil

		// Group=registry.tkestack.io, Version=internalVersion
	case registry.SchemeGroupVersion.WithResource("charts"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().InternalVersion().Charts().Informer()}, nil
	case registry.SchemeGroupVersion.WithResource("chartgroups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().InternalVersion().ChartGroups().Informer()}, nil
	case registry.SchemeGroupVersion.WithResource("configmaps"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().InternalVersion().ConfigMaps().Informer()}, nil
	case registry.SchemeGroupVersion.WithResource("namespaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().InternalVersion().Namespaces().Informer()}, nil
	case registry.SchemeGroupVersion.WithResource("repositories"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Registry().InternalVersion().Repositories().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
