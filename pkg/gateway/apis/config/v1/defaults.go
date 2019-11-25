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

import "k8s.io/apimachinery/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_GatewayConfiguration(obj *GatewayConfiguration) {
	if obj.Components.Platform != nil && obj.Components.Platform.FrontProxy != nil {
		defaultGatewayConfigurationComponent(obj.Components.Platform.FrontProxy)
	}
	if obj.Components.Business != nil && obj.Components.Business.FrontProxy != nil {
		defaultGatewayConfigurationComponent(obj.Components.Business.FrontProxy)
	}
	if obj.Components.Notify != nil && obj.Components.Notify.FrontProxy != nil {
		defaultGatewayConfigurationComponent(obj.Components.Notify.FrontProxy)
	}
	if obj.Components.Monitor != nil && obj.Components.Monitor.FrontProxy != nil {
		defaultGatewayConfigurationComponent(obj.Components.Monitor.FrontProxy)
	}
	if obj.Components.Registry != nil && obj.Components.Registry.FrontProxy != nil {
		defaultGatewayConfigurationComponent(obj.Components.Registry.FrontProxy)
	}
	if obj.Components.Auth != nil && obj.Components.Auth.FrontProxy != nil {
		defaultGatewayConfigurationComponent(obj.Components.Auth.FrontProxy)
	}
}

func defaultGatewayConfigurationComponent(obj *FrontProxyComponent) {
	if obj == nil {
		return
	}
	obj.UsernameHeader = "X-Remote-User"
	obj.GroupsHeader = "X-Remote-Groups"
	obj.ExtraPrefixHeader = "X-Remote-Extra-"
}
