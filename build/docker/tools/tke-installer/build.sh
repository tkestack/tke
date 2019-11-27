#!/usr/bin/env bash

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

REGISTRY_PREFIX=tkestack
VERSION=$(git describe --dirty --always --tags | sed 's/-/./g')
PROVIDER_RES_VERSION=v1.14.6-1
K8S_VERION=${PROVIDER_RES_VERSION%-*}
OUTPUT_DIR=_output/
DST_DIR=$(mktemp -d)
#DST_DIR="/var/folders/20/n4jpmhjs0hd9hjxg80yr2nww0000gn/T/tmp.uN71o6ID"
echo "$DST_DIR" || exit
SCRIPT_DIR=$(dirname "${BASH_SOURCE[0]}")

function usage() {
  cat <<EOF
Usage: ${0} Release TKE
  -h, --help  help
  -q, --quick quick release
EOF
}

function prepare_baremetal_provider() {
  mkdir -p "$DST_DIR"/provider/baremetal/

  cp -rv pkg/platform/provider/baremetal/conf "$DST_DIR"/provider/baremetal
  ls -l "$DST_DIR"/provider/baremetal

  id=$(docker create $REGISTRY_PREFIX/provider-res:"$PROVIDER_RES_VERSION")
  docker cp "$id":/data/res "$DST_DIR"/provider/baremetal/
  docker rm "$id"
}

function prepare_tke_installer() {
  mkdir -p "$DST_DIR"/{bin,conf,data,hooks}

  ls -l "$DST_DIR"

  curl -L https://storage.googleapis.com/kubernetes-release/release/"$K8S_VERION"/bin/linux/amd64/kubectl -o "$DST_DIR"/bin/kubectl
  chmod +x "$DST_DIR"/bin/kubectl

  cp -v "$OUTPUT_DIR"/linux/amd64/tke-installer "$DST_DIR"/bin
  cp -rv cmd/tke-installer/app/installer/manifests "$DST_DIR"
  cp -rv cmd/tke-installer/app/installer/hooks "$DST_DIR"
}

function build_installer_image() {
    docker build --pull -t "$REGISTRY_PREFIX"/tke-installer:"$VERSION" -f "$SCRIPT_DIR"/Dockerfile "$DST_DIR"
    docker push "$REGISTRY_PREFIX"/tke-installer:"$VERSION"
}

function prepare_images() {
  make push VERSION="$VERSION"

  GENERATE_IMAGES_BIN="$OUTPUT_DIR"/$(go env GOOS)/$(go env GOARCH)/tke-generate-images
  make build BINS=tke-generate-images VERSION="$VERSION"

  $GENERATE_IMAGES_BIN
  $GENERATE_IMAGES_BIN | sed "s;^;$REGISTRY_PREFIX/;" | xargs -n1 -I{} sh -c "docker pull {} || exit 1"
  $GENERATE_IMAGES_BIN | sed "s;^;$REGISTRY_PREFIX/;" | xargs docker save | gzip -c >"$DST_DIR"/images.tar.gz
}

pwd

quick=false
while getopts "hq" o; do
    case "${o}" in
        h)
          usage
          ;;
        q)
          quick=true
          ;;
        *)
          usage
          ;;
    esac
done
shift $((OPTIND-1))

make build BINS="tke-installer" OS=linux ARCH=amd64 VERSION="$VERSION"

prepare_baremetal_provider
prepare_tke_installer
if [[ "${quick}" == "false" ]]; then
  prepare_images
fi

build_installer_image
