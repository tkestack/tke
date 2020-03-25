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

	"tkestack.io/tke/pkg/business/controller/platform"

	"k8s.io/apimachinery/pkg/runtime/schema"
	businessv1 "tkestack.io/tke/api/business/v1"
	"tkestack.io/tke/pkg/business/controller/chartgroup"
	"tkestack.io/tke/pkg/business/controller/imagenamespace"
	"tkestack.io/tke/pkg/business/controller/namespace"
	"tkestack.io/tke/pkg/business/controller/project"
)

const (
	namespaceSyncPeriod      = 30 * time.Second
	concurrentNamespaceSyncs = 10

	projectSyncPeriod      = 5 * time.Minute
	concurrentProjectSyncs = 10

	imageNamespaceSyncPeriod      = 30 * time.Second
	concurrentImageNamespaceSyncs = 10

	chartGroupSyncPeriod      = 30 * time.Second
	concurrentChartGroupSyncs = 10

	platformSyncPeriod      = 30 * time.Second
	concurrentPlatformSyncs = 1
)

func startNamespaceController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: businessv1.GroupName, Version: "v1", Resource: "namespaces"}] {
		return nil, false, nil
	}

	ctrl := namespace.NewController(
		ctx.PlatformClient,
		ctx.ClientBuilder.ClientOrDie("namespace-controller"),
		ctx.InformerFactory.Business().V1().Namespaces(),
		namespaceSyncPeriod,
		businessv1.NamespaceFinalize,
	)

	go ctrl.Run(concurrentNamespaceSyncs, ctx.Stop)

	return nil, true, nil
}

func startProjectController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: businessv1.GroupName, Version: "v1", Resource: "projects"}] {
		return nil, false, nil
	}

	ctrl := project.NewController(
		ctx.ClientBuilder.ClientOrDie("project-controller"),
		ctx.InformerFactory.Business().V1().Projects(),
		projectSyncPeriod,
		businessv1.ProjectFinalize,
		ctx.RegistryClient != nil,
	)

	go ctrl.Run(concurrentProjectSyncs, ctx.Stop)

	return nil, true, nil
}

func startImageNamespaceController(ctx ControllerContext) (http.Handler, bool, error) {
	if ctx.RegistryClient == nil {
		return nil, false, nil
	}

	if !ctx.AvailableResources[schema.GroupVersionResource{Group: businessv1.GroupName, Version: "v1", Resource: "imagenamespaces"}] {
		return nil, false, nil
	}

	ctrl := imagenamespace.NewController(
		ctx.RegistryClient,
		ctx.ClientBuilder.ClientOrDie("imagenamespace-controller"),
		ctx.InformerFactory.Business().V1().ImageNamespaces(),
		imageNamespaceSyncPeriod,
		businessv1.ImageNamespaceFinalize,
	)

	go ctrl.Run(concurrentImageNamespaceSyncs, ctx.Stop)

	return nil, true, nil
}

func startChartGroupController(ctx ControllerContext) (http.Handler, bool, error) {
	if ctx.RegistryClient == nil {
		return nil, false, nil
	}

	if !ctx.AvailableResources[schema.GroupVersionResource{Group: businessv1.GroupName, Version: "v1", Resource: "chartgroups"}] {
		return nil, false, nil
	}

	ctrl := chartgroup.NewController(
		ctx.RegistryClient,
		ctx.ClientBuilder.ClientOrDie("chartgroup-controller"),
		ctx.InformerFactory.Business().V1().ChartGroups(),
		chartGroupSyncPeriod,
		businessv1.ChartGroupFinalize,
	)

	go ctrl.Run(concurrentChartGroupSyncs, ctx.Stop)

	return nil, true, nil
}

func startPlatformController(ctx ControllerContext) (http.Handler, bool, error) {
	if ctx.AuthClient == nil {
		return nil, false, nil
	}

	if !ctx.AvailableResources[schema.GroupVersionResource{Group: businessv1.GroupName, Version: "v1", Resource: "platforms"}] {
		return nil, false, nil
	}

	ctrl := platform.NewController(
		ctx.ClientBuilder.ClientOrDie("platform-controller"),
		ctx.AuthClient,
		ctx.InformerFactory.Business().V1().Platforms(),
		platformSyncPeriod,
	)

	go ctrl.Run(concurrentPlatformSyncs, ctx.Stop)

	return nil, true, nil
}
