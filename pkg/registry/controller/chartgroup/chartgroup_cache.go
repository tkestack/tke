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

package chartgroup

import (
	"sync"

	registryv1 "tkestack.io/tke/api/registry/v1"
)

type cachedChartGroup struct {
	// The cached state of the value
	state *registryv1.ChartGroup
}

type chartGroupCache struct {
	mu sync.Mutex
	m  map[string]*cachedChartGroup
}

func (s *chartGroupCache) getOrCreate(name string) *cachedChartGroup {
	s.mu.Lock()
	defer s.mu.Unlock()
	chartGroup, ok := s.m[name]
	if !ok {
		chartGroup = &cachedChartGroup{}
		s.m[name] = chartGroup
	}
	return chartGroup
}

func (s *chartGroupCache) get(name string) (*cachedChartGroup, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	chartGroup, ok := s.m[name]
	return chartGroup, ok
}

func (s *chartGroupCache) Exist(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[name]
	return ok
}

func (s *chartGroupCache) set(name string, chartGroup *cachedChartGroup) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[name] = chartGroup
}

func (s *chartGroupCache) delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, name)
}

func (s *chartGroupCache) listKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

func (s *chartGroupCache) allChartGroups() []*registryv1.ChartGroup {
	s.mu.Lock()
	defer s.mu.Unlock()
	chartGroups := make([]*registryv1.ChartGroup, 0, len(s.m))
	for _, v := range s.m {
		chartGroups = append(chartGroups, v.state)
	}
	return chartGroups
}
