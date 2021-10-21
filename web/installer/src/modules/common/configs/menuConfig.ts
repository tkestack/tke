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
export const businessKey = 'oss';

export const defaultMenu = `/${businessKey}/userinfo`;

export const menuConfig = [
  {
    group: '运营中心',
    menus: [
      {
        title: '概览',
        href: `/${businessKey}/overview`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/overview.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/overview-hover.svg'
        ]
      },
      {
        title: '资源统计',
        href: `/${businessKey}/statistic`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/service.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/service-hover.svg'
        ],
        sub: [
          {
            title: '整体规模',
            href: `/${businessKey}/scale`
          },
          {
            title: '用户资源趋势',
            href: `/${businessKey}/trend`
          },
          {
            title: 'Top数据统计',
            href: `/${businessKey}/top`
          },
          {
            title: '节点可用性',
            href: `/${businessKey}/usability`
          },
          {
            title: '自建K8S情况',
            href: `/${businessKey}/contruct`
          },
          {
            title: '中心资源统计',
            href: `/${businessKey}/paas`
          }
        ]
      },
      {
        title: '数据配置',
        href: `/${businessKey}/data`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/CloudDatabase.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/CloudDatabase-hover.svg'
        ],
        sub: [
          {
            title: '地域配置',
            href: `/${businessKey}/region`
          },
          {
            title: 'K8S版本配置',
            href: `/${businessKey}/k8s`
          },
          {
            title: '镜像配置',
            href: `/${businessKey}/image`
          },
          {
            title: '文档配置',
            href: `/${businessKey}/document`
          }
        ]
      },
      {
        title: '用户配置',
        href: `/${businessKey}/user`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/StorageBlock.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/StorageBlock-hover.svg'
        ],
        sub: [
          {
            title: '配额管理',
            href: `/${businessKey}/quota`
          },
          {
            title: '白名单管理',
            href: `/${businessKey}/whitelist`
          }
        ]
      }
    ]
  },
  {
    group: '运维中心',
    menus: [
      {
        title: '用户资源查询',
        href: `/${businessKey}/userinfo`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/Templates.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/Templates-hover.svg'
        ]
      },
      {
        title: '日志查询',
        href: `/${businessKey}/log`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/log.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/log-hover.svg'
        ],
        sub: [
          {
            title: '云API日志',
            href: `/${businessKey}/apiLog`
          },
          {
            title: 'Dashboard日志',
            href: `/${businessKey}/dashboardLog`
          },
          {
            title: 'GW日志',
            href: `/${businessKey}/gwLog`
          },
          {
            title: 'ApiServer日志',
            href: `/${businessKey}/apiServerLog`
          },
          {
            title: 'Norm日志',
            href: `/${businessKey}/normLog`
          },
          {
            title: 'NightsWatch日志',
            href: `/${businessKey}/nightsWatchLog`
          },
          {
            title: 'EniController日志',
            href: `/${businessKey}/eniControllerLog`
          },
          {
            title: 'EnilbController日志',
            href: `/${businessKey}/enilbControllerLog`
          }
        ]
      },
      {
        title: 'ETCD服务平台',
        href: `/${businessKey}/etcd`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/CLB.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/CLB-hover.svg'
        ],
        sub: [
          {
            title: 'ETCD集群管理',
            href: `/${businessKey}/etcdCluster`
          },
          {
            title: 'ETCD调度策略管理',
            href: `/${businessKey}/etcdHpa`
          },
          {
            title: 'ETCD迁移任务管理',
            href: `/${businessKey}/etcdJob`
          },
          {
            title: 'ETCD备份管理',
            href: `/${businessKey}/etcdBackup`
          },
          {
            title: '集群巡检',
            href: `/${businessKey}/bigCluster`
          },
          {
            title: 'ETCD监控',
            href: `/${businessKey}/etcdMonitor`
          }
        ]
      },
      {
        title: '工具管理',
        href: `/${businessKey}/tools`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/helm.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/helm-hover.svg'
        ],
        sub: [
          {
            title: 'License管理',
            href: `/${businessKey}/license`
          },
          {
            title: '告警规则管理',
            href: `/${businessKey}/alarm`
          },
          {
            title: 'ApiServer监控',
            href: `/${businessKey}/monitor`
          },
          {
            title: '集群Master升级',
            href: `/${businessKey}/master`
          },
          {
            title: '集群Agent升级',
            href: `/${businessKey}/agent`
          },
          {
            title: '集群收集',
            href: `/${businessKey}/collect`
          },
          {
            title: 'Web Shell',
            href: `/${businessKey}/shell`
          }
        ]
      }
    ]
  },
  {
    group: '质量管理中心',
    menus: [
      {
        title: '前台质量管理',
        href: `/${businessKey}/front`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/Flow.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/Flow-hover.svg'
        ],
        sub: [
          {
            title: '控制台访问量统计',
            href: `/${businessKey}/visit`
          },
          {
            title: '前台异常上报',
            href: `/${businessKey}/feAbnormal`
          },
          {
            title: '页面测速上报',
            href: `/${businessKey}/feSpeed`
          },
          {
            title: '接口调用统计',
            href: `/${businessKey}/apiCall`
          }
        ]
      },
      {
        title: '后台异常上报',
        href: `/${businessKey}/endAbnormal`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/BlackstoneMonitor.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/BlackstoneMonitor-hover.svg'
        ]
      }
    ]
  },
  {
    title: '系统管理',
    menus: [
      {
        title: '审计管理',
        href: `/${businessKey}/audit`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/Config.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/monitor/css/img/aside-icon/Config-hover.svg'
        ]
      },
      {
        title: '权限申请',
        href: `/${businessKey}/auth`,
        icon: [
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/helm.svg',
          '//imgcache.qq.com/open_proj/proj_qcloud_v2/mc_2014/docker/css/img/side-icon/helm-hover.svg'
        ]
      }
    ]
  }
];
