#!/usr/bin/env bash

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

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

RPM_DIR=$(mktemp -d)
INSTALL_DIR=$(mktemp -d)

yum install yum-utils -y
yumdownloader --resolve --destdir=${RPM_DIR} ${PKG}
VERSION=$(echo /${RPM_DIR}/${PKG}* | grep -oP '\d+\.\d+\.\d+')

cd ${INSTALL_DIR}
for rpm in ${RPM_DIR}/*.rpm; do rpm2cpio $rpm | cpio -idm; done

tar -cvzf /output/${PKG}-${OS}-${ARCH}-${VERSION}.tar.gz *
