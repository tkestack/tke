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

# tkestack version
version='v1.2.4'

# business helm's or workerload kubernetes's yaml git url
busi_git_url='https://git.xx.com/yy/helm.git'

# business helm's or workerload kubernetes's yaml git branch will be check out
busi_branch='feature/private'

# all servers for set deploy switch
all_servers=("tkestack" "business" "redis" "redis_cluster" "mysql" "prometheus" "kafka" "elk" "nginx_ingress" "minio" "helmtiller" "nfs" "salt_minion" "postgres" "sgikes")

# will be deploy's server set
server_set=("tkestack" "business" "redis" "mysql" "prometheus" "kafka" "elk" "nginx_ingress" "helmtiller" "nfs" "salt_minion" "postgres")

# whether use remote docker registry, if true will be not save business images and copy registry secret; true or false
remote_registry='true'
# remote image registry url, if remote images registry need issue crt, please name: ${remote_img_registry_url}.cert.tar.gz
# on offline-pot-tgz-base dir
remote_img_registry_url='reg.xx.yy.com'

# builder type just support 'all' or 'custom' , default is all; customize will be pack on demand
builder_type='custom'
