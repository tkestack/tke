/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

package cluster

import (
	"fmt"
	"sort"
	"sync"

	"k8s.io/apiserver/pkg/server/mux"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]Provider)
)

// Register makes a provider available by the provided name.
// If Register is called twice with the same name or if provider is nil,
// it panics.
func Register(name string, provider Provider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("cluster: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("cluster: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// RegisterHandler register all provider's hanlder.
func RegisterHandler(mux *mux.PathRecorderMux) {
	for _, p := range providers {
		p.RegisterHandler(mux)
	}
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

// GetProvider returns provider by name
func GetProvider(name string) (Provider, error) {
	providersMu.RLock()
	provider, ok := providers[name]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cluster: unknown provider %q (forgotten import?)", name)

	}

	return provider, nil
}
