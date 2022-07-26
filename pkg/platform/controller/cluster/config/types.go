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

package config

import "time"

// ClusterControllerConfiguration contains elements describing ClusterController.
type ClusterControllerConfiguration struct {
	// clusterSyncPeriod is the period for syncing cluster life-cycle
	// updates.
	ClusterSyncPeriod time.Duration
	// concurrentClusterSyncs is the number of cluster objects that are
	// allowed to sync concurrently.
	ConcurrentClusterSyncs int
	// HealthCheckPeriod is the period for cluster health check
	HealthCheckPeriod time.Duration
	// When the number of clusters is relatively large,
	// the health check of all clusters takes a long time.
	// You can set a random health check time for each cluster
	RandomeRangeLowerLimitForHealthCheckPeriod time.Duration
	RandomeRangeUpperLimitForHealthCheckPeriod time.Duration
	// BucketRateLimiterLimit allows events up to rate r and permits.
	BucketRateLimiterLimit int
	// BucketRateLimiterBurst bursts of at most b tokens.
	BucketRateLimiterBurst int
	// IsCRDMode Whether the controller is using CRD mode
	IsCRDMode bool
}
