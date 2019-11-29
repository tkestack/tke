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

package policy


import (
	"sync"

	v1 "tkestack.io/tke/api/auth/v1"
)

type cachedPolicy struct {
	// The cached state of the cluster
	state *v1.Policy
}

type policyCache struct {
	mu        sync.Mutex // protects policyMap
	policyMap map[string]*cachedPolicy
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *policyCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.policyMap))
	for k := range s.policyMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the policyMap under the given key
func (s *policyCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.policyMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *policyCache) get(keyName string) (*cachedPolicy, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cluster, ok := s.policyMap[keyName]
	return cluster, ok
}

func (s *policyCache) Exist(keyName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.policyMap[keyName]
	return ok
}

func (s *policyCache) getOrCreate(keyName string) *cachedPolicy {
	s.mu.Lock()
	defer s.mu.Unlock()
	cluster, ok := s.policyMap[keyName]
	if !ok {
		cluster = &cachedPolicy{}
		s.policyMap[keyName] = cluster
	}
	return cluster
}

func (s *policyCache) set(keyName string, cluster *cachedPolicy) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.policyMap[keyName] = cluster
}

func (s *policyCache) delete(keyName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.policyMap, keyName)
}

