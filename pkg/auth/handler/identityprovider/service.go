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

package identityprovider

import (
	"fmt"

	"github.com/dexidp/dex/server"
	"github.com/dexidp/dex/storage"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/pkg/auth/authentication/tenant"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util/validation"
)

// Service is responsible for performing OIDC idp crud actions onto the storage backend.
type Service struct {
	dexStorage storage.Storage
	helper     *tenant.Helper
}

// NewIdentidyProviderService create a new idp service object
func NewIdentidyProviderService(store storage.Storage, helper *tenant.Helper) *Service {
	return &Service{dexStorage: store, helper: helper}
}

// CreateIdentityProvier creates a new OIDC idp and returns it.
func (s *Service) CreateIdentityProvier(idpCreate *types.IdentityProvider) (*types.IdentityProvider, error) {
	err := s.dexStorage.CreateConnector(toDexConnector(idpCreate))
	if err != nil {
		return nil, err
	}

	// create tenant admin
	err = s.helper.CreateAdmin(idpCreate.ID)
	if err != nil {
		_ = s.DeleteIdentityProvider(idpCreate.ID)
		return nil, fmt.Errorf("init tenant resource failed: %v", err)
	}

	// create predefine policies and categories
	err = s.helper.LoadResource(idpCreate.ID)
	if err != nil {
		_ = s.DeleteIdentityProvider(idpCreate.ID)
		return nil, fmt.Errorf("init tenant resource failed: %v", err)
	}

	return idpCreate, nil
}

// GetIdentityProvier get a existing OIDC idp by given id..
func (s *Service) GetIdentityProvier(id string) (*types.IdentityProvider, error) {
	conn, err := s.dexStorage.GetConnector(id)
	if err != nil {
		return nil, err
	}
	return fromDexConnector(conn), nil
}

// DeleteIdentityProvider delete a existing OIDC idp by given id..
func (s *Service) DeleteIdentityProvider(id string) error {
	err := s.dexStorage.DeleteConnector(id)
	if err != nil {
		return err
	}
	return nil
}

// ListIdentityProvider get a list of existing OIDC idp.
func (s *Service) ListIdentityProvider(idpType, id, keyword string) (*types.IdentityProviderList, error) {
	idpList := types.IdentityProviderList{}

	if len(id) != 0 {
		conn, err := s.dexStorage.GetConnector(id)
		if err != nil {
			return nil, err
		}

		idpList.Items = append(idpList.Items, fromDexConnector(conn))
		return &idpList, nil
	}

	dexConns, err := s.dexStorage.ListConnectors()
	if err != nil {
		return nil, err
	}

	for _, conn := range dexConns {
		if len(idpType) == 0 && len(keyword) == 0 {
			idpList.Items = append(idpList.Items, fromDexConnector(conn))
			continue
		}

		if (len(idpType) == 0 || conn.Type == idpType) && (len(keyword) == 0 || util.CaseInsensitiveContains(conn.Name, keyword) || util.CaseInsensitiveContains(conn.ID, keyword)) {
			idpList.Items = append(idpList.Items, fromDexConnector(conn))
		}
	}

	return &idpList, nil
}

// UpdateIdentityProvider update a existing OIDC idp.
func (s *Service) UpdateIdentityProvider(idpUpdate *types.IdentityProvider) (*types.IdentityProvider, error) {
	_, err := s.dexStorage.GetConnector(idpUpdate.ID)
	if err != nil {
		return nil, err
	}

	updater := func(current storage.Connector) (storage.Connector, error) {
		current.ID = idpUpdate.ID
		current.Name = idpUpdate.Name
		current.Config = []byte(idpUpdate.Config)
		current.Type = idpUpdate.Type
		current.ResourceVersion = idpUpdate.ResourceVersion
		return current, nil
	}

	err = s.dexStorage.UpdateConnector(idpUpdate.ID, updater)
	if err != nil {
		return nil, err
	}

	return idpUpdate, nil
}

func suppportIDPTypes() []string {
	var supportTypes []string
	for key := range server.ConnectorsConfig {
		supportTypes = append(supportTypes, key)
	}

	supportTypes = append(supportTypes, server.LocalConnector)

	return supportTypes
}

// validateIdentityCreate to validate identityProvider created
func validateIdentityProviderCreate(idpCreate *types.IdentityProvider) error {
	allErrs := field.ErrorList{}

	if err := validation.IsDNS1123Name(idpCreate.ID); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("id"), idpCreate.ID, err.Error()))
	}

	if err := validation.IsDisplayName(idpCreate.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), idpCreate.Name, err.Error()))
	}

	if _, ok := server.ConnectorsConfig[idpCreate.Type]; !ok && idpCreate.Type != server.LocalConnector {
		allErrs = append(allErrs, field.Invalid(field.NewPath("type"), idpCreate.Type, fmt.Sprintf("only support %v", suppportIDPTypes())))
	}

	return allErrs.ToAggregate()
}

func toDexConnector(idp *types.IdentityProvider) storage.Connector {
	return storage.Connector{
		ID:              idp.ID,
		Name:            idp.Name,
		Type:            idp.Type,
		ResourceVersion: idp.ResourceVersion,
		Config:          []byte(idp.Config),
	}
}

func fromDexConnector(conn storage.Connector) *types.IdentityProvider {
	return &types.IdentityProvider{
		ID:              conn.ID,
		Name:            conn.Name,
		Type:            conn.Type,
		ResourceVersion: conn.ResourceVersion,
		Config:          string(conn.Config),
	}
}
