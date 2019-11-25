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
	"strings"
	"time"
	"tkestack.io/tke/pkg/apiserver/authorization/modes"
)

const (
	flagAuthzMode                        = "authorization-mode"
	flagAuthzWebhookConfigFile           = "authorization-webhook-config-file"
	flagAuthzWebhookCacheUnauthorizedTTL = "authorization-webhook-cache-unauthorized-ttl"
	flagAuthzWebhookCacheAuthorizedTTL   = "authorization-webhook-cache-authorized-ttl"
)

const (
	configAuthzMode                        = "authorization.mode"
	configAuthzWebhookConfigFile           = "authorization.webhook_config_file"
	configAuthzWebhookCacheUnauthorizedTTL = "authorization.webhook_cache_unauthorized_ttl"
	configAuthzWebhookCacheAuthorizedTTL   = "authorization.webhook_cache_authorized_ttl"
)

// AuthorizationOptions defines the configuration parameters required to
// include the authorization.
type AuthorizationOptions struct {
	Modes                       []string
	WebhookConfigFile           string
	WebhookCacheAuthorizedTTL   time.Duration
	WebhookCacheUnauthorizedTTL time.Duration
}

// NewAuthorizationOptions creates the default AuthorizationOptions object and
// returns it.
func NewAuthorizationOptions() *AuthorizationOptions {
	return &AuthorizationOptions{
		Modes:                       []string{modes.ModeAlwaysAllow},
		WebhookCacheAuthorizedTTL:   5 * time.Minute,
		WebhookCacheUnauthorizedTTL: 30 * time.Second,
	}
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *AuthorizationOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSlice(flagAuthzMode, o.Modes, ""+
		"Ordered list of plug-ins to do authorization on secure port. Comma-delimited list of: "+
		strings.Join(modes.AuthorizationModeChoices, ",")+".")
	_ = viper.BindPFlag(configAuthzMode, fs.Lookup(flagAuthzMode))

	fs.String(flagAuthzWebhookConfigFile, o.WebhookConfigFile, ""+
		"File with webhook configuration in kubeconfig format, used with --authorization-mode=Webhook. "+
		"The API server will query the remote service to determine access on the API server's secure port.")
	_ = viper.BindPFlag(configAuthzWebhookConfigFile, fs.Lookup(flagAuthzWebhookConfigFile))

	fs.Duration(flagAuthzWebhookCacheAuthorizedTTL, o.WebhookCacheAuthorizedTTL, ""+
		"The duration to cache 'authorized' responses from the webhook authorizer.")
	_ = viper.BindPFlag(configAuthzWebhookCacheAuthorizedTTL, fs.Lookup(flagAuthzWebhookCacheAuthorizedTTL))

	fs.Duration(flagAuthzWebhookCacheUnauthorizedTTL, o.WebhookCacheUnauthorizedTTL,
		"The duration to cache 'unauthorized' responses from the webhook authorizer.")
	_ = viper.BindPFlag(configAuthzWebhookCacheUnauthorizedTTL, fs.Lookup(flagAuthzWebhookCacheUnauthorizedTTL))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *AuthorizationOptions) ApplyFlags() []error {
	var errs []error

	o.Modes = viper.GetStringSlice(configAuthzMode)
	o.WebhookCacheAuthorizedTTL = viper.GetDuration(configAuthzWebhookCacheAuthorizedTTL)
	o.WebhookCacheUnauthorizedTTL = viper.GetDuration(configAuthzWebhookCacheUnauthorizedTTL)
	o.WebhookConfigFile = viper.GetString(configAuthzWebhookConfigFile)

	if len(o.Modes) == 0 {
		errs = append(errs, fmt.Errorf("at least one authorization-mode must be passed"))
	}

	allowedModes := sets.NewString(modes.AuthorizationModeChoices...)
	ms := sets.NewString(o.Modes...)
	for _, mode := range o.Modes {
		if !allowedModes.Has(mode) {
			errs = append(errs, fmt.Errorf("authorization-mode %q is not a valid mode", mode))
		}
		if mode == modes.ModeWebhook {
			if o.WebhookConfigFile == "" {
				errs = append(errs, fmt.Errorf("authorization-mode Webhook's authorization config file not passed"))
			}
		}
	}

	if o.WebhookConfigFile != "" && !ms.Has(modes.ModeWebhook) {
		errs = append(errs, fmt.Errorf("cannot specify --authorization-webhook-config-file without mode Webhook"))
	}

	if len(o.Modes) != len(ms.List()) {
		errs = append(errs, fmt.Errorf("authorization-mode %q has mode specified more than once", o.Modes))
	}

	return errs
}
