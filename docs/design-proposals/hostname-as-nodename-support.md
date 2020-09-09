# Hostname as nodename support


**Author**: Leo Ryu ([@leoryu](https://github.com/leoryu))

**Status** (20200909): Designing

## Abstract

当前创建`Cluster`或添加`Node`逻辑中，kubelet会将`machine IP`作为`--node-name`的参数启动。
导致无法选择以`hostname`作为节点的`nodename`。

但是，现实中hostname可能会包含对用户有价值的信息，例如IDC位置信息。

## Motivation

- 允许在创建`Cluster`时选择以`hostname`或`machine IP`作为`nodename`

## Main proposal

### Stage 1 (Done)

1. 引入可选`Cluser.spec.hostnameAsNode`字段，在创建集群时供用户选择该集群是否以`hostname`作为`nodename`

2. `Cluser.spec.hostnameAsNode`为`boolean`类型。当为`true`时`Node`将以`hostname`作为`nodename`，当为`false`时`Node`将以`machine IP`作为`nodename`

3. 创建`Cluser`和添加`Node`时为每一个`Node`添加一个`lable` `platform.tkestack.io/machine-ip`，其值为`machine IP`

4. 涉及到使用`maachine IP`对`Node`进行操作对部分，首先尝试使用旧方法`machine IP`作为`nodename`的方式获取，获取不到时将使用`labelseletctor`筛选出`platform.tkestack.io/machine-ip = machine IP`的节点 （Stage 3后将完全采用`labelseletctor`的方式获取`Node`）。

PR(s): [#711](https://github.com/tkestack/tke/pull/711)

### Stage 2 (Working)

1. `cluster controller`中对`ClusterMachine`的存量`Node`添加`label` `platform.tkestack.io/machine-ip`，其值为`machine IP`（此部分逻辑未来将在Stage 3中删除）

2. `machine controller`中对`Machine`的存量`Node`添加`label` `platform.tkestack.io/machine-ip`，其值为`machine IP` （此部分逻辑未来将在Stage 3中删除）

3. `web console`中涉及到使用`nodename`获取`Machine`/`ClusterMachine`的部分修改为使用`label` `platform.tkestack.io/machine-ip`获取

### Stage 3

1. 删除`cluster controller`中对存量`Node`加`label`逻辑

2. 删除`machine controller`中对存量`Node`加`label`逻辑

3. 涉及到使用`maachine IP`对`Node`进行操作对部分，删除首先尝试使用旧方法`machine IP`作为`nodename`的方式获取的逻辑
