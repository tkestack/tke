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

package meshmanager

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/mesh"
)

// ValidateMeshManagerName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateMeshManagerName = apimachineryvalidation.NameIsDNSLabel

// ValidateMeshManager tests if required fields in the meshmanager are set.
func ValidateMeshManager(meshmanager *mesh.MeshManager) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&meshmanager.ObjectMeta, false, ValidateMeshManagerName, field.NewPath("metadata"))

	specField := field.NewPath("spec")
	dbFld := specField.Child("dataBase")
	if meshmanager.Spec.DataBase == nil {
		allErrs = append(allErrs, field.Required(dbFld, "database config must specify"))
	} else {
		if len(meshmanager.Spec.DataBase.Host) == 0 {
			allErrs = append(allErrs, field.Required(dbFld.Child("host"), "db host must specify"))
		}

		if meshmanager.Spec.DataBase.Port <= 0 || meshmanager.Spec.DataBase.Port > 65535 {
			allErrs = append(allErrs, field.Required(dbFld.Child("port"), "db port must between 1 and 65536"))
		}

		if len(meshmanager.Spec.DataBase.UserName) == 0 {
			allErrs = append(allErrs, field.Required(dbFld.Child("userName"), "db userName must specify"))
		}

		if len(meshmanager.Spec.DataBase.Password) == 0 {
			allErrs = append(allErrs, field.Required(dbFld.Child("password"), "db password must specify"))
		}

		if len(meshmanager.Spec.DataBase.DbName) == 0 {
			allErrs = append(allErrs, field.Required(dbFld.Child("dbName"), "db name must specify"))
		}
	}

	tracingField := specField.Child("tracingStorageBackend")
	if meshmanager.Spec.TracingStorageBackend == nil {
		allErrs = append(allErrs, field.Required(tracingField, "tracing backend config must specify"))
	} else {
		if len(meshmanager.Spec.TracingStorageBackend.StorageType) == 0 {
			allErrs = append(allErrs, field.Required(tracingField.Child("storageType"), "tracing storage type must specify"))
		}

		if len(meshmanager.Spec.TracingStorageBackend.StorageAddresses) == 0 {
			allErrs = append(allErrs, field.Required(tracingField.Child("storageAddress"), "at least one tracing storage address must specify"))
		}
	}

	metricField := specField.Child("metricStorageBackend")
	if meshmanager.Spec.MetricStorageBackend == nil {
		allErrs = append(allErrs, field.Required(metricField, "metric backend config must specify"))
	} else {
		if len(meshmanager.Spec.MetricStorageBackend.StorageType) == 0 {
			allErrs = append(allErrs, field.Required(metricField.Child("storageType"), "metric storage type must specify"))
		}

		if len(meshmanager.Spec.MetricStorageBackend.StorageAddresses) == 0 {
			allErrs = append(allErrs, field.Required(metricField.Child("storageAddress"), "at least one metric storage address must specify"))
		}
	}

	return allErrs
}

// ValidateMeshManagerUpdate tests if required fields in the meshmanager are set during
// an update.
func ValidateMeshManagerUpdate(meshmanager *mesh.MeshManager, old *mesh.MeshManager) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&meshmanager.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateMeshManager(meshmanager)...)

	if meshmanager.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}

	if meshmanager.Status.Phase == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("status", "phase"), string(meshmanager.Status.Phase)))
	}
	return allErrs
}
