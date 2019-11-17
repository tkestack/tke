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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"time"
	"tkestack.io/tke/api/auth/v1"
	"tkestack.io/tke/pkg/auth/controller/apikey"
)

const (
	apikeySyncPeriod      = 5 * time.Minute
	concurrentApiKeySyncs = 10
)

func startAuthController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: v1.Version, Resource: "apikeys"}] {
		return nil, false, nil
	}

	ctrl := apikey.NewController(
		ctx.ClientBuilder.ClientOrDie("auth-controller"),
		ctx.InformerFactory.Auth().V1().APIKeys(),
		apikeySyncPeriod,
		)

	go ctrl.Run(concurrentApiKeySyncs, ctx.Stop)

	return nil, true, nil
}
