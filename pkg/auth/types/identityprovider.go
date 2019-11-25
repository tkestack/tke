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

package types

// IdentityProvider is an object that contains the metadata about OIDC identity
// provider used to login to TKE.
//+k8s:openapi-gen=true
type IdentityProvider struct {
	// ID that will uniquely identify the connector object and will be used as tenantID.
	ID string `json:"id,omitempty"`
	// The Name of the connector that is used when displaying it to the end user.
	Name string `json:"name,omitempty"`
	// The type of the connector. E.g. 'oidc' or 'ldap'
	Type string `json:"type,omitempty"`
	// ResourceVersion is the static versioning used to keep track of dynamic configuration
	// changes to the connector object made by the API calls.
	ResourceVersion string `json:"resourceVersion,omitempty"`
	// Config holds all the configuration information specific to the connector type. Since there
	// no generic struct we can use for this purpose, it is stored as a json string.
	Config string `json:"config,omitempty"`
}

// IdentityProviderList is the whole list of IdentityProvider.
//+k8s:openapi-gen=true
type IdentityProviderList struct {
	// List of policies.
	Items []*IdentityProvider `json:"items,omitempty"`
}
