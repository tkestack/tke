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
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagMonitorStorageType      = "monitor-storage-type"
	flagMonitorStorageAddresses = "monitor-storage-addresses"
)

const (
	configMonitorStorageType      = "features.monitor_storage_type"
	configMonitorStorageAddresses = "features.monitor_storage_addresses"
)

type FeatureOptions struct {
	MonitorStorageType      string
	MonitorStorageAddresses []string
}

func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{}
}

func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String(flagMonitorStorageType, o.MonitorStorageType,
		"The type of storage for monitor. Support influxdb and elasticsearch.")
	_ = viper.BindPFlag(configMonitorStorageType, fs.Lookup(flagMonitorStorageType))
	fs.StringSlice(flagMonitorStorageAddresses, o.MonitorStorageAddresses,
		"Multiple addresses of storage for monitor. Include username, password and server url.")
	_ = viper.BindPFlag(configMonitorStorageAddresses, fs.Lookup(flagMonitorStorageAddresses))
}

func (o *FeatureOptions) ApplyFlags() []error {
	var errs []error

	o.MonitorStorageAddresses = viper.GetStringSlice(configMonitorStorageAddresses)
	o.MonitorStorageType = viper.GetString(configMonitorStorageType)

	switch o.MonitorStorageType {
	case "":
	case "influxdb", "elasticsearch", "es", "influxDB":
		if len(o.MonitorStorageAddresses) == 0 {
			errs = append(errs, fmt.Errorf("must specify %s server address", o.MonitorStorageType))
		}
	default:
		errs = append(errs, fmt.Errorf("unsupported storage type for monitor"))
	}

	return errs
}
