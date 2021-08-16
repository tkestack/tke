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

CALL_FUN="defaut"

help(){
  echo "show usage:"
  echo "init_and_check: will be init hosts, inistall tke-installer and hosts check"
  echo "dpl_offline_pot: init tke config and deploy offline-pot"
  echo "init_keepalived: just tmp use, when tkestack fix keepalived issue will be remove"
  echo "only_install_tkestack: if you want only install tkestack, please -f parameter pass only_install_tkestack"
  echo "defualt: will be exec dpl_offline_pot and init_keepalived"
  echo "all_func: execute init_and_check, dpl_offline_pot, init_keepalived"
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

INSTALL_DATA_DIR=/opt/tke-installer/data/

init_and_check(){
  sh ./init-and-check.sh
}

# init tke config and deploy offline-pot
dpl_offline_pot(){
  echo "###### deploy offline-pot start ######"
  if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
    # deploy tkestack , base commons and business
    sh ./offline-pot-cmd.sh -s init-tke-config.sh -f init
    docker restart tke-installer
    if [ -f "hosts" ]; then
      installer_ip=`cat hosts | grep -A 1 '\[installer\]' | grep -v installer`
      echo "please exec tail -f ${INSTALL_DATA_DIR}/tke.log or access http://${installer_ip}:8080 check install progress..."
    fi
  elif [ ! -d "../tkestack" ]; then
    # deploy base commons and business on other kubernetes plat
    sh ./post-install
  else
    echo "if first install,please exec init-and-check.sh script, else exec reinstall-offline-pot.sh script" && exit 0
  fi
  echo "###### deploy offline-pot end ######"
}

# just tmp use, when tkestack fix keepalived issue will be remove
init_keepalived(){
  echo "###### init keepalived start  ######"
  if [ -f "${INSTALL_DATA_DIR}/tke.json" ]; then
    if [ `cat ${INSTALL_DATA_DIR}/tke.json | grep -i '"ha"' | wc -l` -gt 0 ]; then
      nohup sh ./init_keepalived.sh 2>&1 > ${INSTALL_DATA_DIR}/dpl-keepalived.log &
    fi
  fi
  echo "###### init keepalived end ######"
}

# only install tkestack
only_install_tkestack(){
  echo "###### install tkestack start ######"
  # change tke components's replicas number
  if [ -f "hosts" ]; then
    sed -i 's/tke_replicas="1"/tke_replicas="2"/g' hosts
  fi
  # hosts init
  if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
    sh ./offline-pot-cmd.sh -s host-init.sh -f sshd_init
    sh ./offline-pot-cmd.sh -s host-init.sh -f selinux_init
    sh ./offline-pot-cmd.sh -s host-init.sh -f remove_devnet_proxy
    sh ./offline-pot-cmd.sh -s host-init.sh -f add_domains
    sh ./offline-pot-cmd.sh -s host-init.sh -f data_disk_init
    sh ./offline-pot-cmd.sh -s host-init.sh -f check_iptables
  else
    echo "please exec install-tke-installer.sh to start tke-installer" && exit 0
  fi
  # start install tkestack
  dpl_offline_pot
  init_keepalived
  echo "###### install tkestack end ######"
}

defaut(){
  # change tke components's replicas number
  if [ -f "hosts" ]; then
    sed -i 's/tke_replicas="2"/tke_replicas="1"/g' hosts
  fi
  # only deploy tkestack
  if [ -d '../tkestack' ] && [ ! -d "../offline-pot-images" ] && [ ! -d "../offline-pot-tgz" ]; then
    only_install_tkestack
  fi
  dpl_offline_pot
  # when deploy tkestack will be init keepalived config
  # if [ -d '../tkestack' ]; then
  #  init_keepalived
  # fi
}

all_func(){
  # change tke components's replicas number
  if [ -f "hosts" ]; then
    sed -i 's/tke_replicas="2"/tke_replicas="1"/g' hosts
  fi
  init_and_check
  defaut
}

main(){
  $CALL_FUN || help
}
main
