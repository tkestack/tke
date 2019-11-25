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

package aggregation

import (
	"github.com/casbin/casbin"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/authorization/union"
	"k8s.io/apiserver/plugin/pkg/authorizer/webhook"
	"tkestack.io/tke/cmd/tke-auth/app/options"
	"tkestack.io/tke/pkg/apiserver/authorization/abac"
	"tkestack.io/tke/pkg/auth/authorization/local"
)

// NewAuthorizer creates a authorizer for subject access review and returns it.
func NewAuthorizer(authorizationOpts *options.AuthorizationOptions, authOpts *options.AuthOptions, enforcer *casbin.SyncedEnforcer) (authorizer.Authorizer, error) {
	var (
		authorizers []authorizer.Authorizer
	)

	if len(authorizationOpts.WebhookConfigFile) != 0 {
		webhookAuthorizer, err := webhook.New(authorizationOpts.WebhookConfigFile,
			authorizationOpts.WebhookCacheAuthorizedTTL,
			authorizationOpts.WebhookCacheUnauthorizedTTL)
		if err != nil {
			return nil, err
		}

		authorizers = append(authorizers, webhookAuthorizer)
	}

	if len(authorizationOpts.PolicyFile) != 0 {
		abacAuthorizer, err := abac.NewABACAuthorizer(authorizationOpts.PolicyFile)
		if err != nil {
			return nil, err
		}
		authorizers = append(authorizers, abacAuthorizer)
	}

	authorizers = append(authorizers, local.NewAuthorizer(enforcer, authOpts.TenantAdmin))

	return union.New(authorizers...), nil
}
