/*
 * Tencent is pleased to support the open source community by making TKEStack available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package util

import (
	"errors"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2/log"
	"github.com/casbin/casbin/v2/rbac"
)

type MatchingFunc func(arg1, arg2 string) bool

// RoleManager provides a default implementation for the RoleManager interface
type RoleManager struct {
	allRoles          *sync.Map
	maxHierarchyLevel int
	matchingFunc      MatchingFunc
}

// NewRoleManager is the constructor for creating an instance of the
// default RoleManager implementation.
func NewRoleManager(maxHierarchyLevel int) rbac.RoleManager {
	rm := RoleManager{}
	rm.allRoles = &sync.Map{}
	rm.maxHierarchyLevel = maxHierarchyLevel
	rm.matchingFunc = func(key1 string, key2 string) bool {
		return key1 == key2
	}

	return &rm
}

// AddMatchingFunc -
// e.BuildRoleLinks must be called after AddMatchingFunc().
func (rm *RoleManager) AddMatchingFunc(matchingFunc MatchingFunc) {
	rm.matchingFunc = func(key1 string, key2 string) bool {
		return key1 == key2
	}
}
func (rm *RoleManager) hasRole(name string) bool {
	var ok bool

	_, ok = rm.allRoles.Load(name)

	return ok
}

func (rm *RoleManager) createRole(name string) *Role {
	role, _ := rm.allRoles.LoadOrStore(name, newRole(name))

	return role.(*Role)
}

// Clear clears all stored data and resets the role manager to the initial state.
func (rm *RoleManager) Clear() error {
	rm.allRoles = &sync.Map{}
	return nil
}

// AddLink adds the inheritance link between role: name1 and role: name2.
// aka role: name1 inherits role: name2.
// domain is a prefix to the roles.
func (rm *RoleManager) AddLink(name1 string, name2 string, domain ...string) error {
	if len(domain) == 1 {
		name1 = domain[0] + "::" + name1
		name2 = domain[0] + "::" + name2
	} else {
		return errors.New("error: domain should be 1 parameter")
	}

	role1 := rm.createRole(name1)
	role2 := rm.createRole(name2)
	role1.addRole(role2)
	return nil
}

// DeleteLink deletes the inheritance link between role: name1 and role: name2.
// aka role: name1 does not inherit role: name2 any more.
// domain is a prefix to the roles.
func (rm *RoleManager) DeleteLink(name1 string, name2 string, domain ...string) error {
	if len(domain) == 1 {
		name1 = domain[0] + "::" + name1
		name2 = domain[0] + "::" + name2
	} else if len(domain) > 1 {
		return errors.New("error: domain should be 1 parameter")
	}

	if !rm.hasRole(name1) || !rm.hasRole(name2) {
		return errors.New("error: name1 or name2 does not exist")
	}

	role1 := rm.createRole(name1)
	role2 := rm.createRole(name2)
	role1.deleteRole(role2)
	return nil
}

// HasLink determines whether role: name1 inherits role: name2.
// domain is a prefix to the roles.
func (rm *RoleManager) HasLink(name1 string, name2 string, domain ...string) (bool, error) {
	if len(domain) != 1 {
		return false, errors.New("error: domain should be 1 parameter")
	}

	if domain[0] == "" {
		domain[0] = "*"
	}

	name1 = domain[0] + "::" + name1
	name2 = domain[0] + "::" + name2

	if !rm.hasRole(name1) || !rm.hasRole(name2) {
		return false, nil
	}

	if _, ok := rm.allRoles.Load(name2); !ok {
		defer rm.allRoles.Delete(name2)
	}

	rm.allRoles.LoadOrStore(name2, newRole(name2))

	role1 := rm.createRole(name1)
	return role1.hasRole(name2, rm.maxHierarchyLevel), nil
}

// GetRoles gets the roles that a subject inherits.
// domain is a prefix to the roles.
func (rm *RoleManager) GetRoles(name string, domain ...string) ([]string, error) {
	var domainName string
	if len(domain) == 1 {
		domainName = domain[0] + "::" + name
	} else {
		return nil, errors.New("error: domain should be 1 parameter")
	}

	if !rm.hasRole(domainName) {
		return []string{}, nil
	}
	res := make([]string, 0)
	simplificateRoles := func(roles []string) {
		for _, role := range roles {
			i := strings.Index(role, "::")
			if role[i+2:] != name {
				res = append(res, role[i+2:])
			}
		}
	}

	baseRoles := rm.createRole(domainName).getRoles()
	for _, role := range baseRoles {
		if rm.matchingFunc(domainName, role) {
			matchroles := rm.createRole(role).getRoles()
			simplificateRoles(matchroles)
		}
	}

	simplificateRoles(baseRoles)

	return res, nil
}

// GetUsers gets the users that inherits a subject.
// domain is an unreferenced parameter here, may be used in other implementations.
func (rm *RoleManager) GetUsers(name string, domain ...string) ([]string, error) {
	var domainName string
	if len(domain) == 1 {
		domainName = domain[0] + "::" + name
	} else {
		return nil, errors.New("error: domain should be 1 parameter")
	}

	if !rm.hasRole(domainName) {
		return nil, errors.New("error: name does not exist")
	}

	names := make([]string, 0)
	rm.allRoles.Range(func(_, value interface{}) bool {
		role := value.(*Role)
		if role.hasDirectRole(domainName, rm.matchingFunc) {
			names = append(names, role.name)
		}
		return true
	})
	var res []string
	for _, role := range names {
		i := strings.Index(role, "::")
		if role[i+2:] != name {
			res = append(res, role[i+2:])
		}
	}

	return res, nil
}

// PrintRoles prints all the roles to log.
func (rm *RoleManager) PrintRoles() error {
	if log.GetLogger().IsEnabled() {
		var sb strings.Builder
		rm.allRoles.Range(func(_, value interface{}) bool {
			if text := value.(*Role).toString(); text != "" {
				if sb.Len() == 0 {
					sb.WriteString(text)
				} else {
					sb.WriteString(", ")
					sb.WriteString(text)
				}
			}
			return true
		})
		log.LogPrint(sb.String())
	}
	return nil
}

// Role represents the data structure for a role in RBAC.
type Role struct {
	name  string
	roles []*Role
}

func newRole(name string) *Role {
	r := Role{}
	r.name = name
	return &r
}

func (r *Role) addRole(role *Role) {
	for _, rr := range r.roles {
		if rr.name == role.name {
			return
		}
	}

	r.roles = append(r.roles, role)
}

func (r *Role) deleteRole(role *Role) {
	for i, rr := range r.roles {
		if rr.name == role.name {
			r.roles = append(r.roles[:i], r.roles[i+1:]...)
			return
		}
	}
}

func (r *Role) hasRole(name string, hierarchyLevel int) bool {
	if r.name == name {
		return true
	}

	if hierarchyLevel <= 0 {
		return false
	}

	for _, role := range r.roles {
		if role.hasRole(name, hierarchyLevel-1) {
			return true
		}
	}

	return false
}

func (r *Role) hasDirectRole(name string, matchingFunc MatchingFunc) bool {
	for _, role := range r.roles {
		if role.name == name {
			return true
		}

		if matchingFunc(name, role.name) {
			return true
		}
	}

	return false
}

func (r *Role) toString() string {
	if len(r.roles) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(r.name)
	sb.WriteString(" < ")
	if len(r.roles) != 1 {
		sb.WriteString("(")
	}

	for i, role := range r.roles {
		if i == 0 {
			sb.WriteString(role.name)
		} else {
			sb.WriteString(", ")
			sb.WriteString(role.name)
		}
	}

	if len(r.roles) != 1 {
		sb.WriteString(")")
	}

	return sb.String()
}

func (r *Role) getRoles() []string {
	names := make([]string, 0)
	for _, role := range r.roles {
		names = append(names, role.name)
	}
	return names
}
