#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="tke_gateway"
hosts="all"

help(){
  echo "show usage:"
  echo "tke_gateway: adjust the number of tke-gateway replicas, -f default value is tke_gateway."
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

# adjust the number of tke-gateway replicas
tke_gateway(){
  echo "###### adjust the number of tke-gateway replicas to 1 start ######"
  ansible-playbook -f 10 -i ../hosts --tags tke_gateway ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### adjust the number of tke-gateway replicas to 1 end ######"
}

main(){
  $CALL_FUN || help
}
main
