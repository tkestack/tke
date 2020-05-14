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
	sets "k8s.io/apimachinery/pkg/util/sets"
	"tkestack.io/tke/api/auth"
)

func InSubjects(subject auth.Subject, slice []auth.Subject) bool {
	for _, s := range slice {
		if subject.ID == s.ID {
			return true
		}
	}
	return false
}

func InSubjectsWithIDOrName(subject auth.Subject, slice []auth.Subject) bool {
	for _, s := range slice {
		if s.ID != "" {
			if subject.ID == s.ID {
				return true
			}
		}

		if s.Name != "" {
			if subject.Name == s.Name {
				return true
			}
		}
	}
	return false
}

func RemoveDuplicateSubjects(slice []auth.Subject) []auth.Subject {
	var ret []auth.Subject
	idSet := sets.String{}
	for _, s := range slice {
		if !idSet.Has(s.ID) {
			ret = append(ret, s)
			idSet.Insert(s.ID)
		}
	}
	return ret
}

func RemoveDuplicateSubjectsByIDOrName(slice []auth.Subject) []auth.Subject {
	var ret []auth.Subject
	idSet := sets.String{}
	nameSet := sets.String{}
	for _, s := range slice {
		if s.ID != "" {
			if !idSet.Has(s.ID) {
				ret = append(ret, s)
				idSet.Insert(s.ID)
			}
			continue
		}
		if s.Name != "" {
			if !nameSet.Has(s.Name) {
				ret = append(ret, s)
				nameSet.Insert(s.Name)
			}
		}
	}
	return ret
}
