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
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagAuthzPolicyFile                  = "authorization-policy-file"
	flagAuthzWebhookConfigFile           = "authorization-webhook-config-file"
	flagAuthzWebhookCacheUnauthorizedTTL = "authorization-webhook-cache-unauthorized-ttl"
	flagAuthzWebhookCacheAuthorizedTTL   = "authorization-webhook-cache-authorized-ttl"
	flagAuthzDebug                       = "authorization-debug"
	flagCasbinModelFile                  = "casbin-model-file"
	flagCasbinReLoadInterval             = "casbin-reload-interval"
)

const (
	configAuthzPolicyFile                  = "authorization.policy_file"
	configAuthzWebhookConfigFile           = "authorization.webhook_config_file"
	configAuthzWebhookCacheUnauthorizedTTL = "authorization.webhook_cache_unauthorized_ttl"
	configAuthzWebhookCacheAuthorizedTTL   = "authorization.webhook_cache_authorized_ttl"
	configAuthzDebug                       = "authorization.debug"
	configCasbinModelFile                  = "casbin.model_file"
	configCasbinReloadInterval             = "casbin.reload_interval"
)

// AuthorizationOptions contains configuration items related to authorization.
type AuthorizationOptions struct {
	CasbinModelFile             string
	CasbinReloadInterval        time.Duration
	Debug                       bool
	PolicyFile                  string
	WebhookConfigFile           string
	WebhookCacheAuthorizedTTL   time.Duration
	WebhookCacheUnauthorizedTTL time.Duration
}

// NewAuthorizationOptions creates a AuthorizationOptions object with default
// parameters.
func NewAuthorizationOptions() *AuthorizationOptions {
	return &AuthorizationOptions{CasbinReloadInterval: 100 * time.Millisecond}
}

// AddFlags adds flags related to authenticate for a specific APIServer to the
// specified FlagSet
func (o *AuthorizationOptions) AddFlags(fs *pflag.FlagSet) {

	fs.String(flagAuthzPolicyFile, o.PolicyFile, ""+
		"File with authorization policy in json line by line format, on the secure port.")
	_ = viper.BindPFlag(configAuthzPolicyFile, fs.Lookup(flagAuthzPolicyFile))

	fs.String(flagAuthzWebhookConfigFile, o.WebhookConfigFile, ""+
		"File with webhook configuration in kubeconfig format. "+
		"The API server will query the remote service to determine access on the API server's secure port.")
	_ = viper.BindPFlag(configAuthzWebhookConfigFile, fs.Lookup(flagAuthzWebhookConfigFile))

	fs.Duration(flagAuthzWebhookCacheAuthorizedTTL, o.WebhookCacheAuthorizedTTL, ""+
		"The duration to cache 'authorized' responses from the webhook authorizer.")
	_ = viper.BindPFlag(configAuthzWebhookCacheAuthorizedTTL, fs.Lookup(flagAuthzWebhookCacheAuthorizedTTL))

	fs.Duration(flagAuthzWebhookCacheUnauthorizedTTL, o.WebhookCacheUnauthorizedTTL,
		"The duration to cache 'unauthorized' responses from the webhook authorizer.")
	_ = viper.BindPFlag(configAuthzWebhookCacheUnauthorizedTTL, fs.Lookup(flagAuthzWebhookCacheUnauthorizedTTL))

	fs.String(flagCasbinModelFile, o.CasbinModelFile,
		"Casbin model file used to store ACL model.")
	_ = viper.BindPFlag(configCasbinModelFile, fs.Lookup(flagCasbinModelFile))

	fs.Duration(flagCasbinReLoadInterval, o.CasbinReloadInterval,
		"The interval of casbin reload policy from backend storage. Default 5s.")
	_ = viper.BindPFlag(configCasbinReloadInterval, fs.Lookup(flagCasbinReLoadInterval))

	fs.Bool(flagAuthzDebug, o.Debug,
		"Enable authorizer to log messages to the Logger.")
	_ = viper.BindPFlag(configAuthzDebug, fs.Lookup(flagAuthzDebug))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *AuthorizationOptions) ApplyFlags() []error {
	var errs []error

	o.PolicyFile = viper.GetString(configAuthzPolicyFile)
	o.CasbinModelFile = viper.GetString(configCasbinModelFile)
	o.CasbinReloadInterval = viper.GetDuration(configCasbinReloadInterval)
	o.Debug = viper.GetBool(configAuthzDebug)
	o.WebhookCacheAuthorizedTTL = viper.GetDuration(configAuthzWebhookCacheAuthorizedTTL)
	o.WebhookCacheUnauthorizedTTL = viper.GetDuration(configAuthzWebhookCacheUnauthorizedTTL)
	o.WebhookConfigFile = viper.GetString(configAuthzWebhookConfigFile)

	return errs
}
