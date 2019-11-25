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
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagContainerRegistryDomain    = "container-registry-domain"
	flagContainerRegistryNamespace = "container-registry-namespace"
)

const (
	configContainerRegistryDomain    = "registry.container_domain"
	configContainerRegistryNamespace = "registry.container_namespace"
)

// RegistryOptions holds the container registry options.
type RegistryOptions struct {
	Domain    string
	Namespace string
}

// NewRegistryOptions creates the default RegistryOptions object.
func NewRegistryOptions() *RegistryOptions {
	return &RegistryOptions{
		Namespace: "tkestack",
	}
}

// AddFlags adds flags related to debugging for controller manager to the specified FlagSet.
func (o *RegistryOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.String(flagContainerRegistryDomain, o.Domain,
		"Sets domain name of the container registry.")
	_ = viper.BindPFlag(configContainerRegistryDomain, fs.Lookup(flagContainerRegistryDomain))

	fs.String(flagContainerRegistryNamespace, o.Namespace,
		"Sets namespace of the container registry.")
	_ = viper.BindPFlag(configContainerRegistryNamespace, fs.Lookup(flagContainerRegistryNamespace))
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *RegistryOptions) ApplyFlags() []error {
	var errs []error

	o.Domain = viper.GetString(configContainerRegistryDomain)
	o.Namespace = viper.GetString(configContainerRegistryNamespace)

	errs = append(errs, o.Validate()...)

	return errs
}

// Validate checks validation of RegistryOptions.
func (o *RegistryOptions) Validate() []error {
	var errs []error

	if o.Namespace == "" {
		errs = append(errs, fmt.Errorf("--%s must be specified", flagContainerRegistryNamespace))
	} else if strings.Contains(o.Namespace, "/") {
		errs = append(errs, fmt.Errorf("container registry namespace is not allowed to contain '/' symbols"))
	}

	return errs
}
