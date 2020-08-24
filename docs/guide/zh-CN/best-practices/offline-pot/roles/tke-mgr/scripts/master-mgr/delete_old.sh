#!/bin/bash
# Author: kubelouislu
# running on HEALTHY_MASTER_IP
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

while getopts ":d:" opt
do
  case $opt in
    d)
    DOWN_MASTER_IP="${OPTARG}"
    ;;
    ?)
    echo "unkown args! just suport -d[IP address of the master will be down] arg!!!"
    exit 0;;
  esac
done

delete_old(){
  kubectl delete node $DOWN_MASTER_IP
}

main(){
    delete_old
}
main
