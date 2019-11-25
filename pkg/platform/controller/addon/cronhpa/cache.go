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

package cronhpa

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedCronHPA struct {
	// The cached state of the CronHPA
	state *v1.CronHPA
}

type cronHPACache struct {
	mu         sync.Mutex // protects cronHPAMap
	cronHPAMap map[string]*cachedCronHPA
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *cronHPACache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.cronHPAMap))
	for k := range s.cronHPAMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the cronHPAMap under the given key
func (s *cronHPACache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.cronHPAMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *cronHPACache) get(cronHPAName string) (*cachedCronHPA, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cronHPA, ok := s.cronHPAMap[cronHPAName]
	return cronHPA, ok
}

func (s *cronHPACache) Exist(cronHPAName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.cronHPAMap[cronHPAName]
	return ok
}

func (s *cronHPACache) getOrCreate(cronHPAName string) *cachedCronHPA {
	s.mu.Lock()
	defer s.mu.Unlock()
	cronHPA, ok := s.cronHPAMap[cronHPAName]
	if !ok {
		cronHPA = &cachedCronHPA{}
		s.cronHPAMap[cronHPAName] = cronHPA
	}
	return cronHPA
}

func (s *cronHPACache) set(cronHPAName string, cronHPA *cachedCronHPA) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cronHPAMap[cronHPAName] = cronHPA
}

func (s *cronHPACache) delete(cronHPAName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.cronHPAMap, cronHPAName)
}
