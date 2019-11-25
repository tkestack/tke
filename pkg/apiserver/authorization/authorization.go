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

package authorization

import (
	"fmt"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericapiserver "k8s.io/apiserver/pkg/server"
	apiserveroptions "tkestack.io/tke/pkg/apiserver/options"
)

// SetupAuthorization to setup the generic apiserver by authorization options.
func SetupAuthorization(genericAPIServerConfig *genericapiserver.Config, authorizationOpts *apiserveroptions.AuthorizationOptions) error {
	var err error
	genericAPIServerConfig.Authorization.Authorizer, genericAPIServerConfig.RuleResolver, err = buildAuthorizer(authorizationOpts)
	if err != nil {
		return fmt.Errorf("invalid authorization config: %v", err)
	}
	return nil
}

// buildAuthorizer constructs the authorizer
func buildAuthorizer(o *apiserveroptions.AuthorizationOptions) (authorizer.Authorizer, authorizer.RuleResolver, error) {
	authorizationConfig := Config{
		AuthorizationModes:          o.Modes,
		WebhookConfigFile:           o.WebhookConfigFile,
		WebhookCacheAuthorizedTTL:   o.WebhookCacheAuthorizedTTL,
		WebhookCacheUnauthorizedTTL: o.WebhookCacheUnauthorizedTTL,
	}
	return authorizationConfig.New()
}
