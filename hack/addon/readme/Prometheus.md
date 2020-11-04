# Prometheus

## Prometheus 介绍

下图是 [Prometheus](https://prometheus.io) 官方提供的架构及其一些相关的生态系统组件：

![架构](../../../docs/images/prometheus-architecture.png)

- Prometheus Server：用于抓取指标、存储时间序列数据
- Exporter：暴露指标让任务来抓
- Pushgateway：push 的方式将指标数据推送到该网关
- Alertmanager：处理报警的报警组件
- PromQL：用于数据查询

Prometheus 常用的获取监控数据的方法：

* **/metrics 接口暴露**：`Prometheus` 的数据指标是通过一个公开的 HTTP(S) 数据接口获取到的，我们不需要单独安装监控的 Agent，只需要暴露一个 Metrics 接口，Prometheus 就会定期去拉取数据；对于一些普通的 HTTP 服务，我们完全可以直接重用这个服务，添加一个 `/metrics` 接口暴露给 Prometheus；而且获取到的指标数据格式是非常易懂的，不需要太高的学习成本。现在很多服务从一开始就内置了一个 `/metrics` 接口，比如 Kubernetes 的各个组件、Istio 服务网格都直接提供了数据指标接口。

* **Exporter**：**有一些服务即使没有原生集成该接口，也完全可以使用一些 [Exporter](https://prometheus.io/docs/instrumenting/exporters/) 来获取到指标数据**，比如 mysqld_exporter、node_exporter，这些 Exporter 就有点类似于传统监控服务中的 Agent，作为一直服务存在，用来收集目标服务的指标数据然后直接暴露给 Prometheus。

*  **kube-state-metrics**：Service 和 Pod 的监控都是应用内部的监控，需要应用本身提供一个 `/metrics` 接口，或者对应的 Exporter 来暴露对应的指标数据，但是在 Kubernetes 集群上 Pod、DaemonSet、Deployment、Job、CronJob 等各种资源对象的状态也需要监控，这也反映了使用这些资源部署的应用的状态。但通过查看前面从集群中拉取的指标(这些指标主要来自 APIServer 和 kubelet 中集成的 cAdvisor)，并没有具体的各种资源对象的状态指标。对于 Prometheus 来说，当然是需要引入新的 Exporter 来暴露这些指标，Kubernetes 提供了一个 [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) 就是我们需要的。将 kube-state-metrics 部署到 Kubernetes 上之后，就会发现 Kubernetes 集群中的 Prometheus 会在 kubernetes-service-endpoints 这个 job 下自动服务发现 kube-state-metrics，并开始拉取  metrics，这是因为部署 kube-state-metrics 的 manifest 定义文件 kube-state-metrics-service.yaml 对 Service 的定义包含`prometheus.io/scrape: 'true'`这样的一个`annotation`，因此 kube-state-metrics 的 endpoint 可以被 Prometheus 自动服务发现。

  关于 kube-state-metrics 暴露的所有监控指标可以参考 kube-state-metrics 的文档 [kube-state-metrics Documentation](https://github.com/kubernetes/kube-state-metrics/tree/master/Documentation)。

### TKEStack 中的 Prometheus

良好的监控环境为 TKEStack 高可靠性、高可用性和高性能提供重要保证。您可以方便为不同资源收集不同维度的监控数据，能方便掌握资源的使用状况，轻松定位故障。

TKEStack 使用 Prometheus 为 Kubernetes 集群提供监控告警服务。旨在降低对容器平台监控告警方案的实现难度，为用户提供开箱即用的监控告警能力，同时提供灵活的扩展能力以满足用户在使用监控告警时的个性化需求。允许用户自定义对接 influxdb，ElasticSearch 等后端存储监控数据。针对在可用性和可扩展性方面，支持使用 thanos 架构提供可靠的细粒度监控和警报服务，构建具有高可用性和可扩展性的细粒度监控能力。

> 指标具体含义可参考：[监控 & 告警指标列表](../FAQ/Platform/alert&monitor-metrics.md)

![image-20201001171647665](../../../docs/images/image-20201001171647665.png)

### Prometheus 使用场景

Prometheus 是 Kubernetes 监控的事实标准，为容器平台提供了一整套的监控告警解决方案。因为相关的组件众多，配置灵活多变，为了降低使用门槛， TKEStack 使用 Prometheus 为用户提供了一键部署的监控告警方案：

1. 方案预置了涵盖 Cluster、Namespace、Workload、Pod、Container 以及**业务**等层级的监控项，包括了 CPU、内存、IO、网络以及 GPU 等监控指标并自动绘制趋势曲线，帮助运维人员全维度的掌握平台运行状态。

2. 利用 AlertManager 的能力，结合 TKEStack 提供的告警模版和渠道，提供灵活的告警策略与告警管理。

3. 方案基于 Prometheus Operator 实现，利用其能力，用户可以非常方便的进行二次开发。

### Prometheus 限制条件

1. 安装 Prometheus 将占用集群 1核 CPU，600MB 内存的资源。同时随着集群规模的扩大，Prometheus 会占用更多的系统资源

2. 仅在1.8版本以上的 kubernetes 集群支持

### 部署在集群内 kubernetes 对象

在集群内部署 Prometheus , 将在集群内部署以下 kubernetes 对象

| kubernetes 对象名称                        | 类型                           | 默认占用资源 | 所属 Namespaces |
| ------------------------------------------| ------------------------------ | ---------- | ------------- |
| [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)<br />收集 k8s 集群内资源对象数据 | Deployment                     | 0.1核CPU,128MB内存      | kube-system  |
| kube-state-metrics                        | ServiceAccount                 | /      | kube-system  |
| kube-state-metrics                        | ClusterRole                    | /      | /            |
| kube-state-metrics                        | ClusterRoleBinding             | /      | /            |
| kube-state-metrics  | Service                        | /      | kube-system  |
| [custom-metrics-apiserver](https://github.com/kubernetes-sigs/custom-metrics-apiserver)<br />可支持任意 Prometheus 采集到的指标，同时也可以实现更多指标的 HPA | Deployment | **目前没有限制** | kube-system |
| [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator)<br />用于部署 prometheus，存储监控数据 | Deployment                     | 0.1核CPU,100MB内存      | kube-system  |
| prometheus-operator                       | ServiceAccount                 | /      | kube-system  |
| prometheus-operator                       | ClusterRole                    | /      | /            |
| prometheus-operator                       | ClusterRoleBinding             | /      | /            |
| prometheus-operator | Service                        | /      | kube-system  |
| [alertmanager-main](https://github.com/prometheus/alertmanager)<br />实现监控报警 | StatefulSet                   | 0.3核CPU,75MB内存      | kube-system  |
| alertmanager                  | Service                        | /      | kube-system  |
| [prometheus-k8s](https://github.com/prometheus/prometheus)<br />Prometheus 主程序 | StatefulSet                   | 0.3核CPU,200MB内存      | kube-system  |
| prometheus-k8s                            | ServiceAccount                 | /      | kube-system  |
| prometheus-k8s                            | ClusterRole                    | /      | /            |
| prometheus-k8s                            | ClusterRoleBinding             | /      | /            |
| prometheus | Service                        | /      | kube-system  |
| [node-exporter](https://github.com/prometheus/node_exporter)<br />收集集群中各节点的数据 | Daemonset                      | 0.1核CPU,128MB内存      | kube-system  |
| alertmanagers.monitoring.coreos.com       | CustomResourceDefinition       | /      | /            |
| podmonitors.monitoring.coreos.com         | CustomResourceDefinition       | /      | /            |
| prometheuses.monitoring.coreos.com        | CustomResourceDefinition       | /      | /            |
| prometheusrules.monitoring.coreos.com     | CustomResourceDefinition       | /      | /            |
| servicemonitors.monitoring.coreos.com     | CustomResourceDefinition       | /      | /            |
| thanosrulers.monitoring.coreos.com | CustomResourceDefinition | / | / |

## Prometheus 使用方法

### 安装

Prometheus 为 TKEStack 扩展组件，需要在集群的 [【基本信息】](../../../docs/guide/zh-CN/products/platform/cluster.md#基本信息) 页里开启 “监控告警”

### 使用

1. 在左侧导航栏中，单击【监控&告警】，进入监控告警管理页面，可管理告警项及通知相关设置。

   ![image-20201021134520950](images/image-20201021134520950.png)

   > 更多请参考 [告警设置](../../../docs/guide/zh-CN/products/platform/monitor&alert/alertsetting.md)

2. 在左侧导航栏中，单击【集群管理】，进入集群管理页面，可查看集群及工作负载的监控数据。

   ![image-20201021134336352](images/image-20201021134336352.png)

   > 更多使用请参考 [利用 Prometheus 监控](../../../docs/guide/zh-CN/features/prometheus.md)
