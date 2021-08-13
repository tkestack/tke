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
  echo "remove_business: remove business"
  echo "remove_base_component: remove base component"
  echo "clean_cluster_nodes: clean cluster nodes, will be delete /data/*, pleace note backup!!!"
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

# remove business
remove_business(){
  echo "###### remove business start ######"
  if [ -d "../tkestack" ]; then
    if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
      if [ -f 'offline-pot-cmd.sh' ]; then
        sh ./offline-pot-cmd.sh -s remove-business.sh -f del_business
      fi
    else
      echo "please exec install-tke-installer.sh to start tke-installer" && exit 0
    fi
  else
    if [ -f "./mgr-scripts/remove-business.sh" ]; then
      sh ./mgr-scripts/remove-business.sh -f del_business
    fi
  fi
  echo "###### remove business end ######"
}

# remove base component
remove_base_component(){
  echo "###### remove base component start ######"
  if [ -d "../tkestack" ]; then
    if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
      if [ -f 'offline-pot-cmd.sh' ]; then
        sh ./offline-pot-cmd.sh -s remove-base-component.sh
      fi
    else
      echo "please exec install-tke-installer.sh to start tke-installer" && exit 0
    fi
  else
    if [ -f "./mgr-scripts/remove-base-component.sh" ]; then
      sh ./mgr-scripts/remove-base-component.sh
    fi
  fi
  echo "###### remove base component end ######"
}

# clean cluster nodes
clean_cluster_nodes(){
  echo "###### clean cluster nodes start ######"
  if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
    if [ -f 'offline-pot-cmd.sh' ] && [ -d "../tkestack" ]; then
      sh ./offline-pot-cmd.sh -s clean-cluster-nodes.sh -f clean_cluster
    fi
  else
    echo "please exec install-tke-installer.sh to start tke-installer" && exit 0
  fi
  echo "###### clean cluster nodes end ######"
}

# execute all function
all_func(){
  remove_business
  remove_base_component
  clean_cluster_nodes
}

main(){
  $CALL_FUN || help
}
main
