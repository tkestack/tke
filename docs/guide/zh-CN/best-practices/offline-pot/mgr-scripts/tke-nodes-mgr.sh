#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="add_tke_nodes"
hosts="all"

help(){
  echo "show usage:"
  echo "add_tke_nodes: add tke nodes, -f default value is add_tke_nodes."
  echo "remove_tke_nodes: remove tke nodes"
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

# tke node init
tke_node_init(){
  echo "###### tke node init start ######"
  ansible-playbook -f 10 -i ../hosts --tags tke_node_init ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### tke node init end ######"
}

# add tke nodes
add_tke_nodes(){
  echo "###### add tke nodes start ######"
  ansible-playbook -f 10 -i ../hosts --tags add_tke_node ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### add tke nodes end ######"
}

# remove tke nodes
remove_tke_nodes(){
  echo "###### remove tke nodes start ######"
  ansible-playbook -f 10 -i ../hosts --tags remove_tke_node ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove tke nodes end ######"
}


main(){
  tke_node_init
  $CALL_FUN || help
}
main
