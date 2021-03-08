/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package meshmanager

import (
	"sync"
	v1 "tkestack.io/tke/api/mesh/v1"
)

type cachedMeshManager struct {
	// The cached state of the meshManager
	state *v1.MeshManager
}

type meshManagerCache struct {
	mu    sync.Mutex // protects lcMap
	lcMap map[string]*cachedMeshManager
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *meshManagerCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.lcMap))
	for k := range s.lcMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the lcMap under the given key
func (s *meshManagerCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.lcMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *meshManagerCache) get(meshManagerName string) (*cachedMeshManager, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	meshManager, ok := s.lcMap[meshManagerName]
	return meshManager, ok
}

func (s *meshManagerCache) Exist(meshManagerName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.lcMap[meshManagerName]
	return ok
}

func (s *meshManagerCache) getOrCreate(meshManagerName string) *cachedMeshManager {
	s.mu.Lock()
	defer s.mu.Unlock()
	collector, ok := s.lcMap[meshManagerName]
	if !ok {
		collector = &cachedMeshManager{}
		s.lcMap[meshManagerName] = collector
	}
	return collector
}

func (s *meshManagerCache) set(meshManagerName string, meshManager *cachedMeshManager) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lcMap[meshManagerName] = meshManager
}

func (s *meshManagerCache) delete(meshManagerName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.lcMap, meshManagerName)
}
