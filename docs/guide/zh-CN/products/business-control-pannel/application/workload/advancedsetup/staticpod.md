# 静态 Pod

在 Kubernetes 中除了我们经常使用到的普通的 Pod 外，还有一种特殊的 Pod，叫做 `Static Pod`，即静态 Pod，静态 Pod 有什么特殊的地方呢？

静态 Pod 直接由特定节点上的 `kubelet` 进程来管理，不通过 Master 节点上的 `APIServer`。无法与我们常用的控制器 `Deployment` 或者 `DaemonSet` 进行关联，它由 `kubelet` 进程自己来监控，当 `Pod` 崩溃时重启该 `Pod`，`kubelete` 也无法对他们进行健康检查。静态 Pod 始终绑定在某一个 `kubelet`，并且始终运行在同一个节点上。

> `kubelet` 会自动为每一个静态 Pod 在 Kubernetes 的 APIServer 上创建一个镜像 Pod（Mirror Pod），因此我们可以在 APIServer 中查询到该 Pod，但是不能通过 APIServer 进行控制（例如不能删除）。

创建静态 Pod 有两种方式：配置文件和 HTTP

## 配置文件

配置文件就是放在特定目录下的标准的 JSON 或 YAML 格式的 Pod 定义文件。用 `kubelet --pod-manifest-path=<the directory>` 来启动 `kubelet` 进程， kubelet 定期的去扫描这个目录，根据这个目录下出现或消失的 YAML/JSON 文件来创建或删除静态 Pod。

比如我们在 Master 节点上用静态 pod 的方式来启动一个 nginx 的服务。我们登录到一个 Master 节点上面，可以通过下面命令找到 kubelet 对应的启动配置文件

```shell
$ systemctl status kubelet
```

![image-20201103153500534](../../../../../../../images/image-20201103153500534.png)

配置文件路径为：

```shell
$ /usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf
```

打开这个文件我们可以看到其中有一条如下的环境变量配置：`Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"`

所以如果我们通过 `kubeadm` 的方式来安装的集群环境，对应的 `kubelet` 已经配置了我们的静态 Pod 文件的路径，那就是 `/etc/kubernetes/manifests`，所以我们只需要在该目录下面创建一个标准的 Pod 的 JSON 或者 YAML 文件即可：

如果你的 kubelet 启动参数中没有配置上面的`--pod-manifest-path`参数的话，那么**添加上这个参数然后重启 kubelet 即可**。

```yaml
$ cat <<EOF >/etc/kubernetes/manifest/static-web.yaml
apiVersion: v1
kind: Pod
metadata:
  name: static-web
  labels:
    app: static
spec:
  containers:
    - name: web
      image: nginx
      ports:
        - name: web
          containerPort: 80
EOF
```

## 通过 HTTP 创建静态 Pods

kubelet 周期地从 `–manifest-url=` 参数指定的地址下载文件，并且把它翻译成 JSON/YAML 格式的 Pod 定义。此后的操作方式与 `–pod-manifest-path=` 相同，kubelet 会不时地重新下载该文件，当文件变化时对应地终止或启动静态 Pod。

## 静态 Pods 的动作行为

kubelet 启动时，由 `--pod-manifest-path= or --manifest-url=` 参数指定的目录下定义的所有 Pod 都会自动创建，例如，我们示例中的 static-web。

```shell
$ docker ps
CONTAINER ID IMAGE         COMMAND  CREATED        STATUS         PORTS     NAMES
f6d05272b57e nginx:latest  "nginx"  8 minutes ago  Up 8 minutes             k8s_web.6f802af4_static-web-fk-node1_default_67e24ed9466ba55986d120c867395f3c_378e5f3c
```

现在我们通过 `kubectl` 工具可以看到这里创建了一个新的镜像 Pod：

```shell
 $ kubectl get pods
NAME                       READY     STATUS    RESTARTS   AGE
static-web-my-node01        1/1       Running   0          2m
```

静态 Pod 的标签会传递给镜像 Pod，可以用来过滤或筛选。 需要注意的是，我们不能通过 API 服务器来删除静态 Pod（例如，通过 [kubectl](https://kubernetes.io/docs/user-guide/kubectl/) 命令），kebelet 不会删除它。

```shell
$ kubectl delete pod static-web-my-node01
$ kubectl get pods
NAME                       READY     STATUS    RESTARTS   AGE
static-web-my-node01        1/1       Running   0          12s
```

我们尝试手动终止容器，可以看到 kubelet 很快就会自动重启容器。

```shell
$ docker ps
CONTAINER ID        IMAGE         COMMAND                CREATED       ...
5b920cbaf8b1        nginx:latest  "nginx -g 'daemon of   2 seconds ago ...
```

## 静态 Pods 的动态增加和删除

运行中的 kubelet 周期扫描配置的目录（我们这个例子中就是  /etc/kubernetes/manifests ）下文件的变化，当这个目录中有文件出现或消失时创建或删除 Pods。

```shell
[root@node01 ~] $ mv /etc/kubernetes/manifests/static-web.yaml /tmp
[root@node01 ~] $ sleep 20
[root@node01 ~] $ docker ps
// no nginx container is running
[root@node01 ~] $ mv /tmp/static-web.yaml  /etc/kubernetes/manifests
[root@node01 ~] $ sleep 20
[root@node01 ~] $ docker ps
CONTAINER ID        IMAGE         COMMAND                CREATED           ...
e7a62e3427f1        nginx:latest  "nginx -g 'daemon of   27 seconds ago
```

其实我们用 kubeadm 安装的集群，Master 节点上面的几个重要组件都是用静态 Pod 的方式运行的，我们登录到 master 节点上查看`/etc/kubernetes/manifests`目录：

```shell
[root@master ~]# ls /etc/kubernetes/manifests/
etcd.yaml  keepalived.yaml  kube-apiserver.yaml  kube-controller-manager.yaml  kube-scheduler.yaml
```

这种方式也为我们将集群的一些组件容器化提供了可能，因为这些 Pod 都不会受到 APIServer 的控制，不然我们这里 `kube-apiserver` 怎么自己去控制自己呢？万一不小心把这个 Pod 删掉了呢？所以只能有 `kubelet` 自己来进行控制，这就是我们所说的静态 Pod。