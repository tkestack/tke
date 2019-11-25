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
	"reflect"
	"testing"

	alertmanager_config "tkestack.io/tke/pkg/platform/controller/addon/prometheus"

	alert_config "github.com/prometheus/alertmanager/config"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestProcessor_Create(t *testing.T) {
	k8sClient, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	t.Logf("With no alert label")
	route := &alert_config.Route{}
	err = p.Create(clusterName, "test", route)
	if err == nil {
		t.Errorf("creation should failed")
		return
	}

	route.Match = map[string]string{
		"alert": "test",
	}
	t.Logf("With duplicate match label")
	err = p.Create(clusterName, "test", route)
	if err == nil {
		t.Errorf("creation should failed")
		return
	}

	route.Match = map[string]string{
		"alert": "abc",
		"foo":   "bar",
	}

	t.Logf("With correct label")
	err = p.Create(clusterName, "abc", route)
	if err != nil {
		t.Errorf("creation should success, code: %s", err)
		return
	}

	t.Logf("Validate persistent data")
	configMap, err := k8sClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(alertmanager_config.AlertManagerConfigMap, metav1.GetOptions{})
	if err != nil {
		t.Errorf("can't get persistent data, %v", err)
		return
	}

	targetConfig := &alert_config.Config{}
	expectConfig := &alert_config.Config{}
	_ = yaml.Unmarshal([]byte(configMap.Data[alertmanager_config.AlertManagerConfigName]), targetConfig)
	_ = yaml.Unmarshal([]byte(exampleAlertConfig), expectConfig)

	expectConfig.Route.Routes = append(expectConfig.Route.Routes, route)

	if !reflect.DeepEqual(targetConfig.Route, expectConfig.Route) {
		t.Errorf("persistent data is not equal, got %+v, expect %+v", targetConfig.Route, expectConfig.Route)
	}
}
