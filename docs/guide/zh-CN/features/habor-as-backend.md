# 使用Harbor仓库作为存储镜像和Helm Chart后端



## 配置第三方Harbor



### 安装Harbor

在K8S安装Harbor可以参考官方Helm Chart安装方法

> https://github.com/goharbor/harbor-helm

需要注意:

- 域名修改成自己真正部署环境所需域名
- xsrfKey需要为32位的字符串，否则可能导致登录的时候出现秘钥出错的的情况
- tkestack对接的Harbor必须是https访问方式



### 配置tkestack

假设Harbor服务在tkestack所在集群或者在外面集群运行正常，那么需要配置下面选项

1，找到Harbor的ca证书，在tkestack的configmap/certs 添加 harbor-ca.crt文件

2，修改tke-registry-api和tke-registry-controller的configmap/tke-registry-api， 在文件 tke-registry-api.yaml 下添加下面2个新字段

```
harborCAFile: "/app/certs/harbor-ca.crt"
harborEnabled: true
```

3，修改文件 tke-registry-api和tke-registry-controller的configmap/tke-registry-api.yaml 下原来的字段

```
domainSuffix: "core.harbor.tcnp.qqa.com"     // harbor的域名，注意需要在tke-registry-api的pod内可以解析
security:
  adminPassword: "Harbor12345"    // harbor的密码
  adminUsername: "admin"          // harbor的用户名
```



### 测试使用Harbor作为镜像后端

1，在tkestack界面创建命名空间

创建完成后，在Harbor会同时创建一个项目，命名方式为 "`租户id`-image-`命名空间名字`"。

例如default租户创建一个命名空间 test，则harbor会同时创建一个 default-image-test 的项目

2，推送镜像

推送镜像和目前tke的镜像仓库上传指引基本上一致，不过目前docker login 所用的用户名和密码是Harbor的用户名密码，其他不变。镜像的tag也还是创建的时候命名空间的名字，下面为例子:

- docker login default.registry.tke.com
- docker tag nginx:latest default.registry.tke.com/test/nginx.latest
- docker push default.registry.tke.com/test/nginx.latest

3，推送镜像后，在 tke console可以看到命名空间下的镜像列表会更新

- 如果镜像不存在的情况下，会新增镜像

- 如果镜像存在的情况下，会新增tag

4，在tke console 删除镜像的同时会把harbor对应的镜像删除

5，在tke console 删除命名空间的时候，会把harbor对应的镜像和项目都删除



### 测试使用Harbor作为Chart后端

1，新建镜像仓库

创建完成后，在Harbor会同时创建一个项目，命名方式为 "`租户id`-chart-`命名空间名字`"。

2，设置helm repo

和现在的指引基本上一致，但路径会需要从chart改为chartrepo，例如:

> helm repo add --ca-file tke-ca.crt tke-admin http://registry.tke.com/chartrepo/admin --username admin --password xxxxxxxx

此处用到的ca文件为tke域名的ca文件，但登录的用户名和密码需要为Harbor的用户名密码

3，安装helm-push 插件并推送helm chart

推送chart成功后，chart会被保存到Harbor，同时console中helm模板会新增一个chart模板

4，在console删除模板版本或者通过删除cr的chart，Harbor中的chart会对应被删除

5，在console删除模板仓库或者通过删除cr的chartGroup，Harbor中的项目会对应被删除，存储的Chart也会被删除





## 限制

目前tke和harbor的数据交互都通过代理监听API的方式处理，因此如果用户单独往harbor创建项目或者上传镜像等操作，tke不会同步该数据。

其次，如果接入tke的不是一个全新空的harbor，那么已有的harbor数据不会同步到tke上，目前tke仅支持全新的harbor接入，并把tke作为harbor数据的单一来源。

