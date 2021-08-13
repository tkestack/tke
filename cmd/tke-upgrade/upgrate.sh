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

export DIR="cmd/tke-upgrade"
export DATADIR="${DIR}/app"
export MANIFESTS="cmd/tke-installer/app/installer/manifests"
export OUTPUT="_output/upgrade/deploy"
export VERSION="v1.2.3"

go run ${DIR}/main.go

cat ${OUTPUT}/* | grep -v 'image: tkestack/' | grep 'image:' | tr '/' ' ' | awk '{print "tkestack/"$4}' | xargs -L1 docker pull
cat ${OUTPUT}/* | grep -v 'image: tkestack/' | grep 'image:' | tr '/' ' ' | awk '{print "tkestack/"$4" "$2"/"$3"/"$4}' | xargs -L1 docker tag
cat ${OUTPUT}/* | grep -v 'image: tkestack/' | grep 'image:' | tr '/' ' ' | awk '{print $2"/"$3"/"$4}' | xargs -L1 docker push
