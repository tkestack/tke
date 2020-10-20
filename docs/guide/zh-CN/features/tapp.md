# TAPP
## TAPP 介绍

Kubernetes 现有应用类型（如：Deployment、StatefulSet等）无法满足很多非微服务应用的需求。比如：操作（升级、停止等）应用中的指定 Pod；应用支持多版本的 Pod。如果要将这些应用改造为适合于这些 Workload 的应用，需要花费很大精力，这将使大多数用户望而却步。

为解决上述复杂应用管理场景，TKEStack 基于 Kubernetes CRD 开发了一种新的应用类型 TAPP，它是一种通用类型的 Workload，同时支持 service 和 batch 类型作业，满足绝大部分应用场景，它能让用户更好的将应用迁移到 Kubernetes 集群。

![tapp picture](../../../images/tapp.png)



### TApp 功能点

功能点 | Deployment | StatefulSet | TAPP
---------------|-------|--------|--------
Pod 唯一性 | 无 | 每个Pod有唯一标识 | 每个Pod有唯一标识
Pod 存储独占 | 仅支持单容器 | 支持 | 支持
存储随Pod迁移 | 不支持 | 支持 | 支持
自动扩缩容 | 支持 | 不支持 | 支持
批量升级 | 支持 | 不支持 | 支持
严格顺序更新 | 不支持 | 支持 | 不支持
自动迁移问题节点 | 支持 | 不支持 | 支持
多版本管理 | 同时只有1个版本 | 可保持2个版本 | 可保持多个版本
Pod原地升级 | 不支持 | 不支持 | 支持

### TApp 运营场景

| 运营场景                       | Workload                                                     |                                                              |                                                              |
| ------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| Deployment                     | Statefulset                                                  | Tapp                                                         |                                                              |
| 发布                           | 并行创建                                                     | 顺序创建各个实例                                             | 并行创建                                                     |
| 扩容/缩容                      | 通过修改instance个数实现扩缩容。但是新增和裁撤的实例无法控制，如无法实现对某个机器上的实例的缩容。 | 只能顺序扩容或者缩容                                         | 扩容时，可以预期新增的实例的pod id，并且可以并发扩容在缩容操作时，可以指定缩容后的实例数，也可以对实例做选择性的缩容操作，即可以任意指定某一个或者一批实例缩掉。 |
| 升级/回退                      | 可以支持滚动升级                                             | 可以实现灰度升级，但需要按照序号升级，同时，也可以支持分段更新，但不能对单个实例操作。 | 除了支持deployment和statefulset类似的升级模式，还可以完全实现灰度升级/回退，即支持真正的灰度升级操作：即允许对个别实例（或者全部实例），指定目标版本的做升级操作，是否对其他实例继续升级或者升级到其他版本，可以由用户选择， |
| 负载均衡（容器平台操直接操作） | 整体绑定或解绑                                               | 整体绑定或解绑                                               | 除了整体绑定和解绑之外，还可以针对一个或多个实例做绑定和解绑操作，方便与其他运维操作结合，实现更稳妥的升级，扩缩容等。 |
| 容灾                           | 保持实例数，但是无法确定“新旧”实例的对应关系。               | 可以保持实例id的不变。                                       | 当由于本地重试或者迁移机器来实现容灾时，可以保证实例id的不变，方便对应，以及对事件、监控等持续性跟踪。 |

如果用Kubernetes的应用类型类比，TAPP ≈ Deployment + StatefulSet + Job ，它包含了Deployment、StatefulSet、Job的绝大部分功能，同时也有自己的特性，并且和原生Kubernetes相同的使用方式完全一致。


1. 实例具有可以标识的 ID

   实例有了 ID，业务就可以将很多状态或者配置逻辑和该 ID 做关联，当容器迁移时，通过 TAPP 的容器实例标识，可以识别该容器原来对应的数据，实现带云硬盘或者数据目录迁移 

1. 每个实例可以绑定自己的存储

   通过 TAPP 的容器实例标识，能很好地支持有状态的作业。在实例发生跨机迁移时，云硬盘能跟随实例一起迁移

1. 实现真正的灰度升级/回退

   Kubernetes 中的灰度升级概念应为滚动升级，kubernetes 将 pod”逐个”的更新，但现实中多业务需要的是稳定的灰度，即同一个 app，需要有多个版本同时稳定长时间的存在，TAPP 解决了此类问题

1. 可以指定实例 ID 做删除、停止、重启等操作

   对于微服务 app 来说，由于没有固定 ID，因此无法对某个实例做操作，而是全部交给了系统去维持足够的实例数

1. 对每个实例都有生命周期的跟踪

   对于一个实例，由于机器故障发生了迁移、重启等操作，很难跟踪和监控其生命周期。通过 TAPP 的容器实例标识，获得了实例真正的生命周期的跟踪，对于判断业务和系统是否正常服务具有特别重要的意义。TAPP 还可以记录事件，以及各个实例的运行时间，状态，高可用设置等。


### TAPP 资源结构

TApp 定义了一种用户自定义资源（CRD），TAPP controller 是 TAPP 对应的 controller/operator，它通过 kube-apiserve r监听 TApp、Pod 相关的事件，根据 TApp spec 和 status 进行相应的操作：创建、删除 Pod 等。

```go
// TApp represents a set of pods with consistent identities.
type TApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired identities of pods in this tapp.
	Spec TAppSpec `json:"spec,omitempty"`

	// Status is the current status of pods in this TApp. This data
	// may be out of date by some window of time.
	Status TAppStatus `json:"status,omitempty"`
}
// A TAppSpec is the specification of a TApp.
type TAppSpec struct {
	// Replicas 指定Template的副本数，尽管共享同一个Template定义，但是每个副本仍有唯一的标识
	Replicas int32 `json:"replicas,omitempty"`

	// 同Deployment的定义，标签选择器，默认为Pod Template上的标签
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// Template 默认模板，描述将要被初始创建/默认缩放的pod的对象，在TApp中可以被添加到TemplatePool中
	Template corev1.PodTemplateSpec `json:"template"`

	// TemplatePool 描述不同版本的pod template， template name --> pod Template
	TemplatePool map[string]corev1.PodTemplateSpec `json:"templatePool,omitempty"`

	// Statuses 用来指定对应pod实例的目标状态，instanceID --> desiredStatus ["Running","Killed"]
	Statuses map[string]InstanceStatus `json:"statuses,omitempty"`

	// Templates 用来指定运行pod实例所使用的Template，instanceID --> template name
	Templates map[string]string `json:"templates,omitempty"`

	// UpdateStrategy 定义滚动更新策略
	UpdateStrategy TAppUpdateStrategy `json:"updateStrategy,omitempty"`

	// ForceDeletePod 定义是否强制删除pod，默认为false
	ForceDeletePod bool `json:"forceDeletePod,omitempty"`

	// 同Statefulset的定义
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
}

// 滚动更新策略
type TAppUpdateStrategy struct {
	// 滚动更新的template name
	Template string `json:"template,omitempty"`
	// 滚动更新时的最大不可用数, 如果不指定此配置，滚动更新时不限制最大不可用数
	MaxUnavailable *int32 `json:"maxUnavailable,omitempty"`
}

// 定义TApp的状态
type TAppStatus struct {
	// most recent generation observed by controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Replicas 描述副本数
	Replicas int32 `json:"replicas"`

	// ReadyReplicas 描述Ready副本数
	ReadyReplicas int32 `json:"readyReplicas"`

	// ScaleSelector 是用于对pod进行查询的标签，它与HPA使用的副本计数匹配
	ScaleLabelSelector string `json:"scaleLabelSelector,omitempty"`

	// AppStatus 描述当前Tapp运行状态, 包含"Pending","Running","Failed","Succ","Killed"
	AppStatus AppStatus `json:"appStatus,omitempty"`

	// Statues 描述实例的运行状态 instanceID --> InstanceStatus ["NotCreated","Pending","Running","Updating","PodFailed","PodSucc","Killing","Killed","Failed","Succ","Unknown"]
	Statuses map[string]InstanceStatus `json:"statuses,omitempty"`
}
```

## 使用示例

本节以一个 TApp 应用部署，配置，升级，扩容以及杀死删除的操作步骤来说明 TApp 的使用。

### 创建 TApp 应用

创建 TApp 应用，副本数为3，TApp-controller 将根据默认模板创建出的 pod

```yaml
$ cat tapp.yaml
apiVersion: apps.tkestack.io/v1
kind: TApp
metadata:
  name: example-tapp
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: example-tapp
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
$ kubect apply -f tapp.yaml  
```

### 查看 TApp 应用

```shell
$ kubectl get tapp XXX
NAME           AGE
example-tapp   20m

$ kubectl descirbe tapp example-tapp
Name:         example-tapp
Namespace:    default
Labels:       app=example-tapp
Annotations:  <none>
API Version:  apps.tkestack.io/v1
Kind:         TApp
...
Spec:
...
Status:
  App Status:            Running
  Observed Generation:   2
  Ready Replicas:        3
  Replicas:              3
  Scale Label Selector:  app=example-tapp
  Statuses:
    0:  Running
    1:  Running
    2:  Running
Events:
  Type    Reason            Age   From             Message
  ----    ------            ----  ----             -------
  Normal  SuccessfulCreate  12m   tapp-controller  Instance: example-tapp-1
  Normal  SuccessfulCreate  12m   tapp-controller  Instance: example-tapp-0
  Normal  SuccessfulCreate  12m   tapp-controller  Instance: example-tapp-2
```

### 升级 TApp 应用

当前3个 pod 实例运行的镜像版本为 nginx:1.7.9，现在要升级其中的一个 pod 实例的镜像版本为 nginx:latest，在 spec.templatPools 中创建模板，然后在 spec.templates 中指定模板pod，指定“1”:“test”表示使用模板 test 创建 pod 1。

如果只更新镜像，Tapp controller 将对 pod 进行原地升级，即仅更新重启对应的容器，否则将按 k8s 原生方式删除 pod 并重新创建它们。

```yaml
apiVersion: apps.tkestack.io/v1
kind: TApp
metadata:
  name: example-tapp
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: example-tapp
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
  templatePool:
    "test":
      metadata:
        labels:
          app: example-tapp
      spec:
        containers:
        - name: nginx
          image: nginx:latest
  templates:
    "1": "test"
```

操作成功后，查看 instanceID 为'1'的 pod 已升级，镜像版本为 nginx:latest
```shell
# kubectl describe tapp example-tapp
Name:         example-tapp
Namespace:    default
Labels:       app=example-tapp
Annotations:  kubectl.kubernetes.io/last-applied-configuration:
                {"apiVersion":"apps.tkestack.io/v1","kind":"TApp","metadata":{"annotations":{},"name":"example-tapp","namespace":"default"},"spec":{"repli...
API Version:  apps.tkestack.io/v1
Kind:         TApp
...
Spec:
...
  Templates:
    1:  test
  Update Strategy:
Status:
  App Status:            Running
  Observed Generation:   4
  Ready Replicas:        3
  Replicas:              3
  Scale Label Selector:  app=example-tapp
  Statuses:
    0:  Running
    1:  Running
    2:  Running
Events:
  Type    Reason            Age   From             Message
  ----    ------            ----  ----             -------
  Normal  SuccessfulCreate  25m   tapp-controller  Instance: example-tapp-1
  Normal  SuccessfulCreate  25m   tapp-controller  Instance: example-tapp-0
  Normal  SuccessfulCreate  25m   tapp-controller  Instance: example-tapp-2
  Normal  SuccessfulUpdate  10m   tapp-controller  Instance: example-tapp-1

# kubectl get pod | grep example-tapp
example-tapp-0         1/1     Running   0          27m
example-tapp-1         1/1     Running   1          27m
example-tapp-2         1/1     Running   0          27m

# kubectl get pod example-tapp-1 -o template --template='{{range .spec.containers}}{{.image}}{{end}}'
nginx:latest

```

上述升级过程可根据实际需求灵活操作，可以指定多个pod的版本，帮助用户实现灵活的应用升级策略
同时可以指定 updateStrategy 升级策略，保证升级是最大不可用数为1，即保证滚动升级时每次仅更新和重启一个容器或 pod

```yaml
# cat tapp.yaml
apiVersion: apps.tkestack.io/v1
kind: TApp
metadata:
  name: example-tapp
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: example-tapp
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
  templatePool:
    "test":
      metadata:
        labels:
          app: example-tapp
      spec:
        containers:
        - name: nginx
          image: nginx:latest
  templates:
    "1": "test"
    "2": "test"
    "0": "test"
  updateStrategy:
    template: test
    maxUnavailable: 1
# kubectl apply -f tapp.yaml
```

### 杀死指定 pod

在 spec.statuses 中指定 pod 的状态，tapp-controller 根据用户指定的状态控制 pod实例，例如，如果 spec.statuses 为“1”:“killed”，tapp 控制器会杀死 pod 1。

```yaml
# cat tapp.yaml
kind: TApp
metadata:
  name: example-tapp
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: example-tapp
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
  templatePool:
    "test":
      metadata:
        labels:
          app: example-tapp
      spec:
        containers:
        - name: nginx
          image: nginx:latest
  templates:
    "1": "test"
    "2": "test"
    "0": "test"
  updateStrategy:
    template: test
    maxUnavailable: 1
  statuses:
    "1": "Killed"
# kubectl apply -f tapp.yaml
```

查看 pod 状态变为 Terminating

```shell
# kubectl get pod
NAME                   READY   STATUS        RESTARTS   AGE
example-tapp-0         1/1     Running       1          59m
example-tapp-1         0/1     Terminating   1          59m
example-tapp-2         1/1     Running       1          59m
```

### 扩容 TApp 应用

如果你想要扩展 TApp 使用默认的 spec.template 模板，只需增加 spec.replicas 的值，否则你需要在 spec.templates 中指定使用哪个模板。kubectl scale 也适用于 TApp。

```shell
kubectl scale --replicas=3 tapp/example-tapp

```

### 删除 TApp 应用

```shell
kubectl delete tapp example-tapp
```

### 其它

Tapp 还支持其他功能，如 HPA、volume templates，它们与 k8s 中的其它工作负载类型类似。