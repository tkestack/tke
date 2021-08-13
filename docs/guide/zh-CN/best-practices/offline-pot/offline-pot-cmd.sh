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

EXEC_SCRIPT=""
CALL_FUN="all_func"
hosts="all"

help(){
  echo "show usage:"
  echo "cmd: will be docker exec tke-installer mgr-scripts's script"
  echo "you can exec script list: "
  echo `ls mgr-scripts/ | grep -v ansible.cfg`
  exit 0
}

while getopts ":s:f:h:" opt
do
  case $opt in
    s)
    EXEC_SCRIPT="${OPTARG}"
    ;;
    f)
    CALL_FUN="${OPTARG}"
    ;;
    h)
    hosts="${OPTARG}"
    ;;
    ?)
    echo "unkown args! just suport -s[mgr-scripts's script] -f[call function] and -h[ansible hosts group] arg!!!"
    exit 0;;
  esac
done

# docker exec tke-installer scripts
cmd(){
  echo "###### exec ${EXEC_SCRIPT} ${CALL_FUN} start ######"
  if [ -f "mgr-scripts/${EXEC_SCRIPT}" ] && [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
    docker exec tke-installer /bin/bash -c "/app/hooks/mgr-scripts/${EXEC_SCRIPT} -f ${CALL_FUN} -h ${hosts}"
  else
    echo "mgr-scripts/${EXEC_SCRIPT} not exist or tke-installer not start, please check!!!"
  fi
  echo "###### exec ${EXEC_SCRIPT} ${CALL_FUN} end ######"
}

main(){
  if [ "x${EXEC_SCRIPT}" == "x" ]; then
    help
  else
    cmd
  fi
}
main
