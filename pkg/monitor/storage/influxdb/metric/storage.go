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

package metric

import (
	"fmt"
	influxclient "github.com/influxdata/influxdb1-client/v2"
	"time"
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
)

type InfluxDB struct {
	clients         []influxclient.Client
	availableClient influxclient.Client
}

func NewStorage(cfg *monitorconfig.InfluxDBStorage) (*InfluxDB, error) {
	var clients []influxclient.Client
	for _, server := range cfg.Servers {
		influxCfg := influxclient.HTTPConfig{
			Addr:               server.Address,
			Username:           server.Username,
			Password:           server.Password,
			UserAgent:          "tke-monitor-api",
			InsecureSkipVerify: true,
		}

		if server.TimeoutSeconds != nil {
			influxCfg.Timeout = time.Duration(*server.TimeoutSeconds) * time.Second
		}

		client, err := influxclient.NewHTTPClient(influxCfg)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		return nil, fmt.Errorf("no available influxDB client")
	}
	return &InfluxDB{
		clients:         clients,
		availableClient: clients[0],
	}, nil
}
