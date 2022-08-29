/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2022 Tencent. All Rights Reserved.
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

package validation

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

type TKEValidateResult struct {
	Name        string          `json:"Name"`
	Description string          `json:"Description"`
	Checked     bool            `json:"Checked"`
	Passed      bool            `json:"Passed"`
	ErrorList   field.ErrorList `json:"-"`
	Detail      string          `json:"Detail"`
}

func (r TKEValidateResult) ToFieldError() *field.Error {
	if len(r.ErrorList) == 0 {
		if r.Checked {
			r.Passed = true
		}
		r.Detail = ""
	} else {
		r.Detail = r.ErrorList.ToAggregate().Error()
	}
	message, _ := json.Marshal(r)
	return &field.Error{
		Type:   field.ErrorTypeInvalid,
		Field:  r.Name,
		Detail: string(message),
	}
}
