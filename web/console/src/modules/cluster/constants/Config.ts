/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import { checkCustomVisible } from '@src/modules/common/components/permission-provider';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
/** ========================= start FFRedux的相关配置 ======================== */
export const FFReduxActionName = {
  CLUSTER: 'cluster',
  REGION: 'region',
  COMPUTER: 'computer',
  MACHINE: 'machine',
  DETAILEVENT: 'DETAILEVENT',

  Resource_Workload: 'Resource_Workload',
  Resource_Detail_Info: 'Resource_Detail_Info',
  LBCF_DRIVER: 'LBCF_DRIVER',

  COMPUTER_WORKLOAD: 'COMPUTER_WORKLOAD'
};
/** ========================= end FFRedux的相关配置 ======================== */

/** ========================= start 集群的相关配置 ======================== */
/** 集群的类型 */
export const ClusterTypeMap = {
  imported: '导入集群',
  baremetal: '独立集群'
};
/** ========================= end 集群的相关配置 ======================== */

/** ========================= resource的相关配置 -start ======================== */

/** pod数量的最大数量限制 */
export const ContainerMaxNumLimit = 65535;

/** pod数量的最小数量限制 */
export const ContainerMinNumLimit = 0;

/** start --- 各种resource的状态展示，在这里做一个统一的入口，因为resourceTablePanel里面去区分类型 */
/**命名空间状态 */
export const NamespaceStatus = {
  Active: {
    text: 'Active',
    classname: 'text-success'
  },
  Available: {
    text: 'Active',
    classname: 'text-success'
  },
  Terminating: {
    text: 'Terminating',
    classname: 'text-restart'
  },
  Failed: {
    text: 'Failed',
    classname: 'text-danger'
  }
};

export const PvcStatus = {
  Available: {
    text: 'Available',
    classname: 'text-success'
  },
  Bound: {
    text: 'Bound',
    classname: 'text-success'
  },
  Released: {
    text: 'Released',
    classname: 'text-restart'
  },
  Failed: {
    text: 'Failed',
    classname: 'text-danger'
  },
  Pending: {
    text: 'Pending',
    classname: 'text-danger'
  },
  Lost: {
    text: 'Lost',
    classname: 'text-danger'
  }
};

export const ResourceStatus = {
  np: NamespaceStatus,
  pvc: PvcStatus,
  pv: PvcStatus
};
/** end --- 各种resource的状态展示，在这里做一个统一的入口，因为resourceTablePanel里面去区分类型 */

/** 是否需要展示resourceLoading */
export const ResourceLoadingIcon = {
  npDelete: ['Terminating']
};

/** 是否需要判断loading状态 */
export const ResourceNeedJudgeLoading = [
  'np',
  'deployment',
  'svc',
  'ingress',
  'pvc',
  'statefulset',
  'daemonset',
  'tapp'
];

/** 创建pvc页面的 云盘数据类型的映射  */
export const DiskTypeName = {
  CLOUD_BASIC: t('普通云硬盘'),
  CLOUD_PREMIUM: t('高性能云硬盘'),
  CLOUD_SSD: t('SSD云硬盘'),
  cbs: t('普通云硬盘')
};

/** container的状态 */
export const ContainerStatusMap = {
  running: {
    text: 'Running',
    classname: 'text-success'
  },
  terminated: {
    text: 'Terminated',
    classname: 'text-danger'
  },
  waiting: {
    text: 'Waiting',
    classname: 'text-restart'
  }
};

/** resource detail当中 日志的 tailList */
export const TailList = [
  {
    value: '100',
    label: t('100条数据')
  },
  {
    value: '200',
    label: t('200条数据')
  },
  {
    value: '500',
    label: t('500条数据')
  },
  {
    value: '1000',
    label: t('1000条数据')
  }
];

/** resource 当中的类型 */
export const ResourceTypeList = [
  {
    value: 'deployment',
    label: t('Deployment（可扩展的部署Pod）')
  },
  {
    value: 'daemonset',
    label: t('DaemonSet（在每个主机上运行Pod）')
  },
  {
    value: 'statefulset',
    label: t('StatefulSet（有状态集的运行Pod）')
  },
  {
    value: 'cronjob',
    label: t('CronJob（按照Cron的计划定时运行）')
  },
  {
    value: 'job',
    label: t('Job（单次任务）')
  },
  {
    value: 'tapp',
    label: t('TApp（可对指定pod进行删除、原地升级、独立挂盘等）')
  }
];

/** 创建workload，hpa的指标选项列表 */
export const HpaMetricsTypeList = [
  {
    value: 'cpuUtilization',
    label: t('CPU利用率')
  },
  {
    value: 'memoryUtilization',
    label: t('内存利用率')
  },
  {
    value: 'cpuAverage',
    label: t('CPU使用量')
  },
  {
    value: 'memoryAverage',
    label: t('内存使用量')
  },
  {
    value: 'inBandwidth',
    label: t('入带宽')
  },
  {
    value: 'outBandwidth',
    label: t('出带宽')
  }
];
/**tapp - 节点异常策略 */
export const NodeAbnormalStrategy = [
  {
    value: 'true',
    text: t('迁移')
  },
  {
    value: 'false',
    text: t('不迁移')
  }
];
/** 创建workload，重启策略的类型 */
export const RestartPolicyTypeList = [
  {
    value: 'OnFailure',
    label: 'OnFailure'
  },
  {
    value: 'Never',
    label: 'Never'
  }
];

/**亲和性调度操作符 */

export const affinityOperator = {
  In: 'In',
  NotIn: 'NotIn',
  Exists: 'Exists',
  DoesNotExits: 'DoesNotExits',
  Gt: 'Gt',
  Lt: 'Lt'
};

//**亲和性调度方式："node" 指定节点调度 "rule" 自定义规则 "unset"

export const affinityType = {
  node: 'node',
  rule: 'rule',
  unset: 'unset'
};
/**
 * 服务调度的操作符
 */
export const affinityRuleOperator = [
  {
    value: 'In',
    tip: t('Label的value在列表中')
  },
  {
    value: 'NotIn',
    tip: t('Label的value不在列表中')
  },
  {
    value: 'Exists',
    tip: t('Label的key存在')
  },
  {
    value: 'DoesNotExist',
    tip: t('Labe的key不存在')
  },
  {
    value: 'Gt',
    tip: t('Label的值大于列表值（字符串匹配）')
  },
  {
    value: 'Lt',
    tip: t('Label的值小于列表值（字符串匹配）')
  }
];

/** 创建 pv的来源设置的列表 */
export const PvCreateSourceList = [
  {
    value: 'static',
    name: t('静态创建')
  },
  {
    value: 'dynamic',
    name: t('动态创建')
  }
];

/** 创建pv的 文件系统 */
export const PvFsTypeList = [
  {
    value: 'ext4',
    label: 'ext4'
  }
];

/** 创建service的时候，workload的列表 */
export const ServiceWorkloadList = [
  {
    value: 'deployment',
    name: 'Deploymemt'
  },
  {
    value: 'statefulset',
    name: 'StatefulSet'
  },
  {
    value: 'daemonset',
    name: 'DaemonSet'
  },
  ...(checkCustomVisible('platform.cluster.service.service_create_tapp')
    ? [
        {
          value: 'tapp',
          name: 'TApp'
        }
      ]
    : []),
  {
    value: 'vmi',
    name: 'VirtualMachines'
  }
];

/** 创建storageClass，云盘的计费方式 */
export const StorageClassCbsPayModeList = [
  {
    value: 'POSTPAID',
    name: t('按量计费')
  },
  {
    value: 'PREPAID',
    name: t('包年包月')
  }
];

/** 创建storageclass，回收策略 */
export const ReclaimPolicyTypeList = [
  {
    value: 'Delete',
    name: t('删除'),
    disabled: false
  },
  {
    value: 'Retain',
    name: t('保留')
  }
];

/** 创建secret 的类型列表 */
export const SecretTypeList = [
  {
    value: 'Opaque',
    name: 'Opaque'
  },
  {
    value: 'kubernetes.io/dockercfg',
    name: 'Dockercfg'
  }
];

/** 创建pvc 读写权限 */
export const PvcAndpvAccessModeList = [
  {
    value: 'ReadWriteOnce',
    name: t('单机读写')
  }
  // {
  //     value: 'ReadOnlyMany',
  //     name: t('多机只读')
  // },
  // {
  //     value: 'ReadWriteMany',
  //     name: t('多机读写')
  // }
];

/** 创建服务- 数据卷类型列表 */
export const VolumeTypeList = [
  {
    value: 'emptyDir',
    label: t('使用临时目录')
  },
  {
    value: 'hostPath',
    label: t('使用主机路径')
  },
  {
    value: 'nfsDisk',
    label: t('使用NFS盘')
  },
  {
    value: 'pvc',
    label: t('使用已有PVC')
  },
  {
    value: 'configMap',
    label: t('使用ConfigMap')
  },
  {
    value: 'secret',
    label: t('使用Secret')
  }
];

/** 创建workload，挂载点的 模式选择 */
export const VolumeMountModeList = [
  {
    value: 'rw',
    label: t('读写')
  },
  {
    value: 'ro',
    label: t('只读')
  }
];

/** 创建workload，健康检查方法 */
export const HealthCheckMethodList = [
  {
    value: 'methodTcp',
    label: t('TCP端口检查')
  },
  {
    value: 'methodHttp',
    label: t('HTTP请求检查')
  },
  {
    value: 'methodCmd',
    label: t('执行命令检查')
  }
];

/** 创建workload，健康检查协议 */
export const HttpProtocolTypeList = [
  {
    value: 'HTTP',
    label: 'HTTP'
  },
  {
    value: 'HTTPS',
    label: 'HTTPS'
  }
];

/** 节点状态 */
export const NodeStatus = {
  AllNormal: {
    text: t('全部正常'),
    classname: 'text-success'
  },
  AllAbnormal: {
    text: t('全部异常'),
    classname: 'text-danger'
  },
  PartialAbnormal: {
    text: t('部分异常'),
    classname: 'text-danger'
  },
  '-': {
    text: '-',
    classname: 'text-restart'
  }
};

/** 全局轮询事件，写在配置文件中，不需要每次手打，容易打错 */
export const PollEventName = {
  resourceDetailEvent: 'pollResourceDetailEvent',
  resourcePodLog: 'pollResourcePodLog',
  resourceLog: 'pollResourceLog',
  resourceEvent: 'pollResourceEvent',
  resourceList: 'pollResourceList',
  resourcePodList: 'pollResourcePodList'
};

/** 创建service当中的 访问方式 */
export const CommunicationTypeList = [
  {
    value: 'ClusterIP',
    label: t('仅在集群内访问'),
    tip: t(
      '将提供一个可以被集群内其他服务或容器访问的入口，支持TCP/UDP协议，数据库类服务如Mysql可以选择集群内访问,来保证服务网络隔离性。'
    )
  },
  {
    value: 'NodePort',
    label: t('主机端口访问'),
    tip: t('提供一个主机端口映射到容器的访问方式，支持TCP&UDP， 可用于业务定制上层LB转发到Node。')
  }
];

/**externalTrafficPolicy */
export const ExternalTrafficPolicy = {
  Cluster: 'Cluster',
  Local: 'Local'
};

export const SessionAffinity = {
  ClientIP: 'ClientIP',
  None: 'None'
};

/** 协议列表 */
export const ProtocolList = [
  {
    value: 'TCP',
    label: 'TCP'
  },
  {
    value: 'UDP',
    label: 'UDP'
  }
];

/** 镜像的更新策略 */
export const ImagePullPolicyList = [
  {
    value: 'Always',
    text: t('Always（总是拉取）')
  },
  {
    value: 'Never',
    text: t('Never（不拉取）')
  },
  {
    value: 'IfNotPresent',
    text: t('IfNotPresent（镜像不存在时拉取）')
  }
];

/** 创建workload的网络模式 */
export enum WorkloadNetworkTypeEnum {
  Overlay = 'overlay',
  FloatingIP = 'floatingip',
  Nat = 'nat',
  Host = 'host'
}

export const WorkloadNetworkType = [
  {
    value: WorkloadNetworkTypeEnum.Overlay,
    text: t('Overlay（虚拟网络）')
  },
  {
    value: WorkloadNetworkTypeEnum.FloatingIP,
    text: t('FloatingIP（浮动IP）')
  },
  {
    value: WorkloadNetworkTypeEnum.Nat,
    text: t('Nat（端口映射）')
  },
  {
    value: WorkloadNetworkTypeEnum.Host,
    text: t('Host（主机网络）')
  }
];

export const FloatingIPReleasePolicy = [
  {
    value: 'immutable',
    text: t('缩容或删除APP时回收')
  },
  {
    value: 'never',
    text: t('永不回收')
  },
  {
    value: 'always',
    text: t('随时回收')
  }
];
/** ========================= resource的相关配置 -end ======================== */

/** ========================= 创建独立集群 的相关配置 start ======================== */
export const k8sVersionList = [
  { text: '1.14.1', value: '1.14.1' },
  {
    text: '1.12.8',
    value: '1.12.8'
  }
];
export const computerRoleList = [
  {
    text: 'Master&Etcd',
    value: 'master_etcd'
  }
];

export const authTypeMapping = {
  password: 'password',
  cert: 'cert'
};
export const authTypeList = [
  {
    text: t('密码认证'),
    value: 'password'
  },
  {
    text: t('密钥认证'),
    value: 'cert'
  }
];
/** ========================= 创建独立集群 的相关配置 end ======================== */

/** 协议列表 */
export const LbcfProtocolList = [
  {
    value: 'TCP',
    text: 'TCP'
  },
  {
    value: 'UDP',
    text: 'UDP'
  }
];

export const LbcfConfig = [
  {
    text: 'CLBID',
    value: 'loadBalancerID',
    input: {
      placeholder: t('不填则自动创建')
    }
  },
  {
    text: t('CLB类型'),
    value: 'loadBalancerType',
    select: {
      //CLB实例类型，OPEN为公网CLB，INTERNAL为内网CLB
      options: [
        {
          text: t('公网CLB'),
          value: 'OPEN'
        },
        {
          text: t('内网CLB'),
          value: 'INTERNAL'
        }
      ]
    },
    defaultValue: 'OPEN'
  },
  {
    text: 'VPCID',
    value: 'vpcID'
  },
  {
    text: t('子网ID'),
    value: 'subnetID'
  },
  {
    text: t('监听器端口'),
    value: 'listenerPort'
  },
  {
    text: t('监听器端口类型'),
    value: 'listenerProtocol',
    select: {
      //TCP,UDP,HTTP,HTTPS
      options: [
        {
          text: 'TCP',
          value: 'TCP'
        },
        {
          text: 'UDP',
          value: 'UDP'
        },
        {
          text: 'HTTP',
          value: 'HTTP'
        },
        {
          text: 'HTTPS',
          value: 'HTTPS'
        }
      ]
    },
    defaultValue: 'TCP'
  },
  {
    text: t('域名'),
    value: 'domain'
  },
  {
    text: t('路径'),
    value: 'url'
  }
];

export const LbcfArgsConfig = [
  {
    text: '证书ID',
    value: 'listenerCertID'
  }
];

export const clearNodeSH = `#!/bin/bash

# common
kubeadm reset -f
rm -fv /root/.kube/config
rm -rfv /etc/kubernetes
rm -rfv /var/lib/kubelet
rm -rfv /var/lib/etcd
rm -rfv /var/lib/cni
rm -rfv /etc/cni
rm -rfv /var/lib/tke-registry-api
rm -rfv /opt/tke-installer
rm -rfv /var/lib/postgresql /etc/core/token /var/lib/redis /storage /chart_storage
ip link del cni0 2>/etc/null

for port in 80 443 2379 2380 6443 8086 8181 9100 30086 31138 31180 31443  {10249..10259} ; do
    fuser -k -9 \${port}/tcp
done

# docker
docker rm -f $(docker ps -aq) 2>/dev/null
systemctl disable docker 2>/dev/null
systemctl stop docker 2>/dev/null
rm -rfv /etc/docker
ip link del docker0 2>/etc/null

# containerd
nerdctl rm -f $(nerdctl ps -aq) 2>/dev/null
ip netns list | cut -d' ' -f 1 | xargs -n1 ip netns delete 2>/dev/null
systemctl disable containerd 2>/dev/null
systemctl stop containerd 2>/dev/null
rm -rfv /var/lib/nerdctl/*

## ip link
ip link delete cilium_net 2>/dev/null
ip link delete cilium_vxlan 2>/dev/null
ip link delete flannel.1 2>/dev/null

## iptables
iptables --flush
iptables --flush --table nat
iptables --flush --table filter
iptables --table nat --delete-chain
iptables --table filter --delete-chain

# reboot
reboot now`;

export enum GPUTYPE {
  PGPU = 'Physical',
  VGPU = 'Virtual'
}

export enum BackendType {
  Pods = 'Pods',
  Service = 'Service',
  Static = 'Static'
}

export const BackendTypeList = [
  { text: BackendType.Pods, value: BackendType.Pods },
  { text: BackendType.Service, value: BackendType.Service },
  { text: BackendType.Static, value: BackendType.Static }
];

export enum CreateICVipType {
  unuse = 'unuse',
  existed = 'existed',
  tke = 'tke'
}

export const CreateICVipTypeOptions = [
  { text: '不使用', value: CreateICVipType.unuse },
  { text: '使用已有', value: CreateICVipType.existed },
  { text: 'TKE提供', value: CreateICVipType.tke }
];

export const CreateICCiliumOptions = [
  { text: 'Galaxy', value: 'Galaxy' },
  { text: 'Cilium', value: 'Cilium' }
];

export const NetworkModeOptions = [
  { text: 'Overlay', value: 'overlay' },
  { text: 'Underlay', value: 'underlay' }
];

export enum MachineStatus {
  Running = 'Running',
  Initializing = 'Initializing',
  Failed = 'Failed',
  Terminating = 'Terminating'
}

/** 增加 Capabilities 选项 */
export const AddCapabilitiesList = [
  'SYS_MODULE',
  'SYS_RAWIO',
  'SYS_PACCT',
  'SYS_ADMIN',
  'SYS_NICE',
  'SYS_RESOURCE',
  'SYS_TIME',
  'SYS_TTY_CONFIG',
  'AUDIT_CONTROL',
  'MAC_OVERRIDE',
  'MAC_ADMIN',
  'NET_ADMIN',
  'SYSLOG',
  'DAC_READ_SEARCH',
  'LINUX_IMMUTABLE',
  'NET_BROADCAST',
  'IPC_LOCK',
  'IPC_OWNER',
  'SYS_PTRACE',
  'SYS_BOOT',
  'LEASE',
  'WAKE_ALARM',
  'BLOCK_SUSPEND',
  'all'
];

/** 删除 Capabilities 选项 */
export const DropCapabilitiesList = [
  'SETPCAP',
  'MKNOD',
  'AUDIT_WRITE',
  'CHOWN',
  'NET_RAW',
  'DAC_OVERRIDE',
  'FOWNER',
  'FSETID',
  'KILL',
  'SETGID',
  'SETUID',
  'NET_BIND_SERVICE',
  'SYS_CHROOT',
  'SETFCAP',
  'all'
];

/** pod远程登录的选项 */
export const podRemoteShellOptions = [
  {
    value: '/bin/bash',
    text: '/bin/bash'
  },
  {
    value: '/bin/zsh',
    text: '/bin/zsh'
  },
  {
    value: '/bin/sh',
    text: '/bin/sh'
  }
];

export enum ContainerRuntimeEnum {
  CONTAINERD = 'containerd',
  DOCKER = 'docker'
}

export const ContainerRuntimeOptions = [
  {
    text: ContainerRuntimeEnum.CONTAINERD,
    value: ContainerRuntimeEnum.CONTAINERD
  },

  {
    text: ContainerRuntimeEnum.DOCKER,
    value: ContainerRuntimeEnum.DOCKER
  }
];

export const ContainerRuntimeTips = {
  [ContainerRuntimeEnum.CONTAINERD]: 'containerd是更为稳定的运行时组件，支持OCI标准，不支持docker api',
  [ContainerRuntimeEnum.DOCKER]: 'dockerd是社区版运行时组件，支持docker api'
};
