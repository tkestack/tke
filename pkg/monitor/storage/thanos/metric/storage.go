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
	"github.com/prometheus/client_golang/api"
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
	"tkestack.io/tke/pkg/util/log"
)

type Thanos struct {
	clients         []api.Client
	availableClient api.Client
}

func NewStorage(cfg *monitorconfig.ThanosStorage) (*Thanos, error) {
	var clients []api.Client
	for _, s := range cfg.Servers {
		client, err := api.NewClient(api.Config{
			Address: s.Address,
		})
		if err != nil {
			log.Errorf("Error creating client: %v", err)
			continue
		}
		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("no available thanos client")
	}
	return &Thanos{
		clients:         clients,
		availableClient: clients[0],
	}, nil
}
