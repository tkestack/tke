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

package app

import (
	"sync"

	applicationv1 "tkestack.io/tke/api/application/v1"
)

type cachedApp struct {
	// The cached state of the value
	state *applicationv1.App
}

type applicationCache struct {
	mu sync.Mutex
	m  map[string]*cachedApp
}

func (s *applicationCache) getOrCreate(name string) *cachedApp {
	s.mu.Lock()
	defer s.mu.Unlock()
	app, ok := s.m[name]
	if !ok {
		app = &cachedApp{}
		s.m[name] = app
	}
	return app
}

func (s *applicationCache) get(name string) (*cachedApp, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	app, ok := s.m[name]
	return app, ok
}

func (s *applicationCache) Exist(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[name]
	return ok
}

func (s *applicationCache) set(name string, app *cachedApp) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[name] = app
}

func (s *applicationCache) delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, name)
}
