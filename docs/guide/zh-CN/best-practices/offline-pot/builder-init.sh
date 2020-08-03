#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="all_func"

help(){
  echo "show usage:"
  echo "builder_init: builder init for get builder onestep offline resource: example base components images etc. "
  exit 0
}

while getopts ":f:" opt
do
  case $opt in
    f)
    CALL_FUN="${OPTARG}"
    ;;
    ?)
    echo "unkown args!just suport -f[call function]  arg!!!"
    exit 0;;
  esac
done


# get builder.sh's parent dir
BUILDER_PDIR=`echo ${BASE_DIR} | awk -Foffline-pot '{print $1}'`

builder_init(){
  echo "##### builder init start #####"
  # start onestep builder init container
  if [ `docker ps -a | grep onestep-builder-init | wc -l ` -eq 0 ]; then
    docker run --name onestep-builder-init -ti -d chenyihua/onestep-builder:v0.1  
  else
    docker rm -f onestep-builder-init
    docker run --name onestep-builder-init -ti -d chenyihua/onestep-builder:v0.1 
  fi
  # copy onestep offline resource
  docker cp onestep-builder-init:/data/base-component-helms ${BUILDER_PDIR}
  docker cp onestep-builder-init:/data/offline-pot-images-base ${BUILDER_PDIR}
  docker cp onestep-builder-init:/data/offline-pot-tgz-base ${BUILDER_PDIR}
  docker cp onestep-builder-init:/data/offline-yum ${BUILDER_PDIR}
  docker cp onestep-builder-init:/data/onestep-builder-init ${BUILDER_PDIR}
  # stop and rm onestep builder init container
  if [ `docker ps -a | grep onestep-builder-init | wc -l` -eq 1 ]; then
    docker stop onestep-builder-init &&  docker rm -f onestep-builder-init
  fi
  echo "##### builder init end #####"
}

all_func(){
   builder_init
}

main(){
  if [ ! -d "${BUILDER_PDIR}/base-component-helms" ] || [ ! -d "${BUILDER_PDIR}/offline-pot-images-base" ] || \
     [ ! -d "${BUILDER_PDIR}/offline-pot-tgz-base" ]; then  
     $CALL_FUN || help
  fi
}
main
