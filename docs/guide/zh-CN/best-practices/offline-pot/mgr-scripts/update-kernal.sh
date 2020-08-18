#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="all_func"
hosts="all"

help(){
  echo "show usage:"
  echo "update_kernel: update centos kernel for es"
  echo "all_func: execute all function, -f default value is all_func !!!"
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

# update centos kernel for es
update_kernel(){
  echo "###### update centos kernel for es start ######"
  ansible-playbook -f 10 -i ../hosts --tags update_kernel ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### update centos kernel for es end ######"
}


# execute all function
all_func(){
  update_kernel
}

main(){
  $CALL_FUN || help
}
main
