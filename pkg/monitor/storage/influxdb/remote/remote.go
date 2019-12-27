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

package remote

import (
	"time"

	influxclient "github.com/influxdata/influxdb1-client/v2"
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
)

type Client struct {
	Address string
	Client  influxclient.Client
}

func NewRemoteClients(cfg *monitorconfig.InfluxDBStorage) ([]Client, error) {
	var clients []Client
	for _, server := range cfg.Servers {
		influxCfg := influxclient.HTTPConfig{
			Addr:               server.Address,
			Username:           server.Username,
			Password:           server.Password,
			UserAgent:          "tke-monitor-controller",
			InsecureSkipVerify: true,
		}

		if server.TimeoutSeconds != nil {
			influxCfg.Timeout = time.Duration(*server.TimeoutSeconds) * time.Second
		}

		client, err := influxclient.NewHTTPClient(influxCfg)
		if err != nil {
			return nil, err
		}
		clients = append(clients, Client{
			Address: server.Address,
			Client:  client,
		})
	}
	return clients, nil
}
