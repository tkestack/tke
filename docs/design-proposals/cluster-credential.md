# Cluster Credential


**Author**: QianChenglong ([@QianChenglong](https://github.com/QianChenglong))

**Status** (20200427): Designing

## Abstract

为了控制集群凭证访问权限，将凭证单独抽离成资源。当前是通过`ClusterCredential.clusterName`来关联到对应集群。
这种实现会存在以下问题：

1. 必须先创建集群，然后才能创建凭证，需要需要获得`clusterName`
2. 对导入集群来说，创建时由于凭证还未创建，无法连接集群，验证集群版本、是否重复导入等


## Motivation

- 支持预创建集群凭证
- 支持创建导入集群时，可以连接集群

## Main proposal

引入可选`Cluser.spec.clusterCredentialRef`字段，指向`ClusterCredential`对象。


创建导入集群，若设置该字段：

1. 校验凭证是否有效。
2. 校验集群版本
3. 校验是否重复导入

创建集群时，若未设置该字段，则自动创建该对象。


创建凭证时，若设置`clusterName`，且集群类型为`Imported`，则继续验证有效性。

cluster controller
1. 存量数据在中，自动同步`Cluser.spec.clusterCredentialRef`
2. 设置`ClusterCredential.clusterName`指向当前`Cluster`
