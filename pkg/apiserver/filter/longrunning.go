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
	"k8s.io/apimachinery/pkg/util/sets"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"net/http"
	"strings"
)

// LongRunningRequestCheck returns true if the given request has one of the
// specified verbs or one of the specified subresources, or is matched path
// prefix request.
func LongRunningRequestCheck(longRunningVerbs, longRunningSubresources sets.String, longRunningPathPrefixes []string) apirequest.LongRunningRequestCheck {
	return func(r *http.Request, requestInfo *apirequest.RequestInfo) bool {
		if longRunningVerbs.Has(requestInfo.Verb) {
			return true
		}
		if requestInfo.IsResourceRequest && longRunningSubresources.Has(requestInfo.Subresource) {
			return true
		}
		if !requestInfo.IsResourceRequest && strings.HasPrefix(requestInfo.Path, "/debug/pprof/") {
			return true
		}
		if !requestInfo.IsResourceRequest && len(longRunningPathPrefixes) > 0 {
			for _, prefix := range longRunningPathPrefixes {
				if strings.HasPrefix(requestInfo.Path, prefix) {
					return true
				}
			}
		}
		return false
	}
}
