[TOC]

# oneStep 开发规范    

**注：除了文档外，其余务必英文。**    

## ansible 使用规范    

##### 1. ansible hosts 配置规范    

- 主机组        
   - 所有子主机组必须包含到all 这个主机组； 
   - 新增组件，只有涉及使用到了local 存储(local storageclass/ host path ) ，主机模式部署 或 docker 方式部署需要固定节点时可以新增子主机组； 
- ansible 全局变量    
    对于ansible 全局变量请在"[all:vars]" 下添加。    
    - 对于新增组件必须在"# switch" 下添加开关变量，变量必须以"deploy\_" 开头(deploy\_\${组件名}), 其他的一些操作(比如压测)推荐归类形成开关项, 从而满足组合式需求；    
- 各个组件的配置项，请添加对应组件注释，然后在组件注释下添加对应配置项变量，比如:        
   
   ```
   # ingress config
  ingress_replica='2' # ingress will be deploy number,default 2
  ingress_host_network='true' # when has LoadBalancer ip for nginx-ingress and just has one master node set false, default true;
  ingress_svc_type='ClusterIP' # when has LoadBalancer ip for nginx-ingress set LoadBalancer, otherwise set ClusterIP
  ingress_lb_ip='' # set LoadBalancer ip for nginx-ingress
   
   ```
   
##### 2. ansible playbooks     

ansible playbooks 目录结构，当前划分为按操作+架构层级划分：主机初始化(hosts-init)，主机检查(hosts-check)，tkestack管理(tke-mgr), 基础组件(base-component), 业务组件(business)。目前基本上已经满足开发需求，基础组件包含了运维基础组件+业务依赖组件。playbooks yml 文件的命名规则： \${目录名}.yml, 比如:  base-component.yml。            

##### 3. ansible roles     

ansible roles 目录结构请和ansible playbooks 保持一致(命名+数量)。    

##### 4. ansible roles 之 handlers    

对于以主机模式启动(如：systemctl 方式启动) 服务，若是涉及配置文件修改需要重启，建议使用handlers实现配置修改后再触发重启等相关操作。    

##### 5. ansible roles 之 tasks     

- tasks 目录下文件命名:    
    tasks目录下的文件命名: \${功能}.yml , 也就是通过yml文件名大概清楚该task 能做什么操作，对于base-component 的tasks 的yml文件请按组件功能命名: \${组件名}-mgr.yml。  

- tasks 目录下文件操作     
    
    对于组件(base-component, business) 需要包括以下能力:    
    - 组件部署    
    - 组件移除    
    - 组件健康检查        

    对于其他的一些操作建议有添加就有移除 原则 。    
    tasks 通用操作规范，包含如下：    
    
    - 充分利用set facts 配置开关项(ansible hosts 上定义), 作为ansible 执行操作条件之一
    - 通过stat 模块 + register 判断执行文件或者一些目录是否存在，作为ansible 执行操作条件之一
    - shell 模块 + register 获取执行结果，作为下一步 ansible执行操作条件之一
    - tasks 的每一个模块操作项必须包含when 进行判断(when 包含执行的主机组，执行条件比如开关)
    - tasks 的每一个模块操作项必须包含tags，作为后续组合执行。比如部署相关采用同一tag，比如dpl_redis, 移除相关用相同tag
    - shell 模块执行shell 命令增加执行条件

- ansible roles 之 templates        
   当部署组件等相关操作涉及到配置文件的配置修改或yaml文件配置修改等通过在ansible hosts 定义变量，ansible templates 作为文件渲染模板，利用ansible template模块进行渲染。    
   
   - 命名    
      \${配置文件名}.j2
   - 基础组件需要按组件名创建模板目录存放模板文件，组件名必须和ansible 开关项的组件名一致(deploy\_\${组件名} 去掉deploy\_)    

对于其他的一些更加有利于解耦，扩展，幂等，简洁的ansible使用规范欢迎加入本规范 !!!    

## shell 脚本    

关于shell 脚本规范如下:    

- 命名    
   命名需要见名知义，也就是看到这个shell 脚本名称基本清楚此脚本作用。

- 统一采用"*#!/bin/bash*"    
- 脚本遇到异常需推出以及配置切换到当前脚本执行路径，具体配置如下:    
   
   ```
   set -e

    BASE_DIR=$(cd `dirname $0` && pwd)
    cd $BASE_DIR
   
   ```
   
- 推荐使用函数方式编写脚本，函数建议包含如下函数:    
    
    - 操作函数, 最终要执行操作的函数    
    - help 函数, 遇到本脚本不支持的函数调用提示当前可以执行哪些函数相关的帮助性信息    
    - main 函数，由main函数作为函数调用统一入口    
    
- 脚本传参推荐使用getopts, 示例:        
    
   ```
   while getopts ":f:h:" opt
  do
    case $opt in
      f)
      CALL_FUN="${OPTARG}"
      ;;
      ?)
      echo "unkown args! just suport -f[call function] arg!!!"
      exit 0;;
    esac
  done
   
   ```    
   
   使用getopts 能更精准传参，但对于参数需要配置一个默认值。

- 条件化调用或执行        
   
   对于函数的调用或者执行的shell 命令强力推荐增加条件判断方式调用或执行。

- 异常处理        
   
   对于遇到异常情况需要做合理处理，输出具有指导性意义的异常信息, 并安全退出。    

- 避免硬编码    

   比如遇到操作某个目录/文件最好组合当前脚本执行目录来获取该目录或文件的路径, 涉及配置通过传参等方式渲染。    

## builder 相关目录命名规范    

在按需打包时，需要根据*builder.cfg* 文件的*server\_set*数组包含的组件名来循环拷贝helm，压缩文件，镜像文件等；为此需要对保存helm，压缩文件，镜像文件的目录命名规范化。    

- builder.cfg 的all\_servers 和 server\_set 数组 规范    

   对于*builder.cfg* 文件的*all\_servers* 和 *server\_set* 数组中的值必须保持和ansible hosts 的开关项的组件名一致(deploy\_\${组件名}, 去掉deploy\_ 即可。)    

- 基础组件helm相关目录命名规范    
   
   基础组件helms打包前处理目录统一放到与*offline-pot* 目录同级的 *base-component-helms* 目录，由于一个组件可能会涉及不少helm chart定义, 并且通过for 循环方式拷贝组件helm，所以规范如下:    
   
   - 创建一个和ansible hosts 的开关项的组件名一致(deploy\_\${组件名}, 去掉deploy\_ 即可。) 目录来存放组件的helm   
   
- 镜像目录命名规范    
   
   组件镜像打包前处理目录统一放到与*offline-pot* 目录同级的*offline-pot-images-base* 目录，为了更清晰明了,按需打包更简洁，规范如下：    
   
   - 创建一个和ansible hosts 开关项的组件名一致(deploy\_\${组件名}, 去掉deploy\_ 即可。) 目录来存放组件相关镜像
   - 镜像文件必须是*.tar 
       
- 压缩包目录命名规范    

   组件相关压缩包打包前处理目录统一放到与*offline-pot* 目录同级的*offline-pot-tgz-base* 目录， 规范如下:    
   
   - 创建一个和ansible hosts 开关项的组件名一致(deploy\_\${组件名}, 去掉deploy\_ 即可。) 目录来存放组件相关压缩包    
   - 对于远程镜像仓库若需要下发证书，证书附件文件目录名为 *builder.cfg* 文件中remote\_img\_registry\_url+.cert, 也就是\${remote\_img\_registry\_url}.cert
   - 镜像证书文件名为: \${remote\_img\_registry\_url}.cert.tar.gz
    
- 远程镜像仓库secrets 目录命名规范    
   
   对于使用远程镜像仓库作为业务镜像仓库，并且需要认证的secret文件，打包前处理目录创建在与 *offline-pot* 目录同级，命名规范如下:    
   
   - *builder.cfg* 文件中remote\_img\_registry\_url+.secrets，也就是\${remote\_img\_registry\_url}.secrets 作为目录名
   - secrets 文件命名：\${中心名}.\${remote\_img\_registry\_url}.yaml 或 \${产品名}.\${remote\_img\_registry\_url}.yaml 这个需要和helmfile 的helm chart 父目录保持一致。    

   
## 业务helmfile 规范    

当前安装包仅支持helmfile方式部署业务，可以自行扩展。当前说下helmfile 规范。    

目录结构如下:    

```    

|-- README.md
|-- charts  # 固定名
|-- helmfile.d # 固定名
`-- ${产品名}/${中心名}

```    

helmfile.d 目录下的目录命名规范：    

```
.
├── commons
│   ├── helmdefaults.yaml
│   └── repos.yaml
├── config
│   ├── cp1-kh1      # 产品名-环境名，也就是${产品名}-${客户简称}, ${产品名}-${app_env_flag}
│   │   ├── environments-values.yaml              # 环境 helmfile values 用于渲染 .gotmpl
│   │   ├── agent-cp1-kh1.yaml.gotmpl    # gotmpl 模板, ${服务名}-${产品名}-${app_env_flag}
│   │   ├── agent-cp1-kh1.yaml              # environments-values.yaml + gotmpl 渲染出的 helm values.yaml
│   │   ├── xpay-cp1-kh1.yaml.gotmpl
│   │   └── xpay-cp1-kh1.yaml
│   ├── cp1-kh2
│   │   ├── ...
│   │   └── web-private-cp1-kh2.yaml.gotmpl
│   └── cp2-kh3
│       ├── environments-values.yaml       # 用于渲染 .gotmpl
│       ├── ...
│       ├── eureka-01-cp2-kh3.yaml.gotmpl
│       └── eureka-01-cp2-kh3.yaml
├── helmfile.yaml
└── releases
    ├── cp2-kh3.yaml   
    ├── cp1-kh1.yaml                    # 产品名-环境名 和 config 目录里对应，也就是${产品名}-${客户简称}, ${产品名}-${app_env_flag}
    ├── ...
    └── cp1-kh2.yaml

```    

## git 忽略提交    

oneStep不能包含基础组件以及业务组件相关的helm，所以需要在.gitignore 忽略掉提交：    

```
cat .gitignore 
# ignore hosts file
hosts
roles/base-component/helms/*
roles/business/helms/*

```   

以上忽略的项必须配置。    

## 本地镜像仓库配置    

涉及到本地镜像仓库域名以及镜像仓库namespace的请按如下规范：    

- 镜像仓库namespace统一为:   library; 也就是在hosts.tpl 上配置如下所示：    

   ```
   POSTGRES_IMAGE='library/${image_name}:${image_tag}' 
   
   OR
   
   POSTGRES_IMAGE='library/${image_name}'
   POSTGRES_TAG='${image_tag}'
   
   这些根据实际情况选择以上两种配置!!!
   
   ```

- 至于本地镜像仓库域名则统一在ansible 模板文件渲染：    

   ```
   cat hosts.tpl
   
   ...
   
   # offline registry config
   registry_domain='registry.tke.com' # offline registry domain, 统一使用此镜像仓库域名作为本地镜像仓库域名
   
   ...
   
   类似：   
   repository: {{ registry_domain }}/{{ minio_mcimg_name }}
   tag: {{ minio_mcimg_tag }}
   
   
   ```

