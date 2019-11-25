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

package file

import (
	"io/ioutil"

	"tkestack.io/tke/pkg/util/bytes"
)

// SetContent set file content by pattern
func SetContent(filename string, pattern string, content string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	newData, err := bytes.SetContent(data, pattern, content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, newData, 0644)
}
