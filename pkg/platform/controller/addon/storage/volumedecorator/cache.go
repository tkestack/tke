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

package volumedecorator

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedVolumeDecorator struct {
	// The cached state of the decorator.
	state *v1.VolumeDecorator
}

type volumeDecoratorCache struct {
	mu    sync.Mutex // protects vdMap
	vdMap map[string]*cachedVolumeDecorator
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *volumeDecoratorCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.vdMap))
	for k := range s.vdMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the vdMap under the given key
func (s *volumeDecoratorCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.vdMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *volumeDecoratorCache) get(decoratorName string) (*cachedVolumeDecorator, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	decorator, ok := s.vdMap[decoratorName]
	return decorator, ok
}

func (s *volumeDecoratorCache) Exist(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.vdMap[name]
	return ok
}

func (s *volumeDecoratorCache) getOrCreate(name string) *cachedVolumeDecorator {
	s.mu.Lock()
	defer s.mu.Unlock()
	decorator, ok := s.vdMap[name]
	if !ok {
		decorator = &cachedVolumeDecorator{}
		s.vdMap[name] = decorator
	}
	return decorator
}

func (s *volumeDecoratorCache) set(name string, decorator *cachedVolumeDecorator) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vdMap[name] = decorator
}

func (s *volumeDecoratorCache) delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.vdMap, name)
}
