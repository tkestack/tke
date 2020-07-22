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

package local

import (
	"time"

	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	versionedinformers "tkestack.io/tke/api/client/informers/externalversions"
	authv1informer "tkestack.io/tke/api/client/informers/externalversions/auth/v1"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"

	"github.com/casbin/casbin/v2"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/tools/cache"
)

type adapterHookHandler struct {
	authClient versionedclientset.Interface

	enforcer       *casbin.SyncedEnforcer
	ruleInformer   authv1informer.RuleInformer
	reloadInterval time.Duration
}

// NewAdapterHookHandler creates a new adapterHookHandler object.
func NewAdapterHookHandler(authClient versionedclientset.Interface, enforcer *casbin.SyncedEnforcer, versionedInformers versionedinformers.SharedInformerFactory, reloadInterval time.Duration) genericapiserver.PostStartHookProvider {
	return &adapterHookHandler{
		authClient:     authClient,
		enforcer:       enforcer,
		reloadInterval: reloadInterval,
		ruleInformer:   versionedInformers.Auth().V1().Rules(),
	}
}

func (d *adapterHookHandler) PostStartHook() (string, genericapiserver.PostStartHookFunc, error) {
	return "create-casbin-adapter", func(context genericapiserver.PostStartHookContext) error {
		log.Info("start create casbin server")
		go d.ruleInformer.Informer().Run(context.StopCh)
		if ok := cache.WaitForCacheSync(context.StopCh, d.ruleInformer.Informer().HasSynced); !ok {
			log.Error("Failed to wait for project caches to sync")
		}

		adpt := util.NewAdapter(d.authClient.AuthV1().Rules(), d.ruleInformer.Lister())
		d.enforcer.SetAdapter(adpt)

		rm := util.NewRoleManager(10)
		d.enforcer.SetRoleManager(rm)
		_ = d.enforcer.LoadPolicy()

		d.enforcer.StartAutoLoadPolicy(d.reloadInterval)
		log.Info("finish start create casbin server")
		return nil
	}, nil
}
