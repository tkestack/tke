#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

CALL_FUN="all_func"
hosts="all"

help(){
  echo "show usage:"
  echo "check_version: check system and kernal version"
  echo "check_nat_module: check host whether had loaded nat module"
  echo "check_disk: check disk meets requirements"
  echo "check_network: check pod network cidr whethere with host network conflict"
  echo "check_extranet: check extranet access for wx"
  echo "check_dns: check dns enable"
  echo "check_perfor: check system base performance"
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

# check system and kernal version
check_version(){
  echo "###### check system and kernal version start ######"
  ansible-playbook -f 10 -i ../hosts --tags check_system_kernal_version ../playbooks/hosts-check/hosts-check.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### check system and kernal version end ######"
}

# check host whether had loaded nat module 
check_nat_module(){
  echo "###### check host whether had loaded nat module start ######"
  ansible-playbook -f 10 -i ../hosts --tags check_nat_module ../playbooks/hosts-check/hosts-check.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### check host whether had loaded nat module end ######"
}

# check disk meets requirements
check_disk(){
  echo "###### check disk meets requirements start ######"
  # check data disk size
  ansible-playbook -f 10 -i ../hosts --tags check_data_disk_size ../playbooks/hosts-check/hosts-check.yml \
  --extra-vars "hosts=${hosts}"
  # check data dir whether create
  ansible-playbook -f 10 -i ../hosts --tags check_data_dir ../playbooks/hosts-check/hosts-check.yml \
  --extra-vars "hosts=${hosts}"
  # check ceph disk size and whether is raw device
  ansible-playbook -f 10 -i ../hosts --tags check_ceph_disk ../playbooks/hosts-check/hosts-check.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### check disk meets requirements end ######"
}

# check pod network cidr whethere with host network conflict
check_network(){
  echo "###### check pod network cidr whethere with host network conflict start ######"
  ansible-playbook -f 10 -i ../hosts --tags check_pod_network ../playbooks/hosts-check/hosts-check.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### check pod network cidr whethere with host network conflict end ######"
}

# check extranet access for wx
check_extranet(){
  echo "###### check extranet access for wx ######"
  if [ -f "../hosts" ]; then
    if [ "x`cat ../hosts | grep ^check_ex_net= | awk -F\' '{print $2}' || echo 'false'`" == "xtrue" ]; then
      ansible-playbook -f 10 -i ../hosts --tags check_access_internet ../playbooks/hosts-check/hosts-check.yml \
      --extra-vars "hosts=${hosts}"
    fi
  fi
  echo "###### check extranet access for wx end ######"
}

# check dns enable
check_dns(){
  echo "###### check dns enable start ######"
  ansible-playbook -f 10 -i ../hosts --tags check_dns ../playbooks/hosts-check/hosts-check.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### check dns enable end ######"
}

# check system base performance
check_perfor(){
  echo "###### check system base performance start ######"
  if [ -f "../hosts" ]; then
    if [ "x`cat ../hosts | grep ^stress_perfor= | awk -F\' '{print $2}' || echo 'false'`" == "xtrue" ]; then
      ansible-playbook -f 10 -i ../hosts --tags check_system_base_perfor ../playbooks/hosts-check/hosts-check.yml \
      --extra-vars "hosts=${hosts}"
    fi
  fi
  echo "###### check system base performance end ######"
}

# execute all function
all_func(){
  check_version
  check_nat_module
  check_disk
  check_network
  check_extranet
  check_dns
  check_perfor
}

main(){
  $CALL_FUN || help
}
main
