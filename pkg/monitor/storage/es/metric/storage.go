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
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
	esclient "tkestack.io/tke/pkg/monitor/storage/es/client"
)

type ES struct {
	clients         []*esclient.Client
	availableClient *esclient.Client
}

func NewStorage(cfg *monitorconfig.ElasticSearchStorage) (*ES, error) {
	var clients []*esclient.Client
	for _, server := range cfg.Servers {
		clients = append(clients, &esclient.Client{
			URL:      server.Address,
			Username: server.Username,
			Password: server.Password,
		})
	}
	if len(clients) == 0 {
		return nil, fmt.Errorf("no available es client")
	}
	return &ES{
		clients:         clients,
		availableClient: clients[0],
	}, nil
}
