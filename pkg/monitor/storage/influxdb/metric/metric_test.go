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
	"encoding/json"
	"github.com/influxdata/influxdb1-client/models"
	"strconv"
	"testing"
	"tkestack.io/tke/pkg/monitor/storage/types"
)

func TestMergeResult(t *testing.T) {
	var input []types.MetricResult
	var res1 types.MetricResult
	series1 := make([]models.Row, 3)
	for i := range series1 {
		series1[i].Name = "k8s_pod_status_ready"
		series1[i].Columns = []string{"time", "k8s_pod_status_ready"}
		series1[i].Tags = make(map[string]string)
		series1[i].Tags["pod_name"] = "ajwdlkj" + strconv.Itoa(i)
		series1[i].Values = make([][]interface{}, 7)
		for j := range series1[i].Values {
			series1[i].Values[j] = []interface{}{json.Number(1565770000 + j), 1}
		}
	}
	var res2 types.MetricResult
	series2 := make([]models.Row, 3)
	for i := range series2 {
		series2[i].Name = "k8s_pod_cpu_core_used"
		series2[i].Columns = []string{"time", "k8s_pod_cpu_core_used"}
		series2[i].Tags = make(map[string]string)
		series2[i].Tags["pod_name"] = "ajwdlkj" + strconv.Itoa(i)
		series2[i].Values = make([][]interface{}, 7)
		for j := range series2[i].Values {
			series2[i].Values[j] = []interface{}{json.Number(1565770000 + j), 0.5}
		}
	}
	res1.Series = series1
	res2.Series = series2
	input = append(input, res1)
	input = append(input, res2)
	res := MergeResult(input, []string{"timestamp(1s)", "pod_name"}, []string{"max(k8s_pod_status_ready)", "mean(k8s_pod_cpu_core_used)"})
	j, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("failed to marshal %v", res)
	}

	expect := `{"columns":["timestamp(1s)","k8s_pod_status_ready_max","k8s_pod_cpu_core_used_mean","pod_name"],"data":[[0,1,0.5,"ajwdlkj0"],[0,1,0.5,"ajwdlkj1"],[0,1,0.5,"ajwdlkj2"],[0,1,0.5,"ajwdlkj0"],[0,1,0.5,"ajwdlkj1"],[0,1,0.5,"ajwdlkj2"],[0,1,0.5,"ajwdlkj0"],[0,1,0.5,"ajwdlkj1"],[0,1,0.5,"ajwdlkj2"],[0,1,0.5,"ajwdlkj0"],[0,1,0.5,"ajwdlkj1"],[0,1,0.5,"ajwdlkj2"],[0,1,0.5,"ajwdlkj0"],[0,1,0.5,"ajwdlkj1"],[0,1,0.5,"ajwdlkj2"],[0,1,0.5,"ajwdlkj0"],[0,1,0.5,"ajwdlkj1"],[0,1,0.5,"ajwdlkj2"],[0,1,0.5,"ajwdlkj0"],[0,1,0.5,"ajwdlkj1"],[0,1,0.5,"ajwdlkj2"]]}`
	if string(j) != expect {
		t.Errorf("expect %s, got %s", expect, string(j))
	}
}
