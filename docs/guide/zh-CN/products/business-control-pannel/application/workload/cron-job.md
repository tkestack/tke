## 简介

一个 CronJob 对象类似于 crontab（cron table）文件中的一行。它根据指定的预定计划周期性地运行一个 Job，格式可以参考 Cron。
Cron 格式说明如下：
```
# 文件格式说明
#  ——分钟（0 - 59）
# |  ——小时（0 - 23）
# | |  ——日（1 - 31）
# | | |  ——月（1 - 12）
# | | | |  ——星期（0 - 6）
# | | | | |
# * * * * *
```

## CronJob 控制台操作指引

### 创建 CronJob
1. 登录TKEStack，切换到【业务管理】控制台，选择左侧导航栏中的【应用管理】。
2. 选择需要创建CronJob的业务下相应的【命名空间】，展开【工作负载】下拉项，进入【CronJob】管理页面。如下图所示：
   ![](../../../../../../images/CronJobNew.png)
3. 单击【创建】按钮，进入 新建Workload页面。
4. 根据实际需求，设置 CronJob 参数。关键参数信息如下：
 - **工作负载名**：输入自定义名称。
 - **标签**：给工作负载添加标签
 - **命名空间**：根据实际需求进行选择。
 - **类型**：选择【CronJob（按照Cron的计划定时运行）】。
 - **执行策略**：根据 Cron 格式设置任务的定期执行策略。
 - **Job设置**
    - **重复执行次数**：Job 管理的 Pod 需要重复执行的次数。
    - **并行度**：Job 并行执行的 Pod 数量。
    - **失败重启策略**：Pod下容器异常推出后的重启策略。
        - **Never**：不重启容器，直至 Pod 下所有容器退出。
        - **OnFailure**：Pod 继续运行，容器将重新启动。
 - **数据卷**：根据需求，为负载添加数据卷为容器提供存，目前支持临时路径、主机路径、云硬盘数据卷、文件存储NFS、配置文件、PVC，还需挂载到容器的指定路径中
   - **临时目录**：主机上的一个临时目录，生命周期和Pod一致
   - **主机路径**：主机上的真实路径，可以重复使用，不会随Pod一起销毁
   - **NFS盘**：挂载外部NFS到Pod，用户需要指定相应NFS地址，格式：127.0.0.1:/data
   - **ConfigMap**：用户在业务Namespace下创建的[ConfigMap](../configurations/ConfigMap.md)
   - **Secret**：用户在业务namespace下创建的[Secret](../configurations/Secret.md)
   - **PVC**：用户在业务namespace下创建的[PVC](../storage/persistent-volume-claim.md)
 - **实例内容器**：根据实际需求，为 CronJob 的一个 Pod 设置一个或多个不同的容器。
    - **名称**：自定义。
    - **镜像**：根据实际需求进行选择。
    - **镜像版本（Tag）**：根据实际需求进行填写。
    - **CPU/内存限制**：可根据 [Kubernetes 资源限制](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/) 进行设置 CPU 和内存的限制范围，提高业务的健壮性。
    - **GPU限制**：如容器内需要使用GPU，此处填GPU需求
    - **环境变量**：用于设置容器内的变量，变量名只能包含大小写字母、数字及下划线，并且不能以数字开头
       * **新增变量**：自己设定变量键值对
       * **引用ConfigMap/Secret**：引用已有键值对
      - **高级设置**：可设置 “**工作目录**”、“**运行命令**”、“**运行参数**”、“**镜像更新策略**”、“**容器健康检查**”和“**特权级**”等参数。这里介绍一下镜像更新策略。
       * **镜像更新策略**：提供以下3种策略，请按需选择
         若不设置镜像拉取策略，当镜像版本为空或 `latest` 时，使用 Always 策略，否则使用 IfNotPresent 策略
         * **Always**：总是从远程拉取该镜像
         * **IfNotPresent**：默认使用本地镜像，若本地无该镜像则远程拉取该镜像
         * **Never**：只使用本地镜像，若本地没有该镜像将报异常
 - **imagePullSecrets**：镜像拉取密钥，用于拉取用户的私有镜像
 - **节点调度策略**：根据配置的调度规则，将Pod调度到预期的节点。支持指定节点调度和条件选择调度
 - **注释（Annotations）**：给Pod添加相应Annotation，如用户信息等
 - **网络模式**：选择Pod网络模式
    * **OverLay（虚拟网络）**：基于 IPIP 和 Host Gateway 的 Overlay 网络方案
    * **FloatingIP（浮动 IP）**：支持容器、物理机和虚拟机在同一个扁平面中直接通过IP进行通信的 Underlay 网络方案。提供了 IP 漂移能力，支持 Pod 重启或迁移时 IP 不变
    * **NAT（端口映射）**：Kubernetes 原生 NAT 网络方案
    * **Host（主机网络）**：Kubernetes 原生 Host 网络方案
5. 单击【创建Workload】，完成创建。

### 查看 CronJob 状态

1. 登录TKEStack，切换到【业务管理】控制台，选择左侧导航栏中的【应用管理】。
2. 选择需要创建CronJob的【业务】下相应的【命名空间】，展开【工作负载】下拉项，进入【CronJob】管理页面。
3. 单击需要查看状态的 CronJob 名称，即可查看 CronJob 详情。

## Kubectl 操作 CronJob 指引

<span id="YAMLSample"></span>
### YAML 示例

```Yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            args:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```
- kind：标识 CronJob 资源类型。
- metadata：CronJob 的名称、Label等基本信息。
- metadata.annotations：对 CronJob 的额外说明，可通过该参数设置腾讯云 TKE 的额外增强能力。
- spec.schedule：CronJob 执行的 Cron 的策略。
- spec.jobTemplate：Cron 执行的 Job 模板。

### 创建 CronJob

#### 方法一
1. 参考 [YAML 示例](#YAMLSample)，准备 CronJob YAML 文件。
2. 安装 Kubectl，并连接集群。操作详情请参考 [通过 Kubectl 连接集群](https://cloud.tencent.com/document/product/457/8438)。
3. 执行以下命令，创建 CronJob YAML 文件。
```shell
kubectl create -f CronJob YAML 文件名称
```
例如，创建一个文件名为 cronjob.yaml 的 CronJob YAML 文件，则执行以下命令：
```shell
kubectl create -f cronjob.yaml
```

#### 方法二
1. 通过执行`kubectl run`命令，快速创建一个 CronJob。
例如，快速创建一个不需要写完整配置信息的 CronJob，则执行以下命令：
```shell
kubectl run hello --schedule="*/1 * * * *" --restart=OnFailure --image=busybox -- /bin/sh -c "date; echo Hello"
```
2. 执行以下命令，验证创建是否成功。
```shell+-
kubectl get cronjob [NAME]
```
返回类似以下信息，即表示创建成功。
```
NAME      SCHEDULE    SUSPEND   ACTIVE    LAST SCHEDULE   AGE
cronjob   * * * * *   False     0         <none>          15s
```

### 删除 CronJob
>!
> - 执行此删除命令前，请确认是否存在正在创建的 Job，否则执行该命令将终止正在创建的 Job。
> - 执行此删除命令时，已创建的 Job 和已完成的 Job 均不会被终止或删除。
> -  如需删除 CronJob 创建的 Job，请手动删除。
> 
执行以下命令，删除 CronJob。
```
kubectl delete cronjob [NAME]
```
