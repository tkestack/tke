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
	businessclient "tkestack.io/tke/api/client/clientset/versioned/typed/business/v1"
	"tkestack.io/tke/api/monitor"
	monitorconfig "tkestack.io/tke/pkg/monitor/apis/config"
	esmetric "tkestack.io/tke/pkg/monitor/storage/es/metric"
	esproject "tkestack.io/tke/pkg/monitor/storage/es/project"
	influxdbmetric "tkestack.io/tke/pkg/monitor/storage/influxdb/metric"
	influxdbproject "tkestack.io/tke/pkg/monitor/storage/influxdb/project"
	"tkestack.io/tke/pkg/monitor/storage/types"
)

type MetricStorage interface {
	Query(query *monitor.MetricQuery) (*types.MetricMergedResult, error)
}

type ProjectStorage interface {
	Collect()
}

func NewMetricStorage(storageConfig *monitorconfig.Storage) (MetricStorage, error) {
	if storageConfig.InfluxDB != nil {
		return influxdbmetric.NewStorage(storageConfig.InfluxDB)
	} else if storageConfig.ElasticSearch != nil {
		return esmetric.NewStorage(storageConfig.ElasticSearch)
	}
	return nil, fmt.Errorf("unregistered metric data storage type")
}

func NewProjectStorage(storageConfig *monitorconfig.Storage, businessClient businessclient.BusinessV1Interface) (ProjectStorage, error) {
	if storageConfig.InfluxDB != nil {
		return influxdbproject.NewStorage(storageConfig.InfluxDB, businessClient)
	} else if storageConfig.ElasticSearch != nil {
		return esproject.NewStorage(storageConfig.ElasticSearch, businessClient)
	}
	return nil, fmt.Errorf("unregistered project data storage type")
}
