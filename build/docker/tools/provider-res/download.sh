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

#CNI_VERSION=v0.7.5
CNI_VERSION=v0.8.5
DOCKER_VERSION=18.09.9
## Update Kubeadm version may need to 
## update Kubernetes at the same time
KUBEADM_VERSION=v1.15.1
KUBERNETES_VERSION=(v1.14.10 v1.16.6)
NVIDIA_DRIVER_VERSION=440.31

OS_LIST=(linux)
ARCH_LIST=(amd64 arm64)

cd "$DST_DIR" || exit

for os in "${OS_LIST[@]}"
do
  for arch in "${ARCH_LIST[@]}"
  do
    platform=${os}-${arch}
    mkdir -p ${platform}
    cd ./${platform}
    echo "Switch to ${platform}"

    ## CNI Plugins
    wget https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-${platform}-${CNI_VERSION}.tgz \
          -O cni-plugins-${platform}-${CNI_VERSION}.tar.gz

    ## Docker
    if [ x${arch} == x"amd64" ]; then
      docker_arch=x86_64
    elif [ x${arch} == x"arm64" ]; then
      docker_arch=aarch64
    else
      echo "[ERROR] Fail to get docker ${DOCKER_VERSION} on ${platform} platform."
      exit
    fi
    wget https://download.docker.com/${os}/static/stable/${docker_arch}/docker-${DOCKER_VERSION}.tgz \
          -O docker-${platform}-${DOCKER_VERSION}.tar.gz

    ## Kubernetes
    for version in "${KUBERNETES_VERSION[@]}"
    do
      wget https://dl.k8s.io/${version}/kubernetes-node-${platform}.tar.gz \
            -O kubernetes-node-${platform}-${version}.tar.gz
    done

    ## Kubeadm
    wget https://storage.googleapis.com/kubernetes-release/release/${KUBEADM_VERSION}/bin/${os}/${arch}/kubeadm
    chmod +x kubeadm
    tar cvzf kubeadm-${platform}-${KUBEADM_VERSION}.tar.gz kubeadm
    rm kubeadm

    ## NVIDIA driver
    if [ x${os} == x"linux" ] && [ x${arch} == x"amd64" ]; then
      wget https://us.download.nvidia.cn/XFree86/Linux-x86_64/${NVIDIA_DRIVER_VERSION}/NVIDIA-Linux-x86_64-${NVIDIA_DRIVER_VERSION}.run
      chmod +x NVIDIA-Linux-x86_64-${NVIDIA_DRIVER_VERSION}.run
      mv NVIDIA-Linux-x86_64-${NVIDIA_DRIVER_VERSION}.run NVIDIA-${platform}-${NVIDIA_DRIVER_VERSION}.run
    else
      ## NVIDIA driver only has arm32 version. 
      ## NVIDIA_DRIVER_ARM32_VERSION=390.132
      ## https://us.download.nvidia.cn/XFree86/Linux-x86-ARM/${NVIDIA_DRIVER_ARM32_VERSION}/NVIDIA-Linux-armv7l-gnueabihf-${NVIDIA_DRIVER_ARM32_VERSION}.run
      echo "[WARN] Cannot get NVIDIA driver ${NVIDIA_DRIVER_VERSION} on ${platform} platform."
    fi
    
    ## 
    cd ../
    echo "Done. Exit ${platform}"
  done ## End of {ARCH_LIST}
done ## End of {OS_LIST}

echo "Finish to download binaries."
