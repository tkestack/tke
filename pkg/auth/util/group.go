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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"tkestack.io/tke/api/auth"
	authinternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/auth/internalversion"
	"tkestack.io/tke/pkg/util"
	"tkestack.io/tke/pkg/util/log"
)

func GroupKey(tenantID string, name string) string {
	return fmt.Sprintf("%s%s", GroupPrefix(tenantID), name)
}

func GroupPrefix(tenantID string) string {
	return fmt.Sprintf("%s##group##", tenantID)
}

func GetPoliciesFromGroupExtra(group *auth.LocalGroup) ([]string, bool) {
	var policies []string
	if len(group.Spec.Extra) == 0 {
		return policies, false
	}
	str, exists := group.Spec.Extra[PoliciesKey]
	if !exists {
		return policies, false
	}

	extra := group.Spec.Extra
	delete(extra, PoliciesKey)
	group.Spec.Extra = extra

	splits := strings.Split(str, ",")

	for _, p := range splits {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "pol-") {
			policies = append(policies, p)
		}
	}

	return policies, true
}

func BindGroupPolicies(authClient authinternalclient.AuthInterface, group *auth.LocalGroup, policies []string) error {
	var errs []error
	for _, p := range policies {
		binding := auth.Binding{}
		binding.Groups = append(binding.Groups, auth.Subject{ID: group.Name, Name: group.Spec.DisplayName})
		pol := &auth.Policy{}
		err := authClient.RESTClient().Post().
			Resource("policies").
			Name(p).
			SubResource("binding").
			Body(&binding).
			Do().Into(pol)
		if err != nil {
			log.Error("bind policy for group failed", log.String("group", group.Name),
				log.String("policy", p), log.Err(err))
			errs = append(errs, err)
		}
	}

	return errors.NewAggregate(errs)
}

func UnBindGroupPolicies(authClient authinternalclient.AuthInterface, group *auth.LocalGroup, policies []string) error {
	var errs []error
	for _, p := range policies {
		binding := auth.Binding{}
		binding.Groups = append(binding.Groups, auth.Subject{ID: group.Name, Name: group.Spec.DisplayName})
		pol := &auth.Policy{}
		err := authClient.RESTClient().Post().
			Resource("policies").
			Name(p).
			SubResource("unbinding").
			Body(&binding).
			Do().Into(pol)
		if err != nil {
			log.Error("unbind policy for group failed", log.String("group", group.Spec.DisplayName),
				log.String("policy", p), log.Err(err))
			errs = append(errs, err)
		}
	}

	return errors.NewAggregate(errs)
}

func HandleGroupPoliciesUpdate(authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer, group *auth.LocalGroup) error {
	newPolicies, needHandlePolicy := GetPoliciesFromGroupExtra(group)
	if !needHandlePolicy {
		return nil
	}

	roles := enforcer.GetRolesForUserInDomain(GroupKey(group.Spec.TenantID, group.Name), "")
	var oldPolicies []string
	for _, r := range roles {
		if strings.HasPrefix(r, "pol-") {
			oldPolicies = append(oldPolicies, r)
		}
	}

	added, removed := util.DiffStringSlice(oldPolicies, newPolicies)

	log.Info("handler group policies ", log.Strings("added", added), log.Strings("removed", removed))
	berr := BindGroupPolicies(authClient, group, added)
	if berr != nil {
		log.Error("bind group policies failed", log.String("group", group.Spec.Username), log.Strings("policies", added), log.Err(berr))
	}

	uerr := UnBindGroupPolicies(authClient, group, removed)
	if berr != nil {
		log.Error("un bind group policies failed", log.String("group", group.Spec.Username), log.Strings("policies", removed), log.Err(uerr))
	}

	return errors.NewAggregate([]error{berr, uerr})
}

func FillGroupPolicies(authClient authinternalclient.AuthInterface, enforcer *casbin.SyncedEnforcer, groupList *auth.LocalGroupList) {
	if enforcer == nil || enforcer.GetRoleManager() == nil || enforcer.GetAdapter() == nil {
		return
	}

	policyDisplayNameMap := make(map[string]string)
	for i, item := range groupList.Items {
		roles := enforcer.GetRolesForUserInDomain(GroupKey(item.Spec.TenantID, item.Name), "")
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
				pol, err := authClient.Policies().Get(p, v1.GetOptions{})
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
			log.Error("Marshal policy map for group failed", log.String("group", item.Spec.Username), log.Err(err))
			continue
		}

		if groupList.Items[i].Spec.Extra == nil {
			groupList.Items[i].Spec.Extra = make(map[string]string)
		}

		groupList.Items[i].Spec.Extra[PoliciesKey] = string(b)
	}
}
