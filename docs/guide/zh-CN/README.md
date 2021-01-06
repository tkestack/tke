# TKEStack - Tencent Kubernetes Engine Stack

<img align="right" width="100px" src="https://avatars0.githubusercontent.com/u/57258287?s=200&v=4">

![TKEStack Logo](https://github.com/tkestack/tke/workflows/build/badge.svg?branch=master)
![build-web](https://github.com/tkestack/tke/workflows/build-web/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/tkestack.io/tke)](https://goreportcard.com/report/tkestack.io/tke)
[![Release](https://img.shields.io/github/release/tkestack/tke.svg?style=flat-square)](https://github.com/tkestack/tke/releases)

> 在线文档地址：https://tkestack.github.io/docs/

***TKEStack*** 是一个**开源的多集群管理平台**，可以为生产环境中部署容器提供全方面能力。 ***TKEStack*** 可以简化部署和使用 [Kubernetes](https://github.com/kubernetes/kubernetes)，满足 IT 要求，并增强 [DevOps](https://en.wikipedia.org/wiki/DevOps) 团队的能力。

## 特点

* **多集群管理**
  * 提供 Web 控制台和命令行客户端，用于集中管理多个 Kubernetes 集群
  * 可与现有的身份验证机制集成，包括 [LDAP](https://en.wikipedia.org/wiki/Lightweight_Directory_Access_Protocol)，[OIDC](https://en.wikipedia.org/wiki/OpenID_Connect)，front proxy 和 public OAuth providers（例如GitHub）
  * 统一授权管理，不仅在集群管理级别，甚至在 Kubernetes 资源级别
  * 多租户支持，包括团队和用户对容器、构建和网络通信的隔离
* **应用程序工作负载管理**
     * 提供直观的 UI 界面，以支持可视化、YAML 导入、其他资源创建和编辑方法，使用户无需预先学习所有 Kubernetes 概念即可运行容器
     * 抽象的项目级资源容器，以支持跨多个集群的多个名称空间管理和部署应用程序
* **运维管理**
     * 集成的系统监控和应用程序监控
     * 支持对接外部存储，以实现持久化 Kubernetes 的事件、审计、日志
     * 限制、跟踪和管理平台上的开发人员和团队
* **插件支持和管理**
     * Authentication identity provider 插件
     * Authorization provider 插件
     * 事件持久化存储插件
     * 系统和应用程序日志持久化存储插件
     * [Galaxy](https://github.com/tkestack/galaxy)（网络）、[TApp](https://github.com/tkestack/tapp)（工作负载）、[GPUManage](https://github.com/tkestack/gpu-manager)（GPU）、[CronHPA](https://github.com/tkestack/cron-hpa)（扩缩容）、[LBCF](https://github.com/tkestack/lb-controlling-framework)（负载均衡）

## 架构

![Architecture Of TKE](../../images/TKEStackHighLevelArchitecture@2x.png)

> 请参考[架构&能力](../zh-CN/installation/installation-architecture.md)获取更多 TKEStack 产品架构与能力说明

## 安装

### 最小化安装需求

* **硬件最低配置**
  * 8核 CPU
  * 16 GB 内存
  * 100 GB 硬盘
* **操作系统**
  * ≥ Ubuntu 16.04/18.04  LTS (64-bit)
  * ≥ CentOS Linux 7.6 (64-bit)
  * ≥ Tencent Linux 2.2 

### 快速安装

1. **需求检查：** 请首先确认[安装要求](installation/installation-requirement.md)

2. **配置 Installer：** 请在您的 **Installer** 节点的终端中执行以下命令

   ```shell
   ＃ 根据安装节点的 CPU 架构选择安装软件包[amd64，arm64]
   arch=amd64 version=v1.4.0 && wget https://tke-release-1251707795.cos.ap-guangzhou.myqcloud.com/tke-installer-linux-$arch-$version.run{,.sha256} && sha256sum --check --status tke-installer-linux-$arch-$version.run.sha256 && chmod +x tke-installer-linux-$arch-$version.run && ./tke-installer-linux-$arch-$version.run
   ```

3. **配置控制台和 Global 集群：** 浏览器访问：`http://【INSTALLER-NODE-IP】:8080/index.html `，Web GUI 将指导您初始化和安装 TKEStack 的 **Global 集群 和 控制台**，您可以参考[安装步骤](installation/installation-procedures.md)

4. **使用 TKEStack：** 浏览器访问：http://console.tke.com

   > TKEStack使用 tke-installer 工具进行部署。有关更多信息，请参考 [tke-installer](../../user/tke-installer/introduction.md)
   >
   > 如果在安装过程中遇到问题，可以参考[安装部分的 FAQ](FAQ/Installation)


### 使用 TKEStack

[TKEStack 中文文档 ](https://tkestack.github.io/docs/)

## 开发

如果您想为 TKEStack 做出贡献，请首先参考 [贡献指南](../../../CONTRIBUTING.md)

在开发 TKEStack 之前，请确保已安装 [Git-LFS](https://github.com/git-lfs/git-lfs)

如果您有合格的开发环境，则只需执行以下操作：

```shell
mkdir -p ~/tkestack
cd ~/tkestack
git clone https://github.com/tkestack/tke
cd tke
make
```

可参考 [开发手册](../../devel/development.md) 以获取更多信息

## 社区

如果有使用问题、发现bug、有新的需求，我们都非常欢迎您通过的 GitHub [issues](https://github.com/tkestack/tke/issues/new/choose) or [pull requests](https://github.com/tkestack/tke/pulls) 进行交流。


## Licensing

TKEStack is licensed under the Apache License, Version 2.0. See [LICENSE](../../../LICENSE) for the full license text.