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

package machine

import (
	"errors"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"k8s.io/apimachinery/pkg/api/validation"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/log/hclog"
)

// Client is a wraper of machine provider instance
type Client struct {
	Name   string
	Client *plugin.Client
	Provider
}

// NewClient creates machine provider client
func NewClient(providerPluginFile string, providerConfigFile string) (*Client, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMaps,
		Cmd:             exec.Command(providerPluginFile),
		Logger:          hclog.NewHCLogger(log.ZapLogger()),
	})

	needCloseClient := false
	defer func() {
		if needCloseClient {
			client.Kill()
		}
	}()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		needCloseClient = true
		log.Error("Failed to create the machine provider plugin client", log.Err(err))
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		needCloseClient = true
		log.Error("Failed to dispense the machine provider plugin by given name", log.Err(err))
		return nil, err
	}

	machineProvider, ok := raw.(Provider)
	if !ok {
		needCloseClient = true
		log.Error("Dispensed machine provider plugin cannot cast to provider interface")
		return nil, errors.New("dispensed machine provider plugin cannot cast to provider interface")
	}

	err = machineProvider.Init(providerConfigFile)
	if err != nil {
		needCloseClient = true
		log.Error("Failed to call machineProvider.Init", log.Err(err))
		return nil, err
	}

	name, err := machineProvider.Name()
	if err != nil {
		needCloseClient = true
		log.Error("Failed to call machineProvider.Name", log.Err(err))
		return nil, err
	}
	if errs := validation.NameIsDNSLabel(name, false); len(errs) > 0 {
		needCloseClient = true
		return nil, errors.New("the provider name of plugin does not conform to the DNS name specification")
	}

	return &Client{
		Name:     name,
		Client:   client,
		Provider: machineProvider,
	}, nil
}

func (c *Client) Close() {
	if c.Client != nil {
		c.Client.Kill()
	}
}
