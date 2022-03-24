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

package policy

import (
	"context"
	"encoding/json"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubectl/pkg/util/rbac"
	"tkestack.io/tke/api/authz"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/apiserver/authentication"
	authzprovider "tkestack.io/tke/pkg/authz/provider"
	"tkestack.io/tke/pkg/util/log"
	namesutil "tkestack.io/tke/pkg/util/names"
)

// Strategy implements verification logic for configmap.
type Strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	platformClient platformversionedclient.PlatformV1Interface
}

const (
	NamePrefix = "pol-"
)

var _ rest.RESTCreateStrategy = &Strategy{}
var _ rest.RESTUpdateStrategy = &Strategy{}
var _ rest.RESTDeleteStrategy = &Strategy{}

// NewStrategy creates a strategy that is the default logic that applies when
// creating and updating namespace set objects.
func NewStrategy(platformClient platformversionedclient.PlatformV1Interface) *Strategy {
	return &Strategy{authz.Scheme, namesutil.Generator, platformClient}
}

func ShouldDeleteDuringUpdate(ctx context.Context, key string, obj, existing runtime.Object) bool {
	pol, ok := obj.(*authz.Policy)
	if !ok {
		log.Errorf("unexpected object, key:%s", key)
		return false
	}
	return len(pol.Finalizers) == 0 && registry.ShouldDeleteDuringUpdate(ctx, key, obj, existing)
}

// DefaultGarbageCollectionPolicy returns the default garbage collection behavior.
func (Strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

// NamespaceScoped is false for namespaceSets
func (Strategy) NamespaceScoped() bool {
	return true
}

// Export strips fields that can not be set by the user.
func (Strategy) Export(ctx context.Context, obj runtime.Object, exact bool) error {
	return nil
}

// PrepareForCreate is invoked on create before validation to normalize
// the object.
func (Strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	username, _ := authentication.UsernameAndTenantID(ctx)
	tenantID := request.NamespaceValue(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	policy := obj.(*authz.Policy)
	policy.TenantID = tenantID
	if policy.Username == "" {
		policy.Username = username
	}
	if policy.Name == "" && policy.GenerateName == "" {
		policy.GenerateName = NamePrefix
	}
	policy.Rules = compactRules(policy.Rules)
	region := authentication.GetExtraValue("region", ctx)
	log.Debugf("region '%v'", region)
	if len(region) != 0 {
		annotations := policy.Annotations
		if len(annotations) == 0 {
			annotations = map[string]string{}
		}
		annotations[authz.GroupName+"/region"] = region[0]
		policy.Annotations = annotations
	}
}

func compactRules(rules []rbacv1.PolicyRule) []rbacv1.PolicyRule {
	if len(rules) != 0 {
		for _, rule := range rules {
			apiGroups := rule.APIGroups
			for j := range rule.APIGroups {
				if apiGroups[j] == "" || apiGroups[j] == "\"\"" || apiGroups[j] == "'\"\"'" {
					apiGroups[j] = ""
				}
			}
		}
		compactedRules, err := rbac.CompactRules(rules)
		if err != nil {
			marshal, _ := json.Marshal(rules)
			log.Errorf("unexpected object, rules:%s", marshal)
		} else {
			return compactedRules
		}
	}
	return rules
}

// PrepareForUpdate is invoked on update before validation to normalize the
// object.
func (Strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	oldPolicy := old.(*authz.Policy)
	policy, _ := obj.(*authz.Policy)
	if policy.TenantID != oldPolicy.TenantID {
		log.Warnf("Unauthorized update policy tenantID '%s'", oldPolicy.TenantID)
		policy.TenantID = oldPolicy.TenantID
	}
	policy.Rules = compactRules(policy.Rules)
}

// Validate validates a new configmap.
func (s Strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	policy := obj.(*authz.Policy)
	provider, err := authzprovider.GetProvider(policy.Annotations)
	if err == nil {
		if fieldErr := provider.Validate(context.TODO(), policy, s.platformClient); fieldErr != nil {
			return field.ErrorList{fieldErr}
		}
	}
	return ValidatePolicy(policy, s.platformClient)
}

// AllowCreateOnUpdate is false for persistent events
func (Strategy) AllowCreateOnUpdate() bool {
	return false
}

// AllowUnconditionalUpdate returns true if the object can be updated
// unconditionally (irrespective of the latest resource version), when there is
// no resource version specified in the object.
func (Strategy) AllowUnconditionalUpdate() bool {
	return false
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (Strategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (Strategy) Canonicalize(obj runtime.Object) {
}

// ValidateUpdate is the default update validation for an end namespace set.
func (s Strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return ValidatePolicyUpdate(ctx, obj.(*authz.Policy), old.(*authz.Policy), s.platformClient)
}

// WarningsOnUpdate returns warnings for the given update.
func (Strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
