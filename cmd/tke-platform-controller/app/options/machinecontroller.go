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

	machineconfig "tkestack.io/tke/pkg/platform/controller/machine/config"
)

const (
	flagMachineSyncPeriod       = "machine-sync-period"
	flagConcurrentMachineSyncs  = "concurrent-machine-syncs"
	flagMachineRateLimiterLimit = "machine-rate-limiter-limit"
	flagMachineRateLimiterBurst = "machine-rate-limiter-burst"
)

const (
	configMachineSyncPeriod       = "controller.machine_sync_period"
	configConcurrentMachineSyncs  = "controller.concurrent_machine_syncs"
	configMachineRateLimiterLimit = "controller.machine_rate_limiter_limit"
	configMachineRateLimiterBurst = "controller.machine_rate_limiter_burst"
)

// MachineControllerOptions holds the MachineController options.
type MachineControllerOptions struct {
	*machineconfig.MachineControllerConfiguration
}

// NewMachineControllerOptions creates a new Options with a default config.
func NewMachineControllerOptions() *MachineControllerOptions {
	return &MachineControllerOptions{
		&machineconfig.MachineControllerConfiguration{
			MachineSyncPeriod:      defaultSyncPeriod,
			ConcurrentMachineSyncs: defaultConcurrentSyncs,
			BucketRateLimiterLimit: defaultBucketRateLimiterLimit,
			BucketRateLimiterBurst: defaultBucketRateLimiterBurst,
		},
	}
}

// AddFlags adds flags related to MachineController for controller manager to the specified FlagSet.
func (o *MachineControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.DurationVar(&o.MachineSyncPeriod, flagMachineSyncPeriod, o.MachineSyncPeriod, "The period for syncing machine life-cycle updates")
	_ = viper.BindPFlag(configMachineSyncPeriod, fs.Lookup(flagMachineSyncPeriod))
	fs.IntVar(&o.ConcurrentMachineSyncs, flagConcurrentMachineSyncs, o.ConcurrentMachineSyncs, "The number of machine objects that are allowed to sync concurrently. Larger number = more responsive machine termination, but more CPU (and network) load")
	_ = viper.BindPFlag(configConcurrentMachineSyncs, fs.Lookup(flagConcurrentMachineSyncs))
	fs.IntVar(&o.BucketRateLimiterLimit, flagMachineRateLimiterLimit, o.BucketRateLimiterLimit, "The number of allows events up to rate r and permits.")
	_ = viper.BindPFlag(configMachineRateLimiterLimit, fs.Lookup(flagMachineRateLimiterLimit))
	fs.IntVar(&o.BucketRateLimiterBurst, flagMachineRateLimiterBurst, o.BucketRateLimiterBurst, "The number of bursts of at most b tokens.")
	_ = viper.BindPFlag(configMachineRateLimiterBurst, fs.Lookup(flagMachineRateLimiterBurst))
}

// ApplyTo fills up MachineController config with options.
func (o *MachineControllerOptions) ApplyTo(cfg *machineconfig.MachineControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.MachineSyncPeriod = o.MachineSyncPeriod
	cfg.ConcurrentMachineSyncs = o.ConcurrentMachineSyncs
	cfg.BucketRateLimiterLimit = o.BucketRateLimiterLimit
	cfg.BucketRateLimiterBurst = o.BucketRateLimiterBurst

	return nil
}

// Validate checks validation of MachineControllerOptions.
func (o *MachineControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *MachineControllerOptions) ApplyFlags() []error {
	o.MachineSyncPeriod = viper.GetDuration(configMachineSyncPeriod)
	o.ConcurrentMachineSyncs = viper.GetInt(configConcurrentMachineSyncs)
	o.BucketRateLimiterLimit = viper.GetInt(configMachineRateLimiterLimit)
	o.BucketRateLimiterBurst = viper.GetInt(configMachineRateLimiterBurst)
	return nil
}
