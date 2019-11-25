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
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/sets"
	genericserveroptions "k8s.io/apiserver/pkg/server/options"
	"net"
	"os"
	"tkestack.io/tke/pkg/util/log"
)

const (
	flagAdvertiseAddress            = "advertise-address"
	flagCORSAllowedOrigins          = "cors-allowed-origins"
	flagTargetRAMMb                 = "target-ram-mb"
	flagExternalHostname            = "external-hostname"
	flagExternalPort                = "external-port"
	flagExternalScheme              = "external-scheme"
	flagExternalCAFile              = "external-ca-file"
	flagMasterServiceNamespace      = "master-service-namespace"
	flagMinRequestTimeout           = "min-request-timeout"
	flagRequestTimeout              = "request-timeout"
	flagMaxMutatingRequestsInflight = "max-mutating-requests-inflight"
	flagMaxRequestsInflight         = "max-requests-inflight"
)

const (
	configAdvertiseAddress            = "generic.advertise_address"
	configCORSAllowedOrigins          = "generic.cors_allowed_origins"
	configTargetRAMMb                 = "generic.target_ram_mb"
	configExternalHostname            = "generic.external_hostname"
	configExternalPort                = "generic.external_port"
	configExternalScheme              = "generic.external_scheme"
	configExternalCAFile              = "generic.external_ca_file"
	configMinRequestTimeout           = "generic.min_request_timeout"
	configRequestTimeout              = "generic.request_timeout"
	configMaxMutatingRequestsInflight = "generic.max_mutating_requests_inflight"
	configMaxRequestsInflight         = "generic.max_requests_inflight"
)

// GenericOptions contains the options while running a generic api server.
type GenericOptions struct {
	*genericserveroptions.ServerRunOptions
	ExternalPort   int
	ExternalScheme string
	ExternalCAFile string
}

// NewGenericOptions creates a Options object with default parameters.
func NewGenericOptions() *GenericOptions {
	return &GenericOptions{
		ServerRunOptions: genericserveroptions.NewServerRunOptions(),
		ExternalScheme:   "https",
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *GenericOptions) AddFlags(fs *pflag.FlagSet) {
	o.ServerRunOptions.AddUniversalFlags(fs)

	_ = fs.MarkHidden(flagMasterServiceNamespace)

	fs.Int(flagExternalPort, o.ExternalPort,
		"The port to use when generating externalized URLs for this server.")
	_ = viper.BindPFlag(configExternalPort, fs.Lookup(flagExternalPort))

	fs.String(flagExternalScheme, o.ExternalScheme,
		"The scheme to use when generating externalized URLs for this server.")
	_ = viper.BindPFlag(configExternalScheme, fs.Lookup(flagExternalScheme))

	fs.String(flagExternalCAFile, o.ExternalCAFile,
		"The CA file to use when generating externalized URLs for this server.")
	_ = viper.BindPFlag(configExternalCAFile, fs.Lookup(flagExternalCAFile))

	_ = viper.BindPFlag(configAdvertiseAddress, fs.Lookup(flagAdvertiseAddress))
	_ = viper.BindPFlag(configCORSAllowedOrigins, fs.Lookup(flagCORSAllowedOrigins))
	_ = viper.BindPFlag(configTargetRAMMb, fs.Lookup(flagTargetRAMMb))
	_ = viper.BindPFlag(configExternalHostname, fs.Lookup(flagExternalHostname))
	_ = viper.BindPFlag(configRequestTimeout, fs.Lookup(flagRequestTimeout))
	_ = viper.BindPFlag(configMaxMutatingRequestsInflight, fs.Lookup(flagMaxMutatingRequestsInflight))
	_ = viper.BindPFlag(configMaxRequestsInflight, fs.Lookup(flagMaxRequestsInflight))
	_ = viper.BindPFlag(configMinRequestTimeout, fs.Lookup(flagMinRequestTimeout))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *GenericOptions) ApplyFlags() []error {
	var errs []error

	o.AdvertiseAddress = net.ParseIP(viper.GetString(configAdvertiseAddress))
	o.CorsAllowedOriginList = viper.GetStringSlice(configCORSAllowedOrigins)
	o.TargetRAMMB = viper.GetInt(configTargetRAMMb)
	o.ExternalHost = viper.GetString(configExternalHostname)
	o.ExternalPort = viper.GetInt(configExternalPort)
	o.ExternalScheme = viper.GetString(configExternalScheme)
	o.ExternalCAFile = viper.GetString(configExternalCAFile)
	o.RequestTimeout = viper.GetDuration(configRequestTimeout)
	o.MaxMutatingRequestsInFlight = viper.GetInt(configMaxMutatingRequestsInflight)
	o.MaxRequestsInFlight = viper.GetInt(configMaxRequestsInflight)
	o.MinRequestTimeout = viper.GetInt(configMinRequestTimeout)

	if validateErrs := o.Validate(); len(validateErrs) > 0 {
		errs = append(errs, validateErrs...)
	}

	return errs
}

// DefaultAdvertiseAddress sets the field AdvertiseAddress if unset. The field
// will be set based on the SecureSecureServingOptions.
func (o *GenericOptions) DefaultAdvertiseAddress(secure *SecureServingOptions) error {
	if o.AdvertiseAddress == nil || o.AdvertiseAddress.IsUnspecified() {
		hostIP, err := netutil.ChooseBindAddress(secure.BindAddress)
		if err != nil {
			return fmt.Errorf("unable to find suitable network address.error='%v'. Try to set the AdvertiseAddress directly or provide a valid BindAddress to fix this", err)
		}
		o.AdvertiseAddress = hostIP
	}
	return nil
}

// DefaultAdvertiseAddressWithInsecure sets the field AdvertiseAddress if unset.
// The field will be set based on the SecureSecureServingOptions and
// InsecureServingOptions.
func (o *GenericOptions) DefaultAdvertiseAddressWithInsecure(secure *SecureServingOptions, insecure *InsecureServingOptions) error {
	if o.AdvertiseAddress == nil || o.AdvertiseAddress.IsUnspecified() {
		if secure.BindPort > 0 {
			hostIP, err := netutil.ChooseBindAddress(secure.BindAddress)
			if err != nil {
				return fmt.Errorf("unable to find suitable network address.error='%v'. Try to set the AdvertiseAddress directly or provide a valid secure BindAddress to fix this", err)
			}
			o.AdvertiseAddress = hostIP
		} else if insecure.BindPort > 0 {
			hostIP, err := netutil.ChooseBindAddress(insecure.BindAddress)
			if err != nil {
				return fmt.Errorf("unable to find suitable network address.error='%v'. Try to set the AdvertiseAddress directly or provide a valid insecure BindAddress to fix this", err)
			}
			o.AdvertiseAddress = hostIP
		}
	}
	return nil
}

// CompleteGenericAndSecureOptions is used to initialize the parameter values
// that are not filled in the common settings.
func CompleteGenericAndSecureOptions(genericOpts *GenericOptions, secureServingOpts *SecureServingOptions) error {
	// set defaults
	if err := genericOpts.DefaultAdvertiseAddress(secureServingOpts); err != nil {
		return err
	}
	if err := secureServingOpts.MaybeDefaultWithSelfSignedCerts(genericOpts.AdvertiseAddress.String(), []string{"localhost", "localhost.localdomain"}, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	if len(genericOpts.ExternalHost) == 0 {
		if len(genericOpts.AdvertiseAddress) > 0 {
			genericOpts.ExternalHost = genericOpts.AdvertiseAddress.String()
		} else {
			if hostname, err := os.Hostname(); err == nil {
				genericOpts.ExternalHost = hostname
			} else {
				return fmt.Errorf("error finding host name: %v", err)
			}
		}
		log.Infof("External host was not specified, using %v", genericOpts.ExternalHost)
	}

	if genericOpts.ExternalPort == 0 {
		genericOpts.ExternalPort = secureServingOpts.BindPort
		log.Infof("External port was not specified, using binding port %d", secureServingOpts.BindPort)
	}

	if genericOpts.ExternalScheme == "" {
		genericOpts.ExternalScheme = "https"
		log.Info("External scheme was not specified, using default scheme `HTTPS`")
	} else {
		schemes := sets.NewString("http", "https")
		if !schemes.Has(genericOpts.ExternalScheme) {
			return fmt.Errorf("error matching external scheme: %s, must be http or https", genericOpts.ExternalScheme)
		}
	}

	if genericOpts.ExternalScheme == "http" {
		if genericOpts.ExternalCAFile != "" {
			return fmt.Errorf("cannot set CA file when external exposure is HTTP")
		}
	} else {
		if genericOpts.ExternalCAFile == "" {
			log.Infof("External CA file was not specified, using server certificate file: %s", secureServingOpts.ServerCert.CertKey.CertFile)
			genericOpts.ExternalCAFile = secureServingOpts.ServerCert.CertKey.CertFile
		}
	}

	return nil
}
