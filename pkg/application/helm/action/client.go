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

package action

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/registry"
	"tkestack.io/tke/pkg/application/helm/config"
	"tkestack.io/tke/pkg/util/log"
)

// Client is a client used to manage helm release
type Client struct {
	helmDriver       string
	restClientGetter *config.RESTClientGetter
}

// NewClient return a client that will work on helm release
func NewClient(helmDriver string, restClientGetter *config.RESTClientGetter) *Client {
	return &Client{
		helmDriver:       helmDriver,
		restClientGetter: restClientGetter,
	}
}

func (c *Client) buildActionConfig(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(c.restClientGetter, namespace, c.helmDriver, log.Debugf)
	if err != nil {
		return nil, err
	}
	actionConfig.RegistryClient, err = registry.NewClient()
	if err != nil {
		return nil, err
	}
	return actionConfig, nil
}
