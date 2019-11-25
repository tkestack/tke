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

package localidentity

import (
	"fmt"
	"github.com/pborman/uuid"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"time"
	"tkestack.io/tke/pkg/auth/authorization/enforcer"
	policyService "tkestack.io/tke/pkg/auth/handler/policy"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/registry/localidentity"
	"tkestack.io/tke/pkg/auth/registry/policy"
	"tkestack.io/tke/pkg/auth/registry/role"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/validation"
)

var (
	reservedNames = map[string]bool{"system": true}
	maskPassword  = "******"
)

// Service is responsible for performing local identity crud actions onto the storage backend.
type Service struct {
	store          *localidentity.Storage
	policyStore    *policy.Storage
	roleStore      *role.Storage
	policyEnforcer *enforcer.PolicyEnforcer
}

// NewLocalIdentityService creates a new local identity service object
func NewLocalIdentityService(registry *registry.Registry, policyEnforcer *enforcer.PolicyEnforcer) *Service {
	return &Service{
		store:          registry.LocalIdentityStorage(),
		policyStore:    registry.PolicyStorage(),
		roleStore:      registry.RoleStorage(),
		policyEnforcer: policyEnforcer,
	}
}

// CreateLocalIdentity create a new TKE User and returns it.
func (s *Service) CreateLocalIdentity(identityCreate *types.LocalIdentity) (*types.LocalIdentity, error) {
	identityCreate.UID = uuid.New()
	identityCreate.CreateAt = time.Now()
	identityCreate.UpdateAt = identityCreate.CreateAt
	if identityCreate.Status == nil {
		identityCreate.Status = &types.LocalIdentityStatus{Locked: false}
	}

	bcryptedPasswd, err := util.BcryptPassword(identityCreate.Spec.HashedPassword)
	if err != nil {
		log.Error("Bcrypt hash password failed", log.Err(err))
		return nil, err
	}

	identityCreate.Spec.HashedPassword = bcryptedPasswd
	err = s.store.Create(identityCreate)
	if err != nil {
		return nil, err
	}
	return identityCreate, nil
}

// UpdateLocalIdentity to update a existing user.
func (s *Service) UpdateLocalIdentity(identityUpdate *types.LocalIdentity) (*types.LocalIdentity, error) {
	// UserName can not be changed
	oldIdentity, err := s.store.Get(identityUpdate.Spec.TenantID, identityUpdate.Name)
	if err != nil {
		return nil, err
	}

	updater := func(current types.LocalIdentity) (types.LocalIdentity, error) {
		current.Spec.TenantID = identityUpdate.Spec.TenantID
		for key, val := range identityUpdate.Spec.Extra {
			current.Spec.Extra[key] = val
		}

		if len(identityUpdate.Spec.HashedPassword) != 0 {
			current.Spec.HashedPassword, err = util.BcryptPassword(identityUpdate.Spec.HashedPassword)
			if err != nil {
				return types.LocalIdentity{}, err
			}
		}
		current.UpdateAt = time.Now()
		current.CreateAt = oldIdentity.CreateAt

		return current, nil
	}

	err = s.store.Update(identityUpdate.Spec.TenantID, identityUpdate.Name, updater)
	if err != nil {
		return nil, err
	}

	identity, err := s.store.Get(identityUpdate.Spec.TenantID, identityUpdate.Name)
	if err != nil {
		return nil, err
	}

	identity.Spec.HashedPassword = maskPassword
	return identity, nil
}

// UpdateLocalIdentityStatus to update a existing user status.
func (s *Service) UpdateLocalIdentityStatus(identityUpdate *types.LocalIdentity) (*types.LocalIdentity, error) {
	_, err := s.store.Get(identityUpdate.Spec.TenantID, identityUpdate.Name)
	if err != nil {
		return nil, err
	}

	updater := func(old types.LocalIdentity) (types.LocalIdentity, error) {
		old.Status.Locked = identityUpdate.Status.Locked
		return old, nil
	}

	err = s.store.Update(identityUpdate.Spec.TenantID, identityUpdate.Name, updater)
	if err != nil {
		return nil, err
	}

	identity, err := s.store.Get(identityUpdate.Spec.TenantID, identityUpdate.Name)
	if err != nil {
		return nil, err
	}

	identity.Spec.HashedPassword = maskPassword
	return identity, nil
}

// UpdateLocalIdentityPassword to update a existing user password.
func (s *Service) UpdateLocalIdentityPassword(identityUpdate *types.LocalIdentity) (*types.LocalIdentity, error) {
	localIdentity, err := s.store.Get(identityUpdate.Spec.TenantID, identityUpdate.Name)
	if err != nil {
		return nil, err
	}
	updater := func(current types.LocalIdentity) (types.LocalIdentity, error) {
		err := util.VerifyDecodedPassword(identityUpdate.Spec.OriginalPassword, localIdentity.Spec.HashedPassword)
		if err != nil {
			log.Error("Invalid original password", log.String("original password", identityUpdate.Spec.OriginalPassword), log.Err(err))
			return types.LocalIdentity{}, fmt.Errorf("verify original password failed: %v", err)
		}

		if len(identityUpdate.Spec.HashedPassword) != 0 {
			current.Spec.HashedPassword, err = util.BcryptPassword(identityUpdate.Spec.HashedPassword)
			if err != nil {
				return types.LocalIdentity{}, fmt.Errorf("bcrypt password failed: %v", err)
			}
		}
		current.UpdateAt = time.Now()
		current.CreateAt = localIdentity.CreateAt

		return current, nil
	}

	err = s.store.Update(identityUpdate.Spec.TenantID, identityUpdate.Name, updater)
	if err != nil {
		return nil, err
	}

	identity, err := s.store.Get(identityUpdate.Spec.TenantID, identityUpdate.Name)
	if err != nil {
		return nil, err
	}

	identity.Spec.HashedPassword = maskPassword
	return identity, nil
}

// DeleteLocalIdentity to delete a user by given name.
func (s *Service) DeleteLocalIdentity(tenantID, name string) error {
	if err := validateIdentityName(name); err != nil {
		return err
	}
	_ = s.policyEnforcer.RemoveAllPermsForUser(tenantID, name)

	return s.store.Delete(tenantID, name)
}

// GetLocalIdentity to get a existing user by given name.
func (s *Service) GetLocalIdentity(tenantID, name string) (*types.LocalIdentity, error) {
	if err := validateIdentityName(name); err != nil {
		return nil, err
	}

	identity, err := s.store.Get(tenantID, name)
	if err != nil {
		return nil, err
	}

	identity.Spec.HashedPassword = maskPassword
	return identity, nil
}

// ListLocalIdentity to get all existing user.
func (s *Service) ListLocalIdentity(tenantID string) (*types.LocalIdentityList, error) {
	identityList := &types.LocalIdentityList{Items: []*types.LocalIdentity{}}
	identityList, err := s.store.List(tenantID)
	if err != nil {
		return nil, err
	}

	for _, ident := range identityList.Items {
		ident.Spec.HashedPassword = maskPassword
	}
	return identityList, nil
}

// ListUserPolicies to get all policies of user related to.
func (s *Service) ListUserPolicies(tenantID string, userName string) (*types.PolicyList, error) {
	policyIDs, err := s.policyEnforcer.ListUserPolicies(tenantID, userName)
	if err != nil {
		return nil, err
	}
	policyList := &types.PolicyList{Items: []*types.Policy{}}
	for _, id := range policyIDs {
		pol, err := s.policyStore.Get(tenantID, id)
		if err != nil {
			log.Error("Get policy failed", log.String("policy", id), log.Err(err))
			continue
		}
		policyList.Items = append(policyList.Items, pol)
	}

	return policyList, nil
}

// ListUserRoles to get all roles of user related to.
func (s *Service) ListUserRoles(tenantID string, userName string) (*types.RoleList, error) {
	roleIDs, err := s.policyEnforcer.ListUserRoles(tenantID, userName)
	if err != nil {
		return nil, err
	}
	roleList := &types.RoleList{Items: []*types.Role{}}
	for _, id := range roleIDs {
		rl, err := s.roleStore.Get(tenantID, id)
		if err != nil {
			log.Error("Get role failed", log.String("role", id), log.Err(err))
			continue
		}

		// not return policies
		rl.Policies = []*types.PolicyMeta{}
		roleList.Items = append(roleList.Items, rl)
	}

	return roleList, nil
}

// ListUserPerms to get all permissions allowed in roles and rules of user related to.
func (s *Service) ListUserPerms(tenantID string, userName string) (*types.Permission, error) {
	policyIDs, err := s.policyEnforcer.ListUserPolicies(tenantID, userName)
	if err != nil {
		return nil, err
	}
	policyIDSet := sets.NewString(policyIDs...)

	// get policies of the roles user related to
	roleIDs, err := s.policyEnforcer.ListUserRoles(tenantID, userName)
	if err != nil {
		return nil, err
	}
	for _, id := range roleIDs {
		policyIDs, err := s.policyEnforcer.ListRolePolicies(tenantID, id)
		if err != nil {
			log.Error("Get role policy failed", log.String("role", id), log.Err(err))
			continue
		}
		for _, policyID := range policyIDs {
			policyIDSet.Insert(policyID)
		}
	}

	policyList := &types.PolicyList{}
	for _, id := range policyIDSet.List() {
		pol, err := s.policyStore.Get(tenantID, id)
		if err != nil {
			log.Error("Get policy failed", log.String("policy", id), log.Err(err))
			continue
		}
		policyList.Items = append(policyList.Items, pol)
	}
	permits := &types.Permission{}
	permits.AllowPerms, permits.DenyPerms = mergePoliciesStatement(policyList)

	return permits, nil
}

func validateLocalIdentityCreate(identityCreate *types.LocalIdentity) error {
	allErrs := field.ErrorList{}

	if _, ok := reservedNames[identityCreate.Name]; ok {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), identityCreate.Name, "name is reserved"))
	}

	if err := validation.IsDNS1123Name(identityCreate.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), identityCreate.Name, err.Error()))
	}

	if len(identityCreate.Spec.HashedPassword) == 0 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), identityCreate.Spec.HashedPassword, "hashPassword must be specified"))
	} else {
		_, err := util.BcryptPassword(identityCreate.Spec.HashedPassword)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), identityCreate.Spec.HashedPassword, err.Error()))
		}
	}

	if err := validation.IsDisplayName(identityCreate.Spec.Extra["displayName"]); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "extra", "displayName"), identityCreate.Spec.Extra["displayName"], err.Error()))
	}

	if email, ok := identityCreate.Spec.Extra["email"]; ok && len(email) != 0 {
		if err := validation.IsEmail(email); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "extra", "email"), identityCreate.Spec.Extra["email"], err.Error()))
		}
	}

	if phoneNumber, ok := identityCreate.Spec.Extra["phoneNumber"]; ok && len(phoneNumber) != 0 {
		if err := validation.IsPhoneNumber(phoneNumber); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "extra", "phoneNumber"), identityCreate.Spec.Extra["phoneNumber"], err.Error()))
		}
	}

	return allErrs.ToAggregate()
}

func validateLocalIdentityUpdate(identityUpdate *types.LocalIdentity) error {
	allErrs := field.ErrorList{}

	if _, ok := reservedNames[identityUpdate.Name]; ok {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), identityUpdate.Name, `name is reserved`))
	}
	if err := validation.IsDNS1123Name(identityUpdate.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), identityUpdate.Name, err.Error()))
	}

	if len(identityUpdate.Spec.HashedPassword) != 0 {
		_, err := util.BcryptPassword(identityUpdate.Spec.HashedPassword)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), identityUpdate.Spec.HashedPassword, err.Error()))
		}
	}

	if displayName, ok := identityUpdate.Spec.Extra["displayName"]; ok && len(displayName) != 0 {
		if err := validation.IsDisplayName(displayName); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "extra", "displayName"), identityUpdate.Spec.Extra["displayName"], err.Error()))
		}
	}

	if email, ok := identityUpdate.Spec.Extra["email"]; ok && len(email) != 0 {
		if err := validation.IsEmail(email); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "extra", "email"), identityUpdate.Spec.Extra["email"], err.Error()))
		}
	}

	if phoneNumber, ok := identityUpdate.Spec.Extra["phoneNumber"]; ok && len(phoneNumber) != 0 {
		if err := validation.IsPhoneNumber(phoneNumber); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "extra", "phoneNumber"), identityUpdate.Spec.Extra["phoneNumber"], err.Error()))
		}
	}

	return allErrs.ToAggregate()
}

func validateIdentityName(name string) error {
	allErrs := field.ErrorList{}

	if err := validation.IsDNS1123Name(name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), name, err.Error()))
	}
	return allErrs.ToAggregate()
}

// mergePoliciesStatement merges policies statement and returns which actions can be allowed on resources.
func mergePoliciesStatement(policies *types.PolicyList) (map[string][]string, map[string][]string) {
	allowSets := map[string]sets.String{}
	denySets := map[string]sets.String{}
	for _, pol := range policies.Items {
		if pol.Statement.Effect == policyService.EffectAllow {
			for _, act := range pol.Statement.Actions {
				if _, ok := allowSets[act]; ok {
					allowSets[act].Insert(pol.Statement.Resource)
				} else {
					allowSets[act] = sets.NewString(pol.Statement.Resource)
				}
			}
		} else {
			for _, act := range pol.Statement.Actions {
				if _, ok := denySets[act]; ok {
					denySets[act].Insert(pol.Statement.Resource)
				} else {
					denySets[act] = sets.NewString(pol.Statement.Resource)
				}
			}
		}
	}

	allowList := map[string][]string{}
	for act, resSet := range allowSets {
		allowList[act] = resSet.List()
	}

	denyList := map[string][]string{}
	for act, resSet := range denySets {
		denyList[act] = resSet.List()
	}

	return allowList, denyList
}
