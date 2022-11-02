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

package message

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/pkg/notify/registry/receiver"
)

// ValidateMessageName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateMessageName = apimachineryvalidation.NameIsDNSLabel

// ValidateMessage tests if required fields in the message are set.
func ValidateMessage(message *notify.Message) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&message.ObjectMeta, false, ValidateMessageName, field.NewPath("metadata"))

	if message.Spec.ReceiverName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "receiverName"), "must specify receiver name"))
	}

	if message.Spec.ReceiverChannel == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "receiverChannel"), "must specify receiver channel"))
	} else if !receiver.IsStandardReceiverChannel(string(message.Spec.ReceiverChannel)) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "receiverChannel"), message.Spec.ReceiverChannel, "must be a standard channel receiver"))
	}

	if message.Spec.Identity == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "identity"), "must specify receiver identity"))
	}

	return allErrs
}

// ValidateMessageUpdate tests if required fields in the message are set during
// an update.
func ValidateMessageUpdate(message *notify.Message, old *notify.Message) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&message.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateMessage(message)...)

	if message.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}
	return allErrs
}
