# Smoke Test and Release Test

**Author**: yongfeijia (@[yongfeijia](https://github.com/JiaYongfei))

**Status** (20201020): In development

## Abstract

项目平均3个月发布一次大版本，1个月发布一次bug fix版本，平均每天3-5次代码PR。由于缺少完善的自动化测试机制，很难保证发布的版本的质量。

## Motivation

提供基于云原生的，简单易用，统一高效的，社区用户可参与的，与开发流程紧密结合的自动化测试系统，并通过Smoke Test和Release Test来保证版本质量。

## Main proposal

### Smoke Test

每次pr或commit时执行，检查关键功能是否正确，执行时间控制在30分钟内。

用例路径：tke/test/e2e/platform/

#### 1. 创建Baremetal独立集群

目标：验证创建独立集群功能正常

步骤：

1. 申请创建cvm，准备2master，1worker节点，规格4C8G，操作系统：CentOS，Ubuntu，tLinux
2. 调用接口创建独立集群，等待集群创建成功

期望：

1. 检查集群状态正常，检查cluster接口返回正确的参数
2. 获取集群凭证，通过Kubectl 检查K8s集群工作正常

#### 2. 为群添加节点

目标：验证节点管理功能正常

步骤：

1. 申请创建cvm，准备1worker节点，规格1C2G
2. 调用接口为集群添加该节点，等等节点添加成功

期望：

1. 检查节点列表返回正常，检查节点被正确添加
2. 获取集群凭证，通过Kubectl 检查K8s下节点正常

#### 3. 节点驱逐，标签，封锁功能

目标：验证节点驱逐，标签，封锁功能正常

步骤：

1. 对worker节点添加标签，检查返回
2. 对worker节点封锁，检查返回
3. 对worker节点驱逐，检查返回

期望：

1. 检查节点返回信息正确
2. 通过Kubectl检查K8s下节点信息正确


#### 4. 为集群删除节点

目标：验证节点删除功能正常

步骤：

1. 调用接口为集群删除该节点

期望：

1. 检查节点列表返回正常，检查节点被正确删除
2. 通过Kubectl 检查K8s下节点被删除


#### 5. 删除Baremetal独立集群

目标：验证删除独立集群功能正常

步骤：

1. 保存集群访问凭证
2. 调用接口删除独立集群，等待集群删除成功

期望：

1. 查询cluster接口，检查集群被正确删除
2. 利用保存的访问凭证，通过Kubectl 检查K8s集群仍工作正常


#### 6. 导入Import集群

目标：验证导入集群功能正常

步骤：

1. 利用保存的访问凭证，调用接口导入集群，等待集群导入成功

期望：

1. 检查集群状态正常，检查cluster接口返回正确的参数
2. 获取集群凭证，通过Kubectl 检查K8s集群工作正常


#### 7. 为集群开启扩展组件

目标：验证扩展组件addon功能正常

步骤：

1. 调用addon接口，开启tapp，等待返回正确状态
2. 开启IPAM，等待返回正确状态
3. 开启CronHPA，等待返回正确状态
4. 开启Prometheus扩展组件，等待返回正确状态

期望：

1. 检查各个addon组件返回正确信息
2. 通过Kubectl 检查K8s集群下所有 pod 工作正常

### Release Test

每晚定时在release build后执行，验证通过tke-installer安装部署集群及更多组件功能，时间控制在4小时内。

TBD
