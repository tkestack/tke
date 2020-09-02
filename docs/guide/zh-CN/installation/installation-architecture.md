#  产品部署架构



## 总体架构

TKEStack 产品架构如下图所示：
![](https://github.com/tkestack/tke/blob/master/docs/images/TKEStackHighLevelArchitecture@2x.png?raw=true)



## 架构说明

TKEStack 采用了 Kubernetes on Kubernetes 的设计理念。即节点仅运行 Kubelet 进程，其他组件均采用容器化部署，由 Kubernetes 进行管理。

架构上分为 Global 集群和业务集群。Global 集群运行整个容器服务开源版平台自身所需要的组件，业务集群运行用户业务。在实际的部署过程中，可根据实际情况进行调整。



## 模块说明



* **Installer**: 运行 tke-installer 安装器的节点，用于提供 Web UI 指导用户在 Global 集群部署TKEStacl控制台；
* **Global Cluster**: 统管业务集群（Cluster），并且是运行的 TKEStack 控制台的 Kubernetes 集群；
* **Cluster**: 运行业务的 Kubernetes 集群，可以通过 TKEStack 控制台创建或导入，由 Global Cluster 统一管理；

其中 Global Cluster 提供容器云平台的支撑环境和运行自身所需的各种组件，包括业务管理组件、平台管理组件、权限认证组件、监控和告警组件、registry 镜像仓库组件以及 gateway 前端页面网关组件等等。各个组件以 workload 的形式灵活部署在 global 集群中，各组件多副本高可用方式部署，单个组件异常或者主机节点掉线等故障不会影响global集群的正常运行，TKEStack仍可提供的管理功能，用户正常的业务访问不受影响。

* Auth: 权限认证组件，提供用户鉴权、权限对接相关功能；
* Gateway: 网关组件，实现集群后台统一入口、统一鉴权相关的功能，并运行控制台的 Web 界面服务；
* Platform: 集群管理组件，提供 Global 集群管理多个业务集群相关功能；
* Business: 业务管理组件，提供平台业务管理相关功能的后台服务；
* Galaxy：网络服务组件，为集群提供多种网络模式服务；
* Monitor: 监控服务组件，提供监控采集、上报、告警相关服务；
* Notify: 通知功能组件，提供消息通知相关的功能；
* Registry: 镜像服务组件，提供平台镜像仓库和 charts 仓库服务；
* Logagent: 日志管理组件，为平台提供日志管理相关服务；
* Audit: 设计组件，提供审计服务功能；
