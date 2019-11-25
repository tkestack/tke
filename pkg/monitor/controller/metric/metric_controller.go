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

package metric

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
	"tkestack.io/tke/pkg/monitor/storage"
	"tkestack.io/tke/pkg/util/log"
)

// Controller is responsible for performing actions dependent upon a message request controller phase.
type Controller struct {
	storage storage.ProjectStorage
	stopCh  <-chan struct{}
}

func NewController(storage storage.ProjectStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers.
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting metric controller")
	defer log.Info("Shutting down metric controller")

	c.stopCh = stopCh

	go wait.Until(c.storage.Collect, 1*time.Minute, stopCh)

	<-stopCh
}
