#!/usr/bin/env bash

# Tencent is pleased to support the open source community by making TKEStack
# available.
#
# Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
PROVIDER_RES_VERSION=v1.21.4-5
K8S_VERSION=${PROVIDER_RES_VERSION%-*}
DOCKER_VERSION=19.03.14
CONTAINERD_VERSION=1.5.4
NERDCTL_VERSION=0.11.0
REGISTRY_VERSION=2.7.1
OSS=(linux)
ARCHS=(amd64 arm64)
OUTPUT_DIR=_output
DST_DIR=$(mktemp -d)
echo "${DST_DIR}" || exit
INSTALLER_DIR=$(mktemp -d)
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
  local -r os="$1"
  local -r arch="$2"

  mkdir -p "${DST_DIR}"/{bin,conf,data,hooks}

  ls -l "${DST_DIR}"

  curl -L "https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/${os}/${arch}/kubectl" -o "${DST_DIR}/bin/kubectl"
  chmod +x "${DST_DIR}/bin/kubectl"

  curl -L "https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/public.charts.tar.gz" -o "${DST_DIR}/public.charts.tar.gz"

  make build BINS="tke-installer" OS="${os}" ARCH="${arch}" VERSION="${VERSION}"

  cp -v "$OUTPUT_DIR/${os}/${arch}/tke-installer" "${DST_DIR}/bin"
  cp -rv cmd/tke-installer/app/installer/manifests "${DST_DIR}"
  cp -rv cmd/tke-installer/app/installer/hooks "${DST_DIR}"
  cp -rv "${SCRIPT_DIR}/certs" "${DST_DIR}"
  cp -rv "${SCRIPT_DIR}/.docker" "${DST_DIR}"
  cp -rv "${SCRIPT_DIR}/run.sh" "${DST_DIR}"
  cp -rv "${SCRIPT_DIR}/dockerd-entrypoint.sh" "${DST_DIR}"

  make web.build.installer
  cp -rv web/installer/build  "${DST_DIR}/assets"
}

function build::installer_image() {
  local -r arch="$1"

  docker build --platform="${arch}" --build-arg ENV_ARCH="${arch}" --pull -t "${REGISTRY_PREFIX}/tke-installer-${arch}:$VERSION" -f "${SCRIPT_DIR}/Dockerfile" "${DST_DIR}"
}

function build::installer() {
    local -r os="$1"
    local -r arch="$2"
    local -r target_platform="${os}-${arch}"
    local -r installer="tke-installer-${target_platform}-${VERSION}.run"

    echo "build ${installer} in dir: ${INSTALLER_DIR}" || exit

    mkdir -p "${INSTALLER_DIR}/res"

    cp -v build/docker/tools/tke-installer/{build.sh,init_installer.sh,install.sh} "${INSTALLER_DIR}"/

    cp -v "${DST_DIR}/provider/baremetal/res/${target_platform}/docker-${target_platform}-${DOCKER_VERSION}.tar.gz" \
          "${INSTALLER_DIR}/res/docker.tgz"
    cp -v pkg/platform/provider/baremetal/conf/docker/docker.service "${INSTALLER_DIR}/res/"
    cp -v build/docker/tools/tke-installer/daemon.json "${INSTALLER_DIR}/res/"
    cp -v "${DST_DIR}/provider/baremetal/res/${target_platform}/containerd-${target_platform}-${CONTAINERD_VERSION}.tar.gz" "${INSTALLER_DIR}/res/containerd.tar.gz"
    cp -v "${DST_DIR}/provider/baremetal/res/${target_platform}/nerdctl-${target_platform}-${NERDCTL_VERSION}.tar.gz" "${INSTALLER_DIR}/res/nerdctl.tar.gz"

    docker save "${REGISTRY_PREFIX}/tke-installer-${arch}:$VERSION" -o "${INSTALLER_DIR}/res/tke-installer.tar"
    docker --config=${DOCKER_PULL_CONFIG} pull "${REGISTRY_PREFIX}/registry-${arch}:$REGISTRY_VERSION"
    docker save "${REGISTRY_PREFIX}/registry-${arch}:$REGISTRY_VERSION" -o "${INSTALLER_DIR}/res/registry.tar"

    sed -i "s;VERSION=.*;VERSION=$VERSION;g" "${INSTALLER_DIR}/install.sh"
    sed -i "s;REGISTRY_VERSION=.*;REGISTRY_VERSION=$REGISTRY_VERSION;g" "${INSTALLER_DIR}/install.sh"

    "${INSTALLER_DIR}/build.sh" "${installer}"
    cp -v "${INSTALLER_DIR}/${installer}" $OUTPUT_DIR

    echo "build tke-installer success! OUTPUT => $OUTPUT_DIR/${installer}"
    (cd $OUTPUT_DIR && sha256sum "${installer}" > "${installer}.sha256")

    echo "current builder is ${BUILDER}"
    if [[ "${BUILDER}" == "tke" ]]; then
      echo "start upload"
      max_tries=3
      for i in $(seq 1 $max_tries); do
        coscmd upload "${INSTALLER_DIR}/$installer" "$installer" && coscmd upload "$OUTPUT_DIR/$installer.sha256" "$installer.sha256"
        result=$?
        if [[ $result -eq 0 ]]; then
          echo "upload successful"
          break
        else
          echo "upload failed"
          sleep 1
        fi
      done
    fi
}

function prepare::images() {
  GENERATE_IMAGES_BIN="$OUTPUT_DIR/$(go env GOOS)/$(go env GOARCH)/generate-images"
  make build BINS=generate-images VERSION="$VERSION"

  $GENERATE_IMAGES_BIN
  $GENERATE_IMAGES_BIN | sed "s;^;${REGISTRY_PREFIX}/;" | xargs -n1 -I{} sh -c "docker --config=${DOCKER_PULL_CONFIG} pull {} || exit 255"
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

if [[ "${quick}" == "false" ]]; then
  prepare::images
fi

prepare::baremetal_provider
for os in "${OSS[@]}"
do
  for arch in "${ARCHS[@]}"
  do
    prepare::tke_installer "${os}" "${arch}"
    build::installer_image "${arch}"
    build::installer "${os}" "${arch}"
  done
done

rm -rf "${DST_DIR}"
rm -rf "${INSTALLER_DIR}"
