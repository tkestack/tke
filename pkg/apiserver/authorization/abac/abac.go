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

// Portions Copyright 2014 The Kubernetes Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package abac authorizes Kubernetes API actions using an Attribute-based access control scheme.
package abac

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/howeyc/fsnotify"

	"tkestack.io/tke/pkg/util/log"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
)

const (
	watchDebounceDelay = 100 * time.Millisecond
)

type abacAuthorizer struct {
	policyList []*policy

	mutex   sync.RWMutex
	watcher *fsnotify.Watcher
}

// NewABACAuthorizer creates a attribute-base authorizer from policy file config,
// Forked from k8s.io/kubernetes/pkg/auth/authorizer/abac/abac.go
func NewABACAuthorizer(policyFile string) (authorizer.Authorizer, error) {
	policyList, err := loadPolicyConfig(policyFile)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watchDir, _ := filepath.Split(policyFile)
	if err := watcher.Watch(watchDir); err != nil {
		return nil, fmt.Errorf("could not watch %v: %v", watchDir, err)
	}

	abacAuthz := &abacAuthorizer{policyList: policyList, watcher: watcher}

	go abacAuthz.loadPolicy(policyFile)
	return abacAuthz, nil
}

// Authorize implements authorizer.Authorize
func (az *abacAuthorizer) Authorize(a authorizer.Attributes) (authorizer.Decision, string, error) {
	az.mutex.RLock()
	defer az.mutex.RUnlock()
	for _, p := range az.policyList {
		if matches(*p, a) {
			log.Debug("ABAC authorize success", log.Any("attr", a), log.Any("policy", *p))
			return authorizer.DecisionAllow, "", nil
		}
	}

	return authorizer.DecisionNoOpinion, "No policy matched.", nil
}

// loadPolicy watch abac policy file from the file system and reload policy list of the authorizer.
func (az *abacAuthorizer) loadPolicy(policyFile string) {
	defer az.watcher.Close()
	var timerC <-chan time.Time

	for {
		select {
		case <-timerC:
			timerC = nil
			log.Infof("Abac policy file may be changed, reload it")
			policyList, err := loadPolicyConfig(policyFile)
			if err != nil {
				log.Errorf("Parse policy file failed: %v", err)
				break
			}
			az.mutex.Lock()
			az.policyList = policyList
			az.mutex.Unlock()
		case event := <-az.watcher.Event:
			// use a timer to debounce configuration updates
			if (event.IsModify() || event.IsCreate()) && timerC == nil {
				timerC = time.After(watchDebounceDelay)
			}
		case err := <-az.watcher.Error:
			log.Errorf("Watcher error: %v", err)
		}
	}
}

// loadPolicyConfig attempts to create a policy list from the given file.
func loadPolicyConfig(path string) ([]*policy, error) {
	// File format is one map per line.  This allows easy concatenation of files,
	// comments in files, and identification of errors by line number.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	pl := make([]*policy, 0)

	i := 0
	for scanner.Scan() {
		i++
		p := &policy{}
		b := scanner.Bytes()

		// skip comment lines and blank lines
		trimmed := strings.TrimSpace(string(b))
		if len(trimmed) == 0 || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if err := json.Unmarshal(b, p); err != nil {
			log.Panic("failed", log.Err(err))
		}
		pl = append(pl, p)

	}

	if err := scanner.Err(); err != nil {
		return nil, policyLoadError{path, -1, nil, err}
	}

	return pl, nil
}

func matches(p policy, a authorizer.Attributes) bool {
	if subjectMatches(p, a.GetUser()) {
		if verbMatches(p, a) {
			// Resource and non-resource requests are mutually exclusive, at most one will match a policy
			if resourceMatches(p, a) {
				return true
			}

			if nonResourceMatches(p, a) {
				return true
			}

		}
	}
	return false
}

// subjectMatches returns true if specified user and group properties in the policy match the attributes
func subjectMatches(p policy, user user.Info) bool {
	matched := false

	if user == nil {
		return false
	}

	username := user.GetName()
	groups := user.GetGroups()

	// If the policy specified a user, ensure it matches
	if len(p.Spec.User) > 0 {
		if p.Spec.User == "*" {
			matched = true
		} else if regexMatch(p.Spec.User, username) {
			matched = true
		} else {
			matched = p.Spec.User == username
			if !matched {
				return false
			}
		}
	}

	// If the policy specified a group, ensure it matches
	if len(p.Spec.Group) > 0 {
		if p.Spec.Group == "*" {
			matched = true
		} else {
			matched = false
			for _, group := range groups {
				if p.Spec.Group == group || regexMatch(p.Spec.Group, group) {
					matched = true
				}
			}
			if !matched {
				return false
			}
		}
	}

	return matched
}

func verbMatches(p policy, a authorizer.Attributes) bool {

	// If policy specify verb, match verb first.
	if p.Spec.Verb != "" {
		if regexMatch(p.Spec.Verb, a.GetVerb()) {
			return true
		}
		return false
	}

	// All policies allow read only requests
	if a.IsReadOnly() {
		return true
	}

	// Allow if policy is not readonly
	if !p.Spec.Readonly {
		return true
	}

	return false
}

func nonResourceMatches(p policy, a authorizer.Attributes) bool {
	// A non-resource policy cannot match a resource request
	if !a.IsResourceRequest() {
		// Allow wildcard match
		if p.Spec.NonResourcePath == "*" {
			return true
		}
		// Allow exact match
		if p.Spec.NonResourcePath == a.GetPath() {
			return true
		}
		// Allow a trailing * subpath match
		if strings.HasSuffix(p.Spec.NonResourcePath, "*") && strings.HasPrefix(a.GetPath(), strings.TrimRight(p.Spec.NonResourcePath, "*")) {
			return true
		}
	}
	return false
}

func resourceMatches(p policy, a authorizer.Attributes) bool {
	// A resource policy cannot match a non-resource request
	if a.IsResourceRequest() {
		if p.Spec.Namespace == "*" || p.Spec.Namespace == a.GetNamespace() {
			if p.Spec.Resource == "*" || p.Spec.Resource == a.GetResource() || regexMatch(p.Spec.Resource, a.GetResource()) {
				if p.Spec.APIGroup == "*" || p.Spec.APIGroup == a.GetAPIGroup() || regexMatch(p.Spec.APIGroup, a.GetAPIGroup()) {
					return true
				}
			}
		}
	}
	return false
}

func regexMatch(pattern, str string) bool {
	pattern = strings.Replace(pattern, "*", ".*", -1)
	res, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}

	return res
}
