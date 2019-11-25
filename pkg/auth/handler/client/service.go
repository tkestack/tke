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

package client

import (
	"fmt"

	"tkestack.io/tke/pkg/auth/util"

	"github.com/dexidp/dex/storage"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/secret"
	"tkestack.io/tke/pkg/util/validation"
)

// Service is responsible for performing client crud actions onto the storage backend.
type Service struct {
	dexStorage storage.Storage
}

// NewClientService create a new client service object
func NewClientService(store storage.Storage) *Service {
	return &Service{dexStorage: store}
}

// CreateClient to create a new OAuth2 client.
func (s *Service) CreateClient(clientCreate *types.Client) (*types.Client, error) {
	if clientCreate.Secret == "" {
		clientCreate.Secret = secret.CreateRandomPassword(32)
	}

	err := s.dexStorage.CreateClient(toDexClient(clientCreate))
	if err != nil {
		return nil, err
	}
	return clientCreate, nil
}

// GetClient get a existing OAuth2 client by given id.
func (s *Service) GetClient(id string) (*types.Client, error) {
	client, err := s.dexStorage.GetClient(id)
	if err != nil {
		return nil, err
	}

	return fromDexClient(client), nil
}

// DeleteClient delete a existing OAuth2 client by given id.
func (s *Service) DeleteClient(id string) error {
	err := s.dexStorage.DeleteClient(id)
	if err != nil {
		return err
	}

	return nil
}

// ListClient get a list of existing OAuth2 clients.
func (s *Service) ListClient(id, keyword string) (*types.ClientList, error) {
	clientList := types.ClientList{}
	if len(id) != 0 {
		cli, err := s.dexStorage.GetClient(id)
		if err != nil {
			return nil, err
		}
		clientList.Items = append(clientList.Items, fromDexClient(cli))
		return &clientList, nil
	}

	dexClients, err := s.dexStorage.ListClients()
	if err != nil {
		return nil, err
	}

	for _, cli := range dexClients {
		if len(keyword) == 0 {
			clientList.Items = append(clientList.Items, fromDexClient(cli))
			continue
		}

		if util.CaseInsensitiveContains(cli.Name, keyword) || util.CaseInsensitiveContains(cli.ID, keyword) {
			clientList.Items = append(clientList.Items, fromDexClient(cli))
		}
	}

	return &clientList, nil
}

// UpdateClient update a existing OAuth2 client.
func (s *Service) UpdateClient(clientUpdate *types.Client) (*types.Client, error) {
	_, err := s.dexStorage.GetClient(clientUpdate.ID)
	if err != nil {
		return nil, err
	}

	updater := func(current storage.Client) (storage.Client, error) {
		current = toDexClient(clientUpdate)
		return current, nil
	}

	err = s.dexStorage.UpdateClient(clientUpdate.ID, updater)
	if err != nil {
		return nil, err
	}

	return clientUpdate, nil
}

func validateClientCreate(clientCreate *types.Client) error {
	allErrs := field.ErrorList{}
	if err := validation.IsDNS1123Name(clientCreate.ID); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("id"), clientCreate.ID, err.Error()))
	}
	if clientCreate.Secret != "" {
		if err := validation.IsDisplayName(clientCreate.Secret); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("secret"), clientCreate.Secret, "length must be less than 64"))
		}
	}
	for i, uri := range clientCreate.RedirectUris {
		if err := validation.IsURL(uri); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath(fmt.Sprintf("redirect_uris[%d]", i)), uri, err.Error()))
		}
	}
	if err := validation.IsDisplayName(clientCreate.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), clientCreate.Name, err.Error()))
	}
	if clientCreate.LogoURL != "" {
		if err := validation.IsURL(clientCreate.LogoURL); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("logo_url"), clientCreate.LogoURL, err.Error()))
		}
	}
	return allErrs.ToAggregate()
}

func toDexClient(client *types.Client) storage.Client {
	return storage.Client{
		ID:           client.ID,
		Secret:       client.Secret,
		RedirectURIs: client.RedirectUris,
		TrustedPeers: client.TrustedPeers,
		Public:       client.Public,
		Name:         client.Name,
		LogoURL:      client.LogoURL,
	}
}

func fromDexClient(client storage.Client) *types.Client {
	return &types.Client{
		ID:           client.ID,
		Secret:       client.Secret,
		RedirectUris: client.RedirectURIs,
		TrustedPeers: client.TrustedPeers,
		Public:       client.Public,
		Name:         client.Name,
		LogoURL:      client.LogoURL,
	}
}
