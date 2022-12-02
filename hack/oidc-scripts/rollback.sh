#!/usr/bin/env bash

function backup() {
  echo "==== 0. Execute OIDC rollback backup doing ===="
  api=$(kubectl -n tke get -o=name deployment | grep api)
  cm=$(kubectl -n tke get -o=name configmap | grep api)
  ds=$(kubectl -n tke get -o=name daemonset | grep gateway)
  gw=$(kubectl -n tke get -o=name configmap | grep gateway)
  r=("$api $cm $ds $gw")
  echo ======= backup current resources =======
  dir=/opt/oidcbackup_beforrollback
  mkdir -p ${dir}/deployment.apps ${dir}/daemonset.apps ${dir}/configmap
  for n in $r; do
    kubectl -n tke get -o=yaml $n >${dir}/$n.yaml
  done
  echo "==== 0. Execute OIDC rollback backup success ===="
}

function rollbackConfigMap() {
  echo ======= rollback configmap $1 =======

  file=$1.yaml

  kubectl -n tke get cm $1 -o yaml >$file

  start=$(sed -n '/last-applied-configuration/=' $file)
  if [ "$start" != "" ]; then
    end=$(($start + 1))
    sed -i "$start, $end d" $file
  fi

  if [ "$file" == "tke-auth-api.yaml" ]; then

    start=$(sed -n '/authentication.oidc/=' $file)
    if [ "$start" != "" ]; then
      end=$(($start + 8))
      sed -i "$start, $end d" $file
    fi

    sed -i '/init_tenant_type/d' $file
    sed -i '/init_tenant_id/d' $file
    sed -i '/cloudindustry_config_file/d' $file
    line=$(sed -n '/init_client_id/=' $file)
    if [ "$line" != "" ]; then
      sed -i "${line}c \ \ \ \ init_client_id = \"default\"" $file
    fi
    line=$(sed -n '/init_client_secret/=' $file)
    if [ "$line" != "" ]; then
      sed -i "${line}c \ \ \ \ init_client_secret = \"2HWJXNnGagpGvnSBQ6Y2P8xJylu\"" $file
    fi
  else
    #    sed -i '/client_secret/d' $file
    line=$(sed -n '/client_secret/=' $file)
    if [ "$line" != "" ]; then
      sed -i "${line}c \ \ \ \ \ \ client_secret = \"2HWJXNnGagpGvnSBQ6Y2P8xJylu\"" $file
    fi
    line=$(sed -n '/client_id/=' $file)
    if [ "$line" != "" ]; then
      sed -i "${line}c \ \ \ \ \ \ client_id = \"default\"" $file
    fi
    line=$(sed -n '/\<issuer_url/=' $file)
    if [ "$line" != "" ]; then
      sed -i "${line}c \ \ \ \ \ \ issuer_url = \"https://tke-auth-api/oidc\"" $file
      cat $file | grep external_issuer_url
      if [ $? != 0 ]; then
        sed -i "/issuer_url/a\      external_issuer_url = \"https://tke-auth-api/oidc\"" $file
      fi
    fi
    line=$(sed -n '/^      ca_file/=' $file)
    if [ "$line" == "" ]; then
      sed -i '/username_prefix/i\      ca_file = "\/app\/certs\/ca.crt"' $file
    else
      sed -i "${line}c \ \ \ \ \ \ ca_file = \"\/app\/certs\/ca.crt\"" $file
    fi
    if [ "$file" == "tke-gateway.yaml" ]; then
      sed -i '/disableOIDCProxy: true/d' $file
      sed -i '/generic/d' $file
      sed -i '/external_hostname/d' $file
    fi
  fi

  sed -i '/resourceVersion/d' $file
  sed -i '/uid/d' $file

  kubectl apply -f $file
  rm $file

}

function rollbackWorkload() {
  resourceType=$1
  resource=$2
  echo ======= rollback $resourceType $resource =======
  kubectl -n tke annotate $1 $2 tkeanywhere/oidc-

  script="
import json; import sys

obj=json.load(sys.stdin)

obj['spec']['template']['spec'].pop('hostAliases', None)

for i, v in enumerate(obj['spec']['template']['spec']['containers'][0]['volumeMounts']):
    if v['name'] == 'oidc-ca-volume' :
      obj['spec']['template']['spec']['containers'][0]['volumeMounts'].pop(i)

for i, v in enumerate(obj['spec']['template']['spec']['containers'][0]['volumeMounts']):
    if v['name'] == 'cloudindustry-config-volume':
      obj['spec']['template']['spec']['containers'][0]['volumeMounts'].pop(i)

for i, v in enumerate(obj['spec']['template']['spec']['volumes']):
    if v['name'] == 'oidc-ca-volume':
      obj['spec']['template']['spec']['volumes'].pop(i)

for i, v in enumerate(obj['spec']['template']['spec']['volumes']):
    if v['name'] == 'cloudindustry-config-volume':
      obj['spec']['template']['spec']['volumes'].pop(i)

print json.dumps(obj, indent=2)
"
  kubectl -n tke get $resourceType $resource -o json | python -c """$script""" >$resource-$resourceType.yaml

  kubectl apply -f $resource-$resourceType.yaml
  rm $resource-$resourceType.yaml
}

function removeuser() {
  echo "==== 3. Execute OIDC rollback removeuser doing ===="
  kubectl get  platforms.business.tkestack.io platform-default -oyaml > default.yaml
  sed -i '/- admin/{n;d}' default.yaml
  kubectl apply -f default.yaml
  echo "==== 3. Execute OIDC rollback removeuser success ===="
}

function rollbackall() {
  echo "==== 1. Execute OIDC rollback rollbackConfigMap doing ===="
  rollbackConfigMap tke-gateway
  for cm in $(kubectl get cm -n tke -o=name | grep api | awk '{print $1}' | awk -F '/' '{print $2}'); do
    rollbackConfigMap $cm
  done
  echo "==== 1. Execute OIDC rollback rollbackConfigMap success ===="

  kubectl delete idp default

  echo "==== 2. Execute OIDC rollback rollback doing ===="
  for deploy in $(kubectl get deployment -n tke -o=name | grep api | awk '{print $1}' | awk -F '/' '{print $2}'); do
    rollbackWorkload deploy $deploy
  done
  rollbackWorkload daemonSet tke-gateway
  echo "==== 2. Execute OIDC rollback rollback success ===="
}
function checkall() {
  echo "==== 4. Execute OIDC rollback checkall doing ===="
  check_daemonset tke-gateway
  for deploy in $(kubectl get deployment -n tke -o=name | grep api | awk '{print $1}' | awk -F '/' '{print $2}'); do
    check_deployment $deploy
  done
  echo "==== 4. Execute OIDC rollback checkall success ===="
}

function check_deployment {
  check_workload_status_args deployment $1 tke .status.readyReplicas
  check=true
  while [ $check = true ]
  do
  check=false
  result=$(kubectl get deployment $1 -n tke --ignore-not-found -o go-template --template='{{if or (ne (.status.replicas) (.status.readyReplicas)) (ne (.status.replicas) (.status.updatedReplicas))}}false{{else}}true{{end}}')
  if [ $check = false ] && ([ -z $result ] || [ $result != true ]); then
    check=true
    echo "waiting $1"
    sleep 1
  fi
  done
}

function check_daemonset {
  check_workload_status_args daemonset $1 tke .status.numberReady
  check=true
  while [ $check = true ]
  do
  check=false
  result=$(kubectl get daemonset $1 -n tke --ignore-not-found -o go-template --template='{{if or (ne (.status.desiredNumberScheduled) (.status.numberReady)) (ne (.status.desiredNumberScheduled) (.status.updatedNumberScheduled))}}false{{else}}true{{end}}')
  if [ $check = false ] && ([ -z $result ] || [ $result != true ]); then
    check=true
    echo "waiting $1"
    sleep 1
  fi
  done
}

function check_workload_status_args {
    check=true
    while [ $check = true ]
    do
    check=false
    result=$(kubectl get $1 $2 -n $3 --ignore-not-found -o go-template --template='{{if not ('$4')}}false{{else}}true{{end}}')
    if [ $check = false ] && ([ -z $result ] || [ $result != true ]); then
        check=true
        echo "waiting $2 val $4"
        sleep 1
    fi
    done
}

backup
rollbackall
removeuser
checkall
