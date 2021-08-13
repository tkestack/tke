#!/bin/bash

# Tencent is pleased to support the open source community by making TKEStack
# available.
#
# Copyright (C) 2012-2021 Tencent. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use
# this file except in compliance with the License. You may obtain a copy of the
# License at
#
# https://opensource.org/licenses/Apache-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OF ANY KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations under the License.

# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="all_func"
hosts="all"

help(){
  echo "show usage:"
  echo "sshd_init: ssh config init"
  echo "selinux_init: disable selinux"
  echo "remove_devnet_proxy: remove devnet proxy"
  echo "add_domains: add domains"
  echo "dpl_yum_repo: deploy offline yum repo"
  echo "install_base_tools: install base tools,example: helm tools install"
  echo "issue_img_crt: issue iamges registry crt"
  echo "data_disk_init: data disk init"
  echo "check_iptables: check iptables"
  echo "time_sync: check time sync service whether deploy"
  echo "install_stress_tools: install stress test tools"
  echo "registry_influxdb_init: create local registry and influxdb data director soft link"
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

# sshd init
sshd_init(){
  echo "###### ssh init start ######"
  ansible-playbook -f 10 -i ../hosts --tags sshd_config_init ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### ssh init end ######"
}

# selinux init
selinux_init(){
  echo "###### selinux init start ######"
  ansible-playbook -f 10 -i ../hosts --tags selinux_init ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### selinux init end ######"
}

# remove devnetcloud proxy
remove_devnet_proxy(){
  echo "###### remove devnetcloud proxy start ######"
  ansible-playbook -f 10 -i ../hosts --tags remove_devnetcloud_docker_proxy ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  ansible-playbook -f 10 -i ../hosts --tags remove_devnet_proxy ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove devnetcloud proxy end ######"
}

# add domains
add_domains(){
  echo "###### add domains start ######"
  ansible-playbook -f 10 -i ../hosts --tags add_domains ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### add domains end ######"
}

# deploy offlie yum repo
dpl_yum_repo(){
  echo "###### deploy offlie yum repo start ######"
  ansible-playbook -f 10 -i ../hosts --tags deploy_offline_yum_repo ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy offlie yum repo end ######"
}

# install base tools
install_base_tools(){
  echo "###### install base tools start ######"
  ansible-playbook -f 10 -i ../hosts --tags install_base_tools ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### install base tools end ######"
}

# issue images registry crt
issue_img_crt(){
  echo "###### issue images registry crt start ######"
  ansible-playbook -f 10 -i ../hosts --tags issue_docker_crt ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### issue images registry crt end ######"
}

# data disk init
data_disk_init(){
  echo "###### data disk init start ######"
  ansible-playbook -f 10 -i ../hosts --tags data_disk_init ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### data disk init end ######"
}

# check iptables
check_iptables(){
  echo "###### check iptables start ######"
  ansible-playbook -f 10 -i ../hosts --tags check_iptables ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### check iptables  end ######"
}

# check time sync service
time_sync(){
  echo "###### check time sync service start ######"
  ansible-playbook -f 10 -i ../hosts --tags check_time_syn ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### check time sync service end ######"
}

# install stress tools
install_stress_tools(){
  echo "###### check time sync service start ######"
  ansible-playbook -f 10 -i ../hosts --tags install_stress_tools ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### check time sync service end ######"
}

# create local registry and influxdb data director soft link
registry_influxdb_init(){
  echo "###### create local registry and influxdb data director soft link start ######"
  ansible-playbook -f 10 -i ../hosts --tags registry_influxdb_init ../playbooks/hosts-init/hosts-init.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### create local registry and influxdb data director soft link end ######"
}

# execute all function
all_func(){
  sshd_init
  selinux_init
  remove_devnet_proxy
  add_domains
  dpl_yum_repo
  data_disk_init
  check_iptables
  time_sync
  install_stress_tools
  install_base_tools
  issue_img_crt
  registry_influxdb_init
}

main(){
  $CALL_FUN || help
}
main
