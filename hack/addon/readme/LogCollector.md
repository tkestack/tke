## LogCollector说明

### 组件介绍

日志收集功能是容器服务为用户提供的集群内日志收集工具，可以将集群内服务或集群节点特定路径文件的日志发送至 Kafka 的指定 topic 或者 日志服务 CLS 的指定日志主题。

日志收集功能需要为每个集群手动开启。日志收集功能开启后，日志收集 Agent 会在集群内以 Daemonset 的形式运行。用户可以通过日志收集规则配置日志的采集源和消费端，日志收集 Agent 会从用户配置的采集源进行日志收集，并将日志内容发送至用户指定的消费端。

需要注意的是，使用日志收集功能需要您确认 Kubernetes 集群内节点能够访问日志消费端。

### 部署在集群内kubernetes对象

在集群内部署LogCollector Add-on , 将在集群内部署以下kubernetes对象

| kubernetes对象名称 | 类型 | 默认占用资源 | 所属Namespaces |
| ----------------- | --- | ---------- | ------------- |
| log-collector |DaemonSet |每节点0.3核CPU, 250MB内存|kube-system|

## LogCollector使用场景

日志收集功能适用于需要对 Kubernetes 集群内服务日志进行存储和分析的用户。用户可以通过配置日志收集规则进行集群内日志的收集并将收集到的日志发送至 Kafka 的指定 Topic 或 日志服务 CLS 的指定日志主题以供用户的其它基础设施进行消费。

## LogCollector限制条件

## LogCollector使用方法

### 安装

1. 登录[容器服务控制台](https://console.qcloud.com/tke2)。

2. 在左侧导航栏中，单击【扩展组件】，进入扩展组件管理页面。

3. 选择需要安装的LogCollector集群，点击【新建】，如图：

![](https://main.qcloudimg.com/raw/aed9a5e2549d865e37f6c77affcca582.png)

### 设置日志采集规则

1. 登录[容器服务控制台](https://console.qcloud.com/tke2)。

2. 在左侧导航栏中，单击【日志采集】，选择进行日志采集的集群。新建日志采集规则。

![](https://main.qcloudimg.com/raw/f714a15be03073c772ab52ddd8853bb3.png)
