#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="all_func"
hosts="all"

help(){
  echo "show usage:"
  echo "helm_tiller_dpl_check: helm tiller deploy check!"
  echo "redis_health_check: redis health check!"
  echo "mysql_health_check: mysql health check!"
  echo "postgres_health_check: postgres health check!"
  echo "prometheus_dpl_check: prometheus deploy check!"
  echo "nginx_ingress_health_check: nginx ingress controller health check!"
  echo "kafka_dpl_check: kafka deploy check!"
  echo "elk_dpl_check: elk deploy check!"
  echo "nfs_health_check: nfs health check!"
  echo "minio_dpl_check: minio deploy check!"
  echo "sgikes_dpl_check: sgikes deploy check!"
  echo "busi_dpl_check: business deploy check!"
  echo "harbor_dpl_check: harbor deploy check!, default will not be check!"
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

# helm tiller deploy check 
helm_tiller_dpl_check(){
  echo "###### helm tiller deploy check start ######"
  ansible-playbook -f 10 -i ../hosts --tags helmtiller_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### helm tiller deploy check end ######"
}

# redis health check
redis_health_check(){
  echo "###### redis health check start ######"
  ansible-playbook -f 10 -i ../hosts --tags redis_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # redis cluster health check
  ansible-playbook -f 10 -i ../hosts --tags redis_cluster_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### redis health check end ######"
}

# mysql health check
mysql_health_check(){
  echo "###### mysql health check start ######"
  ansible-playbook -f 10 -i ../hosts --tags mysql_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### mysql health check end ######"
}

# postgres health check
postgres_health_check(){
  echo "###### postgres health check start ######"
  ansible-playbook -f 10 -i ../hosts --tags pgsql_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### postgres health check end ######"
}

# prometheus deploy check
prometheus_dpl_check(){
  echo "###### prometheus deploy check start ######"
  ansible-playbook -f 10 -i ../hosts --tags prometheus_health_ckeck ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### prometheus deploy check end ######"
}

# nginx ingress controller health check
nginx_ingress_health_check(){
  echo "###### nginx ingress controller health check start ######"
  ansible-playbook -f 10 -i ../hosts --tags ingress_controller_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### nginx ingress controller health check end ######"
} 

# kafka deploy check
kafka_dpl_check(){
  echo "###### kafka deploy check start ######"
  ansible-playbook -f 10 -i ../hosts --tags kafka_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### kafka deploy check end ######"
} 

# elk deploy check
elk_dpl_check(){
  echo "###### elk deploy check start ######"
  ansible-playbook -f 10 -i ../hosts --tags elk_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### elk deploy check end ######"
}

# nfs health check
nfs_health_check(){
  echo "###### nfs health check start ######"
  ansible-playbook -f 10 -i ../hosts --tags nfs_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### nfs health check end ######"
}

# minio deploy check
minio_dpl_check(){
  echo "###### minio deploy check start ######"
  ansible-playbook -f 10 -i ../hosts --tags minio_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### minio deploy check end ######"
}

# sgikes deploy check
sgikes_dpl_check(){
  echo "###### sgikes deploy check start ######"
  ansible-playbook -f 10 -i ../hosts --tags sgikes_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### sgikes deploy check end ######"
}

# business deploy check
busi_dpl_check(){
  echo "###### business deploy check start ######"
  ansible-playbook -f 10 -i ../hosts --tags busi_health_check ../playbooks/business/business.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### business deploy check end ######"
}

# harbor deploy check
harbor_dpl_check(){
  echo "###### harbor deploy check start ######"
  ansible-playbook -f 10 -i ../hosts --tags harbor_health_check ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### harbor deploy check end ######"
}

all_func(){
  helm_tiller_dpl_check
  redis_health_check
  mysql_health_check
  postgres_health_check
  prometheus_dpl_check
  nginx_ingress_health_check
  kafka_dpl_check
  elk_dpl_check
  nfs_health_check
  minio_dpl_check
  sgikes_dpl_check
  busi_dpl_check
}

main(){
  $CALL_FUN || help
}
main
