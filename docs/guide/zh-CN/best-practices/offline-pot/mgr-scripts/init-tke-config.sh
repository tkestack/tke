#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="init"
hosts="all"

help(){
  echo "show usage:"
  echo "init: init tke deploy config, -f default value is init."
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

# init tke deploy config
init(){
  echo "###### init tke deploy config start ######"
  ansible-playbook -f 10 -i ../hosts --tags init_tke_cfg ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### init tke deploy config end ######"
}


main(){
  $CALL_FUN || help
}
main
