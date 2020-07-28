#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="del_business"
hosts="all"

help(){
  echo "show usage:"
  echo "del_business: remove business, -f default value is del_business"
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

# remove business
del_business(){
  echo "###### remove business start ######"
  ansible-playbook -f 10 -i ../hosts --tags remove_business ../playbooks/business/business.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove business end ######"
}

all_func(){
  del_business
}

main(){
  $CALL_FUN || help
}
main
