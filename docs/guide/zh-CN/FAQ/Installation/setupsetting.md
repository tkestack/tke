# 修改组件的启动参数

## 修改 master 的 apiserver 的启动参数

TKEStack 是通过 Kubeadm 搭建集群的， Kubernetes apiserver 是由 static pod 启动，其 yaml 文件的位置在master节点的 `/etc/kubernetes/manifest/kube-apiserver.yaml` 这个路径下。如下所示，在`spec.containers.command`里的内容，从`- --advertise-address=10.0.222.102` 开始全是apiserver的启动配置。同样的，可以在`/etc/kubernetes/manifest`该路径下查看 master 节点其他组件的 yaml 文件。

```shell
# /etc/kubernetes/ 目录下是Kubernetes的相关配置
[root@VM-222-102-centos opt]# cd /etc/kubernetes/
[root@VM-222-102-centos kubernetes]# ls
admin.conf               kubeadm-config.yaml  pki                           tke-authz-webhook.yaml
controller-manager.conf  kubelet.conf         scheduler-policy-config.json
known_tokens.csv         manifests            scheduler.conf

# /etc/kubernetes/ 目录下是master组件的相关配置
[root@VM-222-102-centos kubernetes]# cd manifests/
[root@VM-222-102-centos manifests]# ls
etcd.yaml  kube-apiserver.yaml  kube-controller-manager.yaml  kube-scheduler.yaml
[root@VM-222-102-centos manifests]# cat kube-apiserver.yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    scheduler.alpha.kubernetes.io/critical-pod: ""
    tke.prometheus.io/scrape: "true"
    prometheus.io/scheme: "https"
    prometheus.io/port: "6443"
  annotations:
    kubeadm.kubernetes.io/kube-apiserver.advertise-address.endpoint: 10.0.222.102:6443
  creationTimestamp: null
  labels:
    component: kube-apiserver
    tier: control-plane
  name: kube-apiserver
  namespace: kube-system
spec:
  containers:
  - command:
    - kube-apiserver
    # 以下全是 kube-apiserver 的启动配置
    - --advertise-address=10.0.222.102
    - --allow-privileged=true
    - --authorization-mode=Node,RBAC,Webhook
    - --authorization-webhook-config-file=/etc/kubernetes/tke-authz-webhook.yaml
    - --client-ca-file=/etc/kubernetes/pki/ca.crt
    - --enable-admission-plugins=NodeRestriction
    - --enable-bootstrap-token-auth=true
    - --etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt
    - --etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt
    - --etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key
    - --etcd-servers=https://127.0.0.1:2379
    - --insecure-port=0
    ....
```

Static Pod 的配置文件被修改后，立即生效。

- Kubelet 会监听该文件的变化，当您修改了 `/etc/kubenetes/manifest/kube-apiserver.yaml` 文件之后，kubelet 将自动终止原有的 kube-apiserver-{Node Name} 的 Pod，并自动创建一个使用了新配置参数的 Pod 作为替代。
- 如果您有多个 Kubernetes Master 节点，您需要在每一个 Master 节点上都修改该文件，并使各节点上的参数保持一致。

## 修改 Kubelet 的启动参数

kubelet 组件是通过 systemctl 来管理的，因此可以在`/etc/systemd/system`
或`/usr/lib/systemd/system`下查找相关配置文件

```shell
找到kubelet对应服务的配置文件目录
# cd /etc/systemd/system/kubelet.service.d/   

查看原文件内容
# cat kubelet.service
[Unit]
Description=kubelet: The Kubernetes Node Agent
Documentation=https://kubernetes.io/docs/

[Service]
User=root
ExecStart=/usr/bin/kubelet
Restart=always
StartLimitInterval=0
RestartSec=10

[Install]
WantedBy=multi-user.target

查看kubelet的相关启动参数
# ps -ef | grep kubelet
root     16199     1  2 Sep27 ?        00:29:19 /usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --cgroup-driver=cgroupfs --hostname-override=10.0.222.102 --network-plugin=cni --node-labels=platform.tkestack.io/machine-ip=10.0.222.102 --pod-infra-container-image=registry.tke.com/library/pause:3.1
......

添加一个新的参数 –config
# vim kubelet.service


执行如下命令使新增参数生效
# systemctl stop kubelet
# systemctl daemon-reload
# systemctl start kubelet

检查新增参数是否已经生效
# ps -ef | grep kubelet
/usr/bin/kubelet --config=/etc/kubelet.d/ --kubeconfig=/etc/kubernetes/kubelet.conf --require-kubeconfig=true --pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true --network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin --cluster-dns=10.12.0.10 --cluster-domain=cluster.local --v=4
```

