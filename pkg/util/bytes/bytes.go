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

package bytes

import (
	"fmt"
	"regexp"
)

// SetContent set content by pattern
func SetContent(data []byte, pattern string, content string) ([]byte, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	var newData string
	if re.Match(data) {
		newData = re.ReplaceAllString(string(data), content)
	} else {
		newData = fmt.Sprintf("%s\n%s", data, content)
	}

	return []byte(newData), nil
}
