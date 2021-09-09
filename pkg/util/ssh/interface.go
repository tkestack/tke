/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

package ssh

import (
	"io"
)

type Interface interface {
	Ping() error

	CombinedOutput(cmd string) ([]byte, error)
	Execf(format string, a ...interface{}) (stdout string, stderr string, exit int, err error)
	Exec(cmd string) (stdout string, stderr string, exit int, err error)

	CopyFile(src, dst string) error
	WriteFile(src io.Reader, dst string) error
	ReadFile(filename string) ([]byte, error)
	ReadDir(dirname string) (string, error)
	Exist(filename string) (bool, error)

	LookPath(file string) (string, error)
}
