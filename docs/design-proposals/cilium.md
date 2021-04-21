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
root@VM-0-3-ubuntu:~/ipv6-test# kubectl create -f cluster-ipv6-ds.json 
cluster.platform.tkestack.io/cls-7q46p7mt created
root@VM-0-3-ubuntu:~/ipv6-test#
```
2. After cluster is ready, apply cilium network plug-in 
```
root@VM-0-3-ubuntu:~/ipv6-test# kubectl get cluster 
NAME           TYPE        VERSION   STATUS    AGE
cls-7q46p7mt   Baremetal   1.18.3    Running   3m51s
global         Baremetal   1.18.3    Running   60m
root@VM-0-3-ubuntu:~/ipv6-test#
```
```
root@VM-0-67-ubuntu:~/ipv6-test# kubectl get po -A
NAMESPACE     NAME                                  READY   STATUS    RESTARTS   AGE
kube-system   coredns-bbc9b5888-5lhtp               0/1     Pending   0          3m19s
kube-system   coredns-bbc9b5888-hvb5v               0/1     Pending   0          3m19s
kube-system   etcd-172.22.0.67                      1/1     Running   0          2m29s
kube-system   kube-apiserver-172.22.0.67            1/1     Running   0          2m34s
kube-system   kube-controller-manager-172.22.0.67   1/1     Running   0          3m17s
kube-system   kube-proxy-q4v98                      1/1     Running   0          3m19s
kube-system   kube-scheduler-172.22.0.67            1/1     Running   0          3m17s
root@VM-0-67-ubuntu:~/ipv6-test# 
```
```
root@VM-0-67-ubuntu:~/ipv6-test# kubectl create -f ../calico-v3.16/calicov6.yaml 
configmap/calico-config created
customresourcedefinition.apiextensions.k8s.io/bgpconfigurations.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/bgppeers.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/blockaffinities.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/clusterinformations.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/felixconfigurations.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/globalnetworkpolicies.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/globalnetworksets.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/hostendpoints.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/ipamblocks.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/ipamconfigs.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/ipamhandles.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/ippools.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/kubecontrollersconfigurations.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/networkpolicies.crd.projectcalico.org created
customresourcedefinition.apiextensions.k8s.io/networksets.crd.projectcalico.org created
clusterrole.rbac.authorization.k8s.io/calico-kube-controllers created
clusterrolebinding.rbac.authorization.k8s.io/calico-kube-controllers created
clusterrole.rbac.authorization.k8s.io/calico-node created
clusterrolebinding.rbac.authorization.k8s.io/calico-node created
daemonset.apps/calico-node created
serviceaccount/calico-node created
deployment.apps/calico-kube-controllers created
serviceaccount/calico-kube-controllers created
```
```
root@VM-0-67-ubuntu:~# kubectl get po -A
NAMESPACE     NAME                                       READY   STATUS    RESTARTS   AGE
kube-system   calico-kube-controllers-866f6f96b5-28zpw   1/1     Running   0          78s
kube-system   calico-node-tjgx9                          1/1     Running   0          78s
kube-system   coredns-bbc9b5888-5lhtp                    1/1     Running   0          19m
kube-system   coredns-bbc9b5888-hvb5v                    1/1     Running   0          19m
kube-system   etcd-172.22.0.67                           1/1     Running   0          18m
kube-system   kube-apiserver-172.22.0.67                 1/1     Running   0          18m
kube-system   kube-controller-manager-172.22.0.67        1/1     Running   0          19m
kube-system   kube-proxy-q4v98                           1/1     Running   0          19m
kube-system   kube-scheduler-172.22.0.67                 1/1     Running   0          19m
root@VM-0-67-ubuntu:~# 
```
