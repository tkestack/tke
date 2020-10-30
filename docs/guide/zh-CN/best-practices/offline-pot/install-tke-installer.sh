#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

# get offline-pot parent dir
OFFLINE_POT_PDIR=`echo ${BASE_DIR} | awk -Foffline-pot '{print $1}'`

INSTALL_DIR=/opt/tke-installer
DATA_DIR=${INSTALL_DIR}/data
HOOKS=${OFFLINE_POT_PDIR}offline-pot
IMAGES_DIR="${OFFLINE_POT_PDIR}offline-pot-images"
TGZ_DIR="${OFFLINE_POT_PDIR}offline-pot-tgz"
REPORTS_DIR="${OFFLINE_POT_PDIR}perfor-reports"
version=v1.2.4
arch=amd64

init_tke_installer(){
  if [ `docker images | grep tke-installer | grep ${version} | wc -l` -eq 0 ]; then
    if [ `docker ps -a | grep tke-installer | wc -l` -gt 0 ]; then
      docker rm -f tke-installer
    fi
    if [ `docker images | grep tke-installer | wc -l` -gt 0 ]; then
      docker rmi -f `docker images | grep tke-installer | awk '{print $3}'`
    fi 
    cd ${OFFLINE_POT_PDIR}tkestack
    if [ `echo "$version" | awk -Fv '{print $2}' | awk -F. '{print $1$2$3}'` -lt 130 ]; then
      if [ -d "${OFFLINE_POT_PDIR}tkestack/tke-installer-x86_64-${version}.run.tmp" ]; then
        rm -rf ${OFFLINE_POT_PDIR}tkestack/tke-installer-x86_64-${version}.run.tmp
      fi
      sha256sum --check --status tke-installer-x86_64-$version.run.sha256 && \
      chmod +x tke-installer-x86_64-$version.run && ./tke-installer-x86_64-$version.run
    else
      if [ -d "${OFFLINE_POT_PDIR}tkestack/tke-installer-linux-$arch-${version}.run.tmp" ]; then
        rm -rf ${OFFLINE_POT_PDIR}tkestack/tke-installer-linux-$arch-${version}.run.tmp
      fi
      sha256sum --check --status tke-installer-linux-$arch-$version.run.sha256 && \
      chmod +x tke-installer-linux-$arch-$version.run && ./tke-installer-linux-$arch-$version.run
    fi
  fi
}

reinstall_tke_installer(){
  if [ -d "${REPORTS_DIR}" ]; then
    mkdir -p ${REPORTS_DIR}
  fi
  if [ `docker ps -a | grep tke-installer | wc -l` -eq 1 ]; then
    docker rm -f tke-installer
    rm -rf /opt/tke-installer/data
  fi
  if [ `echo "$version" | awk -Fv '{print $2}' | awk -F. '{print $1$2$3}'` -lt 130 ]; then
    TKT_INSTALLER_IMAGE="tkestack/tke-installer:${version}"
  else
    TKT_INSTALLER_IMAGE="tkestack/tke-installer-${arch}:${version}"
  fi
  docker run --restart=always --name tke-installer -d --privileged --net=host -v/etc/hosts:/app/hosts \
  -v/etc/docker:/etc/docker -v/var/run/docker.sock:/var/run/docker.sock -v$DATA_DIR:/app/data \
  -v$INSTALL_DIR/conf:/app/conf -v$HOOKS:/app/hooks -v$IMAGES_DIR:${IMAGES_DIR} -v${TGZ_DIR}:${TGZ_DIR} \
  -v${REPORTS_DIR}:${REPORTS_DIR} ${TKT_INSTALLER_IMAGE}
  if [ -f "hosts" ]; then
    # set hosts file's dpl_dir variable
    sed -i 's#^dpl_dir=.*#dpl_dir=\"'"${HOOKS}"'\"#g' hosts
    installer_ip=`cat hosts | grep -A 1 '\[installer\]' | grep -v installer`
    echo "please exec install-offline-pot.sh or access http://${installer_ip}:8080 to install offline-pot"
  fi
}

main(){
  init_tke_installer
  reinstall_tke_installer
}
main
