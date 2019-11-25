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

package validation

import (
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"testing"
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
)

func TestValidateMonitorConfiguration(t *testing.T) {
	successCase := &monitorconfig.MonitorConfiguration{
		Storage: monitorconfig.Storage{
			InfluxDB: &monitorconfig.InfluxDBStorage{
				Servers: []monitorconfig.InfluxDBStorageServer{
					{
						Address: "https://127.0.0.1:8080",
					},
				},
			},
		},
	}
	if allErrors := ValidateMonitorConfiguration(successCase); allErrors != nil {
		t.Errorf("expect no errors, got %v", allErrors)
	}

	errorCase := &monitorconfig.MonitorConfiguration{
		Storage: monitorconfig.Storage{
			InfluxDB: &monitorconfig.InfluxDBStorage{
				Servers: []monitorconfig.InfluxDBStorageServer{
					{
						Address: "https://127.0.0.1:8080",
					},
				},
			},
			ElasticSearch: &monitorconfig.ElasticSearchStorage{
				Servers: []monitorconfig.ElasticSearchStorageServer{
					{
						Address: "https://127.0.0.1:8080",
					}, {
						Username: "fake",
					},
				},
			},
		},
	}
	const numErrs = 2
	if allErrors := ValidateMonitorConfiguration(errorCase); len(allErrors.(utilerrors.Aggregate).Errors()) != numErrs {
		t.Errorf("expect %d errors, got %v", numErrs, len(allErrors.(utilerrors.Aggregate).Errors()))
	}
}
