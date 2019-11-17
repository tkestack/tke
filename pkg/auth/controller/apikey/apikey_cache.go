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

package apikey


import (
	"sync"

	v1 "tkestack.io/tke/api/auth/v1"
)

type cachedAPIKey struct {
	// The cached state of the cluster
	state *v1.APIKey
}

type apiKeyCache struct {
	mu        sync.Mutex // protects apiKeyMap
	apiKeyMap map[string]*cachedAPIKey
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *apiKeyCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.apiKeyMap))
	for k := range s.apiKeyMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the apiKeyMap under the given key
func (s *apiKeyCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.apiKeyMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *apiKeyCache) get(keyName string) (*cachedAPIKey, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cluster, ok := s.apiKeyMap[keyName]
	return cluster, ok
}

func (s *apiKeyCache) Exist(keyName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.apiKeyMap[keyName]
	return ok
}

func (s *apiKeyCache) getOrCreate(keyName string) *cachedAPIKey {
	s.mu.Lock()
	defer s.mu.Unlock()
	cluster, ok := s.apiKeyMap[keyName]
	if !ok {
		cluster = &cachedAPIKey{}
		s.apiKeyMap[keyName] = cluster
	}
	return cluster
}

func (s *apiKeyCache) set(keyName string, cluster *cachedAPIKey) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apiKeyMap[keyName] = cluster
}

func (s *apiKeyCache) delete(keyName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.apiKeyMap, keyName)
}

