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

package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	//TODO switch name back to prometheus_status_fail when rm platform addon platform
	prometheusStatusFail = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "prometheus_status_fail_monitor",
			Help: "prometheus addon status fail or not",
		},
		[]string{"tenant_id", "cluster_name", "prometheus_name"})
)

func init() {
	prometheus.MustRegister(prometheusStatusFail)
}

func UpdateMetricPrometheusStatusFail(tenantID string, clusterName string, prometheusName string, failed bool) {
	labels := map[string]string{"tenant_id": tenantID, "cluster_name": clusterName, "prometheus_name": prometheusName}
	if failed {
		prometheusStatusFail.With(labels).Set(1)
	} else {
		prometheusStatusFail.With(labels).Set(0)
	}
}

func DeleteMetricPrometheusStatusFail(tenantID string, clusterName string, prometheusName string) {
	labels := map[string]string{"tenant_id": tenantID, "cluster_name": clusterName, "prometheus_name": prometheusName}
	prometheusStatusFail.Delete(labels)
}
