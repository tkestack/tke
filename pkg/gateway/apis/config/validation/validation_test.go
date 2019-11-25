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

package validation

import (
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"testing"
	gatewayconfig "tkestack.io/tke/pkg/gateway/apis/config"
)

func TestValidateGatewayConfiguration(t *testing.T) {
	successCase := &gatewayconfig.GatewayConfiguration{
		DisableOIDCProxy: true,
		Components: gatewayconfig.Components{
			Platform: &gatewayconfig.Component{
				Address:     "https://127.0.0.1:8080",
				Passthrough: &gatewayconfig.PassthroughComponent{},
			},
			Business: &gatewayconfig.Component{
				Address: "https://127.0.0.1:8080",
				FrontProxy: &gatewayconfig.FrontProxyComponent{
					CAFile:            "/fake.ca.pem",
					ClientCertFile:    "/fake.client.ca.pem",
					ClientKeyFile:     "/fake.client.key.pem",
					UsernameHeader:    "X-Remote-User",
					GroupsHeader:      "X-Remote-Groups",
					ExtraPrefixHeader: "X-Remote-Extra-",
				},
			},
		},
	}
	if allErrors := ValidateGatewayConfiguration(successCase); allErrors != nil {
		t.Errorf("expect no errors, got %v", allErrors)
	}

	errorCase := &gatewayconfig.GatewayConfiguration{
		Components: gatewayconfig.Components{
			Platform: &gatewayconfig.Component{
				Passthrough: &gatewayconfig.PassthroughComponent{},
			},
			Business: &gatewayconfig.Component{
				Address: "https://127.0.0.1:8080",
				FrontProxy: &gatewayconfig.FrontProxyComponent{
					CAFile: "/fake.ca.pem",
				},
			},
		},
	}
	const numErrs = 6
	if allErrors := ValidateGatewayConfiguration(errorCase); len(allErrors.(utilerrors.Aggregate).Errors()) != numErrs {
		t.Errorf("expect %d errors, got %v", numErrs, len(allErrors.(utilerrors.Aggregate).Errors()))
	}
}
