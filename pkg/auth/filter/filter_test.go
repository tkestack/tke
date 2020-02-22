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

package filter

import (
	"context"
	"fmt"
	"testing"

	"tkestack.io/tke/api/business"

	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericrequest "k8s.io/apiserver/pkg/endpoints/request"
)

type testCase struct {
	ctx    context.Context
	attr   authorizer.Attributes
	expect authorizer.Attributes
}

func TestConvertTKEAttributes(t *testing.T) {
	testCases := []testCase{
		{
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "list",
				Resource:        "policies",
				Subresource:     "",
				ResourceRequest: true,
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listPolicies",
				Resource: "policy:*",
			},
		},
		{
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				Resource:        "policies",
				Name:            "policy-default-123",
				ResourceRequest: true,
				Subresource:     "",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "getPolicy",
				Resource: "policy:policy-default-123",
			},
		},
		{
			ctx: contextWithCluster(context.Background()),
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				Namespace:       "demo",
				ResourceRequest: true,
				Resource:        "namespaces",
				Name:            "demo",
				Subresource:     "",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "getNamespace",
				Resource: fmt.Sprintf("cluster:%s/namespace:demo", clusterName),
			},
		},
		{
			ctx: contextWithCluster(context.Background()),
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				ResourceRequest: true,
				Resource:        "namespaces",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listNamespaces",
				Resource: fmt.Sprintf("cluster:%s/namespace:*", clusterName),
			},
		},
		{
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				Namespace:       "demo",
				APIGroup:        business.GroupName,
				ResourceRequest: true,
				Resource:        "namespaces",
				Name:            "demo",
				Subresource:     "",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "getNamespace",
				Resource: "namespace:demo/namespace:demo",
			},
		},
		{
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "list",
				Namespace:       "demo",
				APIGroup:        business.GroupName,
				ResourceRequest: true,
				Resource:        "namespaces",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listNamespaces",
				Resource: "namespace:demo/namespace:*",
			},
		},
		{
			ctx: contextWithCluster(context.Background()),
			attr: &authorizer.AttributesRecord{
				Verb:            "list",
				Namespace:       "demo",
				ResourceRequest: true,
				Resource:        "deployments",
				Subresource:     "",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listDeployments",
				Resource: fmt.Sprintf("cluster:%s/namespace:demo/deployment:*", clusterName),
			},
		},
		{
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				Resource:        "clusters",
				ResourceRequest: true,
				Name:            "cls-82qkvzgp",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "getCluster",
				Resource: "cluster:cls-82qkvzgp",
			},
		},
		{
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				Resource:        "clusters",
				ResourceRequest: true,
				Name:            "cls-82qkvzgp",
				Subresource:     "alarmpolicies",
				Path:            "/api/v1/clusters/cls-82qkvzgp/alarmpolicies/test",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "getAlarmpolicy",
				Resource: "cluster:cls-82qkvzgp/alarmpolicy:test",
			},
		},
		{
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				Resource:        "clusters",
				ResourceRequest: true,
				Name:            "cls-82qkvzgp",
				Subresource:     "alarmpolicies",
				Path:            "/api/v1/clusters/cls-82qkvzgp/alarmpolicies",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listAlarmpolicies",
				Resource: "cluster:cls-82qkvzgp/alarmpolicy:*",
			},
		},
		{
			ctx: contextWithCluster(context.Background()),
			attr: &authorizer.AttributesRecord{
				Verb:            "list",
				ResourceRequest: false,
				Resource:        "/healthz",
				Path:            "/healthz",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "list",
				Resource: "cluster:cls-82qkvzgp//healthz",
			},
		},
	}

	for _, testCase := range testCases {
		result := ConvertTKEAttributes(testCase.ctx, testCase.attr)
		if !compareAttributesVerbAndRes(result, testCase.expect) {
			t.Errorf("expect attributes %v, but got %v", testCase.expect, result)
		}
	}
}

func compareAttributesVerbAndRes(a authorizer.Attributes, b authorizer.Attributes) bool {
	if a.GetVerb() == b.GetVerb() && a.GetResource() == b.GetResource() {
		return true
	}

	return false
}

const clusterContextKey = "clusterName"
const clusterName = "cls-82qkvzgp"

func contextWithCluster(ctx context.Context) context.Context {
	return genericrequest.WithValue(ctx, clusterContextKey, clusterName)
}
