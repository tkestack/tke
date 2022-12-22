#!/usr/bin/env bash
#set -x

function setConfigEnvs() {
  echo "==== 0. Execute OIDC config setConfigEnvs doing ===="
  kubectl get cm oidc-config -n tke
  if [ $? -ne 0 ]; then
    echo "Ignore ERROR: configmap oidc-config not found, creating oidc-config configmap for persistent oidc value in tkeanywhere ..."
    if [ "$1" != "" ]; then
      cp $1 oidc.conf
    fi
    sed -i s/[[:space:]]//g oidc.conf
    kubectl create cm oidc-config -n tke --from-file=./oidc.conf
  else
    kubectl get cm oidc-config -n tke -ojson | python -c "import json; import sys; obj=json.load(sys.stdin); print obj['data']['oidc.conf']" >oidc.conf
  fi

  . oidc.conf
  vip=$(kubectl get cls global -o=jsonpath='{.spec.features.ha.tke.vip}')
  region="ap-guangzhou"
  rm oidc.conf
  echo "==== 0. Execute OIDC config setConfigEnvs success ===="
}

function validateInput() {
  echo "==== 1. Execute OIDC config validateInput doing ===="
  if [ "$ca_crt" == "" ]; then
    curl -v $issuer_url/.well-known/openid-configuration
    if [ $? -ne 0 ]; then
      echo "Validate ERROR: validate issuer URL $issuer_url failed."
      echo "Please check use cmd: curl -v $issuer_url/.well-known/openid-configuration"
      exit 1
    fi
  else
    curl --cacert $ca_crt $issuer_url/.well-known/openid-configuration
    if [ $? -ne 0 ]; then
      echo "Validate ERROR: validate $ca_crt and $issuer_url failed."
      echo "Please check use cmd: curl --cacert $ca_crt $issuer_url/.well-known/openid-configuration"
      exit 1
    fi
  fi
  echo "==== 1. Execute OIDC config validateInput success ===="
}

function backup() {
  echo "==== 2. Execute OIDC config backup doing ===="
  api=$(kubectl -n tke get -o=name deployment | grep api)
  cm=$(kubectl -n tke get -o=name configmap | grep api)
  gw=$(kubectl -n tke get -o=name configmap | grep gateway)
  ds=$(kubectl -n tke get -o=name daemonset | grep gateway)
  r=("$api $cm $ds $gw")
  echo ======= backup current resources =======
  dir=/opt/oidcbackup
  mkdir -p ${dir}/deployment.apps ${dir}/daemonset.apps ${dir}/configmap
  for n in $r; do
    kubectl -n tke get -o=yaml $n >${dir}/$n.yaml
  done
  echo "==== 2. Execute OIDC config backup success ===="
}

function createcms() {
  echo "==== 3. Execute OIDC config createcms doing ===="
  # create oidc-ca configmap
  kubectl get cm oidc-ca -n tke
  if [ $? -ne 0 ] && [ "$ca_crt" != "" ]; then
    echo "Ignore ERROR: configmap oidc-ca not found, creating oidc-ca configmap ..."
    kubectl create cm oidc-ca -n tke --from-file=$ca_crt
  fi

  # create cloudindustry-config configmap
  cat <<EOF | kubectl apply -f -
kind: ConfigMap
metadata:
  name: cloudindustry-config
  namespace: tke
apiVersion: v1
data:
  config: |
    {
        "secret_id": "$secret_id",
        "secret_key": "$secret_key",
        "endpoint": "$endpoint",
        "region": "$region",
        "master_id": "$master_id"
    }
EOF
  echo "==== 3. Execute OIDC config createcms success ===="
}

function createOIDCAuthCMTmpFile() {
  content=$(
    cat <<EOF
\n
[authentication.oidc]\n
client_secret = "$secret_key"\n
client_id = "$secret_id"\n
issuer_url = "$issuer_url"\n
username_prefix ="-"\n
username_claim = "name"\n
groups_claim = "groups"\n
tenantid_claim = "federated_claims"
EOF
  )
  kubectl get cm oidc-ca -n tke
  if [ $? -eq 0 ]; then
    content=$(
      cat <<EOF
${content}\n
ca_file = "/app/oidc/ca.crt"
EOF
    )

  fi
  echo -e $content >oidc_auth_tmp.txt
  sed -i 's/^/     /g' oidc_auth_tmp.txt
  sed -i '1d' oidc_auth_tmp.txt
  sed -i '1{x;p;x;}' oidc_auth_tmp.txt
}

function removeOIDCAuthCMTmpFile() {
  rm oidc_auth_tmp.txt
}

function modifyConfigMap() {
  echo ======= modify configmap $1 =======
  file=./$1-cm.yaml
  kubectl -n tke get cm $1 -o yaml >$file

  start=$(sed -n '/last-applied-configuration/=' $file)
  if [ "$start" != "" ]; then
    end=$(($start + 1))
    sed -i "$start, $end d" $file
  fi

  if [ "$1" = "tke-auth-api" ]; then
    line=$(sed -n '/authentication.oidc/=' $file)
    if [ "$line" = "" ]; then
      sed -i "/privileged_username/r oidc_auth_tmp.txt" $file
    fi
    line=$(sed -n '/init_tenant_type/=' $file)
    if [ "$line" = "" ]; then
      sed -i '/assets_path/a\    init_tenant_type = "cloudindustry"' $file
    fi
    line=$(sed -n '/init_tenant_id/=' $file)
    if [ "$line" = "" ]; then
      sed -i '/assets_path/a\    init_tenant_id = "default"' $file
    fi
    line=$(sed -n '/cloudindustry_config_file/=' $file)
    if [ "$line" = "" ]; then
      sed -i '/assets_path/a\    cloudindustry_config_file = "/app/cloudindustry/config"' $file
    fi

    line=$(sed -n '/init_client_id/=' $file)
    if [ "$line" != "" ]; then
      sed -i "${line}c \ \ \ \ init_client_id = \"${secret_id}\"" $file
    fi
    line=$(sed -n '/init_client_secret/=' $file)
    if [ "$line" != "" ]; then
      sed -i "${line}c \ \ \ \ init_client_secret = \"${secret_key}\"" $file
    fi
  else
    kubectl get cm oidc-ca -n tke
    if [ $? -ne 0 ]; then
      echo "Ignore ERROR: configmap oidc-ca not found, ignore ca.crt ..."
      sed -i '/^      ca_file = "\/app\/certs\/ca.crt"$/d' $file
    else
      line=$(sed -n '/^      ca_file/=' $file)
      if [ "$line" != "" ]; then
        sed -i "${line}c \ \ \ \ \ \ ca_file = \"\/app\/oidc\/ca.crt\"" $file
      fi
    fi

    if [ "$1" = "tke-gateway" ]; then
      sed -i '/disableOIDCProxy/d' $file
      sed -i "/kind: GatewayConfiguration/a\    disableOIDCProxy: true" $file
      sed -i '/generic/d' $file
      sed -i "/port = 80/a\    [generic]" $file
      sed -i '/external_hostname/d' $file
      sed -i "/generic/a\    external_hostname = \"${tke_domain_name}\"" $file
    fi

    sed -i '/external_issuer_url/d' $file

    sed -i '/client_secret/d' $file
    sed -i "/authentication.oidc/a\      client_secret = \"${secret_key}\"" $file

    sed -i '/client_id/d' $file
    sed -i "/authentication.oidc/a\      client_id = \"${secret_id}\"" $file

    sed -i '/issuer_url/d' $file
    sed -i "/authentication.oidc/a\      issuer_url = \"${issuer_url}\"" $file
  fi

  sed -i '/resourceVersion/d' $file
  sed -i '/uid/d' $file
  kubectl apply -f $file
  rm $file
}

function createOIDCVolumeTmpFiles() {
  kubectl get cm oidc-ca -n tke
  if [ $? -ne 0 ]; then
    cat >>./oidc_volumeMounts_tmp.txt <<EOF
        - mountPath: /app/cloudindustry
          name: cloudindustry-config-volume
EOF

    cat >>./oidc_volumes_tmp.txt <<EOF
      - configMap:
          defaultMode: 420
          name: cloudindustry-config
        name: cloudindustry-config-volume
EOF

  else
    cat >>./oidc_volumeMounts_tmp.txt <<EOF
        - mountPath: /app/oidc
          name: oidc-ca-volume
        - mountPath: /app/cloudindustry
          name: cloudindustry-config-volume
EOF

    cat >>./oidc_volumes_tmp.txt <<EOF
      - configMap:
          defaultMode: 420
          name: oidc-ca
        name: oidc-ca-volume
      - configMap:
          defaultMode: 420
          name: cloudindustry-config
        name: cloudindustry-config-volume
EOF

  fi
}

function removeOIDCVolumeTmpFiles() {
  rm oidc_volumeMounts_tmp.txt
  rm oidc_volumes_tmp.txt
}

function createOIDCHostAliasTmpFile() {
  if [ "$hostnames" != "" ] && [ "$vip" != "" ]; then
    file=oidc_hostAlias_tmp.txt
    rm -rf $file
    echo "      hostAliases:" >>$file
    echo "      - hostnames:" >>$file
    arr=(${hostnames//,/ })
    for name in "${arr[@]}"; do
      #    remove ','
      name=$(echo "$name" | sed -e "s/,$//")
      #    remove '"'
      name=$(echo "$name" | sed -e "s/^\"//" -e "s/\"$//")
      echo "****** ${name}"
      echo "        - ${name}" >>$file
    done
    echo "        ip: ${vip}" >>$file
  fi
}

function rmOIDCHostAliasTmpFile() {
  rm -f oidc_hostAlias_tmp.txt
}

function modifyResource() {
  echo ======= modify $1 $2 =======
  file=./$2-$1.yaml

  kubectl -n tke annotate $1 $2 tkeanywhere/oidc="true" --overwrite=true
  kubectl -n tke get $1 $2 -o yaml >$file

  start=$(sed -n '/last-applied-configuration/=' $file)
  if [ "$start" != "" ]; then
    end=$(($start + 1))
    sed -i "$start, $end d" $file
  fi

  line=$(sed -n '/hostAliases/=' $file)
  if [ "$line" = "" ] && [ "$hostnames" != "" ] && [ "$vip" != "" ]; then
    sed -i "/dnsPolicy/r oidc_hostAlias_tmp.txt" $file
  fi

  line=$(sed -n '/oidc-ca-volume/=' $file)
  if [ "$line" = "" ]; then
    sed -i "/volumes/r oidc_volumes_tmp.txt" $file
    sed -i "/volumeMounts/r oidc_volumeMounts_tmp.txt" $file
  fi

  kubectl apply -f $file
  rm $file
}

function adduser() {
  echo "==== 3. Execute OIDC config adduser doing ===="
  kubectl get platforms.business.tkestack.io platform-default -oyaml >default.yaml
  start=$(sed -n '/last-applied-configuration/=' default.yaml)
  if [ "$start" != "" ]; then
    end=$(($start + 1))
    sed -i "$start, $end d" default.yaml
  fi
  cat default.yaml | grep $username
  if [ $? -ne 0 ]; then
    sed -i "/- admin/a\  - ${username}" default.yaml
    kubectl apply -f default.yaml
  fi
  rm default.yaml
  echo "==== 3. Execute OIDC config adduser success ===="
}

function configAll() {
  echo "==== 4. Execute OIDC config modifyConfigMap doing ===="
  createOIDCAuthCMTmpFile
  anno=$(kubectl -n tke get ds tke-gateway -o=jsonpath='{.metadata.annotations.tkeanywhere/oidc}')
  if [ "$anno" == "" ]; then
    modifyConfigMap tke-gateway
  fi
  for cm in $(kubectl get cm -n tke -o=name | grep api | awk '{print $1}' | awk -F '/' '{print $2}'); do
    anno=$(kubectl -n tke get deploy $cm -o=jsonpath='{.metadata.annotations.tkeanywhere/oidc}')
    if [ "$anno" == "" ]; then
      modifyConfigMap $cm
    fi
  done
  echo "==== 4. Execute OIDC config modifyConfigMap success ===="
  kubectl delete idp default
  removeOIDCAuthCMTmpFile

  createOIDCVolumeTmpFiles
  createOIDCHostAliasTmpFile
  echo "==== 5. Execute OIDC config modifyResource doing ===="
  for deploy in $(kubectl get deployment -n tke -o=name | grep api | awk '{print $1}' | awk -F '/' '{print $2}'); do
    anno=$(kubectl -n tke get deploy $cm -o=jsonpath='{.metadata.annotations.tkeanywhere/oidc}')
    if [ "$anno" == "" ]; then
      modifyResource deploy $deploy
    fi
  done

  anno=$(kubectl -n tke get ds tke-gateway -o=jsonpath='{.metadata.annotations.tkeanywhere/oidc}')
  if [ "$anno" == "" ]; then
    modifyResource daemonset tke-gateway
  fi
  removeOIDCVolumeTmpFiles
  rmOIDCHostAliasTmpFile
  echo "==== 5. Execute OIDC config modifyResource success ===="
}

function checkall() {
  echo "==== 6. Execute OIDC config checkall doing ===="
  check_daemonset tke-gateway
  for deploy in $(kubectl get deployment -n tke -o=name | grep api | awk '{print $1}' | awk -F '/' '{print $2}'); do
    check_deployment $deploy
  done
  echo "==== 6. Execute OIDC config checkall success ===="
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

setConfigEnvs $1
validateInput
backup
createcms
configAll
adduser
checkall
