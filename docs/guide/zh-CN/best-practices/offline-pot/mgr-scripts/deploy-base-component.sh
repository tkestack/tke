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

# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="all_func"
hosts="all"

help(){
  echo "show usage:"
  echo "load_and_push_images: load and push images"
  echo "dpl_redis: deploy redis"
  echo "dpl_mysql: deploy mysql"
  echo "dpl_pgsql: deploy postgres"
  echo "dpl_prometheus: deploy prometheus"
  echo "dpl_helm_tiller: deploy helm tiller"
  echo "dpl_nginx_ingress: deploy nginx ingress"
  echo "dpl_kafka: deploy kafka"
  echo "dpl_elk: deploy elk"
  echo "dpl_nfs: deploy nfs"
  echo "dpl_minio: deploy nfs"
  echo "dpl_sgikes: deploy sgikes for search"
  echo "dpl_harbor: deploy harbor, default will not be deploy harbor"
  echo "all_func: execute all function, -f default value is all_func !!!"
  exit 0
}

while getopts ":f:h:" opt
do
  case $opt in
    f)
    CALL_FUN="${OPTARG}"
    ;;
    h)
    hosts="${OPTARG}"
    ;;
    ?)
    echo "unkown args! just suport -f[call function] and -h[ansible hosts group] arg!!!"
    exit 0;;
  esac
done

# load and push images
load_and_push_images(){
  echo "###### load and push images start ######"
  ansible-playbook -f 10 -i ../hosts --tags load_and_push_images ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### load and push images end ######"
}

# deploy redis
dpl_redis(){
  echo "###### deploy redis start ######"
  # redis init
  ansible-playbook -f 10 -i ../hosts --tags redis_mgr_init ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # deploy redis
  ansible-playbook -f 10 -i ../hosts --tags deploy_redis ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # deploy redis cluster
  ansible-playbook -f 10 -i ../hosts --tags dpl_redis_cluster ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy redis end ######"
}

# deploy mysql
dpl_mysql(){
  echo "###### deploy mysql start ######"
  # mysql init
  ansible-playbook -f 10 -i ../hosts --tags mysql_mgr_init ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # deploy mysql
  ansible-playbook -f 10 -i ../hosts --tags deploy_mysql ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy mysql end ######"
}

# deploy prometheus
dpl_prometheus(){
  echo "###### deploy prometheus start ######"
  # prometheus init
  ansible-playbook -f 10 -i ../hosts --tags prometheus_init ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # deploy prometheus
  ansible-playbook -f 10 -i ../hosts --tags deploy_prometheus ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy prometheus end ######"
}

# deploy helm-tiller
dpl_helm_tiller(){
  echo "###### deploy helm-tiller start ######"
  # helm tiller init
  ansible-playbook -f 10 -i ../hosts --tags helmtiller_init ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # deploy helm-tiller
  ansible-playbook -f 10 -i ../hosts --tags deploy_helmtiller ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy helm-tiller end ######"
}

# deploy nginx ingress
dpl_nginx_ingress(){
  echo "###### deploy nginx ingress start ######"
  ansible-playbook -f 10 -i ../hosts --tags deploy_nginx_ingress ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy nginx ingress end ######"
}

# deploy kafka
dpl_kafka(){
  echo "###### deploy kafka start ######"
  # kafka init
  ansible-playbook -f 10 -i ../hosts --tags kafka_init ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # deploy kafka
  ansible-playbook -f 10 -i ../hosts --tags deploy_kafka ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy kafka end ######"
}

# deploy elk
dpl_elk(){
  echo "###### deploy elk start ######"
  # elk init
  ansible-playbook -f 10 -i ../hosts --tags elk_init ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # deploy elk
  ansible-playbook -f 10 -i ../hosts --tags deploy_elk ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy elk end ######"
}

# deploy postgres
dpl_pgsql(){
  echo "###### deploy postgres start ######"
  ansible-playbook -f 10 -i ../hosts --tags dpl_postgres ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy postgres end ######"
}

# deploy nfs
dpl_nfs(){
  echo "###### deploy nfs start ######"
  ansible-playbook -f 10 -i ../hosts --tags dpl_nfs ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy nfs end ######"
}

# deploy minio
dpl_minio(){
  echo "###### deploy minio start ######"
  # minio init
  ansible-playbook -f 10 -i ../hosts --tags minio_init ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # minio deploy
  ansible-playbook -f 10 -i ../hosts --tags deploy_minio ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy minio end ######"
}

# deploy sgikes
dpl_sgikes(){
  echo "###### deploy sgikes start ######"
  # sgikes init
  ansible-playbook -f 10 -i ../hosts --tags sgikes_init ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # deploy sgikes
  ansible-playbook -f 10 -i ../hosts --tags dpl_sgikes ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy sgikes end ######"
}

# deploy harbor, default will not be deploy
dpl_harbor(){
  echo "###### deploy harbor start ######"
  # deploy harbor
  ansible-playbook -f 10 -i ../hosts --tags dpl_harbor ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### deploy harbor end ######"
}

# execute all function
all_func(){
  load_and_push_images
  dpl_redis
  dpl_mysql
  dpl_pgsql
  dpl_prometheus
  dpl_helm_tiller
  dpl_nginx_ingress
  dpl_kafka
  dpl_elk
  dpl_sgikes
  dpl_nfs
  dpl_minio
}

main(){
  $CALL_FUN || help
}
main
