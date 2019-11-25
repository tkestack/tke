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

package route

import (
	alertconfig "github.com/prometheus/alertmanager/config"
	"reflect"
	"strings"
	"testing"
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
        url: 'http://127.0.0.1:5005/'
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
)

func TestRouteOperation(t *testing.T) {
	operator := NewRouteOperator()
	if err := operator.Load(strings.NewReader(exampleAlertConfig)); err != nil {
		t.Logf("can't create example route operator")
		return
	}

	insertedRoute := &alertconfig.Route{
		Match: map[string]string{
			"alert": "inserted",
		},
	}
	updatedRoute := &alertconfig.Route{
		Match: map[string]string{
			"alert": "updated",
		},
	}
	duplicateRoute := &alertconfig.Route{
		Match: map[string]string{
			"alert": "test",
		},
	}

	t.Logf("Test list routes")
	routes, err := operator.ListRoute()
	if err != nil {
		t.Errorf("can't list routes, %v", err)
		return
	}

	if len(routes) != 2 {
		t.Errorf("routes should be 2, got %d", len(routes))
		return
	}

	t.Logf("Insert with no label")
	_, err = operator.InsertRoute(&alertconfig.Route{})
	if err == nil {
		t.Errorf("shouldn't insert a route with no label")
		return
	}

	t.Logf("Insert duplicate route")
	_, err = operator.InsertRoute(duplicateRoute)
	if err == nil {
		t.Errorf("shouldn't insert a duplicate route")
		return
	}

	t.Logf("Insert a route")
	_, err = operator.InsertRoute(insertedRoute)
	if err != nil {
		t.Errorf("should insert a route, got %v", err)
		return
	}

	t.Logf("Update non-existed route")
	_, err = operator.UpdateRoute("non-existed", updatedRoute)
	if err == nil {
		t.Errorf("shouldn't update route")
		return
	}

	t.Logf("Update a route with no label")
	_, err = operator.UpdateRoute("test", &alertconfig.Route{})
	if err == nil {
		t.Errorf("shouldn't update route")
		return
	}

	t.Logf("Update a route")
	_, err = operator.UpdateRoute("test", updatedRoute)
	if err != nil {
		t.Errorf("should update route, got %v", err)
		return
	}

	t.Logf("Get non-existed route")
	_, err = operator.GetRoute("non-existed")
	if err == nil {
		t.Errorf("shouldn't get route")
		return
	}

	t.Logf("Get updated route")
	targetRoute, err := operator.GetRoute("updated")
	if err != nil {
		t.Errorf("should get route, got %v", err)
		return
	}

	if !reflect.DeepEqual(targetRoute, updatedRoute) {
		t.Errorf("updated route not equal, got %v, expect %v", targetRoute, updatedRoute)
		return
	}

	t.Logf("Delete non-existed route")
	_, err = operator.DeleteRoute("abc")
	if err == nil {
		t.Errorf("shouldn't delete abc")
		return
	}

	t.Logf("Delete updated route")
	_, err = operator.DeleteRoute("updated")
	if err != nil {
		t.Errorf("should delete updated, got %v", err)
		return
	}

	t.Logf("Get original route")
	_, err = operator.GetRoute("updated")
	if err == nil {
		t.Errorf("shouldn't get route, got %v", err)
		return
	}
}
