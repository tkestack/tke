# GPU-Manager

## 组件介绍

GPU Manager 提供一个 All-in-One 的 GPU 管理器, 基于 Kubernets Device Plugin 插件系统实现, 该管理器提供了分配并共享 GPU, GPU 指标查询, 容器运行前的 GPU 相关设备准备等功能, 支持用户在 Kubernetes 集群中使用 GPU 设备。

管理器包含如下功能:

- **拓扑分配**：提供基于 GPU 拓扑分配功能, 当用户分配超过1张 GPU 卡的应用, 可以选择拓扑连接最快的方式分配GPU设备

- **GPU 共享**：允许用户提交小于1张卡资源的的任务, 并提供 QoS 保证

- **应用 GPU 指标的查询**：用户可以访问主机的端口(默认为5678)的/metrics 路径,可以为 Prometheus 提供 GPU 指标的收集功能, /usage 路径可以提供可读性的容器状况查询

## 部署在集群内 kubernetes 对象

在集群内部署 GPU-Manager Add-on , 将在集群内部署以下 kubernetes 对象：

| kubernetes 对象名称       | 类型         | 建议预留资源 | 所属 Namespaces |
| --------------------- | ---------- | ------ | ------------ |
| gpu-manager-daemonset | DaemonSet  | 每节点1核CPU, 1Gi内存 | kube-system  |
| gpu-quota-admission   | Deployment | 1核CPU, 1Gi内存      | kube-system  |

## GPU-Manager 使用场景

在 Kubernetes 集群中运行 GPU 应用时, 可以解决 AI 训练等场景中申请独立卡造成资源浪费的情况，让计算资源得到充分利用。

## GPU-Manager 限制条件

1. 该组件基于 Kubernetes DevicePlugin 实现, 只能运行在支持 DevicePlugin 的 TKE 的 1.10 kubernetes 版本之上。

2. 每张 GPU 卡一共有100个单位的资源, 仅支持0-1的小数卡,以及1的倍数的整数卡设置. 显存资源是以 256MiB 为最小的一个单位的分配显存。

3. 使用 GPU-Manager 要求集群内包含 GPU 机型节点。

## GPU-Manager 使用方法

1. 集群的主机有 GPU，并且在创建时有勾选 **GPU**，已安装 GPU 插件

2. 集群安装 GPU-Manager 扩展组件

3. 在安装了 GPU-Manager 扩展组件的集群中，创建工作负载

4. 创建工作负载设置 GPU 限制，如图：

   ![](https://main.qcloudimg.com/raw/c06872ddc0fafbf92345c0d9f26e4ecd.png)


### yaml 创建

如果使用yaml创建工作负载，提交的时候需要在 yaml 为容器设置 GPU 的使用资源, 核资源需要在 resource 上填写`tencent.com/vcuda-core`, 显存资源需要在 resource 上填写`tencent.com/vcuda-memory`,

- 使用1张卡

```yaml

apiVersion: v1

kind: Pod

...

spec:

containers:

- name: gpu

resources:

tencent.com/vcuda-core: 100
```

- 使用0.3张卡, 5GiB显存的应用（20*256MB）

```yaml

apiVersion: v1

kind: Pod

...

spec:

containers:

- name: gpu

resources:

tencent.com/vcuda-core: 30

tencent.com/vcuda-memory: 20
```
