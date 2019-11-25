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

package lbcf

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedLBCF struct {
	// The cached state of the lbcf
	state *v1.LBCF
}

type lbcfCache struct {
	mu      sync.Mutex // protects lbcfMap
	lbcfMap map[string]*cachedLBCF
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *lbcfCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.lbcfMap))
	for k := range s.lbcfMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the lbcfMap under the given key
func (s *lbcfCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.lbcfMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *lbcfCache) get(lbcfName string) (*cachedLBCF, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	lbcf, ok := s.lbcfMap[lbcfName]
	return lbcf, ok
}

func (s *lbcfCache) Exist(lbcfName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.lbcfMap[lbcfName]
	return ok
}

func (s *lbcfCache) getOrCreate(lbcfName string) *cachedLBCF {
	s.mu.Lock()
	defer s.mu.Unlock()
	lbcf, ok := s.lbcfMap[lbcfName]
	if !ok {
		lbcf = &cachedLBCF{}
		s.lbcfMap[lbcfName] = lbcf
	}
	return lbcf
}

func (s *lbcfCache) set(lbcfName string, lbcf *cachedLBCF) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lbcfMap[lbcfName] = lbcf
}

func (s *lbcfCache) delete(lbcfName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.lbcfMap, lbcfName)
}
