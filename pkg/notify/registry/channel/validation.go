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

package channel

import (
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"tkestack.io/tke/api/notify"
)

// ValidateChannelName is a ValidateNameFunc for names that must be a DNS
// subdomain.
var ValidateChannelName = apimachineryvalidation.NameIsDNSLabel

// ValidateChannel tests if required fields in the channel are set.
func ValidateChannel(channel *notify.Channel) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(&channel.ObjectMeta, false, ValidateChannelName, field.NewPath("metadata"))

	if channel.Spec.DisplayName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec", "displayName"), "must specify display name"))
	}

	channelCount := 0
	if channel.Spec.TencentCloudSMS != nil {
		channelCount++

		if channel.Spec.TencentCloudSMS.AppKey == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "tencentCloudSMS", "appKey"), "must specify appKey of tencent cloud sms gateway"))
		}
	}

	if channel.Spec.Wechat != nil {
		channelCount++

		if channel.Spec.Wechat.AppID == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "wechat", "appID"), "must specify appID of wechat public platform"))
		}

		if channel.Spec.Wechat.AppSecret == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "wechat", "appSecret"), "must specify appSecret of wechat public platform"))
		}
	}

	if channel.Spec.SMTP != nil {
		channelCount++

		if channel.Spec.SMTP.Email == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "smtp", "email"), "must specify email of smtp sender"))
		}

		if channel.Spec.SMTP.Password == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "smtp", "password"), "must specify password of smtp sender"))
		}

		if channel.Spec.SMTP.SMTPHost == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "smtp", "smtpHost"), "must specify host name of smtp sender"))
		}

		if channel.Spec.SMTP.SMTPPort == 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "smtp", "smtpPort"), "must specify host port of smtp sender"))
		}
	}

	if channelCount == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec"), "must specify one of channel type: `tencentCloudSMS`, `wechat` or `smtp`"))
	} else if channelCount > 1 {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "may not specify more than 1 channel type: `tencentCloudSMS`, `wechat` or `smtp`"))
	}

	return allErrs
}

// ValidateChannelUpdate tests if required fields in the channel are set during
// an update.
func ValidateChannelUpdate(channel *notify.Channel, old *notify.Channel) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(&channel.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateChannel(channel)...)

	if channel.Spec.TenantID != old.Spec.TenantID {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "tenantID"), "disallowed change the tenant"))
	}
	return allErrs
}
