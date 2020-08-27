[TOC]

# 目录结构

```
.
|-- charts                                      ---	底层依赖，一般不需要修改
|   |-- bootstrap				---	基于 common, 渲染 deployment，service，ingress 模版
|   |   |-- Chart.lock
|   |   |-- charts
|   |   |   `-- common-0.0.5.tgz
|   |   |-- Chart.yaml
|   |   |-- templates
|   |   |   |-- deployment.yaml
|   |   |   |-- ingress.yaml
|   |   |   `-- service.yaml
|   |   `-- values.yaml
|   `-- common                                   ---	底层模版，函数，通常对全局适用通用配置会放到这里
|       |-- Chart.yaml
|       |-- README.md
|       |-- templates
|       |   |-- _chartref.tpl
|       |   |-- _config_center_client_containers.yaml
|       |   |-- _configmap.yaml
|       |   |-- _container.yaml
|       |   |-- _deployment.yaml
|       |   |-- _envvar.tpl
|       |   |-- _fullname.tpl
|       |   |-- _ingress.yaml
|       |   |-- _metadata_annotations.tpl
|       |   |-- _metadata_labels.tpl
|       |   |-- _metadata.yaml
|       |   |-- _name.tpl
|       |   |-- _persistentvolumeclaim.yaml
|       |   |-- _secret.yaml
|       |   |-- _service.yaml
|       |   |-- _skywalking_agent_containers.yaml
|       |   |-- _util.tpl
|       |   `-- _volume.tpl
|       `-- values.yaml
|-- docker                                      ---   helm 镜像制作
|   |-- Dockerfile
|   |-- Dockerstart
|   `-- run.sh                                  ---   业务 helm 打包
|-- helmfile.d                                  ---   helmfile 配置目录
|   |-- commons                                 ---   helmfile 通用配置
|   |   |-- helmdefaults.yaml                   
|   |   `-- repos.yaml
|   |-- config                                  ---   value.yaml 和 value 模板 .gotmpl 存放目录
|   |   |-- product1-demo
|   |   |   |-- environment-values.yaml
|   |   |   |-- product1-app1-demo.yaml         ---   value.yaml
|   |   |   |-- product1-app2-demo.yaml
|   |   |   `-- product1-helm-demo.yaml
|   |   |-- product1-dev
|   |   |   |-- environment-values.yaml
|   |   |   |-- product1-app1-dev.yaml
|   |   |   |-- product1-app2-dev.yaml
|   |   |   `-- product1-helm-dev.yaml
|   |   |-- product1-prod
|   |   |   |-- environment-values.yaml
|   |   |   |-- product1-app1-prod.yaml
|   |   |   |-- product1-app2-prod.yaml
|   |   |   `-- product1-helm-prod.yaml
|   |   |-- product1-stage
|   |   |   |-- environment-values.yaml
|   |   |   |-- product1-app1-stage.yaml
|   |   |   |-- product1-app2-stage.yaml
|   |   |   `-- product1-helm-stage.yaml
|   |   |-- product1-test
|   |   |   |-- environment-values.yaml
|   |   |   |-- product1-app1-test.yaml
|   |   |   |-- product1-app2-test.yaml
|   |   |   `-- product1-helm-test.yaml
|   |   |-- product2-demo
|   |   |   |-- environment-values.yaml
|   |   |   |-- product2-app1-demo.yaml.gotmpl           ---   go 模板， 取 environment-values.yaml 里的值渲染出真正的 value.yaml 
|   |   |   |-- product2-app2-demo.yaml.gotmpl
|   |   |   `-- product2-helm-demo.yaml.gotmpl
|   |   |-- product2-dev
|   |   |   |-- environment-values.yaml
|   |   |   |-- product2-app1-dev.yaml.gotmpl
|   |   |   |-- product2-app2-dev.yaml.gotmpl
|   |   |   `-- product2-helm-dev.yaml.gotmpl
|   |   |-- product2-prod
|   |   |   |-- environment-values.yaml
|   |   |   |-- product2-app1-prod.yaml.gotmpl
|   |   |   |-- product2-app2-prod.yaml.gotmpl
|   |   |   `-- product2-helm-prod.yaml.gotmpl
|   |   |-- product2-stage
|   |   |   |-- environment-values.yaml
|   |   |   |-- product2-app1-stage.yaml.gotmpl
|   |   |   |-- product2-app2-stage.yaml.gotmpl
|   |   |   `-- product2-helm-stage.yaml.gotmpl
|   |   `-- product2-test
|   |       |-- environment-values.yaml
|   |       |-- product2-app1-test.yaml.gotmpl
|   |       |-- product2-app2-test.yaml.gotmpl
|   |       `-- product2-helm-test.yaml.gotmpl
|   |-- helmfile.yaml
|   `-- releases                                       -- helmfile releases
|       |-- product1-demo.yaml                         -- 产品名-环境名 和 config 目录里一一对应
|       |-- product1-dev.yaml
|       |-- product1-prod.yaml
|       |-- product1-stage.yaml
|       |-- product1-test.yaml
|       |-- product2-demo.yaml
|       |-- product2-dev.yaml
|       |-- product2-prod.yaml
|       |-- product2-stage.yaml
|       `-- product2-test.yaml
|-- product1                                   ---   product1 部署配置
|   |-- product1-app1                          ---   product1-app1 chart
|   |   |-- Chart.lock
|   |   |-- charts
|   |   |   `-- bootstrap-0.1.0.tgz
|   |   |-- Chart.yaml
|   |   `-- values
|   |       `-- demo.yaml
|   |-- product1-app2                          ---   product1-app2 chart
|   |   |-- Chart.lock
|   |   |-- charts
|   |   |   `-- bootstrap-0.1.0.tgz
|   |   |-- Chart.yaml
|   |   `-- values
|   |       `-- demo.yaml
|   |-- product1-helm
|   |   |-- Chart.lock
|   |   |-- charts
|   |   |   `-- bootstrap-0.1.0.tgz
|   |   `-- Chart.yaml
|   `-- tools                                  ---   部署辅助工具集
|       |-- create_chart.sh
|       |-- helmfile.sh                        ---   helmfile 部署辅助工具
|       |-- helm.sh                            ---   helm 部署辅助工具
|       |-- product1-demo                      ---   product1 demo chart
|       |   |-- charts
|       |   |   `-- bootstrap-0.1.0.tgz
|       |   |-- Chart.yaml
|       |   `-- values
|       |       `-- demo.yaml
|       `-- update_dependencies.sh
|-- product2
|   |-- product2-app1
|   |   |-- Chart.lock
|   |   |-- charts
|   |   |   `-- bootstrap-0.1.0.tgz
|   |   |-- Chart.yaml
|   |   `-- values
|   |       `-- demo.yaml
|   |-- product2-app2
|   |   |-- Chart.lock
|   |   |-- charts
|   |   |   `-- bootstrap-0.1.0.tgz
|   |   |-- Chart.yaml
|   |   `-- values
|   |       `-- demo.yaml
|   |-- product2-helm
|   |   |-- Chart.lock
|   |   |-- charts
|   |   |   `-- bootstrap-0.1.0.tgz
|   |   `-- Chart.yaml
|   `-- tools
|       |-- create_chart.sh
|       |-- helmfile.sh
|       |-- helm.sh
|       |-- product2-demo
|       |   |-- charts
|       |   |   `-- bootstrap-0.1.0.tgz
|       |   |-- Chart.yaml
|       |   `-- values
|       |       `-- demo.yaml
|       `-- update_dependencies.sh
`-- README.md

```
# 新增服务和服务变更
新增服务和服务变更推荐使用 `./tools/helm.sh` 脚本操作
```sh
./tools/helm.sh -h
Usage:
  ./tools/helm.sh $action $namespace $group $chart_name $env [$tag]

Params:
  namespace    命名空间，如： default、haproxy-ingress
  action       操作命令，如： template、upgrade、install
  group        项目分组，如：product1、product2、product3
  env          环境名，会根据该值取对应的values文件（如values文件名为product1-app1-demo.yaml，则对应的env为ty）
  tag          镜像的tag，如：20190929101615（upgrade操作可省略，省略时取环境中该服务最新tag作为该参数值）

Example:
  ./tools/helm.sh install default product1 product1-app1 demo 20190929101615
  ./tools/helm.sh upgrade default product1 product1-app1 demo
  ./tools/helm.sh template default product1 product1-app1 demo
```

对于 product2 这种通过 .gotmpl 渲染生成 value 的需要先使用 tools/helmfile.sh 进行渲染
```sh
./tools/helmfile.sh -h
Usage:
  ./tools/helmfile.sh  template|sync product2-app1 demo 20190929101615 

Example:
  ./tools/helmfile.sh template product2-app1 demo 20190929101615
  
完成后可以看到 value 文件已生成： product2-app1-demo.yaml 
# ll ../helmfile.d/config/product2-demo/product2-app1-demo.yaml
-rw------- 1 root root 1307 Jul 22 17:15 ../helmfile.d/config/product2-demo/product2-app1-demo.yaml
# diff ../helmfile.d/config/product2-demo/product2-app1-demo.yaml*
2c2
<   replicaCount: "1"
---
>   replicaCount: "{{ .Environment.Values.global | getOrNil "replicaCount" | default "1" }}"
19,26c19,24
<     - name: APP_ENV_FLAG
<       value: "demo"
<     - name: HOST
<       value: "0.0.0.0"
<     - name: LOG_PATH
<       value: "/data/logs"
<     - name: PORT0
<       value: "80"
---
>    {{- range $key, $val := .Environment.Values.global.ENV }}
>    {{- if and (ne $key "execute_key1") (ne $key "execute_key2") }}
>     - name: {{ $key }}
>       value: {{ $val | quote }}
>    {{- end }}
>    {{- end }}
```



# 自定义插件

## 定义环境变量
编辑 /root/.zshrc
export HELMPATH=/root/helm

## 插件安装
helm plugin install $HELMPATH/plugins/tools

## 插件更新
helm plugin remove tools
helm plugin install $HELMPATH/plugins/tools

## 插件使用
helm tools template product1 product1-app1 demo xxxxxxx


# 之前的操作方式
因为tools本做了对其他服务的兼容，推荐统一使用 `./tools/helm.sh` 操作，之前的操作暂时先注释处理
```
>
># 常用操作
>
>注意 `helm` 是 helm 2 还是 helm 3
>
>- 更新依赖
>`helm dep up xxx`
>
>- 渲染模版 - 调试
>`helm template xxx -f config/xxx/productx-appx-demo.yaml`
>
>- 更新
>`helm upgrade product1-app1-demo product1-app1 -f config/product1-demo/product1-app1-demo.yaml --set bootstrap.image.tag=xxx`
>
># tools 的使用
>
>```
>cd product1
>
>* ./tools/helm.sh --- 更新 chart
>Usage: ./tools/helm.sh template|upgrade product1-admin ty 20190929101615
>
>* ./tools/update_dependencies.sh --- 更新所有依赖
>
>* ./tools/create_chart.sh product1-new --- 新建 chart, 如 product1-new
>```

# 使用 helmfile 方式部署
以 product1 在 demo 环境部署为例：
```shell
# helmfile -e demo -f releases/product1-demo.yaml sync
```

```shell
# helmfile -f releases/product1-demo.yaml list
NAME                    NAMESPACE       INSTALLED       LABELS
product1-app1-default                   true            app:product1-app1,chart:product1-app1,component:product1-app1
product1-app2-default                   true            app:product1-app2,chart:product1-app2,component:product1-app2
product1-helm-default                   true            chart:product1-helm,component:product1-helm,app:product1-helm

# helm list
NAME                            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                           APP VERSION
product1-app1-demo              default         1               2020-07-22 17:01:27.555893238 +0800 CST deployed        product1-app1-0.1.0             1.16.0
product1-app2-demo              default         1               2020-07-22 17:01:27.556033504 +0800 CST deployed        product1-app2-0.1.0             1.16.0
product1-helm-demo              default         1               2020-07-22 17:01:27.55151997 +0800 CST  deployed        product1-helm-0.1.0             1.16.0
```
```shell
# kubectl get pods -n demo
NAME                                                 READY   STATUS             RESTARTS   AGE
product-helm-685b8f4cf5-qkwh8                        0/1     ImagePullBackOff   0          111s
product1-app1-6f57f48cdd-f5vj8                       1/1     Running            0          110s
product1-app2-7dd495659b-wq2rp                       1/1     Running            0          110s
```
