/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package util

import (
	"strconv"
	"strings"

	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/fields"

	"tkestack.io/tke/api/auth"
)

const (
	defaultQueryLimit = 50
)

func ParseQueryKeywordAndLimit(options *metainternal.ListOptions) (string, int) {
	keyword := ""
	limit := defaultQueryLimit
	if options.FieldSelector != nil {
		keyword, _ = options.FieldSelector.RequiresExactMatch(auth.KeywordQueryTag)
		limitStr, _ := options.FieldSelector.RequiresExactMatch(auth.QueryLimitTag)
		if li, err := strconv.Atoi(limitStr); err == nil && li >= 0 {
			limit = li
		}
	}

	return keyword, limit
}

func InterceptKeyword(options *metainternal.ListOptions) string {
	keyword := ""
	found := false
	if options.FieldSelector != nil {
		keyword, found = options.FieldSelector.RequiresExactMatch(auth.KeywordQueryTag)
		if found {
			removeKeywordFromField(options)
		}
	}

	return keyword
}

func removeKeywordFromField(options *metainternal.ListOptions) {
	strs := strings.Split(options.FieldSelector.String(), ",")
	var remain []string
	for _, str := range strs {
		s, _ := fields.ParseSelector(str)
		_, found := s.RequiresExactMatch(auth.KeywordQueryTag)
		if !found {
			remain = append(remain, str)
		}
	}

	selector, _ := fields.ParseSelector(strings.Join(remain, ","))
	options.FieldSelector = selector
}
