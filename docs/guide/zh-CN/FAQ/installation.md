

# 常见问题列表：

[如何规划部署资源](#如何规划部署资源)  

[如何使用存储  ](#如何使用存储)  

[常见报错解决方案  ](#常见报错解决方案)  

[如何重新部署集群  ](#如何重新部署集群)  

### 如何规划部署资源

TKEStack支持使用物理机或虚拟机部署，采用kubernetes on kubernetes架构部署，在主机上只拥有一个物理机进程kubelet，其他kubernetes组件均为容器。架构上分为global集群和业务集群。global集群，运行整个TKEStack平台自身所需要的组件，业务集群运行用户业务。在实际的部署过程中，可根据实际情况进行调整。

安装TKEStack，需要提供两种角色的 Server：

Installer server 1台，用以部署集群安装器，安装完成后可以回收。

Global server，若干台，用以部署 Globa 集群，常见的部署模式分为三种：

1. **All in one 模式**，1台server部署 Global集群，global集群同时也充当业务集群的角色，即运行平台基础组件，又运行业务容器。global集群会默认设置taint不可调度，使用此模式时，需要手工在golbal集群【节点管理】-【更多】-【编辑Taint】中去除不可调度设置。(关于taint，[了解更多](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/))。由于此种模式不具有高可用能力，不建议在生产环境中使用。
2. **Global 与业务集群混部的高可用模式**，3台Server部署global集群，global集群同时也充当业务集群的角色，即运行平台基础组件，又运行业务容器。global集群会默认设置taint不可调度，使用此模式时，需要手工在golbal集群【节点管理】-【更多】-【编辑Taint】中去除不可调度设置。(关于taint，[了解更多](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/))。由于此种模式有可能因为业务集群资源占用过高而影响global集群，不建议在生产环境中使用。
3. **Global 与业务集群分别部署的高可用模式**，3台Server部署global集群，仅运行平台自身组件，业务集群单独在TKEStack控制台上创建（建议3台以上），此种模式下，业务资源占有与平台隔离，建议在生产环境中使用此种模式。

集群节点主机配置，请参考[资源需求](../安装部署/资源需求.md)。

### 如何使用存储

TKEStack 没有提供存储服务，Global集群中的镜像仓库、ETCD、InfluxDB等数据组件，均使用本地磁盘存储数据。如果您需要使用存储服务，建议使用[ROOK](https://rook.io/)或者[chubaoFS](https://chubao.io/)，部署一套容器化的分布式存储服务。

### 常见报错解决方案

#### 1.密码安装报错

错误情况：使用密码安装Global集群报 ssh:unable to authenticate 错误。

解决方案：将Global集群节点/etc/ssh/sshd_config配置文件中的PasswordAuthentication设为yes，重启sshd服务。

注：建议配置SSH key的方式安装Global集群。

#### 2.

### 如何重新部署集群

1. 使用如下脚本：

```shell
#!/bin/bash

rm -rf /etc/kubernetes

systemctl stop kubelet 2>/dev/null

docker rm -f $(docker ps -aq) 2>/dev/null
systemctl stop docker 2>/dev/null

ip link del cni0 2>/etc/null

for port in 80 2379 6443 8086 {10249..10259} ; do
    fuser -k -9 ${port}/tcp
done

rm -rfv /etc/kubernetes
rm -rfv /etc/docker
rm -fv /root/.kube/config
rm -rfv /var/lib/kubelet
rm -rfv /var/lib/cni
rm -rfv /etc/cni
rm -rfv /var/lib/etcd
rm -rfv /var/lib/postgresql /etc/core/token /var/lib/redis /storage /chart_storage

systemctl start docker 2>/dev/null
```

2. 清理 installe r节点 /opt/tke-installer/data 目录下的文件，重启 tke-installer 容器后，重新打开安装页面即可。

