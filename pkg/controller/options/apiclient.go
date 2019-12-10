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

// APIServerClientOptions holds the platform apiserver client options.
// If it is not required, the parameter verification will not determine whether
// the address of the apiserver or the configuration file address has been
// specified.
type APIServerClientOptions struct {
	Server             string
	ServerClientConfig string
	ContentType        string
	QPS                float32
	Burst              int32

	Name     string
	Required bool

	flagAPIClientContentType        string
	flagAPIClientQPS                string
	flagAPIClientTimeout            string
	flagAPIClientBurst              string
	flagAPIClientServer             string
	flagAPIClientServerClientConfig string

	configAPIClientContentType        string
	configAPIClientQPS                string
	configAPIClientTimeout            string
	configAPIClientBurst              string
	configAPIClientServer             string
	configAPIClientServerClientConfig string
}

// NewAPIServerClientOptions creates the default APIServerClientOptions object.
func NewAPIServerClientOptions(name string, required bool) *APIServerClientOptions {
	return &APIServerClientOptions{
		ContentType: "application/vnd.kubernetes.protobuf",
		QPS:         20.0,
		Burst:       30,

		Name:     name,
		Required: required,

		flagAPIClientContentType:        fmt.Sprintf("%s-api-content-type", name),
		flagAPIClientQPS:                fmt.Sprintf("%s-api-qps", name),
		flagAPIClientTimeout:            fmt.Sprintf("%s-api-timeout", name),
		flagAPIClientBurst:              fmt.Sprintf("%s-api-burst", name),
		flagAPIClientServer:             FlagAPIClientServer(name),
		flagAPIClientServerClientConfig: FlagAPIClientServerClientConfig(name),

		configAPIClientContentType:        fmt.Sprintf("client.%s.api_content_type", name),
		configAPIClientQPS:                fmt.Sprintf("client.%s.api_qps", name),
		configAPIClientTimeout:            fmt.Sprintf("client.%s.api_timeout", name),
		configAPIClientBurst:              fmt.Sprintf("client.%s.api_burst", name),
		configAPIClientServer:             fmt.Sprintf("client.%s.api_server", name),
		configAPIClientServerClientConfig: fmt.Sprintf("client.%s.api_server_client_config", name),
	}
}

// AddFlags adds flags related to debugging for API client to the specified FlagSet.
func (o *APIServerClientOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.String(o.flagAPIClientServer, o.Server,
		"The address of the "+o.Name+" apiserver (overrides any value in "+o.flagAPIClientServerClientConfig+").")
	_ = viper.BindPFlag(o.configAPIClientServer, fs.Lookup(o.flagAPIClientServer))
	fs.String(o.flagAPIClientServerClientConfig, o.ServerClientConfig,
		"Path to config file with authorization and "+o.Name+" apiserver location information.")
	_ = viper.BindPFlag(o.configAPIClientServerClientConfig, fs.Lookup(o.flagAPIClientServerClientConfig))
	fs.String(o.flagAPIClientContentType, o.ContentType,
		"Content type of requests sent to "+o.Name+" apiserver.")
	_ = viper.BindPFlag(o.configAPIClientContentType, fs.Lookup(o.flagAPIClientContentType))
	fs.Float32(o.flagAPIClientQPS, o.QPS,
		"QPS to use while talking with "+o.Name+" apiserver.")
	_ = viper.BindPFlag(o.configAPIClientQPS, fs.Lookup(o.flagAPIClientQPS))
	fs.Int32(o.flagAPIClientBurst, o.Burst,
		"Burst to use while talking with "+o.Name+" apiserver.")
	_ = viper.BindPFlag(o.configAPIClientBurst, fs.Lookup(o.flagAPIClientBurst))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *APIServerClientOptions) ApplyFlags() []error {
	var errs []error

	o.ServerClientConfig = viper.GetString(o.configAPIClientServerClientConfig)
	o.Server = viper.GetString(o.configAPIClientServer)
	o.Burst = viper.GetInt32(o.configAPIClientBurst)
	o.QPS = float32(viper.GetFloat64(o.configAPIClientQPS))
	o.ContentType = viper.GetString(o.configAPIClientContentType)

	if o.Required {
		if o.ServerClientConfig == "" && o.Server == "" {
			errs = append(errs, fmt.Errorf("must specify either `%s` or `%s`", FlagAPIClientServer(o.Name), FlagAPIClientServerClientConfig(o.Name)))
		}
	}

	return errs
}

// FlagAPIClientServer returns API client server flag by given name.
func FlagAPIClientServer(name string) string {
	return fmt.Sprintf("%s-api-server", name)
}

// FlagAPIClientServerClientConfig returns API client server config flag by
// given name.
func FlagAPIClientServerClientConfig(name string) string {
	return fmt.Sprintf("%s-api-server-client-config", name)
}
