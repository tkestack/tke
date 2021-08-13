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

rm -rf /etc/kubernetes

systemctl stop kubelet 2>/dev/null

docker rm -f $(docker ps -aq) 2>/dev/null
systemctl stop docker 2>/dev/null

ip link del cni0 2>/etc/null

for port in 80 2379 6443 8086 {10249..10259} ; do
    fuser -k -9 ${port}/tcp
done

if [ -d "/var/lib/kubelet/pods" ]; then
  umount $(df -HT | grep '/var/lib/kubelet/pods' | awk '{print $7}')
elif [ -d "/data/kubelet/pods" ]; then
  umount $(df -HT | grep '/data/kubelet/pods' | awk '{print $7}')
fi

rm -rfv /etc/kubernetes || echo "not exist"
rm -rfv /etc/docker || echo "not exist"
rm -rfv /root/.kube/config || echo "not exist"
rm -rfv /var/lib/kubelet || echo "not exist"
rm -rfv /var/lib/cni || echo "not exist"
rm -rfv /etc/cni || echo "not exist"
rm -rfv /var/lib/etcd || echo "not exist"
rm -rfv /tmp/etc/kubernetes/ || echo "not exist"
rm -rfv /usr/libexec/kubernetes/ || echo "not exist"
rm -rfv /etc/keepalived/ || echo "not exist"
rm -rfv /var/lib/postgresql || echo "not exist"
rm -rfv /etc/core/token || echo "not exist"
rm -rfv /var/lib/redis || echo "not exist"
rm -rfv /storage || echo "not exist"
rm -rfv /chart_storage || echo "not exist"

if [ `ls /data/ | wc -l` -gt 0 ]; then
  if [ -L "/storage" ]; then
    rm -rf /storage
  fi
  if [ -L "/var/lib/influxdb" ]; then
    rm -rf /var/lib/influxdb
  fi
  rm -rf /data/*
fi
systemctl start docker 2>/dev/null
