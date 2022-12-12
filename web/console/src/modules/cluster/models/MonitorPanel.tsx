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
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { TransformField } from '@tencent/tchart';

interface FieldType {
  expr: string;
  alias: string;
  unit?: string;
  valueLabels?: (number, from) => string;
  units?: string[];
  thousands?: number;
}

function valueLabels1024(number) {
  return TransformField(number, 1024, 3, ['', 'K', 'M', 'G', 'T', 'P']);
}

interface ColumnType {
  key: string;
  name: string;
  render?: any;
}

export interface ChartPanelProps {
  tables: {
    table: string;
    fields: Array<FieldType>;
    conditions?: Array<Array<string>>;
    groupBy?: Array<{ value: string }>;
    columns?: Array<ColumnType>;
  }[];
  instance?: {
    // 类table中column和data
    columns: Array<ColumnType>;
    list: Array<any>;
  };
  groupBy?: Array<{ value: string }>;
}

export interface MonitorPanelProps extends ChartPanelProps {
  title: string;
  subTitle?: string;
  headerExtraDOM?;
}

export const networkFields = table => [
  {
    expr: `mean(k8s_${table}_network_receive_bytes_bw)`,
    alias: t('网络入带宽'),
    unit: 'Bps',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_network_transmit_bytes_bw)`,
    alias: t('网络出带宽'),
    unit: 'Bps',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_network_receive_bytes)`,
    alias: t('网络入流量'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_network_transmit_bytes)`,
    alias: t('网络出流量'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_network_receive_packets)`,
    alias: t('网络入包量'),
    unit: t('个/s')
  },
  {
    expr: `mean(k8s_${table}_network_transmit_packets)`,
    alias: t('网络出包量'),
    unit: t('个/s')
  }
];

const cpuMemFields = table => [
  {
    expr: `mean(k8s_${table}_cpu_core_used)`,
    alias: t('CPU使用量'),
    unit: t('核')
  },
  {
    expr: `mean(k8s_${table}_rate_cpu_core_used_node)`,
    alias: t('CPU利用率(占主机)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_cpu_core_used_request)`,
    alias: t('CPU利用率(占Request)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_cpu_core_used_limit)`,
    alias: t('CPU利用率(占Limit)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_mem_usage_bytes)`,
    alias: t('内存使用量'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_mem_no_cache_bytes)`,
    alias: t('内存使用量(不包含cache)'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_rate_mem_usage_node)`,
    alias: t('内存利用率(占主机)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_mem_no_cache_node)`,
    alias: t('内存利用率(占主机,不包含cache)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_mem_usage_request)`,
    alias: t('内存利用率(占Request)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_mem_no_cache_request)`,
    alias: t('内存利用率(占Request,不包含cache)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_mem_usage_limit)`,
    alias: t('内存利用率(占Limit)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_mem_no_cache_limit)`,
    alias: t('内存利用率(占Limit,不含cache)'),
    unit: '%'
  }
];

const gpuMemFields = table => [
  {
    expr: `mean(k8s_${table}_gpu_used)`,
    alias: t('GPU使用量'),
    unit: t('卡'),
    valueTransform: value => {
      if (value) return value / 100;

      return value;
    }
  },
  {
    expr: `mean(k8s_${table}_rate_gpu_used_node)`,
    alias: t('GPU使用率(占节点)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_gpu_used_request)`,
    alias: t('GPU使用率(占Request)'),
    unit: '%',
    valueTransform: value => {
      if (value) return value * 100;

      return value;
    }
  },
  {
    expr: `mean(k8s_${table}_gpu_memory_used)`,
    alias: t('GPU内存使用量'),
    unit: 'MiB'
    // thousands: 1024, valueLabels:valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_rate_gpu_memory_used_node)`,
    alias: t('GPU内存使用率(占节点)'),
    unit: '%'
  },
  {
    expr: `mean(k8s_${table}_rate_gpu_memory_used_request)`,
    alias: t('GPU内存利用率(占Request)'),
    unit: '%'
  }
];

const fsFields = table => [
  {
    expr: `mean(k8s_${table}_fs_write_bytes)`,
    alias: t('硬盘写流量'),
    unit: 'B/s',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_fs_read_bytes)`,
    alias: t('硬盘读流量'),
    unit: 'B/s',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_${table}_fs_read_times)`,
    alias: t('硬盘读 IOPS'),
    unit: '次/s'
  },
  {
    expr: `mean(k8s_${table}_fs_write_times)`,
    alias: t('硬盘写 IOPS'),
    unit: '次/s'
  }
];

const statusReadyField = table => ({
  expr: `max(k8s_${table}_status_ready)`,
  alias: t('异常状态'),
  // chartType: 'area',
  // colors: ['#006eff'],
  valueTransform: v => (v !== null ? +!v : v),
  valueLabels: (v, from) => {
    return from === 'yAxis'
      ? v
      : {
          0: `正常`,
          1: `<span class="text-danger">异常</span>`
        }[v];
  }
});

export const resourceMonitorFields = [
  {
    expr: 'sum(k8s_workload_pod_restart_total)',
    alias: t('Pod重启次数'),
    unit: t('次')
  },
  {
    expr: 'mean(k8s_workload_cpu_core_used)',
    alias: t('CPU使用量'),
    unit: t('核')
  },
  {
    expr: 'mean(k8s_workload_rate_cpu_core_used_cluster)',
    alias: t('CPU利用率(占集群)'),
    unit: '%'
  },
  {
    expr: 'mean(k8s_workload_mem_usage_bytes)',
    alias: t('内存使用量'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: 'mean(k8s_workload_rate_mem_usage_bytes_cluster)',
    alias: t('内存利用率(占集群)'),
    unit: '%'
  },

  {
    expr: 'mean(k8s_workload_gpu_used)',
    alias: t('GPU使用量'),
    unit: t('卡')
  },

  {
    expr: 'mean(k8s_workload_rate_gpu_used_cluster)',
    alias: t('GPU利用率(占集群)'),
    unit: '%'
  },

  {
    expr: 'mean(k8s_workload_gpu_memory_used)',
    alias: t('GPU内存使用量'),
    unit: 'MiB'
    // thousands: 1024, valueLabels:valueLabels1024
  },

  {
    expr: 'mean(k8s_workload_rate_gpu_memory_used_cluster)',
    alias: t('GPU内存利用率(占集群)'),
    unit: '%'
  },
  ...fsFields('workload'),
  ...networkFields('workload')
];

export const nodeMonitorFields = [
  {
    expr: 'sum(k8s_node_pod_restart_total)',
    alias: t('Pod重启次数'),
    unit: t('次')
  },
  {
    expr: 'max(k8s_node_pod_num)',
    alias: t('Pod数量'),
    unit: t('个')
  },
  statusReadyField('node'),
  {
    expr: 'mean(k8s_node_cpu_usage)',
    alias: t('CPU利用率'),
    unit: '%'
  },
  {
    expr: 'mean(k8s_node_mem_usage_no_cache)',
    alias: t('内存利用率'),
    unit: '%'
  },
  {
    expr: 'mean(k8s_node_gpu_usage)',
    alias: t('GPU利用率'),
    unit: '%'
  },
  {
    expr: 'mean(k8s_node_gpu_memory_usage)',
    alias: t('GPU内存利用率'),
    unit: '%'
  },
  {
    expr: 'mean(k8s_node_disk_space_rate)',
    alias: t('硬盘利用率'),
    unit: '%'
  },
  ...fsFields('node'),
  {
    expr: `mean(k8s_node_network_receive_bytes_bw)`,
    alias: t('网络入带宽'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: `mean(k8s_node_network_transmit_bytes_bw)`,
    alias: t('网络出带宽'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  }
];

export const podMonitorFields = [
  statusReadyField('pod'),
  ...cpuMemFields('pod'),
  ...gpuMemFields('pod'),
  ...fsFields('pod'),
  ...networkFields('pod')
];
export const containerMonitorFields = [
  ...cpuMemFields('container'),
  ...gpuMemFields('container'),
  ...fsFields('container')
];

export const projectFields = [
  { expr: 'mean(project_capacity_cpu)', alias: t('业务的cpu总量'), unit: t('核') },
  {
    expr: 'mean(project_capacity_memory)',
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务的内存总量')
  },
  { expr: 'mean(project_capacity_gpu)', alias: t('业务的gpu总量'), unit: t('卡') },
  { expr: 'mean(project_capacity_gpu_memory)', unit: '块', alias: t('业务的gpu memroy总量') },
  {
    expr: 'mean(project_capacity_cluster_cpu)',
    alias: t('业务在所有集群的cpu分配总量之和(分配给namespace)'),
    unit: t('核')
  },
  {
    expr: 'mean(project_capacity_cluster_memory)',
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在所有集群的内存分配总量之和(分配给namespace)')
  },
  {
    expr: 'mean(project_capacity_cluster_gpu)',
    alias: t('业务在所有集群的gpu分配总量之和(分配给namespace)'),
    unit: t('卡')
  },
  {
    expr: 'mean(project_capacity_cluster_gpu_memory)',
    unit: '块',
    // thousands: 1024, valueLabels:valueLabels1024,
    alias: t('业务在所有集群的gpu memory分配总量之和(分配给namespace)')
  },
  { expr: 'mean(project_allocated_cpu)', alias: t('业务在所有集群的cpu已分配给业务的总量(分配给pod)'), unit: t('核') },
  {
    expr: 'mean(project_allocated_memory)',
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在所有集群的内存已分配给业务的总量(分配给pod)')
  },
  { expr: 'mean(project_allocated_gpu)', alias: t('业务在所有集群的gpu已分配给业务的总量(分配给pod)'), unit: t('卡') },
  {
    expr: 'mean(project_allocated_gpu_memory)',
    unit: t('块'),
    // thousands: 1024, valueLabels:valueLabels1024,
    alias: t('业务在所有集群的gpu memory已分配给业务的总量(分配给pod)')
  },
  {
    expr: 'mean(project_cluster_capacity_gpu)',
    alias: t('业务在各个集群的gpu总量之和(分配给namespace)'),
    unit: t('卡')
  },
  {
    expr: 'mean(project_cluster_capacity_gpu_memory)',
    unit: t('块'),
    // thousands: 1024, valueLabels:valueLabels1024,
    alias: t('业务在各个集群的gpu memory总量之和(分配给namespace)')
  },
  { expr: 'mean(project_cluster_allocated_gpu)', alias: t('业务在各个集群的gpu分配量之和(分配给pod)'), unit: t('卡') },
  {
    expr: 'mean(project_cluster_allocated_gpu_memory)',
    unit: t('块'),
    // thousands: 1024, valueLabels:valueLabels1024,
    alias: t('业务在各个集群的gpu memory分配量之和(分配给pod)')
  },
  {
    expr: 'mean(project_cluster_cpu_core_used)',
    alias: t('业务在各个集群的cpu实际使用量'),
    unit: t('核')
  },
  {
    expr: 'mean(project_cluster_memory_usage_bytes)',
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在各个集群的内存实际使用')
  },
  {
    expr: 'mean(project_cluster_gpu_used)',
    alias: t('业务在各个集群的gpu实际使用量'),
    unit: t('卡')
  },
  {
    expr: 'mean(project_cluster_gpu_memory_used)',
    unit: '块',
    // thousands: 1024, valueLabels:valueLabels1024,
    alias: t('业务在各个集群的gpu memory实际使用量')
  },
  {
    expr: 'mean(project_cluster_network_receive_bytes_bw)',
    unit: 'Bps',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在各个集群的网络接收带宽实际使用量')
  },
  {
    expr: 'mean(project_cluster_network_transmit_bytes_bw)',
    unit: 'Bps',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在各个集群的网络发送带宽实际使用量')
  },
  {
    expr: 'mean(project_cluster_network_receive_bytes)',
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在各个集群的网络接收速率实际使用量')
  },
  {
    expr: 'mean(project_cluster_network_transmit_bytes)',
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在各个集群的网络发送速率实际使用量')
  },
  {
    expr: 'mean(project_cluster_fs_read_bytes)',
    unit: 'B/s',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在各个集群的磁盘读速率实际使用量')
  },
  {
    expr: 'mean(project_cluster_fs_write_bytes)',
    unit: 'B/s',
    thousands: 1024,
    valueLabels: valueLabels1024,
    alias: t('业务在各个集群的磁盘写速率实际使用量')
  }
].filter(m => !m.expr.includes('gpu'));

export function getClusterTables(clusterId) {
  return [
    {
      table: 'apiserver_request_latencies_summary',
      fields: [
        {
          expr: `mean(k8s_component_apiserver_request_latency)`,
          alias: 'Apiserver时延',
          unit: 'ms',
          valueLabels: v => TransformField(v, 1000, 3, ['µ', 'm', '', 'k', 'M', 'G', 'T', 'P']),
          thousands: 1000
        }
      ],
      groupBy: [{ value: 'node' }],
      columns: [{ key: 'node', name: 'node' }],
      conditions: [['tke_cluster_instance_id', '=', clusterId]]
    },
    {
      table: 'k8s_component_scheduler',
      fields: [
        {
          expr: `mean(k8s_component_scheduler_scheduling_latency)`,
          alias: 'Scheduler时延',
          unit: 's',
          valueLabels: v => TransformField(v, 1000, 3, ['µ', 'm', '', 'k', 'M', 'G', 'T', 'P']),
          // units: ['µ', 'm', '', 'k', 'M', 'G', 'T', 'P'],
          thousands: 1000
        }
      ],
      groupBy: [{ value: 'node' }],
      conditions: [['tke_cluster_instance_id', '=', clusterId]]
    },
    {
      table: 'etcd_cluster',
      fields: [
        {
          expr: `mean(etcd_debugging_mvcc_db_total_size_in_bytes)`,
          alias: 'Etcd 存储量',
          unit: 'B',
          thousands: 1024,
          valueLabels: valueLabels1024
        }
      ],
      groupBy: [],
      conditions: [['tke_cluster_instance_id', '=', clusterId]]
    },
    {
      table: 'k8s_cluster',
      fields: [
        {
          expr: `mean(k8s_cluster_rate_cpu_core_used_cluster)`,
          alias: 'CPU利用率',
          unit: '%'
        },
        {
          expr: `mean(k8s_cluster_rate_cpu_core_request_cluster)`,
          alias: 'CPU分配率(Request)',
          unit: '%'
        },
        {
          expr: `mean(k8s_cluster_rate_mem_usage_bytes_cluster)`,
          alias: '内存利用率',
          unit: '%'
        },
        {
          expr: `mean(k8s_cluster_rate_mem_request_bytes_cluster)`,
          alias: '内存分配率(Request)',
          unit: '%'
        },
        {
          expr: `mean(k8s_cluster_rate_gpu_used_cluster)`,
          alias: 'GPU利用率',
          unit: '%'
        },
        {
          expr: `mean(k8s_cluster_rate_gpu_request_cluster)`,
          alias: 'GPU分配率(Request)',
          unit: '%'
        },
        {
          expr: `mean(k8s_cluster_rate_gpu_memory_used_cluster)`,
          alias: 'GPU内存利用率',
          unit: '%'
        },
        {
          expr: `mean(k8s_cluster_rate_gpu_memory_request_cluster)`,
          alias: 'GPU内存分配率(Request)',
          unit: '%'
        },
        ...networkFields('cluster')
      ],
      conditions: [['tke_cluster_instance_id', '=', clusterId]]
    }
  ];
}

export const meshMonitorFields = [
  { expr: 'sum(count)', alias: t('请求数') },
  { expr: 'mean(duration)', alias: t('请求耗时'), unit: 'ms' },
  {
    expr: 'mean(request_total_size)',
    alias: t('TCP请求字节'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: 'mean(response_total_size)',
    alias: t('TCP接收字节'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  { expr: 'mean(request_size)', alias: t('请求大小'), unit: 'B', thousands: 1024, valueLabels: valueLabels1024 },
  { expr: 'mean(response_size)', alias: t('接收大小'), unit: 'B', thousands: 1024, valueLabels: valueLabels1024 }
];

export const meshTcpMonitorFields = [
  { expr: 'sum(count)', alias: t('请求数') },
  {
    expr: 'mean(connection_sent_bytes_total)',
    alias: t('发送字节总数'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  },
  {
    expr: 'mean(connection_received_bytes_total)',
    alias: t('接收字节总数'),
    unit: 'B',
    thousands: 1024,
    valueLabels: valueLabels1024
  }
];
