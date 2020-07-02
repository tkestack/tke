[TOC]

# tkestack 私有化部署最佳实践

## 背景    

随着私有化项目越来越多，简单快捷部署交付需求日益强烈。在开源协同大行其道，私有化一键部署如何合理利用开源协同力量做到用好80%，做好20% 值得思而深行而简。本文将揭秘私有化一键部署结合tkestack 实现私有化一键部署最佳实践面纱。

## 方案选型

虽然开源协同大行其道，但**合适才是最好**基本原则仍然需要贯彻执行 -- 量体裁衣；结合开发难度，开发效率，后期代码维护，组件维护，实施人员使用难度，实施人员现场修改难度等多维度进行考量。方案如下：    

![](./tkestack-pic/proposal.jpg)    

- **kubeadm+ansible:**     
  **优点：**对kubernetes版本，网络插件，docker 版本自主可控；        
    **缺点：**维护成本高，特别kubernetes升级,，无法专注于业务层面的一键部署;    
- **tkestack+ansible:**     
  **优点：**对于kubernetes的集成，维护无需关心，只需要用好tkestack也就掌握tkestack的基本原理做好应急；专注于业务层面的一键部署集成即可。通过ansible进行机器批量初始化，部署方便快捷对于实施人员要求低，可以随时现场修改;      
    **缺点：** kubernetes 版本，网络插件，docker 版本不自主可控;    
- **tkestack+operator:**      
  **优点：**私有化一键部署产品化，平台化；        
  **缺点：**operator开发成本高，客户环境多变复杂，出现问题没法现场修改;      
  
综合上述方案考虑维护kubernetes成本有点高，另外tkestack+ansible通过hooks方式进行扩展，能实现快速集成，可以专注于业务组件集成即可；ansible 入手容易，降低实施人员学习成本，并且可以随时根据现场环境随时修改适配；综合考量选择tkestack+ansible 模式。

## 需求     

**功能性需求：**        

| 功能 | 说明 | 
| :-------- | :-------- | 
|    主机初始化     |  安装前进行主机初始化，比如添加域名hosts，安装压测工具，离线yum源等   | 
|    主机检查     |  检查当前主机的性能是否符合需求，磁盘大小是否符合需求,操作系统版本，内核版本，性能压测等  |
|    tkestack部署   |  部署kuberntes和tkestack  |     
|    业务依赖组件部署   |  部署业务依赖组件，比如redis，mysql，部署运维组件elk，prometheus等 |  
|    业务部署   |  部署业务服务  |  


**非功能性需求：**        

| 功能 | 说明 | 
| :-------- | :-------- | 
|    解耦     |  针对一些已有kubernetes/tke 平台，此时需要只部署业务依赖组件及业务，所以需要和tkestack解耦   | 
|    扩展性    |  业务依赖组件不同项目需要采用不同的依赖组件，需要快捷集成新的组件  |
|    幂等   |  部署及卸载时可以重复执行  |     



## 实现    

### 1. 走进tkesack

![](./tkestack-pic/tkestack-arti.png)    

从tkestack git 获取到的架构图可以看出tkestack分为installer， Global，cluster这三种角色；其中installer 负责tkestack Global集群的安装，当前提供命令行安装模式和图形化安装模式；cluster 角色是作为业务集群，通Global集群纳管。当前我们只需要部署一个Global集群作为业务集群即可满足需求，cluster集群只是为了提供给客户使用tkestack多集群管理使用。    

tkestack 在installer 以hooks 方式实现用户自定义扩展， 有如下hook脚本：    

-   pre-installer: 主要集群部署前的一些自定义初始化操作    
-   post-cluster-ready: 种子集群ready后针对tkestack 部署前的初始化操作    
-   post-install: tkestack 部署完毕，部署自定义扩展

默认tkestack部署流程如下：    

![](./tkestack-pic/install-proccess.jpg)    

由于installer 节点在tkestack 设计上计划安装完毕直接废弃，所以tkestack 会在global集群重新部署一个镜像仓库作为后续业务使用，当然也会将tkestack 平台的镜像重新repush到集群内的镜像仓库。所以容器化的自定义扩展主件的部署需要放到post-install 脚本进行触发。    

### 2. 魔改部署配置    

-  tkestack git 使用手册给出了两种部署模式一种是web页面配置模式，一种是命令行模式；经过使用发现tkestack有个亮点特性就是配置文件记录了一个step 安装步骤，可以在安装失败后解决问失败原因直接重启tke-installer 即可根据当前step 步骤继续进行安装部署；我们利用这特性实现web页面配置模式也可以命令行模式部署。具体操作是先通过页面配置得到配置文件，把配置文件做成模板; 部署时候通过ansible templet 模块进行渲染。 当前抽取出来配置模板有:    
  - tke-ha-lb.json.j2 对应web页面的使用已有，也就是采用负载均衡ip地址作为tkestack集群高可用
  - tke-ha-keepalived.json.j2 对应web页面的TKE提供，采用vip通过keepalived 浮动漂移实现高可用
  - tke-sigle.json.j2 对应web页面的不设置场景， 也就是单master版场景
  以下以tke-ha-lb.json.j2 为例：    
  
  ```    
  {
 "config": {
  "ServerName": "tke-installer",
  "ListenAddr": ":8080",
  "NoUI": false,
  "Config": "conf/tke.json",
  "Force": false,
  "SyncProjectsWithNamespaces": false,
  "Replicas": {{ tke_replicas }}
 },
 "para": {
  "cluster": {
   "kind": "Cluster",
   "apiVersion": "platform.tkestack.io/v1",
   "metadata": {
    "name": "global"
   },
   "spec": {
    "finalizers": [
     "cluster"
    ],
    "tenantID": "default",
    "displayName": "TKE",
    "type": "Baremetal",
    "version": "{{ k8s_version }}",
    "networkDevice": "{{ net_interface }}",
    "clusterCIDR": "{{ cluster_cidr }}",
    "dnsDomain": "cluster.local",
    "features": {
     "ipvs": {{ ipvs }},
     "enableMasterSchedule": true,
     "ha": {
      "thirdParty": {
       "vip": "{{ tke_vip }}",
       "vport": {{ tke_vport}}
      }
     }
    },
    "properties": {
     "maxClusterServiceNum": {{ max_cluster_service_num }},
     "maxNodePodNum": {{ max_node_pod_num }}
    },
    "machines": [
      {
       "ip": "{{ groups['masters'][0] }}",
       "port": {{ ansible_port }},
       "username": "{{ ansible_ssh_user }}",
       "password": "{{ ansible_ssh_pass_base64 }}"
      },
      {
       "ip": "{{ groups['masters'][1] }}",
       "port": {{ ansible_port }},
       "username": "{{ ansible_ssh_user }}",
       "password": "{{ ansible_ssh_pass_base64 }}"
      },
      {
        "ip": "{{ groups['masters'][2] }}",
        "port": {{ ansible_port }},
        "username": "{{ ansible_ssh_user }}",
        "password": "{{ ansible_ssh_pass_base64 }}"
      }
    ],
    "dockerExtraArgs": {
     "data-root": "{{ docker_data_root }}"
    },
    "kubeletExtraArgs": {
     "root-dir": "{{ kubelet_root_dir }}"
    },
    "apiServerExtraArgs": {
     "runtime-config": "apps/v1beta1=true,apps/v1beta2=true,extensions/v1beta1/daemonsets=true,extensions/v1beta1/deployments=true,extensions/v1beta1/replicasets=true,extensions/v1beta1/networkpolicies=true,extensions/v1beta1/podsecuritypolicies=true"
    }
   }
  },
  "Config": {
   "basic": {
    "username": "{{ tke_admin_user }}",
    "password": "{{ tke_pwd_base64 }}"
   },
   "auth": {
    "tke": {
     "tenantID": "default",
     "username": "{{ tke_admin_user }}",
     "password": "{{ tke_pwd_base64 }}"
    }
   },
   "registry": {
    "tke": {
     "domain": "{{ tke_registry_domain }}",
     "namespace": "library",
     "username": "{{ tke_admin_user }}",
     "password": "{{ tke_pwd_base64 }}"
    }
   },
   "business": {},
   "monitor": {
    "influxDB": {
     "local": {}
    }
   },
   "ha": {
    "thirdParty": {
      "vip": "{{ tke_vip }}",
      "vport": {{ tke_vport}}
    }
   },
   "gateway": {
    "domain": "{{ tke_console_domain }}",
    "cert": {
     "selfSigned": {}
    }
   }
  }
 },
 "cluster": {
  "kind": "Cluster",
  "apiVersion": "platform.tkestack.io/v1",
  "metadata": {
   "name": "global"
  },
  "spec": {
   "finalizers": [
    "cluster"
   ],
   "tenantID": "default",
   "displayName": "TKE",
   "type": "Baremetal",
   "version": "{{ k8s_version }}",
   "networkDevice": "{{ net_interface }}",
   "clusterCIDR": "{{ cluster_cidr }}",
   "dnsDomain": "cluster.local",
   "features": {
    "ipvs": {{ ipvs }},
    "enableMasterSchedule": true,
    "ha": {
     "thirdParty": {
       "vip": "{{ tke_vip }}",
       "vport": {{ tke_vport}}
     }
    }
   },
   "properties": {
    "maxClusterServiceNum": {{ max_cluster_service_num }},
    "maxNodePodNum": {{ max_node_pod_num }}
   },
   "machines": [
     {
       "ip": "{{ groups['masters'][0] }}",
       "port": {{ ansible_port }},
       "username": "{{ ansible_ssh_user }}",
       "password": "{{ ansible_ssh_pass_base64 }}"
     },
     {
       "ip": "{{ groups['masters'][1] }}",
       "port": {{ ansible_port }},
       "username": "{{ ansible_ssh_user }}",
       "password": "{{ ansible_ssh_pass_base64 }}"
     },
     {
       "ip": "{{ groups['masters'][2] }}",
       "port": {{ ansible_port }},
       "username": "{{ ansible_ssh_user }}",
       "password": "{{ ansible_ssh_pass_base64 }}"
     }
   ],
   "dockerExtraArgs": {
    "data-root": "{{ docker_data_root }}"
   },
   "kubeletExtraArgs": {
    "root-dir": "{{ kubelet_root_dir }}"
   },
   "apiServerExtraArgs": {
    "runtime-config": "apps/v1beta1=true,apps/v1beta2=true,extensions/v1beta1/daemonsets=true,extensions/v1beta1/deployments=true,extensions/v1beta1/replicasets=true,extensions/v1beta1/networkpolicies=true,extensions/v1beta1/podsecuritypolicies=true"
   }
  }
 },
 "step": 0 # 重启tke-installer 后会按此步骤执行继续的安装，当前设置为0意味着从零开始
}
  
  ```    
  
  为了实现此方式安装，我们的安装脚本如下：    
  
  ```    
  #!/bin/bash
  # Author: yhchen
  set -e

  BASE_DIR=$(cd `dirname $0` && pwd)
  cd $BASE_DIR

  # get offline-pot parent dir
  OFFLINE_POT_PDIR=`echo ${BASE_DIR} | awk -Foffline-pot '{print $1}'`

  INSTALL_DIR=/opt/tke-installer
  DATA_DIR=${INSTALL_DIR}/data
  HOOKS=${OFFLINE_POT_PDIR}offline-pot
  IMAGES_DIR="${OFFLINE_POT_PDIR}offline-pot-images"
  TGZ_DIR="${OFFLINE_POT_PDIR}offline-pot-tgz"
  REPORTS_DIR="${OFFLINE_POT_PDIR}perfor-reports"
  version=v1.2.4

  init_tke_installer(){
    if [ `docker images | grep tke-installer | grep ${version} | wc -l` -eq 0 ]; then
      if [ `docker ps -a | grep tke-installer | wc -l` -gt 0 ]; then
        docker rm -f tke-installer
      fi
      if [ `docker images | grep tke-installer | wc -l` -gt 0 ]; then
        docker rmi -f `docker images | grep tke-installer | awk '{print $3}'`
      fi 
      cd ${OFFLINE_POT_PDIR}tkestack
      if [ -d "${OFFLINE_POT_PDIR}tkestack/tke-installer-x86_64-${version}.run.tmp" ]; then
        rm -rf ${OFFLINE_POT_PDIR}tkestack/tke-installer-x86_64-${version}.run.tmp
      fi
      sha256sum --check --status tke-installer-x86_64-$version.run.sha256 && \
      chmod +x tke-installer-x86_64-$version.run && ./tke-installer-x86_64-$version.run
    fi
  }

  reinstall_tke_installer(){
    if [ -d "${REPORTS_DIR}" ]; then
      mkdir -p ${REPORTS_DIR}
    fi
    if [ `docker ps -a | grep tke-installer | wc -l` -eq 1 ]; then
      docker rm -f tke-installer
      rm -rf /opt/tke-installer/data
    fi
    docker run --restart=always --name tke-installer -d --privileged --net=host -v/etc/hosts:/app/hosts \
    -v/etc/docker:/etc/docker -v/var/run/docker.sock:/var/run/docker.sock -v$DATA_DIR:/app/data \
    -v$INSTALL_DIR/conf:/app/conf -v$HOOKS:/app/hooks -v$IMAGES_DIR:${IMAGES_DIR} -v${TGZ_DIR}:${TGZ_DIR} \
    -v${REPORTS_DIR}:${REPORTS_DIR} tkestack/tke-installer:$version
    if [ -f "hosts" ]; then
      # set hosts file's dpl_dir variable
      sed -i 's#^dpl_dir=.*#dpl_dir=\"'"${HOOKS}"'\"#g' hosts
      installer_ip=`cat hosts | grep -A 1 '\[installer\]' | grep -v installer`
      echo "please exec install-offline-pot.sh or access http://${installer_ip}:8080 to install offline-pot"
    fi
  }

  main(){
    init_tke_installer # 此函数是为了实现当前节点尚未安装过tke-installer, 进行第一次安装实现初始化
    reinstall_tke_installer # 此函数是实现自定义安装tke-installer, 主要是为了将扩展的hooks脚本挂载到tke-installer，以及hooks脚本调用到的整个一键部署脚本。
  }
  main
  
  ```    
  
  最终实现开始部署tkestack脚本如下：    
  
  ```    
  #!/bin/bash
  # Author: yhchen
  set -e

  BASE_DIR=$(cd `dirname $0` && pwd)
  cd $BASE_DIR

  CALL_FUN="defaut"

  help(){
    echo "show usage:"
    echo "init_and_check: will be init hosts, inistall tke-installer and hosts check"
    echo "dpl_offline_pot: init tke config and deploy offline-pot"
    echo "init_keepalived: just tmp use, when tkestack fix keepalived issue will be remove"
    echo "only_install_tkestack: if you want only install tkestack, please -f parameter pass only_install_tkestack"
    echo "defualt: will be exec dpl_offline_pot and init_keepalived"
    echo "all_func: execute init_and_check, dpl_offline_pot, init_keepalived"
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

  INSTALL_DATA_DIR=/opt/tke-installer/data/

  init_and_check(){
    sh ./init-and-check.sh
  }

  # init tke config and deploy offline-pot
  dpl_offline_pot(){
    echo "###### deploy offline-pot start ######"
    if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
      # deploy tkestack , base commons and business
      sh ./offline-pot-cmd.sh -s init-tke-config.sh -f init
      docker restart tke-installer
      if [ -f "hosts" ]; then
        installer_ip=`cat hosts | grep -A 1 '\[installer\]' | grep -v installer`
        echo "please exec tail -f ${INSTALL_DATA_DIR}/tke.log or access http://${installer_ip}:8080 check install progress..."
      fi
    elif [ ! -d "../tkestack" ]; then
      # deploy base commons and business on other kubernetes plat
      sh ./post-install
    else
      echo "if first install,please exec init-and-check.sh script, else exec reinstall-offline-pot.sh script" && exit 0
    fi
    echo "###### deploy offline-pot end ######"
  }

  # just tmp use, when tkestack fix keepalived issue will be remove
  init_keepalived(){
    echo "###### init keepalived start  ######"
    if [ -f "${INSTALL_DATA_DIR}/tke.json" ]; then
      if [ `cat ${INSTALL_DATA_DIR}/tke.json | grep -i '"ha"' | wc -l` -gt 0 ]; then
        nohup sh ./init_keepalived.sh 2>&1 > ${INSTALL_DATA_DIR}/dpl-keepalived.log &
      fi
    fi
    echo "###### init keepalived end ######"
  }

  # only install tkestack
  only_install_tkestack(){
    echo "###### install tkestack start ######"
    # change tke components's replicas number
    if [ -f "hosts" ]; then 
      sed -i 's/tke_replicas="1"/tke_replicas="2"/g' hosts
    fi
    # hosts init
    if [ `docker ps | grep tke-installer | wc -l` -eq 1 ]; then
      sh ./offline-pot-cmd.sh -s host-init.sh -f sshd_init
      sh ./offline-pot-cmd.sh -s host-init.sh -f selinux_init
      sh ./offline-pot-cmd.sh -s host-init.sh -f remove_devnet_proxy
      sh ./offline-pot-cmd.sh -s host-init.sh -f add_domains
      sh ./offline-pot-cmd.sh -s host-init.sh -f data_disk_init
      sh ./offline-pot-cmd.sh -s host-init.sh -f check_iptables
    else
      echo "please exec install-tke-installer.sh to start tke-installer" && exit 0
    fi
    # start install tkestack
    dpl_offline_pot
    init_keepalived
    echo "###### install tkestack end ######"
  }

  defaut(){
    # change tke components's replicas number
    if [ -f "hosts" ]; then 
      sed -i 's/tke_replicas="2"/tke_replicas="1"/g' hosts
    fi
    # only deploy tkestack
    if [ -d '../tkestack' ] && [ ! -d "../offline-pot-images" ] && [ ! -d "../offline-pot-tgz" ]; then
      only_install_tkestack
    fi
    dpl_offline_pot
    # when deploy tkestack will be init keepalived config
    if [ -d '../tkestack' ]; then
      init_keepalived
    fi
  }

  all_func(){
    # change tke components's replicas number
    if [ -f "hosts" ]; then 
      sed -i 's/tke_replicas="2"/tke_replicas="1"/g' hosts
    fi
    init_and_check
    defaut
  }

  main(){
    $CALL_FUN || help
  }
  main
  
  ```    
  
  此脚本主要是判断当前部署是否需要部署tkestack或者是否单独部署tkestack，若是部署tkestack则生成tkestack 所需的配置文件,然后通过docker restart tke-installer 即可出发tkestack部署以及业务依赖组件，业务部署。    
  
- 添加worker节点    
  
-  增加自定义参数使集群更稳，更强。主要增加自定义参数如下：    
    
    ```    
    1. dockerExtraArgs data-root 制定docker 目录到数据盘，避免系统盘太小导致节点磁盘使用率很快到达节点压力阈值以至于节点处于not ready状态
    2. kubeletExtraArgs kubelete自定义参数 root-dir 和docker data-root 参数作用一致
    3. kubeletExtraArgs  kube-apiserver runtime-config apps/v1beta1=true,apps/v1beta2=true,extensions/v1beta1/daemonsets=true,extensions/v1beta1/deployments=true,extensions/v1beta1/replicasets=true,extensions/v1beta1/networkpolicies=true,extensions/v1beta1/podsecuritypolicies=true 增加工作负载deployment
，statefulset的api version兼容性
    
    ```    
     
当前通过ansible set facts 方式，ansible when 条件执行，以及shell 命令增加判断方式实现幂等；通过设置开关+hooks+ansible tag方式实现扩展性和解耦。      
最终私有化一键部署流程如下：    

![](./tkestack-pic/onstep-install-process.jpg)    

合理利用tkestack特性（用好80%），结合自身业务场景做出满足需求私有化一键部署(做好20%)。



## 不足

- 在测试tkestack keepalived 采用单播模式的高可用方案时出现部署kubernetes过程中由于keepalived 选举切换发生vip 网络抖动最终导致部署失败；
- 当前所有镜像都是打成tar附件模式打包安装包，使得安装包有点大；同时部署集群镜像仓库时还需要从installer节点的镜像仓库重新将镜像推送至集群镜像仓库，这个耗时很大；建议将出包时将镜像推送到离线镜像仓库，然后将离线镜像仓库持久化目录打包这样合理利用镜像特性缩减安装包大小；部署时拷贝镜像仓库持久化数据到对应目录并挂载，加速部署。