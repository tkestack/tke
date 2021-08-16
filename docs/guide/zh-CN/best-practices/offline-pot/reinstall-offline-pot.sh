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

# clean cluster
clean_cluster(){
  sh ./clean-cluster.sh
}

# reinstall tke-installer
redpl_tke_installer(){
  if [ -d "../tkestack" ]; then
    sh ./install-tke-installer.sh
  fi
}

# init tke config and deploy offline-pot
dpl_offline_pot(){
  sh ./install-offline-pot.sh
}


main(){
  clean_cluster
  redpl_tke_installer
  dpl_offline_pot
}
main
