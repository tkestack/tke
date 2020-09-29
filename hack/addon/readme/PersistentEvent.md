## PersistentEvent说明

### 组件介绍

Kubernetes Events 包括了 Kuberntes 集群的运行和各类资源的调度情况，对维护人员日常观察资源的变更以及定位问题均有帮助。TKE 支持为您的所有集群配置事件持久化功能，开启本功能后，会将您的集群事件实时导出到配置的存储端。TKE 还支持使用腾讯云提供的 PAAS 服务或开源软件对事件流水进行检索。

### 部署在集群内kubernetes对象

在集群内部署PersistentEvent Add-on , 将在集群内部署以下kubernetes对象

| kubernetes对象名称 | 类型 | 默认占用资源 | 所属Namespaces |
| ----------------- | --- | ---------- | ------------- |
|tke-persistent-event|deployment|0.2核CPU,100MB内存|kube-system|

## PersistentEvent使用场景

Kubernetes事件是集群内部资源生命周期、资源调度、异常告警等情况产生的记录，可以通过事件深入了解集群内部发生的事情，例如调度程序做出的决策或者某些pod从节点中被逐出的原因。

kubernetes默认仅提供保留一个小时的kubernetes事件。 PersistentEvent 提供了将Kubernetes 事件持久化存储的前置功能，允许您通过PersistentEvent 将集群内事件导出到您自有的存储端。

## PersistentEvent限制条件

1. 安装PersistentEvent 将占用集群0.2核CPU,100MB内存的资源。

2. 仅在1.8版本以上的kubernetes集群支持。