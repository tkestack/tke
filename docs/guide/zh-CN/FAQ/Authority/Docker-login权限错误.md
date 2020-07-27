### Docker login 权限错误

在Tkestack选用用了自建证书，需要用户在客户端手动导入，docker login 权限报错：certificate signed by unknown authority。

##### 方法一

在 Global 集群上执行 kubectl get cm certs -n tke -o yaml
将 ca.crt 内容保存到客户端节点的/etc/docker/certs.d/******/ca.crt ( 为镜像仓库地址)
重启docker即可

##### 方法二：

  在/etc/docker/daemon.json文件里添加insecure-registries，如下：
  {
        "insecure-registries": [
         "xxx","xxx"
        ]
  }
（*** 为镜像仓库地址）

  重启docker即可