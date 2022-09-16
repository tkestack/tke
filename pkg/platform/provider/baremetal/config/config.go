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
	PlatformAPIClientConfig    string            `yaml:"platformAPIClientConfig"`
	ApplicationAPIClientConfig string            `yaml:"applicationAPIClientConfig"`
	Registry                   Registry          `yaml:"registry"`
	Audit                      Audit             `yaml:"audit"`
	Feature                    Feature           `yaml:"feature"`
	Docker                     Docker            `yaml:"docker"`
	Kubelet                    Kubelet           `yaml:"kubelet"`
	APIServer                  APIServer         `yaml:"apiServer"`
	ControllerManager          ControllerManager `yaml:"controllerManager"`
	Scheduler                  Scheduler         `yaml:"scheduler"`
	AuthzWebhook               AuthzWebhook      `yaml:"authzWebhook"`
	Business                   Business          `yaml:"business"`
	SupportOSList              []string          `yaml:"supportOSList"`
}

func (c *Config) Save(filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	y := yaml.NewEncoder(f)
	return y.Encode(c)
}

func (c *Config) AuditEnabled() bool {
	return c.Audit.Address != ""
}

func (c *Config) BusinessEnabled() bool {
	return c.Business.Enabled
}

type Registry struct {
	Prefix    string `yaml:"prefix"`
	IP        string `yaml:"ip"`
	Domain    string `yaml:"-"`
	Namespace string `yaml:"-"`
}

type Audit struct {
	Address string `yaml:"address"`
}

type Feature struct {
	SkipConditions []string `yaml:"skipConditions"`
}

type Docker struct {
	ExtraArgs map[string]string `yaml:"extraArgs"`
}

type Kubelet struct {
	ExtraArgs map[string]string `yaml:"extraArgs"`
}

type APIServer struct {
	ExtraArgs map[string]string `yaml:"extraArgs"`
}

type ControllerManager struct {
	ExtraArgs map[string]string `yaml:"extraArgs"`
}

type Scheduler struct {
	ExtraArgs map[string]string `yaml:"extraArgs"`
}

type AuthzWebhook struct {
	Endpoint string `yaml:"endpoint"`
}
type Business struct {
	Enabled bool `yaml:"enabled"`
}
