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

REGISTRY_PREFIX=${REGISTRY_PREFIX:-tkestack}
BUILDER=${BUILDER:-default}
VERSION=${VERSION:-$(git describe --dirty --always --tags | sed 's/-/./g')}
INSTALLER=tke-installer-x86_64-$VERSION.run
PROVIDER_RES_VERSION=v1.16.6-4
K8S_VERION=${PROVIDER_RES_VERSION%-*}
DOCKER_VERSION=18.09.9
TARGET_OS=linux
TARGET_ARCH=amd64
TARGET_PLATFORM=${TARGET_OS}-${TARGET_ARCH}
OUTPUT_DIR=_output
DST_DIR=$(mktemp -d)
echo "${DST_DIR}" || exit
SCRIPT_DIR=$(dirname "${BASH_SOURCE[0]}")

function usage() {
  cat <<EOF
Usage: ${0} Release TKE
  -h, help
  -q, quick release
  -t, tke release
EOF
}

function prepare::baremetal_provider() {
  mkdir -p "${DST_DIR}/provider/baremetal/"

  cp -rv pkg/platform/provider/baremetal/conf "${DST_DIR}/provider/baremetal"
  cp -rv pkg/platform/provider/baremetal/manifests "${DST_DIR}/provider/baremetal"
  ls -l "${DST_DIR}/provider/baremetal"

  id=$(docker create "${REGISTRY_PREFIX}/provider-res:${PROVIDER_RES_VERSION}")
  docker cp "$id":/data/res "${DST_DIR}/provider/baremetal/"
  docker rm "$id"
}

function prepare::tke_installer() {
  mkdir -p "${DST_DIR}"/{bin,conf,data,hooks}

  ls -l "${DST_DIR}"

  curl -L "https://storage.googleapis.com/kubernetes-release/release/${K8S_VERION}/bin/${TARGET_OS}/${TARGET_ARCH}/kubectl" -o "${DST_DIR}/bin/kubectl"
  chmod +x "${DST_DIR}/bin/kubectl"

  cp -v "$OUTPUT_DIR/${TARGET_OS}/${TARGET_ARCH}/tke-installer" "${DST_DIR}/bin"
  cp -rv cmd/tke-installer/app/installer/manifests "${DST_DIR}"
  cp -rv cmd/tke-installer/app/installer/hooks "${DST_DIR}"
  cp -rv "${SCRIPT_DIR}/certs" "${DST_DIR}"
  cp -rv "${SCRIPT_DIR}/.docker" "${DST_DIR}"
}

function build::installer_image() {
  docker build --pull -t "${REGISTRY_PREFIX}/tke-installer:$VERSION" -f "${SCRIPT_DIR}/Dockerfile" "${DST_DIR}"
}

function build::installer() {
    installer_dir=$(mktemp -d)
    echo "installer dir: ${installer_dir}" || exit

    mkdir -p "${installer_dir}/res"

    cp -v build/docker/tools/tke-installer/{build.sh,init_installer.sh,install.sh} "${installer_dir}"/
    cp -v "${DST_DIR}/provider/baremetal/res/${TARGET_PLATFORM}/docker-${TARGET_PLATFORM}-${DOCKER_VERSION}.tar.gz" \
          "${installer_dir}/res/docker.tgz"
    cp -v pkg/platform/provider/baremetal/conf/docker/docker.service "${installer_dir}/res/"
    docker save "${REGISTRY_PREFIX}/tke-installer:$VERSION" | gzip -c > "${installer_dir}/res/tke-installer.tgz"

    sed -i "s;VERSION=.*;VERSION=$VERSION;g" "${installer_dir}/install.sh"

    "${installer_dir}/build.sh" "$INSTALLER"
    cp -v "${installer_dir}/$INSTALLER" $OUTPUT_DIR

    echo "build tke-installer success! OUTPUT => $OUTPUT_DIR/$INSTALLER"
    (cd $OUTPUT_DIR && sha256sum "$INSTALLER" > "$INSTALLER.sha256")

    if [[ "${BUILDER}" == "tke" ]]; then
      coscmd upload "${installer_dir}/$INSTALLER" "$INSTALLER"
      coscmd upload "$OUTPUT_DIR/$INSTALLER.sha256" "$INSTALLER.sha256"
    fi

    rm -rf "${installer_dir}"
}

function prepare::images() {
  GENERATE_IMAGES_BIN="$OUTPUT_DIR/$(go env GOOS)/$(go env GOARCH)/generate-images"
  make build BINS=generate-images VERSION="$VERSION"

  $GENERATE_IMAGES_BIN
  $GENERATE_IMAGES_BIN | sed "s;^;${REGISTRY_PREFIX}/;" | xargs -n1 -I{} sh -c "docker pull {} || exit 255"
  $GENERATE_IMAGES_BIN | sed "s;^;${REGISTRY_PREFIX}/;" | xargs docker save | gzip -c >"${DST_DIR}"/images.tar.gz
}

pwd

quick=false
while getopts "hq" o; do
    case "${o}" in
        h)
          usage
          ;;
        q)
          quick=true
          ;;
        *)
          usage
          ;;
    esac
done
shift $((OPTIND-1))

make build BINS="tke-installer" OS=${TARGET_OS} ARCH=${TARGET_ARCH} VERSION="$VERSION"

prepare::baremetal_provider
prepare::tke_installer
if [[ "${quick}" == "false" ]]; then
  prepare::images
fi

build::installer_image
build::installer

rm -rf "${DST_DIR}"
