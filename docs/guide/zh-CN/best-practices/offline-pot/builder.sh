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
app_env_flag='' # Named after current customer
builder_cfg='builder.cfg' #  builder config file name

help(){
  echo "show usage:"
  echo "get_business_helms: get business helms"
  echo "get_base_component_helms: get base component helms"
  echo "get_business_images: get business images"
  echo "build_to_tgz: build to tgz"
  echo "all_func: execute all function, -f default value is all_func !!!"
  echo "app_env_flag arg(-a) cloud not empty, app_env_flag named after current customer!!!"
  exit 0
}

while getopts ":f:a:b:" opt
do
  case $opt in
    f)
    CALL_FUN="${OPTARG}"
    ;;
    a)
    app_env_flag="${OPTARG}"
    ;;
    b)
    builder_cfg="${OPTARG}"
    ;;
    ?)
    echo "unkown args!just suport -f[call function],-a[Named after current customer],-b[builder config file name]  arg!!!"
    exit 0;;
  esac
done

# reference builder config file
. ${builder_cfg}

# get builder.sh's parent dir
BUILDER_PDIR=`echo ${BASE_DIR} | awk -Foffline-pot '{print $1}'`

# get business helms
get_business_helms(){
  echo "###### get business helms start ######"
  if [ -d "${BASE_DIR}/roles/business/helms" ]; then
     rm -rf ${BASE_DIR}/roles/business/helms
  fi
  cd ${BASE_DIR}/roles/business/
  # get business helms or workerload yaml from git repo
  git clone ${busi_git_url}
  busi_git_dir=`echo ${busi_git_url} | awk -F\/ '{print $NF}' | awk -F. '{print $1}'`
  if [ -d "${BASE_DIR}/roles/business/${busi_git_dir}" ]; then
    cd ${BASE_DIR}/roles/business/${busi_git_dir}
    git checkout ${busi_branch}
    mv ${BASE_DIR}/roles/business/${busi_git_dir} ${BASE_DIR}/roles/business/helms
  else
    echo "git clone failed,please check!!!" && exit 0
  fi
  cd ${BASE_DIR}/roles/business/helms
  if [ -d "${BASE_DIR}/roles/business/helms/helmfile.d" ]; then
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
        if [ -f "${BUILDER_PDIR}${remote_img_registry_url}.secrets/${s}.${remote_img_registry_url}.yaml" ]; then
          mkdir -p ${BASE_DIR}/roles/business/helms/secrets
          \cp ${BUILDER_PDIR}${remote_img_registry_url}.secrets/${s}.${remote_img_registry_url}.yaml ${BASE_DIR}/roles/business/helms/secrets/
        else
          echo "current business use remote registry: "
          echo "${BUILDER_PDIR}${remote_img_registry_url}.secrets/${s}.${remote_img_registry_url}.yaml not exist, please check!!!" && exit 0
        fi
      done
    fi
    rm -rf ${BASE_DIR}/roles/business/helms/.git
  fi
  echo "###### get business helms end ######"
}

# get base-component-helms
get_base_component_helms(){
  echo "###### get base component helms start ######"
  if [ `ls ${BUILDER_PDIR}base-component-helms/helms | wc -l` -gt 0 ]; then
    if [ -d "${BASE_DIR}/roles/base-component/helms/" ]; then
      rm -rf ${BASE_DIR}/roles/base-component/helms/
    fi
    # create base-component's helms dir
    if [ ! -d "${BASE_DIR}/roles/base-component/helms/" ]; then
      mkdir -p ${BASE_DIR}/roles/base-component/helms/
    fi
    for s in ${server_set[@]};
    # copy base commons's helms
    do
      if [ -d "${BUILDER_PDIR}base-component-helms/helms/${s}" ]; then
        if [ "${s}" == "redis" ]; then
          continue
        elif [ "${s}" == "redis_cluster" ]; then
          \cp -r ${BUILDER_PDIR}base-component-helms/helms/redis/* ${BASE_DIR}/roles/base-component/helms/
        fi
        \cp -r ${BUILDER_PDIR}base-component-helms/helms/${s}/* ${BASE_DIR}/roles/base-component/helms/
      fi
    done
  else
    echo "${BUILDER_PDIR}base-component-helms/helms is empty, please check!!!" && exit 0
  fi
  echo "###### get base component helms end ######"
}

# get business images
get_business_images(){
  echo "###### get business images start ######"
  if [ -d "${BASE_DIR}/roles/business/helms" ]; then
    cd ${BASE_DIR}/roles/business/helms
  else
    echo "${BASE_DIR}/roles/business/helms not exist, please check!!!" && exit 0
  fi
  if [ ! -d "${BUILDER_PDIR}offline-pot-images/" ]; then
    mkdir -p ${BUILDER_PDIR}offline-pot-images
  elif [ `ls ${BUILDER_PDIR}offline-pot-images | wc l ` -gt 0 ]; then
    rm -rf ${BUILDER_PDIR}offline-pot-images/*.tar
  fi
    # get image list
    if [ -d "helmfile.d/config/" ]; then
      CONFIG_DIRS=`ls helmfile.d/config/`
      for c in `echo ${CONFIG_DIRS}`;
      do
        if [ -d "helmfile.d/config/${c}" ]; then
          # get images,pull and save
          for f in `ls helmfile.d/config/${c}/`;
          do
            if [ -f "helmfile.d/config/${c}/$f" ]; then
              repository=`cat helmfile.d/config/${c}/$f | grep -v '#' | grep repository | awk -F: '{print $2}' | tr -d '\r'`
              tag=`cat helmfile.d/config/${c}/$f | grep -v '#' | grep 'tag:' | awk -F: '{print $2}' | awk -F\" '{print $2}'`
              docker pull ${repository}:${tag}
              repo_fname=`echo ${repository} | awk -F/ '{print $NF}'`
              if [ -f "${BUILDER_PDIR}offline-pot-images/${repo_fname}.${tag}.tar" ]; then
                rm -f ${BUILDER_PDIR}offline-pot-images/${repo_fname}.${tag}.tar
              fi
              docker save ${repository}:${tag} > ${BUILDER_PDIR}offline-pot-images/${repo_fname}.${tag}.tar
            fi
          done
        fi
      done
    fi
  echo "###### get business images end ######"
}

# build to tgz
build_to_tgz(){
  echo "###### build to tgz start ######"

  if [ "${builder_type}" == "all" ] || [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep -v ^business$ | grep -v ^tkestack$ | wc -l ` -gt 0 ]; then
    # check offline-pot-images-base whether is empty dir
    if [ -d "${BUILDER_PDIR}offline-pot-images-base" ]; then
      if [ ! -d "${BUILDER_PDIR}offline-pot-images" ]; then
        mkdir -p ${BUILDER_PDIR}offline-pot-images
      elif [ "x$remote_registry" == "xtrue" ]; then
        rm -rf ${BUILDER_PDIR}offline-pot-images/*.tar
      fi
    else
      echo "${BUILDER_PDIR}offline-pot-images-base not exist,please check !!!" && exit 0
    fi
    # check offline-pot-tgz-base is empty dir
    if [ -d "${BUILDER_PDIR}offline-pot-tgz-base" ]; then
      if [ ! -d "${BUILDER_PDIR}offline-pot-tgz" ]; then
        mkdir -p ${BUILDER_PDIR}offline-pot-tgz
      else
        rm -rf ${BUILDER_PDIR}offline-pot-tgz
        mkdir -p ${BUILDER_PDIR}offline-pot-tgz
      fi
    else
      echo "${BUILDER_PDIR}offline-pot-tgz-base not exist,please check !!!" && exit 0
    fi
  fi
  # get commons's images or commons's tgz
  if [ "${builder_type}" == "all" ]; then
    # copy all  commons's images
    for i in `ls ${BUILDER_PDIR}offline-pot-images-base/`;
    do
      \cp -r ${i}/* ${BUILDER_PDIR}offline-pot-images/
    done
    # copy all commons's tgz
    for i in `ls ${BUILDER_PDIR}offline-pot-tgz-base/`;
    do
      \cp -r ${i}/* ${BUILDER_PDIR}offline-pot-tgz/
    done
  elif [[ "${builder_type}" == "custom" ]]; then
    if [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep -v ^business$ | grep -v ^tkestack$ | wc -l ` -gt 0 ]; then
      # get commons's images
      for s in ${server_set[@]};
      do
        # copy server group's images
        if [ -d "${BUILDER_PDIR}offline-pot-images-base/${s}" ]; then
          if [ "${s}" == "redis" ]; then
            continue
          elif [ "${s}" == "redis_cluster" ]; then
            \cp -r ${BUILDER_PDIR}offline-pot-images-base/redis/* ${BUILDER_PDIR}offline-pot-images/
          fi
          \cp -r ${BUILDER_PDIR}offline-pot-images-base/${s}/* ${BUILDER_PDIR}offline-pot-images/
        fi
      done
    fi
    if [ `ls ${BUILDER_PDIR}offline-pot-images-base/yum-repo/ | wc -l ` -gt 0 ]; then
      \cp -r ${BUILDER_PDIR}offline-pot-images-base/yum-repo/* ${BUILDER_PDIR}offline-pot-images/
    else
      echo -e "\033[33m yum-repo commons's images is empty, please check whether need offline yum-repo to init base tools !!! \033[0m"
    fi
    if [ `ls ${BUILDER_PDIR}offline-pot-images-base/debug-img/ | wc -l ` -gt 0 ]; then
      \cp -r ${BUILDER_PDIR}offline-pot-images-base/debug-img/* ${BUILDER_PDIR}offline-pot-images/
    else
      echo -e "\033[33m debug-img commons's images is empty, please check whether need debug image for debug network or dns!!! \033[0m"
    fi
    # get server group's tgz
    for s in ${server_set[@]};
    do
      if [ -d "${BUILDER_PDIR}offline-pot-tgz-base/${s}" ]; then
        \cp -r ${BUILDER_PDIR}offline-pot-tgz-base/${s}/* ${BUILDER_PDIR}offline-pot-tgz/
      fi
    done
    # config remote registry url and get remote registry cert files
    if [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep ^business$ | wc -l ` -eq 1 ] && [ "x$remote_registry" == "xtrue" ] ; then
      # set remote_img_registry_url on hosts.tpl
      sed -i "s/remote_img_registry_url='.*'/remote_img_registry_url='${remote_img_registry_url}'/g" ${BUILDER_PDIR}offline-pot/hosts.tpl
      # when remote registry need cert' , must be had cert file .
      if [ -f "${BUILDER_PDIR}offline-pot-tgz-base/${remote_img_registry_url}.cert/${remote_img_registry_url}.cert.tar.gz" ]; then
        \cp ${BUILDER_PDIR}offline-pot-tgz-base/${remote_img_registry_url}.cert/${remote_img_registry_url}.cert.tar.gz  ${BUILDER_PDIR}offline-pot-tgz/
      else
        echo -e "\033[33m ${remote_img_registry_url}.cert.tar.gz tgz not exist , please check remote registry whether config registry's cert !!! \033[0m"
      fi
    fi
    # get helm's binary files
    if [ -d "${BUILDER_PDIR}offline-pot-tgz-base/helms/" ]; then
      \cp ${BUILDER_PDIR}offline-pot-tgz-base/helms/* ${BUILDER_PDIR}offline-pot-tgz/
    fi
    # get harbor install file when dose not install tkestack
    if [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep ^tkestack$ | wc -l ` -eq 0 ] && [ `ls ${BUILDER_PDIR}offline-pot-images/ | wc -l` -gt 0 ]; then
      if [ -d "${BUILDER_PDIR}offline-pot-tgz-base/harbor" ]; then
        \cp -r ${BUILDER_PDIR}offline-pot-tgz-base/harbor/* ${BUILDER_PDIR}offline-pot-tgz/
      else
        echo -e "\033[33m harbor install pkg not exist, will be influences offline image's  load !!! \033[0m" && exit 0
      fi
    fi
  fi
  # get tkestack
  if [ "${builder_type}" == "all" ] || [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep ^tkestack$ | wc -l ` -eq 1 ]; then
    if [ ! -d "${BUILDER_PDIR}tkestack/" ]; then
      mkdir -p ${BUILDER_PDIR}tkestack/
    else
      if [ `ls ${BUILDER_PDIR}tkestack/tke-installer* | wc -l` -eq 0 ]; then
        cd ${BUILDER_PDIR}tkestack/
        # get tkestack pkg
        if [ `echo "$version" | awk -Fv '{print $2}' | awk -F. '{print $1$2$3}'` -lt 130 ]; then
           wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-x86_64-$version.run{,.sha256}
        else
          wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-linux-$arch-$version.run{,.sha256}
        fi
      elif [ `ls ${BUILDER_PDIR}tkestack/tke-installer* | wc -l` -gt 0 ] && [ `ls ${BUILDER_PDIR}tkestack/*${version}* | wc -l` -lt 2 ]; then
        cd ${BUILDER_PDIR}tkestack/
        rm -f ${BUILDER_PDIR}tkestack/tke-installer*
        # get tkestack pkg
        if [ `echo "$version" | awk -Fv '{print $2}' | awk -F. '{print $1$2$3}'` -lt 130 ]; then
          wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-x86_64-$version.run{,.sha256}
        else
          wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-linux-$arch-$version.run{,.sha256}
        fi
      fi
    fi
    sed -i 's/^version=.*/version='"${version}"'/g' ${BUILDER_PDIR}offline-pot/install-tke-installer.sh
    sed -i 's/^arch=.*/arch='"${arch}"'/g' ${BUILDER_PDIR}offline-pot/install-tke-installer.sh
  fi
  # build tgz
  mkdir -p ${BUILDER_PDIR}offline-pkg/
  cd ${BUILDER_PDIR}
  btime=`date "+%Y%m%d%H%M%S"`
  if [ "${builder_type}" == "custom" ] && [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep ^tkestack$ | wc -l ` -eq 0 ] && [ ${#server_set[@]} -gt 0 ]; then
    if [ `ls offline-pot-images/ | wc -l ` -gt 0 ] && [ `ls offline-pot-tgz/ | wc -l ` -eq 0 ]; then
      tar -zcf ${BUILDER_PDIR}offline-pkg/offline-pot-${app_env_flag}.${btime}.tar.gz offline-pot/ offline-pot-images/ --exclude offline-pot/hosts.tpl.bak
    elif [ `ls offline-pot-images/ | wc -l ` -eq 0 ] && [ `ls offline-pot-tgz/ | wc -l ` -gt 0 ]; then
      tar -zcf ${BUILDER_PDIR}offline-pkg/offline-pot-${app_env_flag}.${btime}.tar.gz offline-pot/ offline-pot-tgz/ --exclude offline-pot/hosts.tpl.bak
    elif [ `ls offline-pot-images/ | wc -l ` -eq 0 ] && [ `ls offline-pot-tgz/ | wc -l ` -eq 0 ]; then
      tar -zcf ${BUILDER_PDIR}offline-pkg/offline-pot-${app_env_flag}.${btime}.tar.gz offline-pot/ --exclude offline-pot/hosts.tpl.bak
    elif [ `ls offline-pot-images/ | wc -l ` -gt 0 ] && [ `ls offline-pot-tgz/ | wc -l ` -gt 0 ]; then
      tar -zcf ${BUILDER_PDIR}offline-pkg/offline-pot-${app_env_flag}.${btime}.tar.gz offline-pot/ offline-pot-images/ offline-pot-tgz/ --exclude offline-pot/hosts.tpl.bak
    fi
  elif [ "${builder_type}" == "custom" ] && [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep ^tkestack$ | wc -l ` -eq 1 ] && [ ${#server_set[@]} -eq 1 ]; then
    tar -zcf ${BUILDER_PDIR}offline-pkg/offline-pot-${app_env_flag}.${btime}.tar.gz offline-pot/  tkestack/ --exclude offline-pot/hosts.tpl.bak
  else
    tar -zcf ${BUILDER_PDIR}offline-pkg/offline-pot-${app_env_flag}.${btime}.tar.gz offline-pot/ offline-pot-images/ offline-pot-tgz/ tkestack/ --exclude offline-pot/hosts.tpl.bak
  fi
  md5sum ${BUILDER_PDIR}offline-pkg/offline-pot-${app_env_flag}.${btime}.tar.gz > ${BUILDER_PDIR}offline-pkg/offline-pot-${app_env_flag}.${btime}.md5.txt
  echo "###### build to tgz end ######"
}

# set commons deploy's switch
set_switch(){
  echo "###### commons deploy's switch sart ######"
  # backup hosts.tpl
  if [ -f "${BUILDER_PDIR}offline-pot/hosts.tpl" ]; then
    \cp ${BUILDER_PDIR}offline-pot/hosts.tpl ${BUILDER_PDIR}offline-pot/hosts.tpl.bak
  else
    echo "${BUILDER_PDIR}offline-pot/hosts.tpl not exist, please check !!!" && exit 0
  fi
  for s in ${all_servers[@]};
  do
    if [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep ^${s}$ | wc -l ` -eq 1 ]; then
      if [ "${s}" == "redis" ]; then
        sed -i "s/deploy_${s}='.*'/deploy_${s}='true'/g" ${BUILDER_PDIR}offline-pot/hosts.tpl
        sed -i "s/redis_mode='.*'/redis_mode='master-slave'/g" ${BUILDER_PDIR}offline-pot/hosts.tpl
      elif [[ "${s}" == "redis_cluster" ]]; then
        sed -i "s/deploy_redis='.*'/deploy_redis='true'/g" ${BUILDER_PDIR}offline-pot/hosts.tpl
        sed -i "s/redis_mode='.*'/redis_mode='cluster'/g" ${BUILDER_PDIR}offline-pot/hosts.tpl
      elif [[ "${s}" == "tkestack" ]]; then
        continue
      else
        sed -i "s/deploy_${s}='.*'/deploy_${s}='true'/g" ${BUILDER_PDIR}offline-pot/hosts.tpl
      fi
    else
      sed -i "s/deploy_${s}='.*'/deploy_${s}='false'/g" ${BUILDER_PDIR}offline-pot/hosts.tpl
    fi
  done
  echo "###### commons deploy's switch end ######"
}

# get builder offline resources
builder_init(){
  if [ -f "./builder-init.sh" ]; then
    sh ./builder-init.sh
  fi
}

all_func(){
  builder_init
  if [ "${builder_type}" == "all" ] || [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep ^business$ | wc -l ` -eq 1 ]; then
    get_business_helms
    if [ "x${remote_registry}" == "xfalse" ]; then
      get_business_images
    fi
  fi
  if [ "${builder_type}" == "all" ] || [ `echo ${server_set[@]} | sed 's/ /\n/g' | grep -v ^business$ | grep -v ^tkestack$ | wc -l ` -gt 0 ]; then
    get_base_component_helms
  fi
  set_switch
  build_to_tgz
}

main(){
  if [ "${builder_type}" != "all" ] && [ "${builder_type}" != "custom" ] ; then
    echo "unkwon builder type!!,builder type just support all or custom. please reconfig on builder's configfile !!!" && exit 0
  fi
  if [ "x$app_env_flag" != "x" ]; then
    $CALL_FUN || help
  else
    help
  fi
  # recover hosts.tpl
  if [ -f "${BUILDER_PDIR}offline-pot/hosts.tpl.bak" ]; then
    mv ${BUILDER_PDIR}offline-pot/hosts.tpl.bak ${BUILDER_PDIR}offline-pot/hosts.tpl
  fi
  # clean offline-pot-images and offline-pot-tgz dir
  if [ -d "${BUILDER_PDIR}offline-pot-images" ]; then
    rm -rf ${BUILDER_PDIR}offline-pot-images
  fi
  if [ -d "${BUILDER_PDIR}offline-pot-tgz" ]; then
    rm -rf ${BUILDER_PDIR}offline-pot-tgz
  fi
}
main
