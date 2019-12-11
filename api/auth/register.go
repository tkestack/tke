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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	// Scheme is the default instance of runtime.Scheme to which types in the TKE API are already registered.
	Scheme = runtime.NewScheme()
	// Codecs provides access to encoding and decoding for the scheme
	Codecs = serializer.NewCodecFactory(Scheme)
	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

// GroupName is group name used to register these schema
const GroupName = "auth.tkestack.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns back a Group qualified
// GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder collects functions that add things to a scheme.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme applies all the stored functions to the scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// addKnownTypes adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&LocalIdentity{},
		&LocalIdentityList{},
		&PasswordReq{},
		&APIKey{},
		&APIKeyList{},
		&APIKeyReq{},
		&APIKeyReqPassword{},
		&APISigningKey{},
		&APISigningKeyList{},
		&Category{},
		&CategoryList{},
		&Policy{},
		&PolicyList{},
		&Rule{},
		&RuleList{},
		&Binding{},
		&Role{},
		&RoleList{},
		&PolicyBinding{},
		&SubjectAccessReview{},
		&Group{},
		&GroupList{},

		&ConfigMap{},
		&ConfigMapList{})

	return nil
}
