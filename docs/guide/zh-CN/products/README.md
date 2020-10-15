# 产品使用指南

TKEStack 为了更方便用户使用容器平台，采取管理和使用分离的设计理念，提供两种管控平台：[平台管理控制台](platform) 主要负责全平台管理，面向管理员；[业务管理控制台](business-control-pannel) 主要负载创建/管理应用，面向使用者，让用户聚焦于其业务本身。用户也可以身兼两个角色，如果用户既是平台管理员，又是平台使用者，两种平台的切换方式也十分方便，请见 [切换控制台](controlpannel.md)。
* [平台管理控制台](platform)
  * 使用 [概览](platform/overview.md) 可以获取 TKEStack 管理的所有集群基本信息，例如所有纳管的集群数量、节点数量、负载数量、业务数等
  * 使用 [集群管理](platform/cluster.md) 可以对集群全生命周期管理，包括对集群的 CUDR，对集群资源的管理、集群监控、查看日志、事件等
  * 使用 [业务管理](platform/business.md) 可以对业务全生命周期管理，包括对业务的CUDR，业务的监控、配额、成员、以及命名空间等的管理
  * 使用 [扩展组件](platform/extender.md) 可以选择性增强集群功能，可以通过在集群里安装组件增强集群功能
  * 使用 [组织资源](platform/resource) 可以管理镜像仓库和应用商店，让用户更方便部署/管理应用
  * 使用 [访问管理](platform/accessmanagement) 可以管理平台用户和策略，可以控制用户权限
  * 使用 [监控&告警](platform/monitor&alert) 可以配置集群/节点/负载告警信息，让用户第一时间获取集群告警
  * 使用 [运维中心](platform/operation) 可以管理集群的日志、事件、审计
* [业务管理控制台](business-control-pannel)
  * 使用 [应用管理](business-control-pannel/application) 可以对 Kubernetes 资源全生命周期管理，包括命名空间、工作负载、自动伸缩、服务于路由、配置管理、存储、日志、事件等
  * 使用 [业务管理](business-control-pannel/business-manage.md) 可以查看业务详细信息、业务监控、业务成员等
  * 使用 [组织资源](business-control-pannel/resource) 可以管理镜像仓库和应用商店，让用户更方便部署/管理应用
  * 使用 [监控&告警](business-control-pannel/monitor&alert) 可以配置集群/节点/负载告警信息，第一时间获取集群告警
  * 使用 [运维管理](business-control-pannel/operation) 可以管理集群的日志、事件
*  [切换控制台](controlpannel.md)