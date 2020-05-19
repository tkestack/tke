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
)

func TestProcessor_List(t *testing.T) {
	_, p, clusterName, err := createProcessorServer()
	if err != nil {
		t.Errorf("can't create processor server, %v", err)
		return
	}

	t.Logf("List all routes")
	targetRoutes, err := p.List(context.Background(), clusterName)
	if err != nil {
		t.Errorf("list should success, code: %s", err)
		return
	}

	expectConfig := &alertconfig.Config{}
	_ = yaml.Unmarshal([]byte(exampleAlertConfig), expectConfig)

	if !reflect.DeepEqual(targetRoutes, expectConfig.Route.Routes) {
		t.Errorf("persistent data is not equal, got %+v, expect %+v", targetRoutes, expectConfig.Route.Routes)
	}
}
