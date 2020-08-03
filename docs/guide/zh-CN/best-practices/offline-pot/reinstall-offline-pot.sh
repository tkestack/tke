#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

# clean cluster
clean_cluster(){
  sh ./clean-cluster.sh
}

# reinstall tke-installer
redpl_tke_installer(){
  if [ -d "../tkestack" ]; then
    sh ./install-tke-installer.sh
  fi
}

# init tke config and deploy offline-pot
dpl_offline_pot(){
  sh ./install-offline-pot.sh
}


main(){
  clean_cluster
  redpl_tke_installer
  dpl_offline_pot
}
main
