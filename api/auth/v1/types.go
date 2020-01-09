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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// KeywordQueryTag is a field tag to query object that contains the keyword.
	KeywordQueryTag string = "keyword"

	// QueryLimitTag is a field tag to query a maximum number of objects for a list call.
	QueryLimitTag string = "limit"

	// IssuerName is the name of issuer location.
	IssuerName = "oidc"
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
	GroupFinalize FinalizerName = "localgroup"

	// RoleFinalize is an internal finalizer values to Role.
	RoleFinalize FinalizerName = "role"
)

// LocalIdentitySpec is a description of an identity.
type LocalIdentitySpec struct {
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,11,rep,name=finalizers,casttype=FinalizerName"`

	Username       string `json:"username" protobuf:"bytes,7,opt,name=name"`
	DisplayName    string `json:"displayName" protobuf:"bytes,8,opt,name=displayName"`
	Email          string `json:"email" protobuf:"bytes,9,opt,name=email"`
	PhoneNumber    string `json:"phoneNumber" protobuf:"bytes,10,opt,name=phone"`
	HashedPassword string `json:"hashedPassword,omitempty" protobuf:"bytes,4,opt,name=hashedPassword"`
	// +optional
	TenantID string `json:"tenantID,omitempty" protobuf:"bytes,2,opt,name=tenantID"`

	Groups []string `json:"groups,omitempty" protobuf:"bytes,6,rep,name=groups"`
	// +optional
	Extra map[string]string `json:"extra,omitempty" protobuf:"bytes,3,rep,name=extra"`
}

// LocalIdentityPhase defines the phase of LocalIdentity construct.
type LocalIdentityPhase string

const (
	// LocalIdentityDeleting means the local identity is undergoing graceful termination.
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PasswordReq contains info to update password for a localIdentity
type PasswordReq struct {
	metav1.TypeMeta `json:",inline"`

	HashedPassword   string `json:"hashedPassword,omitempty" protobuf:"bytes,1,opt,name=hashedPassword"`
	OriginalPassword string `json:"originalPassword,omitempty" protobuf:"bytes,2,opt,name=originalPassword"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalGroup represents a group of users.
type LocalGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of group document in this set.
	Spec LocalGroupSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status LocalGroupStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalGroupList is the whole list of all groups.
type LocalGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of LocalGroup.
	Items []LocalGroup `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// GroupPhase defines the phase of group constructor.
type GroupPhase string

const (
	GroupActive GroupPhase = "Active"
	// GroupTerminating means the group is undergoing graceful termination.
	GroupTerminating GroupPhase = "Terminating"
)

// LocalGroupSpec is a description of group.
type LocalGroupSpec struct {
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=FinalizerName"`

	DisplayName string `json:"displayName" protobuf:"bytes,2,opt,name=displayName"`
	TenantID    string `json:"tenantID" protobuf:"bytes,3,opt,name=tenantID"`

	Username    string `json:"username" protobuf:"bytes,4,opt,name=username"`
	Description string `json:"description" protobuf:"bytes,5,opt,name=description"`
}

// LocalGroupStatus represents information about the status of a group.
type LocalGroupStatus struct {
	// +optional
	Phase GroupPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=GroupPhase"`

	// Users represents the members of the group.
	Users []Subject `json:"users" protobuf:"bytes,2,rep,name=users"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// User is an object that contains the metadata about identify about tke local idp or third-party idp.
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec defines the desired identities of identity in this set.
	Spec UserSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// UserSpec is a description of an user.
type UserSpec struct {
	ID string `json:"id" protobuf:"bytes,1,opt,name=id"`

	// Name must be unique in the same tenant.
	Name        string            `json:"name" protobuf:"bytes,2,opt,name=name"`
	DisplayName string            `json:"displayName,omitempty" protobuf:"bytes,3,opt,name=displayName"`
	Email       string            `json:"email,omitempty" protobuf:"bytes,4,opt,name=email"`
	PhoneNumber string            `json:"phoneNumber,omitempty" protobuf:"bytes,5,opt,name=phoneNumber"`
	TenantID    string            `json:"tenantID,omitempty" protobuf:"bytes,6,opt,name=tenantID"`
	Extra       map[string]string `json:"extra,omitempty" protobuf:"bytes,7,rep,name=extra"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UserList is the whole list of all users.
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of User.
	Items []User `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Group is an object that contains the metadata about identify about tke local idp or third-party idp.
type Group struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec defines the desired identities of group in this set.
	Spec GroupSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	Status GroupStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// GroupSpec is a description of a Group.
type GroupSpec struct {
	ID          string `json:"id" protobuf:"bytes,1,opt,name=id"`
	DisplayName string `json:"displayName" protobuf:"bytes,2,opt,name=displayName"`
	TenantID    string `json:"tenantID" protobuf:"bytes,3,opt,name=tenantID"`
	Description string `json:"description" protobuf:"bytes,4,opt,name=description"`
}

// GroupStatus represents information about the status of a group.
type GroupStatus struct {
	// Users represents the members of the group.
	Users []Subject `json:"users" protobuf:"bytes,2,rep,name=users"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GroupList is the whole list of all groups.
type GroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of Group.
	Items []Group `json:"items" protobuf:"bytes,2,rep,name=items"`
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
	// Expired represents whether the apikey has been expired.
	Expired bool `json:"expired" protobuf:"varint,2,opt,name=expired"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReq contains expiration time used to apply the api key.
type APIKeyReq struct {
	metav1.TypeMeta `json:",inline"`

	// Expire is required, holds the duration of the api key become invalid. By default, 168h(= seven days)
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
	// Users represents the users the policy applies to.
	Users []Subject `json:"users" protobuf:"bytes,2,rep,name=users"`

	// +optional
	// Groups represents the groups the policy applies to.
	Groups []Subject `json:"groups" protobuf:"bytes,3,rep,name=groups"`
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

	// Users holds references to the objects the policy applies to.
	// +optional
	Users []Subject `json:"users, omitempty" protobuf:"bytes,1,rep,name=users"`

	// Groups holds references to the groups the policy applies to.
	// +optional
	Groups []Subject `json:"groups,omitempty" protobuf:"bytes,2,rep,name=groups"`
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

	// Username is Creator
	Username    string `json:"username" protobuf:"bytes,4,opt,name=username"`
	Description string `json:"description" protobuf:"bytes,5,opt,name=description"`

	Policies []string `json:"policies" protobuf:"bytes,6,rep,name=policies"`
}

// RoleStatus represents information about the status of a role.
type RoleStatus struct {
	// +optional
	Phase RolePhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=RolePhase"`

	// +optional
	// Users represents the users the role applies to.
	Users []Subject `json:"users" protobuf:"bytes,2,rep,name=users"`

	// +optional
	// Groups represents the groups the role applies to.
	Groups []Subject `json:"groups" protobuf:"bytes,3,rep,name=groups"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyBinding references the request to bind or unbind policies to the role.
type PolicyBinding struct {
	metav1.TypeMeta `json:",inline"`

	// Policies holds the policies will bind or unbind to the role.
	// +optional
	Policies []string `json:"policies" protobuf:"bytes,1,rep,name=policies"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SubjectAccessReview checks whether or not a user or group can perform an action.  Not filling in a
// spec.namespace means "in all namespaces".
type SubjectAccessReview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec holds information about the request being evaluated
	Spec SubjectAccessReviewSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// Status is filled in by the server and indicates whether the request is allowed or not
	Status SubjectAccessReviewStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// SubjectAccessReviewSpec is a description of the access request.  Exactly one of ResourceAttributes
// and NonResourceAttributes must be set
type SubjectAccessReviewSpec struct {
	// ResourceAttributes describes information for a resource access request
	ResourceAttributes *ResourceAttributes `json:"resourceAttributes,omitempty" protobuf:"bytes,7,opt,name=resourceAttributes"`

	// ResourceAttributesList describes information for multi resource access request.
	ResourceAttributesList []*ResourceAttributes `json:"resourceAttributesList,omitempty" protobuf:"bytes,8,rep,name=resourceAttributesList"`

	// NonResourceAttributes describes information for a non-resource access request
	NonResourceAttributes *NonResourceAttributes `json:"nonResourceAttributes,omitempty" protobuf:"bytes,9,opt,name=nonResourceAttributes"`

	// User is the user you're testing for.
	// If you specify "User" but not "Groups", then is it interpreted as "What if User were not a member of any groups
	// +optional
	User string `json:"user,omitempty" protobuf:"bytes,3,opt,name=user"`
	// Groups is the groups you're testing for.
	// +optional
	Groups []string `json:"groups,omitempty" protobuf:"bytes,4,rep,name=groups"`
	// Extra corresponds to the user.Info.GetExtra() method from the authenticator.  Since that is input to the authorizer
	// it needs a reflection here.
	// +optional
	Extra map[string]ExtraValue `json:"extra,omitempty" protobuf:"bytes,5,rep,name=extra"`
	// UID information about the requesting user.
	// +optional
	UID string `json:"uid,omitempty" protobuf:"bytes,6,opt,name=uid"`
}

// ExtraValue masks the value so protobuf can generate
// +protobuf.nullable=true
// +protobuf.options.(gogoproto.goproto_stringer)=false
type ExtraValue []string

func (t ExtraValue) String() string {
	return fmt.Sprintf("%v", []string(t))
}

// ResourceAttributes includes the authorization attributes available for resource requests to the Authorizer interface
type ResourceAttributes struct {
	// Namespace is the namespace of the action being requested.  Currently, there is no distinction between no namespace and all namespaces
	// "" (empty) is defaulted for LocalSubjectAccessReviews
	// "" (empty) is empty for cluster-scoped resources
	// "" (empty) means "all" for namespace scoped resources from a SubjectAccessReview or SelfSubjectAccessReview
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	// Verb is a kubernetes resource API verb, like: get, list, watch, create, update, delete, proxy.  "*" means all.
	// +optional
	Verb string `json:"verb,omitempty" protobuf:"bytes,2,opt,name=verb"`
	// Group is the API Group of the Resource.  "*" means all.
	// +optional
	Group string `json:"group,omitempty" protobuf:"bytes,3,opt,name=group"`
	// Version is the API Version of the Resource.  "*" means all.
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,4,opt,name=version"`
	// Resource is one of the existing resource types.  "*" means all.
	// +optional
	Resource string `json:"resource,omitempty" protobuf:"bytes,5,opt,name=resource"`
	// Subresource is one of the existing resource types.  "" means none.
	// +optional
	Subresource string `json:"subresource,omitempty" protobuf:"bytes,6,opt,name=subresource"`
	// Name is the name of the resource being requested for a "get" or deleted for a "delete". "" (empty) means all.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,7,opt,name=name"`
}

// NonResourceAttributes includes the authorization attributes available for non-resource requests to the Authorizer interface
type NonResourceAttributes struct {
	// Path is the URL path of the request
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
	// Verb is the standard HTTP verb
	// +optional
	Verb string `json:"verb,omitempty" protobuf:"bytes,2,opt,name=verb"`
}

// SubjectAccessReviewStatus represents the current state of a SubjectAccessReview.
type SubjectAccessReviewStatus struct {
	// Allowed is required. True if the action would be allowed, false otherwise.
	Allowed bool `json:"allowed" protobuf:"varint,1,opt,name=allowed"`
	// Denied is optional. True if the action would be denied, otherwise
	// false. If both allowed is false and denied is false, then the
	// authorizer has no opinion on whether to authorize the action. Denied
	// may not be true if Allowed is true.
	// +optional
	Denied bool `json:"denied,omitempty" protobuf:"varint,4,opt,name=denied"`
	// Reason is optional.  It indicates why a request was allowed or denied.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	// EvaluationError is an indication that some error occurred during the authorization check.
	// It is entirely possible to get an error and be able to continue determine authorization status in spite of it.
	// For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.
	// +optional
	EvaluationError string `json:"evaluationError,omitempty" protobuf:"bytes,3,opt,name=evaluationError"`

	// AllowedList is the allowed response for batch authorization request.
	AllowedList []*AllowedStatus `json:"allowedList,omitempty" protobuf:"bytes,5,rep,name=allowedList"`
}

// AllowedStatus includes the resource access request and response.
// +k8s:openapi-gen=true
type AllowedStatus struct {
	// Resource is the resource of request
	Resource string `json:"resource" protobuf:"bytes,1,opt,name=resource"`
	// Verb is the verb of request
	Verb string `json:"web" protobuf:"bytes,2,opt,name=web"`

	// Allowed is required. True if the action would be allowed, false otherwise.
	Allowed bool `json:"allowed" protobuf:"varint,3,opt,name=allowed"`
	// Denied is optional. True if the action would be denied, otherwise
	// false. If both allowed is false and denied is false, then the
	// authorizer has no opinion on whether to authorize the action. Denied
	// may not be true if Allowed is true.
	Denied bool `json:"denied,omitempty" protobuf:"varint,4,opt,name=denied"`
	// Reason is optional.  It indicates why a request was allowed or denied.
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// EvaluationError is an indication that some error occurred during the authorization check.
	// It is entirely possible to get an error and be able to continue determine authorization status in spite of it.
	// For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.
	EvaluationError string `json:"evaluationError,omitempty" protobuf:"bytes,6,opt,name=evaluationError"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IdentityProvider is an object that contains the metadata about identify
// provider used to login to TKE.
type IdentityProvider struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of identity provider in this set.
	Spec IdentityProviderSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IdentityProviderList is the whole list of all identity providers.
type IdentityProviderList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of identity providers.
	Items []IdentityProvider `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// IdentityProviderSpec is a description of an identity provider.
type IdentityProviderSpec struct {
	// The Name of the connector that is used when displaying it to the end user.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// The type of the connector. E.g. 'oidc' or 'ldap'
	Type string `json:"type" protobuf:"bytes,2,opt,name=type"`
	// The admins means the users is super admin for the idp.
	Admins []string `json:"admin" protobuf:"bytes,3,opt,name=admins"`
	// Config holds all the configuration information specific to the connector type. Since there
	// no generic struct we can use for this purpose, it is stored as a json string.
	Config string `json:"config" protobuf:"bytes,4,opt,name=config"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Client represents an OAuth2 client.
type Client struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired identities of identity provider in this set.
	Spec ClientSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClientList is the whole list of OAuth2 client.
type ClientList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of identity providers.
	Items []Client `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ClientSpec is a description of an client.
type ClientSpec struct {
	ID           string   `json:"id,omitempty" protobuf:"bytes,1,opt,name=id"`
	Secret       string   `json:"secret,omitempty" protobuf:"bytes,2,opt,name=secret"`
	RedirectUris []string `json:"redirect_uris,omitempty" protobuf:"bytes,3,rep,name=redirect_uris,json=redirectUris"`
	// TrustedPeers are a list of peers which can issue tokens on this client's behalf using the dynamic "oauth2:server:client_id:(client_id)" scope.
	TrustedPeers []string `json:"trusted_peers,omitempty" protobuf:"bytes,4,rep,name=trusted_peers,json=trustedPeers"`
	// Public clients must use either use a redirectURL 127.0.0.1:X or "urn:ietf:wg:oauth:2.0:oob".
	Public  bool   `json:"public,omitempty" protobuf:"varint,5,opt,name=public"`
	Name    string `json:"name,omitempty" protobuf:"bytes,6,opt,name=name"`
	LogoURL string `json:"logo_url,omitempty" protobuf:"bytes,7,opt,name=logo_url,json=logoUrl"`
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
