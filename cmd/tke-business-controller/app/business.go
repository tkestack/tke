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

package app

import (
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	v1 "tkestack.io/tke/api/business/v1"
	"tkestack.io/tke/pkg/business/controller/namespace"
	"tkestack.io/tke/pkg/business/controller/project"
)

const (
	namespaceSyncPeriod      = 30 * time.Second
	concurrentNamespaceSyncs = 10

	projectSyncPeriod      = 5 * time.Minute
	concurrentProjectSyncs = 10
)

func startNamespaceController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: "v1", Resource: "namespaces"}] {
		return nil, false, nil
	}

	ctrl := namespace.NewController(
		ctx.PlatformClient,
		ctx.ClientBuilder.ClientOrDie("namespace-controller"),
		ctx.InformerFactory.Business().V1().Namespaces(),
		namespaceSyncPeriod,
		v1.NamespaceFinalize,
	)

	go ctrl.Run(concurrentNamespaceSyncs, ctx.Stop)

	return nil, true, nil
}

func startProjectController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: "v1", Resource: "projects"}] {
		return nil, false, nil
	}

	ctrl := project.NewController(
		ctx.ClientBuilder.ClientOrDie("project-controller"),
		ctx.InformerFactory.Business().V1().Projects(),
		projectSyncPeriod,
		v1.ProjectFinalize,
	)

	go ctrl.Run(concurrentProjectSyncs, ctx.Stop)

	return nil, true, nil
}
