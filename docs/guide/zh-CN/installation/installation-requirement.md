#  部署环境要求



## 硬件要求

* **最小化部署配置：**

<table>
    <tr>
        <td><strong>安装/业务集群</strong></td>
        <td><strong>节点/集群 </td>
        <td><strong>CPU核数 </td>
        <td><strong>内存</td>
        <td><strong>系统盘</td>
        <td><strong>数量</td>
    </tr>
    <tr>
        <td rowspan="2">安装</td>
        <td>Installer节点</td>
        <td>1</td>
        <td>2G</td>
        <td>50G</td>
        <td>1</td>
    </tr>
    <tr>
        <td>Global集群</td>
        <td>8</td>
        <td>16G</td>
        <td>100G</td>
        <td>1</td>
    </tr>
    <tr>
        <td rowspan="2">业务集群</td>
        <td>Master & etcd</td>
        <td>4</td>
        <td>8G</td>
        <td>100G</td>
        <td>1</td>
    </tr>
    <tr>
        <td>Node</td>
        <td>8</td>
        <td>16G</td>
        <td>100G</td>
        <td>3</td>
    </tr>


</table>



* **推荐配置：**



<table>
    <tr>
        <td><strong>安装/业务集群</strong></td>
        <td><strong>节点/集群 </td>
        <td><strong>CPU核数 </td>
        <td><strong>内存</td>
        <td><strong>系统盘</td>
        <td><strong>数量</td>
    </tr>
    <tr>
        <td rowspan="2">安装</td>
        <td>Installer节点</td>
        <td>1</td>
        <td>2G</td>
        <td>50G</td>
        <td>1</td>
    </tr>
    <tr>
        <td>Global集群</td>
        <td>8</td>
        <td>16G</td>
        <td>100G SSD</td>
        <td>3</td>
    </tr>
    <tr>
        <td rowspan="2">业务集群</td>
        <td>Master & etcd</td>
        <td>16</td>
        <td>32G</td>
        <td>300G SSD</td>
        <td>3</td>
    </tr>
    <tr>
        <td>Node</td>
        <td>16</td>
        <td>32G</td>
        <td>系统盘：100G<br>数据盘：300G （/var/lib/docker） </td>
        <td>>3</td>
    </tr>


</table>


## 软件要求

| 需求项      | 具体要求                                                     |
| ----------- | ------------------------------------------------------------ |
| 操作系统    | Ubuntu 16.04/18.04 LTS (64-bit) <br>CentOS Linux 7.6 (64-bit)<br>Tencent Linux 2.2 |
| kernel 版本 | >= Kernel 3.10.0-957.10.1.el7.x86_64                    |
| ssh         | 确保<br> Installer 节点及其容器<br>Global 集群节点及其容器<br>业务集群节点及其容器<br>之间能够 ssh 互联 |
| swap        | 关闭，如果不满足，系统会有一定几率出现 io 飙升，造成 docker 卡死 |
| 防火墙      | 关闭。Kubernetes 官方要求                                    |
| selinux     | 关闭。Kubernetes 官方要求                                    |
| 时间同步    | 所有服务器要求时间必须同步，误差不得超过 2 秒                |
| 时区        | 所有服务器时区必须统一，建议设置为 `Asia/Shanghai`           |
| 端口        | 所有节点防火墙必须放通放通 SSH（默认22）、80、8080、443、6443 端口 |