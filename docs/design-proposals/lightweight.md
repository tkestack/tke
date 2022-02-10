# 轻量化方案

## 公共类工作

1. 有国内公网仓库可在线拉取镜像
2. TKEStack 系统组件实现chart化
3. 物料库与 installer 剥离
4. 引入 ingress 做最外层网关

## 必要组件

## 核心组件群

tke-auth：TKEStack 的认证鉴权模块，网络上需要所有集群都能够访问到tke-auth的31138地址设为鉴权webhook，第一优先安装

tke-platform：集群声明周期管理，及集群访问代理模块，第二优先安装

tke-gateway：TKEStack 的 UI 及网关模块，第三优先安装

## 可选组件群

可选组件依赖必要组件部署完成后才能正常部署
 
### 物料组件群

tke-registry：本地化镜像及 helm chart 存储模块，支持S3协议存储、hostpath

### 应用管理组件群

tke-application：基于 helm 的集群应用分发模块

### 监控组件群

influxdb：自带的监控存储，内容存储在 hostpath
tke-monitor：基于 prometheus 的监控分发采集模块，网络上需要目标集群可以访问到 influxdb 所在 node 节点（master0）的 8086 端口
tke-notify：告警通知模块，需要配合 tke-monitor 使用

### 日志组件群

新方案待定

### 审计组件群

tke-audit：需要外部提供 es 地址，存储审计内容

### mesh 组件群

目前没有客户，建议废弃