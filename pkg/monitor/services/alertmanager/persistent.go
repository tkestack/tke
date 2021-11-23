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

package alertmanager

import (
	"context"
	"strings"
	"time"

	"tkestack.io/tke/pkg/util/log"

	"tkestack.io/tke/pkg/monitor/util"
	"tkestack.io/tke/pkg/monitor/util/route"

	"github.com/pkg/errors"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	alertManagerConfigName = "alertmanager.yml"
	alertManagerConfigMap  = "alertmanager-config"
)

func (h *processor) loadConfig(ctx context.Context, clusterName string) (util.GenericRouteOperator, error) {
	k8sClient, err := util.GetClusterClient(ctx, clusterName, h.platformClient)
	if err != nil {
		return nil, err
	}

	configMap, err := k8sClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(ctx, alertManagerConfigMap, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	configData, ok := configMap.Data[alertManagerConfigName]
	if !ok {
		return nil, errors.Errorf("%s(%s) is not found", clusterName, alertManagerConfigName)
	}

	log.Infof("Load rule from configMap %s(%s)", clusterName, alertManagerConfigName)
	routeOp := route.NewRouteOperator()
	err = routeOp.Load(strings.NewReader(configData))
	if err != nil {
		return nil, err
	}

	return routeOp, nil
}

func (h *processor) saveConfig(ctx context.Context, clusterName string, data string) error {
	k8sClient, err := util.GetClusterClient(ctx, clusterName, h.platformClient)
	if err != nil {
		return err
	}

	log.Infof("Save rule to configMap %s(%s)", clusterName, alertManagerConfigMap)

	return wait.PollImmediate(time.Second, time.Second*5, func() (done bool, err error) {
		configMap, getErr := k8sClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(ctx, alertManagerConfigMap, metav1.GetOptions{})
		if getErr != nil {
			return false, getErr
		}

		configMap.Data[alertManagerConfigName] = data
		_, err = k8sClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Update(ctx, configMap, metav1.UpdateOptions{})
		if err == nil {
			return true, nil
		}

		if apierror.IsConflict(err) {
			return false, nil
		}

		return false, err
	})
}
