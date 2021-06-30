#!/bin/bash

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

rm -rf /etc/kubernetes

systemctl stop kubelet 2>/dev/null

crictl rm -f $(docker ps -aq) 2>/dev/null
systemctl stop containerd 2>/dev/null

ip link del cni0 2>/etc/null

for port in 80 2379 6443 8086 {10249..10259} ; do
    fuser -k -9 ${port}/tcp
done

rm -rfv /etc/kubernetes
rm -rfv /etc/containerd
rm -fv /root/.kube/config
rm -rfv /var/lib/kubelet
rm -rfv /var/lib/cni
rm -rfv /etc/cni
rm -rfv /var/lib/etcd
rm -rfv /var/lib/postgresql /etc/core/token /var/lib/redis /storage /chart_storage

systemctl start containerd 2>/dev/null
