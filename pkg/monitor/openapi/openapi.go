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
	definitions["tkestack.io/tke/pkg/monitor/services/rest.AlarmPolicy"] = schemaAlarmPolicyRequest(ref)
	definitions["tkestack.io/tke/pkg/monitor/services/rest.Response"] = schemaAlarmPolicyResponse(ref)
	return definitions
}

func schemaAlarmPolicyRequest(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Alarm policy request for creation, update",
				Properties: map[string]spec.Schema{
					"AlarmPolicySettings": {
						SchemaProps: spec.SchemaProps{
							Ref:         ref("tkestack.io/tke/pkg/monitor/services/rest.AlarmPolicySettings"),
							Description: "AlarmPolicySettings defines alarm policy details, including name, type and metrics",
						},
					},
					"NotifySettings": {
						SchemaProps: spec.SchemaProps{
							Ref:         ref("tkestack.io/tke/pkg/monitor/services/rest.NotifySettings"),
							Description: "NotifySettings contains notification info of alarm policy, including receiver groups and notify ways",
						},
					},
					"Namespace": {
						SchemaProps: spec.SchemaProps{
							Type:        []string{"string"},
							Description: "Namespace of alarm object, be empty if policy apply to all objects",
						},
					},
					"WorkloadType": {
						SchemaProps: spec.SchemaProps{
							Type:        []string{"string"},
							Description: "Workload type of alarm object, be empty if policy apply to all objects",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"tkestack.io/tke/pkg/monitor/services/rest.AlarmPolicySettings",
			"tkestack.io/tke/pkg/monitor/services/rest.NotifySettings",
		},
	}
}

func schemaAlarmPolicyResponse(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Prometheus rule response for creation, update, get, deletion and list",
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
							Type: []string{"interface{}"},
							Description: "if request is get, this field will return tkestack.io/tke/pkg/monitor/services/rest.AlarmPolicy;" +
								"if request is list, this field will return tkestack.io/tke/pkg/monitor/services/rest.AlarmPolicies;",
						},
					},
				},
			},
		},
		Dependencies: []string{},
	}
}
