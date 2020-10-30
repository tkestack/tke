[TOC]

# oneStep Builder 使用手册    

### builder 配置解析    

```
#!/bin/bash
# Author: yhchen
set -e

BASE_DIR=$(cd `dirname $0` && pwd)
cd $BASE_DIR

# tkestack version, 指定tkestack 版本
version='v1.2.4' 

# business helm's or workerload kubernetes's yaml git url, 指定业务helm git仓库
busi_git_url='https://git.xx.com/yy/helm.git'

# business helm's or workerload kubernetes's yaml git branch will be check out, 指定业务helm git 分支
busi_branch='feature/private'

# all servers for set deploy switch, 只有在新增了组件需要配置
all_servers=("tkestack" "business" "redis" "redis_cluster" "mysql" "prometheus" "kafka" "elk" "nginx_ingress" "minio" "helmtiller" "nfs" "salt_minion" "postgres" "sgikes")

# will be deploy's server set，配置需要部署的组件，它是all_servers 的子集
server_set=("tkestack" "business" "redis" "mysql" "prometheus" "kafka" "elk" "nginx_ingress" "helmtiller" "nfs" "salt_minion" "postgres")

# whether use remote docker registry, if true will be not save business images and copy registry secret; true or false
# 业务是否使用远程镜像仓库
remote_registry='true'
# remote image registry url, if remote images registry need issue crt, please name: ${remote_img_registry_url}.cert.tar.gz 
# on offline-pot-tgz-base dir
# 远程镜像仓库域名/ip
remote_img_registry_url='reg.xx.yy.com'

# builder type just support 'all' or 'custom' , default is all; customize will be pack on demand
# 按需打包/全量打包，默认按需打包
builder_type='custom'

```    

### 修改builder 配置    

- 拷贝配置文件    

   ```
   cp builder.cfg.tpl builder.cfg
   
   ```

-  修改配置文件    

    按上述解析进行builder.cfg 修改。    
    
### 开始打包    

```
sh ./builder.sh -a ${app_env_flag} [-f ${call_function}] [-b ${builder_config_file_name}]
命令解析：    
-a: 和helmfile上定义的app_env_flag 一致，通常为客户的简称；必传参数
-f: 要调用的函数名，默认不传
-b: builder.cfg 文件名，不传时使用builder.cfg；默认不需要传

```

