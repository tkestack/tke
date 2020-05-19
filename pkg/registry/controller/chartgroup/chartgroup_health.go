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

package chartgroup

import (
	"context"
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/pkg/util/log"
)

type chartGroupHealth struct {
	mu          sync.Mutex
	chartGroups sets.String
}

func (s *chartGroupHealth) Exist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.chartGroups.Has(key)
}

func (s *chartGroupHealth) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chartGroups.Delete(key)
}

func (s *chartGroupHealth) Set(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chartGroups.Insert(key)
}

func (c *Controller) startChartGroupHealthCheck(ctx context.Context, key string) {
	if !c.health.Exist(key) {
		c.health.Set(key)
		go func() {
			if err := wait.PollImmediateUntil(1*time.Minute, c.watchChartGroupHealth(ctx, key), c.stopCh); err != nil {
				log.Error("Failed to wait poll immediate until", log.Err(err))
			}
		}()
		log.Info("ChartGroup phase start new health check", log.String("chartGroup", key))
	} else {
		log.Info("ChartGroup phase health check exit", log.String("chartGroup", key))
	}
}

// for PollImmediateUntil, when return true ,an err while exit
func (c *Controller) watchChartGroupHealth(ctx context.Context, key string) func() (bool, error) {
	return func() (bool, error) {
		log.Debug("Check chartGroup health", log.String("key", key))

		if !c.health.Exist(key) {
			return true, nil
		}

		_, chartGroupName, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			log.Error("Failed to split meta chartGroup key", log.String("key", key))
			c.health.Del(key)
			return true, nil
		}

		chartGroup, err := c.client.RegistryV1().ChartGroups().Get(ctx, chartGroupName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			log.Error("ChartGroup not found, to exit the health check loop",
				log.String("chartGroupName", chartGroupName))
			c.health.Del(key)
			return true, nil
		}
		if err != nil {
			log.Error("Check chartGroup health, chartGroup get failed",
				log.String("chartGroupName", chartGroupName), log.Err(err))
			return false, nil
		}
		// if status is terminated,to exit the  health check loop
		if chartGroup.Status.Phase == registryv1.ChartGroupTerminating || chartGroup.Status.Phase == registryv1.ChartGroupPending {
			log.Warn("ChartGroup status is terminated, to exit the health check loop",
				log.String("chartGroupName", chartGroupName))
			c.health.Del(key)
			return true, nil
		}

		if err := c.checkChartGroupHealth(ctx, chartGroup); err != nil {
			log.Error("Failed to check chartGroup health",
				log.String("chartGroupName", chartGroupName), log.Err(err))
		}
		return false, nil
	}
}

func (c *Controller) checkChartGroupHealth(ctx context.Context, chartGroup *registryv1.ChartGroup) error {
	var errs []error
	for _, p := range chartGroup.Spec.Projects {
		_, err := c.businessClient.ChartGroups(p).Get(ctx, chartGroup.Spec.Name, metav1.GetOptions{})
		switch chartGroup.Status.Phase {
		case registryv1.ChartGroupAvailable:
			if err != nil {
				// DONOT DO THIS! We might delete the project first, cascading delete the chartgroups.business.
				// If chartgroups.business have been removed, then it maight update chartGroup's status phase to failed
				//
				// chartGroup.Status.Phase = registryv1.ChartGroupFailed
				// chartGroup.Status.Message = "GetBusinessChartGroup failed"
				// chartGroup.Status.Reason = fmt.Sprintf("BusinessChartGroup may have been removed, %s", err.Error())
				// chartGroup.Status.LastTransitionTime = metav1.Now()
				// err = c.persistUpdate(chartGroup)
				// if err != nil {
				// 	errs = append(errs, err)
				// }
				errs = append(errs, fmt.Errorf("BusinessChartGroup may have been removed, %s", err.Error()))
			}
		default:
			errs = append(errs, fmt.Errorf("internal error, checkChartGroupHealth(%s) found unexpected status %s",
				chartGroup.Name, chartGroup.Status.Phase))
		}
	}
	return utilerrors.NewAggregate(errs)
}
