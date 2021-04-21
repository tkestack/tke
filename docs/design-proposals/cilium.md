#  Enable Cilium on TKEStack clusters 


**Author**: listai([@hmtai](https://github.com/hmtai))

**Status** (20210421): Designing

## Summary

Cilium is open source software for providing and transparently securing network connectivity and loadbalancing between application workloads such as application containers or processes.A new Linux kernel technology called eBPF is at the foundation of Cilium.Cilium is integrated into common orchestration frameworks such as Kubernetes.

## Background

As a Kubernetes CNI plug-in, Cilium was designed from the very beginning for large-scale and highly dynamic container environments, and brought API-level-aware network security management functions. By using a new technology based on Linux kernel features-BPF, Provided based on service/pod/container as an identifier, instead of a traditional IP address, to define and strengthen the security strategy of the network layer and application layer between the container and the Pod.

Cilium not only decouples security control and addressing to simplify the application of security policies in highly dynamic environments, but also provides traditional network layer 3 and 4 isolation functions, and isolation control based on the http layer to provide stronger security isolation.

eBPF is at the foundation of Cilium.Because BPF can dynamically insert programs that control the Linux system, it achieves powerful security visualization functions, and these changes can take effect without updating the application code or restarting the application service itself, because BPF runs in the system kernel.

The above features enable Cilium to be highly scalable, visualized, and secure in large-scale container environments.

## Motivation

TKEStack can support CNI Cilium.

User can manully set field "tunnel" "enable-policy" which cilium exposed.

## Scope

 **In-Scope**: 
 1. (**P0**) Install Cilium on TKEStack baremetal clusters when user enabledCilium,otherwise default CNI is galaxy. 

**Out-Of-Scope**: 
 1. Other clusters does not support install Cilium.
 
## Main proposal

1. Prepare Cilium install yaml.Follow the community Cilium installation yaml to install Cilium on TKEStack.
2. Create EnsureCilium function to install Cilium.
3. Pass the Cilium configuration args by clusterSpec.NetworkArgs, then use go-template to overwrite the Cilium yaml.
4. Add EnsuerCilium switch that indicate CNI is cilium or Galaxy in cluster object and make Galaxy as default CNI. 

## User case

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
root@VM-0-46-ubuntu:~# kubectl apply -f quick.yaml
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
root@VM-0-46-ubuntu:~# kubectl get pods -n kube-system
NAME                                     READY   STATUS    RESTARTS   AGE
cilium-6lrvq                             1/1     Running   0          38s
cilium-operator-654456485c-wvkxx         1/1     Running   0          38s
coredns-745589f8f6-t5dj5                 1/1     Running   0          5m34s
coredns-745589f8f6-wth5g                 1/1     Running   0          5m34s
etcd-vm-0-46-ubuntu                      1/1     Running   0          6m11s
kube-apiserver-vm-0-46-ubuntu            1/1     Running   0          6m5s
kube-controller-manager-vm-0-46-ubuntu   1/1     Running   0          6m5s
kube-proxy-xclms                         1/1     Running   0          5m34s
kube-scheduler-vm-0-46-ubuntu            1/1     Running   0          6m5s
```
3. Test Cilium can work
```
kubectl apply -f connectivity-check.yaml
```
Check the test result.
```
root@VM-0-55-ubuntu:~# kubectl get pods
NAME                                                    READY   STATUS    RESTARTS   AGE
echo-a-9c5d8bfcf-nml2k                                  1/1     Running   0          68m
echo-b-79c6c76fb4-vmbbj                                 1/1     Running   0          68m
host-to-b-multi-node-clusterip-78ffcc7449-2bjr5         0/1     Pending   0          68m
host-to-b-multi-node-headless-6dcb4d494c-lflx4          0/1     Pending   0          68m
pod-to-a-allowed-cnp-9f5cf94c4-jcqqd                    1/1     Running   0          68m
pod-to-a-external-1111-76c557fc56-9htgk                 1/1     Running   0          68m
pod-to-a-f747cbc86-mgkl9                                1/1     Running   0          68m
pod-to-a-l3-denied-cnp-6f6c68d6d4-9h47d                 1/1     Running   0          68m
pod-to-b-intra-node-fd66d747-7ms9l                      1/1     Running   0          68m
pod-to-b-multi-node-clusterip-77cc47f747-zjm9g          0/1     Pending   0          68m
pod-to-b-multi-node-headless-64b6d4fc95-9gz5q           0/1     Pending   0          68m
pod-to-external-fqdn-allow-baidu-cnp-67568c4d96-2984d   1/1     Running   0          68m
```
