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

	// PolicyFinalize is an internal finalizer values to Policy.
	PolicyFinalize FinalizerName = "policy"

	// PolicyFinalize is an internal finalizer values to LocalGroup.
	LocalGroupFinalize FinalizerName = "localgroup"

	// RoleFinalize is an internal finalizer values to Role.
	RoleFinalize FinalizerName = "role"
)

// LocalIdentitySpec is a description of an identity.
type LocalIdentitySpec struct {
	Finalizers []FinalizerName

	Username       string
	DisplayName    string
	Email          string
	PhoneNumber    string
	HashedPassword string
	TenantID       string
	Groups         []string
	Extra          map[string]string
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PasswordReq contains info to update password for a localIdentity
type PasswordReq struct {
	metav1.TypeMeta

	HashedPassword   string
	OriginalPassword string
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalGroup represents a group of users.
type LocalGroup struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Spec defines the desired identities of group document in this set.
	Spec LocalGroupSpec

	// +optional
	Status LocalGroupStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LocalGroupList is the whole list of all groups.
type LocalGroupList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of localgroup.
	Items []LocalGroup
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
	Finalizers []FinalizerName

	DisplayName string
	TenantID    string

	// Username is Creator
	Username    string
	Description string
}

// LocalGroupStatus represents information about the status of a group.
type LocalGroupStatus struct {
	// +optional
	Phase GroupPhase

	// Users represents the members of the group.
	Users []Subject
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// User is an object that contains the metadata about identify about tke local idp or third-party idp.
type User struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	// Spec defines the desired identities of identity in this set.
	Spec UserSpec
}

// UserSpec is a description of an user.
type UserSpec struct {
	ID string

	// Name must be unique in the same tenant.
	Name        string
	DisplayName string
	Email       string
	PhoneNumber string
	TenantID    string
	Extra       map[string]string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UserList is the whole list of all users.
type UserList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of User.
	Items []User
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Group is an object that contains the metadata about identify about tke local idp or third-party idp.
type Group struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	// Spec defines the desired identities of group in this set.
	Spec GroupSpec

	Status GroupStatus
}

// GroupSpec is a description of an Group.
type GroupSpec struct {
	ID          string
	DisplayName string
	TenantID    string
	Description string
}

// GroupStatus represents information about the status of a group.
type GroupStatus struct {
	// Users represents the members of the group.
	Users []Subject
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GroupList is the whole list of all groups.
type GroupList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of group.
	Items []Group
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

// APIKeyList is the whole list of all identities.
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
	// Expired represents whether the apikey has been expired.
	Expired bool `json:"expired"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIKeyReq contains expiration time used to apply the api key.
type APIKeyReq struct {
	metav1.TypeMeta

	// Expire is required, holds the duration of the api key become invalid. By default, 168h(= seven days)
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

	DisplayName string
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
	// Users represents the users the policy applies to.
	Users []Subject

	// +optional
	// Groups represents the groups the policy applies to.
	Groups []Subject
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

// RuleList is the whole list of all rules.
type RuleList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of rules.
	Items []Rule
}

// RuleSpec is a description of a rule.
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

	// Users holds references to the objects the policy applies to.
	// +optional
	Users []Subject

	// Groups holds references to the groups the policy applies to.
	// +optional
	Groups []Subject
}

// Subject references a user can specify by id or name.
type Subject struct {
	ID   string
	Name string
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Role is a collection with multiple policies.
type Role struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Spec defines the desired identities of role document in this set.
	Spec RoleSpec

	// +optional
	Status RoleStatus
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RoleList is the whole list of policy.
type RoleList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of rules.
	Items []Role
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
	Finalizers []FinalizerName

	DisplayName string
	TenantID    string

	// Username is Creator
	Username    string
	Description string

	Policies []string
}

// RoleStatus represents information about the status of a role.
type RoleStatus struct {
	// +optional
	Phase RolePhase

	// Users represents the users of the applies to.
	Users []Subject

	// +optional
	// Groups represents the groups the policy applies to.
	Groups []Subject
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyBinding references the request to bind or unbind policies to the role.
type PolicyBinding struct {
	metav1.TypeMeta

	// Policies holds the policies will bind or unbind to the role.
	// +optional
	Policies []string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SubjectAccessReview checks whether or not a user or group can perform an action.  Not filling in a
// spec.namespace means "in all namespaces".
type SubjectAccessReview struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Spec holds information about the request being evaluated
	Spec SubjectAccessReviewSpec

	// Status is filled in by the server and indicates whether the request is allowed or not
	Status SubjectAccessReviewStatus
}

// SubjectAccessReviewSpec is a description of the access request.  Exactly one of ResourceAttributes
// and NonResourceAttributes must be set
type SubjectAccessReviewSpec struct {
	// ResourceAttributes describes information for a resource access request
	ResourceAttributes *ResourceAttributes

	// ResourceAttributesList describes information for multi resource access request.
	ResourceAttributesList []*ResourceAttributes

	// NonResourceAttributes describes information for a non-resource access request
	NonResourceAttributes *NonResourceAttributes

	// User is the user you're testing for.
	// If you specify "User" but not "Group", then is it interpreted as "What if User were not a member of any groups
	User string
	// Groups is the groups you're testing for.
	Groups []string
	// Extra corresponds to the user.Info.GetExtra() method from the authenticator.  Since that is input to the authorizer
	// it needs a reflection here.
	Extra map[string]ExtraValue
	// UID information about the requesting user.
	UID string
}

// ExtraValue masks the value so protobuf can generate
// +protobuf.nullable=true
type ExtraValue []string

// ResourceAttributes includes the authorization attributes available for resource requests to the Authorizer interface
type ResourceAttributes struct {
	// Namespace is the namespace of the action being requested.  Currently, there is no distinction between no namespace and all namespaces
	// "" (empty) is defaulted for LocalSubjectAccessReviews
	// "" (empty) is empty for cluster-scoped resources
	// "" (empty) means "all" for namespace scoped resources from a SubjectAccessReview or SelfSubjectAccessReview
	Namespace string
	// Verb is a kubernetes resource API verb, like: get, list, watch, create, update, delete, proxy.  "*" means all.
	Verb string
	// Group is the API Group of the Resource.  "*" means all.
	Group string
	// Version is the API Version of the Resource.  "*" means all.
	Version string
	// Resource is one of the existing resource types.  "*" means all.
	Resource string
	// Subresource is one of the existing resource types.  "" means none.
	Subresource string
	// Name is the name of the resource being requested for a "get" or deleted for a "delete". "" (empty) means all.
	Name string
}

// NonResourceAttributes includes the authorization attributes available for non-resource requests to the Authorizer interface
type NonResourceAttributes struct {
	// Path is the URL path of the request
	Path string
	// Verb is the standard HTTP verb
	Verb string
}

// SubjectAccessReviewStatus represents the current state of a SubjectAccessReview.
type SubjectAccessReviewStatus struct {
	// Allowed is required. True if the action would be allowed, false otherwise.
	Allowed bool
	// Denied is optional. True if the action would be denied, otherwise
	// false. If both allowed is false and denied is false, then the
	// authorizer has no opinion on whether to authorize the action. Denied
	// may not be true if Allowed is true.
	Denied bool
	// Reason is optional.  It indicates why a request was allowed or denied.
	Reason string
	// EvaluationError is an indication that some error occurred during the authorization check.
	// It is entirely possible to get an error and be able to continue determine authorization status in spite of it.
	// For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.
	EvaluationError string

	// AllowedList is the allowed response for batch authorization request.
	AllowedList []*AllowedStatus
}

// AllowedStatus includes the resource access request and response.
// +k8s:openapi-gen=true
type AllowedStatus struct {
	// Resource is the resource of request
	Resource string
	// Verb is the verb of request
	Verb string

	// Allowed is required. True if the action would be allowed, false otherwise.
	Allowed bool
	// Denied is optional. True if the action would be denied, otherwise
	// false. If both allowed is false and denied is false, then the
	// authorizer has no opinion on whether to authorize the action. Denied
	// may not be true if Allowed is true.
	Denied bool
	// Reason is optional.  It indicates why a request was allowed or denied.
	Reason string
	// EvaluationError is an indication that some error occurred during the authorization check.
	// It is entirely possible to get an error and be able to continue determine authorization status in spite of it.
	// For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.
	EvaluationError string
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IdentityProvider is an object that contains the metadata about identify
// provider used to login to TKE.
type IdentityProvider struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Spec defines the desired identities of identity provider in this set.
	Spec IdentityProviderSpec
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IdentityProviderList is the whole list of all identity providers.
type IdentityProviderList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of identity providers.
	Items []IdentityProvider
}

// IdentityProviderSpec is a description of an identity provider.
type IdentityProviderSpec struct {
	// The Name of the connector that is used when displaying it to the end user.
	Name string
	// The type of the connector. E.g. 'oidc' or 'ldap'
	Type string

	// The admins means the users is super admin for the idp.
	Administrators []string
	// Config holds all the configuration information specific to the connector type. Since there
	// no generic struct we can use for this purpose, it is stored as a json string.
	Config string
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Client represents an OAuth2 client.
type Client struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Spec defines the desired identities of identity provider in this set.
	Spec ClientSpec
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClientList is the whole list of OAuth2 client.
type ClientList struct {
	metav1.TypeMeta
	metav1.ListMeta
	// List of identity providers.
	Items []Client
}

// ClientSpec is a description of an client.
type ClientSpec struct {
	ID           string
	Secret       string
	RedirectUris []string
	// TrustedPeers are a list of peers which can issue tokens on this client's behalf using the dynamic "oauth2:server:client_id:(client_id)" scope.
	TrustedPeers []string
	// Public clients must use either use a redirectURL 127.0.0.1:X or "urn:ietf:wg:oauth:2.0:oob".
	Public  bool
	Name    string
	LogoURL string
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
