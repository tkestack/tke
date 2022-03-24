/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package options

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	appconfig "tkestack.io/tke/pkg/application/config"
)

const (
	flagConcurrentSyncs   = "concurrent-app-syncs"
	flagSyncAppPeriod     = "sync-app-period"
)

const (
	configSyncAppPeriod      = "controller.sync_app_period"
	configConcurrentAppSyncs = "controller.concurrent_app_syncs"
)


// ControllerOptions contains configuration items related to application attributes.
type ControllerOptions struct {
	SyncAppPeriod      time.Duration
	ConcurrentAppSyncs int
}

// FeatureOptions contains configuration items related to application attributes.
type FeatureOptions struct {
	Controller ControllerOptions
}

// NewFeatureOptions creates a FeatureOptions object with default parameters.
func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{
		Controller: ControllerOptions{
			SyncAppPeriod:      defaultSyncPeriod,
			ConcurrentAppSyncs: defaultconcurrentSyncs,
		},
	}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&o.Controller.SyncAppPeriod, flagSyncAppPeriod, o.Controller.SyncAppPeriod, "The period for app health checks")
	_ = viper.BindPFlag(configSyncAppPeriod, fs.Lookup(flagSyncAppPeriod))

	fs.IntVar(&o.Controller.ConcurrentAppSyncs, flagConcurrentSyncs, o.Controller.ConcurrentAppSyncs, "The number of app objects that are allowed to sync concurrently. Larger number = more responsive app termination, but more CPU (and network) load")
	_ = viper.BindPFlag(configConcurrentAppSyncs, fs.Lookup(flagConcurrentSyncs))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *FeatureOptions) ApplyFlags() []error {
	var errs []error

	o.Controller.SyncAppPeriod = viper.GetDuration(configSyncAppPeriod)
	o.Controller.ConcurrentAppSyncs = viper.GetInt(configConcurrentAppSyncs)
	return errs
}


// ApplyTo fills up Debugging config with options.
func (o *ControllerOptions) ApplyTo(cfg *appconfig.AppControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentSyncs = o.ConcurrentAppSyncs
	cfg.SyncPeriod = o.SyncAppPeriod

	return nil
}
