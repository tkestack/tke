/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package hooks

import (
	"time"
	"tkestack.io/tke/pkg/auth/util"

	"github.com/casbin/casbin/v2"
	"k8s.io/client-go/tools/cache"
	"tkestack.io/tke/pkg/util/log"

	genericapiserver "k8s.io/apiserver/pkg/server"

	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	authv1informer "tkestack.io/tke/api/client/informers/externalversions/auth/v1"
)

type adapterHookHandler struct {
	authClient versionedclientset.Interface

	enforcer *casbin.SyncedEnforcer

	ruleInformer authv1informer.RuleInformer
}

// NewAdapterHookHandler creates a new adapterHookHandler object.
func NewAdapterHookHandler(authClient versionedclientset.Interface, enforcer *casbin.SyncedEnforcer, versionedInformers versionedinformers.SharedInformerFactory) genericapiserver.PostStartHookProvider {
	return &adapterHookHandler{
		authClient:   authClient,
		enforcer:     enforcer,
		ruleInformer: versionedInformers.Auth().V1().Rules(),
	}
}

func (d *adapterHookHandler) PostStartHook() (string, genericapiserver.PostStartHookFunc, error) {
	return "create-casbin-adapter", func(context genericapiserver.PostStartHookContext) error {

		go d.ruleInformer.Informer().Run(context.StopCh)
		if ok := cache.WaitForCacheSync(context.StopCh, d.ruleInformer.Informer().HasSynced); !ok {
			log.Error("Failed to wait for project caches to sync")
		}
		adpt := util.NewAdapter(d.authClient.AuthV1().Rules(), d.ruleInformer.Lister())
		d.enforcer.SetAdapter(adpt)
		d.enforcer.StartAutoLoadPolicy(2 * time.Second)
		return nil
	}, nil
}
