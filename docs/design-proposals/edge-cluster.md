# 支持TKE创建边缘集群

**Author**: [attlee wang](https://github.com/attlee-wang) 

**Status** (20220127): Done

[TOC]

## 1. Abstract

用户在中心有了TKE原生的K8s集群之后，还有很多网络和TKE中心不可达的节点，如果能把这些节点和设备在TKE中心的控制面上统一管控起来，将充分的提高资源的利用率，降低运维的难度。

特别是在一些AI和IOT的边缘场景中，通过边缘K8s集群，用户可以在TKE中心的控制面统一管控洒落在各地的边缘节点，统管众多的小站点和边缘应用，为边缘AI和IOT云边一体化统一赋能。为此我们将在TKE中默认支持[SuperEdge](https://github.com/superedge/superedge)边缘集群，为用户云边端一朵云而奋斗。

## 2. Motivation

 TKE在集成[SuperEdge](https://github.com/superedge/superedge)的边缘集群之后将会有哪些功能：

-   一键创建生产级的SuperEdge边缘K8s集群；
-   一键纳管任何位置的边缘节点，只要这个节点能访问到边缘集群的kube-apiserver，就能被纳管；
-   云边协同能力，可以从中心管理和运维边缘节点和边缘应用；
-   边缘节点具有边缘自治能力，云边断网不影响应用运行，断电重启应用可被自动拉起；
-   支持ServiceGroup能力，可接入边缘多个站点，进行一键化多站点应用；

## 3. Proposal

### 3.1 目标方案

目标方案如下图：

<img src="../images/TKE_SuperEdge_ARCH.png" alt="TKE_SuperEdge_ARCH" style="zoom: 50%;" />

-   用户通过客户运维中心，登录到边缘机房部署

<1>.用户在部署完TKEStack的控制面之后，控制面会多一个**新建边缘集群**，就是创建SuperEdge的边缘集群的入口：

<img src="../images/create_edge_cluster.png" alt="TKE_SuperEdge_ARCH" style="zoom: 80%;" />

<2>. 点击创建集群界面和新建独立集群基本完全一致，需要填下的参数入下界面，填下相关信息后遍可创建出SuperEdge的边缘集群。

<img src="../images/ClusterInfo.png" style="zoom:33%;" />

<3>. 集群创建完之后，可以用edgeadm CLI在界面上添加边缘节点，节点->添加节点 页面只会显示如下信息

```powershell
./edgeadm join <Master Public/Intranet IP Or Domain>:Port --token xxxx --discovery-token-ca-cert-hash sha256:xxxxxxxxxx --install-pkg-path <edgeadm kube-* install package address path> --enable-edge=true
```

用户复制其命令，在自己的边缘节点上执行便可把边缘边缘节点添加上来。

>   批量添加边缘节点二期在做集成。

### 3.2 背后的实现逻辑

总体思路为 TKEStack + SuperEdge，对TKEStack创建原生独立集群的逻辑不做任何修改，代码中采用引源码来进行，在之后的Steps添加创建边缘独立集群SuperEdge的逻辑。

具体代码位置在`tke/pkg/platform/provider/edge/:`

```http
    CreateHandlers: []clusterprovider.Handler{
            ## TKEStack 创建TKEStack的逻辑, step按需要引用, 期间可能对部分函数会改写
            p.EnsureCopyFiles,             
			p.EnsurePreClusterInstallHook,
			p.EnsurePreInstallHook,
			...
			p.EnsurePostInstallHook,
			p.EnsurePostClusterInstallHook,
			
			// Addon SuperEdge 组件的逻辑
			....
			// 准备添加边缘节点的逻辑
	},
```

## 4. Plan

|    时间    | 关键节点                                             |   相关人员   |  进度  |
| :--------: | ---------------------------------------------------- | :----------: | :----: |
| 2022-01-27 | 输出Proposals, review 方案                           | @attlee-wang | 已完成 |
| 2022-02-11 | 输出提交代码框架，分出steps，创建tasks               | @attlee-wang |        |
| 2022-02-18 | 边缘独立集群能够创建出来                             |              |        |
| 2022-02-25 | 边缘节点能够添加                                     |              |        |
| 2022-03-04 | 能够把SuperEdge 打入TKEStack的Releases包             |              |        |
| 2022-03-11 | 提供边缘集群的部署使用方式，补充入TKEStack的使用文档 |              |        |
|            | 添加必要的单元测试和e2e测试                          |              |        |
|            | 支持TKEStack 创建边缘集群页面                        |              |        |

