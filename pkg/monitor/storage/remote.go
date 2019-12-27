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

package storage

import (
	"fmt"

	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
	esclient "tkestack.io/tke/pkg/monitor/storage/es/client"
	esremote "tkestack.io/tke/pkg/monitor/storage/es/remote"
	influxdbremote "tkestack.io/tke/pkg/monitor/storage/influxdb/remote"
)

// RemoteClient wrap influxdb and es client.
type RemoteClient struct {
	InfluxDB []influxdbremote.Client
	ES       []esclient.Client
}

// NewRemoteClient creates RemoteClient object by given monitor storage
// configuration and return it.
func NewRemoteClient(storageConfig *monitorconfig.Storage) (*RemoteClient, error) {
	remoteClient := &RemoteClient{}
	if storageConfig.InfluxDB != nil {
		clients, err := influxdbremote.NewRemoteClients(storageConfig.InfluxDB)
		if err != nil {
			return nil, err
		}
		remoteClient.InfluxDB = clients
	} else if storageConfig.ElasticSearch != nil {
		clients := esremote.NewRemoteClients(storageConfig.ElasticSearch)
		remoteClient.ES = clients
	}

	if len(remoteClient.ES) == 0 || len(remoteClient.InfluxDB) == 0 {
		return nil, fmt.Errorf("unregistered remote monitor data storage type")
	}
	return remoteClient, nil
}
