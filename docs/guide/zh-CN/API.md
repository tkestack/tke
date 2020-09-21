# API使用指引

## 1. 调用方式
### 1.1. 创建访问凭证

访问【平台管理】控制台，在左侧找到【组织资源】，选择【访问凭证】，新建一个访问凭证。

### 1.2. 访问API

TKEStack上各种资源的接口均以K8S原生API的形式提供，所有接口使用统一的前缀: `http://console.tke.com:8080/platform`，请求中需要将上一步申请的访问凭证以`"Authorization: Bearer ${访问凭证}"`的形式放入header。

以查询集群信息为例，使用的请求如下 

```shell
curl -H "Authorization: Bearer xxxxxxx" \
"http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/clusters"
```



### 1.3. 查看特定集群的 Namespace

查看集群所包含的 Namespace 需要传递 "X-TKE-ClusterName: cls-xxx" 的 header，cls-xxx 为特定集群 id

```shell
curl -H "Authorization: Bearer xxxxxxx" \
-H "X-TKE-ClusterName: cls-xxx" \
"http://console.tke.com:8080/platform/api/v1/namespaces"
```



## 2. 通过API创建应用

### 2.1. 非TApp应用（deployment，statefulset，daemonset）

```shell
curl 'http://console.tke.com:8080/platform/apis/apps/v1/namespaces/命名空间/工作负载类型/工作负载名称' 

-X PATCH

-H 'Content-Type:application/strategic-merge-patch+json' 

-H 'X-TKE-ClusterName:所属集群'

-H 'Authorization: Bearer 访问凭证'

-d '{"spec":{"template":{"spec":{"containers":[{"name":"容器名称","image":"容器镜像"}]}}}}'
```

工作负载类型： 选择需要更新的工作负载类型（deployment，statefulset， daemonset）
所属集群：填写所要更新容器所属集群。
命名空间：填写所要更新容器所属的命名空间。
工作负载名称：填写所要更新容器的工作负载名称。
容器名称：填写所要更新容器的名称。
访问凭证：填写访问该容器资源的访问凭证，可以在“tkestack-组织资源-访问凭证“中获取该信息（**访问凭证有过期时间，如过期需要重新创建**）。
容器镜像：填写所要更新的Docker镜像

### 2.2.  TApp

tapp是自研的应用类型，更新镜像需要两步，首先获取当前的容器spec，调整镜像名后在调用更新接口

#### 2.2.1.  获取tapp spec

```shell
curl 'http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/clusters/所属集群/tapps?namespace=命名空间&name=工作负载名称'

-X GET

-H 'Authorization: Bearer 访问凭证'
```


返回值示例：

`{"apiVersion":"apps.tkestack.io/v1","kind":"TApp","metadata":{"creationTimestamp":"2020-06-10T13:35:54Z","generation":8,"labels":{"k8s-app":"kevintest","qcloud-app":"kevintest"},"name":"kevintest","namespace":"default","resourceVersion":"13925571","selfLink":"/apis/apps.tkestack.io/v1/namespaces/default/tapps/kevintest","uid":"0269fb69-fa87-42f8-9c3a-e1f96cef40f1"},"spec":{"forceDeletePod":true,"replicas":1,"selector":{"matchLabels":{"k8s-app":"kevintest","qcloud-app":"kevintest"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"k8s-app":"kevintest","qcloud-app":"kevintest","tapp_template_hash_key":"9636164821252331163","tapp_uniq_hash_key":"9518255606018677371"}},"spec":{"containers":[{"image":"mirrors.tencent.com/elsanli/devops-demo:62","imagePullPolicy":"Always","livenessProbe":{"failureThreshold":10,"periodSeconds":10,"successThreshold":1,"tcpSocket":{"port":8888},"timeoutSeconds":2},"name":"test","readinessProbe":{"failureThreshold":10,"periodSeconds":30,"successThreshold":1,"tcpSocket":{"port":8888},"timeoutSeconds":2},"resources":{"limits":{"cpu":"100m","memory":"48Mi"},"requests":{"cpu":"100m","memory":"25Mi"}}}],"restartPolicy":"Always"}},"updateStrategy":{}},"status":{"appStatus":"Running","observedGeneration":7,"readyReplicas":0,"replicas":1,"scaleLabelSelector":"k8s-app=kevintest,qcloud-app=kevintest","statuses":{"0":"Pending"}}}`

### 2.3. 更新tapp镜像

从上一步返回值中获取想要更新的整个容器的spec，替换其中的image字段，这样做是为了避免将其他字段覆盖为空

```shell
curl ''http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/clusters/所属集群/tapps?namespace=命名空间&name=工作负载名称'

-X PATCH

-H 'Content-Type:application/merge-patch+json' 

-H 'X-TKE-ClusterName:所属集群'

-H 'Authorization: Bearer 访问凭证'

-d '{"spec":{"template":{"spec":{"containers":[{"name":"容器名称","image":"容器镜像","resources":{"limits":{"cpu":"100m","memory":"48Mi"},"requests":{"cpu":"100m","memory":"25Mi"}},"livenessProbe":{"tcpSocket":{"port":8888},"timeoutSeconds":2,"periodSeconds":10,"successThreshold":1,"failureThreshold":10},"readinessProbe":{"tcpSocket":{"port":8888},"timeoutSeconds":2,"periodSeconds":30,"successThreshold":1,"failureThreshold":10},"imagePullPolicy":"Always"}]}},"templates":null}}
```

所属集群：填写所要更新容器所属集群。
命名空间：填写所要更新容器所属的命名空间。
工作负载名称：填写所要更新容器的工作负载名称。
容器名称：填写所要更新容器的名称。
访问凭证：填写访问该容器资源的访问凭证，可以在“tkestack-组织资源-访问凭证“中获取该信息（访问凭证有过期时间，如过期需要重新创建）。
容器镜像：填写所要更新的Docker镜像

## 3. 通过API增删集群节点

只能对独立集群的节点进行增删操作，不可操作导入集群。

### 3.1. 增加节点

URL:  http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/machines

Method: POST

Headers:

1. Content-Type: application/json
2. Authorization: Bearer xxx

按照以下命令的格式，将中文部分替换成实际值，发送请求。请求成功后，会返回被创建的Machine对象。

```shell
curl -X POST \
"http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/machines" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer 你的访问凭证" \
-d '
{
    "kind": "Machine",
    "apiVersion": "platform.tkestack.io/v1",
    "metadata": {
        "generateName": "mc-"
    },
    "spec": {
        "finalizers": [
            "machine"
        ],
        "tenantID": "租户ID（联系平台管理员获取）",
        "clusterName": "集群ID，可通过页面查看（不是集群名称）",
        "type": "Baremetal",
        "ip": "节点IP",
        "port": 节点SSH端口（int）,
        "username": "root",
        "password": "节点root密码(需经base64编码)"
    }
}'
```


password base64编码：

`echo -n $PASSWORD | base64`
假设password原文为123456，则生成的base64编码为MTIzNDU2

> PS: 使用echo命令时一定加上-n参数

### 3.2. 查看节点

URL:  http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/machines/${machine.metadata.name}

Method: GET

Headers:

1. Authorization: Bearer xxx

假设平台中有name为mc-brd44nzd的Machine对象：

```shell
{
  "kind": "Machine",
  "apiVersion": "platform.tkestack.io/v1",
  "metadata": {
    "name": "mc-brd44nzd",
    "generateName": "mc-",
    "selfLink": "/apis/platform.tkestack.io/v1/machines/mc-brd44nzd",
    "uid": "9ef7c08f-c535-4e99-b11d-9f7d02be19f5",
    "resourceVersion": "343953553",
    "creationTimestamp": "2020-02-27T00:25:02Z"
  },
  "spec": {
    "finalizers": [
      "machine"
    ],
    "tenantID": "default",
    "clusterName": "xxxx",
    "type": "Baremetal",
    "ip": "xxxxxx",
    "port": 36000,
    "username": "root",
    "password": "xxxxxx"
  }
}
```


则查看该Machine部署进度的请求为：

```shell
curl "http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/machines/mc-brd44nzd" \
-H "Authorization: Bearer 你的访问凭证"
```

### 3.3. 删除节点

URL:  http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/machines/${machine.metadata.name}

Method: DELETE

Headers:

1. Authorization: Bearer xxx

假设平台中有name为mc-brd44nzd的Machine对象，则删除节点的请求为：

```shell
curl -X DELETE "http://console.tke.com:8080/platform/apis/platform.tkestack.io/v1/machines/mc-brd44nzd" \
-H "Authorization: Bearer 你的访问凭证"
```


## 4. 通过API获取业务信息

### 4.1. 查看自身所在业务

```shell
curl 'http://console.tke.com:8080/business/apis/business.tkestack.io/v1/portal' \
-X GET \
-H "Authorization: Bearer 访问凭证"
```

### 4.2. 查看特定业务包含的 Namespace 信息

```shell
curl 'http://console.tke.com:8080/business/apis/business.tkestack.io/v1/namespaces/prj-xxx/namespaces' \
-X GET \
-H "Authorization: Bearer 访问凭证"
prj-xxx 为业务 id
```

### 4.3. 查看特定业务信息

```shell
curl 'http://console.tke.com:8080/business/apis/business.tkestack.io/v1/projects/prj-xxx' \
-X GET \
-H 'Authorization: Bearer 访问凭证'
prj-xxx为业务id
```

