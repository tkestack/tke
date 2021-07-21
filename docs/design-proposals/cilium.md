#  Enable Cilium on TKEStack clusters 

**Author**: listai([@hmtai](https://github.com/hmtai))

**Status** (20210421): Done

## Summary

Cilium is open source software for providing and transparently securing network connectivity and loadbalancing between application workloads such as application containers or processes.A new Linux kernel technology called eBPF is at the foundation of Cilium.Cilium is integrated into common orchestration frameworks such as Kubernetes.

## Background

As a Kubernetes CNI plug-in, Cilium was designed from the very beginning for large-scale and highly dynamic container environments, and brought API-level-aware network security management functions. By using a new technology based on Linux kernel features-BPF, Provided based on service/pod/container as an identifier, instead of a traditional IP address, to define and strengthen the security strategy of the network layer and application layer between the container and the Pod.

Cilium not only decouples security control and addressing to simplify the application of security policies in highly dynamic environments, but also provides traditional network layer 3 and 4 isolation functions, and isolation control based on the http layer to provide stronger security isolation.

eBPF is at the foundation of Cilium.Because BPF can dynamically insert programs that control the Linux system, it achieves powerful security visualization functions, and these changes can take effect without updating the application code or restarting the application service itself, because BPF runs in the system kernel.

The above features enable Cilium to be highly scalable, visualized, and secure in large-scale container environments.

## Motivation

TKEStack can support Cilium on network mode overlay and underlay.

User can manully set the following configurations when install Cilium on TKEStack from UI:

### Install Cilium network mode overlay

```
CNI      Cilium
网络模式 Overlay
```

### Install Cilium network mode underlay

```
CNI        Cilium
网络模式   Underlay
自治系统号(输入自制系统号) 
交换机IP(输入交换机IP)
```

## Scope

**In-Scope**:
- Support Cilium on TencentOS Server 3.1 for network mode overlay and underlay.
- Install Cilium on TKEStack baremetal clusters when user enabledCilium,otherwise default CNI is galaxy.
- Only support install cilium to store all required state using Kubernetes custom resource definitions (CRDs).
- Support tkestack HA topology include CLB and VIP.

**Out-Of-Scope**: 
- Imported clusters does not support install Cilium.
- Not support install Cilium with etcd.
- Not support IPV6.
- Cilium can not replace the kube-proxy.
- Not support cluster mesh.
- Not support "enable-l7-proxy".

### System Requirements

1. In order for the eBPF feature to be enabled properly, the following kernel configuration options must be enabled. This is typically the case with distribution kernels. When an option can be built as a module or statically linked, either choice is valid.
```
CONFIG_BPF=y
CONFIG_BPF_SYSCALL=y
CONFIG_NET_CLS_BPF=y
CONFIG_BPF_JIT=y
CONFIG_NET_CLS_ACT=y
CONFIG_NET_SCH_INGRESS=y
CONFIG_CRYPTO_SHA1=y
CONFIG_CRYPTO_USER_API_HASH=y
CONFIG_CGROUPS=y
CONFIG_CGROUP_BPF=y
```
2. The following ports should also be available on each node:

| Port Range/Protocol | Description |
| :--------:          |    :-----:  |
| 4240/tcp | cluster health checks (cilium-health) |
| 9876/tcp | cilium-agent health status API |
| 9890/tcp | cilium-agent gops server (listening on 127.0.0.1) |
| 9891/tcp | operator gops server (listening on 127.0.0.1) |

## Main proposal

1. Add EnsuerCilium switch that indicate CNI is cilium or Galaxy in cluster object and make Galaxy as default CNI.
2. Create EnsureCilium function to install Cilium.
3. Check the linux kernel version when the verison does not meet the requirement then installation will finished and throw a message for users.
4. Pass the Cilium configuration args `networkMode` by clusterSpec.NetworkArgs, then use go-template to overwrite the Cilium yaml.
5. Set the controller-manager flags: `--configure-cloud-routes=false`,`--allocate-node-cidrs=false` to disable the cloud provider set pod CIDR for nodes when deploy cilium on overlay network mode. Use ipam to allocate pod CIDR to nodes.
6. Add node labels: `infra.tce.io/as`, `infra.tce.io/switch-ip` on the first master for underlay network mode.Pass the label configuration args by clusterSpec.NetworkArgs.
7. Pass configuration args `cluster-cidrs`and `base-pod-number` for ipamd deployment.Pass configuration args `NonMasqueradeCIDRs` for ip-masq-agent.

## User cases

### Case 1 Deploy Cilium on overlay network mode

1. Deploy a TKEStack environment then create a new baremetal cluster with rebuid tke-platform-controller without galaxy. 
```
root@VM-0-20-ubuntu:~# kubectl create -f clusterCilium.yaml
cluster.platform.tkestack.io/cls-smt66nk6 created
```
2. After cluster is ready, apply cilium network plug-in 
```
root@VM-0-20-ubuntu:~# kubectl get cluster
NAME           TYPE        VERSION   STATUS    AGE
cls-smt66nk6   Baremetal   1.19.7    Running   3m46s
global         Baremetal   1.19.7    Running   14d
root@VM-0-20-ubuntu:~#
```
```
root@VM-0-16-ubuntu:~# kubectl apply -f quick.yaml
serviceaccount/cilium created
serviceaccount/cilium-operator created
configmap/cilium-config created
clusterrole.rbac.authorization.k8s.io/cilium created
clusterrole.rbac.authorization.k8s.io/cilium-operator created
clusterrolebinding.rbac.authorization.k8s.io/cilium created
clusterrolebinding.rbac.authorization.k8s.io/cilium-operator created
daemonset.apps/cilium created
deployment.apps/cilium-operator created
```
```
root@VM-0-16-ubuntu:~# kubectl get pods -n kube-system
NAME                                     READY   STATUS    RESTARTS   AGE
cilium-92mx2                             1/1     Running   0          3m11s
cilium-operator-59df9b5857-vmzf7         1/1     Running   0          3m11s
coredns-ccc77fb9d-jq7b8                  1/1     Running   0          3m11s
coredns-ccc77fb9d-khtn4                  1/1     Running   0          3m11s
etcd-10.0.0.16                           1/1     Running   0          83s
kube-apiserver-10.0.0.16                 1/1     Running   0          77s
kube-controller-manager-10.0.0.16        1/1     Running   0          109s
kube-proxy-cvhb2                         1/1     Running   0          3m11s
kube-scheduler-10.0.0.16                 1/1     Running   0          88s
metrics-server-v0.3.6-794ccd69c8-nr8g2   2/2     Running   0          54s
tke-eni-ipamd-98f6f5d55-bhf6s            1/1     Running   0          3m11s
```
3. Test Cilium can work. It will deploy a series of deployments which will use various connectivity paths to connect to each other. Connectivity paths include with and without service load-balancing and various network policy combinations. The pod name indicates the connectivity variant and the readiness and liveness gate indicates success or failure of the test. Make sure you have at least two avaliable nodes to use.
```
root@VM-0-16-ubuntu:~# kubectl apply -f connectivity-check.yaml
service/echo-a created
deployment.apps/echo-a created
service/echo-b created
service/echo-b-headless created
deployment.apps/echo-b created
deployment.apps/host-to-b-multi-node-clusterip created
deployment.apps/host-to-b-multi-node-headless created
deployment.apps/pod-to-a-allowed-cnp created
ciliumnetworkpolicy.cilium.io/pod-to-a-allowed-cnp created
deployment.apps/pod-to-a-l3-denied-cnp created
ciliumnetworkpolicy.cilium.io/pod-to-a-l3-denied-cnp created
deployment.apps/pod-to-a created
deployment.apps/pod-to-b-intra-node created
deployment.apps/pod-to-b-multi-node-clusterip created
deployment.apps/pod-to-b-multi-node-headless created
deployment.apps/pod-to-a-external-1111 created
deployment.apps/pod-to-external-fqdn-allow-baidu-cnp created
ciliumnetworkpolicy.cilium.io/pod-to-external-fqdn-allow-baidu-cnp created
```
Check the test results.
```
root@VM-0-16-ubuntu:~# kubectl get pods
NAME                                                    READY   STATUS    RESTARTS   AGE
echo-a-9c5d8bfcf-65vhk                                  1/1     Running   0          106s
echo-b-79c6c76fb4-qn9v6                                 1/1     Running   0          106s
host-to-b-multi-node-clusterip-78ffcc7449-nr6gq         1/1     Running   0          106s
host-to-b-multi-node-headless-6dcb4d494c-k9ctj          1/1     Running   1          106s
pod-to-a-allowed-cnp-9f5cf94c4-wmgz2                    1/1     Running   0          106s
pod-to-a-external-1111-76c557fc56-5rnrg                 1/1     Running   0          105s
pod-to-a-f747cbc86-xgmsz                                1/1     Running   0          105s
pod-to-a-l3-denied-cnp-6f6c68d6d4-58j8w                 1/1     Running   0          106s
pod-to-b-intra-node-fd66d747-qckbm                      1/1     Running   0          105s
pod-to-b-multi-node-clusterip-77cc47f747-bptlb          1/1     Running   0          105s
pod-to-b-multi-node-headless-64b6d4fc95-568rc           1/1     Running   1          105s
pod-to-external-fqdn-allow-baidu-cnp-67568c4d96-xqx8j   1/1     Running   0          104s
```

### Case 2 deploy Cilium on underlay network mode
1. Create a cluster with NetworkArgs{ans:64512;switch-ip:192.168.60.2}

```
root@VM-0-20-ubuntu:~#kubectl create -f clusterCilium0.yaml
```
```
root@VM-0-20-ubuntu:~# kubectl get cluster
NAME           TYPE        VERSION   STATUS    AGE
cls-jz7clcth   Baremetal   1.19.7    Running   51m
cls-kz8xchfx   Baremetal   1.19.7    Running   22m
global         Baremetal   1.19.7    Running   15d
```
2. Check the Cilium pods and labels:
```
[root@TENCENT64 ~]# kubectl get pods -n kube-system
NAME                                     READY   STATUS    RESTARTS   AGE
cilium-ksg69                             1/1     Running   1          9m52s
cilium-operator-59df9b5857-hxbxn         1/1     Running   1          9m52s
cilium-router-vvggb                      1/1     Running   0          63s
coredns-ccc77fb9d-444pp                  1/1     Running   0          9m52s
coredns-ccc77fb9d-q8bcm                  1/1     Running   0          9m52s
etcd-192.168.50.6                        1/1     Running   3          7m21s
kube-apiserver-192.168.50.6              1/1     Running   4          6m18s
kube-controller-manager-192.168.50.6     1/1     Running   3          7m29s
kube-proxy-lx8pr                         1/1     Running   0          9m52s
kube-scheduler-192.168.50.6              1/1     Running   3          9m23s
metrics-server-v0.3.6-794ccd69c8-2nsqg   2/2     Running   0          4m
tke-eni-ipamd-5c6bb5f777-75j4g           1/1     Running   0          5m46s
```
Check the node labels:
```
[root@TENCENT64 ~]# kubectl get nodes 192.168.50.6 --show-labels
NAME           STATUS   ROLES    AGE     VERSION         LABELS
192.168.50.6   Ready    master   8m40s   v1.20.4-tke.1   beta.kubernetes.io/arch=amd64,beta.kubernetes.io/os=linux,infra.tce.io/as=64512,infra.tce.io/switch-ip=192.168.60.2,kubernetes.io/arch=amd64,kubernetes.io/hostname=192.168.50.6,kubernetes.io/os=linux,node-role.kubernetes.io/master=,platform.tkestack.io/machine-ip=192.168.50.6
```
Through the results the NetworkArgs has passed into kubernetes node labels.

### Case 3
1. Install Cilium on linux kernel version 4.10:
When ensure cilium the installation breakdown because the Linux kernel doesn't meet the requirement minimal  version 4.11.
```
EnsureCilium	失败	2021-04-27 18:07:37	FailedInit
```
The error log:
```
install cilium error: 10.0.0.46: [preflight] Some fatal errors occurred: [ERROR KernelCheck-4-10]: kernel version(4.10.0-118-generic) must not lower than 4.11
```

### Case 4

1. Create cluster with build-in HA(keepalived + vip + TencentOS 3.1)

> This case depends on ARP function, please make sure it is able in your router.

Create cluster in global:

```sh
wget https://github.com/tkestack/tke/blob/master/docs/yamls/cilium/cls-vip.json
## edit json to fulfill your vip, machine ip and ssh info
kubectl -f cls-vip.json
```

2. After HA cluster is running, create cillium in HA cluster:

```sh
kubectl apply -f https://github.com/tkestack/tke/blob/master/docs/yamls/cilium/cilium-ha.yaml
```

Check pod status through `kubectl get pod -A`, all pods will be running status.

3. Confirm HA is working. Shutdown the node which is binding vip, will find vip is binded on another node, everthing wokrs fine through UI and kubectl.


### Case 5

1. Create cluster with third-party HA(CLB + TencentOS 3.1)

Create cluster in global:

```sh
wget https://github.com/tkestack/tke/blob/master/docs/yamls/cilium/cls-clb.json
## edit json to fulfill your clb ip, machine ip and ssh info
kubectl -f cls-clb.json
```

2. After HA cluster is running, create cillium in HA cluster:

```sh
kubectl apply -f https://github.com/tkestack/tke/blob/master/docs/yamls/cilium/cilium-ha.yaml
```

Check pod status through `kubectl get pod -A`, all pods will be running status.

3. Confirm HA is working. Shutdown any one node, everthing works fine through UI and kubectl.
