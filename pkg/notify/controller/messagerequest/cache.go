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

package messagerequest

import (
	"sync"

	v1 "tkestack.io/tke/api/notify/v1"
)

type cachedMessageRequest struct {
	// The cached state of the message request
	state *v1.MessageRequest
}

type messageRequestCache struct {
	mu                sync.Mutex // protects message request Map
	messageRequestMap map[string]*cachedMessageRequest
}

// ListKeys implements the interface required by DeltaFIFO to list the keys we
// already know about.
func (s *messageRequestCache) ListKeys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.messageRequestMap))
	for k := range s.messageRequestMap {
		keys = append(keys, k)
	}
	return keys
}

// GetByKey returns the value stored in the messageRequestMap under the given key
func (s *messageRequestCache) GetByKey(key string) (interface{}, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.messageRequestMap[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (s *messageRequestCache) get(messageRequestName string) (*cachedMessageRequest, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	messageRequest, ok := s.messageRequestMap[messageRequestName]
	return messageRequest, ok
}

func (s *messageRequestCache) Exist(messageRequestName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.messageRequestMap[messageRequestName]
	return ok
}

func (s *messageRequestCache) getOrCreate(messageRequestName string) *cachedMessageRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	messageRequest, ok := s.messageRequestMap[messageRequestName]
	if !ok {
		messageRequest = &cachedMessageRequest{}
		s.messageRequestMap[messageRequestName] = messageRequest
	}
	return messageRequest
}

func (s *messageRequestCache) set(messageRequestName string, messageRequest *cachedMessageRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageRequestMap[messageRequestName] = messageRequest
}

func (s *messageRequestCache) delete(messageRequestName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.messageRequestMap, messageRequestName)
}
