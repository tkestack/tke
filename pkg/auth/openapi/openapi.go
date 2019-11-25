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

package openapi

import (
	"github.com/go-openapi/spec"
	"k8s.io/kube-openapi/pkg/common"
	openapicommon "tkestack.io/tke/api/openapi"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	definitions := openapicommon.GetOpenAPIDefinitions(ref)

	definitions["tkestack.io/tke/pkg/auth/types.APIKeyData"] = schemaAPIKeyData(ref)
	definitions["tkestack.io/tke/pkg/auth/types.APIKeyList"] = schemaAPIKeyList(ref)
	definitions["tkestack.io/tke/pkg/auth/types.APIKeyReq"] = schemaAPIKeyReq(ref)
	definitions["tkestack.io/tke/pkg/auth/types.APIKeyReqPassword"] = schemaAPIKeyReqPassword(ref)
	definitions["tkestack.io/tke/pkg/auth/types.Action"] = schemaAction(ref)
	definitions["tkestack.io/tke/pkg/auth/types.AllowedResponse"] = schemaAllowedResponse(ref)
	definitions["tkestack.io/tke/pkg/auth/types.AttachInfo"] = schemaAttachInfo(ref)
	definitions["tkestack.io/tke/pkg/auth/types.Category"] = schemaCategory(ref)
	definitions["tkestack.io/tke/pkg/auth/types.CategoryList"] = schemaCategoryList(ref)
	definitions["tkestack.io/tke/pkg/auth/types.Client"] = schemaClient(ref)
	definitions["tkestack.io/tke/pkg/auth/types.ClientList"] = schemaClientList(ref)
	definitions["tkestack.io/tke/pkg/auth/types.Duration"] = schemaDuration(ref)
	definitions["tkestack.io/tke/pkg/auth/types.IdentityProvider"] = schemaIdentityProvider(ref)
	definitions["tkestack.io/tke/pkg/auth/types.IdentityProviderList"] = schemaIdentityProviderList(ref)
	definitions["tkestack.io/tke/pkg/auth/types.LocalIdentity"] = schemaLocalIdentity(ref)
	definitions["tkestack.io/tke/pkg/auth/types.LocalIdentityList"] = schemaLocalIdentityList(ref)
	definitions["tkestack.io/tke/pkg/auth/types.LocalIdentitySpec"] = schemaLocalIdentitySpec(ref)
	definitions["tkestack.io/tke/pkg/auth/types.LocalIdentityStatus"] = schemaLocalIdentityStatus(ref)
	definitions["tkestack.io/tke/pkg/auth/types.NonResourceAttributes"] = schemaNonResourceAttributes(ref)
	definitions["tkestack.io/tke/pkg/auth/types.Permission"] = schemaPermission(ref)
	definitions["tkestack.io/tke/pkg/auth/types.Policy"] = schemaPolicy(ref)
	definitions["tkestack.io/tke/pkg/auth/types.PolicyCreate"] = schemaPolicyCreate(ref)
	definitions["tkestack.io/tke/pkg/auth/types.PolicyList"] = schemaPolicyList(ref)
	definitions["tkestack.io/tke/pkg/auth/types.PolicyMeta"] = schemaPolicyMeta(ref)
	definitions["tkestack.io/tke/pkg/auth/types.PolicyOption"] = schemaPolicyOption(ref)
	definitions["tkestack.io/tke/pkg/auth/types.ResourceAttributes"] = schemaResourceAttributes(ref)
	definitions["tkestack.io/tke/pkg/auth/types.Role"] = schemaRole(ref)
	definitions["tkestack.io/tke/pkg/auth/types.RoleList"] = schemaRoleList(ref)
	definitions["tkestack.io/tke/pkg/auth/types.RoleOption"] = schemaRoleOption(ref)
	definitions["tkestack.io/tke/pkg/auth/types.Statement"] = schemaStatement(ref)
	definitions["tkestack.io/tke/pkg/auth/types.SubjectAccessReview"] = schemaSubjectAccessReview(ref)
	definitions["tkestack.io/tke/pkg/auth/types.SubjectAccessReviewSpec"] = schemaSubjectAccessReviewSpec(ref)
	definitions["tkestack.io/tke/pkg/auth/types.SubjectAccessReviewStatus"] = schemaSubjectAccessReviewStatus(ref)
	definitions["tkestack.io/tke/pkg/auth/types.TokenReviewRequest"] = schemaTokenReviewRequest(ref)
	definitions["tkestack.io/tke/pkg/auth/types.TokenReviewResponse"] = schemaTokenReviewResponse(ref)
	definitions["tkestack.io/tke/pkg/auth/types.TokenReviewSpec"] = schemaTokenReviewSpec(ref)
	definitions["tkestack.io/tke/pkg/auth/types.TokenReviewStatus"] = schemaTokenReviewStatus(ref)
	definitions["tkestack.io/tke/pkg/auth/types.TokenReviewUser"] = schemaTokenReviewUser(ref)

	return definitions
}

func schemaAPIKeyData(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "APIKeyData contains expiration time used to apply the api key.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"apiKey": {
						SchemaProps: spec.SchemaProps{
							Description: "APIkey is the jwt token used to authenticate user, and contains user info and sign.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"disabled": {
						SchemaProps: spec.SchemaProps{
							Description: "Disabled represents whether the apikey has been disabled.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"deleted": {
						SchemaProps: spec.SchemaProps{
							Description: "Deleted represents whether the apikey has been deleted.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Description: "Description describes api keys usage.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"issue_at": {
						SchemaProps: spec.SchemaProps{
							Description: "IssueAt is the created time for api key",
							Type:        []string{"string"},
							Format:      "date-time",
						},
					},
					"expire_at": {
						SchemaProps: spec.SchemaProps{
							Description: "ExpireAt is the expire time for api key",
							Type:        []string{"string"},
							Format:      "date-time",
						},
					},
				},
			},
		},
	}
}

func schemaAPIKeyList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "APIKeyList is the whole list of APIKeyData.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"items": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.APIKeyData"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.APIKeyData"},
	}
}

func schemaAPIKeyReq(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "APIKeyReq contains expiration time used to apply the api key.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"expire": {
						SchemaProps: spec.SchemaProps{
							Description: "Exipre holds the duration of the api key become invalid. By default, 168h(= seven days)",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Description: "Description describes api keys usage.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schemaAPIKeyReqPassword(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "APIKeyReqPassword contains userinfo and expiration time used to apply the api key.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Description: "TenantID for user",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"username": {
						SchemaProps: spec.SchemaProps{
							Description: "UserName",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"password": {
						SchemaProps: spec.SchemaProps{
							Description: "Password (encoded by base64)",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"expire": {
						SchemaProps: spec.SchemaProps{
							Description: "Exipre holds the duration of the api key become invalid. By default, 168h(= seven days)",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Description: "Description describes api keys usage.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schemaAction(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Action defines a action verb for authorization.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Description: "Name represents user access review request verb.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Description: "Description describes the action.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schemaAllowedResponse(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AllowedResponse includes the resource access request and response.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"resource": {
						SchemaProps: spec.SchemaProps{
							Description: "Path is the URL path of the request",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"verb": {
						SchemaProps: spec.SchemaProps{
							Description: "Verb is the standard HTTP verb",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"allowed": {
						SchemaProps: spec.SchemaProps{
							Description: "Allowed is required. True if the action would be allowed, false otherwise.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"denied": {
						SchemaProps: spec.SchemaProps{
							Description: "Denied is optional. True if the action would be denied, otherwise false. If both allowed is false and denied is false, then the authorizer has no opinion on whether to authorize the action. Denied may not be true if Allowed is true.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"reason": {
						SchemaProps: spec.SchemaProps{
							Description: "Reason is optional.  It indicates why a request was allowed or denied.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"evaluationError": {
						SchemaProps: spec.SchemaProps{
							Description: "EvaluationError is an indication that some error occurred during the authorization check. It is entirely possible to get an error and be able to continue determine authorization status in spite of it. For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"allowed"},
			},
		},
	}
}

func schemaAttachInfo(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AttachInfo contains info to attach/detach users to/from policy or role.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"id": {
						SchemaProps: spec.SchemaProps{
							Description: "role or policy id bond",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"userNames": {
						SchemaProps: spec.SchemaProps{
							Description: "name of users",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"policyIDs": {
						SchemaProps: spec.SchemaProps{
							Description: "id of policies",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Description: "id of tenant",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"id", "userNames", "policyIDs", "tenantID"},
			},
		},
	}
}

func schemaCategory(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Category defines a category of actions for policy.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Description: "Name identifies policy category",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"displayName": {
						SchemaProps: spec.SchemaProps{
							Description: "DisplayName used to display category name",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"actions": {
						SchemaProps: spec.SchemaProps{
							Description: "Actions represents a series of actions work on the policy category",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.Action"),
									},
								},
							},
						},
					},
					"createAt": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "date-time",
						},
					},
					"updateAt": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "date-time",
						},
					},
				},
				Required: []string{"actions"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.Action"},
	}
}

func schemaCategoryList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "CategoryList is the whole list of policy Category.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"items": {
						SchemaProps: spec.SchemaProps{
							Description: "List of category.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.Category"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.Category"},
	}
}

func schemaClient(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Client represents an OAuth2 client.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"id": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"secret": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"redirect_uris": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"trusted_peers": {
						SchemaProps: spec.SchemaProps{
							Description: "TrustedPeers are a list of peers which can issue tokens on this client's behalf using the dynamic \"oauth2:server:client_id:(client_id)\" scope.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"public": {
						SchemaProps: spec.SchemaProps{
							Description: "Public clients must use either use a redirectURL 127.0.0.1:X or \"urn:ietf:wg:oauth:2.0:oob\".",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"logo_url": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
			},
		},
	}
}

func schemaClientList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ClientList is the whole list of OAuth2 client.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"items": {
						SchemaProps: spec.SchemaProps{
							Description: "List of policies.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.Client"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.Client"},
	}
}

func schemaDuration(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Duration implements json marshal func for time.Duration.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"Duration": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int64",
						},
					},
				},
				Required: []string{"Duration"},
			},
		},
	}
}

func schemaIdentityProvider(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "IdentityProvider is an object that contains the metadata about OIDC identity provider used to login to TKE.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"id": {
						SchemaProps: spec.SchemaProps{
							Description: "ID that will uniquely identify the connector object and will be used as tenantID.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"name": {
						SchemaProps: spec.SchemaProps{
							Description: "The Name of the connector that is used when displaying it to the end user.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"type": {
						SchemaProps: spec.SchemaProps{
							Description: "The type of the connector. E.g. 'oidc' or 'ldap'",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"resourceVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "ResourceVersion is the static versioning used to keep track of dynamic configuration changes to the connector object made by the API calls.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"config": {
						SchemaProps: spec.SchemaProps{
							Description: "Config holds all the configuration information specific to the connector type. Since there no generic struct we can use for this purpose, it is stored as a json string.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schemaIdentityProviderList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "IdentityProviderList is the whole list of IdentityProvider.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"items": {
						SchemaProps: spec.SchemaProps{
							Description: "List of policies.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.IdentityProvider"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.IdentityProvider"},
	}
}

func schemaLocalIdentity(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LocalIdentity is an object that contains the metadata about identify used to login to TKE.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"uid": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"createAt": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "date-time",
						},
					},
					"updateAt": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "date-time",
						},
					},
					"Spec": {
						SchemaProps: spec.SchemaProps{
							Description: "Spec defines the desired identities of identity in this set.",
							Ref:         ref("tkestack.io/tke/pkg/auth/types.LocalIdentitySpec"),
						},
					},
					"Status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("tkestack.io/tke/pkg/auth/types.LocalIdentityStatus"),
						},
					},
				},
				Required: []string{"Spec", "Status"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.LocalIdentitySpec", "tkestack.io/tke/pkg/auth/types.LocalIdentityStatus"},
	}
}

func schemaLocalIdentityList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LocalIdentityList is the whole list of all identities.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"items": {
						SchemaProps: spec.SchemaProps{
							Description: "List of identities.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.LocalIdentity"),
									},
								},
							},
						},
					},
				},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.LocalIdentity"},
	}
}

func schemaLocalIdentitySpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LocalIdentitySpec is a description of an identity.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"hashedPassword": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"originalPassword": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"groups": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"extra": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func schemaLocalIdentityStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LocalIdentityStatus is a description of an identity status.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"locked": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"boolean"},
							Format: "",
						},
					},
				},
			},
		},
	}
}

func schemaNonResourceAttributes(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NonResourceAttributes includes the authorization attributes available for non-resource requests to the Authorizer interface.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"path": {
						SchemaProps: spec.SchemaProps{
							Description: "Path is the URL path of the request.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"verb": {
						SchemaProps: spec.SchemaProps{
							Description: "Verb is the standard HTTP verb.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schemaPermission(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Permission defines a series of action on resource can be done or not.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"allow": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Type:   []string{"string"},
													Format: "",
												},
											},
										},
									},
								},
							},
						},
					},
					"deny": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Type:   []string{"string"},
													Format: "",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Required: []string{"allow", "deny"},
			},
		},
	}
}

func schemaPolicy(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Policy defines a data structure containing a authorization strategy.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"id": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"service": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"statement": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("tkestack.io/tke/pkg/auth/types.Statement"),
						},
					},
					"userName": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"createAt": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "date-time",
						},
					},
					"updateAt": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "date-time",
						},
					},
					"type": {
						SchemaProps: spec.SchemaProps{
							Description: "Type defines policy is created by system(1), or created by user(0)",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
				},
				Required: []string{"name", "id", "tenantID", "service", "statement", "userName", "description", "createAt", "updateAt", "type"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.Statement"},
	}
}

func schemaPolicyCreate(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PolicyCreate defines the policy create request.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"service": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"statement": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("tkestack.io/tke/pkg/auth/types.Statement"),
						},
					},
					"userName": {
						SchemaProps: spec.SchemaProps{
							Description: "UserName claims users attached to the policy created and split by ','. e.g: user1,user2.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"name", "tenantID", "service", "statement", "userName", "description"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.Statement"},
	}
}

func schemaPolicyList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PolicyList is the whole list of policy.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"items": {
						SchemaProps: spec.SchemaProps{
							Description: "List of policies.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.Policy"),
									},
								},
							},
						},
					},
				},
				Required: []string{"items"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.Policy"},
	}
}

func schemaPolicyMeta(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PolicyMeta contains metadata of Policy used for in roles.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"id": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"service": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"type": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int32",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"name", "id", "tenantID", "service", "type", "description"},
			},
		},
	}
}

func schemaPolicyOption(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PolicyOption is option for listing polices.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"id": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"userName": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"keyword": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"scope": {
						SchemaProps: spec.SchemaProps{
							Description: "PolicyStorage list scope: local, system, all",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"id", "userName", "name", "keyword", "tenantID", "scope"},
			},
		},
	}
}

func schemaResourceAttributes(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ResourceAttributes includes the authorization attributes available for resource requests to the Authorizer interface. Only verb and resource fields could be considered, tke-auth will ignore other fields right now.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"namespace": {
						SchemaProps: spec.SchemaProps{
							Description: "Namespace is the namespace of the action being requested.  Currently, there is no distinction between no namespace and all namespaces \"\" (empty) is defaulted for LocalSubjectAccessReviews \"\" (empty) is empty for cluster-scoped resources \"\" (empty) means \"all\" for namespace scoped resources from a SubjectAccessReview or SelfSubjectAccessReview",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"verb": {
						SchemaProps: spec.SchemaProps{
							Description: "Verb is a kubernetes resource API verb, like: get, list, watch, create, update, delete, proxy.  \"*\" means all.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"group": {
						SchemaProps: spec.SchemaProps{
							Description: "Group is the API Group of the Resource.  \"*\" means all.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Description: "Version is the API Version of the Resource.  \"*\" means all.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"resource": {
						SchemaProps: spec.SchemaProps{
							Description: "Resource is one of the existing resource types.  \"*\" means all.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"subresource": {
						SchemaProps: spec.SchemaProps{
							Description: "Subresource is one of the existing resource types.  \"\" means none.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"name": {
						SchemaProps: spec.SchemaProps{
							Description: "Name is the name of the resource being requested for a \"get\" or deleted for a \"delete\". \"\" (empty) means all.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}

func schemaRole(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Role is a collection with multiple policies.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"id": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"userName": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"policies": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.PolicyMeta"),
									},
								},
							},
						},
					},
					"createAt": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "date-time",
						},
					},
					"updateAt": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "date-time",
						},
					},
					"type": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int32",
						},
					},
				},
				Required: []string{"id", "name", "description", "tenantID", "userName", "policies", "createAt", "updateAt", "type"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.PolicyMeta"},
	}
}

func schemaRoleList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "RoleList is the whole list of policy.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"items": {
						SchemaProps: spec.SchemaProps{
							Description: "List of policies.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.Role"),
									},
								},
							},
						},
					},
				},
				Required: []string{"items"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.Role"},
	}
}

func schemaRoleOption(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "RoleOption is option for listing polices.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"id": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"userName": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"keyword": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"scope": {
						SchemaProps: spec.SchemaProps{
							Description: "RoleStorage list scope: local, system, all",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"id", "userName", "tenantID", "name", "keyword", "scope"},
			},
		},
	}
}

func schemaStatement(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Statement defines a series of action on resource can be done or not.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"action": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"resource": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"effect": {
						SchemaProps: spec.SchemaProps{
							Description: "Effect indicates action on the resource is allowed or not, can be \"allow\" or \"deny\"",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"action", "resource", "effect"},
			},
		},
	}
}

func schemaSubjectAccessReview(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SubjectAccessReview checks whether or not a user or group can perform an action.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Description: "Spec holds information about the request being evaluated.",
							Ref:         ref("tkestack.io/tke/pkg/auth/types.SubjectAccessReviewSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is filled in by the server and indicates whether the request is allowed or not.",
							Ref:         ref("tkestack.io/tke/pkg/auth/types.SubjectAccessReviewStatus"),
						},
					},
				},
				Required: []string{"spec"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.SubjectAccessReviewSpec", "tkestack.io/tke/pkg/auth/types.SubjectAccessReviewStatus"},
	}
}

func schemaSubjectAccessReviewSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SubjectAccessReviewSpec is a description of the access request.  Exactly one of ResourceAuthorizationAttributes and NonResourceAuthorizationAttributes must be set.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"resourceAttributes": {
						SchemaProps: spec.SchemaProps{
							Description: "ResourceAuthorizationAttributes describes information for a resource access request.",
							Ref:         ref("tkestack.io/tke/pkg/auth/types.ResourceAttributes"),
						},
					},
					"resourceAttributesList": {
						SchemaProps: spec.SchemaProps{
							Description: "ResourceAttributesList describes information for multi resource access request.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.ResourceAttributes"),
									},
								},
							},
						},
					},
					"nonResourceAttributes": {
						SchemaProps: spec.SchemaProps{
							Description: "NonResourceAttributes describes information for a non-resource access request.",
							Ref:         ref("tkestack.io/tke/pkg/auth/types.NonResourceAttributes"),
						},
					},
					"user": {
						SchemaProps: spec.SchemaProps{
							Description: "User is the user you're testing for.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"extra": {
						SchemaProps: spec.SchemaProps{
							Description: "Extra corresponds to the user.Info.GetExtra() method from the authenticator.  Since that is input to the authorizer it needs a reflection here.",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Type:   []string{"string"},
													Format: "",
												},
											},
										},
									},
								},
							},
						},
					},
					"uid": {
						SchemaProps: spec.SchemaProps{
							Description: "UID information about the requesting user.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.NonResourceAttributes", "tkestack.io/tke/pkg/auth/types.ResourceAttributes"},
	}
}

func schemaSubjectAccessReviewStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SubjectAccessReviewStatus indicates whether the request is allowed or not",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"allowed": {
						SchemaProps: spec.SchemaProps{
							Description: "Allowed is required. True if the action would be allowed, false otherwise.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"denied": {
						SchemaProps: spec.SchemaProps{
							Description: "Denied is optional. True if the action would be denied, otherwise false. If both allowed is false and denied is false, then the authorizer has no opinion on whether to authorize the action. Denied may not be true if Allowed is true.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"reason": {
						SchemaProps: spec.SchemaProps{
							Description: "Reason is optional.  It indicates why a request was allowed or denied.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"evaluationError": {
						SchemaProps: spec.SchemaProps{
							Description: "EvaluationError is an indication that some error occurred during the authorization check. It is entirely possible to get an error and be able to continue determine authorization status in spite of it. For instance, RBAC can be missing a role, but enough roles are still present and bound to reason about the request.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"allowedList": {
						SchemaProps: spec.SchemaProps{
							Description: "AllowedList is the allowed response for batch authorization request.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("tkestack.io/tke/pkg/auth/types.AllowedResponse"),
									},
								},
							},
						},
					},
				},
				Required: []string{"allowed"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.AllowedResponse"},
	}
}

func schemaTokenReviewRequest(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TokenReviewRequest attempts to authenticate a token to a known user.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion is default as \"authentication.k8s.io/v1beta1\".",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is default as \"TokenReview\".",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Description: "Spec holds information about the request being evaluated.",
							Ref:         ref("tkestack.io/tke/pkg/auth/types.TokenReviewSpec"),
						},
					},
				},
				Required: []string{"apiVersion", "kind", "spec"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.TokenReviewSpec"},
	}
}

func schemaTokenReviewResponse(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TokenReviewResponse is response info for authenticating a token to a known user.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion is default as \"authentication.k8s.io/v1beta1\".",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is default as \"TokenReview\".",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is filled in by the server and indicates whether the request can be authenticated.",
							Ref:         ref("tkestack.io/tke/pkg/auth/types.TokenReviewStatus"),
						},
					},
				},
				Required: []string{"apiVersion", "kind", "status"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.TokenReviewStatus"},
	}
}

func schemaTokenReviewSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TokenReviewSpec is a description of the token authentication request.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"token": {
						SchemaProps: spec.SchemaProps{
							Description: "Token is the opaque bearer token.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"token"},
			},
		},
	}
}

func schemaTokenReviewStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TokenReviewStatus is the result of the token authentication request.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"authenticated": {
						SchemaProps: spec.SchemaProps{
							Description: "Authenticated indicates that the token was associated with a known user.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"user": {
						SchemaProps: spec.SchemaProps{
							Description: "User is the UserInfo associated with the provided token.",
							Ref:         ref("tkestack.io/tke/pkg/auth/types.TokenReviewUser"),
						},
					},
				},
				Required: []string{"authenticated"},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/auth/types.TokenReviewUser"},
	}
}

func schemaTokenReviewUser(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TokenReviewUser holds the information about the user needed to implement the",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"username": {
						SchemaProps: spec.SchemaProps{
							Description: "The name that uniquely identifies this user among all active users.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"uid": {
						SchemaProps: spec.SchemaProps{
							Description: "A unique value that identifies this user across time. If this user is deleted and another user by the same name is added, they will have different UIDs.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"groups": {
						SchemaProps: spec.SchemaProps{
							Description: "The names of groups this user is a part of.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"extra": {
						SchemaProps: spec.SchemaProps{
							Description: "Any additional information provided by the authenticator.",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type: []string{"array"},
										Items: &spec.SchemaOrArray{
											Schema: &spec.Schema{
												SchemaProps: spec.SchemaProps{
													Type:   []string{"string"},
													Format: "",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Required: []string{"username", "uid", "groups", "extra"},
			},
		},
	}
}
