#!/bin/bash

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

#!/bin/bash

# common
kubeadm reset -f
rm -fv /root/.kube/config
rm -rfv /etc/kubernetes
rm -rfv /var/lib/kubelet
rm -rfv /var/lib/etcd
rm -rfv /var/lib/cni
rm -rfv /etc/cni
rm -rfv /var/lib/tke-registry-api
rm -rfv /opt/tke-installer
rm -rfv /var/lib/postgresql /etc/core/token /var/lib/redis /storage /chart_storage
ip link del cni0 2>/etc/null

for port in 80 443 2379 2380 6443 8086 8181 9100 30086 31138 31180 31443  {10249..10259} ; do
    fuser -k -9 ${port}/tcp
done

# docker
docker rm -f $(docker ps -aq) 2>/dev/null
systemctl disable docker 2>/dev/null
systemctl stop docker 2>/dev/null
rm -rfv /etc/docker
ip link del docker0 2>/etc/null

# containerd
nerdctl rm -f $(nerdctl ps -aq) 2>/dev/null
ip netns list | cut -d' ' -f 1 | xargs -n1 ip netns delete 2>/dev/null
systemctl disable containerd 2>/dev/null
systemctl stop containerd 2>/dev/null
rm -rfv /var/lib/nerdctl/*

## ip link
ip link delete cilium_net 2>/dev/null
ip link delete cilium_vxlan 2>/dev/null
ip link delete flannel.1 2>/dev/null

## iptables
iptables --flush
iptables --flush --table nat
iptables --flush --table filter
iptables --table nat --delete-chain
iptables --table filter --delete-chain

# reboot
reboot now