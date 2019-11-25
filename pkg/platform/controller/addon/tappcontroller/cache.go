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

package tappcontroller

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedTappController struct {
	// The cached state of the tapp controller
	state *v1.TappController
}

type tappControllerCache struct {
	mu                sync.Mutex // protects tappControllerMap
	tappControllerMap map[string]*cachedTappController
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *tappControllerCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.tappControllerMap))
	for k := range s.tappControllerMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the tappControllerMap under the given key
func (s *tappControllerCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.tappControllerMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *tappControllerCache) get(tappControllerName string) (*cachedTappController, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tappController, ok := s.tappControllerMap[tappControllerName]
	return tappController, ok
}

func (s *tappControllerCache) Exist(tappControllerName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tappControllerMap[tappControllerName]
	return ok
}

func (s *tappControllerCache) getOrCreate(tappControllerName string) *cachedTappController {
	s.mu.Lock()
	defer s.mu.Unlock()
	tappController, ok := s.tappControllerMap[tappControllerName]
	if !ok {
		tappController = &cachedTappController{}
		s.tappControllerMap[tappControllerName] = tappController
	}
	return tappController
}

func (s *tappControllerCache) set(tappControllerName string, tappController *cachedTappController) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tappControllerMap[tappControllerName] = tappController
}

func (s *tappControllerCache) delete(tappControllerName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tappControllerMap, tappControllerName)
}
