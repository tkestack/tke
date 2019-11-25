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

## PersistentEvent使用方法

### 安装并设置存储端

1. 登录[容器服务控制台](https://console.qcloud.com/tke2)。

2. 在左侧导航栏中，单击【扩展组件】，进入扩展组件管理页面。

3. 选择需要安装的PersistentEvent集群，点击【新建】，如图：

![](https://main.qcloudimg.com/raw/2c1d974b5b5437ad823b83eae565ec95.png)

4. 配置事件持久化存储端。

### 更新存储端

1. 登录[容器服务控制台](https://console.qcloud.com/tke2)。

2. 在左侧导航栏中，单击【扩展组件】，进入扩展组件管理页面。

3. 选择需要更新的PersistentEvent集群，选择PersistentEvent点击【更新配置】，如图：

4. 配置事件持久化存储端。

### 在CLS控制台检索事件

1. 登录[日志服务控制台](hhttps://console.qcloud.com/cls)。

2. 在左侧导航栏中，单击【日志集管理】，选择PersistentEvent配置的日志集，打开日志检索功能。

![](https://main.qcloudimg.com/raw/e5509745ffa52df39272a7c97197a8d8.png)

3. 选择该日志集，点击检索事件，如图：

![](https://main.qcloudimg.com/raw/2707f519c5f682671909e0315878b575.png)

