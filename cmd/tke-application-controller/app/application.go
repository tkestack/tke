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
	applicationv1 "tkestack.io/tke/api/application/v1"
	"tkestack.io/tke/pkg/application/controller/app"
)

const (
	applicationSyncPeriod      = 30 * time.Second
	concurrentApplicationSyncs = 10
)

func startAppController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: applicationv1.GroupName, Version: "v1", Resource: "apps"}] {
		return nil, false, nil
	}

	ctrl := app.NewController(
		ctx.ClientBuilder.ClientOrDie("app-controller"),
		ctx.PlatformClient,
		ctx.Repo,
		ctx.InformerFactory.Application().V1().Apps(),
		applicationSyncPeriod,
		applicationv1.AppFinalize,
	)

	go ctrl.Run(concurrentApplicationSyncs, ctx.Stop)

	return nil, true, nil
}
