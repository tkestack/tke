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
	"github.com/go-openapi/inflect"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/audit"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	genericfilters "k8s.io/apiserver/pkg/endpoints/filters"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog"
	"net/http"
	"strings"
	genericoidc "tkestack.io/tke/pkg/apiserver/authentication/authenticator/oidc"
	commonapiserverfilter "tkestack.io/tke/pkg/apiserver/filter"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	"tkestack.io/tke/pkg/util/log"
	"unicode"
)

const (
	// Annotation key names set in advanced audit
	decisionAnnotationKey = "authorization.auth.tke.com/decision"
	reasonAnnotationKey   = "authorization.auth.tke.com/reason"

	// Annotation values set in advanced audit
	decisionAllow  = "allow"
	decisionForbid = "forbid"
	reasonError    = "internal error"
)

// WithTKEAuthorization passes all tke-auth authorized requests on to handler, and returns a forbidden error otherwise.
func WithTKEAuthorization(handler http.Handler, a authorizer.Authorizer, s runtime.NegotiatedSerializer, ignoreAuthPathPrefixes []string) http.Handler {
	if a == nil {
		klog.Warningf("TKE Authorization is disabled")
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
		ae := request.AuditEventFrom(ctx)
		attributes, err := genericfilters.GetAuthorizerAttributes(ctx)
		if err != nil {
			responsewriters.InternalError(w, req, err)
			return
		}

		var (
			authorized authorizer.Decision
			reason     string
		)
		// first check if user is admin
		tkeAttributes := ConvertTKEAttributes(ctx, attributes)
		authorized = UnprotectedAuthorized(tkeAttributes)
		if authorized != authorizer.DecisionAllow {
			log.Debug("Convert to tke tkeAttributes", log.String("user name", tkeAttributes.GetUser().GetName()),
				log.String("resource", tkeAttributes.GetResource()), log.String("resource", tkeAttributes.GetName()),
				log.String("verb", tkeAttributes.GetVerb()))
			authorized, reason, err = a.Authorize(tkeAttributes)
		}

		// an authorizer like RBAC could encounter evaluation errors and still allow the request, so authorizer decision is checked before error here.
		if authorized == authorizer.DecisionAllow {
			audit.LogAnnotation(ae, decisionAnnotationKey, decisionAllow)
			audit.LogAnnotation(ae, reasonAnnotationKey, reason)
			handler.ServeHTTP(w, req)
			return
		}
		if err != nil {
			audit.LogAnnotation(ae, reasonAnnotationKey, reasonError)
			responsewriters.InternalError(w, req, err)
			return
		}

		log.Infof("Forbidden: %#v %#v, Reason: %q", req.Method, req.RequestURI, reason)
		audit.LogAnnotation(ae, decisionAnnotationKey, decisionForbid)
		audit.LogAnnotation(ae, reasonAnnotationKey, reason)
		responsewriters.Forbidden(ctx, tkeAttributes, w, req, reason, s)
	})
}

var (
	unprotectedVerbSets = sets.NewString("listPortal", "createApikey", "listApikey", "updateApikey")
)

// UnprotectedAuthorized checks a request attribute has privileged to pass authorization.
func UnprotectedAuthorized(attributes authorizer.Attributes) authorizer.Decision {
	info := attributes.GetUser()
	if info == nil {
		return authorizer.DecisionNoOpinion
	}
	extras := info.GetExtra()
	tenantID, ok := extras[genericoidc.TenantIDKey]
	if info.GetName() == "admin" && (!ok || len(tenantID) == 0) {
		return authorizer.DecisionAllow
	}

	verb := attributes.GetVerb()
	if unprotectedVerbSets.Has(verb) {
		return authorizer.DecisionAllow
	}

	return authorizer.DecisionNoOpinion
}

// specialUpdateResources contains resources which update verb may be considered as add, verb will be add
var specialVerbUpdateResources = sets.NewString("roles", "policies")

// specialSubResources contains resources which get verb use get instead of list
var specialSubResources = sets.NewString("status", "log", "finalize")

// ConvertTKEAttributes converts attributes parsed by apiserver compatible with casbin enforcer
func ConvertTKEAttributes(ctx context.Context, attr authorizer.Attributes) authorizer.Attributes {
	tkeAttribs := attr.(*authorizer.AttributesRecord)

	log.Debug("Attr parsed by k8s resolver", log.Any("attr", tkeAttribs))
	resourceType := attr.GetResource()
	subResource := attr.GetSubresource()
	resourceName := attr.GetName()

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
		if len(subResource) != 0 && specialVerbUpdateResources.Has(resourceType) {
			verb = "add"
		}
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

	// if not specify resource name in path, set it to "*" (all)
	if len(resourceName) == 0 {
		resourceName = "*"
	}

	// URL forms: GET /users/jack/policies,  parsed into verb: getUserPolicies, resource: users:jack/policies:*
	tkeAttribs.Verb = fmt.Sprintf("%s%s%s", verb, upperFirst(resourceType), upperFirst(subResource))
	tkeAttribs.Resource = fmt.Sprintf("%s:%s", resourceTypeSingle, resourceName)

	if tkeAttribs.Namespace != "" {
		tkeAttribs.Resource = fmt.Sprintf("namespace:%s/%s", tkeAttribs.Namespace, tkeAttribs.Resource)
	}

	if ctx != nil && len(filter.ClusterFrom(ctx)) != 0 {
		clusterName = filter.ClusterFrom(ctx)
	}

	if len(clusterName) != 0 {
		tkeAttribs.Resource = fmt.Sprintf("cluster:%s/%s", clusterName, tkeAttribs.Resource)
	}

	tkeAttribs.Subresource = subResource
	tkeAttribs.Name = resourceName

	log.Debug("Convert to tke attributes", log.Any("tke attributes", tkeAttribs))
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
