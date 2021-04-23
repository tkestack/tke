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
set -o xtrace

SCRIPT_DIR=$(realpath $(dirname "${BASH_SOURCE[0]}"))

# PKGS for extract packages in centos:7
PKGS="conntrack-tools"

declare -A archMap=(
  [amd64]=x86_64
  [arm64]=arm64
)

cd "$DST_DIR" || exit

function download::cni_plugins() {
  for version in ${CNI_PLUGINS_VERSIONS}; do
    wget -c "https://github.com/containernetworking/plugins/releases/download/${version}/cni-plugins-${platform}-${version}.tgz" \
      -O "cni-plugins-${platform}-${version}.tar.gz"
  done
}

function download::docker() {
  if [ "${arch}" == "amd64" ]; then
    docker_arch=x86_64
  elif [ "${arch}" == "arm64" ]; then
    docker_arch=aarch64
  else
    echo "[ERROR] Fail to get docker ${arch} on ${platform} platform."
    exit 255
  fi

  for version in ${DOCKER_VERSIONS}; do
    wget -c "https://download.docker.com/${os}/static/stable/${docker_arch}/docker-${version}.tgz" \
      -O "docker-${platform}-${version}.tar.gz"
  done
}

function download::kubernetes() {
  for version in ${K8S_VERSIONS}; do
    result_tke=$(echo ${version} | grep "tke")
    if [ -z "$result_tke" ]; then
      wget -c "https://dl.k8s.io/${version}/kubernetes-node-${platform}.tar.gz" \
        -O "kubernetes-node-${platform}-${version}.tar.gz"
    else
      wget -c "https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/kubernetes-node-linux-amd64-${version}.tar.gz" \
        -O "kubernetes-node-linux-amd64-${version}.tar.gz"
    fi
  done
}

function download::nvidia_driver() {
  if ! { [ "${os}" == "linux" ] && [ "${arch}" == "amd64" ]; }; then
    return
  fi

  for version in ${NVIDIA_DRIVER_VERSIONS}; do
    wget -c "https://us.download.nvidia.cn/XFree86/Linux-x86_64/${version}/NVIDIA-Linux-x86_64-${version}.run" \
      -O "NVIDIA.run"
    chmod +x "NVIDIA.run"
    GZIP=-n tar cvzf "NVIDIA-${platform}-${version}.tar.gz" "NVIDIA.run"
    rm "NVIDIA.run"
  done
}

function download::nvidia_container_runtime() {
  if ! { [ "${os}" == "linux" ] && [ "${arch}" == "amd64" ]; }; then
    return
  fi

  for version in ${NVIDIA_CONTAINER_RUNTIME_VERSIONS}; do
    wget -c "https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/res/${os}/${arch}/nvidia-container-runtime-${platform}-${version}.tgz" \
      -O "nvidia-container-runtime-${platform}-${version}.tar.gz"
  done
}

function download::pkgs() {
  if [ -z ${archMap[${arch}]+unset} ]; then
    echo "ERROR: unsupport arch ${arch}"
    exit 1
  fi
  docker_arch=${archMap[${arch}]}
  docker pull --platform=${docker_arch} centos:7
  for pkg in ${PKGS}; do
    docker run --platform="${docker_arch}" -e OS="${os}" -e ARCH="${arch}" -e PKG="${pkg}" --rm -v"${SCRIPT_DIR}":/tmp/bin -v$(realpath $(pwd)):/output centos:7 /tmp/bin/run.sh
  done
}

echo "Starting to download resources..."

for os in ${OSS}; do
  for arch in ${ARCHS}; do
    platform=${os}-${arch}
    mkdir -p "${platform}"
    cd "${platform}"

    download::cni_plugins
    download::docker
    download::kubernetes
    download::nvidia_driver
    download::nvidia_container_runtime
    download::pkgs

    cd -
  done
done

echo "Finish to download resources."
