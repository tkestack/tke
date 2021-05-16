# Thanos Support To TKE-Installer

**Author**: dengqiang([@t1mt](https://github.com/t1mt))

**Status** (20210303): Designing

[TOC]

## Abstract

落地服务网格（ServiceMesh），支持应用服务从单集群扩展到多集群。因此，通过使用 Thanos 提供多集群全局监控查询视图。

## Background

在《[产品架构 & 能力说明](https://github.com/tkestack/tke/blob/master/docs/guide/zh-CN/installation/installation-architecture.md#monitor--notify-%E7%BB%84%E4%BB%B6)》中，tke-monitor 模块已经支持 thanos，但 tke-installer 暂未提供 thanos 安装选项。

![image](https://raw.githubusercontent.com/tkestack/tke/master/docs/images/monitor.png)

![image](https://user-images.githubusercontent.com/2208260/108940746-df605c80-768e-11eb-8c8d-99919866bea0.png)

## Motivation

在 tke-installer 接口 POST /api/cluster 中支持 thanos 的安装配置参数。

- 安装 thanos 组件

  | 组件     | 使用范围                          | Port : NodePort            |
  | -------- | --------------------------------- | -------------------------- |
  | Query    | global 集群，tke-monitor          | 10901 : <br/>9090 : <br/>  |
  | Store    | global 集群                       | 10901 : <br>10902 : <br>   |
  | Compact  | global 集群                       | 10902 : <br>               |
  | Rule     | global 集群                       | 10901 : <br/>10902 : <br/> |
  | Receiver | 业务集群，Prometheus Remote Write | 19291 : 31141<br>          |

- 配置 tke-monitor/tke-platform-controller 为 thanos 对应组件地址

  | TKEStack 组件           | 配置修改                                       |
  | ----------------------- | ---------------------------------------------- |
  | tke-monitor-api         | http://thanos-query.tke.svc.cluster.local:9090 |
  | tke-monitor-controller  | http://thanos-query.tke.svc.cluster.local:9090 |
  | tke-platform-controller | thanos-receiver NodePort address               |

## Main proposal

### 功能范围

包括

- 在 API 层面添加 Thanos 安装配置参数（CreateClusterPara.Config.Monitor.ThanosMonitor）
- 原来 tke-monitor 对应的监控查询配置改为 Thanos-Query 接口
- 原来 tke-platform-controller 对应的监控存储配置改为 Thanos-Store 接口，用户集群使用 global 节点 NodePort 访问

不包括

- tke-installer 安装选择 Thanos 交互界面

### 接口修改

新增 config.monitor.thanos 请求参数

#### 接口示例

```json
// POST http://<tke-installer>/api/cluster
{
    "cluster": {
        "apiVersion": "platform.tkestack.io/v1",
        "kind": "Cluster",
        "spec": {
            "networkDevice": "eth1",
            "features": {
                "enableMetricsServer": true
            },
            "dockerExtraArgs": {},
            "kubeletExtraArgs": {},
            "apiServerExtraArgs": {},
            "controllerManagerExtraArgs": {},
            "schedulerExtraArgs": {},
            "clusterCIDR": "192.168.0.0/16",
            "properties": {
                "maxClusterServiceNum": 256,
                "maxNodePodNum": 256
            },
            "type": "Baremetal",
            "machines": [
                {
                    "ip": "<YOUR_MACHINE_IP>",
                    "port": <SSH_PORT>,
                    "username": "<LOGIN_USER>",
                    "password": "<PASSWORD | base64>"
                }
            ]
        }
    },
    "config": {
        "basic": {
            "username": "<TKE_LOGIN_NAME>",
            "password": "<TKE_LOGIN_PASSWORD | base64>"
        },
        "auth": {
            "tke": {}
        },
        "registry": {
            "tke": {
                "domain": "registry.tke.com"
            }
        },
        "application": {},
        "business": {},
        "monitor": {
            "thanos": {
                "bucketConfig": {
                    "type": "s3",
                    "config": {
                        "access_key": "<ACCESS_KEY>",
                        "bucket": "<BUCKET_NAME>",
                        "endpoint": "<IP:PORT>",
                        "insecure": true,
                        "secret_key": "<SECRET_KEY>",
                        "signature_version2": true
                    }
                }
            }
        },
        "logagent": {},
        "gateway": {
            "domain": "console.tke.com",
            "cert": {
                "selfSigned": {}
            }
        }
    }
}
```

### 配置文件修改

#### tke-monitor-config.yaml

修改安装配置 storage.thanos

```yaml
tke-monitor-config.yaml: |
  apiVersion: monitor.config.tkestack.io/v1
  kind: MonitorConfiguration
  storage:
    thanos:
      servers:
        - address: http://thanos-query.tke.svc.cluster.local:9090
```

#### tke-platform-controller.toml

修改安装配置 features.monitor_storage_type 和 features.monitor_storage_address

```yaml
tke-platform-controller.toml: |
  [secure_serving]
  tls_cert_file = "/app/certs/server.crt"
  tls_private_key_file = "/app/certs/server.key"

  [client]
    [client.platform]
    api_server = "https://tke-platform-api"
    api_server_client_config = "/app/conf/tke-platform-config.yaml"

  [registry]
  container_domain = "xxxx"
  container_namespace = "xxxx"

  [features]
  monitor_storage_type = "thanos"
  monitor_storage_addresses = "http://< GLOBAL_NODE_IP >:31141"
```

## Limition

1. 目前 tke-installer 不会默认对象存储，因此，安装 thanos 前，需要用户另外准备对象存储等参数；
2. tke-installer 默认安装 thanos 组件都为单副本，安装后，用户需要自行修改 thanos 各组件的副本数。
