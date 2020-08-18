#!/bin/bash

action=$1
app=$2
env=$3
tag=$4

if [ ! -n "$1" ] || [ ! -n "$2" ] || [ ! -n "$3" ]; then
    echo "Usage: ./tools/helmfile.sh template|sync product1-app1 demo 20190929101615"
    exit
fi

# patch for qcloud(pro) env
if [ "$env" == 'qcloud' ]; then
    release=${app}
else
    release=${app}-${env}
fi

# get current deployment image tag as tag var
if [ ! -n "$tag" ]; then
    tag="`kubectl get deploy ${release} -o=jsonpath='{.spec.template.spec.containers[0].image}' | awk -F':' '{print $2}'`"
    echo "get current tag($tag) as default tag"
fi

if [ ! -n "$tag" ]; then
    echo 'failed to get current deployment tag, please input it manually'
    echo "Usage: ./tools/helmfile.sh template|sync 20190929101615 ty 20190929101615"
    exit 1
fi

echo "helmfile --log-level=debug --namespace default -e ${env} --selector app=${app} -f ../helmfile.d/releases/product1-${env}.yaml template --args \"--debug --dry-run\"  2>&1| grep  'exec: helm template'|grep 'dry-run:' |awk '{print "mv", $9, "../helmfile.d/config/"releases"/"$4".yaml"}' releases="product1-$env" |sh"


if [ "$action" == 'template' ]; then
    helmfile --log-level=debug --namespace default -e ${env} --selector app=${app} -f ../helmfile.d/releases/product1-${env}.yaml ${action} --args "--debug --dry-run"  2>&1| grep  'exec: helm template'|grep 'dry-run:' |awk '{print "mv", $9, "../helmfile.d/config/"releases"/"$4".yaml"}' releases="product1-$env" |sh
else
    helmfile --namespace default -e ${env} --selector app=${app} -f ../helmfile.d/releases/product1-${env}.yaml ${action}
fi


### Helmfile install example
#helmfile --namespace default -e demo --selector app=${app} -f ../helmfile.d/releases/product1-demo.yaml sync
