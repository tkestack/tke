# suppport influxdb-cluster

**Author**: Zhiguo Lee[@rootdeep](https://github.com/rootdeep)

**Status** (202208031): Designing

## Summary
At present, an opensource influxdb cluster solution comes out (refer to: https://github.com/chengshiwen/influxdb-cluster). It's to be an alternative to InfluxDB Enterprise. Infludb cluster supports built-in HTTP API and SQL-like query language.

## Background
In tkestack, a single influxdb is provided as the default time series database. It is accepted for test enviroment, but not a good choice in production for missing HA feature.

## Motivation 
Introducing influxdb cluster in tkestack, replace the single influxdb with influxdb cluster. 


## Main proposal 

components manifests
| statefulset |  replicas |   Port                          | role      |
| ------------|-----------|---------------------------------|-----------|
| meta pod    |     3     |  8091,8089                      | data nodes|           
| data pod    |     3     |  8086,8088,8089,2003,4242,25826 | meta nodes|

|        service        |  type     |             role                                                           |
| ----------------------|-----------|----------------------------------------------------------------------------|
| influxdb-cluster-meta |  headless | give each meta node a unique name with {pod name}.{svc service} in cluster |
| influxdb-cluster-data |  headless | 1. give each data node a unique name 2. export a storage address         |


**In Scope**
- provide a step by step way to install the influxdb cluster background
- provide a write endpoint with http protocal 


**Out Scope**
- not support to install the cluster in tke-installer
- not provide images in registry 

**Integration Step**

- download chart infuxdb-cluster to a master node. (chart repo:  https://github.com/influxtsdb/helm-charts/tree/master/charts/influxdb-cluster)
- install helm tool
- modify values.yaml
  - replace storageclass name with nfs-client-provisioner (this is the defalut storageclass name)
  - modify replicas of data node to 3

- install chart, compoments will be installed in namespace default
  ```
  helm install  influxdb-cluster influxdb-cluster
  ```
- join nodes in a cluster
  ```
  kubectl exec -it influxdb-cluster-meta-0 bash
  influxd-ctl add-meta influxdb-cluster-meta-0.influxdb-cluster-meta:8091
  influxd-ctl add-meta influxdb-cluster-meta-1.influxdb-cluster-meta:8091
  influxd-ctl add-meta influxdb-cluster-meta-2.influxdb-cluster-meta:8091
  influxd-ctl add-data influxdb-cluster-data-0.influxdb-cluster-data:8088
  influxd-ctl add-data influxdb-cluster-data-1.influxdb-cluster-data:8088
  influxd-ctl add-data influxdb-cluster-data-0.influxdb-cluster-data:8088
  influxd-ctl show
  ```
- modify influxdb address to "influxdb-cluster-data.default:8086"

  - kubectl edit configmap -n tke tke-platform-controller 
    ```
    [features]
    monitor_storage_type = "influxdb"
    monitor_storage_addresses = "http://influxdb-cluster-data.default:8086"
    ```
  - kubectl edit configmap -n tke tke-monitor-controller
    ```
    tke-monitor-config.yaml: |
      apiVersion: monitor.config.tkestack.io/v1
      kind: MonitorConfiguration
      storage:
        influxDB:
          retentionDays: 45
          servers:
            - address: http://influxdb-cluster-data.default:8086
              username: <no value>
              password: <no value>
              timeoutSeconds: 10
    tke-monitor-controller.toml: |
      ...
      [client]
	    ...
        [features]
        monitor_storage_type = "influxdb"
        monitor_storage_addresses = "http://influxdb-cluster-data.default:8086"
    ```
  - kubectl edit configmap -n tke tke-monitor-api
    ```
      tke-monitor-config.yaml: |
        apiVersion: monitor.config.tkestack.io/v1
        kind: MonitorConfiguration
        storage:
          influxDB:
            servers:
              - address: http://influxdb-cluster-data.default:8086
                username: <no value>
                password: <no value>
                timeoutSeconds: 10
    ```
- restart pods
  ```
  kubectl rollout restart deployment -n tke tke-monitor-api
  kubectl rollout restart deployment -n tke tke-monitor-controller
  kubectl rollout restart deployment -n tke tke-platform-controller
  ```
- switch on monitor for global cluster from web

- change replica policy

  After switching on monitor, a projects and global database will be created in influxdb cluster, logining one of data nodes with kubectl, and change the replica policy.
  ```
	kubectl exec -it influxdb-cluster-data-0  bash
	root@influxdb-cluster-data-0:/#:  influx
	Connected to http://localhost:8086 version 1.8.10-c1.0.0
	InfluxDB shell version: 1.8.10-c1.0.0
	> alter retention policy tke on projects replication 3;
	> alter retention policy tke on global replication 3;
	> use projects
	> show retention policies;
    name    duration  shardGroupDuration replicaN default
	----    --------  ------------------ -------- -------
	autogen 0s        168h0m0s           3        false
	tke     1080h0m0s 24h0m0s            3        true
	> use global
	> show retention policies;
	name    duration  shardGroupDuration replicaN default
	----    --------  ------------------ -------- -------
	autogen 0s        168h0m0s           3        false
	tke     1080h0m0s 24h0m0s            3        true
  ```

## Test case
Bellow test case have been taken
- delete one data node and re-schedulded again
- delete one meta node and re-schedulded again
- shutdown a vm node which runing both a data node and a mete node for 30 minutes, and the pods are in terminating state

## Reference
[1] https://github.com/chengshiwen/influxdb-cluster/wiki/Home-Eng#community--communication

[2] https://github.com/chengshiwen/influxdb-cluster

[3] https://github.com/chengshiwen/influxdb-cluster/wiki

[4] https://docs.influxdata.com/enterprise_influxdb/v1.8/concepts/clustering/
