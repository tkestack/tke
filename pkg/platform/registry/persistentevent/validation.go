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

package persistentevent

import (
	"reflect"

	apiMachineryValidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/platform"
)

// ValidateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateName = apiMachineryValidation.ValidateNamespaceName

// ValidatePersistentEvent tests if required fields in the cluster are set.
func ValidatePersistentEvent(persistentEvent *platform.PersistentEvent) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMeta(&persistentEvent.ObjectMeta, false, ValidateName, field.NewPath("metadata"))

	if len(persistentEvent.Spec.ClusterName) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "clusterName"), "must specify a cluster name"))
	}

	v := reflect.ValueOf(persistentEvent.Spec.PersistentBackEnd)
	count := v.NumField()
	var backendCount int
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if !f.IsNil() {
			backendCount++
		}
	}
	if backendCount != 1 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "persistentBackend"), "must specify one backend"))
	}

	if persistentEvent.Spec.PersistentBackEnd.CLS != nil {
		if persistentEvent.Spec.PersistentBackEnd.CLS.LogSetID == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "persistentBackend", "CLS", "logSetID"), "must specify CLS logSetID"))
		}

		if persistentEvent.Spec.PersistentBackEnd.CLS.TopicID == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "persistentBackend", "CLS", "topicID"), "must specify CLS topicID"))
		}
	}

	if persistentEvent.Spec.PersistentBackEnd.ES != nil {
		if persistentEvent.Spec.PersistentBackEnd.ES.IndexName == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "persistentBackend", "ES", "indexName"), "must specify ES indexName"))
		}

		if persistentEvent.Spec.PersistentBackEnd.ES.IP == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "persistentBackend", "ES", "ip"), "must specify ES IP"))
		}

		if persistentEvent.Spec.PersistentBackEnd.ES.Scheme == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "persistentBackend", "ES", "scheme"), "must specify ES scheme (http or https)"))
		}
	}

	return allErrs
}

// ValidatePersistentEventUpdate tests if required fields in the namespace set are
// set during an update.
func ValidatePersistentEventUpdate(persistentEvent *platform.PersistentEvent, old *platform.PersistentEvent) field.ErrorList {
	allErrs := apiMachineryValidation.ValidateObjectMetaUpdate(&persistentEvent.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidatePersistentEvent(persistentEvent)...)

	if persistentEvent.Spec.ClusterName != old.Spec.ClusterName {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "clusterName"), persistentEvent.Spec.ClusterName, "disallowed change the cluster name"))
	}

	if persistentEvent.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "tenantID"), persistentEvent.Spec.TenantID, "disallowed change the tenant"))
	}

	return allErrs
}
