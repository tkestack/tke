# etcdbackup

## Introduction

`etcdbackup` is responsible for backup etcd data.

## Set values in values.yaml before install the chart
Create secret to get conf.certs.cacrt(etcd-ca.crt)、conf.certs.clientcrt(etcd.crt)、conf.certs.clientkey(etcd.key).
```console
kubectl create secret generic etcdbackup --from-file=/opt/tke-installer/data/etcd-ca.crt --from-file=/opt/tke-installer/data/etcd.crt --from-file=/opt/tke-installer/data/etcd.key
```
Set the conf.storageClsName if you use remote storage type.
```console
kubectl get sc
```
Set the conf.hostPath if you use hostPath type.For example create the path in all nodes:
```console
mkdir -p /data/etcdbackup
chmod 777 /data/etcdbackup
```
Set the conf.image if you deploy in idc.Pull image `ccr.ccs.tencentyun.com/tdccimages/etcd-operator:v0.0.1` then retag the image and push into idc registry.You should use `registry.tke.com/library/etcd-operator:v0.0.1`image.
```console
nerdctl pull ccr.ccs.tencentyun.com/tdccimages/etcd-operator:v0.0.1
nerdctl tag ccr.ccs.tencentyun.com/tdccimages/etcd-operator:v0.0.1 registry.tke.com/library/etcd-operator:v0.0.1
nerdctl --insecure-registry login -u <username> -p <password> registry.tke.com
nerdctl --insecure-registry push registry.tke.com/library/etcd-operator:v0.0.1
```

## Installing the Chart

```console
helm install etcdbackuper etcdbackuper
```
## Deploy customer resource
Create the `etcdbackup.yaml`.Set the `etcdEndpoints`and `backupPolicy`.`backupIntervalInSecond`indicate how long to backup the data.
`maxBackups` the max backup numbers, `timeoutInSecond`every time etcdbackup could take.
```
apiVersion: etcd.database.coreos.com/v1beta2
kind: EtcdBackup
metadata:
  annotations:
  generation: 1
  labels:
    clusterName: gz-vpc-etcd-03
    region: gz
    source: etcd-life-cycle-operator
  name: gz-vpc-etcd-03
  namespace: etcd-ops
spec:
  backupPolicy:
    backupIntervalInSecond: 36000
    maxBackups: 3
    timeoutInSecond: 600
  clientTLSSecret: etcd-v3-secret
  hostPath:
    path: /data/
  etcdEndpoints:
  - https://<etcdip>:2379
  - https://<etcdip>:2379
  - https://<etcdip>:2379
  insecureSkipVerify: false
  storageType: HostPath
```
Apply `etcdbackup.yaml`
```
kubectl apply -f etcdbackup.yaml
```
## Check the backup data
For storage class type
```
$ kubectl get pods -n etcd-ops
$ kubectl exec -it <etcdbackuper_pods> sh -n etcd-ops
 / $ mount | grep ceph-fuse
ceph-fuse on /data type fuse.ceph-fuse (rw,nosuid,nodev,relatime,user_id=0,group_id=0,allow_other)
/ $ ls /data/
etcdbackup_v514592_2022-06-02-11:11:15  etcdbackup_v514742_2022-06-02-11:11:35  etcdbackup_v514886_2022-06-02-11:11:55

```
For hostPath type
```
ls /data/etcdbackup/
```


