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

package project

import (
	"sync"
	"tkestack.io/tke/api/business/v1"
)

type cachedProject struct {
	// The cached state of the value
	state *v1.Project
}

type projectCache struct {
	mu sync.Mutex
	m  map[string]*cachedProject
}

func (s *projectCache) getOrCreate(name string) *cachedProject {
	s.mu.Lock()
	defer s.mu.Unlock()
	project, ok := s.m[name]
	if !ok {
		project = &cachedProject{}
		s.m[name] = project
	}
	return project
}

func (s *projectCache) get(name string) (*cachedProject, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	project, ok := s.m[name]
	return project, ok
}

func (s *projectCache) Exist(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[name]
	return ok
}

func (s *projectCache) set(name string, project *cachedProject) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[name] = project
}

func (s *projectCache) delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, name)
}
