/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package helm

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"gopkg.in/yaml.v2"
)

type Values struct {
	Audit    Audit    `yaml:"audit"`
	Auth     Auth     `yaml:"auth"`
	Business Business `yaml:"business"`
	Gateway  Gateway  `yaml:"gateway"`
	Logagent Logagent `yaml:"logagent"`
	Monitor  Monitor  `yaml:"monitor"`
	Registry Registry `yaml:"registry"`

	CertsName        string     `yaml:"certsName"`
	CACert           string     `yaml:"caCert"`
	CAKey            string     `yaml:"caKey"`
	ServerCert       string     `yaml:"serverCert"`
	ServerKey        string     `yaml:"serverKey"`
	FrontProxyCACert string     `yaml:"frontProxyCACert"`
	OIDCCACert       string     `yaml:"oidcCACert"`
	ETCDValues       ETCDValues `yaml:"etcd"`

	OIDCClientClientID string `yaml:"oidcClientID"`
	OIDCClientSecret   string `yaml:"oidcClientSecret"`
	OIDCIssuerURL      string `yaml:"oidcIssuerURL"`
	USEOIDCCA          bool   `yaml:"useOIDCCA"`

	RegistryPreifx string `yaml:"registryPrefix"`
	Token          string `yaml:"token"`
	TenantID       string `yaml:"tenantID"`
	AdminUsername  string `yaml:"adminUsername"`
	AdminPassword  string `yaml:"adminPassword"`
}

type ETCDValues struct {
	Hosts      []string `yaml:"hosts"`
	CACert     string   `yaml:"caCert"`
	ClientCert string   `yaml:"clientCert"`
	ClientKey  string   `yaml:"clientKey"`
}

type Audit struct {
	Eanbled            bool               `yaml:"enabled"`
	AuditElasticSearch AuditElasticSearch `yaml:"elasticSearch"`
}

type AuditElasticSearch struct {
	Address     string `yaml:"address"`
	ReserveDays int    `yaml:"reserveDays"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
}

type Auth struct {
	Eanbled       bool     `yaml:"enabled"`
	RedirectHosts []string `yaml:"redirectHosts"`
}

type Business struct {
	Eanbled bool `yaml:"enabled"`
}

type Gateway struct {
	Eanbled bool `yaml:"enabled"`
}

type Logagent struct {
	Eanbled bool `yaml:"enabled"`
}

type Monitor struct {
	Eanbled         bool     `yaml:"enabled"`
	InfluxDB        InfluxDB `yaml:"influxdb"`
	StorageType     string   `yaml:"storageType"`
	StorageAddress  string   `yaml:"storageAddress"`
	StorageUsername string   `yaml:"storageUsername"`
	StoragePassword string   `yaml:"storagePassword"`
}

type InfluxDB struct {
	Eanbled  bool   `yaml:"enabled"`
	NodeName string `yaml:"nodeName"`
}

type Registry struct {
	Eanbled   bool   `yaml:"enabled"`
	Domain    string `yaml:"domain"`
	Namespace string `yaml:"namespace"`
	NodeName  string `yaml:"nodeName"`
}

const (
	ValuesFile = "helmfile.d/environments/default/values.yaml"
)

func Install(namespace string, values *Values) error {
	err := save(values)
	if err != nil {
		return err
	}

	err = sync(namespace)
	if err != nil {
		return err
	}

	return nil
}

func save(values *Values) error {
	err := saveValues(ValuesFile, values)
	if err != nil {
		return err
	}

	return nil
}

func saveValues(filename string, values interface{}) error {
	_ = os.MkdirAll(path.Dir(filename), 0755)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	y := yaml.NewEncoder(f)
	return y.Encode(values)
}

func sync(namespace string) error {
	cmdString := fmt.Sprintf("helmfile --namespace=%s sync", namespace)
	cmd := exec.Command("sh", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
