#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="clean_cluster"
hosts="all"

help(){
  echo "show usage:"
  echo "clean_cluster: clean cluster nodes,will be delete /data/* dir,note: backup, -f default value is clean_cluster."
  exit 0
}

while getopts ":f:h:" opt
do
  case $opt in
    f)
    CALL_FUN="${OPTARG}"
    ;;
    h)
    hosts="${OPTARG}"
    ;;
    ?)
    echo "unkown args! just suport -f[call function] and -h[ansible hosts group] arg!!!"
    exit 0;;
  esac
done

# clean cluster nodes
clean_cluster(){
  echo "###### clean cluster nodes start ######"
  ansible-playbook -f 10 -i ../hosts --tags remove_cluster ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### clean cluster nodes end ######"
}

all_func(){
  clean_cluster
}

main(){
  $CALL_FUN || help
}
main
