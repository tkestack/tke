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
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"

	genericfilter "tkestack.io/tke/pkg/apiserver/filter"

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

// the kube verb associated with API requests (this includes get, list, watch, create, update, patch, delete, deletecollection, and proxy),
// or the lowercased HTTP verb associated with non-API requests (this includes get, put, post, patch, and delete)
var verbMap = sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "proxy", "put", "post", "deletecollection")

// Authorizer implement the authorize interface that use local repository to
// authorize the subject access review.
type Authorizer struct {
	privilegedUsername string

	authClient authinternalclient.AuthInterface
	enforcer   *casbin.SyncedEnforcer
}

// NewAuthorizer creates a local repository authorizer and returns it.
func NewAuthorizer(authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer, privilegedUsername string) *Authorizer {
	return &Authorizer{
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
		tenantID  string
		projectID string
		debug     bool
	)
	extra := attr.GetUser().GetExtra()
	if len(extra) > 0 {
		if tenantIDs, ok := extra[genericoidc.TenantIDKey]; ok {
			if len(tenantIDs) > 0 {
				tenantID = tenantIDs[0]
			} else {
				tenantID = "default"
			}
		}

		if debugs, ok := extra[debugKey]; ok {
			if len(debugs) > 0 && debugs[0] == "true" {
				debug = true
			}
		}
	}
	find := false
	if tenantID == "" {
		find, tenantID = genericfilter.FindValueFromGroups(attr.GetUser().GetGroups(), "tenant")
		if find && tenantID == "" {
			tenantID = "default"
		}
	}
	projectID = genericfilter.GetValueFromGroups(attr.GetUser().GetGroups(), "project")
	log.Debug("Authorize", log.String("subject", subject), log.String("action", action),
		log.String("resource", resource), log.String("project", projectID), log.String("tenant", tenantID))

	// First check if user is privileged
	if subject == a.privilegedUsername {
		log.Debug("privileged user", log.String("subject", subject))
		return authorizer.DecisionAllow, "", nil
	}

	// Second check if user is a admin of the identity provider for tenant.
	if tenantID != "" {
		idp, err := a.authClient.IdentityProviders().Get(ctx, tenantID, metav1.GetOptions{})
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
		log.Debug("UnprotectedAuthorized", log.String("action", action))
		return authorizer.DecisionAllow, "", nil
	}

	if debug {
		perms, err := a.enforcer.GetImplicitPermissionsForUser(authutil.UserKey(tenantID, subject), projectID)
		if err != nil {
			log.Error("Get permissions for user failed", log.String("user", authutil.UserKey(tenantID, subject)), log.String("projectID", projectID), log.Err(err))
		} else {
			log.Info("Authorize get user perms", log.String("user", authutil.UserKey(tenantID, subject)), log.Any("user perm", perms))
			data, _ := json.Marshal(perms)
			reason = string(data)
		}
	}

	if projectID != "" && attr.GetName() != "*" && resource == fmt.Sprintf("project:%s", attr.GetName()) && attr.GetName() != projectID {
		return authorizer.DecisionDeny, fmt.Sprintf("unmatched projectIDs: %v %v", attr.GetName(), projectID), nil
	}

	if !strings.HasPrefix(resource, "/api") {
		ns := genericfilter.GetValueFromGroups(attr.GetUser().GetGroups(), "namespace")
		if ns != "" && ns != attr.GetNamespace() {
			log.Errorf("want to access namespace '%s', but only allowed to access %s, attr: %v", attr.GetNamespace(), ns, attr)
			return authorizer.DecisionDeny, fmt.Sprintf("can NOT access namespace other than %v", ns), nil
		}
	}

	if tenantID != "" && verbMap.Has(action) {
		record := &authorizer.AttributesRecord{
			User:            attr.GetUser(),
			Verb:            attr.GetVerb(),
			Namespace:       attr.GetNamespace(),
			APIGroup:        attr.GetAPIGroup(),
			Resource:        attr.GetResource(),
			Subresource:     attr.GetSubresource(),
			Name:            attr.GetName(),
			ResourceRequest: attr.IsResourceRequest(),
			Path:            attr.GetPath(),
		}
		tkeAttributes := filter.ConvertTKEAttributes(ctx, record)
		attrStr, _ := json.Marshal(attr)
		tkeAttributesStr, _ := json.Marshal(tkeAttributes)
		log.Debugf("Attribute '%s' converted to TKEAttributes '%s'", string(attrStr), string(tkeAttributesStr))
		attr = tkeAttributes
	}
	return a.casbinDecision(attr, tenantID, subject, projectID, attr.GetResource(), attr.GetVerb(), reason, debug)
}

func (a *Authorizer) casbinDecision(attr authorizer.Attributes, tenantID, subject, projectID, resource, action, reason string, debug bool) (authorizer.Decision, string, error) {
	allow, err := a.enforcer.Enforce(authutil.UserKey(tenantID, subject), projectID, resource, action)
	if err != nil {
		log.Error("Casbin enforcer failed", log.Any("att", attr), log.String("projectID", projectID), log.String("subj", subject), log.String("act", action), log.String("res", resource), log.Err(err))
		return authorizer.DecisionDeny, "", err
	}
	if !allow {
		allowAll, err := a.enforcer.Enforce(authutil.UserKey(tenantID, authutil.DefaultAll), projectID, resource, action)
		if err == nil && allowAll {
			return authorizer.DecisionAllow, reason, nil
		}
		log.Info("Casbin enforcer: ", log.Any("att", attr), log.String("projectID", projectID), log.String("subj", subject), log.String("act", action), log.String("res", resource), log.String("allow", "false"))
		if debug {
			return authorizer.DecisionDeny, reason, nil
		}
		return authorizer.DecisionDeny, fmt.Sprintf("permission for %s on %s not verify", action, resource), nil
	}
	log.Debug("Casbin enforcer: ", log.Any("att", attr), log.String("projectID", projectID), log.String("subj", subject), log.String("act", action), log.String("res", resource), log.String("allow", "true"))
	return authorizer.DecisionAllow, reason, nil
}
