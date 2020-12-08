# 自定义k8s版本升级

用户可以通过向TKEStack平台提供自定义版本的k8s，以允许集群升级到非内置的版本。本文将以v1.16.15版本的k8s作为例子演示用户如何将集群升级到自定义版本。本文中只以amd64环境作为示例，如果用户希望自己的物料镜像可以支持`multi-CPU architecture`，请在制作镜像和推送镜像阶段参考[Leverage multi-CPU architecture support](https://docs.docker.com/docker-for-mac/multi-arch/)和[构建多CPU架构支持的Docker镜像](https://blog.csdn.net/dev_csdn/article/details/79138424)。

## 制作provider-res镜像

provider镜像用于存储kubeadm、kubelet和kubectl的二进制文件。

执行下面命令为环境设置好版本号，并从官方下载好二进制文件并压缩，若遇到网络问题请通过其他途径下载对应二进制文件：

```sh
export RELEASE=v1.16.15 && \
curl -L --remote-name-all https://storage.googleapis.com/kubernetes-release/release/$RELEASE/bin/linux/amd64/{kubeadm,kubelet,kubectl} && \
chmod +x kubeadm kubectl kubelet && \
mkdir -p kubernetes/node/bin/ && \
cp kubelet kubectl kubernetes/node/bin/ && \
tar -czvf kubeadm-linux-amd64-$RELEASE.tar.gz kubeadm && \
tar -czvf kubernetes-node-linux-amd64-$RELEASE.tar.gz kubernete
```

执行下面命令生成dockerfiel:

```sh
cat << EOF >Dockerfile
FROM tkestack/provider-res:v1.18.3-2

WORKDIR /data

COPY kubernetes-*.tar.gz   res/linux-amd64/
COPY kubeadm-*.tar.gz      res/linux-amd64/

ENTRYPOINT ["sh"]
EOF
```

制作provider-res镜像：

```sh
docker build -t registry.tke.com/library/provider-res:myversion .
```

此处使用了默认的registry.tke.com作为registry的domian，如未使用默认的domain请修改为自定义的domain，下文中如遇到registry.tke.com也做相同处理。

## 为平台准备必要镜像

从官方下载k8s组件镜像，如遇到网络问题请通过其他途径下载：

```sh
docker pull k8s.gcr.io/kube-scheduler:$RELEASE && \
docker pull k8s.gcr.io/kube-controller-manager:$RELEASE && \
docker pull k8s.gcr.io/kube-apiserver:$RELEASE && \
docker pull k8s.gcr.io/kube-proxy:$RELEAS
```

重新为镜像为镜像打标签：

```sh
docker tag k8s.gcr.io/kube-proxy:$RELEASE registry.tke.com/library/kube-proxy:$RELEASE && \
docker tag k8s.gcr.io/kube-apiserver:$RELEASE registry.tke.com/library/kube-apiserver:$RELEASE && \
docker tag k8s.gcr.io/kube-controller-manager:$RELEASE registry.tke.com/library/kube-controller-manager:$RELEASE && \
docker tag k8s.gcr.io/kube-scheduler:$RELEASE registry.tke.com/library/kube-scheduler:$RELEASE
```

导出镜像：

```sh
docker save -o kube-proxy.tar registry.tke.com/library/kube-proxy:$RELEASE && \
docker save -o kube-apiserver.tar registry.tke.com/library/kube-apiserver:$RELEASE && \
docker save -o kube-controller-manager.tar registry.tke.com/library/kube-controller-manager:$RELEASE && \
docker save -o kube-scheduler.tar registry.tke.com/library/kube-scheduler:$RELEASE && \
docker save -o provider-res.tar registry.tke.com/library/provider-res:myversion
```

发送到global集群节点上：

```sh
scp kube*.tar provider-res.tar root@your_global_node:/root/
```

## 在global集群上导入物料

注意在此之后执行到命令都是发生在global集群节点上，为了方便首先在环境中设置版本号：

```sh
export RELEASE=v1.16.15
```

加载镜像：

```sh
docker load -i kube-apiserver.tar && \
docker load -i kube-controller-manager.tar && \
docker load -i kube-proxy.tar && \
docker load -i kube-scheduler.tar && \
docker load -i provider-res.tar
```

登陆registry：

```sh
docker login registry.tke.co
```

此处会提示输入用户名密码，如果默认使用了内置registry，用户名密码为admin的用户名密码，如果配置了第三方镜像仓库，请使用第三方镜像仓库的用户名密码。

登陆成功后推送镜像到registry：

```sh
docker push registry.tke.com/library/kube-apiserver:$RELEASE && \
docker push registry.tke.com/library/kube-controller-manager:$RELEASE && \
docker push registry.tke.com/library/kube-proxy:$RELEASE && \
docker push registry.tke.com/library/kube-scheduler:$RELEASE && \
docker push registry.tke.com/library/provider-res:myversion
```

为使得导入物料可以被平台使用，首先需要修改tke-platform-controller的deployment：

```sh
kubectl edit -n tke deployments tke-platform-controller
```

修改`spec.template.spec.initContainers[0].image`中的内容为刚刚制作的provider-res镜像`registry.tke.com/library/provider-res:myversion`。

其次需要修改cluster-info：

```sh
kubectl edit -n kube-public configmaps cluster-info
```

在data.k8sValidVersions内容中添加`"1.16.15"`。

## 升级集群到自定义版本

触发集群升级需要在global集群上修改cluster资源对象内容：

```sh
kubectl edit cluster cls-yourcluster
```

修改`spec.version`中的内容为`1.16.15`。

更详细的升级相关文档请参考：[K8S 版本升级说明](https://github.com/tkestack/tke/blob/master/docs/guide/zh-CN/best-practices/cluster-upgrade-guide.md)。

> 目前Web UI不允许补丁版本升级，会导致可以在UI升级选项中可以看到`1.16.15`版本，但是提示无法升级到该版本，后续版本中将会修复。当前请使用kubectl修改cluster资源对象内容升级自定义版本。
