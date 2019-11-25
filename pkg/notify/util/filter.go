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

package util

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"tkestack.io/tke/api/notify"
	"tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
)

// FilterChannel is used to filter channels that do not belong to the tenant.
func FilterChannel(ctx context.Context, channel *notify.Channel) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if channel.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("channel"), channel.ObjectMeta.Name)
	}
	return nil
}

// FilterReceiver is used to filter receiver that do not belong to the tenant.
func FilterReceiver(ctx context.Context, receiver *notify.Receiver) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if receiver.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("receiver"), receiver.ObjectMeta.Name)
	}
	return nil
}

// FilterReceiverGroup is used to filter receiver group that do not belong to the tenant.
func FilterReceiverGroup(ctx context.Context, receiverGroup *notify.ReceiverGroup) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if receiverGroup.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("receivergroup"), receiverGroup.ObjectMeta.Name)
	}
	return nil
}

// FilterMessage is used to filter message that do not belong to the tenant.
func FilterMessage(ctx context.Context, message *notify.Message) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if message.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("message"), message.ObjectMeta.Name)
	}
	return nil
}

// FilterMessageRequest is used to filter message request that do not belong to the tenant.
func FilterMessageRequest(ctx context.Context, messageRequest *notify.MessageRequest) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if messageRequest.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("message"), messageRequest.ObjectMeta.Name)
	}
	return nil
}

// FilterTemplate is used to filter template that do not belong to the tenant.
func FilterTemplate(ctx context.Context, tpl *notify.Template) error {
	_, tenantID := authentication.GetUsernameAndTenantID(ctx)
	if tenantID == "" {
		return nil
	}
	if tpl.Spec.TenantID != tenantID {
		return errors.NewNotFound(v1.Resource("template"), tpl.ObjectMeta.Name)
	}
	return nil
}
