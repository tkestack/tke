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

package chart

import (
	"context"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/util/log"
)

type chartHealth struct {
	mu     sync.Mutex
	charts sets.String
}

func (s *chartHealth) Exist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.charts.Has(key)
}

func (s *chartHealth) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.charts.Delete(key)
}

func (s *chartHealth) Set(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.charts.Insert(key)
}

func (c *Controller) startChartHealthCheck(ctx context.Context, key string) {
	if !c.health.Exist(key) {
		c.health.Set(key)
		go func() {
			if err := wait.PollImmediateUntil(1*time.Minute, c.watchChartHealth(ctx, key), c.stopCh); err != nil {
				log.Error("Failed to wait poll immediate until", log.Err(err))
			}
		}()
		log.Info("Chart phase start new health check", log.String("chart", key))
	} else {
		log.Info("Chart phase health check exit", log.String("chart", key))
	}
}

// for PollImmediateUntil, when return true ,an err while exit
func (c *Controller) watchChartHealth(ctx context.Context, key string) func() (bool, error) {
	return func() (bool, error) {
		log.Debug("Check chart health", log.String("key", key))

		if !c.health.Exist(key) {
			return true, nil
		}

		chartGroupName, chartName, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			log.Error("Failed to split meta chart key", log.String("key", key))
			c.health.Del(key)
			return true, nil
		}

		chart, err := c.client.RegistryV1().Charts(chartGroupName).Get(ctx, chartName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			log.Error("Chart not found, to exit the health check loop",
				log.String("chartName", chartName))
			c.health.Del(key)
			return true, nil
		}
		if err != nil {
			log.Error("Check chart health, chart get failed",
				log.String("chartName", chartName), log.Err(err))
			return false, nil
		}
		// if status is terminated,to exit the  health check loop
		if chart.Status.Phase == registryv1.ChartTerminating || chart.Status.Phase == registryv1.ChartPending {
			log.Warn("Chart status is terminated, to exit the health check loop",
				log.String("chartName", chartName))
			c.health.Del(key)
			return true, nil
		}
		return false, nil
	}
}
