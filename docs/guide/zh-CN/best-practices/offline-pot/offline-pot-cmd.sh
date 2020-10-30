#!/bin/bash
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
