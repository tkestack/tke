/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
 *
 */

package errors

import (
	apierror "k8s.io/apimachinery/pkg/api/errors"
)

func HandleAPIError(err error) (isErrStatus bool, code int, message string) {
	if err != nil {
		var s *apierror.StatusError
		s, isErrStatus = err.(*apierror.StatusError)
		if isErrStatus {
			st := s.Status()
			code = int(st.Code)
			message = st.Message
		}
	}
	return
}
