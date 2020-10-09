# Helm

## Helm 介绍

Helm 是管理 Kubernetes 应用程序的打包工具, 更多详情请查看 [Helm 官网文档](https://helm.sh/)

TKEStack 集成 Helm 相关功能，提供了 Helm Chart 在指定集群内图形化的增删改查

### Helm 限制条件

1. 仅支持1.8版本以上的 kubernetes 版本
2. 将占用集群 0.28 核 CPU 180Mi 的资源

### 部署在集群内 kubernetes 对象

在集群内部署 Helm ，将在集群内部署以下 kubernetes 对象

| kubernetes 对象名称 | 类型       | 默认占用资源        | 所属 Namespaces |
| ------------------- | ---------- | ------------------- | --------------- |
| swift               | deployment | 0.03核CPU, 20Mi内存 | kube-system     |
| tiller-deploy       | deployment | 0.15核CPU, 80Mi内存 | kube-system     |
