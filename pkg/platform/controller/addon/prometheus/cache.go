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

package prometheus

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedPrometheus struct {
	// The cached state of the prometheus
	state *v1.Prometheus
}

type prometheusCache struct {
	mu            sync.Mutex // protects prometheusMap
	prometheusMap map[string]*cachedPrometheus
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *prometheusCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.prometheusMap))
	for k := range s.prometheusMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the prometheusMap under the given key
func (s *prometheusCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.prometheusMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *prometheusCache) get(prometheusName string) (*cachedPrometheus, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	prometheus, ok := s.prometheusMap[prometheusName]
	return prometheus, ok
}

func (s *prometheusCache) Exist(prometheusName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.prometheusMap[prometheusName]
	return ok
}

func (s *prometheusCache) getOrCreate(prometheusName string) *cachedPrometheus {
	s.mu.Lock()
	defer s.mu.Unlock()
	prometheus, ok := s.prometheusMap[prometheusName]
	if !ok {
		prometheus = &cachedPrometheus{}
		s.prometheusMap[prometheusName] = prometheus
	}
	return prometheus
}

func (s *prometheusCache) set(prometheusName string, prometheus *cachedPrometheus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prometheusMap[prometheusName] = prometheus
}

func (s *prometheusCache) delete(prometheusName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.prometheusMap, prometheusName)
}
