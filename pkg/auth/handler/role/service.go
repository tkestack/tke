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

package role

import (
	"fmt"
	"github.com/pborman/uuid"
	"time"

	"tkestack.io/tke/pkg/auth/util"

	"tkestack.io/tke/pkg/auth/registry"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"tkestack.io/tke/pkg/auth/authorization/enforcer"
	"tkestack.io/tke/pkg/auth/registry/policy"
	"tkestack.io/tke/pkg/auth/registry/role"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/validation"
)

// Service is responsible for performing role crud actions onto the storage backend.
type Service struct {
	store          *role.Storage
	policyStore    *policy.Storage
	policyEnforcer *enforcer.PolicyEnforcer
}

// NewRoleService creates a new role service object
func NewRoleService(registry *registry.Registry, policyEnforcer *enforcer.PolicyEnforcer) *Service {
	policyEnforcer.StartSyncRoles()
	return &Service{store: registry.RoleStorage(), policyStore: registry.PolicyStorage(), policyEnforcer: policyEnforcer}
}

// CreateRole to create a new role with policies.
func (s *Service) CreateRole(roleCreate *types.Role) (*types.Role, error) {
	log.Debug("Create RoleStorage", log.Any("role", roleCreate))

	roleCreate.ID = fmt.Sprintf("%s%s-%s", types.RoleIDPrefix, roleCreate.TenantID, uuid.New())
	roleCreate.CreateAt = time.Now()
	roleCreate.UpdateAt = roleCreate.CreateAt

	if err := s.store.Create(roleCreate); err != nil {
		return nil, err
	}

	var pids []string
	for _, pol := range roleCreate.Policies {
		pids = append(pids, pol.ID)
	}
	err := s.policyEnforcer.AddRolePolicies(roleCreate.ID, pids)
	if err != nil {
		return nil, err
	}

	s.fillRolePolicies(roleCreate)

	return roleCreate, nil
}

// GetRole to return a role by given ID.
func (s *Service) GetRole(tenantID, ID string) (*types.Role, error) {
	rl, err := s.store.Get(tenantID, ID)
	if err != nil {
		return nil, err
	}

	s.fillRolePolicies(rl)
	return rl, nil
}

// ListRoles to return roles for given options.
func (s *Service) ListRoles(option *types.RoleOption) (*types.RoleList, error) {
	result := &types.RoleList{}
	if len(option.ID) != 0 {
		rl, err := s.GetRole(option.TenantID, option.ID)
		if err != nil {
			log.Debug("list roles failed", log.String("role ID", option.ID), log.Err(err))
			return result, nil
		}
		result.Items = append(result.Items, rl)
		return result, nil
	}
	allRoles, err := s.store.List(option.TenantID)
	if err != nil {
		log.Error("list all policies failed", log.Err(err))
		return result, nil
	}

	for _, rl := range allRoles.Items {
		// If specify username, match creator
		if option.UserName != "" && option.UserName != rl.UserName {
			continue
		}

		switch option.Scope {
		case types.ListScopeLocal:
			if rl.Type != types.UserDefine {
				continue
			}
		case types.ListScopeSystem:
			if rl.Type != types.PreDefine {
				continue
			}
		}

		//Specify name, return role
		if len(option.Name) != 0 {
			if rl.Name == option.Name {
				s.fillRolePolicies(rl)
				result.Items = append(result.Items, rl)
			}
			continue
		}

		//Keyword search
		if len(option.Keyword) != 0 {
			if util.CaseInsensitiveContains(rl.Name, option.Keyword) {
				s.fillRolePolicies(rl)
				result.Items = append(result.Items, rl)
				continue
			}
		} else {
			//Return all roles for owner, include pre-define role
			s.fillRolePolicies(rl)
			result.Items = append(result.Items, rl)
		}
	}

	return result, nil
}

// UpdateRole to update a existing role metadata.
func (s *Service) UpdateRole(roleUpdate *types.Role) (*types.Role, error) {
	updater := func(old types.Role) (types.Role, error) {
		old.Description = roleUpdate.Description
		old.Name = roleUpdate.Name
		old.UpdateAt = time.Now()
		return old, nil
	}

	err := s.store.Update(roleUpdate.TenantID, roleUpdate.ID, updater)
	if err != nil {
		return nil, err
	}

	rl, err := s.store.Get(roleUpdate.TenantID, roleUpdate.ID)
	if err != nil {
		return nil, err
	}

	s.fillRolePolicies(rl)
	return nil, err
}

// DeleteRole to delete a existing role.
func (s *Service) DeleteRole(tenantID, id string) error {
	if err := s.store.Delete(tenantID, id); err != nil {
		log.Error("delete role failed", log.String("role ID", id))
		return err
	}

	err := s.policyEnforcer.DeleteRole(id)
	if err != nil {
		return err
	}

	return nil
}

// AttachRolePolicies to bind polices for role.
func (s *Service) AttachRolePolicies(info *types.AttachInfo) error {
	_, err := s.store.Get(info.TenantID, info.ID)
	if err != nil {
		return err
	}

	var validPolicies []string
	for _, pid := range info.PolicyIDs {
		_, err := s.policyStore.Get(info.TenantID, pid)
		if err == nil {
			validPolicies = append(validPolicies, pid)
		}
	}

	if err := s.policyEnforcer.AddRolePolicies(info.ID, validPolicies); err != nil {
		log.Error("attach policies to role failed", log.String("role", info.ID), log.Err(err))
		return err
	}

	return nil
}

// DetachRolePolicies to unbind polices for role.
func (s *Service) DetachRolePolicies(info *types.AttachInfo) error {
	_, err := s.store.Get(info.TenantID, info.ID)
	if err != nil {
		return err
	}

	if err := s.policyEnforcer.RemoveRolePolicies(info.ID, info.PolicyIDs); err != nil {
		log.Error("detach policies to role failed", log.String("role", info.ID), log.Err(err))
		return err
	}

	return nil
}

// ListRolePolicies to return all policies attached to the policy.
func (s *Service) ListRolePolicies(tenantID, id string) (*types.PolicyList, error) {
	_, err := s.store.Get(tenantID, id)
	if err != nil {
		return nil, err
	}
	info := &types.AttachInfo{ID: id, TenantID: tenantID}
	policyIDs, err := s.policyEnforcer.ListRolePolicies(tenantID, id)
	if err != nil {
		return nil, err
	}
	info.PolicyIDs = policyIDs
	policyList := types.PolicyList{Items: []*types.Policy{}}
	for _, pid := range policyIDs {
		pol, err := s.policyStore.Get(tenantID, pid)
		if err == nil {
			policyList.Items = append(policyList.Items, pol)
		}
	}

	return &policyList, nil
}

// AttachUsersRole to bind role for users.
func (s *Service) AttachUsersRole(info *types.AttachInfo) error {
	err := s.policyEnforcer.AddUsersPolicy(info.TenantID, info.ID, info.UserNames)
	if err != nil {
		log.Error("attach role to user failed", log.String("role", info.ID), log.Err(err))
		return err
	}

	return nil
}

// DetachUsersRole to unbind role for users
func (s *Service) DetachUsersRole(info *types.AttachInfo) error {
	_, err := s.store.Get(info.TenantID, info.ID)
	if err != nil {
		return err
	}

	err = s.policyEnforcer.RemovePolicyUsers(info.TenantID, info.ID, info.UserNames)
	if err != nil {
		log.Error("detach role to user failed", log.String("role", info.ID), log.Err(err))
		return err
	}

	return nil
}

// ListRoleUsers to list all users attached to the role.
func (s *Service) ListRoleUsers(tenantID, id string) (*types.AttachInfo, error) {
	_, err := s.store.Get(tenantID, id)
	if err != nil {
		return nil, err
	}
	info := &types.AttachInfo{ID: id}
	userNames, err := s.policyEnforcer.ListPolicyUsers(tenantID, id)
	if err != nil {
		log.Error("list users attached to role failed", log.String("role", id), log.Err(err))
		return nil, err
	}

	info.UserNames = userNames
	return info, nil
}

// validateRoleCreate to validate role name and policies
func (s *Service) validateRoleCreate(role *types.Role) error {
	allErrs := field.ErrorList{}

	if err := validation.IsDisplayName(role.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), role.Name, err.Error()))
	}

	for i, pol := range role.Policies {
		if p, err := s.policyStore.Get(role.TenantID, pol.ID); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath(fmt.Sprintf("polices[%d]", i)), pol, err.Error()))
		} else if p.TenantID != role.TenantID {
			allErrs = append(allErrs, field.Invalid(field.NewPath(fmt.Sprintf("polices[%d]", i)), pol, "policy and role not in same tenant"))
		}
	}
	return allErrs.ToAggregate()
}

func (s *Service) fillRolePolicies(role *types.Role) {
	policyIDs, err := s.policyEnforcer.ListRolePolicies(role.TenantID, role.ID)
	if err != nil {
		log.Error("list role policies failed", log.Err(err))
		return
	}

	role.Policies = []*types.PolicyMeta{}
	for _, pid := range policyIDs {
		if pol, err := s.policyStore.Get(role.TenantID, pid); err == nil {
			policyInfo := types.PolicyMeta{}
			policyInfo.Name = pol.Name
			policyInfo.Type = pol.Type
			policyInfo.TenantID = pol.TenantID
			policyInfo.Service = pol.Service
			role.Policies = append(role.Policies, &policyInfo)
		}
	}
}
