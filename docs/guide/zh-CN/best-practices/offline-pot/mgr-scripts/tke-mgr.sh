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

CALL_FUN="add_tke_nodes"
hosts="all"

help(){
  echo "show usage:"
  echo "init: init tke deploy config, -f default value is init."
  echo "add_tke_nodes: add tke nodes, -f default value is add_tke_nodes."
  echo "remove_tke_nodes: remove tke nodes"
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

# init tke config
init_tke_cfg(){
  if [ -f "./init-tke-config.sh" ]; then
    sh ./init-tke-config.sh -f init
  else
    echo "init-tke-config.sh not exist, please check !!!" && exit 0
  fi
}


# add tke nodes
add_tke_nodes(){
  if [ -f "./tke-nodes-mgr.sh" ]; then
    sh ./tke-nodes-mgr.sh -f add_tke_nodes
  else
    echo "tke-nodes-mgr.sh not exist, please check !!!" && exit 0
  fi
}

# remove tke nodes
remove_tke_nodes(){
  if [ -f "./tke-nodes-mgr.sh" ]; then
    sh ./tke-nodes-mgr.sh -f remove_tke_nodes
  else
    echo "tke-nodes-mgr.sh not exist, please check !!!" && exit 0
  fi
}

# adjust the number of tke-gateway replicas
tke_gateway(){
  if [ -f "./tke-gateway-mgr.sh" ]; then
    sh ./tke-gateway-mgr.sh -f tke_gateway
  else
    echo "tke-gateway-mgr.sh not exist, please check !!!" && exit 0
  fi
}


main(){
  $CALL_FUN || help
}
main
