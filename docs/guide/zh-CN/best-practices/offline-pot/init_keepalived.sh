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

INSTALL_DATA_DIR=/opt/tke-installer/data/

# init keepalived, is tmp script; when tkestack keepalived issue fix will be remove
init_keepalived(){

  for i in {1..100000};
  do
    if [ `cat ${INSTALL_DATA_DIR}/tke.log | grep 'EnsureKubeadmInitKubeConfigPhase' | grep -i 'Success' | wc -l` -gt 0 ]; then
      sh ./offline-pot-cmd.sh -s init-keepalived.sh -f init
      exit 0
    fi
  done
}


main(){
  init_keepalived
}
main
