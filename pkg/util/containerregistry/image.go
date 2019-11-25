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

package containerregistry

import (
	"bytes"
	"path"
)

var (
	registryDomain    string
	registryNamespace string
)

func Init(domain string, namespace string) {
	registryDomain = domain
	registryNamespace = namespace
}

type Image struct {
	Name string
	Tag  string
}

func (i Image) BaseName() string {
	b := new(bytes.Buffer)
	b.WriteString(i.Name)
	if i.Tag != "" {
		b.WriteString(":" + i.Tag)
	}
	return b.String()
}

func (i Image) FullName() string {
	return path.Join(registryDomain, registryNamespace, i.BaseName())
}

func GetImagePrefix(name string) string {
	return path.Join(registryDomain, registryNamespace, name)
}

func GetPrefix() string {
	return path.Join(registryDomain, registryNamespace)
}
