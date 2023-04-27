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
import { checkCustomVisible } from '@common/components/permission-provider';
import { virtualMachineAPI } from '@src/webApi';
import { t } from '@tencent/tea-app/lib/i18n';

export const firstRouterNameMap = {
  overview: t('概览'),
  cluster: t('集群'),
  mesh: t('服务网格')
};

/** 一些操作，create、update的一些header的名称映射 */
export const typeMapName = {
  create: t('新建'),
  modify: t('更新')
};

/** 二级导航栏的配置文件
 * @param sub   当前的一级导航
 */
export const subRouterConfig = (module = 'cluster'): any => {
  if (module === 'mesh') {
    return meshSubRouterConfig;
  }
  return clusterSubRouterConfig;
};

const clusterSubRouterConfig = [
  {
    name: t('基本信息'),
    path: 'basic',
    basicUrl: 'info'
  },
  {
    name: t('节点管理'),
    path: 'nodeManage',
    sub: [
      {
        name: t('节点'),
        path: 'node'
      }
    ]
  },
  {
    name: t('命名空间'),
    path: 'namespace',
    basicUrl: 'np'
  },
  {
    name: t('工作负载'),
    path: 'resource', // 用于判断哪个二级菜单栏需要展开
    sub: [
      {
        name: 'Deployment',
        path: 'deployment'
      },
      {
        name: 'StatefulSet',
        path: 'statefulset'
      },
      {
        name: 'DaemonSet',
        path: 'daemonset'
      },
      {
        name: 'Job',
        path: 'job'
      },
      {
        name: 'CronJob',
        path: 'cronjob'
      },
      {
        name: 'VirtualMachine',
        path: 'virtual-machine',
        visible: ({ clusterId }) => virtualMachineAPI.checkVmEnable(clusterId)
      }
    ]
  },
  {
    name: t('自动伸缩'),
    path: 'scale',
    sub: [
      {
        name: 'HPA',
        path: 'hpa'
      },
      {
        name: 'CronHPA',
        path: 'cronhpa',
        visible: () => checkCustomVisible('platform.cluster.cronHpa_menu')
      }
    ]
  },
  {
    name: t('服务'),
    path: 'service', // 用于判断哪个二级菜单栏需要展开
    sub: [
      {
        name: 'Service',
        path: 'svc'
      },
      {
        name: 'Ingress',
        path: 'ingress'
      }
      // {
      //   name: t('负载均衡'),
      //   path: 'lbcf'
      // }
    ]
  },
  {
    name: t('配置管理'),
    path: 'config',
    sub: [
      {
        name: 'ConfigMap',
        path: 'configmap'
      },
      {
        name: 'Secret',
        path: 'secret'
      }
    ]
  },
  {
    name: t('存储'),
    path: 'storage',
    sub: [
      {
        name: 'PersistentVolume',
        path: 'pv'
      },
      {
        name: 'PersistentVolumeClaim',
        path: 'pvc'
      },
      {
        name: 'StorageClass',
        path: 'sc'
      }
    ]
  },
  {
    name: t('日志'),
    path: 'k8sLog',
    basicUrl: 'log'
  },
  {
    name: t('事件'),
    basicUrl: 'event',
    path: 'k8sEvent'
  }
];

const meshSubRouterConfig = [
  {
    name: t('基本信息'),
    basicUrl: 'info',
    path: 'basic'
  },
  {
    name: t('网格拓扑'),
    basicUrl: 'topo',
    path: 'dashboard'
  },
  {
    name: t('服务'),
    basicUrl: 'svc',
    path: 'service'
  },
  {
    name: t('Virtual Service'),
    basicUrl: 'virtualservice',
    path: 'virtualservice'
  },
  {
    name: t('Gateway'),
    basicUrl: 'gateway',
    path: 'gateway'
  },
  {
    name: t('组件管理'),
    basicUrl: 'controlPlane',
    path: 'plane'
  }
];

export const notifySubRouter = [
  {
    name: t('通知渠道'),
    id: 'channel',
    basicUrl: 'channel'
  },
  {
    name: t('通知模版'),
    id: 'template',
    basicUrl: 'template'
  },
  {
    name: t('接收人'),
    id: 'receiver',
    basicUrl: 'receiver'
  },
  {
    name: t('接收组'),
    id: 'receiverGroup',
    basicUrl: 'receiverGroup'
  }
];
