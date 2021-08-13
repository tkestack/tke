#!/bin/bash

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

action=$1
app=$2
env=$3
tag=$4

if [ ! -n "$1" ] || [ ! -n "$2" ] || [ ! -n "$3" ]; then
    echo "Usage: ./tools/helm.sh template|upgrade|install product2-demo ty 20190929101615"
    exit
fi

release=${app}-${env}

# get current deployment image tag as tag var
if [ ! -n "$tag" ]; then
    tag="`kubectl get deploy ${release} -o=jsonpath='{.spec.template.spec.containers[0].image}' | awk -F':' '{print $2}'`"
    echo "get current tag($tag) as default tag"
fi

if [ ! -n "$tag" ]; then
    echo 'failed to get current deployment tag, please input it manually'
    echo "Usage: ./tools/helm.sh template|upgrade|install product2-demo ty 20190929101615"
    exit 1
fi

echo "helm dep up ${app}"
#echo "helm -n default ${action} -f ${app}/values/${env}.yaml --set bootstrap.image.tag=${tag} ${release} ${app}"
echo "helm -n default ${action} -f ../helmfile.d/config/product2-${env}/${release}.yaml --set bootstrap.image.tag=${tag} ${release} ${app}"


helm dep up ${app}
#helm -n default ${action} -f ${app}/values/${env}.yaml --set bootstrap.image.tag=${tag} ${release} ${app}
#helmfile --log-level=debug --namespace default -e ${env} -f ../helmfile.d/releases/${release}.yaml template --args "--debug --dry-run"  2>&1| grep  'exec: helm template'|grep 'dry-run:' |awk '{print "mv", $9, "config/"releases"/"$4".yaml"}' releases="$releases" |sh
helm -n default ${action} -f ../helmfile.d/config/product2-${env}/${release}.yaml --set bootstrap.image.tag=${tag} ${release} ${app}



### Helm install example
#helm install product2-web-job1 product2-web -f product2-web/values/job1.yaml

