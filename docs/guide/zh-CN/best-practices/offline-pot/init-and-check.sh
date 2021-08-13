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

help(){
  echo "show usage:"
  echo "installer_selinux_init: installer selinux disable"
  echo "install_tkeinstaller: install tke-installer"
  echo "host_init: host init"
  echo "hosts_check: hosts check"
  echo "all_func: execute all function, -f default value is all_func !!!"
  exit 0
}

while getopts ":f:h:" opt
do
  case $opt in
    f)
    CALL_FUN="${OPTARG}"
    ;;
    ?)
    echo "unkown args! just suport -f[call function] arg!!!"
    exit 0;;
  esac
done

# installer selinux init
installer_selinux_init(){
  echo "###### installer selinux disable start ######"
  if [ -f "/etc/selinux/config" ]; then
    sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
    setenforce 0 || echo "is disabled !!!"
  fi
  echo "###### installer selinux disable end ######"
}

# install tke-installer
install_tkeinstaller(){
  echo "###### install tke-installer start ######"
  if [ -f './install-tke-installer.sh' ]; then
    sh ./install-tke-installer.sh
  fi
  echo "###### install tke-installer end ######"
}

# host init
host_init(){
  if [ ! -d "/opt/tke-installer/data" ]; then
    mkdir -p /opt/tke-installer/data
  fi
  if [ -d "../tkestack" ]; then
    if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
      if [ -f 'offline-pot-cmd.sh' ]; then
        # whether issue docker crt
        if [ `ls ../offline-pot-tgz/${remote_img_registry_url}.cert.tar.gz | wc -l` -eq 0 ]; then
          sed -i 's/issue_docker_crt/unissue_docker_crt/g' mgr-scripts/host-init.sh
        else
          sed -i 's/unissue_docker_crt/issue_docker_crt/g' mgr-scripts/host-init.sh
        fi
        # host init
        sh ./offline-pot-cmd.sh -s host-init.sh 2>&1 > /opt/tke-installer/data/host-init.log
      fi
    else
      echo "please exec install-tke-installer.sh to start tke-installer" && exit 0
    fi
  else
    if [ -f "./mgr-scripts/host-init.sh" ]; then
      # when not deploy tkestack, will be just deploy helm, mysql need's tools
      sh ./mgr-scripts/host-init.sh -f add_domains 2>&1 > /opt/tke-installer/data/host-init.log
      sh ./mgr-scripts/host-init.sh -f dpl_yum_repo 2>&1 >> /opt/tke-installer/data/host-init.log
      sh ./mgr-scripts/host-init.sh -f install_base_tools 2>&1 >> /opt/tke-installer/data/host-init.log
      if [ -f 'hosts' ]; then
        remote_img_registry_url=`cat hosts | grep ^remote_img_registry_url | awk -F\' '{print $2}'`
      else
        echo "hosts file not exist, please check!!!" && exit 0
      fi
      if [ `ls ../offline-pot-tgz/${remote_img_registry_url}.cert.tar.gz | wc -l` -gt 0 ]; then
        sh ./mgr-scripts/host-init.sh -f issue_img_crt 2>&1 >> /opt/tke-installer/data/host-init.log
      fi
    fi
  fi
}

# hosts check
hosts_check(){
  if [ ! -d "/opt/tke-installer/data" ]; then
    mkdir -p /opt/tke-installer/data
  fi
  if [ -d "../tkestack" ]; then
    if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
      if [ -f 'offline-pot-cmd.sh' ]; then
        sh ./offline-pot-cmd.sh -s hosts-check.sh 2>&1 > /opt/tke-installer/data/hosts-check.log
      fi
    else
      echo "please exec install-tke-installer.sh to start tke-installer" && exit 0
    fi
  else
    if [ -f "./mgr-scripts/hosts-check.sh" ]; then
      sh ./mgr-scripts/hosts-check.sh 2>&1 > /opt/tke-installer/data/hosts-check.log
    fi
  fi
}

all_func(){
  installer_selinux_init
  if [ -d "../tkestack" ]; then
    install_tkeinstaller
  fi
  host_init
  if [ -d "../tkestack" ]; then
    hosts_check
  fi
}

main(){
  $CALL_FUN || help
}
main
