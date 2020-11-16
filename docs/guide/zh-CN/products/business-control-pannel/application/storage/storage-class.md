# StorageClass

**StorageClass**： 通过 PVC 请求到一定的存储空间也很有可能不足以满足应用对于存储设备的各种需求，而且不同的应用程序对于存储性能的要求可能也不尽相同，比如读写速度、并发性能等，为了解决这一问题，Kubernetes 又引入了一个新的资源对象：StorageClass，通过 StorageClass 的定义，管理员可以将存储资源定义为某种类型的资源，比如快速存储、慢速存储等，用户根据 StorageClass 的描述就可以非常直观的知道各种存储资源的具体特性了，这样就可以根据应用的特性去申请合适的存储资源了。StorageClass 用于描述存储的类型，集群管理员可以为集群定义不同的存储类别。可通过  StorageClass 配合 PersistentVolumeClaim 可以动态创建需要的存储资源。

> TKEStack 没有提供存储服务，Global 集群中的镜像仓库、ETCD、InfluxDB 等数据组件，均使用**本地磁盘存储数据**。如果您需要使用存储服务，建议使用 [ROOK](https://rook.io/) 或者 [chubaoFS](https://chubao.io/)，部署一套容器化的分布式存储服务。
>
> 您可阅读 [在 TKEStack 上使用存储的最佳实践](../../../../best-practices/storage.md) 作为参考。

## 创建

要使用 StorageClass，就得安装对应的自动配置程序，比如这里存储后端使用的是 NFS，那么就需要使用到一个 nfs-client 的自动配置程序，也可以叫称之为 Provisioner，这个程序使用已经配置好的 NFS 服务器来自动创建持久卷，也就是自动帮创建 PV。

- 自动创建的 PV 以 `${namespace}-${pvcName}-${pvName}` 这样的命名格式创建在 NFS 服务器上的共享数据目录中
- 而当这个 PV 被回收后会以 `archieved-${namespace}-${pvcName}-${pvName}` 这样的命名格式存在 NFS 服务器上

当然在部署 `nfs-client` 之前，需要先参照上一节安装 [NFS 服务器](pv&pvc.md#安装 NFS server（服务端）)，然后部署 nfs-client 即可，也可以直接参考 nfs-client 的文档：https://github.com/kubernetes-incubator/external-storage/tree/master/nfs-client，进行安装。

1. 配置 Deployment，将里面的对应的参数替换成自己的 NFS 配置（nfs-client.yaml）

  ```yaml
  kind: Deployment
  apiVersion: apps/v1
  metadata:
    name: nfs-client-provisioner
  spec:
    replicas: 1
    strategy:
      type: Recreate
    selector:
      matchLabels:
        app: nfs-client-provisioner
    template:
      metadata:
        labels:
          app: nfs-client-provisioner
      spec:
        serviceAccountName: nfs-client-provisioner
        containers:
          - name: nfs-client-provisioner
            image: quay.io/external_storage/nfs-client-provisioner:latest
            volumeMounts:
              - name: nfs-client-root
                mountPath: /persistentvolumes
            env:
              - name: PROVISIONER_NAME
                value: fuseim.pri/ifs
              - name: NFS_SERVER
                value: 42.194.158.74
              - name: NFS_PATH
                value: /data/k8s
        volumes:
          - name: nfs-client-root
            nfs:
              server: 42.194.158.74
              path: /data/k8s
  ```

2. 将环境变量 NFS_SERVER 和 NFS_PATH 替换，当然也包括下面的 NFS 配置，可以看到这里使用了一个名为 nfs-client-provisioner 的`serviceAccount`，所以也需要创建一个 SA，然后绑定上对应的权限：（nfs-client-sa.yaml）

  ```yaml
  apiVersion: v1
  kind: ServiceAccount
  metadata:
    name: nfs-client-provisioner

  ---
  kind: ClusterRole
  apiVersion: rbac.authorization.k8s.io/v1
  metadata:
    name: nfs-client-provisioner-runner
  rules:
    - apiGroups: [""]
      resources: ["persistentvolumes"]
      verbs: ["get", "list", "watch", "create", "delete"]
    - apiGroups: [""]
      resources: ["persistentvolumeclaims"]
      verbs: ["get", "list", "watch", "update"]
    - apiGroups: ["storage.k8s.io"]
      resources: ["storageclasses"]
      verbs: ["get", "list", "watch"]
    - apiGroups: [""]
      resources: ["events"]
      verbs: ["list", "watch", "create", "update", "patch"]
    - apiGroups: [""]
      resources: ["endpoints"]
      verbs: ["create", "delete", "get", "list", "watch", "patch", "update"]

  ---
  kind: ClusterRoleBinding
  apiVersion: rbac.authorization.k8s.io/v1
  metadata:
    name: run-nfs-client-provisioner
  subjects:
    - kind: ServiceAccount
      name: nfs-client-provisioner
      namespace: default
  roleRef:
    kind: ClusterRole
    name: nfs-client-provisioner-runner
    apiGroup: rbac.authorization.k8s.io
  ```

  > 这里新建的一个名为 nfs-client-provisioner 的 `ServiceAccount`，然后绑定了一个名为 nfs-client-provisioner-runner 的 `ClusterRole`，而该 `ClusterRole` 声明了一些权限，其中就包括对 `persistentvolumes` 的增、删、改、查等权限，所以可以利用该 `ServiceAccount` 来自动创建 PV。

3. nfs-client 的 Deployment 声明完成后，就可以来创建一个`StorageClass`对象了：（nfs-client-class.yaml）

  ```yaml
  apiVersion: storage.k8s.io/v1
  kind: StorageClass
  metadata:
    name: course-nfs-storage
  provisioner: fuseim.pri/ifs # or choose another name, must match deployment's env PROVISIONER_NAME'
  ```

  > 声明了一个名为 course-nfs-storage 的`StorageClass`对象，注意下面的`provisioner`对应的值一定要和上面的`Deployment`下面的 PROVISIONER_NAME 这个环境变量的值一样。

4. 创建这些资源对象吧：

  ```shell
  $ kubectl create -f nfs-client.yaml
  $ kubectl create -f nfs-client-sa.yaml
  $ kubectl create -f nfs-client-class.yaml
  ```

5. 创建完成后查看下资源状态：

  ```
  $ kubectl get pods
  NAME                                             READY     STATUS             RESTARTS   AGE
  ...
  nfs-client-provisioner-7648b664bc-7f9pk          1/1       Running            0          7h
  ...
  $ kubectl get storageclass
  NAME                 PROVISIONER      AGE
  course-nfs-storage   fuseim.pri/ifs   11s
  ```

## 新建

上面把`StorageClass`资源对象创建成功了，接下来通过一个示例测试下动态 PV，首先创建一个 PVC 对象：(test-pvc.yaml)

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: test-pvc
spec:
  accessModes:
  - ReadWriteMany
  resources:
    requests:
      storage: 1Mi
```

这里声明了一个 PVC 对象，采用 ReadWriteMany 的访问模式，请求 1Mi 的空间，但是可以看到上面的 PVC 文件我们没有标识出任何和 StorageClass 相关联的信息，那么如果我们现在直接创创建一个合适的 PV:

- 第一种方法：在这个 PVC 对象中添加一个声明 StorageClass 对象的标识，这里可以利用一个 annotations 属性来标识，如下：

  ```yaml
  apiVersion: v1
  kind: PersistentVolumeClaim
  metadata:
    name: test-pvc
    annotations:
      volume.beta.kubernetes.io/storage-class: "course-nfs-storage"
  spec:
    accessModes:
    - ReadWriteMany
    resources:
      requests:
        storage: 1Mi
  ```

- 第二种方法：可以设置这个 course-nfs-storage 的 StorageClass 为 Kubernetes 的默认存储后端，可以用 kubectl patch 命令来更新：

  ```shell
  $ kubectl patch storageclass course-nfs-storage -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
  ```

上面这两种方法都是可以的，当然为了不影响系统的默认行为，这里还是采用第一种方法，直接创建即可：

  ```shell
  $ kubectl create -f test-pvc.yaml
  persistentvolumeclaim "test-pvc" created
  $ kubectl get pvc
  NAME         STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS          AGE
  ...
  test-pvc     Bound     pvc-73b5ffd2-8b4b-11e8-b585-525400db4df7   1Mi        RWX            course-nfs-storage    2m
  ...
  ```

可以看到一个名为 test-pvc 的 PVC 对象创建成功了，状态已经是 Bound 了，是不是也产生了一个对应的 VOLUME 对象，最重要的一栏是 STORAGECLASS，现在是不是也有值了，就是刚刚创建的 StorageClass 对象 course-nfs-storage。

然后查看下 PV 对象：

  ```shell
  $ kubectl get pv
  NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM                STORAGECLASS          REASON    AGE
  ...
  pvc-73b5ffd2-8b4b-11e8-b585-525400db4df7   1Mi        RWX            Delete           Bound       default/test-pvc     course-nfs-storage              8m
  ...
  ```

可以看到是不是自动生成了一个关联的 PV 对象，访问模式是 RWX，回收策略是 Delete，这个 PV 对象并不是我们手动创建的吧，这是通过上面的 StorageClass 对象自动创建的。这就是 StorageClass 的创建方法。

## 测试

接下来还是用一个简单的示例来测试下上面用 StorageClass 方式声明的 PVC 对象吧：(test-pod.yaml)

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: test-pod
spec:
  containers:
  - name: test-pod
    image: busybox
    imagePullPolicy: IfNotPresent
    command:
    - "/bin/sh"
    args:
    - "-c"
    - "touch /mnt/SUCCESS && exit 0 || exit 1"
    volumeMounts:
    - name: nfs-pvc
      mountPath: "/mnt"
  restartPolicy: "Never"
  volumes:
  - name: nfs-pvc
    persistentVolumeClaim:
      claimName: test-pvc
```

上面这个 Pod 非常简单，就是用一个 busybox 容器，在 /mnt 目录下面新建一个 SUCCESS 的文件，然后把 /mnt 目录挂载到上面新建的 test-pvc 这个资源对象上面了，要验证很简单，只需要去查看下 NFS 服务器上面的共享数据目录下面是否有 SUCCESS 这个文件即可：

```shell
$ kubectl create -f test-pod.yaml
pod "test-pod" created
```

然后可以在 NFS 服务器的共享数据目录下面查看下数据：

```shell
$ ls /data/k8s/
default-test-pvc-pvc-73b5ffd2-8b4b-11e8-b585-525400db4df7
```

可以看到下面有名字很长的文件夹，这个文件夹的命名方式是不是和上面的规则：**${namespace}-${pvcName}-${pvName}** 是一样的，再看下这个文件夹下面是否有其他文件：

```shell
$ ls /data/k8s/default-test-pvc-pvc-73b5ffd2-8b4b-11e8-b585-525400db4df7
SUCCESS
```

可以看到下面有一个 SUCCESS 的文件，就证明上面的验证是成功的。

另外可以看到这里是手动创建的一个 PVC 对象，在实际工作中，使用 StorageClass 更多的是 StatefulSet 类型的服务，StatefulSet 类型的服务也可以通过一个 volumeClaimTemplates 属性来直接使用 StorageClass，如下：(test-statefulset-nfs.yaml)

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: nfs-web
spec:
  serviceName: "nginx"
  replicas: 3
  selector:
    matchLabels:
      app: nfs-web
  template:
    metadata:
      labels:
        app: nfs-web
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
          name: web
        volumeMounts:
        - name: www
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: www
      annotations:
        volume.beta.kubernetes.io/storage-class: course-nfs-storage
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

实际上 volumeClaimTemplates 下面就是一个 PVC 对象的模板，就类似于这里 StatefulSet 下面的 template，实际上就是一个 Pod 的模板，不单独创建成 PVC 对象，而用这种模板就可以动态的去创建了对象了，这种方式在 StatefulSet 类型的服务下面使用得非常多。

直接创建上面的对象：

```shell
$ kubectl create -f test-statefulset-nfs.yaml
statefulset.apps "nfs-web" created
$ kubectl get pods
NAME                                             READY     STATUS              RESTARTS   AGE
...
nfs-web-0                                        1/1       Running             0          1m
nfs-web-1                                        1/1       Running             0          1m
nfs-web-2                                        1/1       Running             0          33s
...
```

创建完成后可以看到上面的3个 Pod 已经运行成功，然后查看下 PVC 对象：

```shell
$ kubectl get pvc
NAME            STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS          AGE
...
www-nfs-web-0   Bound     pvc-cc36b3ce-8b50-11e8-b585-525400db4df7   1Gi        RWO            course-nfs-storage    2m
www-nfs-web-1   Bound     pvc-d38285f9-8b50-11e8-b585-525400db4df7   1Gi        RWO            course-nfs-storage    2m
www-nfs-web-2   Bound     pvc-e348250b-8b50-11e8-b585-525400db4df7   1Gi        RWO            course-nfs-storage    1m
...
```

可以看到是不是也生成了3个 PVC 对象，名称由模板名称 name 加上 Pod 的名称组合而成，这3个 PVC 对象也都是 绑定状态了，很显然查看 PV 也可以看到对应的3个 PV 对象：

```shell
$ kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM                   STORAGECLASS          REASON    AGE
...                                                        1d
pvc-cc36b3ce-8b50-11e8-b585-525400db4df7   1Gi        RWO            Delete           Bound       default/www-nfs-web-0   course-nfs-storage              4m
pvc-d38285f9-8b50-11e8-b585-525400db4df7   1Gi        RWO            Delete           Bound       default/www-nfs-web-1   course-nfs-storage              4m
pvc-e348250b-8b50-11e8-b585-525400db4df7   1Gi        RWO            Delete           Bound       default/www-nfs-web-2   course-nfs-storage              4m
...
```

查看 NFS 服务器上面的共享数据目录：

```shell
$ ls /data/k8s/
...
default-www-nfs-web-0-pvc-cc36b3ce-8b50-11e8-b585-525400db4df7
default-www-nfs-web-1-pvc-d38285f9-8b50-11e8-b585-525400db4df7
default-www-nfs-web-2-pvc-e348250b-8b50-11e8-b585-525400db4df7
...
```