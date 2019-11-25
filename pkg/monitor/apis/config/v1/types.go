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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MonitorConfiguration contains the configuration for the Monitor
type MonitorConfiguration struct {
	metav1.TypeMeta

	Storage Storage `json:"storage"`
}

type Storage struct {
	// +optional
	InfluxDB *InfluxDBStorage `json:"influxDB,omitempty"`
	// +optional
	ElasticSearch *ElasticSearchStorage `json:"elasticSearch,omitempty"`
}

type InfluxDBStorage struct {
	Servers []InfluxDBStorageServer `json:"servers"`
}

type InfluxDBStorageServer struct {
	Address string `json:"address"`
	// +optional
	Username string `json:"username,omitempty"`
	// +optional
	Password string `json:"password,omitempty"`
	// +optional
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`
}

type ElasticSearchStorage struct {
	Servers []ElasticSearchStorageServer `json:"servers"`
}

type ElasticSearchStorageServer struct {
	Address string `json:"address"`
	// +optional
	Username string `json:"username,omitempty"`
	// +optional
	Password string `json:"password,omitempty"`
}
