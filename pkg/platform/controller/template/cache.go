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

package template

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedTemplate struct {
	// The cached state of the value
	state *v1.Template
}

type templateCache struct {
	mu sync.Mutex // protects m
	m  map[string]*cachedTemplate
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *templateCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the m under the given key
func (s *templateCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.m[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *templateCache) Exist(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[name]
	return ok
}

func (s *templateCache) getOrCreate(name string) *cachedTemplate {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.m[name]
	if !ok {
		value = &cachedTemplate{}
		s.m[name] = value
	}
	return value
}

func (s *templateCache) set(name string, value *cachedTemplate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[name] = value
}

func (s *templateCache) delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, name)
}

func (s *templateCache) get(key string) (*cachedTemplate, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	template, ok := s.m[key]
	return template, ok
}
