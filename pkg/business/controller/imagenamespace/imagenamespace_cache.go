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

package imagenamespace

import (
	"sync"

	v1 "tkestack.io/tke/api/business/v1"
)

type cachedImageNamespace struct {
	// The cached state of the value
	state *v1.ImageNamespace
}

type imageNamespaceCache struct {
	mu sync.Mutex
	m  map[string]*cachedImageNamespace
}

func (s *imageNamespaceCache) getOrCreate(name string) *cachedImageNamespace {
	s.mu.Lock()
	defer s.mu.Unlock()
	imageNamespace, ok := s.m[name]
	if !ok {
		imageNamespace = &cachedImageNamespace{}
		s.m[name] = imageNamespace
	}
	return imageNamespace
}

func (s *imageNamespaceCache) get(name string) (*cachedImageNamespace, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	imageNamespace, ok := s.m[name]
	return imageNamespace, ok
}

func (s *imageNamespaceCache) Exist(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[name]
	return ok
}

func (s *imageNamespaceCache) set(name string, imageNamespace *cachedImageNamespace) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[name] = imageNamespace
}

func (s *imageNamespaceCache) delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, name)
}
