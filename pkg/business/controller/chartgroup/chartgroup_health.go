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
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	businessv1 "tkestack.io/tke/api/business/v1"
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

		projectName, chartGroupName, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			log.Error("Failed to split meta chartGroup key", log.String("key", key))
			c.health.Del(key)
			return true, nil
		}

		chartGroup, err := c.client.BusinessV1().ChartGroups(projectName).Get(ctx, chartGroupName, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			log.Error("ChartGroup not found, to exit the health check loop",
				log.String("projectName", projectName), log.String("chartGroupName", chartGroupName))
			c.health.Del(key)
			return true, nil
		}
		if err != nil {
			log.Error("Check chartGroup health, chartGroup get failed",
				log.String("projectName", projectName), log.String("chartGroupName", chartGroupName), log.Err(err))
			return false, nil
		}
		// if status is terminated,to exit the  health check loop
		if chartGroup.Status.Phase == businessv1.ChartGroupTerminating || chartGroup.Status.Phase == businessv1.ChartGroupPending {
			log.Warn("ChartGroup status is terminated, to exit the health check loop",
				log.String("projectName", projectName), log.String("chartGroupName", chartGroupName))
			c.health.Del(key)
			return true, nil
		}

		if err := c.checkChartGroupHealth(ctx, chartGroup); err != nil {
			log.Error("Failed to check chartGroup health",
				log.String("projectName", projectName), log.String("chartGroupName", chartGroupName), log.Err(err))
		}
		return false, nil
	}
}

func (c *Controller) checkChartGroupHealth(ctx context.Context, chartGroup *businessv1.ChartGroup) error {
	chartGroupList, err := c.registryClient.ChartGroups().List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.tenantID=%s,spec.name=%s", chartGroup.Spec.TenantID, chartGroup.Spec.Name),
	})
	if err != nil {
		return err
	}

	switch chartGroup.Status.Phase {
	case businessv1.ChartGroupAvailable:
		if len(chartGroupList.Items) == 0 {
			chartGroup.Status.Phase = businessv1.ChartGroupFailed
			chartGroup.Status.Message = "ListChartGroup failed"
			chartGroup.Status.Reason = "ChartGroup may have been removed."
			chartGroup.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdate(ctx, chartGroup)
		}
		chartGroupObject := chartGroupList.Items[0]
		if chartGroupObject.Status.Locked != nil && *chartGroupObject.Status.Locked {
			chartGroup.Status.Phase = businessv1.ChartGroupLocked
			chartGroup.Status.Message = "ChartGroup locked"
			chartGroup.Status.Reason = "ChartGroup has been locked."
			chartGroup.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdate(ctx, chartGroup)
		}
	case businessv1.ChartGroupLocked:
		if len(chartGroupList.Items) == 0 {
			chartGroup.Status.Phase = businessv1.ChartGroupFailed
			chartGroup.Status.Message = "ListChartGroup failed"
			chartGroup.Status.Reason = "ChartGroup may have been removed."
			chartGroup.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdate(ctx, chartGroup)
		}
		chartGroupObject := chartGroupList.Items[0]
		if chartGroupObject.Status.Locked == nil || !*chartGroupObject.Status.Locked {
			chartGroup.Status.Phase = businessv1.ChartGroupAvailable
			chartGroup.Status.Message = ""
			chartGroup.Status.Reason = ""
			chartGroup.Status.LastTransitionTime = metav1.Now()
			return c.persistUpdate(ctx, chartGroup)
		}
	default:
		return fmt.Errorf("internal error, checkChartGroupHealth(%s/%s) found unexpected status %s",
			chartGroup.Namespace, chartGroup.Name, chartGroup.Status.Phase)
	}
	return nil
}
