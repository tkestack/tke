#!/usr/bin/env bash

export DIR="cmd/tke-upgrade"
export DATADIR="${DIR}/app"
export MANIFESTS="cmd/tke-installer/app/installer/manifests"
export OUTPUT="_output/upgrade/deploy"
export VERSION="v1.2.3"

go run ${DIR}/main.go

cat ${OUTPUT}/* | grep -v 'image: tkestack/' | grep 'image:' | tr '/' ' ' | awk '{print "tkestack/"$4}' | xargs -L1 docker pull
cat ${OUTPUT}/* | grep -v 'image: tkestack/' | grep 'image:' | tr '/' ' ' | awk '{print "tkestack/"$4" "$2"/"$3"/"$4}' | xargs -L1 docker tag
cat ${OUTPUT}/* | grep -v 'image: tkestack/' | grep 'image:' | tr '/' ' ' | awk '{print $2"/"$3"/"$4}' | xargs -L1 docker push
