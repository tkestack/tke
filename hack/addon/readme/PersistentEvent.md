# PersistentEvent

## PersistentEvent 介绍

Kubernetes Events 包括了 Kuberntes 集群的运行和各类资源的调度情况，对维护人员日常观察资源的变更以及定位问题均有帮助。TKEStack 支持为您的所有集群配置事件持久化功能，开启本功能后，会将您的集群事件实时导出到 ElasticSearch 的指定索引。

### PersistentEvent 使用场景

Kubernetes 事件是集群内部资源生命周期、资源调度、异常告警等情况产生的记录，可以通过事件深入了解集群内部发生的事情，例如调度程序做出的决策或者某些pod从节点中被逐出的原因。

kubernetes 默认仅提供保留一个小时的 kubernetes 事件到集群的 ETCD 里。 PersistentEvent 提供了将 Kubernetes 事件持久化存储的前置功能，允许您通过PersistentEvent 将集群内事件导出到您自有的存储端。

### PersistentEvent 限制条件

1. **注意：当前只支持版本号为5的 ElasticSearch，且未开启 ElasticSearch 集群的用户登录认证**
2. 安装 PersistentEvent 将占用集群0.2核 CPU,100MB 内存的资源
3. 仅在1.8版本以上的 kubernetes 集群支持

### 部署在集群内kubernetes对象

在集群内部署PersistentEvent Add-on , 将在集群内部署以下kubernetes对象

| kubernetes对象名称 | 类型 | 默认占用资源 | 所属Namespaces |
| ----------------- | --- | ---------- | ------------- |
|tke-persistent-event|deployment|0.2核CPU,100MB内存|kube-system|

## PersistentEvent 使用方法

### 在 扩展组件 里使用

  1. 登录 TKEStack

  2. 切换至【平台管理】控制台，选择 【扩展组件】，选择需要安装事件持久化组件的集群，安装 PersistentEvent 组件，注意安装 PersistentEvent 时需要在页面下方指定 ElasticSearch 的地址和索引

     > 注意：当前只支持版本号为5，且未开启用户登录认证的 ES 集群

### 在 运维中心 里使用

  1. 登录 TKEStack

  2. 切换至【平台管理】控制台，选择 【运维中心】->【事件持久化】，查看事件持久化列表

  3. 单击列表最右侧【设置】按钮，如下图所示：
     ![事件持久化设置](../../../docs/images/事件持久化设置.png)

  4. 在“设置事件持久化”页面填写持久化信息

     + **事件持久化存储：** 是否进行持久化存储

       > 注意：当前只支持版本号为5，且未开启用户登录认证的 ES 集群

     + **Elasticsearch地址：** ES 地址，如：http://190.0.0.1:200

     + **索引：** ES索引，最长60个字符，只能包含小写字母、数字及分隔符("-"、"_"、"+")，且必须以小写字母开头

  5. 单击【完成】按钮