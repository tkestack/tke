#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="init"
hosts="all"

help(){
  echo "show usage:"
  echo "init: init tke keepalived config, -f default value is init."
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

# init keepalived config, just tmp use.
init(){
  echo "###### init keepalived config start ######"
  ansible-playbook -f 10 -i ../hosts --tags init_keepalived ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### init keepalived config end ######"
}


main(){
  $CALL_FUN || help
}
main
