#!/bin/bash
# Author: kubelouislu
# running on HEALTHY_MASTER_IP
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

get_args(){
  ORIGINAL_NAME="$(hostnamectl status|grep 'Static hostname' | awk -F':' '{print $2}')"
  TOKEN="$(kubeadm token list|grep forever|awk '{print $1}')"
  HASH="$(openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //')"
  hostnamectl set-hostname temp-node-name
  if [ -f "temp" ]; then
    rm -rf temp
  fi
  kubeadm init phase upload-certs --upload-certs >temp
  CERT="$(tail -1 temp)"
  rm -rf temp
  if [ -f "/root/args" ]; then
    rm -rf /root/args
  fi
  mkdir -p /root/args
  echo $CERT > /root/args/cert
  echo $TOKEN > /root/args/token
  echo $HASH > /root/args/hash
  hostnamectl set-hostname $ORIGINAL_NAME
}

main(){
    get_args
}
main
