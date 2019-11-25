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
	"time"

	"tkestack.io/tke/pkg/monitor/services"
	"tkestack.io/tke/pkg/monitor/util"
	alertmanager_config "tkestack.io/tke/pkg/platform/controller/addon/prometheus"
	"tkestack.io/tke/pkg/util/log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	exampleAlertConfig = `
global:
  resolve_timeout: 5m
route:
  group_by:
    - alertname
  group_wait: 1s
  group_interval: 1s
  repeat_interval: 1h
  receiver: web.hook
  routes:
    -
      match:
        alert: test
    -
      match_re:
        service: app
      receiver: web.hook
receivers:
  -
    name: web.hook
    webhook_configs:
      -
        url: 'http://10.12.91.240:5005/'
inhibit_rules:
  -
    source_match:
      severity: critical
    target_match:
      severity: warning
    equal:
      - alertname
      - dev
      - instance
`
	testClusterName = "fake"
)

func init() {
	logOpts := log.NewOptions()
	logOpts.EnableCaller = true
	logOpts.Level = log.InfoLevel
	log.Init(logOpts)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Flush()
		}
	}()
}

func createProcessorServer() (kubernetes.Interface, services.RouteProcessor, string, error) {
	k8sClient := fake.NewSimpleClientset()
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: alertmanager_config.AlertManagerConfigMap,
		},
		Data: map[string]string{
			alertmanager_config.AlertManagerConfigName: exampleAlertConfig,
		},
	}

	util.ClusterNameToClient.Store(testClusterName, k8sClient)
	_, _ = k8sClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Create(configMap)
	// Because we have set kubernetes client, so set nil is ok
	p := NewProcessor(nil)

	return k8sClient, p, testClusterName, nil
}
