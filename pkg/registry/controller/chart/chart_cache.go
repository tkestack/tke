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

package chart

import (
	"sync"

	registryv1 "tkestack.io/tke/api/registry/v1"
)

type cachedChart struct {
	// The cached state of the value
	state *registryv1.Chart
}

type chartCache struct {
	mu sync.Mutex
	m  map[string]*cachedChart
}

func (s *chartCache) getOrCreate(name string) *cachedChart {
	s.mu.Lock()
	defer s.mu.Unlock()
	chart, ok := s.m[name]
	if !ok {
		chart = &cachedChart{}
		s.m[name] = chart
	}
	return chart
}

func (s *chartCache) get(name string) (*cachedChart, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	chart, ok := s.m[name]
	return chart, ok
}

func (s *chartCache) Exist(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[name]
	return ok
}

func (s *chartCache) set(name string, chart *cachedChart) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[name] = chart
}

func (s *chartCache) delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, name)
}
