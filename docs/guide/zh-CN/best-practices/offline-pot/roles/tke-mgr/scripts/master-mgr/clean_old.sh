#!/bin/bash
# Author: kubelouislu
# running on DOWN_MASTER_IP
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

clean_old(){
  kubeadm reset -f
  iptables -F && iptables -t nat -F && iptables -t mangle -F
  rm -rf $HOME/.kube/config
}

main(){
    clean_old
}
main
