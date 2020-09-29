# HPA

HPA会基于cpu、内存等指标对负载的pod数量动态调控，达到工作负载稳定的目的。

依赖：metrics-server（当前global集群自带metrics-server，导入集群需要检查其是否安装）

## 安装依赖

### 安装metrics-server

Kubernetes Metrics Server是一个集群范围的资源使用数据聚合器，是Heapster的继承者。metrics服务器通过从kubernet.summary_api收集数据收集节点和pod的CPU和内存使用情况。summary API是一个内存有效的API，用于将数据从Kubelet/cAdvisor传递到metrics server，下图为HPA和kubectl等调用metrics-server获取相关信息的原理图。

![image-20200929172542934](../../../../../../images/image-20200929172542934.png)

metrics-server yaml参考 https://github.com/kubernetes-sigs/metrics-server/releases 如下 components.yaml

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:aggregated-metrics-reader
  labels:
    rbac.authorization.k8s.io/aggregate-to-view: "true"
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
rules:
- apiGroups: ["metrics.k8s.io"]
  resources: ["pods", "nodes"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: metrics-server:system:auth-delegator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: metrics-server-auth-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
---
apiVersion: apiregistration.k8s.io/v1beta1
kind: APIService
metadata:
  name: v1beta1.metrics.k8s.io
spec:
  service:
    name: metrics-server
    namespace: kube-system
  group: metrics.k8s.io
  version: v1beta1
  insecureSkipTLSVerify: true
  groupPriorityMinimum: 100
  versionPriority: 100
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: metrics-server
  namespace: kube-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: metrics-server
  namespace: kube-system
  labels:
    k8s-app: metrics-server
spec:
  selector:
    matchLabels:
      k8s-app: metrics-server
  template:
    metadata:
      name: metrics-server
      labels:
        k8s-app: metrics-server
    spec:
      serviceAccountName: metrics-server
      volumes:
      # mount in tmp so we can safely use from-scratch images and/or read-only containers
      - name: tmp-dir
        emptyDir: {}
      containers:
      - name: metrics-server
        image: mirrors.tencent.com/capdtke/metrics-server-amd64:v0.3.6
        imagePullPolicy: IfNotPresent
        args:
          - --cert-dir=/tmp
          - --secure-port=4443
        command:
          - /metrics-server
          - --kubelet-preferred-address-types=InternalIP
          - --kubelet-insecure-tls
        ports:
        - name: main-port
          containerPort: 4443
          protocol: TCP
        securityContext:
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 1000
        volumeMounts:
        - name: tmp-dir
          mountPath: /tmp
      nodeSelector:
        user: master
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                    - prometheus
              topologyKey: "kubernetes.io/hostname"
---
apiVersion: v1
kind: Service
metadata:
  name: metrics-server
  namespace: kube-system
  labels:
    kubernetes.io/name: "Metrics-server"
    kubernetes.io/cluster-service: "true"
spec:
  selector:
    k8s-app: metrics-server
  ports:
  - port: 443
    protocol: TCP
    targetPort: main-port
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:metrics-server
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - nodes
  - nodes/stats
  - namespaces
  - configmaps
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:metrics-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:metrics-server
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
```

具体请安装配置参考metrics-server git地址 https://github.com/kubernetes-sigs/metrics-server

## 使用HPA

TKEStack已经支持在页面多处位置为负载配置HPA

1. 新建负载页（负载包括Deployment，StatefulSet，TApp）这里新建负载时将会同时新建与负载同名的HPA对象。

![image-20200929173056091](../../../../../../images/image-20200929173056091.png)

2. 负载列表页（负载包括Deployment，StatefulSet，TApp）

   ![image-20200929173209190](../../../../../../../../../../Typora/images/image-20200929173209190.png)

   * 点击“更新实例数量”，进入配置界面如图所示，这里将会同时新建与负载同名的HPA对象。

     ![image-20200929173300650](../../../../../../images/image-20200929173300650.png)

3. 自动伸缩的HPA列表页。此处可以查看/修改/新建HPA。

   ![image-20200929173933713](../../../../../../images/image-20200929173933713.png)

   * 点击上图中的新建，出现新建HPA页面，如下图所示。

   ![image-20200929173834852](../../../../../../images/image-20200929173834852.png)