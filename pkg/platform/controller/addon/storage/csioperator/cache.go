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

package csioperator

import (
	"sync"

	v1 "tkestack.io/tke/api/platform/v1"
)

type cachedCSIOperator struct {
	// The cached state of the operator
	state *v1.CSIOperator
}

type csiOperatorCache struct {
	mu    sync.Mutex // protects coMap
	coMap map[string]*cachedCSIOperator
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *csiOperatorCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.coMap))
	for k := range s.coMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the coMap under the given key
func (s *csiOperatorCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.coMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *csiOperatorCache) get(operatorName string) (*cachedCSIOperator, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	operator, ok := s.coMap[operatorName]
	return operator, ok
}

func (s *csiOperatorCache) Exist(operatorName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.coMap[operatorName]
	return ok
}

func (s *csiOperatorCache) getOrCreate(operatorName string) *cachedCSIOperator {
	s.mu.Lock()
	defer s.mu.Unlock()
	operator, ok := s.coMap[operatorName]
	if !ok {
		operator = &cachedCSIOperator{}
		s.coMap[operatorName] = operator
	}
	return operator
}

func (s *csiOperatorCache) set(operatorName string, operator *cachedCSIOperator) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.coMap[operatorName] = operator
}

func (s *csiOperatorCache) delete(operatorName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.coMap, operatorName)
}
