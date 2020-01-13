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
	"context"
	"encoding/json"
	"fmt"

	"github.com/casbin/casbin/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authorization/authorizer"

	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	"tkestack.io/tke/pkg/auth/filter"
	authutil "tkestack.io/tke/pkg/auth/util"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
)

var (
	debugKey = "debug"
)

// Authorizer implement the authorize interface that use local repository to
// authorize the subject access review.
type Authorizer struct {
	tenantAdmin        string
	privilegedUsername string

	authClient authinternalclient.AuthInterface
	enforcer   *casbin.SyncedEnforcer
}

// NewAuthorizer creates a local repository authorizer and returns it.
func NewAuthorizer(authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer, tenantAdmin string, privilegedUsername string) *Authorizer {
	return &Authorizer{
		tenantAdmin:        tenantAdmin,
		privilegedUsername: privilegedUsername,
		authClient:         authClient,
		enforcer:           enforcer,
	}
}

// Authorize to determine the subject access.
func (a *Authorizer) Authorize(ctx context.Context, attr authorizer.Attributes) (authorized authorizer.Decision, reason string, err error) {
	subject := attr.GetUser().GetName()
	action := attr.GetVerb()
	resource := attr.GetResource()
	var (
		tenantID string
		debug    bool
	)
	extra := attr.GetUser().GetExtra()
	if len(extra) > 0 {
		if tenantIDs, ok := extra[genericoidc.TenantIDKey]; ok {
			if len(tenantIDs) > 0 {
				tenantID = tenantIDs[0]
			}
		}

		if debugs, ok := extra[debugKey]; ok {
			if len(debugs) > 0 && debugs[0] == "true" {
				debug = true
			}
		}
	}

	// First check if user is tenantAdmin or privileged
	if subject == a.tenantAdmin || subject == a.privilegedUsername {
		return authorizer.DecisionAllow, "", nil
	}

	// Second check if user is a admin of the identity provider for tenant.
	if tenantID != "" {
		idp, err := a.authClient.IdentityProviders().Get(tenantID, metav1.GetOptions{})
		if err != nil {
			log.Error("Get identity provider for tenant failed", log.String("tenant", tenantID), log.Err(err))
			return authorizer.DecisionDeny, "", err
		}

		if util.InStringSlice(idp.Spec.Administrators, subject) {
			return authorizer.DecisionAllow, "", nil
		}
	}

	authorized = filter.UnprotectedAuthorized(attr)
	if authorized == authorizer.DecisionAllow {
		return authorizer.DecisionAllow, "", nil
	}

	if debug {
		perms, err := a.enforcer.GetImplicitPermissionsForUser(authutil.UserKey(tenantID, subject))
		if err != nil {
			log.Error("Get permissions for user failed", log.String("user", authutil.UserKey(tenantID, subject)), log.Err(err))
		} else {
			log.Info("Authorize get user perms", log.String("user", subject), log.Any("user perm", perms))
			data, _ := json.Marshal(perms)
			reason = string(data)
		}
	}

	allow, err := a.enforcer.Enforce(authutil.UserKey(tenantID, subject), resource, action)
	if err != nil {
		log.Error("Casbin enforcer failed", log.Any("att", attr), log.String("subj", subject), log.String("act", action), log.String("res", resource), log.Err(err))
		return authorizer.DecisionDeny, "", err
	}
	if !allow {
		log.Info("Casbin enforcer: ", log.Any("att", attr), log.String("subj", subject), log.String("act", action), log.String("res", resource), log.String("allow", "false"))
		if debug {
			return authorizer.DecisionDeny, reason, nil
		}
		return authorizer.DecisionDeny, fmt.Sprintf("permission for %s on %s not verify", action, resource), nil
	}

	log.Debug("Casbin enforcer: ", log.Any("att", attr), log.String("subj", subject), log.String("act", action), log.String("res", resource), log.String("allow", "true"))
	return authorizer.DecisionAllow, reason, nil
}
