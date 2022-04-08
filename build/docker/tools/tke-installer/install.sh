#! /usr/bin/env bash

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

umask 0022
unset IFS
unset OFS
unset LD_PRELOAD
unset LD_LIBRARY_PATH

export PATH='/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'

VERSION=latest

INSTALL_DIR=/opt/tke-installer
DATA_DIR=$INSTALL_DIR/data
REGISTRY_DIR=$INSTALL_DIR/registry
REGISTRY_VERSION=2.7.1
OPTIONS="--name tke-installer -d --privileged --net=host --restart=always
-v /etc/hosts:/app/hosts
-v $DATA_DIR:/app/data
-v /var/run/containerd/containerd.sock:/var/run/containerd/containerd.sock
-v /run/containerd/:/run/containerd/
-v $INSTALL_DIR/conf:/app/conf
-v registry-certs:/app/certs
-v tke-installer-bin:/app/bin
-v /tmp:/tmp 
-v /lib/modules/:/lib/modules/
"

RegistryHTTPOptions="--name registry-http -d --net=host --restart=always -p 80:5000
-e REGISTRY_HTTP_ADDR=0.0.0.0:80 \
-v $REGISTRY_DIR:/var/lib/registry
"
RegistryHTTPSOptions="--name registry-https -d --net=host --restart=always -p 443:443
-v $REGISTRY_DIR:/var/lib/registry
-v registry-certs:/certs
-e REGISTRY_HTTP_ADDR=0.0.0.0:443
-e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/server.crt
-e REGISTRY_HTTP_TLS_KEY=/certs/server.key
"

declare -A archMap=(
  [x86_64]=amd64
  [aarch64]=arm64
)

if [ -z ${archMap[$(arch)]+unset} ]; then
  echo "ERROR: unsupport arch $(arch)"
fi
ARCH=${archMap[$(arch)]}

function preflight() {
  echo "Step.1 preflight"

  check::root
  check::docker
  check::disk '/opt' 30
  check::disk '/var/lib' 20
}

function check::root() {
  if [ "root" != "$(whoami)" ]; then
    echo "only root can execute this script"
    exit 1
  fi
  echo "root: yes"
}

function check::disk() {
  local -r path=$1
  local -r size=$2

    disk_avail=$(df -BG "$path" | tail -1 | awk '{print $4}' | grep -oP '\d+')
  if ((disk_avail < size)); then
    echo "available disk space for $path needs be greater than $size GiB"
    exit 1
  fi

  echo "available disk space($path):  $disk_avail GiB"
}

function check::docker() {
  echo "check docker status"

  if systemctl is-active --quiet docker; then
    echo "docker is running,please stop docker before using containerd runtime"
    echo "if the migration case please follow the migration doc:"
    echo "https://tkestack.github.io/web/zh/blog/2021/09/01/container-runtime-migraion/"
    exit 1
  fi
}

function ensure_containerd() {
  echo "Step.2 check containerd status"

  if ! [ -x "$(command -v nerdctl)" ]; then
    echo "command nerdctl not find"
    install_containerd
  fi
  if ! systemctl is-active --quiet containerd; then
    echo "containerd status is not running"
    install_containerd
  fi
}

function install_containerd() {
  echo "install containerd [in process]"

  # Install containerd exclude cni binaries and cni config file.
  tar xvaf "res/containerd.tar.gz" -C / --exclude=etc/cni --exclude=opt
  tar xvaf "res/nerdctl.tar.gz" -C /usr/local/bin/

  systemctl daemon-reload

  # becuase first start docker may be restart some times
  systemctl start containerd || :
  maxSecond=60
  for i in $(seq 1 $maxSecond); do
    if systemctl is-active --quiet containerd; then
      break
    fi
    sleep 1
  done
  if ((i == maxSecond)); then
    echo "start containerd failed, please check containerd service."
    exit 1
  fi

  echo "install containerd [done]"
}

function load_image() {
  echo "Step.3 load tke-installer image [in process]"

  nerdctl load -i res/tke-installer.tar
  nerdctl load -i res/registry.tar

  echo "Step.3 load tke-installer image [done]"
}

function clean_old_data() {
  echo "Step.4 clean old data [in process]"

  nerdctl stop tke-installer >/dev/null 2>&1 && nerdctl rm tke-installer >/dev/null 2>&1 || :
  nerdctl stop registry-http >/dev/null 2>&1 && nerdctl rm registry-http >/dev/null 2>&1 || :
  nerdctl stop registry-https >/dev/null 2>&1 && nerdctl rm registry-https >/dev/null 2>&1 || :
  nerdctl volume rm tke-installer-bin >/dev/null 2>&1 || :
  nerdctl volume rm registry-certs >/dev/null 2>&1 || :

  if  [ -d  "$DATA_DIR" ]; then
    rm -f $DATA_DIR/tke.json >/dev/null 2>&1 || :
    rm -f $DATA_DIR/tke.log >/dev/null 2>&1 || :
  fi

  echo "Step.4 clean old data [done]"
}

function start_installer() {
  echo "Step.5 start tke-installer [in process]"
  mkdir -p $DATA_DIR
  mkdir -p $INSTALL_DIR/conf
  nerdctl run $OPTIONS "tkestack/tke-installer-${ARCH}:$VERSION" $@

  echo "Step.5 start tke-installer [done]"
}

function start_registry() {
  echo "Step.6 start regisry [in process]"
  checkupgrade="--upgrade"
  if [[ $@ =~ $checkupgrade ]]; then
    echo "Step.6 upgrade will not start local registry"
    echo "Step.6 start registry [skip]"
  else
    mkdir -p $REGISTRY_DIR
    nerdctl run $RegistryHTTPOptions "tkestack/registry-${ARCH}:$REGISTRY_VERSION"
    nerdctl run $RegistryHTTPSOptions "tkestack/registry-${ARCH}:$REGISTRY_VERSION"
    echo "Step.6 start registry [done]"
  fi
}

function check_installer() {
  s=10
  for i in $(seq 1 $s)
  do
    echo "Step.6 check tke-installer status [in process]"
    url="http://127.0.0.1:8080/index.html"
    if ! curl -sSf "$url" >/dev/null 2>&1; then
      sleep 3
      echo "Step.6 retries left $(($s-$i))"
      continue
    else
      echo "Step.6 check tke-installer status [done]"
      echo "Please use your browser which can connect this machine to open $url for install TKE!"
      exit 0
    fi
  done
  echo "check installer status error"
  nerdctl logs tke-installer
  exit 1
}

preflight
ensure_containerd
load_image
clean_old_data
start_installer $@
start_registry $@
check_installer
