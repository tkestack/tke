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

package receiver

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/notify"
)

// ValidateReceiverName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateReceiverName = apimachineryvalidation.NameIsDNSLabel

// ValidateReceiver tests if required fields in the receiver are set.
func ValidateReceiver(receiver *notify.Receiver) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&receiver.ObjectMeta, false, ValidateReceiverName, field.NewPath("metadata"))

	if receiver.Spec.DisplayName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "displayName"), "must specify display name"))
	}

	for rc, value := range receiver.Spec.Identities {
		if !IsStandardReceiverChannel(string(rc)) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "identities").Key(string(rc)), value, "must be a standard channel receiver"))
		}
	}

	return allErrs
}

// ValidateReceiverUpdate tests if required fields in the receiver are set during
// an update.
func ValidateReceiverUpdate(receiver *notify.Receiver, old *notify.Receiver) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&receiver.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateReceiver(receiver)...)

	if receiver.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}
	return allErrs
}

var standardReceiverChannel = sets.NewString(
	string(notify.ReceiverChannelEmail),
	string(notify.ReceiverChannelMobile),
	string(notify.ReceiverChannelWechatOpenID),
)

// IsStandardReceiverChannel returns true if the receiver channel is known to
// the system.
func IsStandardReceiverChannel(str string) bool {
	return standardReceiverChannel.Has(str)
}
