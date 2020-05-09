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
	flagDomain           = "domain"
	flagNamespace        = "namespace"
)

const (
	configDomain           = "features.domain"
	configNamespace        = "features.namespace"
)

type FeatureOptions struct {
	Domain           string
	Namespace        string
}

func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{

	}
}

// AddFlags adds flags for console to the specified FlagSet object.
func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.String(flagDomain, o.Domain,"registry domain")
	_ = viper.BindPFlag(configDomain,fs.Lookup(flagDomain))

	fs.String(flagNamespace, o.Namespace, "registry namespace")
	_ = viper.BindPFlag(configNamespace, fs.Lookup(flagNamespace))

}


// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *FeatureOptions) ApplyFlags() []error {
	var errs []error
	o.Domain = viper.GetString(configDomain)
	o.Namespace = viper.GetString(configNamespace)

	if len(o.Namespace) == 0 {
		errs = append(errs, fmt.Errorf("%s must be specified", configNamespace))
	}
	return errs
}
