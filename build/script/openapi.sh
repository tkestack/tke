#!/usr/bin/env bash

# Tencent is pleased to support the open source community by making TKEStack
# available.
#
# Copyright (C) 2012-2019 Tencent. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use
# this file except in compliance with the License. You may obtain a copy of the
# License at
#
# https://opensource.org/licenses/Apache-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OF ANY KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

API_PACKAGE="tkestack.io/tke/api"

API_MACHINERY_DIR=$(go list -f '{{ .Dir }}' -m k8s.io/apimachinery)
API_DIR=$(go list -f '{{ .Dir }}' -m k8s.io/api)

# kubernetes api machinery
api_machinery=$(
  grep --color=never -rl '+k8s:openapi-gen=' "${API_MACHINERY_DIR}" | \
  xargs -n1 dirname | \
  sed "s,^${API_MACHINERY_DIR}/,k8s.io/apimachinery/," | \
  sort -u
)

# kubernetes api
api=$(
  grep --color=never --exclude-dir=origin -rl '+k8s:openapi-gen=' "${API_DIR}" | \
  xargs -n1 dirname | \
  sed "s,^${API_DIR}/,k8s.io/api/," | \
  sort -u
)

input_dirs=(
  ${api_machinery}
  ${api}
  "${API_PACKAGE}"/platform/v1
  "${API_PACKAGE}"/business/v1
  "${API_PACKAGE}"/notify/v1
  "${API_PACKAGE}"/registry/v1
  "${API_PACKAGE}"/monitor/v1
)

echo "$(IFS=,; echo "${input_dirs[*]}")"
