## Helm说明

### 组件介绍

Helm 是管理 Kubernetes 应用程序的打包工具, 更多详情请查看Helm 官网文档

腾讯云容器服务（Tencent Kubernetes Engine，TKE）集成 Helm 相关功能，提供了 Helm Chart 在指定集群内图形化的增删改查

### 部署在集群内kubernetes对象

在集群内部署Helm Add-on , 将在集群内部署以下kubernetes对象

| kubernetes对象名称 | 类型         | 默认占用资源           | 所属Namespaces |
| -------------- | ---------- | ---------------- | ------------ |
| swift          | deployment | 0.03核CPU, 20Mi内存 | kube-system  |
| tiller-deploy  | deployment | 0.15核CPU, 80Mi内存 | kube-system  |


## Helm限制条件

1. 仅支持1.8版本以上的kubernetes版本
2. 将占用集群 0.28 核CPU 180Mi 的资源

## Helm使用方法

详情见[Helm应用管理](https://cloud.tencent.com/document/product/457/32730)
