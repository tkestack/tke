import { t, Trans } from '@tencent/tea-app/lib/i18n';

export const AlarmObjectsType = {
  node: [
    { value: 'all', text: t('全部选择'), tip: t('（包括后续新增加节点）') }
    // {
    //   value: 'k8sLabel',
    //   text: '按 Kubernetes Label 选择',
    //   tip: '（对满足所有 Label 条件的节点生效，包括后续新增加节点）'
    // }
  ],
  pod: [
    { value: 'part', text: t('按工作负载选择'), tip: t('（包括后续新增加Pod）') },
    /// #if tke
    { value: 'all', text: t('全部选择'), tip: t('（包括后续新增加Pod）') }
    /// #endif

    // {
    //   value: 'k8sLabel',
    //   text: '按 Kubernetes Label 选择',
    //   tip: '（对满足所有 Label 条件的Pod生效，包括后续新增加Pod）'
    // }
  ]
};

export const AlarmNotifyWay = [
  {
    value: 'SMS'
  },
  {
    value: 'EMAIL'
  },
  {
    value: 'WECHAT'
  },
  {
    value: 'CALL'
  }
];

export const AlarmPolicyType = [
  {
    text: t('集群'),
    value: 'cluster'
  },
  {
    text: t('节点'),
    value: 'node'
  },
  {
    text: 'pod',
    value: 'pod'
  }
];
//统计周期
export const AlarmPolicyMetricsStatisticsPeriod = [
  { value: '1' },
  { value: '2' },
  { value: '3' },
  { value: '4' },
  { value: '5' }
];

//指标操作
export const AlarmPolicyMetricsEvaluatorType = [
  {
    text: '>',
    value: 'gt'
  },
  { text: '<', value: 'lt' }
];

//阈值
export const AlarmPolicyMetricsEvaluatorValue = [
  { text: 'True', value: 'true' },
  { text: 'False', value: 'false' }
];

//持续周期
export const AlarmPolicyMetricsContinuePeriod = [
  { value: '1' },
  { value: '2' },
  { value: '3' },
  { value: '4' },
  { value: '5' }
];

export const AlarmPolicyPhoneCircleTimes = [
  { value: '1' },
  { value: '2' },
  { value: '3' },
  { value: '4' },
  { value: '5' }
];

export const AlarmPolicyPhoneInterval = [
  { value: '1' },
  { value: '2' },
  { value: '3' },
  { value: '4' },
  { value: '5' },
  { value: '6' },
  { value: '7' },
  { value: '8' },
  { value: '9' },
  { value: '10' },
  { value: '11' },
  { value: '12' },
  { value: '13' },
  { value: '14' },
  { value: '15' }
];

export const workloadTypeList = [
  {
    value: 'deployment',
    label: 'Deployment'
  },
  {
    value: 'daemonset',
    label: 'DaemonSet'
  },
  {
    value: 'statefulset',
    label: 'StatefulSet'
  }
];

export const AlarmPolicyMetrics = {
  cluster: [
    {
      enable: true,
      measurement: 'k8s_cluster',
      statisticsPeriod: 1,
      metricName: 'k8s_cluster_rate_cpu_core_used_cluster',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('CPU利用率'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_cluster',
      statisticsPeriod: 1,
      metricName: 'k8s_cluster_rate_mem_usage_bytes_cluster',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('内存利用率'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_cluster',
      statisticsPeriod: 1,
      metricName: 'k8s_cluster_rate_cpu_core_request_cluster',
      evaluatorType: 'gt',
      evaluatorValue: '95',
      metricDisplayName: t('CPU分配率'),
      continuePeriod: 5,
      type: 'percent',
      tip: t('所有容器Request之和/集群总可分配资源'),
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_cluster',
      statisticsPeriod: 1,
      metricName: 'k8s_cluster_rate_mem_request_bytes_cluster',
      evaluatorType: 'gt',
      evaluatorValue: '95',
      metricDisplayName: t('内存分配率'),
      continuePeriod: 5,
      type: 'percent',
      tip: t('所有容器Request之和/集群总可分配资源'),
      // metricType: 'metric',
      unit: '%'
    }
  ],
  independentClusetr: [
    {
      enable: true,
      measurement: 'k8s_cluster',
      statisticsPeriod: 1,
      metricName: 'k8s_cluster_rate_cpu_core_used_cluster',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('CPU利用率'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_cluster',
      statisticsPeriod: 1,
      metricName: 'k8s_cluster_rate_mem_usage_bytes_cluster',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('内存利用率'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_cluster',
      statisticsPeriod: 1,
      metricName: 'k8s_cluster_rate_cpu_core_request_cluster',
      evaluatorType: 'gt',
      evaluatorValue: '95',
      metricDisplayName: t('CPU分配率'),
      continuePeriod: 5,
      type: 'percent',
      tip: t('所有容器Request之和/集群总可分配资源'),
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_cluster',
      statisticsPeriod: 1,
      metricName: 'k8s_cluster_rate_mem_request_bytes_cluster',
      evaluatorType: 'gt',
      evaluatorValue: '95',
      metricDisplayName: t('内存分配率'),
      continuePeriod: 5,
      type: 'percent',
      tip: t('所有容器Request之和/集群总可分配资源'),
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_component',
      statisticsPeriod: 1,
      metricName: 'k8s_component_apiserver_ready',
      evaluatorType: 'eq',
      evaluatorValue: AlarmPolicyMetricsEvaluatorValue[1].value,
      metricDisplayName: t('API Server正常'),
      continuePeriod: 5,
      type: 'boolean',
      tip: '',
      // metricType: 'event',
      unit: ''
    },
    {
      enable: true,
      measurement: 'k8s_component',
      statisticsPeriod: 1,
      metricName: 'k8s_component_etcd_ready',
      evaluatorType: 'eq',
      evaluatorValue: AlarmPolicyMetricsEvaluatorValue[1].value,
      metricDisplayName: t('Etcd正常'),
      continuePeriod: 5,
      type: 'boolean',
      tip: '',
      // metricType: 'event',
      unit: ''
    },
    {
      enable: true,
      measurement: 'k8s_component',
      statisticsPeriod: 1,
      metricName: 'k8s_component_scheduler_ready',
      evaluatorType: 'eq',
      evaluatorValue: AlarmPolicyMetricsEvaluatorValue[1].value,
      metricDisplayName: t('Scheduler正常'),
      continuePeriod: 5,
      type: 'boolean',
      tip: '',
      // metricType: 'event',
      unit: ''
    },
    {
      enable: true,
      measurement: 'k8s_component',
      statisticsPeriod: 1,
      metricName: 'k8s_component_controller_manager_ready',
      evaluatorType: 'eq',
      evaluatorValue: AlarmPolicyMetricsEvaluatorValue[1].value,
      metricDisplayName: t('Controll Manager正常'),
      continuePeriod: 5,
      type: 'boolean',
      tip: '',
      // metricType: 'event',
      unit: ''
    }
  ],
  node: [
    {
      enable: true,
      measurement: 'k8s_node',
      statisticsPeriod: 1,
      metricName: 'k8s_node_cpu_usage',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('CPU利用率'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_node',
      statisticsPeriod: 1,
      metricName: 'k8s_node_mem_usage',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('内存利用率'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_node',
      statisticsPeriod: 1,
      metricName: 'k8s_node_pod_restart_total',
      evaluatorType: 'gt',
      evaluatorValue: '1',
      metricDisplayName: t('节点上Pod重启次数'),
      continuePeriod: 5,
      type: 'times',
      tip: t('该统计为按工作负载或Label条件聚合后的数值'),
      // metricType: 'metric',
      unit: t('次')
    },
    {
      enable: true,
      measurement: 'k8s_node',
      statisticsPeriod: 1,
      metricName: 'k8s_node_status_ready',
      evaluatorType: 'eq',
      evaluatorValue: AlarmPolicyMetricsEvaluatorValue[1].value,
      metricDisplayName: 'Node Ready',
      continuePeriod: 5,
      type: 'boolean',
      tip: '',
      // metricType: 'event',
      unit: ''
    }
  ],
  pod: [
    {
      enable: true,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_rate_cpu_core_used_node',
      evaluatorType: 'gt',
      evaluatorValue: '80',
      metricDisplayName: t('CPU利用率（占节点）'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_rate_mem_usage_node',
      evaluatorType: 'gt',
      evaluatorValue: '80',
      metricDisplayName: t('内存利用率（占节点）'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_rate_mem_no_cache_node',
      evaluatorType: 'gt',
      evaluatorValue: '80',
      metricDisplayName: t('实际内存利用率（占节点）'),
      continuePeriod: 5,
      type: 'percent',
      tip: t('不包括缓存'),
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_rate_cpu_core_used_limit',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('CPU利用率（占limit）'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_rate_mem_usage_limit',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('内存利用率（占limit）'),
      continuePeriod: 5,
      type: 'percent',
      tip: '',
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_rate_mem_no_cache_limit',
      evaluatorType: 'gt',
      evaluatorValue: '90',
      metricDisplayName: t('实际内存利用率（占Limit）'),
      continuePeriod: 5,
      type: 'percent',
      tip: t('不包括缓存'),
      // metricType: 'metric',
      unit: '%'
    },
    {
      enable: true,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_restart_total',
      evaluatorType: 'gt',
      evaluatorValue: '1',
      metricDisplayName: t('Pod重启次数'),
      continuePeriod: 5,
      type: 'times',
      tip: t('该统计为按工作负载或Label条件聚合后的数值'),
      // metricType: 'metric',
      unit: t('次')
    },
    {
      enable: true,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_status_ready',
      evaluatorType: 'eq',
      evaluatorValue: AlarmPolicyMetricsEvaluatorValue[1].value,
      metricDisplayName: 'Pod Ready',
      continuePeriod: 5,
      type: 'boolean',
      tip: '',
      // metricType: 'event',
      unit: ''
    },
    {
      enable: false,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_cpu_core_used',
      evaluatorType: 'gt',
      evaluatorValue: '',
      metricDisplayName: t('CPU使用量'),
      continuePeriod: 5,
      type: 'sumCpu',
      tip: '',
      // metricType: 'metric',
      unit: t('核')
    },
    {
      enable: false,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_mem_usage_bytes',
      evaluatorType: 'gt',
      evaluatorValue: '',
      metricDisplayName: t('内存使用量'),
      continuePeriod: 5,
      type: 'sumMem',
      tip: '',
      // metricType: 'metric',
      unit: 'MB'
    },
    {
      enable: false,
      measurement: 'k8s_pod',
      statisticsPeriod: 1,
      metricName: 'k8s_pod_mem_no_cache_bytes',
      evaluatorType: 'gt',
      evaluatorValue: '',
      metricDisplayName: t('实际内存使用量'),
      continuePeriod: 5,
      type: 'sumMem',
      tip: t('不包括缓存'),
      // metricType: 'metric',
      unit: 'MB'
    }
  ]
};

export const MetricNameMap = {};

for (let key in AlarmPolicyMetrics) {
  for (let metric of AlarmPolicyMetrics[key]) {
    MetricNameMap[metric.metricName] = metric.metricDisplayName;
  }
}
