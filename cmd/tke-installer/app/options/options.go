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
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
)

// Options is the main context object for the TKE apiserver.
type Options struct {
	Log                        *log.Options
	ListenAddr                 *string
	NoUI                       *bool
	Config                     *string
	Force                      *bool
	SyncProjectsWithNamespaces *bool
	Replicas                   *int
}

// NewOptions creates a new Options with a default config.
func NewOptions(serverName string) *Options {
	return &Options{
		Log: log.NewOptions(),
	}
}

// AddFlags adds flags for a specific server to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.Log.AddFlags(fs)

	ip, err := util.GetExternalIP()
	if err != nil {
		panic(err)
	}
	o.ListenAddr = fs.String("listen-addr", fmt.Sprintf("%s:8080", ip), "listen addr")
	o.NoUI = fs.Bool("no-ui", false, "run without web")
	o.Config = fs.String("input", "conf/tke.json", "specify input file")
	o.Force = fs.Bool("force", false, "force run as clean")
	o.SyncProjectsWithNamespaces = fs.Bool("sync-projects-with-namespaces", false, "Enable creating/deleting the corresponding namespace when creating/deleting a project.")
	o.Replicas = fs.Int("replicas", 2, "tke components replicas")
}

// ApplyFlags parsing parameters from the command line or configuration file
// to the options instance.
func (o *Options) ApplyFlags() []error {
	var errs []error

	errs = append(errs, o.Log.ApplyFlags()...)

	return errs
}

// Complete set default Options.
// Should be called after tke-installer flags parsed.
func (o *Options) Complete() error {
	return nil
}
