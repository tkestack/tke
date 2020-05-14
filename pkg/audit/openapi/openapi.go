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
)

// GetOpenAPIDefinitions provide definition for all models used by routes.
func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	definitions := make(map[string]common.OpenAPIDefinition)
	definitions["tkestack.io/tke/pkg/audit/storage/types.Event"] = schemaAlarmPolicyRequest(ref)
	return definitions
}

func schemaAlarmPolicyRequest(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Alarm policy request for creation, update",
				Properties: map[string]spec.Schema{
					"AuditID": {
						SchemaProps: spec.SchemaProps{
							Type:        []string{"string"},
							Description: "Namespace of alarm object, be empty if policy apply to all objects",
						},
					},
					"Stage": {
						SchemaProps: spec.SchemaProps{
							Type:        []string{"string"},
							Description: "Workload type of alarm object, be empty if policy apply to all objects",
						},
					},
				},
			},
		},
	}
}
