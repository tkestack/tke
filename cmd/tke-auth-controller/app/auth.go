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

package app

import (
	"net/http"
	"time"

	"tkestack.io/tke/pkg/auth/controller/config"
	"tkestack.io/tke/pkg/auth/controller/group"
	"tkestack.io/tke/pkg/auth/controller/localidentity"
	"tkestack.io/tke/pkg/auth/controller/projectpolicybinding"
	"tkestack.io/tke/pkg/auth/controller/role"

	"k8s.io/apimachinery/pkg/runtime/schema"
	v1 "tkestack.io/tke/api/auth/v1"
	"tkestack.io/tke/pkg/auth/controller/policy"
)

const (
	policySyncPeriod      = 5 * time.Minute
	concurrentPolicySyncs = 10

	projectPolicySyncPeriod      = 5 * time.Minute
	concurrentProjectPolicySyncs = 10

	localIdentitySyncPeriod      = 5 * time.Minute
	concurrentLocalIdentitySyncs = 5

	groupSyncPeriod      = 5 * time.Minute
	concurrentGroupSyncs = 5

	roleSyncPeriod      = 5 * time.Minute
	concurrentRoleSyncs = 5

	idpSyncPeriod      = 5 * time.Minute
	concurrentIDPSyncs = 5
)

func startPolicyController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "policies"}] {
		return nil, false, nil
	}

	ctrl := policy.NewController(
		ctx.ClientBuilder.ClientOrDie("policy-controller"),
		ctx.InformerFactory.Auth().V1().Policies(),
		ctx.InformerFactory.Auth().V1().Rules(),
		ctx.Enforcer,
		policySyncPeriod,
		v1.PolicyFinalize,
	)

	go ctrl.Run(concurrentPolicySyncs, ctx.Stop)

	return nil, true, nil
}

func startProjectPolicyController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "policies"}] {
		return nil, false, nil
	}

	ctrl := projectpolicybinding.NewController(
		ctx.ClientBuilder.ClientOrDie("projectpolicy-controller"),
		ctx.InformerFactory.Auth().V1().ProjectPolicyBindings(),
		ctx.InformerFactory.Auth().V1().Rules(),
		ctx.Enforcer,
		projectPolicySyncPeriod,
		v1.ProjectPolicyFinalize,
	)

	go ctrl.Run(concurrentProjectPolicySyncs, ctx.Stop)

	return nil, true, nil
}

func startLocalIdentityController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "localidentities"}] {
		return nil, false, nil
	}

	ctrl := localidentity.NewController(
		ctx.ClientBuilder.ClientOrDie("localidentity-controller"),
		ctx.InformerFactory.Auth().V1().LocalIdentities(),
		ctx.InformerFactory.Auth().V1().Rules(),
		ctx.Enforcer,
		localIdentitySyncPeriod,
		v1.LocalIdentityFinalize,
	)

	go ctrl.Run(concurrentLocalIdentitySyncs, ctx.Stop)

	return nil, true, nil
}

func startGroupController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "localgroups"}] {
		return nil, false, nil
	}

	ctrl := group.NewController(
		ctx.ClientBuilder.ClientOrDie("group-controller"),
		ctx.InformerFactory.Auth().V1().LocalGroups(),
		ctx.InformerFactory.Auth().V1().Rules(),
		ctx.Enforcer,
		groupSyncPeriod,
		v1.GroupFinalize,
	)

	go ctrl.Run(concurrentGroupSyncs, ctx.Stop)

	return nil, true, nil
}

func startRoleController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "roles"}] {
		return nil, false, nil
	}

	ctrl := role.NewController(
		ctx.ClientBuilder.ClientOrDie("role-controller"),
		ctx.InformerFactory.Auth().V1().Roles(),
		ctx.InformerFactory.Auth().V1().Rules(),
		ctx.Enforcer,
		roleSyncPeriod,
		v1.RoleFinalize,
	)

	go ctrl.Run(concurrentRoleSyncs, ctx.Stop)

	return nil, true, nil
}

func startConfigController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "categories"}] {
		return nil, false, nil
	}

	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "policies"}] {
		return nil, false, nil
	}

	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "identityproviders"}] {
		return nil, false, nil
	}

	ctrl := config.NewController(
		ctx.ClientBuilder.ClientOrDie("config-controller"),
		ctx.InformerFactory.Auth().V1().IdentityProviders(),
		idpSyncPeriod,
		ctx.PolicyPath,
		ctx.CategoryPath,
		ctx.TenantAdmin,
		ctx.TenantAdminSecret,
	)

	go ctrl.Run(concurrentIDPSyncs, ctx.Stop)

	return nil, true, nil
}
