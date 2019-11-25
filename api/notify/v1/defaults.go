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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ChannelStatus(obj *ChannelStatus) {
	if obj.Phase == "" {
		obj.Phase = ChannelActived
	}
}

func SetDefaults_TemplateSpec(obj *TemplateSpec) {
	if obj.Keys == nil {
		obj.Keys = []string{}
	}
}

func SetDefaults_ReceiverSpec(obj *ReceiverSpec) {
	if obj.Identities == nil {
		obj.Identities = make(map[ReceiverChannel]string)
	}
}

func SetDefaults_ReceiverGroupSpec(obj *ReceiverGroupSpec) {
	if obj.Receivers == nil {
		obj.Receivers = []string{}
	}
}

func SetDefaults_MessageRequestSpec(obj *MessageRequestSpec) {
	if obj.Variables == nil {
		obj.Variables = make(map[string]string)
	}
	if obj.Receivers == nil {
		obj.Receivers = []string{}
	}
	if obj.ReceiverGroups == nil {
		obj.ReceiverGroups = []string{}
	}
}

func SetDefaults_MessageRequestStatus(obj *MessageRequestStatus) {
	if obj.Phase == "" {
		obj.Phase = MessageRequestPending
	}
	if obj.Errors == nil {
		obj.Errors = make(map[string]string)
	}
}

func SetDefaults_MessageStatus(obj *MessageStatus) {
	if obj.Phase == "" {
		obj.Phase = MessageUnread
	}
}

func SetDefaults_ConfigMap(obj *ConfigMap) {
	if obj.Data == nil {
		obj.Data = make(map[string]string)
	}
}
