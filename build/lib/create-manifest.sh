#!/bin/bash

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
    variant="--variant unknown"
  fi

  docker manifest create --amend ${DES_REGISTRY}:${VERSION} \
    ${DES_REGISTRY}-${arch}:${VERSION}

  docker manifest annotate ${DES_REGISTRY}:${VERSION} \
		${DES_REGISTRY}-${arch}:${VERSION} \
		--os ${os} --arch ${arch} ${variant}
done
docker manifest push --purge ${DES_REGISTRY}:${VERSION}
