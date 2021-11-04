/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package openapi

import (
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// GetOpenAPIDefinitions provide definition for all models used by routes.
func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	definitions := make(map[string]common.OpenAPIDefinition)
	definitions["tkestack.io/tke/pkg/mesh/services/rest.Response"] = schemaResponse(ref)
	return definitions
}

func schemaResponse(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "response for creation, update, get, deletion and list",
				Properties: map[string]spec.Schema{
					"result": {
						SchemaProps: spec.SchemaProps{
							Type:        []string{"boolean"},
							Description: "Result of request",
						},
					},
					"err": {
						SchemaProps: spec.SchemaProps{
							Type:        []string{"string"},
							Description: "if result is false, this field will tell you the error message",
						},
					},
					"data": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"string"},
							Description: "if request is get, this field will return single object data," +
								"if request is list, this field will return list data;",
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}
