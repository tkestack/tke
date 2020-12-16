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

VERSION=latest

INSTALL_DIR=/opt/tke-installer
DATA_DIR=$INSTALL_DIR/data
OPTIONS="--name tke-installer -d --privileged --net=host --restart=always
-v /etc/hosts:/app/hosts
-v /etc/docker:/etc/docker
-v /var/run/docker.sock:/var/run/docker.sock
-v $DATA_DIR:/app/data
-v $INSTALL_DIR/conf:/app/conf
-v registry-certs:/app/certs
-v tke-installer-bin:/app/bin
"

declare -A archMap=(
  [x86_64]=amd64
  [aarch64]=arm64
)

if [ -z ${archMap[$(arch)]+unset} ]; then
  echo "ERROR: unsupport arch $(arch)"
fi
ARCH=${archMap[$(arch)]}

function prefight() {
  echo "Step.1 prefight"

  check::root
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

function ensure_docker() {
  echo "Step.2 ensure docker is ok"

  if ! [ -x "$(command -v docker)" ]; then
    echo "command docker not find"
    install_docker
  fi
  if ! systemctl is-active --quiet docker; then
    echo "docker status is not running"
    install_docker
  fi
}

function install_docker() {
  echo "install docker [doing]"

  tar xvaf "res/docker.tgz" -C /usr/bin --strip-components=1
  cp -v res/docker.service /etc/systemd/system
  mkdir -p /etc/docker
  cp -v res/daemon.json /etc/docker/

  systemctl daemon-reload

  # becuase first start docker may be restart some times
  systemctl start docker || :
  maxSecond=60
  for i in $(seq 1 $maxSecond); do
    if systemctl is-active --quiet docker; then
      break
    fi
    sleep 1
  done
  if ((i == maxSecond)); then
    echo "start docker failed, please check docker service."
    exit 1
  fi

  echo "install docker [ok]"
}

function load_image() {
  echo "Step.3 load tke-installer image [doing]"

  docker load -i res/tke-installer.tgz

  echo "Step.3 load tke-installer image [ok]"
}

function clean_old_data() {
  echo "Step.4 clean old data [doing]"

  rm -rf $DATA_DIR >/dev/null 2>&1 || :
  docker rm -f tke-installer >/dev/null 2>&1 || :
  docker volume prune -f >/dev/null 2>&1 || :

  echo "Step.4 clean old data [ok]"
}

function start_installer() {
  echo "Step.5 start tke-installer [doing]"

  docker run $OPTIONS "tkestack/tke-installer-${ARCH}:$VERSION" $@

  echo "Step.5 start tke-installer [ok]"
}


function check_installer() {
  s=10
  for i in $(seq 1 $s)
  do
    echo "Step.6 check tke-installer status [doing]"
    url="http://127.0.0.1:8080/index.html"
    if ! curl -sSf "$url" >/dev/null; then
      sleep 3
      echo "Step.6 retries left $(($s-$i))"
      continue
    else
      echo "Step.6 check tke-installer status [ok]"
      echo "Please use your browser which can connect this machine to open $url for install TKE!"
      exit 0
    fi
  done
  echo "check installer status error"
  docker logs tke-installer
  exit 1 
}

prefight
ensure_docker
load_image
clean_old_data
start_installer $@
check_installer
