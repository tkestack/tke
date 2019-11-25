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

package helm

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedHelm struct {
	// The cached state of the helm
	state *v1.Helm
}

type helmCache struct {
	mu      sync.Mutex // protects helmMap
	helmMap map[string]*cachedHelm
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *helmCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.helmMap))
	for k := range s.helmMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the helmMap under the given key
func (s *helmCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.helmMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *helmCache) get(helmName string) (*cachedHelm, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	helm, ok := s.helmMap[helmName]
	return helm, ok
}

func (s *helmCache) Exist(helmName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.helmMap[helmName]
	return ok
}

func (s *helmCache) getOrCreate(helmName string) *cachedHelm {
	s.mu.Lock()
	defer s.mu.Unlock()
	helm, ok := s.helmMap[helmName]
	if !ok {
		helm = &cachedHelm{}
		s.helmMap[helmName] = helm
	}
	return helm
}

func (s *helmCache) set(helmName string, helm *cachedHelm) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.helmMap[helmName] = helm
}

func (s *helmCache) delete(helmName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.helmMap, helmName)
}
