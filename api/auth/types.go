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

package auth

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentity is an object that contains the metadata about identify used to login
// to TKE.
type LocalIdentity struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	// Spec defines the desired identities of identity in this set.
	Spec   LocalIdentitySpec
	Status LocalIdentityStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentityList is the whole list of all identities.
type LocalIdentityList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of identities.
	Items []LocalIdentity
}

// FinalizerName is the name identifying a finalizer during object lifecycle.
type FinalizerName string

const (
	// LocalIdentityFinalize is an internal finalizer values to LocalIdentity.
	LocalIdentityFinalize FinalizerName = "localidentity"
)

// LocalIdentitySpec is a description of an identity.
type LocalIdentitySpec struct {
	Finalizers []FinalizerName

	Username         string
	DisplayName      string
	Email            string
	PhoneNumber      string
	HashedPassword   string
	OriginalPassword string
	TenantID         string
	Groups           []string
	Extra            map[string]string
}

// LocalIdentityPhase defines the phase of LocalIdentity construct.
type LocalIdentityPhase string

const (
	// LocalIdentityDeleting means the localidentity is undergoing graceful termination.
	LocalIdentityDeleting LocalIdentityPhase = "Deleting"
)

// LocalIdentityStatus is a description of an identity status.
type LocalIdentityStatus struct {
	Locked bool

	Phase LocalIdentityPhase
	// The last time the local identity was updated.
	// +optional
	LastUpdateTime metav1.Time
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKey contains expiration time used to apply the api key.
type APIKey struct {
	metav1.TypeMeta

	// +optional
	metav1.ObjectMeta

	// Spec defines the desired identities of APIkey in this set.
	Spec   APIKeySpec
	Status APIKeyStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentityList is the whole list of all identities.
type APIKeyList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of api keys.
	Items []APIKey
}

// APIKeySpec is a description of an apiKey.
type APIKeySpec struct {
	// APIkey is the jwt token used to authenticate user, and contains user info and sign.
	APIkey string `json:"apiKey,omitempty"`

	TenantID string `json:"tenantID,omitempty"`

	// Creator
	// +optional
	Username string `json:"username,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description"`

	// IssueAt is the created time for api key
	IssueAt metav1.Time `json:"issue_at,omitempty"`

	// ExpireAt is the expire time for api key
	ExpireAt metav1.Time `json:"expire_at,omitempty"`
}

// APIKeyStatus is a description of an api key status.
type APIKeyStatus struct {
	// Disabled represents whether the apikey has been disabled.
	Disabled bool `json:"disabled"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReq contains expiration time used to apply the api key.
type APIKeyReq struct {
	metav1.TypeMeta

	//Exipre is required, holds the duration of the api key become invalid. By default, 168h(= seven days)
	Expire metav1.Duration `json:"expire,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReqPassword contains userinfo and expiration time used to apply the api key.
type APIKeyReqPassword struct {
	metav1.TypeMeta

	// TenantID for user
	TenantID string `json:"tenantID,omitempty"`

	// Username
	Username string `json:"username,omitempty"`

	// Password (encoded by base64)
	Password string `json:"password,omitempty"`

	// Description describes api keys usage.
	Description string `json:"description"`

	// Expire holds the duration of the api key become invalid. By default, 168h(= seven days)
	Expire metav1.Duration `json:"expire,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APISigningKey hold encryption and signing key for api key.
type APISigningKey struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	SigningKey    []byte
	SigningKeyPub []byte
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APISigningKeyList is the whole list of all signing keys.
type APISigningKeyList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []APISigningKey
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Category defines a category of actions for policy.
type Category struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec CategorySpec
}

type CategorySpec struct {
	// CategoryName identifies action category
	CategoryName string
	// DisplayName used to display category name
	DisplayName string
	// +optional
	Description string
	// Actions represents a series of actions work on the policy category
	Actions []Action
}

// Action defines a action verb for authorization.
type Action struct {
	// Name represents user access review request verb.
	Name string
	// Description describes the action.
	Description string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CategoryList is the whole list of policy Category.
type CategoryList struct {
	metav1.TypeMeta
	metav1.ListMeta

	// List of category.
	Items []Category
}

// ProjectPhase defines the phase of policy constructor.
type PolicyPhase string

const (
	// PolicyActive indicates the policy is active.
	PolicyActive PolicyPhase = "Active"
	// PolicyTerminating means the policy is undergoing graceful termination.
	PolicyTerminating PolicyPhase = "Terminating"
)

const (
	// PolicyFinalize is an internal finalizer values to Policy.
	PolicyFinalize FinalizerName = "policy"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Policy represents a policy document for access control.
type Policy struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	// Spec defines the desired identities of policy document in this set.
	Spec PolicySpec

	// +optional
	Status PolicyStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyList is the whole list of all policies.
type PolicyList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of policies.
	Items []Policy
}

// Effect defines the policy effect.
type Effect string

const (
	// Allow is the allow type.
	Allow Effect = "allow"
	// Deny is the deny type.
	Deny Effect = "deny"
)

// PolicyType defines the policy is default or created by user.
type PolicyType string

const (
	PolicyCustom  PolicyType = "custom"
	PolicyDefault PolicyType = "default"
)

// PolicySpec is a description of a policy.
type PolicySpec struct {
	Finalizers []FinalizerName

	PolicyName string
	TenantID    string
	Category    string
	Type        PolicyType
	// Creator
	Username    string
	Description string
	Statement   Statement
	Conditions  []byte
}

// Statement defines a series of action on resource can be done or not.
type Statement struct {
	Actions   []string
	Resources []string
	// Effect indicates action on the resource is allowed or not, can be "allow" or "deny"
	Effect Effect
}

// PolicyStatus represents information about the status of a policy.
type PolicyStatus struct {
	// +optional
	Phase PolicyPhase

	// +optional
	// Subjects represents the objects the policy applies to.
	Subjects []Subject
}

const (
	DefaultRuleModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub)  && keyMatchCustom(r.obj, p.obj) && keyMatchCustom(r.act, p.act)
`
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Rule represents a rule document for access control.
type Rule struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	// Spec defines the desired identities of policy document in this set.
	Spec RuleSpec
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RuleList is the whole list of all policies.
type RuleList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of rules.
	Items []Rule
}

// RuleSpec is a description of a policy.
type RuleSpec struct {
	PType string `json:"ptype"`
	V0    string `json:"v0"`
	V1    string `json:"v1"`
	V2    string `json:"v2"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
	V6    string `json:"v6"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Binding references the objects a policy applies to, but does not contain it.
type Binding struct {
	metav1.TypeMeta

	// Subjects holds references to the objects the policy applies to.
	// +optional
	Subjects []Subject
}

// Subject references a user can specify by id or name.
type Subject struct {
	ID   string
	Name string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMap holds configuration data for tke to consume.
type ConfigMap struct {
	metav1.TypeMeta
	// +optional
	metav1.ObjectMeta

	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string

	// BinaryData contains the binary data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// BinaryData can contain byte sequences that are not in the UTF-8 range.
	// The keys stored in BinaryData must not overlap with the ones in
	// the Data field, this is enforced during validation process.
	// +optional
	BinaryData map[string][]byte
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapList is a resource containing a list of ConfigMap objects.
type ConfigMapList struct {
	metav1.TypeMeta

	// +optional
	metav1.ListMeta

	// Items is the list of ConfigMaps.
	Items []ConfigMap
}
