#!/bin/bash
# Author: kubelouislu
# running on NEW_MASTER_IP
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

init_server(){
  if [ -f "/etc/kubernetes/pki/etcd/" ]; then
    rm -rf /etc/kubernetes/pki/etcd/
  fi
  mkdir -p /etc/kubernetes/pki/etcd/
  if [ -f "/var/lib/etcd/" ]; then
    rm -rf /var/lib/etcd/
  fi
  mkdir -p /var/lib/etcd/
  # copy the conf belonging to kubelet
  kubeadm reset -f
  iptables -F && iptables -t nat -F && iptables -t mangle -F
  rm -rf $HOME/.kube/config
  # reset the new machine if it's a node but just in test env
}

main(){
    init_server
}
main
