#!/usr/bin/env bash

# Tencent is pleased to support the open source community by making TKEStack
# available.
#
# Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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

# generate-groups generates everything for a project with external types only, e.g. a project based
# on CustomResourceDefinitions.

if [[ "$#" -lt 4 ]] || [[ "${1}" == "--help" ]]; then
  cat <<EOF
Usage: $(basename "$0") <generators> <output-package> <internal-apis-package> <extensiona-apis-package> <groups-versions> ...
  <generators>        the generators comma separated to run (e.g. deepcopy-external,defaulter-external,client-external,
                      lister-external,informer-external,deepcopy-internal,defaulter-internal,client-internal,
                      lister-internal,informer-internal or all-external,all-internal,all).
  <output-package>    the output package name (e.g. github.com/example/project/pkg/generated).
  <int-apis-package>  the internal types dir (e.g. github.com/example/project/pkg/apis).
  <ext-apis-package>  the external types dir (e.g. github.com/example/project/pkg/apis or githubcom/example/apis).
  <groups-versions>   the groups and their versions in the format "groupA:v1,v2 groupB:v1 groupC:v2", relative
                      to <api-package>.
  ...                 arbitrary flags passed to all generator binaries.
Examples:
  $(basename "$0") all-external github.com/example/project/pkg/client github.com/example/project/pkg/apis github.com/example/project/apis "foo:v1 bar:v1alpha1,v1beta1"
  $(basename "$0") deepcopy-external,client-external github.com/example/project/pkg/client github.com/example/project/pkg/apis github.com/example/project/apis "foo:v1 bar:v1alpha1,v1beta1"
  $(basename "$0") all-internal github.com/example/project/pkg/client github.com/example/project/pkg/apis github.com/example/project/apis "foo:v1 bar:v1alpha1,v1beta1"
  $(basename "$0") deepcopy-internal,defaulter-internal,conversion-internal github.com/example/project/pkg/client github.com/example/project/pkg/apis github.com/example/project/apis "foo:v1 bar:v1alpha1,v1beta1"
EOF
  exit 0
fi

GENS="$1"
OUTPUT_PKG="$2"
INT_APIS_PKG="$3"
EXT_APIS_PKG="$4"
GROUPS_WITH_VERSIONS="$5"
shift 5

GOPATH=${GOPATH:-/go}
K8S_ROOT=${K8S_ROOT:-/go/src/k8s.io/kubernetes}
K8S_BIN=${K8S_ROOT}/_output/bin
PATH=${K8S_BIN}:${PATH}

function codegen_join() { local IFS="$1"; shift; echo "$*"; }

# Generates types_swagger_doc_generated file for the given group version.
# $1: Name of the group version
# $2: Path to the directory where types.go for that group version exists. This
# is the directory where the file will be generated.
function gen_types_swagger_doc() {
  local group_version=$1
  local gv_dir=$2

  TMP_FILE="${TMPDIR:-/tmp}/types_swagger_doc_generated.$(date +%s).go"

  echo "===========> Generating swagger type docs for ${group_version} at ${gv_dir}"

    {
    echo -e "$(cat /root/boilerplate.go.txt)\n"
    echo "package ${group_version##*/}"
    cat <<EOF

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-generated-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE. DO NOT EDIT.
EOF
  } > "${TMP_FILE}"

  "${K8S_BIN}"/genswaggertypedocs -s \
    "${gv_dir}/types.go" \
    -f - \
    >>  "$TMP_FILE"

  echo "// AUTO-GENERATED FUNCTIONS END HERE" >> "$TMP_FILE"

  gofmt -w -s "$TMP_FILE"
  mv "$TMP_FILE" "${gv_dir}"/types_swagger_doc_generated.go
}

# enumerate group versions
ALL_FQ_APIS=(${ALL_FQ_APIS:-}) # e.g. k8s.io/kubernetes/pkg/apis/apps k8s.io/api/apps/v1
INT_FQ_APIS=(${INT_FQ_APIS:-}) # e.g. k8s.io/kubernetes/pkg/apis/apps
EXT_FQ_APIS=(${EXT_FQ_APIS:-}) # e.g. k8s.io/api/apps/v1
EXT_PB_APIS=(${EXT_PB_APIS:-}) # e.g. k8s.io/api/apps/v1

for GVs in ${GROUPS_WITH_VERSIONS}; do
  IFS=: read -r G Vs <<<"${GVs}"

  if [[ -n "${INT_APIS_PKG}" ]]; then
    ALL_FQ_APIS+=("${INT_APIS_PKG}/${G}")
    INT_FQ_APIS+=("${INT_APIS_PKG}/${G}")
  fi

  # enumerate versions
  for V in ${Vs//,/ }; do
    ALL_FQ_APIS+=("${EXT_APIS_PKG}/${G}/${V}")
    EXT_FQ_APIS+=("${EXT_APIS_PKG}/${G}/${V}")
  done
done

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-external" ]] || grep -qw "deepcopy-external" <<<"${GENS}"; then
  echo "===========> Generating external deepcopy funcs"
  "${K8S_BIN}"/deepcopy-gen \
        --go-header-file /root/boilerplate.go.txt \
        --input-dirs "$(codegen_join , "${EXT_FQ_APIS[@]}")" \
        -O zz_generated.deepcopy \
        --bounding-dirs "${EXT_APIS_PKG}" \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-external" ]] || grep -qw "client-external" <<<"${GENS}"; then
  echo "===========> Generating external clientset for ${GROUPS_WITH_VERSIONS} at ${OUTPUT_PKG}/clientset"
  "${K8S_BIN}"/client-gen \
        --go-header-file /root/boilerplate.go.txt \
        --clientset-name versioned \
        --input-base "" \
        --input "$(codegen_join , "${EXT_FQ_APIS[@]}")" \
        --output-package "${OUTPUT_PKG}"/clientset \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-external" ]] || grep -qw "lister-external" <<<"${GENS}"; then
  echo "===========> Generating external listers for ${GROUPS_WITH_VERSIONS} at ${OUTPUT_PKG}/listers"
  "${K8S_BIN}"/lister-gen \
        --go-header-file /root/boilerplate.go.txt \
        --input-dirs "$(codegen_join , "${EXT_FQ_APIS[@]}")" \
        --output-package "${OUTPUT_PKG}"/listers \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-external" ]] || grep -qw "informer-external" <<<"${GENS}"; then
  echo "===========> Generating external informers for ${GROUPS_WITH_VERSIONS} at ${OUTPUT_PKG}/informers"
  "${K8S_BIN}"/informer-gen \
        --go-header-file /root/boilerplate.go.txt \
        --input-dirs "$(codegen_join , "${EXT_FQ_APIS[@]}")" \
        --versioned-clientset-package "${OUTPUT_PKG}"/clientset/versioned \
        --listers-package "${OUTPUT_PKG}"/listers \
        --output-package "${OUTPUT_PKG}"/informers \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-internal" ]] || grep -qw "deepcopy-internal" <<<"${GENS}"; then
  echo "===========> Generating internal deepcopy funcs"
  "${K8S_BIN}"/deepcopy-gen \
        --go-header-file /root/boilerplate.go.txt \
        --input-dirs "$(codegen_join , "${ALL_FQ_APIS[@]}")" \
        -O zz_generated.deepcopy \
        --bounding-dirs "${INT_APIS_PKG}","${EXT_APIS_PKG}" \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-external" ]] || grep -qw "defaulter-external" <<<"${GENS}"; then
  echo "===========> Generating external defaulters"
  "${K8S_BIN}"/defaulter-gen \
        --go-header-file /root/boilerplate.go.txt \
        --input-dirs "$(codegen_join , "${EXT_FQ_APIS[@]}")" \
        -O zz_generated.defaults \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-external" ]] || grep -qw "conversion-external" <<<"${GENS}"; then
  echo "===========> Generating external conversions"
  "${K8S_BIN}"/conversion-gen \
        --go-header-file /root/boilerplate.go.txt \
        --input-dirs "$(codegen_join , "${ALL_FQ_APIS[@]}")" \
        -O zz_generated.conversion \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-internal" ]] || grep -qw "client-internal" <<<"${GENS}"; then
  echo "===========> Generating internal clientset for ${GROUPS_WITH_VERSIONS} at ${OUTPUT_PKG}/clientset"
  if [[ -n "${INT_APIS_PKG}" ]]; then
    IFS=" " read -r -a APIS <<< "$(printf '%s/ ' "${INT_FQ_APIS[@]}")"
    "${K8S_BIN}"/client-gen \
            --go-header-file /root/boilerplate.go.txt \
            --clientset-name internalversion \
            --input-base "" \
            --input "$(codegen_join , "${APIS[@]}")" \
            --output-package "${OUTPUT_PKG}"/clientset \
            "$@"
  fi
  "${K8S_BIN}"/client-gen \
        --go-header-file /root/boilerplate.go.txt \
        --clientset-name versioned \
        --input-base "" \
        --input "$(codegen_join , "${EXT_FQ_APIS[@]}")" \
        --output-package "${OUTPUT_PKG}"/clientset \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-internal" ]] || grep -qw "lister-internal" <<<"${GENS}"; then
  echo "===========> Generating internal listers for ${GROUPS_WITH_VERSIONS} at ${OUTPUT_PKG}/listers"
  "${K8S_BIN}"/lister-gen \
        --go-header-file /root/boilerplate.go.txt \
        --input-dirs "$(codegen_join , "${ALL_FQ_APIS[@]}")" \
        --output-package "${OUTPUT_PKG}"/listers \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-internal" ]] || grep -qw "informer-internal" <<<"${GENS}"; then
  echo "===========> Generating informers for ${GROUPS_WITH_VERSIONS} at ${OUTPUT_PKG}/informers"
  "${K8S_BIN}"/informer-gen \
        --go-header-file /root/boilerplate.go.txt \
        --input-dirs "$(codegen_join , "${ALL_FQ_APIS[@]}")" \
        --versioned-clientset-package "${OUTPUT_PKG}"/clientset/versioned \
        --internal-clientset-package "${OUTPUT_PKG}"/clientset/internalversion \
        --listers-package "${OUTPUT_PKG}"/listers \
        --output-package "${OUTPUT_PKG}"/informers \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-external" ]] || grep -qw "protobuf-external" <<<"${GENS}"; then
  echo "===========> Generating external protobuf codes"
  "${K8S_BIN}"/go-to-protobuf \
        --go-header-file "/root/boilerplate.go.txt" \
        --proto-import "${K8S_ROOT}/vendor" \
        --proto-import "${K8S_ROOT}/third_party/protobuf" \
        --packages "$(codegen_join , "${EXT_FQ_APIS[@]}" "${EXT_PB_APIS[@]}")" \
        "$@"
fi

if [[ "${GENS}" = "all" ]] || [[ "${GENS}" = "all-external" ]] || grep -qw "swagger-external" <<<"${GENS}"; then
  # To avoid compile errors, remove the currently existing files.
  for group_version in "${EXT_FQ_APIS[@]}"; do
    rm -f "${GOPATH}"/src/"${group_version}"/types_swagger_doc_generated.go
  done
  for group_version in "${EXT_FQ_APIS[@]}"; do
    gen_types_swagger_doc "${group_version}" "${GOPATH}/src/${group_version}"
  done
  gen_types_swagger_doc "${EXT_SWAGGER_API}" "${GOPATH}/src/${EXT_SWAGGER_API}"
fi
