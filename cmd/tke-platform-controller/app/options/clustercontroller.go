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

package options

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	clusterconfig "tkestack.io/tke/pkg/platform/controller/cluster/config"
)

const (
	flagClusterSyncPeriod       = "cluster-sync-period"
	flagConcurrentClusterSyncs  = "concurrent-cluster-syncs"
	flagHealthCheckPeriod       = "healthcheck-period"
	flagClusterRateLimiterLimit = "cluster-rate-limiter-limit"
	flagClusterRateLimiterBurst = "cluster-rate-limiter-burst"
)

const (
	configClusterSyncPeriod       = "controller.cluster_sync_period"
	configConcurrentClusterSyncs  = "controller.concurrent_cluster_syncs"
	configHealthCheckPeriod       = "controller.healthcheck_period"
	configClusterRateLimiterLimit = "controller.cluster_rate_limiter_limit"
	configClusterRateLimiterBurst = "controller.cluster_rate_limiter_burst"
)

// ClusterControllerOptions holds the ClusterController options.
type ClusterControllerOptions struct {
	*clusterconfig.ClusterControllerConfiguration
}

// NewClusterControllerOptions creates a new Options with a default config.
func NewClusterControllerOptions() *ClusterControllerOptions {
	return &ClusterControllerOptions{
		&clusterconfig.ClusterControllerConfiguration{
			ClusterSyncPeriod:      defaultSyncPeriod,
			ConcurrentClusterSyncs: defaultConcurrentSyncs,
			HealthCheckPeriod:      defaultHealthCheckPeriod,
			BucketRateLimiterLimit: defaultBucketRateLimiterLimit,
			BucketRateLimiterBurst: defaultBucketRateLimiterBurst,
		},
	}
}

// AddFlags adds flags related to ClusterController for controller manager to the specified FlagSet.
func (o *ClusterControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.DurationVar(&o.ClusterSyncPeriod, flagClusterSyncPeriod, o.ClusterSyncPeriod, "The period for syncing cluster life-cycle updates")
	_ = viper.BindPFlag(configClusterSyncPeriod, fs.Lookup(flagClusterSyncPeriod))
	fs.IntVar(&o.ConcurrentClusterSyncs, flagConcurrentClusterSyncs, o.ConcurrentClusterSyncs, "The number of cluster objects that are allowed to sync concurrently. Larger number = more responsive cluster termination, but more CPU (and network) load")
	_ = viper.BindPFlag(configConcurrentClusterSyncs, fs.Lookup(flagConcurrentClusterSyncs))
	fs.DurationVar(&o.HealthCheckPeriod, flagHealthCheckPeriod, o.HealthCheckPeriod, "The period for cluster health check")
	_ = viper.BindPFlag(configHealthCheckPeriod, fs.Lookup(flagHealthCheckPeriod))
	fs.IntVar(&o.BucketRateLimiterLimit, flagClusterRateLimiterLimit, o.BucketRateLimiterLimit, "The number of allows events up to rate r and permits.")
	_ = viper.BindPFlag(configClusterRateLimiterLimit, fs.Lookup(flagClusterRateLimiterLimit))
	fs.IntVar(&o.BucketRateLimiterBurst, flagClusterRateLimiterBurst, o.BucketRateLimiterBurst, "The number of bursts of at most b tokens.")
	_ = viper.BindPFlag(configClusterRateLimiterBurst, fs.Lookup(flagClusterRateLimiterBurst))
}

// ApplyTo fills up ClusterController config with options.
func (o *ClusterControllerOptions) ApplyTo(cfg *clusterconfig.ClusterControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ClusterSyncPeriod = o.ClusterSyncPeriod
	cfg.ConcurrentClusterSyncs = o.ConcurrentClusterSyncs
	cfg.HealthCheckPeriod = o.HealthCheckPeriod
	cfg.BucketRateLimiterLimit = o.BucketRateLimiterLimit
	cfg.BucketRateLimiterBurst = o.BucketRateLimiterBurst

	return nil
}

// Validate checks validation of ClusterControllerOptions.
func (o *ClusterControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *ClusterControllerOptions) ApplyFlags() []error {
	o.ClusterSyncPeriod = viper.GetDuration(configClusterSyncPeriod)
	o.ConcurrentClusterSyncs = viper.GetInt(configConcurrentClusterSyncs)
	o.HealthCheckPeriod = viper.GetDuration(configHealthCheckPeriod)
	o.BucketRateLimiterLimit = viper.GetInt(configClusterRateLimiterLimit)
	o.BucketRateLimiterBurst = viper.GetInt(configClusterRateLimiterBurst)
	return nil
}
