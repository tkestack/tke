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

package policy

import (
	"fmt"
	"github.com/pborman/uuid"
	"time"

	util2 "tkestack.io/tke/pkg/util"

	"k8s.io/apimachinery/pkg/util/sets"

	"tkestack.io/tke/pkg/auth/util"

	"tkestack.io/tke/pkg/auth/registry"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"tkestack.io/tke/pkg/auth/authorization/enforcer"
	"tkestack.io/tke/pkg/auth/registry/localidentity"
	"tkestack.io/tke/pkg/auth/registry/policy"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/validation"
)

const (
	// EffectAllow represents action on resource is allowed
	EffectAllow = "allow"

	// EffectDeny represents action on resource is denied
	EffectDeny = "deny"
)

var (
	// specialServices contains policy categories that will bind policy creator to the policy.
	specialServices = sets.NewString("tcr")
)

// Service is responsible for performing policy crud actions onto the storage backend.
type Service struct {
	store           *policy.Storage
	identityStorage *localidentity.Storage
	policyEnforcer  *enforcer.PolicyEnforcer
}

// NewPolicyService creates a new policy service object
func NewPolicyService(registry *registry.Registry, policyEnforcer *enforcer.PolicyEnforcer) *Service {
	policyEnforcer.StartSyncPolicy()
	return &Service{store: registry.PolicyStorage(), identityStorage: registry.LocalIdentityStorage(), policyEnforcer: policyEnforcer}
}

// CreatePolicy to create a new policy.
func (s *Service) CreatePolicy(policyCreate *types.Policy, attachUsers []string) (*types.Policy, error) {
	policyCreate.ID = fmt.Sprintf("%s%s-%s", types.PolicyIDPrefix, policyCreate.TenantID, uuid.New())
	policyCreate.CreateAt = time.Now()
	policyCreate.UpdateAt = policyCreate.CreateAt

	if err := s.store.Create(policyCreate); err != nil {
		return nil, err
	}

	err := s.policyEnforcer.AddPolicy(policyCreate)
	if err != nil {
		return nil, err
	}

	if len(attachUsers) == 0 && specialServices.Has(policyCreate.Service) {
		err = s.policyEnforcer.AddUsersPolicy(policyCreate.TenantID, policyCreate.ID, attachUsers)
		if err != nil {
			return nil, err
		}
	}

	return policyCreate, nil
}

// GetPolicy to return a policy by given id.
func (s *Service) GetPolicy(tenantID, id string) (*types.Policy, error) {
	p, err := s.store.Get(tenantID, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ListPolicies to return policies for given owner.
func (s *Service) ListPolicies(option *types.PolicyOption) (*types.PolicyList, error) {
	result := &types.PolicyList{Items: []*types.Policy{}}
	if len(option.ID) != 0 {
		p, err := s.GetPolicy(option.TenantID, option.ID)
		if err != nil {
			log.Debug("list policy failed", log.String("policy id", option.ID), log.Err(err))
			return result, nil
		}
		result.Items = append(result.Items, p)
		return result, nil
	}

	allPolicies, err := s.store.List(option.TenantID)
	if err != nil {
		log.Error("list all policies failed", log.Err(err))
		return result, nil
	}

	for _, p := range allPolicies.Items {
		// If specify username, match creator
		if option.UserName != "" && option.UserName != p.UserName {
			continue
		}

		switch option.Scope {
		case types.ListScopeLocal:
			if p.Type != types.UserDefine {
				continue
			}
		case types.ListScopeSystem:
			if p.Type != types.PreDefine {
				continue
			}
		}

		// Specify name, return policy
		if len(option.Name) != 0 {
			if p.Name == option.Name {
				result.Items = append(result.Items, p)
			}
			continue
		}

		// Keyword search
		if len(option.Keyword) != 0 {
			if util.CaseInsensitiveContains(p.Name, option.Keyword) {
				result.Items = append(result.Items, p)
			}
		} else {
			result.Items = append(result.Items, p)
		}
	}

	return result, nil
}

// DeletePolicy to delete a policy by given id.
func (s *Service) DeletePolicy(tenantID, id string) error {
	if err := s.store.Delete(tenantID, id); err != nil {
		log.Error("delete policy failed", log.String("policy id", id))
		return err
	}

	err := s.policyEnforcer.DeletePolicy(id)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePolicy to update a existing policy.
func (s *Service) UpdatePolicy(policyUpdate *types.Policy) (*types.Policy, error) {
	updater := func(old types.Policy) (types.Policy, error) {
		old.TenantID = policyUpdate.TenantID
		old.Description = policyUpdate.Description
		old.Statement = policyUpdate.Statement
		old.Name = policyUpdate.Name
		old.UpdateAt = time.Now()
		return old, nil
	}

	err := s.store.Update(policyUpdate.TenantID, policyUpdate.ID, updater)
	if err != nil {
		return nil, err
	}

	err = s.policyEnforcer.UpdatePolicy(policyUpdate)
	if err != nil {
		log.Error("update policy failed, sync rules to casbin failed", log.String("policy", policyUpdate.ID), log.Err(err))
	}

	p, err := s.store.Get(policyUpdate.TenantID, policyUpdate.ID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// AttachPolicyUsers to add policy for users.
func (s *Service) AttachPolicyUsers(request *types.AttachInfo) error {
	err := s.policyEnforcer.AddUsersPolicy(request.TenantID, request.ID, request.UserNames)
	if err != nil {
		log.Error("attach policy to user failed", log.String("policy", request.ID), log.Err(err))
		return err
	}

	return nil
}

// DetachPolicyUsers to remove policy for users.
func (s *Service) DetachPolicyUsers(request *types.AttachInfo) error {
	_, err := s.store.Get(request.TenantID, request.ID)
	if err != nil {
		return err
	}

	err = s.policyEnforcer.RemovePolicyUsers(request.TenantID, request.ID, request.UserNames)
	if err != nil {
		log.Error("detach policy to user failed", log.String("policy", request.ID), log.Err(err))
		return err
	}

	return nil
}

// ListPolicyUsers to list all users attached to the policy.
func (s *Service) ListPolicyUsers(tenantID, id string) (*types.AttachInfo, error) {
	_, err := s.store.Get(tenantID, id)
	if err != nil {
		return nil, err
	}
	info := &types.AttachInfo{ID: id, TenantID: tenantID}
	userNames, err := s.policyEnforcer.ListPolicyUsers(tenantID, id)
	if err != nil {
		log.Error("list users attached to policy failed", log.String("policy", id), log.Err(err))
		return nil, err
	}

	info.UserNames = userNames

	return info, nil
}

// ListUserPolicies to return all policies id of the user.
func (s *Service) ListUserPolicies(tenantID, userName string) ([]string, error) {
	policyIDs, err := s.policyEnforcer.ListUserPolicies(tenantID, userName)
	if err != nil {
		log.Error("list all policies attached to user failed", log.String("user", userName), log.Err(err))
		return nil, err
	}

	return policyIDs, nil
}

// LoadPredefinePolicies creates or updates predefine policies from input params.
func (s *Service) LoadPredefinePolicies(tenantID string, policyList []*types.Policy) error {
	existPolicyMap := map[string]*types.Policy{}

	policyAll, err := s.store.List(tenantID)
	if err != nil {
		log.Error("List all policy for tenant failed", log.String("tenant", tenantID), log.Err(err))
		return err
	}

	for _, pol := range policyAll.Items {
		if pol.Type == types.PreDefine {
			existPolicyMap[pol.Name] = pol
		}
	}
	for _, pol := range policyList {
		pol.TenantID = tenantID
		pol.Type = types.PreDefine
		if polExist, ok := existPolicyMap[pol.Name]; !ok {
			log.Info("Predefine policy not exists, create it", log.String("name", pol.Name))
			if _, err := s.CreatePolicy(pol, nil); err != nil {
				log.Error("Create predefine policy failed", log.String("name", pol.Name), log.Err(err))
				return err
			}
		} else {
			pol.ID = polExist.ID
			isChanged := false
			if polExist.Description != pol.Description || polExist.Statement.Effect != pol.Statement.Effect ||
				polExist.Statement.Resource != pol.Statement.Resource {
				isChanged = true
			}

			if added, _ := util2.DiffStringSlice(polExist.Statement.Actions, pol.Statement.Actions); len(added) > 0 {
				polExist.Statement.Actions = append(polExist.Statement.Actions, added...)
				pol.Statement.Actions = polExist.Statement.Actions
				isChanged = true
			}

			if isChanged {
				_, err := s.UpdatePolicy(pol)
				if err != nil {
					log.Error("Update policy failed", log.Err(err))
				}
			}
		}
	}

	return nil
}

func validatePolicyCreate(policy *types.Policy) error {
	allErrs := field.ErrorList{}
	if err := validation.IsDisplayName(policy.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), policy.Name, `must be specified`))
	}

	if policy.Statement.Effect != EffectAllow && policy.Statement.Effect != EffectDeny {
		allErrs = append(allErrs, field.Invalid(field.NewPath("effect"), policy.Statement.Effect, `must be allow or deny`))
	}
	return allErrs.ToAggregate()
}
