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
	versionedclientset "tkestack.io/tke/api/client/clientset/versioned"
	"tkestack.io/tke/api/client/informers/externalversions"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/pkg/platform/controller/addon/cronhpa"
	"tkestack.io/tke/pkg/platform/controller/addon/logcollector"
	"tkestack.io/tke/pkg/platform/controller/addon/persistentevent"
	"tkestack.io/tke/pkg/platform/controller/addon/prometheus"
	"tkestack.io/tke/pkg/platform/controller/addon/storage/csioperator"
	"tkestack.io/tke/pkg/platform/controller/addon/tappcontroller"
	bootstrapps "tkestack.io/tke/pkg/platform/controller/bootstrapapps"
	clustercontroller "tkestack.io/tke/pkg/platform/controller/cluster"
	"tkestack.io/tke/pkg/platform/controller/machine"
)

const (
	persistentEventSyncPeriod      = 5 * time.Minute
	concurrentPersistentEventSyncs = 5

	eventSyncPeriod = 5 * time.Minute
	concurrentSyncs = 10

	promEventSyncPeriod = 5 * time.Minute
	concurrentPromSyncs = 10
)

func startClusterController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "clusters"}] {
		return nil, false, nil
	}

	ctrl := clustercontroller.NewController(
		ctx.ClientBuilder.ClientOrDie("cluster-controller").PlatformV1(),
		ctx.InformerFactory.Platform().V1().Clusters(),
		ctx.Config.ClusterController,
		platformv1.ClusterFinalize,
	)

	go func() {
		_ = ctrl.Run(ctx.Config.ClusterController.ConcurrentClusterSyncs, ctx.Stop)
	}()

	return nil, true, nil
}

func startMachineController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "machines"}] {
		return nil, false, nil
	}

	ctrl := machine.NewController(
		ctx.ClientBuilder.ClientOrDie("machine-controller").PlatformV1(),
		ctx.InformerFactory.Platform().V1().Machines(),
		ctx.Config.MachineController,
		platformv1.MachineFinalize,
	)

	go func() {
		_ = ctrl.Run(ctx.Config.MachineController.ConcurrentMachineSyncs, ctx.Stop)
	}()

	return nil, true, nil
}

func startPersistentEventController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "persistentevents"}] {
		return nil, false, nil
	}

	ctrl := persistentevent.NewController(
		ctx.ClientBuilder.ClientOrDie("persistentevent-controller"),
		ctx.InformerFactory.Platform().V1().PersistentEvents(),
		persistentEventSyncPeriod,
	)

	go func() {
		_ = ctrl.Run(concurrentPersistentEventSyncs, ctx.Stop)
	}()

	return nil, true, nil
}

func startTappControllerController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "tappcontrollers"}] {
		return nil, false, nil
	}

	ctrl := tappcontroller.NewController(
		ctx.ClientBuilder.ClientOrDie("tapp-controller-controller"),
		ctx.InformerFactory.Platform().V1().TappControllers(),
		eventSyncPeriod,
	)

	go func() {
		_ = ctrl.Run(concurrentSyncs, ctx.Stop)
	}()

	return nil, true, nil
}

func startCronHPAController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "cronhpas"}] {
		return nil, false, nil
	}

	ctrl := cronhpa.NewController(
		ctx.ClientBuilder.ClientOrDie("cron-hpa-controller"),
		ctx.InformerFactory.Platform().V1().CronHPAs(),
		eventSyncPeriod,
	)

	go func() {
		_ = ctrl.Run(concurrentSyncs, ctx.Stop)
	}()

	return nil, true, nil
}

func startCSIOperatorController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "csioperators"}] {
		return nil, false, nil
	}

	ctrl := csioperator.NewController(
		ctx.ClientBuilder.ClientOrDie("csi-operator-controller"),
		ctx.InformerFactory.Platform().V1().CSIOperators(),
		eventSyncPeriod,
	)

	go func() {
		_ = ctrl.Run(concurrentSyncs, ctx.Stop)
	}()

	return nil, true, nil
}

func startLogCollectorController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "logcollectors"}] {
		return nil, false, nil
	}

	ctrl := logcollector.NewController(
		ctx.ClientBuilder.ClientOrDie("log-collector-controller"),
		ctx.InformerFactory.Platform().V1().LogCollectors(),
		eventSyncPeriod,
	)

	go func() {
		_ = ctrl.Run(concurrentSyncs, ctx.Stop)
	}()

	return nil, true, nil
}

func startPrometheusController(ctx ControllerContext) (http.Handler, bool, error) {
	if ctx.RemoteType == "" || len(ctx.RemoteAddresses) == 0 {
		return nil, false, nil
	}

	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "prometheuses"}] {
		return nil, false, nil
	}

	ctrl := prometheus.NewController(
		ctx.ClientBuilder.ClientOrDie("prometheus-controller"),
		ctx.InformerFactory.Platform().V1().Prometheuses(),
		promEventSyncPeriod,

		ctx.RemoteAddresses,
		ctx.RemoteType,
	)

	go func() {
		_ = ctrl.Run(concurrentPromSyncs, ctx.Stop)
	}()

	return nil, true, nil
}

func startBootstrapAppsController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: platformv1.GroupName, Version: "v1", Resource: "clusters"}] ||
		ctx.ApplicationClient == nil {
		return nil, false, nil
	}
	appclientset := versionedclientset.NewForConfigOrDie(ctx.Config.ApplicationAPIServerClientConfig)
	appInformerFactory := externalversions.NewSharedInformerFactory(appclientset, ctx.ResyncPeriod())

	ctrl := bootstrapps.NewBootstrapAppsController(
		ctx.ClientBuilder.ClientOrDie("bootstrap-apps-controller").PlatformV1(),
		appclientset.ApplicationV1(),
		ctx.InformerFactory.Platform().V1().Clusters(),
		appInformerFactory.Application().V1().Apps(),
		ctx.Config.ClusterController,
		platformv1.ClusterFinalize,
	)

	appInformerFactory.Start(ctx.Stop)

	go func() {
		_ = ctrl.Run(ctx.Config.ClusterController.ConcurrentClusterSyncs, ctx.Stop)
	}()

	return nil, true, nil
}
