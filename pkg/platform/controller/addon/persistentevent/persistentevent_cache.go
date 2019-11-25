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

package persistentevent

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedPersistentEvent struct {
	// The cached state of the persistent event
	state *v1.PersistentEvent
}

type persistentEventCache struct {
	mu                 sync.Mutex // protects persistentEventMap
	persistentEventMap map[string]*cachedPersistentEvent
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *persistentEventCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.persistentEventMap))
	for k := range s.persistentEventMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the persistentEventMap under the given key
func (s *persistentEventCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.persistentEventMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *persistentEventCache) get(key string) (*cachedPersistentEvent, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	persistentEvent, ok := s.persistentEventMap[key]
	return persistentEvent, ok
}

func (s *persistentEventCache) Exist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.persistentEventMap[key]
	return ok
}

func (s *persistentEventCache) getOrCreate(key string) *cachedPersistentEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	persistentEvent, ok := s.persistentEventMap[key]
	if !ok {
		persistentEvent = &cachedPersistentEvent{}
		s.persistentEventMap[key] = persistentEvent
	}
	return persistentEvent
}

func (s *persistentEventCache) set(key string, persistentEvent *cachedPersistentEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.persistentEventMap[key] = persistentEvent
}

func (s *persistentEventCache) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.persistentEventMap, key)
}
