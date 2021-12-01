/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package application

import (
	applicationprovider "tkestack.io/tke/pkg/application/provider/application"
	"tkestack.io/tke/pkg/util/log"
)

const (
	name = "Example"
)

// Run your register func in your main.
func RegisterProvider() {
	p, err := NewProvider()
	if err != nil {
		log.Errorf("init application provider error: %s", err)
		return
	}
	applicationprovider.Register(p.Name(), p)
}

type Provider struct {
	*applicationprovider.DelegateProvider
	yourconfig string
}

var _ applicationprovider.Provider = &Provider{}

func NewProvider() (*Provider, error) {
	result := new(Provider)
	result.yourconfig = "yourconfig"
	result.DelegateProvider.ProviderName = name
	return result, nil
}
