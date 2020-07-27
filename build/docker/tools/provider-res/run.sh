#!/usr/bin/env bash

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