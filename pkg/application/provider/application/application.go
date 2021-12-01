/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package application

import (
	"fmt"
	"sort"
	"sync"

	applicationv1 "tkestack.io/tke/api/application/v1"
)

var (
	providersMu sync.RWMutex
	providers   = defaultProviders()
)

const AnnotationProviderNameKey = "application.tkestack.io/provider-name"

// Register makes a provider available by the provided name.
// If Register is called twice with the same name or if provider is nil,
// it panics.
func Register(name string, provider Provider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("application: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("application: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// Providers returns a sorted list of the names of the registered providers.
func Providers() []string {
	providersMu.RLock()
	defer providersMu.RUnlock()
	var list []string
	for name := range providers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// GetProvider will get your provider with the application,
// set an annotation with key, application.tkestack.io/provider-name, and value, the provider will work for your application.
func GetProvider(app *applicationv1.App) (Provider, error) {
	if app == nil {
		return &DelegateProvider{}, nil
	}
	providersMu.RLock()
	provider, ok := providers[app.Annotations[AnnotationProviderNameKey]]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("application: unknown provider %q (forgotten import?)", app.Annotations[app.Annotations[AnnotationProviderNameKey]])

	}

	return provider, nil
}

func defaultProviders() map[string]Provider {
	results := make(map[string]Provider)
	results[""] = &DelegateProvider{}
	return results
}
