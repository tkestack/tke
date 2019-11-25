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

package ipam

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedIPAM struct {
	// The cached state of the ipam
	state *v1.IPAM
}

type ipamCache struct {
	mu      sync.Mutex // protects ipamMap
	ipamMap map[string]*cachedIPAM
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *ipamCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.ipamMap))
	for k := range s.ipamMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the ipamMap under the given key
func (s *ipamCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.ipamMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *ipamCache) get(ipamName string) (*cachedIPAM, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ipam, ok := s.ipamMap[ipamName]
	return ipam, ok
}

func (s *ipamCache) Exist(ipamName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.ipamMap[ipamName]
	return ok
}

func (s *ipamCache) getOrCreate(ipamName string) *cachedIPAM {
	s.mu.Lock()
	defer s.mu.Unlock()
	ipam, ok := s.ipamMap[ipamName]
	if !ok {
		ipam = &cachedIPAM{}
		s.ipamMap[ipamName] = ipam
	}
	return ipam
}

func (s *ipamCache) set(ipamName string, ipam *cachedIPAM) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ipamMap[ipamName] = ipam
}

func (s *ipamCache) delete(ipamName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ipamMap, ipamName)
}
