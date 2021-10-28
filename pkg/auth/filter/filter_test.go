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
	"reflect"
	"testing"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericrequest "k8s.io/apiserver/pkg/endpoints/request"

	"tkestack.io/tke/api/business"
	"tkestack.io/tke/api/registry"
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
				Verb:            "list",
				APIGroup:        business.GroupName,
				ResourceRequest: true,
				Resource:        "projects",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listProjects",
				Resource: "project:*",
			},
		},
		{
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				APIGroup:        business.GroupName,
				ResourceRequest: true,
				Resource:        "projects",
				Name:            "demo",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "getProject",
				Resource: "project:demo",
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
				Resource: "project:demo/namespace:demo",
			},
		},
		{
			//GET -H "X-TKE-ProjectID: xxx" /business.tkestack.io/v1/namespaces/demo/namespace/demo
			ctx: contextWithProject(context.Background()),
			attr: &authorizer.AttributesRecord{
				Verb:            "get",
				Namespace:       "demo",
				APIGroup:        business.GroupName,
				ResourceRequest: true,
				Resource:        "namespaces",
				Name:            "demo",
				Subresource:     "",
				User:            &user.DefaultInfo{},
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "getNamespace",
				Resource: "project:demo/namespace:demo",
				User:     &user.DefaultInfo{},
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
				Resource: "project:demo/namespace:*",
			},
		},

		{
			//GET -H "X-TKE-ProjectID: xxx" /business.tkestack.io/v1/namespaces/demo/namespaces
			ctx: contextWithProject(context.Background()),
			attr: &authorizer.AttributesRecord{
				Verb:            "list",
				Namespace:       "demo",
				APIGroup:        business.GroupName,
				ResourceRequest: true,
				Resource:        "namespaces",
				User:            &user.DefaultInfo{},
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listNamespaces",
				Resource: "project:demo/namespace:*",
				User:     &user.DefaultInfo{},
			},
		},
		{
			//GET /registry.tkestack.io/v1/namespaces/
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "list",
				APIGroup:        registry.GroupName,
				ResourceRequest: true,
				Resource:        "namespaces",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listRegistrynamespaces",
				Resource: "registrynamespace:*",
			},
		},
		{
			//GET /registry.tkestack.io/v1/namespaces/rns-57g8dc9v/repositories
			ctx: context.Background(),
			attr: &authorizer.AttributesRecord{
				Verb:            "list",
				Namespace:       "demo",
				APIGroup:        registry.GroupName,
				ResourceRequest: true,
				Resource:        "repositories",
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listRepositories",
				Resource: "registrynamespace:demo/repository:*",
			},
		},
		{
			//GET -H "X-TKE-ClusterName: xxx" /apps/v1/namespaces/demo/deployments
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
			//GET -H "X-TKE-ClusterName: xxx" -H "X-TKE-ProjectID: xxx" /apps/v1/namespaces/demo/deployments
			ctx: contextWithProject(contextWithCluster(context.Background())),
			attr: &authorizer.AttributesRecord{
				Verb:            "list",
				Namespace:       "demo",
				ResourceRequest: true,
				Resource:        "deployments",
				Subresource:     "",
				User: &user.DefaultInfo{
					Groups: []string{fmt.Sprintf("project:%s", projectID)},
				},
			},
			expect: &authorizer.AttributesRecord{
				Verb:     "listDeployments",
				Resource: fmt.Sprintf("cluster:%s/namespace:demo/deployment:*", clusterName),
				User: &user.DefaultInfo{
					Groups: []string{fmt.Sprintf("project:%s", projectID)},
				},
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
			ctx: context.Background(),
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

	for i, testCase := range testCases {
		result := ConvertTKEAttributes(testCase.ctx, testCase.attr)
		if !compare(result, testCase.expect) {
			t.Fatalf("%d, expect attributes %v, but got %v", i, testCase.expect, result)
		}
	}
}

func compare(a authorizer.Attributes, b authorizer.Attributes) bool {
	if a.GetVerb() == b.GetVerb() && a.GetResource() == b.GetResource() && reflect.DeepEqual(a.GetUser(), b.GetUser()) {
		return true
	}
	fmt.Println(a.GetUser(), b.GetUser())

	return false
}

const (
	clusterContextKey = "clusterName"
	clusterName       = "cls-82qkvzgp"
	projectID         = "prj-82qkvzgp"
	projectContextKey = "projectID"
)

func contextWithCluster(ctx context.Context) context.Context {
	return genericrequest.WithValue(ctx, clusterContextKey, clusterName)
}

func contextWithProject(ctx context.Context) context.Context {
	return genericrequest.WithValue(ctx, projectContextKey, projectID)
}
