#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

INSTALL_DATA_DIR=/opt/tke-installer/data/

# init keepalived, is tmp script; when tkestack keepalived issue fix will be remove
init_keepalived(){
   
  for i in {1..100000}; 
  do
    if [ `cat ${INSTALL_DATA_DIR}/tke.log | grep 'EnsureKubeadmInitKubeConfigPhase' | grep -i 'Success' | wc -l` -gt 0 ]; then
      sh ./offline-pot-cmd.sh -s init-keepalived.sh -f init
      exit 0
    fi
  done
}


main(){
  init_keepalived
}
main
