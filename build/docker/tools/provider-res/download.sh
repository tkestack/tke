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

CNI_VERSION=v0.7.5
DOCKER_VERSION=18.09.9
KUBEADM_VERSION=v1.15.1
KUBERNETES_VERSION=v1.14.6
NVIDIA_DRIVER_VERSION=440.31

cd "$DST_DIR" || exit

wget https://github.com/containernetworking/plugins/releases/download/$CNI_VERSION/cni-plugins-amd64-$CNI_VERSION.tgz

wget https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_VERSION.tgz

wget https://dl.k8s.io/$KUBERNETES_VERSION/kubernetes-node-linux-amd64.tar.gz -O kubernetes-node-linux-amd64-$KUBERNETES_VERSION.tar.gz

wget https://storage.googleapis.com/kubernetes-release/release/$KUBEADM_VERSION/bin/linux/amd64/kubeadm
chmod +x kubeadm
tar cvzf kubeadm-$KUBEADM_VERSION.tar.gz kubeadm
rm kubeadm

wget https://us.download.nvidia.cn/XFree86/Linux-x86_64/$NVIDIA_DRIVER_VERSION/NVIDIA-Linux-x86_64-$NVIDIA_DRIVER_VERSION.run
chmod +x NVIDIA-Linux-x86_64-$NVIDIA_DRIVER_VERSION.run

