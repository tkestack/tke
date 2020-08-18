#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="dpl_business"
hosts="all"

help(){
  echo "show usage:"
  echo "dpl_business: deploy business, -f default value is dpl_business"
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

# deploy business
dpl_business(){
  echo "###### deploy business start ######"
  # deploy business
  ansible-playbook -f 10 -i ../hosts --tags dpl_business ../playbooks/business/business.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy business end ######"
}

all_func(){
  dpl_business
}

main(){
  $CALL_FUN || help
}
main
