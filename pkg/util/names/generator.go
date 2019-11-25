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

package names

import (
	"strings"

	utilrand "k8s.io/apimachinery/pkg/util/rand"
	k8snames "k8s.io/apiserver/pkg/storage/names"
)

// generator generates random names.
type generator struct{}

// Generator is a generator that returns the name plus a random suffix of eight
// alphanumerics when a name is requested.
var Generator k8snames.NameGenerator = generator{}

func (generator) GenerateName(base string) string {
	if strings.HasSuffix(base, "-") {
		return base + utilrand.String(8)
	}
	return base + "-" + utilrand.String(8)
}
