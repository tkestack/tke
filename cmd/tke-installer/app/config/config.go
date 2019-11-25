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

package config

import (
	"tkestack.io/tke/cmd/tke-installer/app/options"
)

// Config is the running configuration structure of the TKE controller manager.
type Config struct {
	ServerName                 string
	ListenAddr                 string
	NoUI                       bool
	Config                     string
	Force                      bool
	SyncProjectsWithNamespaces bool
	Replicas                   int
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given TKE apiserver command line or configuration file option.
func CreateConfigFromOptions(serverName string, opts *options.Options) (*Config, error) {
	return &Config{
		ServerName:                 serverName,
		ListenAddr:                 *opts.ListenAddr,
		NoUI:                       *opts.NoUI,
		Config:                     *opts.Config,
		Force:                      *opts.Force,
		SyncProjectsWithNamespaces: *opts.SyncProjectsWithNamespaces,
		Replicas:                   *opts.Replicas,
	}, nil
}
