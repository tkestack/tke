# 平台所有组件

[CronHPA](CronHPA.md)：可实现定时自动对负载的实例数量扩缩容

[CSIOperator](CSIOperator.md)：用于对接使用存储资源

[GPUManager](GPUManager.md)：用于支持容器使用 GPU 资源，支持给容器绑定非整数张卡

[IPAM](../gal)：用于实现固定容器 IP，容器重启后可以维持其 IP 保持不变

[LBCF](LBCF.md)：用于对接使用负载均衡器

[LogAgent](LogAgent.md)：用于集群日志采集，提供多个维度的日志采集功能，并可以将日志发送给 ElasticSearch 或 Kafka

[PersistentEvent](PersistentEvent.md)：集群资源对象的事件信息默认仅在 ETCD 里存储一小时，PersistentEvent 可以将事件发送到 ElasticSearch，实现事件的持久化存储

[Prometheus](Prometheus.md)：实现集群的监控、告警功能

[TApp](TappController.md)：自研工作负载类型，支持同时部署多种类型任务，支持多种升级发布方式

[VolumeDecorator](VolumeDecorator.md)：为 Kubernetes Volume(PVC) 维护更多运行时信息的控制器，例如应用程序身份，实际使用情况等