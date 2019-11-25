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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
	"time"
	"tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/notify/controller/channel"
	"tkestack.io/tke/pkg/notify/controller/messagerequest"
)

const (
	channelSyncPeriod      = 30 * time.Second
	concurrentChannelSyncs = 10

	messageRequestSyncPeriod      = 5 * time.Minute
	concurrentMessageRequestSyncs = 10
)

func startChannelController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: "v1", Resource: "channels"}] {
		return nil, false, nil
	}

	ctrl := channel.NewController(
		ctx.ClientBuilder.ClientOrDie("channel-controller"),
		ctx.InformerFactory.Notify().V1().Channels(),
		channelSyncPeriod,
		v1.ChannelFinalize,
	)

	go ctrl.Run(concurrentChannelSyncs, ctx.Stop)

	return nil, true, nil
}

func startMessageRequestController(ctx ControllerContext) (http.Handler, bool, error) {
	if !ctx.AvailableResources[schema.GroupVersionResource{Group: v1.GroupName, Version: "v1", Resource: "messagerequests"}] {
		return nil, false, nil
	}

	ctrl := messagerequest.NewController(
		ctx.ClientBuilder.ClientOrDie("message-request-controller"),
		ctx.InformerFactory.Notify().V1().MessageRequests(),
		messageRequestSyncPeriod,
	)

	go ctrl.Run(concurrentMessageRequestSyncs, ctx.Stop)

	return nil, true, nil
}
