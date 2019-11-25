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

package receivergroup

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	notifyinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/notify/internalversion"
	"tkestack.io/tke/api/notify"
)

// ValidateReceiverGroupName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateReceiverGroupName = apimachineryvalidation.NameIsDNSLabel

// ValidateReceiverGroup tests if required fields in the receiverGroup are set.
func ValidateReceiverGroup(receiverGroup *notify.ReceiverGroup, notifyClient *notifyinternalclient.NotifyClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&receiverGroup.ObjectMeta, false, ValidateReceiverGroupName, field.NewPath("metadata"))

	if receiverGroup.Spec.DisplayName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "displayName"), "must specify display name"))
	}

	if len(receiverGroup.Spec.Receivers) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "receivers"), "must specify a receiver"))
	} else {
		for _, receiverName := range receiverGroup.Spec.Receivers {
			receiver, err := notifyClient.Receivers().Get(receiverName, metav1.GetOptions{})
			if err != nil && errors.IsNotFound(err) {
				allErrs = append(allErrs, field.NotFound(field.NewPath("spec", "receivers").Key(receiverName), receiverName))
			} else if err != nil {
				allErrs = append(allErrs, field.InternalError(field.NewPath("spec", "receivers").Key(receiverName), err))
			} else if receiver.Spec.TenantID != receiverGroup.Spec.TenantID {
				allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "receivers").Key(receiverName), fmt.Sprintf("no authorized to add receiver %s to this group", receiverName)))
			}
		}
	}

	return allErrs
}

// ValidateReceiverGroupUpdate tests if required fields in the receiverGroup are set during
// an update.
func ValidateReceiverGroupUpdate(receiverGroup *notify.ReceiverGroup, old *notify.ReceiverGroup, notifyClient *notifyinternalclient.NotifyClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&receiverGroup.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateReceiverGroup(receiverGroup, notifyClient)...)

	if receiverGroup.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}
	return allErrs
}
