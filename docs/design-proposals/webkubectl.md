# Cluster Credential


**Author**: HansenWu ([@HansenWu](https://github.com/jacky68147527))

**Status** (20201124): In development

## Abstract

webkubectl是一个web方式执行`kubectl`工具，通过ServiceAccount授权登陆kubernetes集群，避免登陆宿主机操作，提高效率。

## Motivation

```
1、前端基于xterm.js，后台采用原始client-go开发，无其他依赖，支持websocket连接状态提示
2、通过kubernetes原生的RBAC授权，workload仅需设置serviceAccountName即可使用
3、提供yaml及helm charts的2种安装方式，方便用户快速安装
```

## Main proposal

### 实现原理
<pre>

                                                                                  ┌─────────────────────────────┐
                                                                                  │     kubernetes cluster      │
                                                                                  ├─────────────────────────────┤
                                                                                  │                             │
                                                                                  │     ┌───────────────────┐   │
                                                                                  │     │  service account  │   │
                                                                                  │     └─────────┬─────────┘   │
                                                                                  │               │             │
                                                                                  │               │             │
┌────────────┐                 ┌────────────────────────┐                         │     ┌─────────▼─────────┐   │
│web browser │────────────────▶│        gateway         │─────────────────────────┼────▶│ webkubectl-agent  │   │
└────────────┘                 └────────────────────────┘        websocket        │     └───────────────────┘   │
                                                                                  │                             │
                                                                                  │                             │
                                                                                  │                             │
                                                                                  └─────────────────────────────┘                                                                           │                             │
</pre>

### 安装部署
- yaml方式安装
```
kubeclt create -f https://raw.githubusercontent.com/tkestack/tke/master/docs/design-proposals/webkubectl-files/webkubectl.yaml
```
- helm方式安装
```
1、前提条件：确认是否已经添加helm repo,如果没有可以参考[此处](https://helm.sh/docs/helm/helm_repo_add/)
2、添加chart到repo
    cd webkubectl-files
    helm package webkubectl --save=false
3、安装chart
    helm install webkubectl
```

### 使用方式
浏览器访问：
```
http://${tkestack-domain}/webtty.html?clusterName=${cluster}&podName=webkubectl-agent-0&containerName=webkubectl-agent&namespace=default
```
