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

REGISTRY_PREFIX=${REGISTRY_PREFIX:-"tkestack"}
PLATFORMS=${PLATFORMS:-"linux_amd64 linux_arm64"}

if [ -z ${IMAGE} ]; then
  echo "Please provide IMAGE."
  exit 1
fi

if [ -z ${VERSION} ]; then
  echo "Please provide VERSION."
  exit 1
fi

rm -rf ${HOME}/.docker/manifests/docker.io_${REGISTRY_PREFIX}_${IMAGE}-${VERSION}
DES_REGISTRY=${REGISTRY_PREFIX}/${IMAGE}
for platform in ${PLATFORMS}; do
  os=${platform%_*}
  arch=${platform#*_}
  variant=""
  if [ ${arch} == "arm64" ]; then
    variant="--variant v8"
  fi

  docker manifest create --amend ${DES_REGISTRY}:${VERSION} \
    ${DES_REGISTRY}-${arch}:${VERSION}

  docker manifest annotate ${DES_REGISTRY}:${VERSION} \
		${DES_REGISTRY}-${arch}:${VERSION} \
		--os ${os} --arch ${arch} ${variant}
done
docker manifest push --purge ${DES_REGISTRY}:${VERSION}
