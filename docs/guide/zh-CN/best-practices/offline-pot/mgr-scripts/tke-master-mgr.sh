#!/bin/bash
# Author: kubelouislu
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

DOWN_MASTER_IP=""
NEW_MASTER_IP=""
HEALTHY_MASTER_IP=""
REBIRTH_MASTER_IP=""
CALL_FUN="help"
hosts="all"

help(){
  echo "show usage:"
  echo "###### change master usage ######"
  echo "cmd: will reset the master whose IP is from the arg d and will start a master to hold on whose NEW_MASTER_IP is from the arg n also We need some info from a healthy master whose IP is from the arg h"
  echo "you can exec script by the following: "
  echo "./tke-master-mgr.sh -f change-master -d {DOWN_MASTER_IP} -n {NEW_MASTER_IP} -h {HEALTHY_MASTER_IP}"
  echo "###### rebirth master usage ######"
  echo "cmd: will rebirth a master to hold on whose IP is from the arg r and We need some info from a healthy master whose IP is from the arg h"
  echo "you can exec script by the following: "
  echo "./tke-master-mgr.sh -f rebirth-master -r {REBIRTH_MASTER_IP} -h {HEALTHY_MASTER_IP}"
  exit 0
}

while getopts ":f:d:n:h:r:" opt
do
  case $opt in
    f)
    CALL_FUN="${OPTARG}"
    ;;
    d)
    DOWN_MASTER_IP="${OPTARG}"
    ;;
    n)
    NEW_MASTER_IP="${OPTARG}"
    ;;
    h)
    HEALTHY_MASTER_IP="${OPTARG}"
    ;;
    r)
    REBIRTH_MASTER_IP="${OPTARG}"
    ;;
    ?)
    echo "unkown args! just suport -f[call function: change-master or rebirth-master] -r[IP address of the server will rebirth] -d[IP address of the master will be down] -n[IP address of the server will be up] -h[IP address of the master that is stable in cluster] arg!!!"
    exit 0;;
  esac
done

confirm_change_master(){
  if [ "$DOWN_MASTER_IP" == "" ]; then
      echo "Lack of the arg d"
      exit 1
  fi
  if [ "$NEW_MASTER_IP" == "" ]; then
      echo "Lack of the arg n"
      exit 1
  fi
  if [ "$HEALTHY_MASTER_IP" == "" ]; then
      echo "Lack of the arg h"
      exit 1
  fi
  echo "The master: ${DOWN_MASTER_IP} will be down"
  echo "The server: ${NEW_MASTER_IP} will be up"
  echo "The master: ${HEALTHY_MASTER_IP} will be reference"
  read -r -p "Are You Sure? [Y/n] " input

  case $input in
      [yY][eE][sS]|[yY])
      ;;

      [nN][oO]|[nN])
      exit 1
          ;;

      *)
      echo "Invalid input..."
      exit 1
      ;;
  esac
}

confirm_rebirth_master(){
  if [ "$REBIRTH_MASTER_IP" == "" ]; then
      echo "Lack of the arg r"
      exit 1
  fi
  if [ "$HEALTHY_MASTER_IP" == "" ]; then
      echo "Lack of the arg h"
      exit 1
  fi
  echo "The server: ${REBIRTH_MASTER_IP} will rebirth"
  echo "The master: ${HEALTHY_MASTER_IP} will be reference"
  read -r -p "Are You Sure? [Y/n] " input

  case $input in
      [yY][eE][sS]|[yY])
      ;;

      [nN][oO]|[nN])
      exit 1
          ;;

      *)
      echo "Invalid input..."
      exit 1
      ;;
  esac
}

# init env
init_env(){
  echo "###### init env start ######"
  ansible-playbook -f 10 -i ../hosts --tags tke_node_init ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  ansible-playbook -f 10 -i ../hosts --tags add_tke_node ../playbooks/tke-mgr/tke-mgr.yml \
  --extra-vars "hosts=${hosts}"
  echo "###### init env end ######"
}

# check_over
check_over(){
  echo "###### check over start ######"
  ORIGIN=`ansible -i ../hosts $HEALTHY_MASTER_IP -m command -a "kubectl get node"|grep $NEW_MASTER_IP|wc -l`
  while [ $ORIGIN -eq 0 ]; do
    ORIGIN=`ansible -i ../hosts $HEALTHY_MASTER_IP -m command -a "kubectl get node"|grep $NEW_MASTER_IP|wc -l`
    echo "###### env init doing please wait ######"
    sleep 5s
  done
  echo "###### check over end ######"
}

# init server env
init_server(){
  echo "###### init new server env start ######"
  # init the new one but just make some dirs and clean kube env if it has before
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags init_server_env ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${NEW_MASTER_IP}"
    echo "###### init new server env end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# init rebirth server env
init_rebirth_server(){
  echo "###### init rebirth server env start ######"
  # init the new one but just make some dirs and clean kube env if it has before
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags init_server_env ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${REBIRTH_MASTER_IP}"
    echo "###### init rebirth server env end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# get args from the stable one
get_args(){
  echo "###### get args start ######"
  # get args from the stable one to help join the control plane
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags get_args ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${HEALTHY_MASTER_IP}"
    echo "###### get args end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# join the control plane
join_control_plane(){
  echo "###### join control plane start ######"
  # join the control plane with the args
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags join_control_plane ../playbooks/tke-mgr/tke-mgr.yml \
    -e "HEALTHY_MASTER_IP=${HEALTHY_MASTER_IP} NEW_MASTER_IP=${NEW_MASTER_IP}" --extra-vars "hosts=${NEW_MASTER_IP}"
    echo "###### join control plane end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# join the rebirth control plane
join_rebirth_control_plane(){
  echo "###### join control plane start ######"
  # join the control plane with the args
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags join_control_plane ../playbooks/tke-mgr/tke-mgr.yml \
    -e "HEALTHY_MASTER_IP=${HEALTHY_MASTER_IP} NEW_MASTER_IP=${REBIRTH_MASTER_IP}" --extra-vars "hosts=${REBIRTH_MASTER_IP}"
    echo "###### join control plane end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# cp known tokens files
cp_known_tokens_files(){
  echo "###### copy known tokens start ######"
  # cp the conf file kubelet.service
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags cp_known_tokens_files ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${HEALTHY_MASTER_IP}"
    echo "###### copy known tokens end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# paste known tokens files
paste_known_tokens_files(){
  echo "###### paste known tokens start ######"
  # paste data files
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags paste_known_tokens_files ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${NEW_MASTER_IP}"
    echo "###### paste known tokens end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# paste known tokens files in rebirth
paste_known_tokens_files_rebirth(){
  echo "###### paste known tokens start ######"
  # paste known tokens data files
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags paste_known_tokens_files ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${REBIRTH_MASTER_IP}"
    echo "###### paste known tokens end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# cp scheduler files
cp_scheduler_files(){
  echo "###### copy scheduler start ######"
  # cp the conf file scheduler-policy-config.json
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags cp_scheduler_files ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${HEALTHY_MASTER_IP}"
    echo "###### copy scheduler end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# paste scheduler files rebirth
paste_scheduler_files_rebirth(){
  echo "###### paste scheduler start ######"
  # paste the conf file scheduler-policy-config.json
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags paste_scheduler_files ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${REBIRTH_MASTER_IP}"
    echo "###### paste scheduler end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# clean temp file in new master
clean_new_master(){
  echo "###### clean temp files in new master start ######"
  # clean some files left in new master
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags clean_new_master ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${NEW_MASTER_IP}"
    echo "###### clean temp files in new master end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# clean temp file in new master rebirth
clean_new_master_rebirth(){
  echo "###### clean temp files in new master start ######"
  # clean some files in new master
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags clean_new_master ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${REBIRTH_MASTER_IP}"
    echo "###### clean temp files in new master end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# clean temp file in installer
clean_installer(){
  echo "###### clean temp files in installer start ######"
  if [ -f "../hosts" ]; then
    rm -rf /data/tmp
    echo "###### clean temp files in installer end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# clean the old master
clean_old(){
  echo "###### clean old master start ######"
  # reset the old master
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags clean_old ../playbooks/tke-mgr/tke-mgr.yml \
    --extra-vars "hosts=${DOWN_MASTER_IP}"
    echo "###### clean old master end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

# delete the node
delete_old(){
  echo "###### delete old master start ######"
  # delete the node
  if [ -f "../hosts" ]; then
    ansible-playbook -f 10 -i ../hosts --tags delete_old ../playbooks/tke-mgr/tke-mgr.yml \
    -e "DOWN_MASTER_IP=${DOWN_MASTER_IP}" --extra-vars "hosts=${HEALTHY_MASTER_IP}"
    echo "###### delete old master end ######"
  else
      echo "hosts file not exist, please check!!!" && exit 0
  fi
}

all_over(){
  echo "###### All has been down successfully and now enjoy your kube ######"
}

change-master(){
  confirm_change_master || help
  init_env
  check_over
  init_server
  get_args
  join_control_plane
  cp_known_tokens_files
  paste_known_tokens_files
  cp_scheduler_files
  paste_scheduler_files
  clean_new_master
  clean_installer
  clean_old
  delete_old
  all_over
}

rebirth-master(){
  confirm_rebirth_master || help
  init_rebirth_server
  get_args
  join_rebirth_control_plane
  cp_known_tokens_files
  paste_known_tokens_files_rebirth
  cp_scheduler_files
  paste_scheduler_files_rebirth
  clean_new_master_rebirth
  clean_installer
  all_over
}

main(){
  $CALL_FUN || help
}
main