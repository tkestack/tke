# LBCF说明

## 组件介绍 : Load Balancer Controlling Framework (LBCF)

LBCF是一款部署在Kubernetes内的通用负载均衡控制面框架，旨在降低容器对接负载均衡的实现难度，并提供强大的扩展能力以满足业务方在使用负载均衡时的个性化需求。

## 部署在集群内kubernetes对象

在集群内部署LBCF Add-on , 将在集群内部署以下kubernetes对象

| kubernetes对象名称                                 | 类型                             | 默认占用资源 | 所属Namespaces |
| ---------------------------------------------- | ------------------------------ | ------ | ------------ |
| lbcf-controller                                | Deployment                     | /      | kube-system  |
| lbcf-controller                                | ServiceAccount                 | /      | kube-system  |
| lbcf-controller                                | ClusterRole                    | /      | /            |
| lbcf-controller                                | ClusterRoleBinding             | /      | /            |
| lbcf-controller                                | Secret                         | /      | kube-system  |
| lbcf-controller                                | Service                        | /      | kube-system  |
| backendrecords.lbcf.tkestack.io      | CustomResourceDefinition       | /      | /            |
| backendgroups.lbcf.tkestack.io       | CustomResourceDefinition       | /      | /            |
| loadbalancers.lbcf.tkestack.io       | CustomResourceDefinition       | /      | /            |
| loadbalancerdrivers.lbcf.tkestack.io | CustomResourceDefinition       | /      | /            |
| lbcf-mutate                                    | MutatingWebhookConfiguration   | /      | /            |
| lbcf-validate                                  | ValidatingWebhookConfiguration | /      | /            |

## LBCF使用场景

LBCF对K8S内部晦涩的运行机制进行了封装并以Webhook的形式对外暴露，在容器的全生命周期中提供了多达8种Webhook。通过实现这些Webhook，开发人员可以轻松实现下述功能：

- 对接任意负载均衡/名字服务，并自定义对接过程
- 实现自定义灰度升级策略
- 容器环境与其他环境共享同一个负载均衡
- 解耦负载均衡数据面与控制面

## LBCF使用方法

1. 通过扩展组件安装LBCF
1. 开发或选择安装LBCF Webhook规范的要求实现Webhook服务器
1. 以下按腾讯云CLB开发的webhook服务器为例

详细的使用方法和帮助文档，请参考[lb-controlling-framework](https://github.com/tkestack/lb-controlling-framework)文档

## 使用示例

### 使用已有四层CLB

本例中使用了id为`lb-7wf394rv`的负载均衡实例，监听器为四层监听器，端口号为20000，协议类型TCP。

*注: 程序会以`端口号20000，协议类型TCP`为条件查询监听器，若不存在，会自动创建新的*

```
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

### 创建新的七层CLB

本例在vpc  `vpc-b5hcoxj4`中创建了公网(OPEN)负载均衡实例，并为之创建了端口号为9999的HTTP监听器，最后会在监听器中创建`mytest.com/index.html`的转发规则

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

### 设定backend权重

本例展示了Service NodePort的绑定。被绑定Service的名称为svc-test，service port为80（TCP)，绑定到CLB的每个`Node:NodePort`的权重都是66

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

### 腾讯云CLB LBCF driver

ConfigMap：

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: trusted-tencentcloudapi
  namespace: kube-system
data:
  tencentcloudapi.pem: |
    -----BEGIN CERTIFICATE-----
    .............
    -----END CERTIFICATE-----
```

Deployment

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

Service:

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
