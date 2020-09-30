# 概览

平台概览页面，可查看 TKEStack 控制台管理资源的概览。

如下图所示，在【平台管理】页面点击【概览】，此处可以展现：

![image-20200821171320826](../../../../images/overview.png)

1. 平台的资源概览
   1. 集群：TKEStack 管理的集群数量
   2. 节点：TKEStack 管理的集群下所有节点数量之和
   3. 负载数：TKEStack 管理的集群下所有负载数量，包括集群下所有的Deployment、DaemonSet、StatefulSet、TApp（如果在当前集群的[扩展组件](../platform/extender.md)安装了 TApp 组件）数量之和
   4. 业务：TKEStack 平台已有[业务](../platform/business.md)总和
2. 集群的资源状态
   1. 集群 CPU 和内存的总数、分配数、利用率
   2. 节点：集群节点数量
   3. Workload：集群Workload数量，包括集群下所有的Deployment、DaemonSet、StatefulSet、TApp（如果在当前集群的[扩展组件](../platform/extender.md)安装了 TApp 组件）数量之和
   4. Master&ETCD：检查该组件状态。注意：如果导入一个云厂商的托管集群，是没有该组件的，因此这里会显示异常，但不影响集群的使用
3. 快速入口
   1. 创建独立集群
   2. 创建角色
   3. github-issue：**如有任何平台使用问题，欢迎[提出 Issue](https://github.com/tkestack/tke/issues)，我们会认真对待每个 issue**
4. 实用提示
   1. 平台实验室：体验平台最新功能
   2. 使用指引：通过创建业务，管理集群资源配额来使用平台
5. 右上角可查看 TKEStack 的在线文档



