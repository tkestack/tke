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
	"reflect"
	"testing"

	alertconfig "github.com/prometheus/alertmanager/config"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	alertmanagerconfig "tkestack.io/tke/pkg/platform/controller/addon/prometheus"
)

func TestProcessor_Update(t *testing.T) {
	k8sClient, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	route := &alertconfig.Route{}

	t.Logf("With no alert label")
	err = p.Update(context.Background(), clusterName, "test", route)
	if err == nil {
		t.Errorf("update should failed")
		return
	}

	route.Match = map[string]string{
		"alert": "abc",
		"foo":   "bar",
	}
	t.Logf("With non-existed label")
	err = p.Update(context.Background(), clusterName, "abc", route)
	if err == nil {
		t.Errorf("update should failed")
		return
	}

	t.Logf("With correct label")
	err = p.Update(context.Background(), clusterName, "test", route)
	if err != nil {
		t.Errorf("creation should success, code: %s", err)
		return
	}

	t.Logf("Validate persistent data")
	configMap, err := k8sClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(context.Background(), alertmanagerconfig.AlertManagerConfigMap, metav1.GetOptions{})
	if err != nil {
		t.Errorf("can't get persistent data, %v", err)
		return
	}

	targetConfig := &alertconfig.Config{}
	expectConfig := &alertconfig.Config{}
	_ = yaml.Unmarshal([]byte(configMap.Data[alertmanagerconfig.AlertManagerConfigName]), targetConfig)
	_ = yaml.Unmarshal([]byte(exampleAlertConfig), expectConfig)

	expectConfig.Route.Routes[0] = route

	if !reflect.DeepEqual(targetConfig.Route, expectConfig.Route) {
		t.Errorf("persistent data is not equal, got %+v, expect %+v", targetConfig.Route, expectConfig.Route)
	}
}
