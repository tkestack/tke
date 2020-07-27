#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="cfg_salt_minion"
hosts="all"

help(){
  echo "show usage:"
  echo "cfg_salt_minion: config salt minion, -f default value is cfg_salt_minion"
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

# config salt minion
cfg_salt_minion(){
  echo "###### config salt minion start ######"
  # config salt minion
  if [ `ls ../roles/business/helms/ | grep -v README.md| grep -v charts| grep -v helmfile.d | \
    grep -v secrets | wc -l` -gt 0 ]; then
    center_array=(`ls ../roles/business/helms/ |grep -v README.md|grep -v charts|grep -v helmfile.d|grep -v secrets`)
    for c in ${!center_array[@]}; do
      ansible-playbook -f 10 -i ../hosts --tags cfg_salt ../playbooks/business/business.yml \
      --extra-vars "hosts=salt[$c] center=${center_array[$c]}";
    done
  fi
  echo "###### config salt minion end ######"
}

all_func(){
  cfg_salt_minion
}

main(){
  $CALL_FUN || help
}
main
