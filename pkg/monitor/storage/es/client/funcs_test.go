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

package client

import (
	"testing"
)

func TestGetTables(t *testing.T) {
	indices := []string{
		"prometheusbeat-6.4.1-2019.08.12",
		"prometheusbeat-6.4.1-2019.05.22",
		"prometheusbeat-6.4.1-2020.10.31",
	}
	tableName := "prometheusbeat-6.4.1"
	// millisecond timestamp
	// 2019/5/21 17:16:00
	var startTime int64 = 1558430160000
	// 2019/8/16 19:45:58
	var endTime int64 = 1565955958000

	tables := GetTablesMonitor(indices, tableName, startTime, endTime)
	expectedTables := "prometheusbeat-6.4.1-2019.05.22,prometheusbeat-6.4.1-2019.08.12"

	if tables != expectedTables {
		t.Fatal("GetTables is error")
	}
}
