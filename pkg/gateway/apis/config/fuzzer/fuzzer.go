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

package fuzzer

import (
	"github.com/google/gofuzz"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
)

// Funcs returns the fuzzer functions for the gatewayconfig apis.
func Funcs(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		// provide non-empty values for fields with defaults, so the defaulter doesn't change values during round-trip
		func(obj *gatewayconfig.GatewayConfiguration, c fuzz.Continue) {
			c.FuzzNoCustom(obj)
			obj.DisableOIDCProxy = true
			obj.Components.Platform = &gatewayconfig.Component{
				FrontProxy: &gatewayconfig.FrontProxyComponent{
					UsernameHeader:    "X-Remote-User",
					GroupsHeader:      "X-Remote-Groups",
					ExtraPrefixHeader: "X-Remote-Extra-",
				},
			}
			obj.Components.Business = &gatewayconfig.Component{
				FrontProxy: &gatewayconfig.FrontProxyComponent{
					UsernameHeader:    "X-Remote-User",
					GroupsHeader:      "X-Remote-Groups",
					ExtraPrefixHeader: "X-Remote-Extra-",
				},
			}
			obj.Components.Auth = &gatewayconfig.Component{
				FrontProxy: &gatewayconfig.FrontProxyComponent{
					UsernameHeader:    "X-Remote-User",
					GroupsHeader:      "X-Remote-Groups",
					ExtraPrefixHeader: "X-Remote-Extra-",
				},
			}
			obj.Components.Monitor = &gatewayconfig.Component{
				FrontProxy: &gatewayconfig.FrontProxyComponent{
					UsernameHeader:    "X-Remote-User",
					GroupsHeader:      "X-Remote-Groups",
					ExtraPrefixHeader: "X-Remote-Extra-",
				},
			}
			obj.Components.Notify = &gatewayconfig.Component{
				FrontProxy: &gatewayconfig.FrontProxyComponent{
					UsernameHeader:    "X-Remote-User",
					GroupsHeader:      "X-Remote-Groups",
					ExtraPrefixHeader: "X-Remote-Extra-",
				},
			}
			obj.Components.Registry = &gatewayconfig.Component{
				FrontProxy: &gatewayconfig.FrontProxyComponent{
					UsernameHeader:    "X-Remote-User",
					GroupsHeader:      "X-Remote-Groups",
					ExtraPrefixHeader: "X-Remote-Extra-",
				},
			}
		},
	}
}
