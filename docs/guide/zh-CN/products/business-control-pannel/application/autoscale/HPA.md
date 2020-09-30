# HPA

HPA 会基于 CPU、内存等指标对负载的 Pod 数量动态调控，达到工作负载稳定的目的。

依赖：[metrics-server](https://github.com/kubernetes-sigs/metrics-server)（**当前global集群自带 metrics-server，导入集群需要检查其是否安装**）

## 安装依赖

### 安装 metrics-server

Kubernetes Metrics Server 是一个集群范围的资源使用数据聚合器，是 Heapster 的继承者。metrics-server 通过从 kubernet.summary_api 收集数据收集节点和 Pod 的 CPU 和内存使用情况。Summary API 是一个内存有效的 API，用于将数据从 Kubelet/cAdvisor 传递到 metrics-server，下图为 HPA 和 kubectl 等调用 metrics-server 获取相关信息的原理图。

![image-20200929172542934](../../../../../../images/image-20200929172542934.png)

metrics-server yaml 参考 https://github.com/kubernetes-sigs/metrics-server/releases 

具体请安装配置参考 metrics-server git地址 https://github.com/kubernetes-sigs/metrics-server

## 使用 HPA

TKEStack 已经支持在页面多处位置为负载配置 HPA

1. 新建负载页（负载包括 Deployment，StatefulSet，TApp）这里新建负载时将会同时新建与负载同名的 HPA 对象：

   ![image-20200929173056091](../../../../../../images/image-20200929173056091.png)

2. 负载列表页（负载包括 Deployment，StatefulSet，TApp）

   ![image-20200929173209190](../../../../../../../../../../Typora/images/image-20200929173209190.png)

   * 点击“更新实例数量”，进入配置界面如图所示，这里将会同时新建与负载同名的 HPA 对象：

     ![image-20200929173300650](../../../../../../images/image-20200929173300650.png)

3. 自动伸缩的 HPA 列表页。此处可以查看/修改/新建 HPA：

   ![image-20200929173933713](../../../../../../images/image-20200929173933713.png)

   * 点击上图中的新建，出现新建 HPA 页面，如下图所示：

   ![image-20200929173834852](../../../../../../images/image-20200929173834852.png)