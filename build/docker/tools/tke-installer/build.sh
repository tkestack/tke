#! /usr/bin/env bash

# Tencent is pleased to support the open source community by making TKEStack
# available.
#
# Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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


set -o errexit
set -o nounset
set -o pipefail

umask 0022
unset IFS
unset OFS
unset LD_PRELOAD
unset LD_LIBRARY_PATH

export PATH='/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'

target=${1:-tke-installer}

cd $(dirname $0)

echo "===> Begin to create $target in $(pwd)"

echo "Step.1 cleanup"
rm -f $target

echo "Step.2 prepare package"
chmod +x *.sh
cp -f init_installer.sh $target
tar -czf package.tgz install.sh res

echo "Step.3 generate installer"
cat package.tgz >>$target
chmod +x $target
rm -f package.tgz

echo "===> Success to create $target"
