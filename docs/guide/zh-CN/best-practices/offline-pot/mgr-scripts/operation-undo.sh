#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="all_func"
hosts="all"

help(){
  echo "show usage:"
  echo "sshd_init_undo: ssh config init undo"
  echo "remove_devnet_proxy_undo: remove devnet proxy undo"
  echo "recover_disk_to_raw: recover disk to raw,default will be set data disk to raw!!"
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

# sshd init undo
sshd_init_undo(){
  echo "###### ssh init undo start ######"
  ansible-playbook -f 10 -i ../hosts --tags sshd_config_recover ../playbooks/operation-undo/operation-undo.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### ssh init undo end ######"
}

# remove devnetcloud proxy undo
remove_devnet_proxy_undo(){
  echo "###### remove devnetcloud proxy undo start ######"
  ansible-playbook -f 10 -i ../hosts --tags recover_docker_proxy ../playbooks/operation-undo/operation-undo.yml \
  --extra-vars "hosts=${hosts}"
  ansible-playbook -f 10 -i ../hosts --tags recover_devnet_proxy ../playbooks/operation-undo/operation-undo.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove devnetcloud proxy undo end ######"
}

# recover disk to raw
recover_disk_to_raw(){
  echo "###### recover disk to raw start ######"
  ansible-playbook -f 10 -i ../hosts --tags recover_disk_to_raw ../playbooks/operation-undo/operation-undo.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### recover disk to raw end ######"
}

# execute all function
all_func(){
  sshd_init_undo
  remove_devnet_proxy_undo
  recover_disk_to_raw
}

main(){
  $CALL_FUN || help
}
main
