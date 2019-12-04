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
	"tkestack.io/tke/pkg/auth/controller/localidentity"

	"k8s.io/apimachinery/pkg/runtime/schema"
	v1 "tkestack.io/tke/api/auth/v1"
	"tkestack.io/tke/pkg/auth/controller/policy"
)

const (
	policySyncPeriod      = 5 * time.Minute
	concurrentPolicySyncs = 10

	localIdentitySyncPeriod      = 5 * time.Minute
	concurrentLocalIdentitySyncs = 10
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

func startLocalIdentityController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "policies"}] {
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
