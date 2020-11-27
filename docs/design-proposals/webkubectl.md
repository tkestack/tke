# Cluster Credential


**Author**: HansenWu ([@HansenWu](https://github.com/jacky68147527))

**Status** (20201124): In development

## Abstract

webkubectl 是一个 web 方式执行`kubectl`工具，通过 ServiceAccount 授权登陆 kubernetes 集群，避免登陆宿主机操作，提高效率。

## Motivation

```
1、前端基于 xterm.js，后台采用原始 client-go 开发，无其他依赖，支持 websocket 连接状态提示
2、通过 kubernetes 原生的 RBAC 授权，workload 仅需设置 serviceAccountName 即可使用
3、提供 yaml 及 helm charts 的 2 种安装方式，方便用户快速安装
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
- yaml 方式安装
    ``` sh
    kubeclt create -f https://raw.githubusercontent.com/tkestack/tke/master/docs/design-proposals/webkubectl-files/webkubectl.yaml
    ```
- helm 方式安装

    1 前提条件：确认是否已经添加 helm repo,如果没有可以参考[此处](https://helm.sh/docs/helm/helm_repo_add/)

    2 添加 chart 到 repo
    ``` sh
    cd webkubectl-files
    helm package webkubectl --save=false
    ```
    3 安装 chart

    ``` sh
    helm install webkubectl
    ```

### 使用方式
浏览器访问：
```
http://${tkestack-domain}/webtty.html?clusterName=${cluster}&podName=webkubectl-agent-0&containerName=webkubectl-agent&namespace=default
```
