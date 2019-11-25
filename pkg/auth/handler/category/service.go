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

package category

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"time"
	"tkestack.io/tke/pkg/auth/registry"
	"tkestack.io/tke/pkg/auth/registry/category"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/validation"
)

// Service is responsible for performing policy category and actions operation onto the storage backend.
type Service struct {
	store *category.Storage
}

// NewCategoryService creates a new category service object
func NewCategoryService(registry *registry.Registry) *Service {
	return &Service{store: registry.CategoryStorage()}
}

// CreateCategory creates a new policy action category and returns it.
func (s *Service) CreateCategory(categoryCreate *types.Category) (*types.Category, error) {
	categoryCreate.CreateAt = time.Now()
	categoryCreate.UpdateAt = categoryCreate.CreateAt

	err := s.store.Create(categoryCreate)
	if err != nil {
		return nil, err
	}
	return categoryCreate, nil
}

// ListCategory gets all existing categories for specify tenant.
func (s *Service) ListCategory(tenantID string) (*types.CategoryList, error) {
	categoryList, err := s.store.List(tenantID)
	if err != nil {
		return nil, err
	}

	return categoryList, nil
}

// GetCategory gets a category by a given name.
func (s *Service) GetCategory(tenantID string, name string) (*types.Category, error) {
	cat, err := s.store.Get(tenantID, name)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

// DeleteCategory deletes a category by a given name.
func (s *Service) DeleteCategory(tenantID string, name string) error {
	err := s.store.Delete(tenantID, name)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCategory updates a category metadata, include name and description, not actions.
func (s *Service) UpdateCategory(categoryUpdate *types.Category) (*types.Category, error) {
	_, err := s.store.Get(categoryUpdate.TenantID, categoryUpdate.Name)
	if err != nil {
		return nil, err
	}

	updater := func(current types.Category) (types.Category, error) {
		current.Name = categoryUpdate.Name
		current.DisplayName = categoryUpdate.DisplayName
		current.Description = categoryUpdate.Description

		current.UpdateAt = time.Now()
		return current, nil
	}

	err = s.store.Update(categoryUpdate.TenantID, categoryUpdate.Name, updater)
	if err != nil {
		return nil, err
	}

	categoryUpdated, err := s.store.Get(categoryUpdate.TenantID, categoryUpdate.Name)
	if err != nil {
		return nil, err
	}

	return categoryUpdated, nil
}

// DeleteActions deletes actions in a category.
func (s *Service) DeleteActions(tenantID string, name string, categoryUpdate *types.Category) (*types.Category, error) {
	_, err := s.store.Get(tenantID, name)
	if err != nil {
		return nil, err
	}

	updater := func(current types.Category) (types.Category, error) {
		for act := range categoryUpdate.Actions {
			delete(current.Actions, act)
		}

		current.UpdateAt = time.Now()
		return current, nil
	}

	err = s.store.Update(tenantID, name, updater)
	if err != nil {
		return nil, err
	}

	categoryUpdated, err := s.store.Get(tenantID, name)
	if err != nil {
		return nil, err
	}

	return categoryUpdated, nil
}

// AddActions adds actions in a category.
func (s *Service) AddActions(tenantID string, name string, categoryUpdate *types.Category) (*types.Category, error) {
	_, err := s.store.Get(tenantID, name)
	if err != nil {
		return nil, err
	}

	updater := func(current types.Category) (types.Category, error) {
		for act, desc := range categoryUpdate.Actions {
			current.Actions[act] = desc
		}

		current.UpdateAt = time.Now()
		return current, nil
	}

	err = s.store.Update(tenantID, name, updater)
	if err != nil {
		return nil, err
	}

	categoryUpdated, err := s.store.Get(tenantID, name)
	if err != nil {
		return nil, err
	}

	return categoryUpdated, nil
}

func validateCategoryCreate(categoryCreate *types.Category) error {
	allErrs := field.ErrorList{}
	if err := validation.IsDNS1123Name(categoryCreate.Name); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("name"), categoryCreate.Name, err.Error()))
	}
	if err := validation.IsDisplayName(categoryCreate.DisplayName); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("displayName"), categoryCreate.DisplayName, err.Error()))
	}

	return allErrs.ToAggregate()
}
