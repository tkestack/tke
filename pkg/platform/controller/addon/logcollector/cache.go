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

package logcollector

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedLogCollector struct {
	// The cached state of the collector
	state *v1.LogCollector
}

type logcollectorCache struct {
	mu    sync.Mutex // protects lcMap
	lcMap map[string]*cachedLogCollector
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *logcollectorCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.lcMap))
	for k := range s.lcMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the lcMap under the given key
func (s *logcollectorCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.lcMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *logcollectorCache) get(logCollectorName string) (*cachedLogCollector, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	logCollector, ok := s.lcMap[logCollectorName]
	return logCollector, ok
}

func (s *logcollectorCache) Exist(logCollectorName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.lcMap[logCollectorName]
	return ok
}

func (s *logcollectorCache) getOrCreate(logCollectorName string) *cachedLogCollector {
	s.mu.Lock()
	defer s.mu.Unlock()
	collector, ok := s.lcMap[logCollectorName]
	if !ok {
		collector = &cachedLogCollector{}
		s.lcMap[logCollectorName] = collector
	}
	return collector
}

func (s *logcollectorCache) set(logCollectorName string, collector *cachedLogCollector) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lcMap[logCollectorName] = collector
}

func (s *logcollectorCache) delete(logCollectorName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.lcMap, logCollectorName)
}
