# TkeStack Controlplane HA and Loadbalance


**Author**: jamiezzhao ([@jamiezzhao](https://github.com/kingofstormland))

**Status** (20200529): Designing

## Abstract

在很多实际应用场景下，用户需要TkeStack集群控制平面提供高可用和负载均衡功能，而目前TkeStack安装包只通过Keepalived支持高可用，而不具备自主配置负载均衡的能力，用户想实现高可用负载均衡，只能依赖于外部提供的CLB方案，如腾讯云CLB。

对部分用户来说，自主配置集群控制平面负载均衡有一定的门槛。

## Motivation

- 支持自主配置负载均衡

## Main proposal
高可用集群有两种拓扑结构，即堆叠式（`Stacked`）和拓展式（`External`）。

因为堆叠式易于部署，节点较少可维护性较高，我们的高可用和负载均衡方案基于对堆叠式拓扑结构调研。

### 方案一：keepalivd + haproxy

keepalivd + haproxy 目前是比较流行的高可用负载均衡方案，大概方案如下：

1、keepalived 负责本机上haproxy高可用

2、haproxy 负责集群所有kube-apisever负载均衡（包含健康检查）

其中keepalived和haproxy可通过systemd实现自动拉起，也可通过部署static pod方式由kubelet自动拉起。目前项目中keepalived采用后者实现。

该方案haproxy需静态配置后端服务器ip，不利于动态扩缩容示例：

```
k8s-apiservers backend  # 配置apiserver，端口6443
server master-172.16.64.35 172.16.64.35:6443 check
server master-172.16.64.48 172.16.64.48:6443 check
```
此外haproxy需跟进社区版本升级，不方便维护。


### 方案二：vip dnat to kubernetes service

kube-apiserver启动过程中，自动创建kubenetes service, 示例：

```
kubectl describe  svc kubernetes
Name:              kubernetes
Namespace:         default
Labels:            component=apiserver
                   provider=kubernetes
Annotations:       <none>
Selector:          <none>
Type:              ClusterIP
IP:                10.244.255.1
Port:              https  443/TCP
TargetPort:        6443/TCP
Endpoints:         172.16.64.35:6443,172.16.64.48:6443
Session Affinity:  None
Events:            <none>
```
请求kubernets service，可以实现集群内kube-apiserver负载均衡但由于sevice 的cluster ip只在集群中有效，集群外部无法访问该service 的cluster ip实现负载均衡。

keepalived为master节点绑定了唯一的vip，该vip可被集群内部和集群外部访问，因此，可考虑通过某种规则将vip上的流量转发到cluster ip上。

service路由功能由集群节点中的kube-proxy完成，而kube-proxy主要支持iptables和ipvs两种模式，不同的模式实现负载均衡方式不一样。

#### iptables mod
iptables模式实现负载均衡由一系列iptables规则构成：
```
[root@VM_64_125_centos ~]# iptables -t nat -nxL PREROUTING
Chain PREROUTING (policy ACCEPT)
target     prot opt source               destination         
KUBE-HOSTPORTS  all  --  0.0.0.0/0            0.0.0.0/0            /* kube hostport portals */ ADDRTYPE match dst-type LOCAL
KUBE-SERVICES  all  --  0.0.0.0/0            0.0.0.0/0            /* kubernetes service portals */
DOCKER     all  --  0.0.0.0/0            0.0.0.0/0            ADDRTYPE match dst-type LOCAL

[root@VM_64_125_centos ~]# iptables -t nat -nxL KUBE-SERVICES
Chain KUBE-SERVICES (2 references)
target     prot opt source               destination         
KUBE-SVC-NPX46M4PTMTKRN6Y  tcp  --  0.0.0.0/0            10.244.255.1         /* default/kubernetes:https cluster IP */ tcp dpt:443
[root@VM_64_125_centos ~]# iptables -t nat -nxL KUBE-SVC-NPX46M4PTMTKRN6Y 
Chain KUBE-SVC-NPX46M4PTMTKRN6Y (2 references)
target     prot opt source               destination         
KUBE-SEP-UOR5H5XFFSPBYWEL  all  --  0.0.0.0/0            0.0.0.0/0            statistic mode random probability 0.50000000000
KUBE-SEP-II5EQNCDQYQVIX2G  all  --  0.0.0.0/0            0.0.0.0/0
[root@VM_64_125_centos ~]# iptables -t nat -nxL KUBE-SEP-UOR5H5XFFSPBYWEL
Chain KUBE-SEP-UOR5H5XFFSPBYWEL (1 references)
target     prot opt source               destination         
KUBE-MARK-MASQ  all  --  172.16.64.125        0.0.0.0/0           
DNAT       tcp  --  0.0.0.0/0            0.0.0.0/0            tcp to:172.16.64.125:6443
```
通过增加以下规则实现vip转发到iptables消息的负载均衡（与NodePort对应的Iptables规则类似）：
```
集群外访问规则：

iptables -t nat -A PREROUTING -d 172.16.64.58 -p tcp --dport 6443 -j KUBE-SVC-NPX46M4PTMTKRN6Y
iptables -t nat -I PREROUTING -d 172.16.64.58 -p tcp --dport 6443 -j KUBE-MARK-MASQ

集群内访问规则：

iptables -t nat -A OUTPUT -d 172.16.64.58 -p tcp --dport 6443 -j KUBE-SVC-NPX46M4PTMTKRN6Y
iptables -t nat -I OUTPUT -d 172.16.64.58 -p tcp --dport 6443 -j KUBE-MARK-MASQ
```
其中chainName：`KUBE-SVC-NPX46M4PTMTKRN6Y`生成方式参考

`pkg/proxy/iptables/proxier.go:591`

由代码可知chainName由servicename+port共同hash决定，只要k8s版本中不变更生成方式，chainName即可保持稳定不变且唯一。

iptable的规则由kube-proxy定时根据最新的enpoint刷新，因此若其中某台kube-apiserver挂掉，iptable规则需间隔一段时间之后才能更新，导致VIP消息路由到故障的apiserver。

备注：方案考虑过dnat vip 到cluster ip的方式，测试不通，原因为在一条链上只能做一次dnat，foward到postrouting上，无法找到cluster ip对应的host（cluster ip 没有绑定具体host，无路由信息）

#### ipvs mod
ipvs模式工作原理与iptables类似,当创建了Service之后，kube-proxy会在宿主机上创建一个虚拟网卡（kube-ipvs0）,并未其分配Cluster ip作为ip地址，如下所示：
```
kube-ipvs0: <BROADCAST,NOARP> mtu 1500 qdisc noop state DOWN group default 
    link/ether 3e:e9:a1:8e:d4:55 brd ff:ff:ff:ff:ff:ff
    inet 10.244.255.1/32 brd 10.244.255.1 scope global kube-ipvs0
       valid_lft forever preferred_lft forever
```
kube-proxy为这个ip地址设置三个IPVS虚拟主机，并设置使用轮询模式作为负载均衡策略（可选），通过ipvsadm可查看：
```
[root@VM_64_35_centos Fri May 29 20:56:03 kubernetes]# ipvsadm -ln 
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  10.244.255.1:443 rr
  -> 172.16.64.35:6443            Masq    1      2          0         
  -> 172.16.64.48:6443            Masq    1      0          0     
```
通过增加以下规则实现vip转发到ipvs消息的负载均衡：
```
集群外访问规则：

iptables -t nat -I PREROUTING -d 172.16.64.112 -p tcp --dport 6443 -j DNAT --to-destination 10.244.255.1:443
iptables -t nat -I PREROUTING -d 172.16.64.112 -p tcp --dport 6443 -j KUBE-MARK-MASQ

集群内访问规则：

iptables -t nat -I OUTPUT -d 172.16.64.112 -p tcp --dport 6443 -j DNAT --to-destination 10.244.255.1:443
iptables -t nat -I OUTPUT -d 172.16.64.112 -p tcp --dport 6443 -j KUBE-MARK-MASQ
```
ipvs模式优点：
具备健康检查能力，能够自动从ipvs规则中及时剔除故障的kube-apiserver ip，具备更高的可用性，同时有更多的负载均衡策略可供用户选择使用。

ipvs模式缺点：
宿主机需支持ipvs相关模块。

#### 测试用例
```
1、同vpc下非集群内节点访问vip 10次:curl -s -k https://172.16.64.58:6443

2、同vpc下集群内节点(master0,1)分别访问vip 10次:curl -s -k https://172.16.64.58:6443

3、同vpc下集群内节点(master0,1)访问kubernetest cluster ip 10次：curl -s -k https://10.244.255.1:443
```
### 总结
综合上述考虑，vip dnat to kubernetes service 方案在负载均衡规则上更接近k8s原生实现，且无需附带额外的辅助均衡软件haproxy。

而使用上kube-proxy的iptables方案无需依赖宿主机具备ipvs模块，故考虑在iptables方案下新增iptables规则，完成负载均衡。



