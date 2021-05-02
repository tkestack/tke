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

package expansion

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// TODO: expand this to rend and copy file
//nolint
func copyFile(src, dst string, perm os.FileMode) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if perm == 0 {
		fi, _ := os.Stat(src)
		perm = fi.Mode()
	}
	err = ioutil.WriteFile(dst, input, perm)
	if err != nil {
		return err
	}

	return nil
}

//nolint
func copyFileToDir(src, dstDir string, perm os.FileMode) error {
	filename := path.Base(src)
	return copyFile(src, fmt.Sprintf("%s/%s", dstDir, filename), perm)
}
