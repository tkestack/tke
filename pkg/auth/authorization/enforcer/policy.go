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

package enforcer

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	"tkestack.io/tke/pkg/auth/types"

	"tkestack.io/tke/pkg/util/log"
)

const (
	idIndex = iota
)

var (
	policySyncedInterval = 1 * time.Minute
)

// AddPolicy to create a new policy into casbin.
func (pe *PolicyEnforcer) AddPolicy(policy *types.Policy) error {
	rules := convertPolicy(policy)
	for _, rule := range rules {
		if _, err := pe.enforcer.AddPolicySafe(rule); err != nil {
			// If add one failed, delete policy's all rules
			_ = pe.DeletePolicy(policy.ID)
			log.Error("add policy failed", log.Any("policy", policy), log.Err(err))
			return err
		}
	}

	return nil
}

// UpdatePolicy to update a existing policy, delete old rules and create new rules.
func (pe *PolicyEnforcer) UpdatePolicy(policy *types.Policy) error {
	rules := convertPolicy(policy)
	err := pe.syncPolicy(policy.ID, rules)
	if err != nil {
		return err
	}

	return nil
}

// DeletePolicy to delete a existing policy.
func (pe *PolicyEnforcer) DeletePolicy(pid string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	// Delete role will remove policy and role inheritance
	pe.enforcer.DeleteRole(pid)

	return err
}

// StartSyncPolicy is to ensure actions of policies into casbin.
func (pe *PolicyEnforcer) StartSyncPolicy() {
	go func() {
		for {
			// error intentionally ignored
			_ = pe.syncPolicies()

			time.Sleep(policySyncedInterval)
		}
	}()
}

func (pe *PolicyEnforcer) syncPolicies() error {
	allTenants, err := pe.registry.DexStorage().ListConnectors()
	if err != nil {
		return err
	}

	for _, tenant := range allTenants {
		policyList, err := pe.registry.PolicyStorage().List(tenant.ID)
		if err != nil {
			log.Error("List all policies failed", log.String("tenant", tenant.ID), log.Err(err))
			continue
		}

		existsPolicyID, err := pe.listPolicyIDs(tenant.ID)
		if err != nil {
			log.Error("List all policies id from casbin failed", log.String("tenant", tenant.ID), log.Err(err))
			continue
		}

		neededPolicyIDs := sets.NewString()
		for _, p := range policyList.Items {
			neededPolicyIDs.Insert(p.ID)
			rules := convertPolicy(p)
			_ = pe.syncPolicy(p.ID, rules)
		}

		for _, id := range existsPolicyID {
			if !neededPolicyIDs.Has(id) {
				log.Info("sync policy: policy in casbin is not needed, remove all related rules", log.String("policyID", id))
				if err := pe.DeletePolicy(id); err != nil {
					log.Error("sync policy: delete policy failed", log.Err(err))
				}
			}
		}
	}

	return nil
}

func (pe *PolicyEnforcer) syncPolicy(id string, rules [][]string) error {
	// Add new rule if not exists
	for _, rule := range rules {
		if existsRule := pe.enforcer.GetFilteredPolicy(idIndex, rule...); len(existsRule) == 0 {
			log.Info("sync policy: rule not found, add it to casbin", log.Strings("rule", rule))
			if _, err := pe.enforcer.AddPolicySafe(rule); err != nil {
				log.Errorf("sync policy failed: %+v", err)
				return err
			}
		}
	}

	// Remove rule if not belong to policy
	existsRules := pe.enforcer.GetFilteredPolicy(idIndex, id)
	for _, rule := range existsRules {
		if ruleShouldBeRemoved(rule, rules) {
			log.Info("sync policy: rule not need, will remove it from casbin", log.Strings("rule", rule))
			if _, err := pe.enforcer.RemoveFilteredPolicySafe(idIndex, rule...); err != nil {
				return err
			}
		}
	}

	tenantID := parseTenantID(id)
	if len(tenantID) == 0 {
		log.Warn("Not legal policy id, cannot parse tenantID, delete it", log.String("role", id))
		return fmt.Errorf("cannot parse tenantID from id: %s", id)
	}

	// Todo check user exists with 3rd idp
	// Remove user rule if not exists, comment it to support 3rd idp
	//userRules := pe.enforcer.GetFilteredGroupingPolicy(idIndex+1, id)
	//for _, rule := range userRules {
	//	if strings.HasPrefix(rule[0], types.UserPrefix) {
	//		userName := splitUserPrefix(tenantID, rule[0])
	//		_, err := pe.registry.LocalIdentityStorage().Get(tenantID, userName)
	//		if err != nil && err == etcd.ErrNotFound {
	//			log.Info("User is not exists, unbind it", log.Strings("rule", rule))
	//			if _, err := pe.enforcer.RemoveFilteredGroupingPolicySafe(idIndex, rule...); err != nil {
	//				log.Errorf("sync role failed: %+v", err)
	//				return err
	//			}
	//		}
	//	}
	//}

	return nil
}

// AddUsersPolicy to attach users to the policy.
func (pe *PolicyEnforcer) AddUsersPolicy(tenantID, id string, usernames []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	for _, name := range usernames {
		pe.enforcer.AddRoleForUser(keyUser(tenantID, name), id)
	}

	return nil
}

// RemovePolicyUsers to remove users from the policy.
func (pe *PolicyEnforcer) RemovePolicyUsers(tenantID, id string, usernames []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	for _, name := range usernames {
		pe.enforcer.DeleteRoleForUser(keyUser(tenantID, name), id)
	}

	return nil
}

// ListPolicyUsers to return all users attached to the policy.
func (pe *PolicyEnforcer) ListPolicyUsers(tenantID, policyID string) (userNames []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	rules := pe.enforcer.GetFilteredGroupingPolicy(1, policyID)
	log.Debugf("Get grouping rules for policy: %s, %v", policyID, rules)
	for _, rule := range rules {
		if strings.HasPrefix(rule[0], fmt.Sprintf("%s%s-", types.UserPrefix, tenantID)) {
			userNames = append(userNames, splitUserPrefix(tenantID, rule[0]))
		}
	}

	return
}

// ListUserPolicies to return all policies for user.
func (pe *PolicyEnforcer) ListUserPolicies(tenantID, userName string) (policyIDs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// will get policies and roles that a user has, filter roles.
	roles := pe.enforcer.GetRolesForUser(keyUser(tenantID, userName))
	log.Infof("Get policies or roles for user: %s, %v", userName, roles)
	for _, role := range roles {
		if strings.HasPrefix(role, types.PolicyIDPrefix) {
			if parseTenantID(role) == tenantID {
				policyIDs = append(policyIDs, role)
			}
		}
	}

	return
}

// RemoveAllPermsForUser removes all policies or roles for user.
func (pe *PolicyEnforcer) RemoveAllPermsForUser(tenantID, userName string) error {
	// will get policies and roles that a user has.
	roles := pe.enforcer.GetRolesForUser(keyUser(tenantID, userName))
	log.Infof("Remove policies or roles for user: %s, %v", userName, roles)
	for _, role := range roles {
		if parseTenantID(role) == tenantID {
			if strings.HasPrefix(role, types.PolicyIDPrefix) ||
				strings.HasPrefix(role, types.RoleIDPrefix) {
				if err := pe.RemovePolicyUsers(tenantID, role, []string{userName}); err != nil {
					log.Error("Remove policy or role for user failed", log.String("policy or role", role), log.String("user", userName), log.Err(err))
				}
			}
		}
	}

	return nil
}

// listPolicyIDs returns all policy id in casbin
func (pe *PolicyEnforcer) listPolicyIDs(tenantID string) ([]string, error) {
	idSets := sets.NewString()
	policyRules := pe.enforcer.GetPolicy()

	for _, rule := range policyRules {
		if parseTenantID(rule[0]) == tenantID {
			if strings.HasPrefix(rule[0], types.PolicyIDPrefix) {
				idSets.Insert(rule[0])
			}
		}
	}

	// Policy may be deleted but relation with roles or users may be exists
	roleRules := pe.enforcer.GetAllRoles()
	for _, rule := range roleRules {
		if parseTenantID(rule) == tenantID {
			if strings.HasPrefix(rule, types.PolicyIDPrefix) {
				idSets.Insert(rule)
			}
		}
	}

	return idSets.List(), nil
}

// convertPolicy to convert policy to casbin rule.
func convertPolicy(policy *types.Policy) [][]string {
	var rules [][]string
	for _, act := range policy.Statement.Actions {
		rule := []string{policy.ID, policy.Statement.Resource, act, policy.Statement.Effect}
		rules = append(rules, rule)
	}

	return rules
}

// ruleShouldBeRemoved to check whether rule exists in casbin.
func ruleShouldBeRemoved(rule []string, rules [][]string) bool {
	for _, val := range rules {
		if reflect.DeepEqual(rule, val) {
			return false
		}
	}

	return true
}

func keyUser(tenantID string, name string) string {
	return fmt.Sprintf("%s%s-%s", types.UserPrefix, tenantID, name)
}

func splitUserPrefix(tenantID string, str string) string {
	return strings.TrimPrefix(str, fmt.Sprintf("%s%s-", types.UserPrefix, tenantID))
}

// split tenanID from policy or role id, {prefix}-{tenantID}-{uuid}
func parseTenantID(id string) string {
	array := strings.Split(id, "-")
	if len(array) >= 7 {
		array = array[1 : len(array)-5]
		return strings.Join(array, "-")
	}

	return ""
}
