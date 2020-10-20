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
	"tkestack.io/tke/api/auth"
)

func ConvertPolicyToRuleArray(policy *auth.Policy) [][]string {
	var rules [][]string
	if policy.Spec.Scope != auth.PolicyProject && len(policy.Status.Users) == 0 && len(policy.Status.Groups) == 0 {
		return rules
	}
	for _, act := range policy.Spec.Statement.Actions {
		for _, res := range policy.Spec.Statement.Resources {
			rule := []string{policy.Name, "*", res, act, string(policy.Spec.Statement.Effect)}
			rules = append(rules, rule)
		}
	}

	return rules
}

func ConvertPolicyToRuleArrayUsingRuleName(roleName string, policy *auth.Policy) [][]string {
	var rules [][]string
	for _, act := range policy.Spec.Statement.Actions {
		for _, res := range policy.Spec.Statement.Resources {
			rule := []string{roleName, "*", res, act, string(policy.Spec.Statement.Effect)}
			rules = append(rules, rule)
		}
	}

	return rules
}
