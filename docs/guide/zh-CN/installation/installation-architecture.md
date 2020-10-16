#  产品架构&能力说明

## 总体架构

TKEStack 产品架构如下图所示：
![](https://github.com/tkestack/tke/blob/master/docs/images/TKEStackHighLevelArchitecture@2x.png?raw=true)

## 架构说明

TKEStack 采用了 Kubernetes on Kubernetes 的设计理念。即**节点仅运行 Kubelet 进程，其他组件均采用容器化部署，由 Kubernetes 进行管理**。

架构上分为 Global 集群和业务集群。Global 集群运行整个容器服务开源版平台自身所需要的组件，业务集群运行用户业务。在实际的部署过程中，可根据实际情况进行调整。

## 模块说明

* **Installer**: 运行 tke-installer 安装器的节点，用于提供 Web UI 指导用户在 Global 集群部署 TKEStack 控制台；
* **Global Cluster**: 统管业务集群（Business Cluster），并且是运行的 TKEStack 控制台的 Kubernetes 集群；
* **Business Cluster**: 运行业务的 Kubernetes 集群，如上图中的 Cluster A、Cluster B、Cluster C，可以通过 TKEStack 控制台创建或导入，由 Global Cluster 统一管理。

其中 Global Cluster 提供容器云平台的支撑环境和运行自身所需的各种组件，包括业务管理组件、平台管理组件、权限认证组件、监控和告警组件、registry 镜像仓库组件以及 gateway 前端页面网关组件等等。各个组件以 Workload 的形式灵活部署在 Global 集群中，**各组件多副本高可用方式部署**，单个组件异常或者主机节点掉线等故障不会影响 Global 集群的正常运行，TKEStack 仍可提供的管理功能，用户正常的业务访问不受影响。

* **Auth**: 权限认证组件，提供用户鉴权、权限对接相关功能；
* **Gateway**: 网关组件，实现集群后台统一入口、统一鉴权相关的功能，并运行控制台的 Web 界面服务；
* **Platform**: 集群管理组件，提供 Global 集群管理多个业务集群相关功能；
* **Business**: 业务管理组件，提供平台业务管理相关功能的后台服务；
* **Galaxy**：网络服务组件，为集群提供多种网络模式服务；
* **Monitor**: 监控服务组件，提供监控采集、上报、告警相关服务；
* **Notify**: 通知功能组件，提供消息通知相关的功能；
* **Registry**: 镜像服务组件，提供平台镜像仓库和 Helm charts 仓库服务；
* **Logagent**: 日志管理组件，为平台提供日志管理相关服务；
* **Audit**: 审计组件，提供审计服务功能。

此外还有诸如 **Prometheus** 、**TApp**、**GPUManager** 等组件都可以安装在平台上的任意集群，以增强集群功能。

## 能力说明

- **原生**：TKEStack 兼容了 Kubernetes 原生服务访问模式。

- **产品特色**：TKEStack 扩展 Galaxy（网络）、TAPP（工作负载）、GPUManage（GPU）、CronHPA（扩缩容）、LBCF（负载均衡）等组件，界面化支持，插件化部署。

- **多集群管理**：提供多集群统一管理能力。

- **多租户统一认证**：支持 OIDC 和 LDAP 对接，实现企业租户身份的统一认证。

- **权限管理**：提供多租户统一认证与权限管理能力。不同于 Kubernetes RBAC，TKEStack 权限管理是基于 Casbin 模型。TKEStack 支持平台用户和业务用户，可为用户/用户组配置不同的角色，并绑定对应的策略，从而实现资源共享和访问隔离。

- **仓库管理**：集成 docker registry 和 chartmuseum 能力，支持创建公/私有仓库。支持创建有效时间范围的访问凭证。

- **运维能力**：提供集群、节点、工作负载、Pod、Container 五个粒度的监控数据收集和展示功能；提供短信、微信、邮件三种告警机制；提供容器文件、容器输出、节点文件三种日志采集方式，支持 ES、Kafka 两种消费端。

- **界面化**：大量 YAML 配置转换成可视化配置，降低使用门槛。

- **安全性**：TKEStack 支持 webtty，webtty 的鉴权接入了 TKEStack，减少对 kubeconfig 的依赖，降低集群 hack 风险。

- **版本升级**：TKEStack 采取了 Kubernetes 的代码理念，通过迭代可以不断适配新版本Kubernetes。

- **异构治理**：支持一键部署 x86/arm64 异构容器集群及多型号 GPU 卡异构容器集群。