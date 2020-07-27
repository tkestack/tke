#!/bin/bash

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
