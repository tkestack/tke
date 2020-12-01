/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"tkestack.io/tke/api/auth"
	authv1 "tkestack.io/tke/api/auth/v1"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"

	"github.com/casbin/casbin/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/errors"
)

const (
	PolicyTag        = "policy"
	PoliciesKey      = "policies"
	administratorKey = "administrator"

	PlatformPolicyPattern      = "pol-%s-platform"
	AdministratorPolicyPattern = "pol-%s-administrator"
)

func GetLocalIdentity(ctx context.Context, authClient authinternalclient.AuthInterface, tenantID, username string) (auth.LocalIdentity, error) {
	tenantUserSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.tenantID", tenantID),
		fields.OneTermEqualSelector("spec.username", username))

	localIdentityList, err := authClient.LocalIdentities().List(ctx, v1.ListOptions{FieldSelector: tenantUserSelector.String()})
	if err != nil {
		return auth.LocalIdentity{}, err
	}

	if len(localIdentityList.Items) == 0 {
		return auth.LocalIdentity{}, apierrors.NewNotFound(auth.Resource("localIdentity"), username)
	}

	return localIdentityList.Items[0], nil
}

func GetUserByName(ctx context.Context, authClient authinternalclient.AuthInterface, tenantID, username string) (auth.User, error) {
	tenantUserSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("spec.tenantID", tenantID),
		fields.OneTermEqualSelector("spec.username", username))

	userList, err := authClient.Users().List(ctx, v1.ListOptions{FieldSelector: tenantUserSelector.String()})
	if err != nil {
		return auth.User{}, err
	}

	if len(userList.Items) == 0 {
		return auth.User{}, apierrors.NewNotFound(auth.Resource("users"), username)
	}

	return userList.Items[0], nil

}

func UserKey(tenantID string, name string) string {
	return fmt.Sprintf("%s%s", UserPrefix(tenantID), name)
}

func UserPrefix(tenantID string) string {
	return fmt.Sprintf("%s##user##", tenantID)
}

func ProjectPolicyName(projectID string, policyID string) string {
	return fmt.Sprintf("%s-%s", projectID, policyID)
}

func GetGroupsForUser(ctx context.Context, authClient authinternalclient.AuthInterface, userID string) (auth.LocalGroupList, error) {
	groupList := auth.LocalGroupList{}
	err := authClient.RESTClient().Get().
		Resource("localidentities").
		Name(userID).
		SubResource("groups").Do(ctx).Into(&groupList)

	return groupList, err
}

func ParseTenantAndName(str string) (string, string) {
	parts := strings.Split(str, "##")
	if len(parts) > 1 {
		return parts[0], parts[1]
	}

	return "", str
}

func CombineTenantAndName(tenantID, name string) string {
	return fmt.Sprintf("%s##%s", tenantID, name)
}

func GetPoliciesFromUserExtra(localIdentity *auth.LocalIdentity) ([]string, bool) {
	var policies []string
	if len(localIdentity.Spec.Extra) == 0 {
		return policies, false
	}
	str, exists := localIdentity.Spec.Extra[PoliciesKey]
	if !exists {
		return policies, false
	}

	extra := localIdentity.Spec.Extra
	delete(extra, PoliciesKey)
	localIdentity.Spec.Extra = extra

	splits := strings.Split(str, ",")

	for _, p := range splits {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "pol-") {
			policies = append(policies, p)
		}
	}

	return policies, true
}

func BindUserPolicies(ctx context.Context, authClient authinternalclient.AuthInterface, localIdentity *auth.LocalIdentity, policies []string) error {
	var errs []error
	for _, p := range policies {
		binding := auth.Binding{}
		binding.Users = append(binding.Users, auth.Subject{ID: localIdentity.Name, Name: localIdentity.Spec.Username})
		pol := &auth.Policy{}
		err := authClient.RESTClient().Post().
			Resource("policies").
			Name(p).
			SubResource("binding").
			Body(&binding).
			Do(ctx).Into(pol)
		if err != nil {
			log.Error("bind policy for user failed", log.String("user", localIdentity.Spec.Username),
				log.String("policy", p), log.Err(err))
			errs = append(errs, err)
		}
	}

	return errors.NewAggregate(errs)
}

func UnBindUserPolicies(ctx context.Context, authClient authinternalclient.AuthInterface, localIdentity *auth.LocalIdentity, policies []string) error {
	var errs []error
	for _, p := range policies {
		binding := auth.Binding{}
		binding.Users = append(binding.Users, auth.Subject{ID: localIdentity.Name, Name: localIdentity.Spec.Username})
		pol := &auth.Policy{}
		err := authClient.RESTClient().Post().
			Resource("policies").
			Name(p).
			SubResource("unbinding").
			Body(&binding).
			Do(ctx).Into(pol)
		if err != nil {
			log.Error("unbind policy for user failed", log.String("user", localIdentity.Spec.Username),
				log.String("policy", p), log.Err(err))
			errs = append(errs, err)
		}
	}

	return errors.NewAggregate(errs)
}

func HandleUserPoliciesUpdate(ctx context.Context, authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer, localIdentity *auth.LocalIdentity) error {
	newPolicies, needHandlePolicy := GetPoliciesFromUserExtra(localIdentity)
	if !needHandlePolicy {
		return nil
	}

	roles := enforcer.GetRolesForUserInDomain(UserKey(localIdentity.Spec.TenantID, localIdentity.Spec.Username), DefaultDomain)
	var oldPolicies []string
	for _, r := range roles {
		if strings.HasPrefix(r, "pol-") {
			oldPolicies = append(oldPolicies, r)
		}
	}

	added, removed := util.DiffStringSlice(oldPolicies, newPolicies)

	log.Info("handler user policies ", log.Strings("added", added), log.Strings("removed", removed))
	berr := BindUserPolicies(ctx, authClient, localIdentity, added)
	if berr != nil {
		log.Error("bind user policies failed", log.String("user", localIdentity.Spec.Username), log.Strings("policies", added), log.Err(berr))
	}

	uerr := UnBindUserPolicies(ctx, authClient, localIdentity, removed)
	if berr != nil {
		log.Error("un bind user policies failed", log.String("user", localIdentity.Spec.Username), log.Strings("policies", removed), log.Err(uerr))
	}

	return errors.NewAggregate([]error{berr, uerr})
}

func SetAdministrator(enforcer *casbin.SyncedEnforcer, localIdentity *auth.LocalIdentity, idp *auth.IdentityProvider) {
	if localIdentity.Spec.Extra == nil {
		localIdentity.Spec.Extra = make(map[string]string)
	}
	localIdentity.Spec.Extra[administratorKey] = "false"
	// Use implicit roles to check admin
	roles, err := enforcer.GetImplicitRolesForUser(UserKey(localIdentity.Spec.TenantID, localIdentity.Spec.Username), DefaultDomain)
	if err != nil {
		log.Error("Get implicit roles for user failed", log.String("user", localIdentity.Spec.Username), log.Err(err))
		return
	}
	for _, r := range roles {
		if r == fmt.Sprintf(PlatformPolicyPattern, localIdentity.Spec.TenantID) ||
			r == fmt.Sprintf(AdministratorPolicyPattern, localIdentity.Spec.TenantID) {
			localIdentity.Spec.Extra[administratorKey] = "true"
			return
		}
	}

	if idp != nil {
		for _, name := range idp.Spec.Administrators {
			if name == localIdentity.Spec.Username {
				localIdentity.Spec.Extra[administratorKey] = "true"
				return
			}
		}
	}
}

func IsPlatformAdministrator(user authv1.User) bool {
	if user.Spec.Extra != nil && user.Spec.Extra[administratorKey] == "true" {
		return true
	}
	return false
}

func FillUserPolicies(ctx context.Context, authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer,
	localidentityList *auth.LocalIdentityList) {
	if enforcer == nil || enforcer.GetRoleManager() == nil || enforcer.GetAdapter() == nil {
		return
	}

	idpList, err := authClient.IdentityProviders().List(ctx, v1.ListOptions{})
	if err != nil || idpList == nil {
		return
	}

	idpMap := make(map[string]*auth.IdentityProvider)
	for i := 0; i < len(idpList.Items); i++ {
		idpMap[idpList.Items[i].GetName()] = &idpList.Items[i]
	}

	policyDisplayNameMap := make(map[string]string)
	for i, item := range localidentityList.Items {
		if idp, ok := idpMap[localidentityList.Items[i].Spec.TenantID]; ok {
			SetAdministrator(enforcer, &localidentityList.Items[i], idp)
		} else {
			SetAdministrator(enforcer, &localidentityList.Items[i], nil)
		}

		// Use direct roles to fill policies
		roles, _ := enforcer.GetRoleManager().GetRoles(UserKey(item.Spec.TenantID, item.Spec.Username), DefaultDomain)
		var policies []string
		for _, r := range roles {
			if strings.HasPrefix(r, "pol-") {
				policies = append(policies, r)
			}
		}

		m := make(map[string]string)
		for _, p := range policies {
			displayName, ok := policyDisplayNameMap[p]
			if ok {
				m[p] = displayName
			} else {
				pol, err := authClient.Policies().Get(ctx, p, v1.GetOptions{})
				if err != nil {
					log.Error("get policy failed", log.String("policy", p), log.Err(err))
					continue
				}

				m[p] = pol.Spec.DisplayName
				policyDisplayNameMap[p] = pol.Spec.DisplayName
			}
		}

		b, err := json.Marshal(m)
		if err != nil {
			log.Error("Marshal policy map for user failed", log.String("user", item.Spec.Username), log.Err(err))
			continue
		}

		localidentityList.Items[i].Spec.Extra[PoliciesKey] = string(b)

	}
}
