## LBCF说明

### 组件介绍 : Load Balancer Controlling Framework (LBCF)

LBCF是一款部署在Kubernetes内的通用负载均衡控制面框架，旨在降低容器对接负载均衡的实现难度，并提供强大的扩展能力以满足业务方在使用负载均衡时的个性化需求。

### 部署在集群内kubernetes对象

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

### LBCF使用场景

LBCF对K8S内部晦涩的运行机制进行了封装并以Webhook的形式对外暴露，在容器的全生命周期中提供了多达8种Webhook。通过实现这些Webhook，开发人员可以轻松实现下述功能：

- 对接任意负载均衡/名字服务，并自定义对接过程
- 实现自定义灰度升级策略
- 容器环境与其他环境共享同一个负载均衡
- 解耦负载均衡数据面与控制面

### LBCF限制条件

#### 系统要求：

- K8S 1.10及以上版本
- 开启Dynamic Admission Control，在apiserver中添加启动参数：
  - --enable-admission-plugins=MutatingAdmissionWebhook,ValidatingAdmissionWebhook
- K8S 1.10版本，在apiserver中额外添加参数：
  - --feature-gates=CustomResourceSubresources=true

推荐环境：

在[腾讯云](https://cloud.tencent.com/product/tke)上购买1.12.4版本集群，无需修改任何参数，开箱可用

## LBCF使用方法

1. 通过扩展组件安装LBCF

2. 开发或选择安装LBCF Webhook规范的要求实现Webhook服务器

3. 以下按腾讯云CLB开发的webhook服务器为例

#### LBCF CLB driver

##### 功能列表

- 使用已有负载均衡
- 创建新的负载均衡（四层/七层）
- 绑定Service NodePort
- CLB直通POD(直接绑定Pod至CLB，不通过Service）
- 权重调整
- 能够校验并拒绝非法参数

##### 部署LBCF CLB driver

部署前需修改YAML（文中附录已提供yaml文件，需要向deploy.yaml中填入以下信息

- 镜像信息
- 所在地域
- 所在vpcID （绑定service NodePort时用来查找节点对应的instanceID）
- secret-id
- secret-key

```
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
```

登陆集群，使用以下命令安装YAML

```
kubectl apply -f configmap.yaml
kubectl apply -f deploy.yaml
kubectl apply -f service.yaml
```

##### 使用示例

**使用已有四层CLB**

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

**创建新的七层CLB**

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

**设定backend权重**

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

### 附录

#### 腾讯云CLB LBCF driver

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
    MIIGgDCCBWigAwIBAgIMIsk10+Y2GGXUKTrmMA0GCSqGSIb3DQEBCwUAMGYxCzAJ
    BgNVBAYTAkJFMRkwFwYDVQQKExBHbG9iYWxTaWduIG52LXNhMTwwOgYDVQQDEzNH
    bG9iYWxTaWduIE9yZ2FuaXphdGlvbiBWYWxpZGF0aW9uIENBIC0gU0hBMjU2IC0g
    RzIwHhcNMTgxMjIwMDcxNTA5WhcNMTkxMjIxMDcxNTA5WjCBjDELMAkGA1UEBhMC
    Q04xEjAQBgNVBAgTCWd1YW5nZG9uZzERMA8GA1UEBxMIc2hlbnpoZW4xNjA0BgNV
    BAoTLVRlbmNlbnQgVGVjaG5vbG9neSAoU2hlbnpoZW4pIENvbXBhbnkgTGltaXRl
    ZDEeMBwGA1UEAwwVKi50ZW5jZW50Y2xvdWRhcGkuY29tMIIBIjANBgkqhkiG9w0B
    AQEFAAOCAQ8AMIIBCgKCAQEAyepbdY0laI2rgfm1qe4TUv0kR9r0YJQwTwXN3LM6
    2W75Y5m2k9WcfFilcoah9q4J1ndkbtiDSaLRYJYce7ivObmR79gb4MGCrnVix0eI
    KYW1qiFIBjETxhzTZt4sVztty4LW0F+R4lggraAP8d7qdsbFTyk4y9dKHS1FANQc
    xVkxdFIMCk+WoppMmTI2DpNg9kY6BrL7TiLyjx8gpF1XymKl0CefqYxwZt/+KEaA
    75G/R361h2wi5lFC1ybhGtlPT6t285A6j6avC7AqEhdZQqoAv60iQud2Hj7bmkbf
    OTgE24+5LepekWyK0iEDCHX8aN/wtfKPLqA3oQVlnLLlbQIDAQABo4IDBTCCAwEw
    DgYDVR0PAQH/BAQDAgWgMIGgBggrBgEFBQcBAQSBkzCBkDBNBggrBgEFBQcwAoZB
    aHR0cDovL3NlY3VyZS5nbG9iYWxzaWduLmNvbS9jYWNlcnQvZ3Nvcmdhbml6YXRp
    b252YWxzaGEyZzJyMS5jcnQwPwYIKwYBBQUHMAGGM2h0dHA6Ly9vY3NwMi5nbG9i
    YWxzaWduLmNvbS9nc29yZ2FuaXphdGlvbnZhbHNoYTJnMjBWBgNVHSAETzBNMEEG
    CSsGAQQBoDIBFDA0MDIGCCsGAQUFBwIBFiZodHRwczovL3d3dy5nbG9iYWxzaWdu
    LmNvbS9yZXBvc2l0b3J5LzAIBgZngQwBAgIwCQYDVR0TBAIwADBJBgNVHR8EQjBA
    MD6gPKA6hjhodHRwOi8vY3JsLmdsb2JhbHNpZ24uY29tL2dzL2dzb3JnYW5pemF0
    aW9udmFsc2hhMmcyLmNybDA1BgNVHREELjAsghUqLnRlbmNlbnRjbG91ZGFwaS5j
    b22CE3RlbmNlbnRjbG91ZGFwaS5jb20wHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsG
    AQUFBwMCMB0GA1UdDgQWBBRR63cbhz8Aloch9ZEw6Y4TZKXspjAfBgNVHSMEGDAW
    gBSW3mHxvRwWKVMcwMx9O4MAQOYafDCCAQYGCisGAQQB1nkCBAIEgfcEgfQA8gB3
    AFWB1MIWkDYBSuoLm1c8U/DA5Dh4cCUIFy+jqh0HE9MMAAABZ8p32P4AAAQDAEgw
    RgIhAOBSocmwefb43lFbW9CVd9Kx6P6o35YLoXR5YO6vae2bAiEAwVSFT6xIb7wG
    mQqVwUItRUG9LtqjuQNhfMkhPiCV3zsAdwDuS723dc5guuFCaR+r4Z5mow9+X7By
    2IMAxHuJeqj9ywAAAWfKd9YUAAAEAwBIMEYCIQC8txn4L1STQ9ai4JcWJ6vwNoc4
    5tFfQsKXDGs4CXHaUgIhAJ7PTTgajS5A9xTvTdD0Tw3iH643MjjdLTKH83Qdu2ty
    MA0GCSqGSIb3DQEBCwUAA4IBAQDFu2JcyLG3Bg8YhJi+IqoTljwGsYC98i148hoT
    CwlbwH3UaHrPlR1crX8Hv+XEsHj4Ot3/krdiuYGWEZVhY61e8DT3QovUTXh6pvG+
    R9Q22SfGMuGuwrgTdhfR5QYv4whE/Mj4TqJQDRGBetb9dpPkhhLN6E+h/9/WmGyC
    HObUUZyEP1rTqgPxLk8e5Xyt8yv/loo5eptQXvduY/v4ngpvJScqepDXedSJd3SK
    Muu7gepolidg/fBlZjfpksLWSUGVVuVUS4zT2gaMpTqD/NjxHwC3roIiP9pSnY7w
    GcPWfnp6Xs8ahCmiYdOEwbrMH/QtEdRdonsPyMS2FU3Rv7hD
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
