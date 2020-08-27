[TOC]

# offline-pot 部署手册     


## 主机规划    

**当前offline-pot 命令行部署方式仅支持三master高可用部署(带clb的vip或者float ip：keepalived)以及单master部署。主机规划如下：**    

高可用部署规划：    

| 角色 |  数量  | 主机配置 |  备注 |
| :-----| :---- | :---- |  :---- |
|  installer | 1 | 8C16G 100G: /root 200G: /data |  不能和master，node复用  |
| masters | 3 | 8C16G 100G: /root 500G: /data |  必须：lb绑定三master节点的6443端口或提供float ip(vip)给keepalived  |
| workers | >= 0 | 8C16G 100G: /root 500G: /data | 数量根据业务需要添加, 主机配置根据业务需要调整,但必须包含系统盘和数据盘 |       

单机部署规划：    

| 角色 |  数量  | 主机配置 |  备注 |
| :-----| :---- | :---- |  :---- |
|  installer | 1 | 8C16G 100G: /root 200G: /data |  不能和master，node复用  |
| masters | 1 | 8C16G 100G: /root 500G: /data |    |
| workers | >= 0 | 8C16G 100G: /root 500G: /data | 数量根据业务需要添加, 主机配置根据业务需要调整,但必须包含系统盘和数据盘 |       

其他资源：    
kubernetes 高可用 1 x lb/float ip(vip): 当高可用方案采用lb时， 需要lb tcp模式绑定到所有master机器的6443,;当nginx-ingress 部署到master主机组的第二，第三个master节点时需要将lb的80,443端口绑定到master主机组的第二，第三个master节点节点上，这时第一个master一定不要绑定lb的80，443端口除非对于nginx-ingress有独立的lb另当别论；当高可用方案采用keepalived时采用float ip,**推荐使用lb模式的高可用方案**。针对es，sgikes 等额外持久化数据的建议单独规划磁盘空间，最好单独一块磁盘以提升性能。    



   

## 部署目录说明        

```
/data
|-- offline-pot
|-- offline-pot-images
|-- offline-pot-tgz
|-- perfor-reports
`-- tkestack
```

- **offline-pot** 一键部署安装脚本目录
- **offline-pot-images** 组件、业务镜像目录
- **offline-pot-tgz** 组件（helm、redis、mysql）目录
- **perfor-reports** 压测报告目录
- **tkestack** Tkestack安装包目录    


### 1. offline-pot 当前目录说明

```
offline-pot
|-- ansible.cfg # ansible 配置文件，主要包含ansible roles 目录，ssh时是否检验主机key等ansible相关的配置,不用关心
|-- builder.cfg.tpl # 打包配置模板文件
|-- builder.sh # 打包脚本
|-- clean-cluster.sh # 清理tke集群脚本；包含清理db，redis，安装机除外;会直接删除/data/目录下的所有子目录,注意备份

|-- docs # 文档存放目录
|-- hosts.tpl # oneStep部署配置模板，部署时需要cp hosts.tpl hosts
|-- init-and-check.sh # 主机初始化以及检查主机是否符合要求脚本
|-- init_keepalived.sh # 初始化keepalived 配置脚本，预留
|-- install-offline-pot.sh # 一键部署offline-pot 脚本; tke，db，redis，日志监控，业务部署,或只部署tkestack集群
|-- install-tke-installer.sh # 部署tke-installer 脚本
|-- mgr-scripts # 管理脚本集目录
|-- offline-pot-cmd.sh # 但采用tkestack 方式部署时，用来在安装机执行mgr-scripts目录下的脚本（单项执行）
|-- playbooks # ansible playbooks
|-- post-cluster-ready # tkestack 集群ready hook 脚本, 预留
|-- post-install # tke 已就绪，tkestack自定义基础组件以及业务部署的hook触发脚本; 基础组件，业务部署执行入口
|-- pre-install # tkestack 安装前初始化脚本，预留
|-- reinstall-offline-pot.sh # offline-pot重装脚本；集群清理+offline-pot部署脚本
`-- roles # ansible roles

```    

### 2. mgr-scripts 目录说明    

```
mgr-scripts/
|-- ansible.cfg # ansible 配置文件，主要包含ansible roles 目录，ssh时是否检验主机key等ansible相关的配置,不用关心
|-- cfg-salt-minion.sh # 配置salt minion 脚本
|-- clean-cluster-nodes.sh # 清理k8s 节点脚本，会删除/data下的所有子目录，注意备份
|-- deploy-base-component.sh # 部署基础组件脚本，elk，kafka，helms，镜像加载，prometheus，nfs等
|-- deploy-business.sh # 部署业务脚本
|-- health-check.sh # 健康检查脚本
|-- host-init.sh # 主机初始化脚本
|-- hosts-check.sh # 主机合规检查脚本
|-- init-keepalived.sh # 初始化keepalived 配置脚本，预留
|-- init-tke-config.sh # 初始化tkestack 配置脚本
|-- operation-undo.sh # 对于sshd ，docker代理设置的回滚操作
|-- remove-base-component.sh # 卸载基础组件
|-- remove-business.sh # 卸载业务组件
|-- tke-gateway-mgr.sh # 调整tkestack gateway 副本数
|-- tke-mgr.sh # tkestack 管理脚本
|-- tke-nodes-mgr.sh # tkestack 添加/移除节点脚本
`-- update-kernal.sh # 更新centos 内核脚本

```

### 3. playbooks 目录说明    

```
playbooks/
|-- base-component # 基础组件playbooks
|   `-- base-component.yml
|-- business  # 业务playbooks
|   `-- business.yml
|-- hosts-check # 健康检查playbooks
|   `-- hosts-check.yml
|-- hosts-init # 主机初始化playbooks
|   `-- hosts-init.yml
|-- operation-undo # ssh,docker 代理操作回滚playbooks
|   `-- operation-undo.yml
`-- tke-mgr # tkestack playbooks
    `-- tke-mgr.yml

```    

### 4. roles 目录说明    

#### 1. 主机初始化roles 说明    

```
roles/hosts-init/
|-- README.md
|-- defaults
|   `-- main.yml
|-- handlers
|   `-- main.yml # 主机服务重启，启动，停止 处理
|-- meta
|   `-- main.yml
|-- scripts
|   |-- cdr2mask.sh # crd 转 掩码 脚本
|   `-- mask2cdr.sh # 掩码 转 crd 
|-- tasks
|   |-- add-domains.yml # 添加自定义域名解析
|   |-- check-iptables.yml # 检查是否开启iptables/firewalld, 若开启则增加对应防火墙规则
|   |-- check-time-syn-service.yml # 检查是否已配置了时间同步，若没有则以安装机作为时间同步服务器。
|   |-- data-disk-init.yml # 磁盘初始化， 需要当前节点已经安装了lvm
|   |-- deploy-offlie-yum-repo.yml # 部署离线yum repo
|   |-- install-base-tools.yml # 安装基础工具，比如helm，helm2， helmfile，nettools等
|   |-- install-stress-tools.yml # 安装压测工具
|   |-- main.yml # 总的task 编排
|   |-- registry-influxdb-init.yml # tkestack 镜像仓库以及influxdb 持久化目录软链接处理到数据盘
|   |-- remove-devnetcloud-proxy.yml # 移除devnetcloud 代理
|   |-- selinux-init.yml # 禁用seliunx
|   |-- sshd-init.yml # 开启ssh 密码认证能力
|   `-- update-kernal.yml # 更新centos7 内核
|-- templates
|   |-- ntp.conf.client.j2 # ntp 客户端模板文件
|   |-- ntp.conf.server.j2 # ntp 服务端模板文件
|   `-- offline-yum.repo.j2 # 离线yum repo 配置模板文件
|-- tests
|   |-- inventory
|   `-- test.yml
`-- vars
    `-- main.yml

```

### 2. 主机检查roles 说明    

```
roles/hosts-check/
|-- README.md
|-- defaults
|   `-- main.yml
|-- handlers
|   `-- main.yml 
|-- meta
|   `-- main.yml
|-- tasks
|   |-- check-disk-meets-requirements.yml # 检查磁盘空间大小是否满足需求
|   |-- check-dns-enable.yml # 检查节点是否已配置dns
|   |-- check-extranet-access.yml # 检查是否可以访问外网服务
|   |-- check-nat-module.yml # 检查节点是否安装nat 模块
|   |-- check-perfor.yml # 性能检查, 磁盘和网络
|   |-- check-pod-network-cidr.yml # 检查pod 的网段是否和主机网段冲突
|   |-- check-system-and-kernal-version.yml # 检查当前系统版本以及内核版本是否过低
|   `-- main.yml # 总的task 编排
|-- tests
|   |-- inventory
|   `-- test.yml
`-- vars
    `-- main.yml # 预定义系统版本，内核版本变量

```

### 3. 回滚操作roles 说明    

```
roles/operation-undo/
|-- README.md
|-- defaults
|   `-- main.yml
|-- handlers
|   `-- main.yml #  主机服务重启，启动，停止 处理
|-- meta
|   `-- main.yml
|-- tasks
|   |-- main.yml # 总的task 编排
|   |-- recover-disk-to-raw.yml # 将磁盘恢复成裸设备
|   |-- remove-devnetcloud-proxy-undo.yml # 恢复devnetcloud 代理
|   `-- sshd-init-undo.yml # 恢复修改前ssh 配置
|-- tests
|   |-- inventory
|   `-- test.yml
`-- vars
    `-- main.yml

```

### 4. tkestack roles 说明    

```
roles/tke-mgr/
|-- README.md
|-- defaults
|   `-- main.yml
|-- files
|-- handlers
|   `-- main.yml
|-- meta
|   `-- main.yml
|-- scripts
|   |-- clean-cluster.sh # 清理tkestack 节点(集群清理时调用)
|   `-- clean-nodes.sh # 清理tkestack 节点(单节点清理时调用)
|-- tasks
|   |-- init-keepalived.yml # 初始化keepalived 配置，预留
|   |-- init-tke-cfg.yml # 初始化tkestack 配置
|   |-- main.yml # 总的task 编排
|   |-- manager-tke-node.yml # 添加/移除tkestack节点
|   |-- remove-cluster.yml # 移除tkestack集群
|   `-- tke-gateway-mgr.yml # 调整tkestack gateway 副本数
|-- templates
|   |-- keepalived.conf.j2 # keepalvied 配置模板
|   |-- tke-ha-keepalived.json.j2 # tkestack keepalived 高可用配置模板(对应页面的tkestack支持)
|   |-- tke-ha-lb.json.j2 # tkestack lb 高可用配置模版(对应页面的已有)
|   |-- tke-node.yaml.j2 # 添加tkestack 节点模板
|   `-- tke-sigle.json.j2 # tkestack 单master 配置模板
|-- tests
|   |-- inventory
|   `-- test.yml
`-- vars
    `-- main.yml # ansible_ssh_pass_base64 变量定义

```

### 5. 基础组件 roles 说明     

```
roles/base-component/
|-- README.md
|-- defaults
|   `-- main.yml
|-- files
|-- handlers
|   `-- main.yml
|-- helms # 打包后基础组件helm chart 保存位置
|-- meta
|   `-- main.yml
|-- tasks
|   |-- elk-mgr.yml # elk 部署/卸载
 |      | --  harbor-mgr.yml # harbor 部署/卸载
|   |-- helm-tiller-mgr.yml # helm tiller 部署/卸载
|   |-- kafka-mgr.yml # kafka 部署/卸载
|   |-- load-and-push-images.yml # 加载和推送镜像(自定义)
|   |-- main.yml # 总的task 编排
|   |-- minio-mgr.yml # minio 部署/卸载
|   |-- mysql-manager.yml # mysql 部署/卸载
|   |-- nfs-mgr.yml # nfs 部署
|   |-- nginx-ingress-mgr.yml # nginx ingress controller 部署/卸载
|   |-- postgres-mgr.yml # postgres 部署/卸载
|   |-- prometheus-mgr.yml # prometheus 部署/卸载
|   |-- redis-manager.yml # redis 部署/卸载
|   `-- sgikes-mgr.yml # 分词搜索服务部署/卸载
|-- templates
|   |-- common # 通用模板存放目录
|   |   |-- base-component-tools.sh.j2 # 所有基础组件部署卸载脚本，建议拆分到各自模块实现部署卸载
|   |   `-- local-storage.yaml.j2 # 创建local storage 模板文件
|   |-- elk
|   |   |-- elasticsearch.yaml.j2 # elasticsearch helm value 模板
|   |   |-- es-base-auth-secret.yaml.j2 # es 认证secret 模板文件
|   |   |-- es-crt-secret.yaml.j2 # es 证书secret 模板文件
|   |   |-- es-local-pv.yaml.j2 # es local pv 模板
|   |   |-- filebeat-k8s.yaml.j2 # filebeat 部署模板
|   |   |-- kibana-cert-secret.yaml.j2 # kibana 证书secret 模板
|   |   |-- kibana.tke.com-secret.yaml.j2 # kibana.tke.com 域名secret模板
|   |   |-- kibana.yaml.j2 # kibana helm value 模板
|   |   |-- log-test.yaml.j2 # 日志采集测试服务部署模板
|   |   |-- logstash.yaml.j2 # logstash 部署模板
|   |   `-- telnet-curl-tcpdump.yaml.j2 # telnet curl tpcdump 调试服务部署模板
 |      | -- harbor
 |      |   `-- harbor-hosts.j2 # harbor 配置模板
|   |-- helmtiller
|   |   |-- helm-tiller.yaml.j2 # tkestack 时helmtiller 部署模板
|   |   |-- helmtiller-clusterrole.yaml.j2 # 非tkestack helmtiller cluster role
|   |   |-- helmtiller-deploy.yaml.j2 # 非tkestack helmtiller 部署模板
|   |   `-- helmtiller-sa.yaml.j2 # 非tkestack helmtiller service account
|   |-- kafka
|   |   |-- kafka-local-pv.yaml.j2 # kafka local pv 模板
|   |   |-- kafka_manager.yaml.j2 # kafka manager helm value 模板
|   |   |-- kafka_zk.yaml.j2 # kafka,zk helm value 模板
|   |   `-- zk-local-pv.yaml.j2 # zookeeper local pv 模板
|   |-- minio
|   |   |-- minio-dpl.yaml.j2 # minio helm value 模板
|   |   `-- minio-local-pv.yaml.j2 # minio local pv 模板
|   |-- mysql
|   |   |-- mysql_master_deploy.sh.j2 # mysql master 部署脚本模板
|   |   |-- mysql_slave_deploy.sh.j2 # mysql slave 部署脚本模板
|   |   `-- remove_mysql.sh.j2 # 卸载mysql 脚本模板
|   |-- nfs
|   |   |-- nfs-mgr-tools.sh.j2 # nfs 启动，配置脚本模板
|   |   `-- nfs-pv.yaml.tpl.j2 # 创建nfs pv 模板
|   |-- nginx_ingress
|   |   `-- nginx-ingess.yaml.j2 # nginx ingress controller helm value 模板
|   |-- postgres
|   |   `-- postgres_start.sh.j2 # postgres 启动模板
|   |-- prometheus
|   |   `-- prometheus.platform.tkestack.io.yaml.j2 # tkestack prometheus 部署模板
|   |-- redis
|   |   |-- clean_redis.sh.j2 # 清理主从模式redis 脚本
|   |   |-- redis-cluster-client.yaml.j2 # redis 集群模式 client 部署 模板
|   |   |-- redis-cluster-local-pv.yaml.j2 # redis local pv 模板
|   |   |-- redis-cluster-values.yaml.j2 # redis 集群模式 helm value 模板
|   |   |-- redis_master_deploy.sh.j2 # redis master 部署脚本模板
|   |   `-- redis_slave_deploy.sh.j2 # redis slave 部署脚本模板
|   `-- sgikes
|       |-- sg-ik-es-data-pv-isolate.yaml.j2 # 分词搜索服务数据节点 pv 隔离模板
|       |-- sg-ik-es-data-pv.yaml.j2 # 分词搜索服务数据节点pv 模板
|       |-- sg-ik-es-master-pv.yaml.j2 # 分词搜索服务master pv 模板
|       `-- sg-ik-values.yaml.j2 # 分词搜索服务 helm value 模板
|-- tests
|   |-- inventory
|   `-- test.yml
`-- vars
    `-- main.yml

```    

### 6. 业务 roles 说明     

```
roles/business/
|-- README.md
|-- defaults
|   `-- main.yml
|-- files
|-- handlers
|   `-- main.yml
|-- helms # 打包后业务helm chart 存放目录
|-- meta
|   `-- main.yml
|-- tasks
|   |-- business-mgr.yml # 业务组件部署/卸载 
|   `-- main.yml
|-- templates
|   |-- business-tools.sh.j2 # 业务部署/卸载脚本模板，当前仅支持helmfile方式部署，其他请自行在此扩展；注意条件式执行
|   `-- minion.cfg.j2 # salt minion 配置模板
|-- tests
|   |-- inventory
|   `-- test.yml
`-- vars
    `-- main.yml

```



## 配置文件解析    

以下针对host配置文件配置项进行说明：    

```
# 
# Create an all group that contains the masters,workers and installer groups
[all:children]  # 定义包含所有主机组的主机组
masters
workers
installer
ceph
db
redis
logs
monitor
ingress
nfs
minio
salt
sgikes

# define global variables
[all:vars]
#SSH user, this user should allow ssh based auth without requiring a password
ansible_ssh_user=root # 主机的用户名，目前暂支持用户名登录
ansible_ssh_pass=******* # 主机密码
ansible_port=22 # 登录主机的ssh端口

# 当前客户，此配置和业务的helmfile定义的app_env_flag一致，打包及部署业务均需要用到
app_env_flag='shtu' # Named after current customer

# offline-pot path, exec install-tke-installer.sh will be auto set. don't change!!!
# 部署目录，执行install-tke-installer.sh时会自动设置，不要修改此配置!!!
dpl_dir='/data/offline-pot'

# remote image registry url, exec builder.sh will be auto set. don't change!!!
# 远程镜像仓库域名
remote_img_registry_url='reg.xx.yy.com'

# download enable or disable proxy script domain for devnetcloud
proxy_domain='download.devenv.xxx.com'

# switch ---  开关配置集
# 是否检查数据盘大小开关，默认检查
check_data_disk_size_switch='true' # check data disk size switch,the value must be true or false
# 是否初始化数据盘，当数据盘仍然是裸设备时，请将其设置为true
data_disk_init_swith='true' # when the data disk is raw device,need init the data disk and mount to /data dir,the value must be true or false

# 是否检查可以访问外网服务
check_ex_net='true' # when deploy wx this variables set true,the value must be true or false.

# 是否将磁盘恢复成裸设备
recover_disk_to_raw_switch='false' # recover disk to raw switch,the value must be true or false

# 是否部署ceph
deploy_ceph='false' # when deploy ceph need to check ceph disk whethere is raw device,the value must be true or false
# 是否使用calico,暂时上不支持
use_calico='false' # pod network component whether use calico,default unuse,the value must be true or false.current not support calico
# 是否进行压测(网络/磁盘)
stress_perfor='true' # whether stress performance switch,default true, the value must be true or false.
# 是否部署redis
deploy_redis='true' # whether deploy redis server,default true, the value must be true or false.
# 是否部署mysql
deploy_mysql='true' # whether deploy mysql server,default true, the value must be true or false.
# 是否部署prometheus
deploy_prometheus='true' # whether deploy prometheus,default true, the value must be true or false.
# 是否部署helmtiller
deploy_helmtiller='true' # whether deploy helm tiller,default true, the value must be true or false.
# 是否部署kafka
deploy_kafka='true' # whether deploy kafka, default true, the value must be true or false.
# 是否部署elk
deploy_elk='true'  # whether deploy elk, default true, the value must be true or false; resources are not enough set false
# 是否部署nfs
deploy_nfs='true' # whether deploy nfs, default true, the value must be true or false
# 是否部署minio开关
deploy_minio='false' # whether deploy minio, default true, the value must be true or false
# 是否部署业务开关
deploy_business='true' # whether deploy business, defaut true, the value must be true or false
# 是否部署salt minion agent 作为业务的cicd agent，只有部署了业务后才考虑是否部署
deploy_salt_minion='true' # whether deploy salt-minion for business cicd, default true, the value must be true or false
# 当前部署postgres 为kong提供数据持久化
deploy_postgres='true' # whether dpeloy postgres, default true, the value must be true or false
# 部署sgikes, 此服务提供搜索和分词能力，当前主要在t*pd使用，默认不部署
deploy_sgikes='false' # whether dpeloy sgikes, default false, the value must be true or false; sgikes current for t*pd

# check items variables --- 检查项变量
# 检查数据盘目录或者磁盘初始化的挂载目录，默认是数据盘，其他盘请结合lvm的设置进行修改路径
check_path='/data' # check data disk dir, or for mount disk
# 检查数据盘多大才满足需求
check_data_disk_size='200' # check data disk size whether meets requirements
# 数据盘设备名，或者其他盘需要初始化时也需要修改此变量值
data_disk_name='vdb' # check data disk whethere is raw for data disk init,lsblk to get
# ceph 盘设备名
ceph_disk_name='vdc' # check ceph disk whethere is raw device,lsblk to get
# 检查ceph盘多大才符合需求
check_ceph_disk_size='200' # check ceph disk whether meets requirements.
# 磁盘io延迟 检查值
disk_io_latency='25000' #  disk randread io latency reference
# 磁盘iops 检查参考值
disk_iops='300' # disk io qps reference
# 磁盘压测时压测数据大小
fio_size='10G' # fio test io performance file size
 
# check out server domains, just support two; domain1 is https, domain2 is http
# 检查外部服务域名，当前仅支持两个；域名1是https，域名2是http
out_ser_domain1='api.xx.yy.com'
out_ser_domain2='reg.xx.yy.com'

# create lvm and mount disk variables ---  创建lvm和挂载磁盘变量
# vg 名
vg_name='vg_data' # required, volume group name
lv名
lv_name='lv_data' # required, logical volume name
# 磁盘设备名带/dev
disk_device_name='/dev/vdb' # data disk device name,fdisk -l to get
# 文件系统类型
filesystem='xfs' # optional, default is 'xfs'
# 执行lvm操作命令，有可能需要根据实际情况调整
partition_cmd='n\np\n1\n\n\nt\n8e\nw' # when create partition failed need check system create partition proccess with manual


# offline registry config
# if tkestack not deploy, registry_domain's value must be "registry.tke.com" for harbor
# 本地镜像仓库域名, 当没有部署tkestack时，这变量值必须是 "registry.tke.com"  对于 harbor来说
registry_domain='registry.tke.com' # offline registry domain

# tke config --- tke 配置
# 当部署业务的且nginx-ingress 需要复用master主机组的第二，第三个节点作为统一入口情况下调整为1；仅部署tkestack时默认会被设设置为2（tke组件副本数）
tke_replicas="1" # tke components's replicas number, default 1, Adjust according to the actual situation
# 使用已有高可用方案，也就是使用lb作为k8s 高可用统一入口;需要在lb上tcp方式绑定到三个master节点的6443端口
tke_ha_type="third" # tke ha type , third: will be use lb, tke: will be use keepalived, none: is not ha deploy
# k8s 版本，出包时固定
k8s_version='1.16.6' # kubernetes version, need forllow tkestack install pkg
# 网卡名称, 根据实际情况设置
net_interface='eth1'  # network insterface name, need all machine is the same name
# k8s pod cidr
cluster_cidr='172.16.0.0/19' # tke cluster's pod and service network cidr,default is 172.16.0.0/19.
# lb的ip或者vip（keepalived ip）用来作为k8s高可用
tke_vip='172.17.0.6' # kubernetes master's ha vip address(lb or float ip), ha will be use
# k8s api 端口
tke_vport='6443' # k8s api port, will be config on lb
# 集群能创建svc最大数
max_cluster_service_num=256 # cluster service max number
# 每个节点最多能启动的pod数
max_node_pod_num=256 # per node max pod number
# docker 数据保存目录
docker_data_root='/data/docker' # docker data root dir
# kubelet 数据保存目录
kubelet_root_dir='/data/kubelet'
# 登录tke 控制台用户名
tke_admin_user='admin' # tke controller platform admin user name
# 登录tke 控制台密码
tke_pwd='admin' # tke controller platform admin user password
# 访问tke控制台域名
tke_console_domain='console.tke.com' # tke console domain
# 是否开启ipvs，默认开启
ipvs='true' # whether enable ipvc,default true,must be true or false

# salt_minion configs
# sat minion 配置
salt_master_domain='s.master.tke.com' # salt master's domain
salt_master_port='44506' # salt master's port

# redis configs -- redis配置
# redis 部署模式, 当前支持主从和集群模式，集群模式默认三主三从，数据不落盘
redis_mode='master-slave' # redis deploy mode, support master-slave or cluster.
# redis监听端口
REDIS_PORT='10101' # set redis listen port
# redis 密码
REDIS_PASS='redis_P@s5' # set redis password

# redis-cluster config, redis 集群相关配置信息
redis_image_tag='6.0.5-debian-10-r2' # redis image tag; redis 镜像tag
redis_nodes='6' # redis cluster node number, include master and slave； redis节点数，包含master和slave
redis_replicas='1' # redis cluster replicas number，redis 副本数
use_aof='no' # redis data whethere use AOF persistence, yes or no, 缓存数据是否采用AOF落盘
redis_persistence='false' # redis data whethere use persistence, true or false, 是否持久化，当落盘时需要将此设置为true
redis_data_dir="/data/redis" # redis persistence data dir, redis 持久化数据目录
redis_exporter_img_tag='1.6.1-debian-10-r28' # redis exporter image tag, redis exporter镜像tag
redis_sysctl_img_tag='buster' # redis sysctl image tag, redis initcontainer 镜像tag 
redis_taints='true' # redis node whetere add taints for not allow other server pod  schedule, true not allow; true or false; redis 节点是否打污点不允许别的服务pod调度
redis_cluster_client_img_tag='6.0.5-debian-10-r0' # redis cluster client image tag; redis cluster客户端镜像tag

# mysql configs -- mysql 配置
# mysql buffer 值
MYSQL_BUFFER='521M' # mysql buffer  value
# mysql 监听端口
MYSQL_PORT='3306' # mysql listen port
# mysql 数据保存目录
MYSQL_DATADIR='/data/mysql' # mysql data dir
# mysql root密码
MYSQL_PASS='mysql_P@s5'
# mysql mode, wx please set, wx请设置如下值:
# sql_mode=ONLY_FULL_GROUP_BY,STRICsql_mode=STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
# git please set，git请设置如下值:
# sql_mode=STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
# t*pd , please set, t*pd 请设置如下值:
# sql_mode=
# common's please set ' ', default is wx's value, 通用请设置为单引号里加空格，默认是wx的值
MYSQL_MODE='sql_mode=ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'

# postgres configs, postgres will be deploy db hosts group's second node
# postgres 配置，postgres 将会部署在db主机组的第二个节点也就是db的备份节点
POSTGRES_DATA_DIR='/data/postgres/data' # postgres 数据目录
POSTGRES_PASSWORD='pgdb_P@s5'
POSTGRES_IMAGE='library/postgres:12-alpine' # postgres镜像名

# ingress config -- nginx ingress controller 配置
# ingress controller 副本数(实例数)，当nginx-ingress 需要复用master主机组的第二，第三个节点作为统一入口情况下副本数为2；若无需复用maste主机组可以根据实际情况调整
ingress_replica='2' # ingress will be deploy number,default 2
# 是否共享主机网络，默认true，当有lb的驱动支持svc LoadBalancer绑定可以采用LoadBalancer时，设置为false
ingress_host_network='true' # when has LoadBalancer ip for nginx-ingress and deploy on master nodes set false;
# svc类型，当有lb作为ingress 统一入口，请设置为LoadBalancer
ingress_svc_type='ClusterIP' # when has LoadBalancer ip for nginx-ingress set LoadBalancer, otherwise set ClusterIP
# LoadBalancer ip 地址，作为ingress controller 访问统一入口
ingress_lb_ip='' # set LoadBalancer ip for nginx-ingress

# kafka and zookeeper configs -- kafka 和 zookeeper的配置
# kafka 持久化目录
kafka_data='/data/kafka' # kafka and zookeeper data save dir 
# kafka cpu limit 值,根据规模调整
kafka_limit_cpu='1' # Adjust according to the actual situation,1 eq 1c
# kafka 内存limit 值,根据规模调整
kafka_limit_mem='2Gi' # Adjust according to the actual situation
# kafka cpu request 值,根据规模调整
kafka_request_cpu='500m' # Adjust according to the actual situation, 1000m eq 1c
# kafka 内存request值,根据规模调整
kafka_request_mem='1Gi' # Adjust according to the actual situation
# kafka 堆大小,根据规模调整
kafka_heap_options='-Xmx2G -Xms2G' # Adjust according to the actual situation
# zk 堆大小,根据规模调整
zk_heap_size='2G' # Adjust according to the actual situation
# kafka镜像
kafka_image_name='library/cp-kafka' # Adjust according to the actual situation
# kafka镜像tag
kafka_image_tag='5.0.1' # Adjust according to the actual situation
# zk镜像名
zk_image_name='library/zookeeper' # Adjust according to the actual situation
# zk 镜像tag
zk_image_tag='3.5.5' # Adjust according to the actual situation
# kafka manager 镜像名
kafka_manager_image_name='library/kafka-manager' # Adjust according to the actual situation
# kafka manager  镜像tag
kafka_manager_image_tag='1.3.3.22' # Adjust according to the actual situation
# kafka manager 控制台用户名
kafka_manager_username='admin' # Adjust according to the actual situation
# kafka manager 控制台密码
kafka_manager_pwd='admin@123654' # Adjust according to the actual situation 

# elk configs --- elk 配置
# es持久化目录
es_data='/data/es' # save es data dir, Adjust according to the actual situation
# logstash 副本数，根据规模调整
logstash_replicas='1' # logstash replicas number, Adjust according to the actual situation
# logstash 内存limit值，根据规模调整
logstash_mem_limit='2Gi' # logstash memory limit, Adjust according to the actual situation
# logstash 内存request值，根据规模调整
logstash_mem_req='2Gi' # logstash memory request, Adjust according to the actual situation
# es java 堆大小，根据规模调整
es_java_opts='-Xmx1g -Xms1g' # es jave options, Adjust according to the actual situation
# es cpu request 值，根据规模调整
es_request_cpu='100m' # Adjust according to the actual situation,1000m eq 1c
# es 内存request, 根据规模调整
es_request_mem='2Gi' # Adjust according to the actual situation
# es cpu limit ，根据规模调整
es_limit_cpu='1' # Adjust according to the actual situation,1 eq 1c
# es 内存limit 根据规模调整
es_limit_mem='2Gi' # Adjust according to the actual situation
# kibana  cpu request， 根据规模调整 
kibana_request_cpu='100m' # Adjust according to the actual situation
# kibana 内存request ， 根据规模调整
kibana_request_mem='500Mi' # Adjust according to the actual situation
# kibana cpu limit，根据规模调整
kibana_limit_cpu='1' # Adjust according to the actual situation
# kibana 内存limit ，根据规模调整
kibana_limit_mem='1Gi' # Adjust according to the actual situation
# es(kibana) 密码
es_pwd='2020#happyNY' # Adjust according to the actual situation
# es(kibana) 用户名
es_uname='elastic' # Adjust according to the actual situation

# nfs config -- nfs 配置
# nfs 数据目录
nfs_data='/data/nfsdata' # nfs data dir
# 需要用到nfs业务名
nfs_app_list='("wx-uni" "wx-web")' # need nfs storage app's name,must be shell array
# nfs pv 大小，业务的pvc必须和此值保持一致;不能大于此值
nfs_pv_storage_size="2Gi" # nfs pv storage size, app's pvc must be match this size 
# 是否创建nfs pv给业务，默认不会创建
is_create_pv="false" # default will be not create pv, must be true or false

# minio config
minio_img_name='library/minio'  # minio 镜像名
minio_img_tag='RELEASE.2019-12-17T23-16-33Z' # minio 镜像tag
minio_mcimg_name='library/mc' # minio-mc 镜像名
minio_mcimg_tag='edge' # minio-mc 镜像tag
minio_mount_path='/data/minio' # minio data  dir # minio 数据目录
minio_cpu_request='250m' # minio cpu request
minio_mem_request='256Mi' # minio memory request
minio_domain='minio.pot.tke.com' # 访问minio 域名

# sg-ik-es ,搜索，分词
# sg ik es 镜像url
sg_ik_repository="library" # sg ik es registry uri
sg_ik_busyboxversion="1.29.3" # sg ik busybox image tag 
sg_ik_elkversion="6.8.0" # sg ik elasticsearch image tag
sg_ik_sgversion="25.1.ik" # sg ik searchguard image tag
sg_ik_sgkibanaversion="18.3" # sg ik kibana image tag
sg_ik_heapSize="2g" # sg ik es heap size, please adjust according to the actual situation，es的堆大小，生产环境请调大
sg_ik_cpu_limit="1" # sg ik es cpu limit, please adjust according to the actual situation，es的cpu limit值，生产环境请调大
sg_ik_mem_limit="4Gi" # sg ik es memory limit, please adjust according to the actual situation，es 内存limit值，生产环境请调大
sg_ik_cpu_req="500m" # sg ik es cpu request, please adjust according to the actual situation，es cpu request值，生产环境请调大
sg_ik_mem_req="2Gi" # sg ik es memory request, please adjust according to the actual situation，es 内存 request值，生产环境请调大
data_size="500Gi" # es data node pvc size, please adjust according to the actual situation，es 数据节点数据大小
master_data_size="200Gi" # es master node pvc size, please adjust according to the actual situation，es master节点数据大小
sg_ik_kibana_cpu_limit="500m" # sg ik kibana cpu limit, please adjust according to the actual situation，kibana cpu limit，生产环境请调大
sg_ik_kibana_mem_limit="1Gi" # sg ik kibana memory limit, please adjust according to the actual situation，kibana 内存limit，生产环境请调大
sg_ik_kibana_cpu_req="100m" # sg ik kibana cpu request, please adjust according to the actual situation，kibana cpu request，生产环境请调大
sg_ik_kibana_mem_req="500Mi" # sg ik kibana memory request, please adjust according to the actual situation，kibana 内存 request，生产环境请调大
sg_ik_ingress_class="nginx" # sg ik ingress class ，ingress class，当前私有化仅集成nginx的ingress controller
sg_ik_kibana_domain="sgik-kibana.t*pd.tke.com" # sg ik kibana domain，kibana域名
sg_ik_es_data="/data/sg-ik-es" # sg ik es data dir，es master节点，数据节点数据父目录

# harbor's config
harbor_http="80" # harbor proxy http port
harbor_https="443" # harbor proxy https port
harbor_admin_password="Harbor12345" # harbor admin password
harbor_db_password="root123" # db password
harbor_data_volume="/data/registry" # harbor registry dir

# Create installer group -- 安装机主机组, 请配置为当前安装机的ip地址，而不是127.0.0.1的地址
[installer]
127.0.0.1

# Create masters group -- k8s master 主机组, 1 或 3 节点数
[masters]
172.17.0.44

# Create workers group -- k8s 节点主机组
[workers]

# Create db group # mysql 主机组，必须是两个节点，主备
[db]

# create redis group, when master-slave mode must be two nodes; cluster mode less three nodes and 
# persistence current just support three nodes
# redis 主机组，主从模式时必须两个节点；集群模式建议三个节点做三主三从, 默认不持久化；当持久化时当前默认情况只支持三节点的三主三从模式；可以修改local pv已便支持。
[redis]

# Create ceph group # ceph 主机组，暂时预留尚未支持ceph
[ceph]

# Create logs group # elk，kafka主机组，必须三个节点
[logs]

# create monitor group # 预留
[monitor]

# create ingress controller group # ingress 主机组,当ingress 复用master主机组时，请设置为master主机组的
# 第二和第三个节点作为ingress主机组
[ingress]

# create nfs server group # nfs 主机组，必须1个节点
[nfs]

# create minio group, must be four nodes # minio 主机组，必须4个节点
[minio]

# create salt-minion group, there are several centers to define several nodes,defualt is the first of master's group
# salt-minion node must be can exec kubelet cmd
# 创建salt-minion 主机组，通常情况下当前客户部署了几个中心的业务(如：wx,t*pd等)就需要配置几个salt-minion节点；
# 默认将第一个master节点作为salt-minion的节点，查看有几个中心可以在offline-pot目录下执行：    
# ls roles/business/helms/ |grep -v README.md|grep -v charts|grep -v helmfile.d|grep -v secrets 获取
# salt-minion 必须可以执行当前k8s 集群kubectl 命令的节点(也就是k8s集群的节点);
[salt]
172.17.0.44

# sg ik elasticsearch, current for t*pd,three nodes or six nodes, when six node es master is 0~2, es data is 3~5
# sg ik es, 支持搜索，分词；当前主要应用在t*pd；支持三个节点或6个节点，当六个节点时前三个是master节点；后三个时数据节点
[sgikes]


```

## 开始部署    


##### 1. ssh 登录安装机   

```
mkdir -p /data # 创建存放安装包目录,建议/data/目录

```    

##### 2. 上传安装包    

```
将安装包offline-pot*.tar.gz 上传至安装机/data 目录   
     
```   

##### 3. 解压安装包       

```
安装机： 
cd /data && tar -zxf offline-pot*.tar.gz

```
##### 4. 配置offline-pot

```    
进入安装目录配置offline-pot :  
cd /data/offline-pot/ && cp hosts.tpl  hosts ; 
按上述配置文件解析章节进行配置: 
通常只需要配置主机组，主机用户名密码，ssh端口，数据盘磁盘名，tke的vip地址或lb地址，其他酌情根据实际情况调整

注：针对组件相关的部署开关，请不要再次配置； 因为在打包的时候已配置完毕(也就是以deploy_开头的配置项)。    

```        

##### 5. 主机初始化及主机检查     

```

完成hosts配置文件配置后，执行：
cd /data/offline-pot/
./init-and-check.sh # init-and-check.sh 当不部署tkestack时只做离线yum repo部署和helm工具的安装!!!

此脚本将进行主机安装tke-installer，主机初始化以及主机检查操作 ；查看日志请到/opt/tke-installer/data/ 目录
tke.log # tke安装日志
host-init.log # 主机初始化日志文件
hosts-check.log # 主机检查日志文件

注：主机初始化和主机检查完毕后，若操作系统使centos7.x 请务必执行下, 但升级内核有一定风险：    
cd /data/offline-pot/
./offline-pot-cmd.sh -s update-kernal.sh # 更新内核避免es oom， 本操作会重启出installer节点外的所有节点使内核生效。

```    

##### 6. 安装offline-pot

**注：若部署tkestack请自行在安装机安装ansible！！！**

```
当主机初始化以及主机检查ok时，则开始进行offlin-pot的部署:
cd /data/offline-pot/
./install-offline-pot.sh # 此脚本将开始安装tke，基础组件，以及业务部署或只部署tkestack;安装日志查看:    

tail -f /opt/tke-installer/data/tke.log


```    

##### 7. 重装    

```
当集群由于某种原因安装失败，经修复后执行如下脚本进行重装即可：      
cd /data/offline-pot/
./reinstall-offline-pot.sh

```    

##### 8. 单独执行某项操作    

```
当需要单独执行某项操作(比如部署业务或移除业务):    
cd /data/offline-pot/
./offline-pot-cmd.sh -s deploy-business.sh -f dpl_business

./offline-pot-cmd.sh 通过-s 参数指定 mgr-scripts 目录下要执行的脚本；-f 参数执行mgr-scripts 目录下
某个脚本的具体函数，此脚本只能用在安装机安装了tke-installer 情况下；若未安装tke-installer, 需要单独执行某项操作请直接到
mgr-scripts 目录下执行对应脚本对应函数即可!!!

```    

##### 9. 进入redis cluster cli    

```    
kubectl exec -ti `kubectl get pods -n pot | grep redis-cluster-client | awk '{print $1}'` /bin/bash -n pot

redis-cli -c -h redis-cluster -a $REDIS_PASSWORD

```    


