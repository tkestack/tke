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
OPTIONS="--name tke-installer -d --privileged --net=host
-v /etc/hosts:/app/hosts
-v /etc/docker:/etc/docker
-v /var/run/docker.sock:/var/run/docker.sock
-v $DATA_DIR:/app/data
-v $INSTALL_DIR/conf:/app/conf
-v registry-certs:/app/certs
-v tke-installer-bin:/app/bin
"

function prefight() {
  echo "Step.1 prefight"

  if [ "root" != "$(whoami)" ]; then
    echo "only root can execute this script"
    exit 1
  fi
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
  for i in {1..60}; do
    if systemctl is-active --quiet docker; then
      break
    fi
    sleep 1
  done
  if (( i == 10)); then
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

  rm -rf $DATA_DIR || :
  docker rm -f tke-installer || :
  docker volume prune -f || :

  echo "Step.4 clean old data [ok]"
}

function start_installer() {
  echo "Step.5 start tke-installer [doing]"

  docker run $OPTIONS tkestack/tke-installer:$VERSION

  echo "Step.5 start tke-installer [ok]"
}

function check_installer() {
  echo "Step.6 check tke-installer status [doing]"

  ip=$(ip route get 1 | awk '{print $NF;exit}')
  url="http://$ip:8080/index.html"
  if ! curl -sSf $url  >/dev/null; then
    echo "check installer status error"
    docker logs tke-installer
    exit 1
  fi

  echo "Step.6 check tke-installer status [ok]"

  echo "Please use your browser which can connect this machine to open $url for install TKE!"
}

prefight
ensure_docker
load_image
clean_old_data
start_installer
check_installer
