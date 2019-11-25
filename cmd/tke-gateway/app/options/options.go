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
	"k8s.io/apimachinery/pkg/util/sets"
	"net"
	"os"
	apiserveroptions "tkestack.io/tke/pkg/apiserver/options"
	"tkestack.io/tke/pkg/util/log"
)

const (
	flagGatewayConfig   = "gateway-config"
	configGatewayConfig = "gateway_config"
)

// Options is the main context object for the TKE gateway server.
type Options struct {
	Log             *log.Options
	SecureServing   *apiserveroptions.SecureServingOptions
	InsecureServing *apiserveroptions.InsecureServingOptions
	Generic         *apiserveroptions.GenericOptions
	OIDC            *apiserveroptions.OIDCWithSecretOptions
	// The Gateway will load its initial configuration from this file.
	// The path may be absolute or relative; relative paths are under the Gateway's current working directory.
	GatewayConfig string
}

// NewOptions creates a new Options with a default config.
func NewOptions(serverName string) *Options {
	return &Options{
		Log:             log.NewOptions(),
		SecureServing:   apiserveroptions.NewSecureServingOptions(serverName, 9441),
		InsecureServing: apiserveroptions.NewInsecureServingOptions(9442),
		Generic:         apiserveroptions.NewGenericOptions(),
		OIDC:            apiserveroptions.NewOIDCWithSecretOptions(),
	}
}

// AddFlags adds flags for a specific server to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.Log.AddFlags(fs)
	o.SecureServing.AddFlags(fs)
	o.InsecureServing.AddFlags(fs)
	o.Generic.AddFlags(fs)
	o.OIDC.AddFlags(fs)

	fs.String(flagGatewayConfig, o.GatewayConfig,
		"The Gateway will load its initial configuration from this file. The path may be absolute or relative; relative paths start at the Gateway's current working directory. Omit this flag to use the built-in default configuration values.")
	_ = viper.BindPFlag(configGatewayConfig, fs.Lookup(flagGatewayConfig))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags() []error {
	var errs []error

	errs = append(errs, o.Log.ApplyFlags()...)
	errs = append(errs, o.SecureServing.ApplyFlags()...)
	errs = append(errs, o.InsecureServing.ApplyFlags()...)
	errs = append(errs, o.Generic.ApplyFlags()...)
	errs = append(errs, o.OIDC.ApplyFlags()...)

	o.GatewayConfig = viper.GetString(configGatewayConfig)

	return errs
}

// Complete set default Options.
// Should be called after tke-console flags parsed.
func (o *Options) Complete() error {
	// set defaults
	if err := o.Generic.DefaultAdvertiseAddressWithInsecure(o.SecureServing, o.InsecureServing); err != nil {
		return err
	}
	if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts(o.Generic.AdvertiseAddress.String(), []string{"localhost", "localhost.localdomain"}, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	if len(o.Generic.ExternalHost) == 0 {
		if len(o.Generic.AdvertiseAddress) > 0 {
			o.Generic.ExternalHost = o.Generic.AdvertiseAddress.String()
		} else {
			if hostname, err := os.Hostname(); err == nil {
				o.Generic.ExternalHost = hostname
			} else {
				return fmt.Errorf("error finding host name: %v", err)
			}
		}
		log.Infof("External host was not specified, using %v", o.Generic.ExternalHost)
	}

	if o.Generic.ExternalPort == 0 {
		if o.SecureServing.BindPort > 0 {
			o.Generic.ExternalPort = o.SecureServing.BindPort
			log.Infof("External port was not specified, using secure binding port %d", o.SecureServing.BindPort)
		} else if o.InsecureServing.BindPort > 0 {
			o.Generic.ExternalPort = o.InsecureServing.BindPort
			log.Infof("External port was not specified, using insecure binding port %d", o.InsecureServing.BindPort)
		}
	}

	if o.Generic.ExternalScheme == "" {
		o.Generic.ExternalScheme = "https"
		log.Info("External scheme was not specified, using default scheme `HTTPS`")
	} else {
		schemes := sets.NewString("http", "https")
		if !schemes.Has(o.Generic.ExternalScheme) {
			return fmt.Errorf("error matching external scheme: %s, must be http or https", o.Generic.ExternalScheme)
		}
	}

	return nil
}
