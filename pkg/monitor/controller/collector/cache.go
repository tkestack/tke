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

package collector

import (
	"sync"

	v1 "tkestack.io/tke/api/monitor/v1"
)

type cachedCollector struct {
	// The cached state of the collector
	state *v1.Collector
}

type collectorCache struct {
	mu           sync.Mutex // protects collectorMap
	collectorMap map[string]*cachedCollector
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *collectorCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.collectorMap))
	for k := range s.collectorMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the collectorMap under the given key
func (s *collectorCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.collectorMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *collectorCache) get(collectorName string) (*cachedCollector, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	prometheus, ok := s.collectorMap[collectorName]
	return prometheus, ok
}

func (s *collectorCache) Exist(collectorName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.collectorMap[collectorName]
	return ok
}

func (s *collectorCache) getOrCreate(collectorName string) *cachedCollector {
	s.mu.Lock()
	defer s.mu.Unlock()
	prometheus, ok := s.collectorMap[collectorName]
	if !ok {
		prometheus = &cachedCollector{}
		s.collectorMap[collectorName] = prometheus
	}
	return prometheus
}

func (s *collectorCache) set(collectorName string, prometheus *cachedCollector) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.collectorMap[collectorName] = prometheus
}

func (s *collectorCache) delete(collectorName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.collectorMap, collectorName)
}
