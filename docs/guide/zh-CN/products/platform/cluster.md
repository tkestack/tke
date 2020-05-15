# 集群管理

## 概念
**在这里可以管理你的 Kubernetes 集群。**

## 操作步骤

平台安装之后，可在【平台管理】控制台的【集群管理】中看到global集群。如下图所示：
   ![Global集群](../../../../images/cluster.png)

TKEStack还可以另外**新建独立集群**以及**导入已有集群**实现**多集群的管理**。

> 注意：**新建独立集群**和**导入已有集群**都属于[TKEStack架构](../../installation/installation-architecture.md)中的**业务集群**。

#### 新建独立集群

1. 登录 TKEStack，右上角会出现当前登录的用户名，示例为admin。

2. 切换至【平台管理】控制台。

3. 在“集群管理”页面中，单击【新建独立集群】。如下图所示：
   ![新建独立集群](../../../../images/createCluster.png)

4. 在“新建独立集群”页面，填写集群的基本信息。新建的集群需满足[installation requirements](../../../../docs/guide/zh-CN/installation/installation-requirement.md)的需求，在满足需求之后，TKEStack的集群添加非常便利。如下图所示，只需填写【集群名称】、【目标机器】、【密码】，其他保持默认即可添加新的集群。

   > 注意：若【保存】按钮是灰色，单击附近空白处即可变蓝

   ![集群基本信息0.png](../../../../images/ClusterInfo.png)

- **集群名称：** 支持**中文**，小于60字符即可

+ **Kubernetes版本：** 选择合适的kubernetes版本，各版本特性对比请查看 [Supported Versions of the Kubernetes Documentation](https://kubernetes.io/docs/home/supported-doc-versions/)。（**建议使用默认值**）

+ **网卡名称：** 最长63个字符，只能包含小写字母、数字及分隔符(' - ')，且必须以小写字母开头，数字或小写字母结尾。（**建议使用默认值eth0**）

+ **VIP** ：高可用 VIP 地址。（**按需使用**）

+ **GPU**：选择是否安装 GPU 相关依赖。（**按需使用**）

  + **pGPU**：平台会自动为集群安装 [GPUManager](https://github.com/tkestack/docs/blob/master/features/gpumanager.md) 扩展组件
  + **vGPU**：平台会自动为集群安装 [Nvidia-k8s-device-plugin](https://github.com/NVIDIA/k8s-device-plugin)

+ **容器网络** ：将为集群内容器分配在容器网络地址范围内的 IP 地址，您可以自定义三大私有网段作为容器网络， 根据您选择的集群内服务数量的上限，自动分配适当大小的 CIDR 段用于 kubernetes service；根据您选择 Pod 数量上限/节点，自动为集群内每台云服务器分配一个适当大小的网段用于该主机分配 Pod 的 IP 地址。（**建议使用默认值**）

  + **CIDR**： 集群内 Sevice、 Pod 等资源所在网段。

  + **Pod数量上限/节点**： 决定分配给每个 Node 的 CIDR 的大小。

  + **Service数量上限/集群** ：决定分配给 Sevice 的 CIDR 大小。

+ **目标机器** ：

  + **目标机器**：节点的内网地址。（建议: Master&Etcd 节点配置**4核**及以上的机型）

  + **SSH端口**： 请确保目标机器安全组开放 22 端口和 ICMP 协议，否则无法远程登录和 PING 云服务器。（**建议使用默认值22**）

  + **主机label**：给主机设置Label,可用于指定容器调度。（**按需使用**）

  + **认证方式**：连接目标机器的方式

    +  **密码认证**：
       +  **密码**：目标机器密码
    +  **密钥认证**：
       +  **私钥**：目标机器秘钥
       +  **私钥密码**：目标机器私钥密码，可选填

  + **GPU**： 使用GPU机器需提前安装驱动和runtime。（**按需使用**）

    > 输入以上信息后单击【保存】后还可**继续添加集群的节点**

5. **提交**： 集群信息填写完毕后，【提交】按钮变为可提交状态，单击即可提交。

#### 导入已有集群

1. 登录 TKEStack。
2. 切换至【平台管理】控制台。
3. 在“集群管理”页面，单击【导入集群】。如下图所示：
   ![导入集群](../../../../images/importCluster-1.png)

4. 在“导入集群”页面，填写被导入的集群信息。如下图所示：
   ![导入集群信息](../../../../images/importCluster-2.png)

- **名称**： 被导入集群的名称，最长60字符
- **API Server**： 
  - 被导入集群的API server的域名或IP地址，注意域名不能加上https://
  - 端口，此处用的是https协议，端口应填443。
- **CertFile**： 输入被导入集群的cert 文件内容
- **Token**： 输入被导入集群创建时的token值

5. 单击最下方 【提交】 按钮 。