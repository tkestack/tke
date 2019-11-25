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
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"net"
)

const (
	flagInsecureServingBindAddress = "insecure-bind-address"
	flagInsecureServingBindPort    = "insecure-port"
)

const (
	configInsecureServingBindAddress = "insecure_serving.bind_address"
	configInsecureServingBindPort    = "insecure_serving.port"
)

// InsecureServingOptions contains the options that serve HTTP.
type InsecureServingOptions struct {
	*genericoptions.DeprecatedInsecureServingOptions
}

// NewInsecureServingOptions gives default values for the http server which are
// not the options wanted by "normal" servers running on the platform.
func NewInsecureServingOptions(defaultPort int) *InsecureServingOptions {
	o := &genericoptions.DeprecatedInsecureServingOptions{
		BindAddress: net.ParseIP("0.0.0.0"),
		BindPort:    defaultPort,
	}
	return &InsecureServingOptions{o}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *InsecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.IP(flagInsecureServingBindAddress, o.BindAddress, ""+
		"The IP address on which to serve the --insecure-port (set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	_ = viper.BindPFlag(configInsecureServingBindAddress, fs.Lookup(flagInsecureServingBindAddress))

	fs.Int(flagInsecureServingBindPort, o.BindPort, "The port on which to serve HTTP. If 0, don't serve HTTP at all.")
	_ = viper.BindPFlag(configInsecureServingBindPort, fs.Lookup(flagInsecureServingBindPort))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *InsecureServingOptions) ApplyFlags() []error {
	var errs []error

	o.BindAddress = net.ParseIP(viper.GetString(configInsecureServingBindAddress))
	o.BindPort = viper.GetInt(configInsecureServingBindPort)

	if o.BindPort > 65535 {
		errs = append(errs, fmt.Errorf("insecure port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port", o.BindPort))
	}

	return errs
}
