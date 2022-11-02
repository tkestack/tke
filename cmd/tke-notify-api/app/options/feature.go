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
	"time"
)

const (
	flagMessageTTL    = "message-ttl"
	configMessageTTL  = "features.message_ttl"
	defaultMessageTTL = time.Hour * 24 * 30
)

type FeatureOptions struct {
	MessageTTL time.Duration
}

func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{MessageTTL: defaultMessageTTL}
}

func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&o.MessageTTL, flagMessageTTL, o.MessageTTL,
		"How long to retain messages and messagerequests")
	_ = viper.BindPFlag(configMessageTTL, fs.Lookup(flagMessageTTL))
}

func (o *FeatureOptions) ApplyFlags() []error {
	var errs []error

	o.MessageTTL = viper.GetDuration(configMessageTTL)

	return errs
}
