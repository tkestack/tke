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
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	cliflag "k8s.io/component-base/cli/flag"
	"net"
	"strings"
)

const (
	flagSecureServingBindAddress                  = "bind-address"
	flagSecureServingBindPort                     = "secure-port"
	flagSecureServingCertDir                      = "cert-dir"
	flagSecureServingTLSCertFile                  = "tls-cert-file"
	flagSecureServingTLSKeyFile                   = "tls-private-key-file"
	flagSecureServingTLSCipherSuites              = "tls-cipher-suites"
	flagSecureServingTLSMinVersion                = "tls-min-version"
	flagSecureServingTLSSNICertKeys               = "tls-sni-cert-key"
	flagSecureServingHTTP2MaxStreamsPerConnection = "http2-max-streams-per-connection"
)

const (
	configSecureServingBindAddress                  = "secure_serving.bind_address"
	configSecureServingBindPort                     = "secure_serving.port"
	configSecureServingCertDir                      = "secure_serving.cert_dir"
	configSecureServingTLSCertFile                  = "secure_serving.tls_cert_file"
	configSecureServingTLSKeyFile                   = "secure_serving.tls_private_key_file"
	configSecureServingTLSCipherSuites              = "secure_serving.tls_cipher_suites"
	configSecureServingTLSMinVersion                = "secure_serving.tls_min_version"
	configSecureServingTLSSNICertKeys               = "secure_serving.tls_sni_cert_key"
	configSecureServingHTTP2MaxStreamsPerConnection = "secure_serving.http2_max_streams_per_connection"
)

// SecureServingOptions contains the options that serve HTTPS.
type SecureServingOptions struct {
	*apiserveroptions.SecureServingOptionsWithLoopback
}

// NewSecureServingOptions gives default values for the HTTPS server which are
// not the options wanted by "normal" servers running on the platform.
func NewSecureServingOptions(serverName string, defaultPort int) *SecureServingOptions {
	o := apiserveroptions.SecureServingOptions{
		BindAddress: net.ParseIP("0.0.0.0"),
		BindPort:    defaultPort,
		Required:    true,
		ServerCert: apiserveroptions.GeneratableKeyCert{
			PairName:      serverName,
			CertDirectory: "_output/certificates",
		},
	}
	return &SecureServingOptions{
		o.WithLoopback(),
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	o.SecureServingOptionsWithLoopback.AddFlags(fs)

	_ = viper.BindPFlag(configSecureServingBindAddress, fs.Lookup(flagSecureServingBindAddress))
	_ = viper.BindPFlag(configSecureServingBindPort, fs.Lookup(flagSecureServingBindPort))
	_ = viper.BindPFlag(configSecureServingCertDir, fs.Lookup(flagSecureServingCertDir))
	_ = viper.BindPFlag(configSecureServingTLSCertFile, fs.Lookup(flagSecureServingTLSCertFile))
	_ = viper.BindPFlag(configSecureServingTLSKeyFile, fs.Lookup(flagSecureServingTLSKeyFile))
	_ = viper.BindPFlag(configSecureServingTLSCipherSuites, fs.Lookup(flagSecureServingTLSCipherSuites))
	_ = viper.BindPFlag(configSecureServingTLSMinVersion, fs.Lookup(flagSecureServingTLSMinVersion))
	_ = viper.BindPFlag(configSecureServingTLSSNICertKeys, fs.Lookup(flagSecureServingTLSSNICertKeys))
	_ = viper.BindPFlag(configSecureServingHTTP2MaxStreamsPerConnection, fs.Lookup(flagSecureServingHTTP2MaxStreamsPerConnection))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *SecureServingOptions) ApplyFlags() []error {
	var errs []error

	o.BindAddress = net.ParseIP(viper.GetString(configSecureServingBindAddress))
	o.BindPort = viper.GetInt(configSecureServingBindPort)
	o.ServerCert.CertDirectory = viper.GetString(configSecureServingCertDir)
	o.ServerCert.CertKey.CertFile = viper.GetString(configSecureServingTLSCertFile)
	o.ServerCert.CertKey.KeyFile = viper.GetString(configSecureServingTLSKeyFile)
	o.CipherSuites = viper.GetStringSlice(configSecureServingTLSCipherSuites)
	o.MinTLSVersion = viper.GetString(configSecureServingTLSMinVersion)

	nck := cliflag.NewNamedCertKeyArray(&o.SNICertKeys)
	sniCertKeysString := viper.GetString(configSecureServingTLSSNICertKeys)
	sniCertKeysStringRaw := strings.TrimPrefix(strings.TrimSuffix(sniCertKeysString, "]"), "[")
	sniCertKeysStringArray := strings.Split(sniCertKeysStringRaw, ";")
	if len(sniCertKeysStringArray) > 0 {
		for _, sniCertKey := range sniCertKeysStringArray {
			if sniCertKey != "" {
				if err := nck.Set(sniCertKey); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	o.HTTP2MaxStreamsPerConnection = viper.GetInt(configSecureServingHTTP2MaxStreamsPerConnection)

	if validateErrs := o.Validate(); len(validateErrs) > 0 {
		errs = append(errs, validateErrs...)
	}

	return errs
}
