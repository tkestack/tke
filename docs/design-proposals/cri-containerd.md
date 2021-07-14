# Use standard CRI runtime containerd


**Author**: Tai Hui Min

**Status** (20210419): Designing

## Abstract
CRI（Container runtime interface）是Kubernetes容器运行时接口,是Kubernetes规定的容器运行时标准，由protocol buffer协议和gRPC API构成。Docker、Containerd、CRI-O为实现CRI的容器引擎。
以下Proposal主要介绍了TKEStack支持CRI标准容器引擎containerd的背景、动机和具体的实施步骤。


## Background
  随着Kubernetes社区的蓬勃发展，容器引擎也发展迅速，除了Docker引擎外，出现多种更加轻量级的容器引擎: 比如Containerd[1]、CRI-O[2]、iSulad[3]等；为了统一支持这些运行引擎，Kubernetes V1.5版本引入了CRI(Container runtime interface)。

  同时为了保持对Docker的支持,Kubelet内置了dockershim来支持CRI，但也因此引出了一系列问题，最终Kubernetes社区计划移除dockershim，确保Kubelet和容器引擎的界面是CRI接口[4]: 1.20-1.22标记depracated, 1.23正式移除。后续dockershim由Mirantis维护，需要在Kubelet和Docker之外额外安装[5]。

  ![CRI](../images/cri-runtime.svg)

  此外，Containerd 、CRI-O等容器引擎相比Docker更加轻量级，表现在功能更简单、创建容器时的调用链更短，性能更加健壮并且可移植性强，因此更适合作为Kubernetes的容器运行时。Containerd于2019年已经正式从CNCF毕业，成为行业标准的容器运行引擎，在各大互联网公司也有较广泛的应用 [1]。本文内容集中在TKEStack支持Containerd，Dockershim+Docker以后TKEStack将不再支持。

  Containerd作为纯粹的容器运行时工具，缺少镜像制作等能力，需要引入额外的工具满足TKEStack的需求。下表为TKEStack不同功能阶段对容器引擎的功能需求。

|  功能阶段   | 容器引擎功能点  | Docker+Dockershim | CRI标准接口(Containerd/CRIO) | containerd  + ctr |
|  -----  |   -------      |-|-|-|
| 构建     | 制作镜像        |是|否|否|
| 构建     | push镜像 |是|否|是|
| 构建     | tag镜像        |是|否|是|
| install | load镜像       |是|否|是|
| install | push镜像 |是|否|是|
| install | manifest(multi-arch) |是|否|是|
| running(管理集群: 创建、导入、增加节点、运行负载) | 启动容器        |是|是|是|
| running(管理集群) | 下载镜像        |是|是|是|
| running(管理集群) | 运行容器        |是|是|是|
| running(管理集群) | 运行POD        |是|是|是|
| running(管理集群) | 配置网络        |是|是|是|
| running(管理集群) | 删除容器        |是|是|是|
| running(管理集群) | 删除镜像        |是|是|是|



## Motivation

TKEStack支持配置CRI标准兼容容器引擎；根据TKEStack在不同阶段使用功能的差异，我们把支持CRI接口容器引擎分为四个阶段：

**阶段1  完成running阶段的改造**
  * 支持配置Containerd为容器引擎

**阶段2 完成Install阶段的改造**

  * 支持配置Containerd作为容器引擎
  *  引入其他工具如Podman支持tag、push等操作

**阶段3  完成构建阶段的改造**
  *  支持配置Containerd作为容器引擎
  *  引入其他工具如Podman完成build、tag、push功能

需要明确该步改造的必要性，如无必要则不做改造。

## Main proposal

### 阶段1 完成Running阶段的改造

该阶段容器引擎的调用者是Kubelet，只使用CRI标准定义的接口，所以只需要引入Containerd，不需要引入其他第三方组件。该阶段包含2个部分：

**1. 安装Cluster**

 *  Install过程中安装Global Cluster
 *  安装User Cluster

 两个过程均调用`clusterProvider.OnCreate`函数创建集群，其中的区别是安装Global集群过程中可以配置容器引擎的额外参数[6]。相关的流程以及修改如下

 　![clusterprovider.OnCreate流程](../images/cluster-create.svg)

修改文件：

*   `pkg/platform/provider/baremetal/cluster/provider.go`
*   `pkg/platform/provider/baremetal/cluster/create.go`

**2. user cluster扩容节点**

该过程通过调用`machineProvider.OnCreate`，和安装cluster过程类似，不再赘述。

修改文件：

*   `pkg/platform/provider/baremetal/machine/provider.go`
*   `pkg/platform/provider/baremetal/machine/create.go`

由于增加了Containerd，需要引入配置文件和server文件

* `pkg/platform/provider/baremetal/conf/containerd/config.toml`
* `pkg/platform/provider/baremetal/conf/containerd/containerd.service`

增加cri-tools，需要引入如下配置文件

* pkg/platform/provider/baremetal/conf/critools/crictl.yaml

目前存在的未明确问题：
* 目前Global集群创建过程中，是否有必要自定义容器引擎参数，通过修改配置文件是否足够？ 如果后续支持多种容器引擎且都支持自定义参数，会显得很乱。
* 目前CRI-O和Containerd均支持配置default oci runtime hook，可以替代 EnsureNvidiaContainerRuntime的功能，减少组件的数量，目前是否有必要替代?
* IncludeSelf=true时， installer会使用global cluster的一个节点，这时候需要容器引擎需要满足Install阶段的要求，一般的部署时怎样的？
### 阶段2 完成install阶段的改造
  如果使用第三方镜像仓库，在安装过程中需要搭建临时镜像仓库，并将多种体系架构的镜像push到临时仓库中完成安装。由于Docker 镜像格式原生不支持需要多体系架构，需要额外引入流程解决，过程如下：
![install过程镜像](../images/cri-image-old.svg)

存在步骤2和步骤3的根本原因是Docker的镜像格式不支持多体系架构，image.tar.gz 也同样无法包含同名的不同架构镜像。所以此类镜像加载到容器镜像后，还需要通过docker的manifest命令进行适配。

![install过程镜像](../images/cri-image-new.svg)

目前OCI Image spec通过manifest list天生支持了多体系架构，所以该过程可以简化为上图。需要注意的问题是这步改动已经超过了CRI规范的范畴，要引入Containerd自带的命令进行补充。

### 阶段3 完成构建阶段的改造
  该阶段容器引擎的使用者不是Kubernetes，所以不受Kubelet移除dockershim的影响。改造优先级较低。
  随着社区的发展，容器基础工具集开始符合Unix的设计哲学：一个工具只做好一个事情。镜像制作领域也出现了一系列工具：

 * BuildKit[7]
BuildKit是Docker出品的一款镜像制作工具，目的是分担Docker Engine的功能。
 * buildah[8]
 BuidAh是Redhat出品的容器工具集的一部分，不再作为dameon常驻后台。
 * kaniko[9]
 Kaniko是Google针对k8s场景出品的容器引擎，适合在容器中制作镜像。
### 阶段3 支持多种CRI标准兼容容器引擎

通<TKEStack不同功能阶段对容器引擎的功能需求表>可以知道，CRI标准的运行时尽仅可以满足Running状态所需的能力。所以该阶段计划等Dockershim正式从Kubelet中提出后，再考虑把Docker+dockershim增加到可选的容器引擎中。



## refereneces:
1.  https://containerd.io/
2.  https://github.com/cri-o/cri-o
3.  https://gitee.com/openeuler/iSulad
4. https://github.com/kubernetes/enhancements/blob/master/keps/sig-node/1985-remove-dockershim/README.md
5. https://www.mirantis.com/blog/mirantis-to-take-over-support-of-kubernetes-dockershim-2
6. https://tkestack.github.io/docs/installation/installation-steps.html
7. https://github.com/moby/buildkit
8. https://github.com/containers/buildah
9. https://github.com/GoogleContainerTools/kaniko
