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
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/jinzhu/configor"
)

func New(filename string) (*Config, error) {
	config := &Config{}
	if err := configor.Load(config, filename); err != nil {
		return nil, err
	}

	s := strings.Split(config.Registry.Prefix, "/")
	if len(s) != 2 {
		return nil, errors.New("invalid registry prefix")
	}
	config.Registry.Domain = s[0]
	config.Registry.Namespace = s[1]

	return config, nil
}

type Config struct {
	Registry               Registry               `yaml:"registry"`
	Docker                 Docker                 `yaml:"docker" required:"true"`
	CNIPlugins             CNIPlugins             `yaml:"cniPlugins" required:"true"`
	NvidiaDriver           NvidiaDriver           `yaml:"nvidiaDriver" required:"true"`
	NvidiaContainerRuntime NvidiaContainerRuntime `yaml:"nvidiaContainerRuntime" required:"true"`
}

func (c *Config) Save(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	y := yaml.NewEncoder(f)
	return y.Encode(c)
}

type Registry struct {
	Prefix    string `yaml:"prefix"`
	IP        string `yaml:"ip"`
	Domain    string
	Namespace string
}

func (r *Registry) UseTKE() bool {
	return r.Domain != "docker.io"
}

type Docker struct {
	DefaultVersion string `yaml:"defaultVersion"`
}

type CNIPlugins struct {
	DefaultVersion string `yaml:"defaultVersion"`
}

type NvidiaDriver struct {
	DefaultVersion string `yaml:"defaultVersion"`
}

type NvidiaContainerRuntime struct {
	DefaultVersion string `yaml:"defaultVersion"`
}
