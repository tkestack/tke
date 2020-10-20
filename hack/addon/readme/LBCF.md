# LBCF

## LBCF 介绍 

Load Balancer Controlling Framework (LBCF) 是一款部署在Kubernetes内的**通用负载均衡控制面框架**，旨在降低容器对接负载均衡的实现难度，并提供强大的扩展能力以满足业务方在使用负载均衡时的个性化需求。

### LBCF 使用场景

LBCF 对 Kubernetes 内部晦涩的运行机制进行了封装并以 Webhook 的形式对外暴露，在容器的全生命周期中提供了多达8种 Webhook，通过实现这些 Webhook，开发人员可以轻松实现下述功能：

- 对接任意负载均衡/名字服务，并自定义对接过程
- 实现自定义灰度升级策略
- 容器环境与其他环境共享同一个负载均衡
- 解耦负载均衡数据面与控制面

### LBCF 限制条件

- K8S 1.10及以上版本
- 开启 Dynamic Admission Control，在 APIServer 中添加启动参数：`--enable-admission-plugins=MutatingAdmissionWebhook,ValidatingAdmissionWebhook`
- K8S 1.10版本，在 APIServer 中额外添加参数：`--feature-gates=CustomResourceSubresources=true`

推荐环境：

在[腾讯云](https://cloud.tencent.com/product/tke)上购买1.12.4版本的集群，无需修改任何参数，开箱可用

### 部署在集群内 kubernetes 对象

在集群内部署 LBCF，将在集群内部署以下 kubernetes 对象：

| kubernetes 对象名称                  | 类型                           | 默认占用资源 | 所属 Namespaces |
| ------------------------------------ | ------------------------------ | ------------ | --------------- |
| lbcf-controller                      | Deployment                     | /            | kube-system     |
| lbcf-controller                      | ServiceAccount                 | /            | kube-system     |
| lbcf-controller                      | ClusterRole                    | /            | /               |
| lbcf-controller                      | ClusterRoleBinding             | /            | /               |
| lbcf-controller                      | Secret                         | /            | kube-system     |
| lbcf-controller                      | Service                        | /            | kube-system     |
| backendrecords.lbcf.tkestack.io      | CustomResourceDefinition       | /            | /               |
| backendgroups.lbcf.tkestack.io       | CustomResourceDefinition       | /            | /               |
| loadbalancers.lbcf.tkestack.io       | CustomResourceDefinition       | /            | /               |
| loadbalancerdrivers.lbcf.tkestack.io | CustomResourceDefinition       | /            | /               |
| lbcf-mutate                          | MutatingWebhookConfiguration   | /            | /               |
| lbcf-validate                        | ValidatingWebhookConfiguration | /            | /               |

## LBCF 使用方法

1. 通过 【扩展组件】安装 LBCF

2. 开发或选择安装 LBCF Webhook 规范的要求实现 Webhook 服务器

3. 以下按腾讯云 CLB 开发的 Webhook 服务器为例

#### LBCF CLB driver

功能列表：

- 使用已有负载均衡
- 创建新的负载均衡（四层/七层）
- 绑定 Service NodePort
- CLB 直通 POD（直接绑定 Pod 至 CLB，不通过 Service）
- 权重调整
- 能够校验并拒绝非法参数

##### 部署 LBCF CLB driver

部署前需修改 YAML（文中最后面的[附录](#附录)已提供 YAML 文件，需要向 [Deployment ](#Deployment)文件中填入以下五点信息：

```yaml
    spec:
      priorityClassName: "system-node-critical"
      containers:
        - name: driver
          image: ${image-name}	# 1. 镜像地址
          args:
            - "--region=${your-region}"	# 2. 集群所在地域
            - "--vpc-id=${your-vpc-id}"	# 3. 集群所在 VPC 的 ID （绑定 Service NodePort时用来查找节点对应的 instanceID）
            - "--secret-id=${your-account-secret-id}"	# 4. 腾讯云账号的 SecretID
            - "--secret-key=${your-account-secret-key}"	# 5. 腾讯云账号的 SecretKey
```

登陆集群，使用以下命令安装Y AML

```
kubectl apply -f configmap.yaml
kubectl apply -f deploy.yaml
kubectl apply -f service.yaml
```

##### 使用示例

###### 使用已有四层 CLB

本例中使用了 ID 为`lb-7wf394rv`的负载均衡实例，监听器为四层监听器，端口号为 20000，协议类型 TCP。

*注: 程序会以`端口号20000，协议类型TCP`为条件查询监听器，若不存在，会自动创建新的*

```yaml
apiVersion: lbcf.tkestack.io/v1beta1
kind: LoadBalancer
metadata:
  name: example-of-existing-lb 
  namespace: kube-system
spec:
  lbDriver: lbcf-clb-driver
  lbSpec:
    loadBalancerID: "lb-7wf394rv"
    listenerPort: "20000"
    listenerProtocol: "TCP"
  ensurePolicy:
    policy: Always
```

###### **创建新的七层 CLB**

本例在 VPC  `vpc-b5hcoxj4`中创建了公网(OPEN)负载均衡实例，并为之创建了端口号为 9999 的 HTTP 监听器，最后会在监听器中创建`mytest.com/index.html`的转发规则

```
apiVersion: lbcf.tkestack.io/v1beta1
kind: LoadBalancer
metadata:
  name: example-of-create-new-lb 
  namespace: kube-system
spec:
  lbDriver: lbcf-clb-driver
  lbSpec:
    vpcID: vpc-b5hcoxj4
    loadBalancerType: "OPEN"
    listenerPort: "9999"
    listenerProtocol: "HTTP"
    domain: "mytest.com"
    url: "/index.html"
  ensurePolicy:
    policy: Always
```

###### **设定 backend 权重**

本例展示了 Service NodePort 的绑定。被绑定 Service 的名称为 svc-test，service port 为80（TCP)，绑定到 CLB 的每个`Node:NodePort`的权重都是 66

```
apiVersion: lbcf.tkestack.io/v1beta1
kind: BackendGroup
metadata:
  name: web-svc-backend-group
  namespace: kube-system
spec:
  lbName: test-clb-load-balancer
  service:
    name: svc-test
    port:
      portNumber: 80
  parameters:
    weight: "66"
```

## 附录

### 腾讯云 CLB LBCF driver

#### ConfigMap

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: trusted-tencentcloudapi
  namespace: kube-system
data:
  tencentcloudapi.pem: |
    -----BEGIN CERTIFICATE-----
   
   HERE IS YOUR CERTIFICATE

    -----END CERTIFICATE-----
```

#### Deployment

```yaml
apiVersion: lbcf.tkestack.io/v1beta1
kind: LoadBalancerDriver
metadata:
  name: lbcf-clb-driver
  namespace: kube-system
spec:
  driverType: Webhook
  url: "http://lbcf-clb-driver.kube-system.svc"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: lbcf-clb-driver
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      lbcf.tkestack.io/component: lbcf-clb-driver
  template:
    metadata:
      labels:
        lbcf.tkestack.io/component: lbcf-clb-driver
    spec:
      priorityClassName: "system-node-critical"
      containers:
        - name: driver
          image: ${image-name}
          args:
            - "--region=${your-region}"
            - "--vpc-id=${your-vpc-id}"
            - "--secret-id=${your-account-secret-id}"
            - "--secret-key=${your-account-secret-key}"
          ports:
            - containerPort: 80
              name: insecure
          imagePullPolicy: Always
          volumeMounts:
            - name: trusted-ca
              mountPath: /etc/ssl/certs
              readOnly: true
      volumes:
        - name: trusted-ca
          configMap:
            name: trusted-tencentcloudapi
```

#### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
  name: lbcf-clb-driver
  namespace: kube-system
spec:
  ports:
    - name: insecure
      port: 80
      targetPort: 80
  selector:
    lbcf.tkestack.io/component: lbcf-clb-driver
  sessionAffinity: None
  type: ClusterIP
```
