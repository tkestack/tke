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

package gpumanager

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type gmCachedItem struct {
	// The cached holder of the gpu manager
	holder *v1.GPUManager
}

type gmCache struct {
	mu    sync.Mutex // protects store
	store map[string]*gmCachedItem
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *gmCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the store under the given key
func (s *gmCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.store[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *gmCache) get(key string) (*gmCachedItem, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.store[key]
	return val, ok
}

func (s *gmCache) Exist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.store[key]
	return ok
}

func (s *gmCache) getOrCreate(key string) *gmCachedItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.store[key]
	if !ok {
		val = &gmCachedItem{}
		s.store[key] = val
	}
	return val
}

func (s *gmCache) set(key string, val *gmCachedItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = val
}

func (s *gmCache) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
}
