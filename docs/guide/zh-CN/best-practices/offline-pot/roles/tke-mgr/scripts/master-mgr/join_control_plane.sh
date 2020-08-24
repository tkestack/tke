#!/bin/bash
# Author: kubelouislu
# running on NEW_MASTER_IP
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

while getopts ":n:h:" opt
do
  case $opt in
    n)
    NEW_MASTER_IP="${OPTARG}"
    ;;
    h)
    HEALTHY_MASTER_IP="${OPTARG}"
    ;;
    ?)
    echo "unkown args! just suport -n[IP address of the server will be up] -h[IP address of the master that is stable in cluster] arg!!!"
    exit 0;;
  esac
done

join_control_plane(){
  if [ ! -d "/root/temp/" ]; then
    echo "the config files do not exist please check the step before !!!" && exit 0
  fi
  TOKEN=`cat /root/temp/token`
  HASH=`cat /root/temp/hash`
  CERT=`cat /root/temp/cert`
  IP_ADDR=$NEW_MASTER_IP
  # turn the hostname back
  kubeadm join $HEALTHY_MASTER_IP:6443 --token $TOKEN --discovery-token-ca-cert-hash sha256:$HASH --control-plane --certificate-key $CERT --node-name $IP_ADDR
  # join in the control plane with token hash and secret for certs
  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config
  # set up the master
  # some files are needed for some kube components
}

main(){
    join_control_plane
}
main
