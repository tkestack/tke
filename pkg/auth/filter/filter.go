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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"

	"github.com/go-openapi/inflect"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/audit"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericfilters "k8s.io/apiserver/pkg/endpoints/filters"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericfilter "tkestack.io/tke/pkg/apiserver/filter"

	"tkestack.io/tke/api/business"
	"tkestack.io/tke/api/registry"
	commonapiserverfilter "tkestack.io/tke/pkg/apiserver/filter"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/platform/registry/cluster"
	"tkestack.io/tke/pkg/util/log"
)

const (
	maxCheckClusterNameCount int = 6

	createProjectAction string = "createProject"
	updateProjectAction string = "updateProject"
	createClusterAction string = "createCluster"

	// Annotation key names set in advanced audit
	decisionAnnotationKey = "authorization.auth.tke.com/decision"
	reasonAnnotationKey   = "authorization.auth.tke.com/reason"

	// Annotation values set in advanced audit
	decisionAllow          = "allow"
	decisionForbid         = "forbid"
	reasonError            = "internal error"
	k8sDecisionAllowHeader = "X-TKE-K8sDecisionAllow"
)

var (
	AllowClusterNotFoundActions = sets.NewString(createClusterAction, updateProjectAction)
)

func ExtractClusterNames(ctx context.Context, req *http.Request, resource string) []string {
	clusterNames := sets.NewString()

	clusterName := filter.ClusterFrom(ctx)
	if len(clusterName) > 0 {
		clusterNames.Insert(clusterName)
	}

	filterResourceClusterNames := cluster.ClusterPattern.FindAllString(resource, -1)
	for _, filterClusterName := range filterResourceClusterNames {
		clusterNames.Insert(cluster.NamePattern.FindAllString(filterClusterName, -1)...)
	}

	filterURLClusterNames := cluster.ClusterPattern.FindAllString(req.URL.String(), -1)
	for _, filterClusterName := range filterURLClusterNames {
		clusterNames.Insert(cluster.NamePattern.FindAllString(filterClusterName, -1)...)
	}

	return clusterNames.List()
}

func ForbiddenResponse(ctx context.Context, tkeAttributes authorizer.Attributes,
	w http.ResponseWriter, req *http.Request, s runtime.NegotiatedSerializer, reason string) {
	log.Infof("Forbidden: %#v %#v, Reason: %q", req.Method, req.RequestURI, reason)
	audit.AddAuditAnnotation(ctx, decisionAnnotationKey, decisionForbid)
	audit.AddAuditAnnotation(ctx, reasonAnnotationKey, reason)
	responsewriters.Forbidden(ctx, tkeAttributes, w, req, reason, s)
}

// WithTKEAuthorization passes all tke-auth authorized requests on to handler, and returns a forbidden error otherwise.
func WithTKEAuthorization(handler http.Handler, a authorizer.Authorizer, s runtime.NegotiatedSerializer, ignoreAuthPathPrefixes []string) http.Handler {
	if a == nil {
		log.Warn("TKE Authorization is disabled")
		return handler
	}
	allIgnorePathPrefixes := commonapiserverfilter.MakeAllIgnoreAuthPathPrefixes(ignoreAuthPathPrefixes)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodOptions {
			handler.ServeHTTP(w, req)
			return
		}

		ignorePathPrefix := ""
		reqPath := strings.ToLower(req.URL.Path)
		for _, pathPrefix := range allIgnorePathPrefixes {
			if strings.HasPrefix(reqPath, strings.ToLower(pathPrefix)) {
				ignorePathPrefix = pathPrefix
				break
			}
		}
		if ignorePathPrefix != "" {
			handler.ServeHTTP(w, req)
			return
		}

		ctx := req.Context()
		attributes, err := genericfilters.GetAuthorizerAttributes(ctx)
		if err != nil {
			responsewriters.InternalError(w, req, err)
			return
		}

		var (
			authorized authorizer.Decision
			reason     string
		)

		tenantID := ""
		extra := attributes.GetUser().GetExtra()
		if len(extra) > 0 {
			if tenantIDs, ok := extra[genericoidc.TenantIDKey]; ok {
				if len(tenantIDs) > 0 {
					tenantID = tenantIDs[0]
				} else {
					tenantID = "default"
				}
			}
		}
		find := false
		if tenantID == "" {
			find, tenantID = genericfilter.FindValueFromGroups(attributes.GetUser().GetGroups(), "tenant")
			if find && tenantID == "" {
				tenantID = "default"
			}
		}

		// firstly check if resource is unprotected
		authorized = UnprotectedAuthorized(attributes)
		if authorized != authorizer.DecisionAllow {
			authorized, reason, err = a.Authorize(ctx, attributes)
		}

		// secondly check k8s resource authz result
		if authorized != authorizer.DecisionAllow {
			if tenantID != "" {
				log.Debugf("TKEStack user '%v'", attributes.GetUser())
				attributes = ConvertTKEAttributes(ctx, attributes)
			}
			authorized, reason, err = a.Authorize(ctx, attributes)
		} else {
			setK8sDecision(req, true)
		}

		// finaly check tke casbin resource authz resoult
		// an authorizer like RBAC could encounter evaluation errors and still allow the request, so authorizer decision is checked before error here.
		if authorized == authorizer.DecisionAllow {
			audit.AddAuditAnnotation(ctx, decisionAnnotationKey, decisionAllow)
			audit.AddAuditAnnotation(ctx, reasonAnnotationKey, reason)
			handler.ServeHTTP(w, req)
			return
		}
		if err != nil {
			audit.AddAuditAnnotation(ctx, reasonAnnotationKey, reasonError)
			responsewriters.InternalError(w, req, err)
			return
		}

		ForbiddenResponse(ctx, attributes, w, req, s, reason)
	})
}

func WithInspectors(handler http.Handler, inspectors []Inspector, c *genericapiserver.Config) http.Handler {
	if len(inspectors) > 0 {
		for _, inspector := range inspectors {
			handler = inspector.Inspect(handler, c)
		}
	}
	return handler
}

var (
	unprotectedVerbSets = sets.NewString("listPortal")
)

// UnprotectedAuthorized checks a request attribute has privileged to pass authorization.
func UnprotectedAuthorized(attributes authorizer.Attributes) authorizer.Decision {
	info := attributes.GetUser()
	if info == nil {
		return authorizer.DecisionNoOpinion
	}

	verb := attributes.GetVerb()
	if unprotectedVerbSets.Has(verb) {
		return authorizer.DecisionAllow
	}

	return authorizer.DecisionNoOpinion
}

// specialSubResources contains resources which get verb use get instead of list
var specialSubResources = sets.NewString("status", "log", "finalize", "proxy")

// ConvertTKEAttributes converts attributes parsed by apiserver compatible with casbin enforcer
func ConvertTKEAttributes(ctx context.Context, attr authorizer.Attributes) authorizer.Attributes {
	tkeAttribs := attr.(*authorizer.AttributesRecord)

	resourceType := attr.GetResource()
	subResource := attr.GetSubresource()
	resourceName := attr.GetName()

	if resourceType == "namespaces" && attr.GetAPIGroup() == registry.GroupName {
		resourceType = "registrynamespaces"
	}

	clusterName := ""

	// URL forms: /clusters/{cluster}/{kind}/*, where parts are adjusted to be relative to kind
	if resourceType == "clusters" && len(subResource) != 0 && !specialSubResources.Has(subResource) {
		resourceType = subResource
		clusterName = resourceName
		resourceName = getNextPart(subResource, attr.GetPath())
		subResource = getNextPart(resourceName, attr.GetPath())
	}

	resourceTypeSingle := inflect.Singularize(resourceType)
	if resourceType == "status" {
		resourceTypeSingle = resourceType
	}

	verb := attr.GetVerb()
	switch verb {
	case "list":
		if len(resourceName) != 0 {
			resourceType = resourceTypeSingle
		}
	case "get":
		if len(subResource) != 0 && !specialSubResources.Has(subResource) {
			verb = "list"
		}

		if len(resourceName) == 0 {
			verb = "list"
		} else {
			resourceType = resourceTypeSingle
		}
	case "update":
		resourceType = resourceTypeSingle
	case "patch":
		verb = "update"
		resourceType = resourceTypeSingle
	case "deletecollection":
		verb = "delete"
		resourceType = resourceTypeSingle
	default:
		resourceType = resourceTypeSingle
	}

	if tkeAttribs.ResourceRequest {
		// if not specify resource name in path, set it to "*" (all)
		if len(resourceName) == 0 {
			resourceName = "*"
		}

		// URL forms: GET /users/jack/policies,  parsed into verb: getUserPolicies, resource: users:jack/policies:*
		tkeAttribs.Verb = fmt.Sprintf("%s%s%s", verb, upperFirst(resourceType), upperFirst(subResource))
		tkeAttribs.Resource = fmt.Sprintf("%s:%s", resourceTypeSingle, resourceName)
	} else {
		tkeAttribs.Verb = verb
		tkeAttribs.Resource = resourceType
	}

	if tkeAttribs.Namespace != "" {
		switch attr.GetAPIGroup() {
		case business.GroupName:
			tkeAttribs.Resource = fmt.Sprintf("project:%s/%s", tkeAttribs.Namespace, tkeAttribs.Resource)
		case registry.GroupName:
			tkeAttribs.Resource = fmt.Sprintf("registrynamespace:%s/%s", tkeAttribs.Namespace, tkeAttribs.Resource)
		default:
			if resourceType != "namespace" {
				tkeAttribs.Resource = fmt.Sprintf("namespace:%s/%s", tkeAttribs.Namespace, tkeAttribs.Resource)
			}
		}
	} else {
		// for /apis/platform.tkestack.io/v1/clusters/cls-xxx/lbcfbackendgroups?namespace=ns
		ns := filter.NamespaceFrom(ctx)
		if ns != "" {
			tkeAttribs.Resource = fmt.Sprintf("namespace:%s/%s", ns, tkeAttribs.Resource)
		}
	}

	if ctx != nil && len(filter.ClusterFrom(ctx)) != 0 {
		clusterName = filter.ClusterFrom(ctx)
	}

	if clusterName == "" && attr.GetUser() != nil {
		clusterName = commonapiserverfilter.GetClusterFromGroups(attr.GetUser().GetGroups())
	}

	if clusterName != "" && resourceTypeSingle != "cluster" {
		tkeAttribs.Resource = fmt.Sprintf("cluster:%s/%s", clusterName, tkeAttribs.Resource)
	}

	tkeAttribs.Subresource = subResource
	tkeAttribs.Name = resourceName

	return tkeAttribs
}

// upperFirst makes the first char of a string uppercase
func upperFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// getNextPart returns the next segments of a URL path for a name.
func getNextPart(flag string, path string) string {
	parts := splitPath(path)
	for i, part := range parts {
		if part == flag {
			if i+1 <= len(parts)-1 {
				return parts[i+1]
			}
		}
	}
	return ""
}

// splitPath returns the segments for a URL path.
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

func setK8sDecision(req *http.Request, k8sDecisionAllow bool) {
	if req != nil {
		req.Header.Set(k8sDecisionAllowHeader, strconv.FormatBool(k8sDecisionAllow))
	}
}

func getK8sDecision(req *http.Request) (k8sDecisionAllow bool) {
	if req != nil {
		k8sDecisionAllowString := req.Header.Get(k8sDecisionAllowHeader)
		k8sDecisionAllow, _ = strconv.ParseBool(k8sDecisionAllowString)
	}
	return k8sDecisionAllow
}
