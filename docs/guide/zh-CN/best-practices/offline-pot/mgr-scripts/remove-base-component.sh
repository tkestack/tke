#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="all_func"
hosts="all"

help(){
  echo "show usage:"
  echo "del_redis: remove redis"
  echo "del_mysql: remove mysql"
  echo "del_pgsql: remove postgres"
  echo "del_prometheus: remove prometheus"
  echo "del_helm_tiller: remove helm tiller"
  echo "del_nginx_ingress: remove nginx ingress"
  echo "del_kafka: remove kafka"
  echo "del_elk: remove elk"
  echo "del_nfs: remove nfs"
  echo "del_minio: remove minio"
  echo "del_sgikes: remove sgikes "
  echo "del_harbor: remove harbor, default will not be remove harbor"
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


# remove redis
del_redis(){
  echo "###### remove redis start ######"
  # remove redis
  ansible-playbook -f 10 -i ../hosts --tags remove_redis ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  # remove redis cluster
  ansible-playbook -f 10 -i ../hosts --tags del_redis_cluster ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove redis end ######"
}

# remove mysql
del_mysql(){
  echo "###### remove mysql start ######"
  # remove mysql
  ansible-playbook -f 10 -i ../hosts --tags remove_mysql ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove mysql end ######"
}

# remove prometheus
del_prometheus(){
  echo "###### remove prometheus start ######"
  # remove prometheus
  ansible-playbook -f 10 -i ../hosts --tags remove_prometheus ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove prometheus end ######"
}

# remove helm-tiller
del_helm_tiller(){
  echo "###### remove helm-tiller start ######"
  # remove helm-tiller
  ansible-playbook -f 10 -i ../hosts --tags remove_helmtiller ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove helm-tiller end ######"
}

# remove nginx ingress
del_nginx_ingress(){
  echo "###### remove nginx ingress start ######"
  ansible-playbook -f 10 -i ../hosts --tags remove_nginx_ingress ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove nginx ingress end ######"
}

# remove kafka
del_kafka(){
  echo "###### remove kafka start ######"
  # remove kafka
  ansible-playbook -f 10 -i ../hosts --tags remove_kafka ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove kafka end ######"
}

# remove elk
del_elk(){
  echo "###### remove elk start ######"
  # remove elk
  ansible-playbook -f 10 -i ../hosts --tags remove_elk ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove elk end ######"
}

# remove postgres
del_pgsql(){
  echo "###### remove postgres start ######"
  ansible-playbook -f 10 -i ../hosts --tags del_postgres ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove postgres end ######"
}

# remove nfs
del_nfs(){
  echo "###### remove nfs start ######"
  ansible-playbook -f 10 -i ../hosts --tags remove_nfs ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove nfs end ######"
}

# remove minio
del_minio(){
  echo "###### remove minio start ######"
  ansible-playbook -f 10 -i ../hosts --tags remove_minio ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove minio end ######"
}

# remove sgikes
del_sgikes(){
  echo "###### remove sgikes start ######"
  # remove sgikes
  ansible-playbook -f 10 -i ../hosts --tags remove_sgikes ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove sgikes end ######"
}

# remove harbor, default will not be remove harbor
del_harbor(){
  echo "###### remove harbor start ######"
  # remove harbor
  ansible-playbook -f 10 -i ../hosts --tags del_harbor ../playbooks/base-component/base-component.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### remove harbor end ######"
}

# execute all function
all_func(){
  del_redis
  del_mysql
  del_pgsql
  del_prometheus
  del_nginx_ingress
  del_kafka
  del_elk
  del_sgikes
  del_minio
  del_nfs
  del_helm_tiller

}

main(){
  $CALL_FUN || help
}
main
