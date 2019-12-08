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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentity is an object that contains the metadata about identify used to
// login to TKE.
type LocalIdentity struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of identity in this set.
	// +optional
	Spec LocalIdentitySpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status LocalIdentityStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalIdentityList is the whole list of all identities.
type LocalIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of identities.
	Items []LocalIdentity `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type FinalizerName string

const (
	// LocalIdentityFinalize is an internal finalizer values to LocalIdentity.
	LocalIdentityFinalize FinalizerName = "localidentity"

	// PolicyFinalize is an internal finalizer values to Policy.
	PolicyFinalize FinalizerName = "policy"

	// GroupFinalize is an internal finalizer values to Group.
	GroupFinalize FinalizerName = "group"

	// RoleFinalize is an internal finalizer values to Role.
	RoleFinalize FinalizerName = "role"
)

// LocalIdentitySpec is a description of an identity.
type LocalIdentitySpec struct {
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,11,rep,name=finalizers,casttype=FinalizerName"`

	Username         string `json:"username" protobuf:"bytes,7,opt,name=name"`
	DisplayName      string `json:"displayName" protobuf:"bytes,8,opt,name=displayName"`
	Email            string `json:"email" protobuf:"bytes,9,opt,name=email"`
	PhoneNumber      string `json:"phoneNumber" protobuf:"bytes,10,opt,name=phone"`
	HashedPassword   string `json:"hashedPassword,omitempty" protobuf:"bytes,4,opt,name=hashedPassword"`
	OriginalPassword string `json:"originalPassword,omitempty" protobuf:"bytes,5,opt,name=originalPassword"`
	// +optional
	TenantID string `json:"tenantID,omitempty" protobuf:"bytes,2,opt,name=tenantID"`

	Groups []string `json:"groups,omitempty" protobuf:"bytes,6,rep,name=groups"`
	// +optional
	Extra map[string]string `json:"extra,omitempty" protobuf:"bytes,3,rep,name=extra"`
}

// LocalIdentityPhase defines the phase of LocalIdentity construct.
type LocalIdentityPhase string

const (
	// LocalIdentityDeleting means the localidentity is undergoing graceful termination.
	LocalIdentityDeleting LocalIdentityPhase = "Deleting"
)

// LocalIdentityStatus is a description of an identity status.
type LocalIdentityStatus struct {
	Phase LocalIdentityPhase `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase,casttype=LocalIdentityPhase"`

	// +optional
	Locked bool `json:"locked,omitempty" protobuf:"varint,1,opt,name=locked"`

	// The last time the local identity was updated.
	// +optional
	LastUpdateTime metav1.Time `protobuf:"bytes,2,opt,name=lastUpdateTime"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKey contains expiration time used to apply the api key.
type APIKey struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=objectMeta"`

	// Spec defines the desired identities of APIkey in this set.
	Spec   APIKeySpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status APIKeyStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyList is the whole list of all identities.
type APIKeyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=listMeta"`
	// List of api keys.
	Items []APIKey `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// APIKeySpec is a description of an apiKey.
type APIKeySpec struct {
	// APIkey is the jwt token used to authenticate user, and contains user info and sign.
	APIkey string `json:"apiKey,omitempty" protobuf:"bytes,1,opt,name=apiKey"`

	TenantID string `json:"tenantID,omitempty" protobuf:"bytes,5,opt,name=tenantID"`

	// Username is creator
	// +optional
	Username string `json:"username,omitempty" protobuf:"bytes,6,opt,name=username"`

	// Description describes api keys usage.
	// +optional
	Description string `json:"description" protobuf:"bytes,2,opt,name=description"`

	// IssueAt is the created time for api key
	IssueAt metav1.Time `json:"issue_at,omitempty" protobuf:"bytes,3,opt,name=issue_at,json=issueAt"`

	// ExpireAt is the expire time for api key
	ExpireAt metav1.Time `json:"expire_at,omitempty" protobuf:"bytes,4,opt,name=expire_at,json=expireAt"`
}

// APIKeyStatus is a description of an api key status.
type APIKeyStatus struct {
	// Disabled represents whether the apikey has been disabled.
	// +optional
	Disabled bool `json:"disabled" protobuf:"varint,1,opt,name=disabled"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReq contains expiration time used to apply the api key.
type APIKeyReq struct {
	metav1.TypeMeta `json:",inline"`

	//Exipre is required, holds the duration of the api key become invalid. By default, 168h(= seven days)
	// +optional
	Expire metav1.Duration `json:"expire,omitempty" protobuf:"bytes,2,opt,name=expire"`

	// Description describes api keys usage.
	Description string `json:"description" protobuf:"bytes,3,opt,name=description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReqPassword contains userinfo and expiration time used to apply the api key.
type APIKeyReqPassword struct {
	metav1.TypeMeta `json:",inline"`

	// TenantID for user
	TenantID string `json:"tenantID,omitempty" protobuf:"bytes,2,opt,name=tenantID"`

	// Username
	Username string `json:"username,omitempty" protobuf:"bytes,3,opt,name=username"`

	// Password (encoded by base64)
	Password string `json:"password,omitempty" protobuf:"bytes,4,opt,name=password"`

	// Description describes api keys usage.
	// +optional
	Description string `json:"description" protobuf:"bytes,5,opt,name=description"`

	// Expire holds the duration of the api key become invalid. By default, 168h(= seven days)
	// +optional
	Expire metav1.Duration `json:"expire,omitempty" protobuf:"bytes,6,opt,name=expire"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APISigningKey hold encryption and signing key.
type APISigningKey struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +optional
	SigningKey []byte `json:"signingKey,omitempty" protobuf:"bytes,2,opt,name=signingKey"`
	// +optional
	SigningKeyPub []byte `json:"signingKeyPub,omitempty" protobuf:"bytes,3,opt,name=signingKeyPub"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APISigningKeyList is the whole list of all signing key.
type APISigningKeyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of keys.
	Items []APISigningKey `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Category defines a category of actions for policy.
type Category struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec CategorySpec `protobuf:"bytes,2,opt,name=spec"`
}

// CategorySpec is a description of category.
type CategorySpec struct {
	// CategoryName identifies action category
	CategoryName string `json:"categoryName" protobuf:"bytes,1,opt,name=categoryName"`
	// DisplayName used to display category name
	DisplayName string `json:"displayName" protobuf:"bytes,2,opt,name=displayName"`
	// +optional
	Description string `json:"description" protobuf:"bytes,3,opt,name=description"`
	// Actions represents a series of actions work on the policy category
	Actions []Action `json:"actions" protobuf:"bytes,4,rep,name=actions"`
}

// Action defines a action verb for authorization.
type Action struct {
	// Name represents user access review request verb.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Description describes the action.
	Description string `json:"description" protobuf:"bytes,2,opt,name=description"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CategoryList is the whole list of policy Category.
type CategoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of category.
	Items []Category `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// PolicyPhase defines the phase of policy constructor.
type PolicyPhase string

const (
	// PolicyActive indicates the policy is active.
	PolicyActive PolicyPhase = "Active"
	// PolicyTerminating means the policy is undergoing graceful termination.
	PolicyTerminating PolicyPhase = "Terminating"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Policy represents a policy document for access control.
type Policy struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of policy document in this set.
	// +optional
	Spec PolicySpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status PolicyStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyList is the whole list of all policies.
type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of policies.
	Items []Policy `json:"items" protobuf:"bytes,2,rep,name=items"`
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
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,8,rep,name=finalizers,casttype=FinalizerName"`

	DisplayName string     `json:"displayName" protobuf:"bytes,7,opt,name=displayName"`
	TenantID    string     `json:"tenantID" protobuf:"bytes,1,opt,name=tenantID"`
	Category    string     `json:"category" protobuf:"bytes,9,opt,name=category"`
	Type        PolicyType `json:"type" protobuf:"varint,10,opt,name=type,casttype=PolicyType"`
	Username    string     `json:"username" protobuf:"bytes,2,opt,name=username"`
	// +optional
	Description string `json:"description" protobuf:"bytes,3,opt,name=description"`

	Statement Statement `json:"statement" protobuf:"bytes,5,rep,name=statement"`
	// +optional
	Conditions []byte `json:"conditions,omitempty" protobuf:"bytes,6,rep,name=conditions"`
}

// Statement defines a series of action on resource can be done or not.
type Statement struct {
	Actions   []string `json:"actions" protobuf:"bytes,1,rep,name=actions"`
	Resources []string `json:"resources" protobuf:"bytes,2,rep,name=resources"`
	// Effect indicates action on the resource is allowed or not, can be "allow" or "deny"
	Effect Effect `json:"effect" protobuf:"bytes,3,opt,name=effect,casttype=Effect"`
}

// PolicyStatus represents information about the status of a policy.
type PolicyStatus struct {
	// +optional
	Phase PolicyPhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=PolicyPhase"`
	// +optional
	// Subjects represents the objects the policy applies to.
	Subjects []Subject `json:"subjects" protobuf:"bytes,2,rep,name=subjects"`
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
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec defines the desired identities of policy document in this set.
	Spec RuleSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RuleList is the whole list of all rules.
type RuleList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of rules.
	Items []Rule `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// RuleSpec is a description of a rule.
type RuleSpec struct {
	PType string `json:"ptype" protobuf:"bytes,1,opt,name=ptype"`
	V0    string `json:"v0" protobuf:"bytes,2,opt,name=v0"`
	V1    string `json:"v1" protobuf:"bytes,3,opt,name=v1"`
	V2    string `json:"v2" protobuf:"bytes,4,opt,name=v2"`
	V3    string `json:"v3" protobuf:"bytes,5,opt,name=v3"`
	V4    string `json:"v4" protobuf:"bytes,6,opt,name=v4"`
	V5    string `json:"v5" protobuf:"bytes,7,opt,name=v5"`
	V6    string `json:"v6" protobuf:"bytes,8,opt,name=v6"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Binding is used to bind or unbind the subjects to or from the policy,role or group.
type Binding struct {
	metav1.TypeMeta `json:",inline"`

	// Subjects holds references to the objects the policy applies to.
	// +optional
	Subjects []Subject `json:"subjects" protobuf:"bytes,1,rep,name=subject"`
}

// Subject references a user can specify by id or name.
type Subject struct {
	ID   string `json:"id" protobuf:"bytes,1,opt,name=id"`
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Role is a collection with multiple policies.
type Role struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of role document in this set.
	Spec RoleSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status RoleStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RoleList is the whole list of policy.
type RoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of rules.
	Items []Role `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// RolePhase defines the phase of role constructor.
type RolePhase string

const (
	RoleActive RolePhase = "Active"
	// RoleTerminating means the role is undergoing graceful termination.
	RoleTerminating RolePhase = "Terminating"
)

// RoleSpec is a description of role.
type RoleSpec struct {
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=FinalizerName"`

	DisplayName string `json:"displayName" protobuf:"bytes,2,opt,name=displayName"`
	TenantID    string `json:"tenantID" protobuf:"bytes,3,opt,name=tenantID"`

	//Creator
	Username    string `json:"username" protobuf:"bytes,4,opt,name=username"`
	Description string `json:"description" protobuf:"bytes,5,opt,name=description"`

	Policies []string `json:"policies" protobuf:"bytes,6,rep,name=policies"`
}

// RoleStatus represents information about the status of a role.
type RoleStatus struct {
	// +optional
	Phase RolePhase `protobuf:"bytes,1,opt,name=phase,casttype=RolePhase"`

	// Subjects represents the members of the group.
	Subjects []Subject `protobuf:"bytes,2,rep,name=subjects"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyBinding references the request to bind or unbind policies to the role.
type PolicyBinding struct {
	metav1.TypeMeta `json:",inline"`

	// Policies holds the policies will bind or unbind to the role.
	// +optional
	Policies []string `json:"policies" protobuf:"bytes,1,rep,name=policies"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Group represents a group of users.
type Group struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of group document in this set.
	Spec GroupSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status GroupStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GroupList is the whole list of all groups.
type GroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of rules.
	Items []Group `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// GroupPhase defines the phase of group constructor.
type GroupPhase string

const (
	GroupActive GroupPhase = "Active"
	// GroupTerminating means the group is undergoing graceful termination.
	GroupTerminating GroupPhase = "Terminating"
)

// GroupSpec is a description of group.
type GroupSpec struct {
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=FinalizerName"`

	DisplayName string `json:"displayName" protobuf:"bytes,2,opt,name=displayName"`
	TenantID    string `json:"tenantID" protobuf:"bytes,3,opt,name=tenantID"`

	//Creator
	Username    string `json:"username" protobuf:"bytes,4,opt,name=username"`
	Description string `json:"description" protobuf:"bytes,5,opt,name=description"`
}

// GroupStatus represents information about the status of a group.
type GroupStatus struct {
	// +optional
	Phase GroupPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=GroupPhase"`

	// Subjects represents the members of the group.
	Subjects []Subject `json:"subjects" protobuf:"bytes,2,rep,name=subjects"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMap holds configuration data for tke to consume.
type ConfigMap struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`

	// BinaryData contains the binary data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// BinaryData can contain byte sequences that are not in the UTF-8 range.
	// The keys stored in BinaryData must not overlap with the ones in
	// the Data field, this is enforced during validation process.
	// +optional
	BinaryData map[string][]byte `json:"binaryData,omitempty" protobuf:"bytes,3,rep,name=binaryData"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapList is a resource containing a list of ConfigMap objects.
type ConfigMapList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of ConfigMaps.
	Items []ConfigMap `json:"items" protobuf:"bytes,2,rep,name=items"`
}
