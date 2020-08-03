#!/bin/bash

set -e
BASE_DIR=../
cd $BASE_DIR
app_env_flag=$1 # Named after current customer
version='v1.2.4' # tkestack version
check_out='master' # business helm's git branch  check out

# get business helms
get_business_helms(){
  echo "###### get business helms start ######"
  # remove is not ${app_env_flag} business
  RELEASE_FILE=`ls helmfile.d/releases/*${app_env_flag}*.yaml`
  CENTER_LIST=`cat ${RELEASE_FILE} | grep -v '#' | grep 'chart: ../..' | awk -F/ '{print $3F}' | uniq`
  # remove not include ${app_env_flag} business chart dir base on center
  for c in `echo ${CENTER_LIST}`;
  do
    mv ${c} ${c}-builder
  done
  for d in `ls | grep -v builder | grep -v charts | grep -v helmfile.d`;
  do
    if [ -d "${d}" ]; then
      rm -rf ${d}
    fi
  done
  for c in `ls -d *-builder`;
  do
    mv ${c} `echo ${c} | awk -F-builder '{print $1}'`
  done
  # remove not include ${app_env_flag} business helmfile.d/commons
  if [ -d "helmfile.d/commons/" ]; then
    for c in `echo ${CENTER_LIST}`;
    do
      if [ -d "helmfile.d/commons/${c}" ]; then
        mv helmfile.d/commons/${c} helmfile.d/commons/${c}-builder
      fi
    done
    for i in `ls helmfile.d/commons/ | grep -v 'builder'`;
    do
      if [ -d "helmfile.d/commons/${i}" ]; then
        rm -rf helmfile.d/commons/${i}
      fi
    done
    for i in `ls helmfile.d/commons/ | grep 'builder'`;
    do
      if [ -d "helmfile.d/commons/${i}" ]; then
        mv helmfile.d/commons/${i} helmfile.d/commons/`echo ${i} | awk -F-builder '{print $1}'`
      fi
    done
  fi
  # remove not include ${app_env_flag} business helmfile.d/config
  if [ -d "helmfile.d/config/" ]; then
    for d in `ls helmfile.d/config/ | grep -v ${app_env_flag}`;
    do
      if [ -d "helmfile.d/config/${d}" ]; then
        rm -rf helmfile.d/config/${d}
      fi
    done
  fi
  # remove not include ${app_env_flag} business helmfile.d/environments
  if [ -d "helmfile.d/environments/" ]; then
    for d in `ls helmfile.d/environments/ | grep -v ${app_env_flag}`;
    do
      if [ -d "helmfile.d/environments/${d}" ]; then
        rm -rf helmfile.d/environments/${d}
      fi
    done
  fi
  # remove not include ${app_env_flag} business helmfile.d/releases
  if [ -d "helmfile.d/releases/" ]; then
    for f in `ls helmfile.d/releases/ | grep -v ${app_env_flag}`;
    do
      if [ -f "helmfile.d/releases/${f}" ]; then
        rm -f helmfile.d/releases/${f}
      fi
    done
  fi
  # get chart list
  CHART_LIST=`cat ${RELEASE_FILE} | grep -v '#' | grep 'chart: ../..' | awk -F/ '{print $3"/"$4}' | uniq`
  # remove not include ${app_env_flag} business chart base on center's business
  for c in `echo ${CHART_LIST} | tr -d '\r'`;
  do
    rm -rf ${c}/values
    mv ${c} ${c}-builder
  done
  for c in `echo ${CENTER_LIST}`;
  do
    for d in `ls ${c} | grep -v builder | grep -v tools`;
    do
      if [ -d "${c}/${d}" ]; then
        rm -rf ${c}/${d}
      fi
    done
    for i in `ls -d ${c}/*-builder`;
    do
      mv ${i} `echo ${i} | awk -F-builder '{print $1}'`
    done
  done
  # get registry secrets
  if [ "x${remote_registry}" == "xtrue" ]; then
    for s in `echo ${CENTER_LIST}`;
    do
      if [ -f "${BUILDER_PDIR}reg.mydomain.com.secrets/${s}.reg.mydomain.com.yaml" ]; then
        mkdir -p ${BASE_DIR}/roles/business/helms/secrets
        \cp ${BUILDER_PDIR}reg.mydomain.com.secrets/${s}.reg.mydomain.com.yaml ${BASE_DIR}/roles/business/helms/secrets/
      else
        echo "current business use remote registry: "
        echo "${BUILDER_PDIR}reg.mydomain.com.secrets/${s}.reg.mydomain.com.yaml not exist, please check!!!" && exit 0
      fi
    done
  fi
  rm -rf ${BASE_DIR}/roles/business/helms/.git
#  # tmp use /data/helm-base/charts
#  rm -rf ${BASE_DIR}/roles/business/helms/charts/ && cp -r ${BUILDER_PDIR}helm-base/charts ${BASE_DIR}/roles/business/helms/
#  echo "###### get business helms end ######"
}
get_business_helms
