# 私有化部署FloatingIP配置

## 集群管理员配置FloatingIP网段
 
创建下面格式的configmap，修改其中的IP配置（请确保这些POD IP是配置在物理路由器上的，可以是单独给容器配置的网段，也可以是宿主机网段的部分IP）

```
kind: ConfigMap
apiVersion: v1
metadata:
 name: floatingip-config
 namespace: kube-system
data:
floatingips: '[{"routableSubnet":"10.0.0.0/16","ips":["10.0.70.2~10.0.70.241"],"subnet":"10.0.70.0/24","gateway":"10.0.70.1","vlan":2}]'
```
 
- routableSubnet: kubelet节点网段
- ips: POD可用IP
- subnet: POD可用IP网段
- vlan（可选配置）: 如果POD IP网段与宿主机网段不同vlan，需要配置这个字段
- gateway: POD网段的网关IP（不可以配置为.0的IP，.0的IP无法配置默认路由）

# FloatingIP使用

## 提交APP配置

在podSpec的annotation中增加`k8s.v1.cni.cncf.io/networks`和`k8s.v1.cni.galaxy.io/release-policy` annotation，container的limits和requests resources中增加`tke.cloud.tencent.com/eni-ip: "1"`，配置示例如下：

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
      annotations:
        k8s.v1.cni.cncf.io/networks: galaxy-k8s-vlan
        k8s.v1.cni.galaxy.io/release-policy: never
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        resources:
          limits:
            tke.cloud.tencent.com/eni-ip: "1"
          requests:
            tke.cloud.tencent.com/eni-ip: "1"
```

k8s.v1.cni.cncf.io/networks表示使用哪种CNI插件给POD配置网络，floatingip网络支持两种插件galaxy-k8s-vlan和galaxy-k8s-sriov，galaxy默认配置只开启了galaxy-k8s-vlan

k8s.v1.cni.galaxy.io/release-policy表示IP回收策略，value有三种选择：

- 不配置这个annotation或者配置位空字符串，表示“随时回收”策略，POD重启或跨机迁移，缩容或删除APP时IP均会被回收并重新分配；
- 配置为immutable标示“缩容或删除APP时回收”策略，保证POD重启或者跨机迁移时IP不变，缩容或删除APP时回收IP；
- 配置为never表示“永不回收”策略，IP永不被回收，删除APP后创建同名APP会分配到相同IP。可以通过调用galaxy-ipam的API进行回收。

另外FloatingIP还支持多个Deployment共享同一个IP池的功能，以实现蓝绿发布时IP不变。如果需要使用这个功能，可以给PodSpec再配置一个annotation `tke.cloud.tencent.com/eni-ip-pool: ${poolName}`，poolName可以随便取名。注意配置`tke.cloud.tencent.com/eni-ip-pool`的deployment的IP回收策略均为never，无论是否配置了`k8s.v1.cni.galaxy.io/release-policy`

## API文档

对于永不释放的IP，可通过调用API查询和回收IP，比如：

查询nginx-deployment tapp的ip（galaxy-ipam service的API端口是32761）：

```
curl "http://192.168.99.100:32761/v1/ip?namespace=default&appName=nginx-deployment&appType=tapp"

{
  "last": true,
  "totalElements": 3,
  "totalPages": 1,
  "first": true,
  "numberOfElements": 3,
  "size": 10,
  "number": 0,
  "content": [
   {
    "ip": "10.0.70.142",
    "namespace": "default",
    "appName": "nginx-deployment",
    "podName": "nginx-deployment-2",
    "policy": 2,
    "appType": "tapp",
    "updateTime": "2019-10-10T08:37:13.498446133Z",
    "status": "Deleted",
    "releasable": true
   },
   {
    "ip": "10.0.70.169",
    "namespace": "default",
    "appName": "nginx-deployment",
    "podName": "nginx-deployment-1",
    "policy": 2,
    "appType": "tapp",
    "updateTime": "2019-10-10T08:37:13.447486802Z",
    "status": "Deleted",
    "releasable": true
   },
   {
    "ip": "10.0.70.33",
    "namespace": "default",
    "appName": "nginx-deployment",
    "podName": "nginx-deployment-0",
    "policy": 2,
    "appType": "tapp",
    "updateTime": "2019-10-10T08:37:13.468236091Z",
    "status": "Deleted",
    "releasable": true
   }
  ]
 }
```

回收nginx-deployment-1 pod的ip

```
curl -v -H 'Content-type: application/json' -d '{"ips":[{"ip":"10.0.70.169", "namespace":"default", "appName":"nginx-deployment", "appType":"tapp", "podName": "nginx-deployment-1"}]}' "http://192.168.99.100:32761/v1/ip"
{
  "code": 200,
  "message": ""
}
```

API swagger文档：

```json
{
  "swaggerVersion": "1.2",
  "apiVersion": "",
  "basePath": "http://192.168.99.100:32761",
  "resourcePath": "/v1",
  "info": {
   "title": "",
   "description": ""
  },
  "apis": [
   {
    "path": "/v1/ip",
    "description": "",
    "operations": [
     {
      "type": "api.ListIPResp",
      "method": "GET",
      "summary": "List ips by keyword or params",
      "nickname": "ListIPs",
      "parameters": [
       {
        "type": "string",
        "paramType": "query",
        "name": "keyword",
        "description": "keyword",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "string",
        "paramType": "query",
        "name": "poolName",
        "description": "pool name",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "string",
        "paramType": "query",
        "name": "appName",
        "description": "app name",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "string",
        "paramType": "query",
        "name": "podName",
        "description": "pod name",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "string",
        "paramType": "query",
        "name": "namespace",
        "description": "namespace",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "boolean",
        "paramType": "query",
        "name": "isDeployment",
        "description": "listing deployments or statefulsets. Deprecated, please set appType",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "string",
        "paramType": "query",
        "name": "appType",
        "description": "app type, deployment, statefulset or tapp",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "integer",
        "paramType": "query",
        "name": "page",
        "description": "page number, valid range [0,99999]",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "integer",
        "defaultValue": "10",
        "paramType": "query",
        "name": "size",
        "description": "page size, valid range (0,9999]",
        "required": false,
        "allowMultiple": false
       },
       {
        "type": "string",
        "defaultValue": "ip asc",
        "paramType": "query",
        "name": "sort",
        "description": "sort by which field, supports ip/namespace/podname/policy asc/desc",
        "required": false,
        "allowMultiple": false
       }
      ],
      "responseMessages": [
       {
        "code": 200,
        "message": "request succeed",
        "responseModel": "api.ListIPResp"
       }
      ],
      "produces": [
       "application/json"
      ],
      "consumes": [
       "application/json"
      ]
     },
     {
      "type": "api.ReleaseIPResp",
      "method": "POST",
      "summary": "Release ips",
      "nickname": "ReleaseIPs",
      "parameters": [
       {
        "type": "api.ReleaseIPReq",
        "paramType": "body",
        "name": "body",
        "description": "",
        "required": true,
        "allowMultiple": false
       }
      ],
      "responseMessages": [
       {
        "code": 200,
        "message": "request succeed",
        "responseModel": "api.ReleaseIPResp"
       },
       {
        "code": 202,
        "message": "Unreleased ips have been released or allocated to other pods, or are not within valid range",
        "responseModel": "api.ReleaseIPResp"
       },
       {
        "code": 400,
        "message": "10.0.0.2 is not releasable"
       },
       {
        "code": 500,
        "message": "internal server error"
       }
      ],
      "produces": [
       "application/json"
      ],
      "consumes": [
       "application/json"
      ]
     }
    ]
   },
   {
    "path": "/v1/pool/{name}",
    "description": "",
    "operations": [
     {
      "type": "api.GetPoolResp",
      "method": "GET",
      "summary": "Get pool by name",
      "nickname": "Get",
      "parameters": [
       {
        "type": "string",
        "paramType": "path",
        "name": "name",
        "description": "pool name",
        "required": true,
        "allowMultiple": false
       }
      ],
      "responseMessages": [
       {
        "code": 200,
        "message": "request succeed",
        "responseModel": "api.GetPoolResp"
       },
       {
        "code": 400,
        "message": "pool name is empty"
       },
       {
        "code": 404,
        "message": "pool not found"
       },
       {
        "code": 500,
        "message": "internal server error"
       }
      ],
      "produces": [
       "application/json"
      ],
      "consumes": [
       "application/json"
      ]
     },
     {
      "type": "httputil.Resp",
      "method": "DELETE",
      "summary": "Delete pool by name",
      "nickname": "Delete",
      "parameters": [
       {
        "type": "string",
        "paramType": "path",
        "name": "name",
        "description": "pool name",
        "required": true,
        "allowMultiple": false
       }
      ],
      "responseMessages": [
       {
        "code": 200,
        "message": "request succeed",
        "responseModel": "httputil.Resp"
       },
       {
        "code": 400,
        "message": "pool name is empty"
       },
       {
        "code": 404,
        "message": "pool not found"
       },
       {
        "code": 500,
        "message": "internal server error"
       }
      ],
      "produces": [
       "application/json"
      ],
      "consumes": [
       "application/json"
      ]
     }
    ]
   },
   {
    "path": "/v1/pool",
    "description": "",
    "operations": [
     {
      "type": "httputil.Resp",
      "method": "POST",
      "summary": "Create or update pool",
      "nickname": "CreateOrUpdate",
      "parameters": [
       {
        "type": "api.Pool",
        "paramType": "body",
        "name": "body",
        "description": "",
        "required": true,
        "allowMultiple": false
       }
      ],
      "responseMessages": [
       {
        "code": 200,
        "message": "request succeed",
        "responseModel": "api.UpdatePoolResp"
       },
       {
        "code": 202,
        "message": "No enough IPs",
        "responseModel": "api.UpdatePoolResp"
       },
       {
        "code": 400,
        "message": "pool name is empty"
       },
       {
        "code": 500,
        "message": "internal server error"
       }
      ],
      "produces": [
       "application/json"
      ],
      "consumes": [
       "application/json"
      ]
     }
    ]
   }
  ],
  "models": {
   "api.ListIPResp": {
    "id": "api.ListIPResp",
    "required": [
     "content",
     "last",
     "totalElements",
     "totalPages",
     "first",
     "numberOfElements",
     "size",
     "number"
    ],
    "properties": {
     "content": {
      "type": "array",
      "items": {
       "$ref": "api.FloatingIP"
      }
     },
     "last": {
      "type": "boolean",
      "description": "if this is the last page"
     },
     "totalElements": {
      "type": "integer",
      "format": "int32",
      "description": "total number of elements"
     },
     "totalPages": {
      "type": "integer",
      "format": "int32",
      "description": "total number of pages"
     },
     "first": {
      "type": "boolean",
      "description": "if this is the first page"
     },
     "numberOfElements": {
      "type": "integer",
      "format": "int32",
      "description": "number of elements in this page"
     },
     "size": {
      "type": "integer",
      "format": "int32",
      "description": "page size"
     },
     "number": {
      "type": "integer",
      "format": "int32",
      "description": "page index starting from 0"
     }
    }
   },
   "page.Page.content": {
    "id": "page.Page.content",
    "properties": {}
   },
   "api.FloatingIP": {
    "id": "api.FloatingIP",
    "required": [
     "ip",
     "policy"
    ],
    "properties": {
     "ip": {
      "type": "string",
      "description": "ip"
     },
     "namespace": {
      "type": "string",
      "description": "namespace"
     },
     "appName": {
      "type": "string",
      "description": "deployment or statefulset name"
     },
     "podName": {
      "type": "string",
      "description": "pod name"
     },
     "poolName": {
      "type": "string"
     },
     "policy": {
      "type": "integer",
      "description": "ip release policy"
     },
     "isDeployment": {
      "type": "boolean",
      "description": "deployment or statefulset, deprecated please set appType"
     },
     "appType": {
      "type": "string",
      "description": "deployment, statefulset or tapp"
     },
     "updateTime": {
      "type": "string",
      "format": "date-time",
      "description": "last allocate or release time of this ip"
     },
     "status": {
      "type": "string",
      "description": "pod status if exists"
     },
     "releasable": {
      "type": "boolean",
      "description": "if the ip is releasable. An ip is releasable if it isn't belong to any pod"
     }
    }
   },
   "api.ReleaseIPResp": {
    "id": "api.ReleaseIPResp",
    "required": [
     "code",
     "message"
    ],
    "properties": {
     "code": {
      "type": "integer",
      "format": "int32"
     },
     "message": {
      "type": "string"
     },
     "content": {
      "$ref": "httputil.Resp.content"
     },
     "unreleased": {
      "type": "array",
      "items": {
       "type": "string"
      },
      "description": "unreleased ips, have been released or allocated to other pods, or are not within valid range"
     }
    }
   },
   "httputil.Resp.content": {
    "id": "httputil.Resp.content",
    "properties": {}
   },
   "api.ReleaseIPReq": {
    "id": "api.ReleaseIPReq",
    "required": [
     "ips"
    ],
    "properties": {
     "ips": {
      "type": "array",
      "items": {
       "$ref": "api.FloatingIP"
      }
     }
    }
   },
   "api.GetPoolResp": {
    "id": "api.GetPoolResp",
    "required": [
     "code",
     "message",
     "pool"
    ],
    "properties": {
     "code": {
      "type": "integer",
      "format": "int32"
     },
     "message": {
      "type": "string"
     },
     "content": {
      "$ref": "httputil.Resp.content"
     },
     "pool": {
      "$ref": "api.Pool"
     }
    }
   },
   "api.Pool": {
    "id": "api.Pool",
    "required": [
     "name",
     "size",
     "preAllocateIP"
    ],
    "properties": {
     "name": {
      "type": "string",
      "description": "pool name"
     },
     "size": {
      "type": "integer",
      "format": "int32",
      "description": "pool size"
     },
     "preAllocateIP": {
      "type": "boolean",
      "description": "Set to true to allocate IPs when creating or updating pool"
     }
    }
   },
   "httputil.Resp": {
    "id": "httputil.Resp",
    "required": [
     "code",
     "message"
    ],
    "properties": {
     "code": {
      "type": "integer",
      "format": "int32"
     },
     "message": {
      "type": "string"
     },
     "content": {
      "$ref": "httputil.Resp.content"
     }
    }
   },
   "api.UpdatePoolResp": {
    "id": "api.UpdatePoolResp",
    "required": [
     "code",
     "message",
     "realPoolSize"
    ],
    "properties": {
     "code": {
      "type": "integer",
      "format": "int32"
     },
     "message": {
      "type": "string"
     },
     "content": {
      "$ref": "httputil.Resp.content"
     },
     "realPoolSize": {
      "type": "integer",
      "format": "int32",
      "description": "real num of IPs of this pool after creating or updating"
     }
    }
   }
  }
 }
```