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

package local

import (
	"fmt"

	"github.com/casbin/casbin"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/auth/filter"
	"tkestack.io/tke/pkg/auth/types"
	"tkestack.io/tke/pkg/util/log"
)

// Authorizer implement the authorize interface that use local repository to
// authorize the subject access review.
type Authorizer struct {
	admin    string
	enforcer *casbin.SyncedEnforcer
}

// NewAuthorizer creates a local repository authorizer and returns it.
func NewAuthorizer(enforcer *casbin.SyncedEnforcer, admin string) *Authorizer {
	return &Authorizer{
		enforcer: enforcer,
		admin:    admin,
	}
}

// Authorize to determine the subject access.
func (a *Authorizer) Authorize(attr authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
	subject := attr.GetUser().GetName()
	action := attr.GetVerb()
	resource := attr.GetResource()
	var tenantID string
	extra := attr.GetUser().GetExtra()
	if len(extra) > 0 {
		if tenantIDs, ok := extra[genericoidc.TenantIDKey]; ok {
			if len(tenantIDs) > 0 {
				tenantID = tenantIDs[0]
			}
		}
	}

	// First check if user is admin
	if subject == a.admin {
		return authorizer.DecisionAllow, "", nil
	}

	authorized = filter.UnprotectedAuthorized(attr)
	if authorized == authorizer.DecisionAllow {
		return authorizer.DecisionAllow, "", nil
	}

	log.Debug("Authorize get user perms", log.Any("user perm", a.enforcer.GetImplicitPermissionsForUser(fmt.Sprintf("%s%s-%s", types.UserPrefix, tenantID, subject))))
	if !a.enforcer.Enforce(fmt.Sprintf("%s%s-%s", types.UserPrefix, tenantID, subject), resource, action) {
		log.Info("Casbin enforcer: ", log.Any("att", attr), log.String("subj", subject), log.String("act", action), log.String("res", resource), log.String("allow", "false"))
		return authorizer.DecisionDeny, "permission not verify", nil
	}

	log.Debug("Casbin enforcer: ", log.Any("att", attr), log.String("subj", subject), log.String("act", action), log.String("res", resource), log.String("allow", "true"))
	return authorizer.DecisionAllow, "", nil
}
