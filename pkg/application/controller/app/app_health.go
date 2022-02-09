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
	"context"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	applicationv1 "tkestack.io/tke/api/application/v1"
	"tkestack.io/tke/pkg/util/log"
)

type applicationHealth struct {
	mu           sync.Mutex
	applications sets.String
}

func (s *applicationHealth) Exist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.applications.Has(key)
}

func (s *applicationHealth) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.applications.Delete(key)
}

func (s *applicationHealth) Set(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.applications.Insert(key)
}

func (c *Controller) startAppHealthCheck(ctx context.Context, key string) {
	if !c.health.Exist(key) {
		c.health.Set(key)
		go func() {
			if err := wait.PollImmediateUntil(1*time.Minute, c.watchAppHealth(ctx, key), c.stopCh); err != nil {
				log.Error("Failed to wait poll immediate until", log.Err(err))
			}
		}()
		log.Info("App phase start new health check", log.String("app", key))
	} else {
		log.Info("App phase health check exit", log.String("app", key))
	}
}

// for PollImmediateUntil, when return true ,an err while exit
func (c *Controller) watchAppHealth(ctx context.Context, key string) func() (bool, error) {
	return func() (bool, error) {
		log.Debug("Check app health", log.String("key", key))

		if !c.health.Exist(key) {
			return true, nil
		}

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			log.Error("Failed to split meta app key", log.String("key", key))
			c.health.Del(key)
			return true, nil
		}

		app, err := c.lister.Apps(namespace).Get(name)
		if err != nil && errors.IsNotFound(err) {
			log.Error("App not found, to exit the health check loop",
				log.String("name", name))
			c.health.Del(key)
			return true, nil
		}
		if err != nil {
			log.Error("Check app health, app get failed",
				log.String("name", name), log.Err(err))
			return false, nil
		}
		// if status is terminated,to exit the  health check loop
		if app.Status.Phase == applicationv1.AppPhaseTerminating {
			log.Warn("App status is terminated, to exit the health check loop",
				log.String("name", name))
			c.health.Del(key)
			return true, nil
		}
		return false, nil
	}
}
