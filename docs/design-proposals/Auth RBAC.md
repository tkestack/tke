# Auth RBAC


**Author**:
jason ([@wangao1236](https://github.com/wangao1236))
mingyu  ([@metang326](https://github.com/metang326))


## Abstract

目前tke不支持k8s原生的rbac鉴权模式，无法通过创建ClusterRoleBinding、ClusterRole设置clusters、clustercredentials等资源的访问权限。


## Motivation

- 支持k8s原生的rbac鉴权

## Main proposal

增加一种鉴权方式，如果rbac鉴权通过则允许访问相应资源。

## Solution
1. 创建ClusterRoleBindings，将tke命名空间下的ServiceAccount（name：default）与ClusterRole（name：cluster-admin）进行绑定，使default能获得访问所有资源的权限。
2. 增加k8sinformers，对Roles、RoleBindings、ClusterRoles、ClusterRoleBindings进行watch。
3. 增加rbacAuthorizer，根据rbac进行鉴权，如有相应权限则Allowed为true
4. 由于ServiceAccount本身没有租户信息，会因为tenantID为空导致访问某些资源时会提示invalid tenantID，未进行鉴权就返回了。目前解决方法是对ServiceAccount命名时用后缀标记租户，例如-tenant-default，在代码中会对形如"system:serviceaccount:xx:xxx-tenant-xxx"的username进行解析。

## Test
### step 0
在tmy这个命名空间下创建一个名为rbac-tenant-default的sa用于测试，其基本信息如下：
```
# kubectl get sa -ntmy rbac -oyaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: rbac-tenant-default
  namespace: tmy
secrets:
- name: rbac-token-5ljnn
```
获取其对应的secret后生成一个kubeconfig文件，用于后续的测试。
```
kubectl get secret -ntmy rbac-token-5ljnn -oyaml
echo -n "xxx" | base64 -d
cp ~/.kube/config test.config
vim test.config
```
test.config内容如下，其中token是上面执行echo -n "xxx" | base64 -d的结果。
```
apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://127.0.0.1:6443
  name: default-cluster
contexts:
- context:
    cluster: default-cluster
    user: default-user
  name: default-context
current-context: default-context
kind: Config
preferences: {}
users:
- name: default-user
  user:
    token: xxxx
```
创建clusterrole和clusterrolebinding，允许rbac对clusters进行get、list。
```
# kubectl get clusterrole cls-role -o yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cls-role
rules:
- apiGroups:
  - '*'
  resources:
  - clusters
  verbs:
  - get
  - list

# kubectl get clusterrolebinding cls-bind -o yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cls-bind
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cls-role
subjects:
- kind: ServiceAccount
  name: rbac-tenant-default
  namespace: tmy
```

# step 1
使用rbac的kubeconfig文件执行kubectl命令，clusters资源被授权因此可以访问；clustercredentials资源未被授权，因此不能访问。
```
# kubectl get clusters --kubeconfig=test.config
NAME           CREATED AT
cls-mmxx4mkr   2021-04-22T03:29:24Z
global         2021-02-23T06:49:58Z

# kubectl get clustercredentials  --kubeconfig=test.config
Error from server (Forbidden): clustercredentials.platform.tkestack.io is forbidden: User "system:serviceaccount:tmy:rbac-tenant-default" cannot list resource "clustercredentials" in API group "platform.tkestack.io" at the cluster scope: permission for list on clustercredentials not verify
```

# step 2
通过kubectl edit clusterrole cls-role增加对clustercredentials的访问权限，修改后如下：
```
# kubectl get clusterrole cls-role -o yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cls-role
rules:
- apiGroups:
  - '*'
  resources:
  - clusters
  - clustercredentials
  verbs:
  - get
  - list
```
增加权限后，可访问clustercredentials资源。
```
# kubectl get clustercredentials  --kubeconfig=test.config
NAME          CREATED AT
cc-global     2021-02-23T06:49:58Z
cc-xsvfnrng   2021-04-22T03:29:24Z
```
