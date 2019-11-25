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

# generate-groups generates openapi for a project with external types only.

if [[ "$#" -lt 2 ]]; then
    cat <<EOF
Usage: $(basename "$0") <output-package> <apis-package>
  <output-package>    the output package name (e.g. github.com/example/project/pkg/openapi).
  <apis-package>      the types dir (e.g. github.com/example/project/pkg/apis).
Examples:
  $(basename "$0") openapi github.com/example/project/pkg github.com/example/project/pkg/apis
EOF
  exit 0
fi

OUTPUT_PKG="$1"
APIS_PKG="$2"

GOPATH=${GOPATH:-/go}
K8S_ROOT=${K8S_ROOT:-/go/src/k8s.io/kubernetes}
K8S_BIN=${K8S_ROOT}/_output/bin
PATH=${K8S_BIN}:${PATH}

echo "===========> Generating external openapi codes"
"${K8S_BIN}"/openapi-gen \
    -O zz_generated.openapi \
    --go-header-file "/root/boilerplate.go.txt" \
    --input-dirs "${APIS_PKG}" \
    --output-package "${OUTPUT_PKG}"
