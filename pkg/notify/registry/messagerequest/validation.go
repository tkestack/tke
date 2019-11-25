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

package messagerequest

import (
	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	notifyinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/notify/internalversion"
	"tkestack.io/tke/api/notify"
)

// ValidateMessageRequestName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateMessageRequestName = apimachineryvalidation.NameIsDNSLabel

// ValidateMessageRequest tests if required fields in the messageRequest are set.
func ValidateMessageRequest(messageRequest *notify.MessageRequest, notifyClient *notifyinternalclient.NotifyClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&messageRequest.ObjectMeta, true, ValidateMessageRequestName, field.NewPath("metadata"))

	channel, err := notifyClient.Channels().Get(messageRequest.ObjectMeta.Namespace, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		allErrs = append(allErrs, field.NotFound(field.NewPath("metadata", "namespace"), messageRequest.ObjectMeta.Namespace))
	} else if err != nil {
		allErrs = append(allErrs, field.InternalError(field.NewPath(""), err))
	} else {
		if messageRequest.Spec.TenantID != "" && channel.Spec.TenantID != messageRequest.Spec.TenantID {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "no authorized to create message request in this channel"))
		} else {
			if messageRequest.Spec.TemplateName == "" {
				allErrs = append(allErrs, field.Required(field.NewPath("spec", "templateName"), "must specify template name"))
			} else {
				template, err := notifyClient.Templates(channel.ObjectMeta.Name).Get(messageRequest.Spec.TemplateName, metav1.GetOptions{})
				if err != nil && errors.IsNotFound(err) {
					allErrs = append(allErrs, field.NotFound(field.NewPath("spec", "templateName"), messageRequest.Spec.TemplateName))
				} else if err != nil {
					allErrs = append(allErrs, field.InternalError(field.NewPath("spec", "templateName"), err))
				} else {
					for _, key := range template.Spec.Keys {
						if _, ok := messageRequest.Spec.Variables[key]; !ok {
							allErrs = append(allErrs, field.Required(field.NewPath("spec", "variables").Key(key), "must specify variables"))
						}
					}
				}
			}
		}
	}

	return allErrs
}

// ValidateMessageRequestUpdate tests if required fields in the messageRequest are set during
// an update.
func ValidateMessageRequestUpdate(messageRequest *notify.MessageRequest, old *notify.MessageRequest, notifyClient *notifyinternalclient.NotifyClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&messageRequest.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateMessageRequest(messageRequest, notifyClient)...)

	if messageRequest.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}
	return allErrs
}
