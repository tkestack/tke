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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"time"
)

const (
	flagTokenWebhookConfigFile        = "token-webhook-config-file"
	flagTokenWebhookVersion          = "token-webhook-version"
	flagTokenWebhookCacheTTL          = "token-webhook-cache-ttl"
)

const (
	configTokenWebhookConfigFile         = "authentication.webhook.config_file"
	configTokenWehookVersion          = "authentication.webhook.version"
	configTokenWehookCacheTTL           = "authentication.webhook.cache_ttl"
)

type WebHookOptions struct {
	ConfigFile string
	Version    string
	CacheTTL   time.Duration
}


// NewWebhookOptions creates the default WebHookOptions object.
func NewWebhookOptions() *WebHookOptions {
	return &WebHookOptions{
		Version:  "v1beta1",
		CacheTTL: 2 * time.Minute,
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (w *WebHookOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&w.ConfigFile, flagTokenWebhookConfigFile, w.ConfigFile, ""+
		"File with webhook configuration for token authentication in kubeconfig format. "+
		"The API server will query the remote service to determine authentication for bearer tokens.")
	_ = viper.BindPFlag(configTokenWebhookConfigFile, fs.Lookup(flagTokenWebhookConfigFile))

	fs.StringVar(&w.Version, flagTokenWebhookVersion, w.Version, ""+
		"The API version of the authentication.k8s.io TokenReview to send to and expect from the webhook.")
	_ = viper.BindPFlag(configTokenWehookVersion, fs.Lookup(flagTokenWebhookVersion))

	fs.DurationVar(&w.CacheTTL, flagTokenWebhookCacheTTL, w.CacheTTL,
		"The duration to cache responses from the webhook token authenticator.")
	_ = viper.BindPFlag(configTokenWehookCacheTTL, fs.Lookup(flagTokenWebhookCacheTTL))

}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (w *WebHookOptions) ApplyFlags() []error {
	var errs []error
	w.ConfigFile = viper.GetString(configTokenWebhookConfigFile)
	w.Version = viper.GetString(configTokenWehookVersion)
	w.CacheTTL = viper.GetDuration(configTokenWehookCacheTTL)

	return errs
}

