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
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/log"
)

// AddRolePolicies to add some policies into role.
func (pe *PolicyEnforcer) AddRolePolicies(roleID string, policyIDs []string) error {
	for _, pid := range policyIDs {
		if _, err := pe.enforcer.AddGroupingPolicySafe(roleID, pid); err != nil {
			log.Error("add new policy to role failed", log.Any("policy", pid), log.String("role", roleID), log.Err(err))
			return err
		}
	}

	return nil
}

// RemoveRolePolicies to remove some policy from role.
func (pe *PolicyEnforcer) RemoveRolePolicies(roleID string, policyIDs []string) error {
	for _, id := range policyIDs {
		if _, err := pe.enforcer.RemoveFilteredGroupingPolicySafe(idIndex, roleID, id); err != nil {
			log.Error("remove policy from role failed", log.Any("policy", id), log.String("role", roleID), log.Err(err))
			return err
		}
	}

	return nil
}

// ListRolePolicies to return all policies attached to the policy.
func (pe *PolicyEnforcer) ListRolePolicies(tenantID, roleID string) (policyIDs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	rules := pe.enforcer.GetFilteredGroupingPolicy(idIndex, roleID)
	log.Debugf("Get grouping rules role: %s, %v", roleID, rules)
	for _, rule := range rules {
		if parseTenantID(rule[1]) == tenantID {
			if strings.HasPrefix(rule[1], types.PolicyIDPrefix) {
				policyIDs = append(policyIDs, rule[1])
			}
		}
	}

	return
}

// ListUserRoles to return all roles for user.
func (pe *PolicyEnforcer) ListUserRoles(tenantID, userName string) (roleIDs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// will get policies and roles that a user has, filter roles.
	roles := pe.enforcer.GetRolesForUser(keyUser(tenantID, userName))
	log.Debugf("Get policies or roles for user: %s, %v", userName, roles)
	for _, role := range roles {
		if parseTenantID(role) == tenantID {
			if strings.HasPrefix(role, types.RoleIDPrefix) {
				roleIDs = append(roleIDs, role)
			}
		}
	}

	return
}

// DeleteRole to delete a existing role.
func (pe *PolicyEnforcer) DeleteRole(roleID string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	// Delete role will remove policy and role inheritance
	pe.enforcer.DeleteRole(roleID)

	return err
}

// StartSyncRoles is to ensure policies of the roles into casbin.
func (pe *PolicyEnforcer) StartSyncRoles() {
	go func() {
		for {
			// error intentionally ignored
			_ = pe.syncRoles()

			time.Sleep(policySyncedInterval)
		}
	}()
}

func (pe *PolicyEnforcer) syncRoles() error {
	allTenants, err := pe.registry.DexStorage().ListConnectors()
	if err != nil {
		return err
	}

	for _, tenant := range allTenants {
		roleList, err := pe.registry.RoleStorage().List(tenant.ID)
		if err != nil {
			log.Error("List all roles failed", log.String("tenant", tenant.ID), log.Err(err))
			return err
		}

		existsRoleIDs, err := pe.listRoleIDs(tenant.ID)
		if err != nil {
			log.Error("List all roles id from casbin failed", log.String("tenant", tenant.ID), log.Err(err))
			continue
		}

		neededRoleIDs := sets.NewString()
		for _, role := range roleList.Items {
			neededRoleIDs.Insert(role.ID)
			_ = pe.syncRole(role.ID)
		}

		for _, id := range existsRoleIDs {
			if !neededRoleIDs.Has(id) {
				log.Info("sync role: role in casbin is not needed, remove all related rules", log.String("policyID", id))
				if err := pe.DeleteRole(id); err != nil {
					log.Error("sync role: delete role failed", log.Err(err))
				}
			}
		}
	}

	return nil
}

func (pe *PolicyEnforcer) syncRole(id string) error {
	tenantID := parseTenantID(id)
	if len(tenantID) == 0 {
		log.Warn("Not legal role id, no tenantID", log.String("role", id))
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

// listRoleIDs returns all roles id in casbin.
func (pe *PolicyEnforcer) listRoleIDs(tenantID string) ([]string, error) {
	idSets := sets.NewString()
	rules := pe.enforcer.GetAllRoles()

	for _, rule := range rules {
		if strings.HasPrefix(rule, types.RoleIDPrefix) {
			if parseTenantID(rule) == tenantID {
				idSets.Insert(rule)
			}
		}
	}

	return idSets.List(), nil
}
