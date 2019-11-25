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

package template

import (
	"k8s.io/apimachinery/pkg/api/errors"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	notifyinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/notify/internalversion"
	"tkestack.io/tke/api/notify"
)

// ValidateTemplateName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateTemplateName = apimachineryvalidation.NameIsDNSLabel

// ValidateTemplate tests if required fields in the template are set.
func ValidateTemplate(template *notify.Template, notifyClient *notifyinternalclient.NotifyClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&template.ObjectMeta, true, ValidateTemplateName, field.NewPath("metadata"))

	if template.Spec.DisplayName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("displayName"), "must specify display name"))
	}

	channel, err := notifyClient.Channels().Get(template.ObjectMeta.Namespace, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		allErrs = append(allErrs, field.NotFound(field.NewPath("metadata", "namespace"), template.ObjectMeta.Namespace))
	} else if err != nil {
		allErrs = append(allErrs, field.InternalError(field.NewPath(""), err))
	} else {
		if channel.Spec.TenantID != template.Spec.TenantID {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "no authorized to create template in this channel"))
		} else {
			if channel.Spec.TencentCloudSMS != nil {
				if template.Spec.TencentCloudSMS == nil {
					allErrs = append(allErrs, field.Required(field.NewPath("tencentCloudSMS"), "must specify tencent cloud sms template"))
				} else {
					if template.Spec.TencentCloudSMS.TemplateID == "" {
						allErrs = append(allErrs, field.Required(field.NewPath("tencentCloudSMS", "templateID"), "must specify template id of tencent cloud sms gateway"))
					}
				}
			}

			if channel.Spec.Wechat != nil {
				if template.Spec.Wechat == nil {
					allErrs = append(allErrs, field.Required(field.NewPath("wechat"), "must specify wechat template"))
				} else {
					if template.Spec.Wechat.TemplateID == "" {
						allErrs = append(allErrs, field.Required(field.NewPath("wechat", "templateID"), "must specify template id of wechat public platform"))
					}
				}
			}

			if channel.Spec.SMTP != nil {
				if template.Spec.Text == nil {
					allErrs = append(allErrs, field.Required(field.NewPath("text"), "must specify text template"))
				} else {
					if template.Spec.Text.Body == "" {
						allErrs = append(allErrs, field.Required(field.NewPath("text", "body"), "must specify body of text channel"))
					}
				}
			}
		}
	}

	return allErrs
}

// ValidateTemplateUpdate tests if required fields in the template are set during
// an update.
func ValidateTemplateUpdate(template *notify.Template, old *notify.Template, notifyClient *notifyinternalclient.NotifyClient) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&template.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateTemplate(template, notifyClient)...)

	return allErrs
}
