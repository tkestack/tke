import { t, Trans } from '@tencent/tea-app/lib/i18n';

/** ========================= start FFRedux的相关配置 ======================== */
export const FFReduxActionName = {
  REGION: 'region',
  CLUSTER: 'cluster'
};
/** ========================= end FFRedux的相关配置 ======================== */

/**
 * helm开通状态
 */
export const ClusterHelmStatus = {
  //未开通
  NONE: 'none',
  //正在检查
  CHECKING: 'checking',
  //正在初始化
  INIT: 'initializing',
  //正在重新初始化
  REINIT: 'reinitalizing',
  //开通失败
  ERROR: 'error',
  //已开通
  RUNNING: 'running'
};

export const InstallingStatus = {
  INSTALLING: 0,
  ERROR: 1
};

export const InstallingStatusText = {
  [InstallingStatus.INSTALLING]: {
    text: t('创建中'),
    classname: 'text-success'
  },
  [InstallingStatus.ERROR]: {
    text: t('失败'),
    classname: 'text-danger'
  },
  '-': {
    text: '-',
    classname: 'text-restart'
  }
};

/**
 * helm状态
 */
export const helmStatus = {
  DEPLOYED: {
    text: t('正常'),
    classname: 'text-success'
  },
  DELETED: {
    text: t('已删除'),
    classname: 'text-restart'
  },
  DELETING: {
    text: t('正在删除'),
    classname: 'text-restart'
  },
  SUPERSEDED: {
    text: t('已废弃'),
    classname: 'text-restart'
  },
  FAILED: {
    text: t('异常'),
    classname: 'text-danger'
  },
  '-': {
    text: '-',
    classname: 'text-restart'
  }
};

export const HelmResource = {
  TencentHub: 'TencentHub',
  Helm: 'Helm',
  Other: 'Other'
};
export const helmResourceList = [
  // {
  //   value: HelmResource.TencentHub,
  //   name: 'TencentHub'
  // },
  // {
  //   value: HelmResource.Helm,
  //   name: "Helm官方"
  // },
  {
    value: HelmResource.Other,
    name: t('其他')
  }
];

export const TencentHubType = {
  Public: 'tencenthub',
  Private: 'private'
};

export const tencentHubTypeList = [
  {
    value: TencentHubType.Public,
    name: t('公有')
  },
  {
    value: TencentHubType.Private,
    name: t('私有')
  }
];

export const OtherType = {
  Public: 'public',
  Private: 'private'
};

export const otherTypeList = [
  {
    value: OtherType.Public,
    name: t('公有')
  },
  {
    value: OtherType.Private,
    name: t('私有')
  }
];

export const ResourceUrl = {
  Deployment: (rid, clusterId) => {
    return `/tke/cluster/sub/list/resource/deployment?rid=${rid}&clusterId=${clusterId}`;
  },

  StatefulSet: (rid, clusterId) => {
    return `/tke/cluster/sub/list/resource/statefulset?rid=${rid}&clusterId=${clusterId}`;
  },

  DaementSet: (rid, clusterId) => {
    return `/tke/cluster/sub/list/resource/daemonset?rid=${rid}&clusterId=${clusterId}`;
  },
  Job: (rid, clusterId) => {
    return `/tke/cluster/sub/list/resource/job?rid=${rid}&clusterId=${clusterId}`;
  },
  CronJob: (rid, clusterId) => {
    return `/tke/cluster/sub/list/resource/cronjob?rid=${rid}&clusterId=${clusterId}`;
  },
  Service: (rid, clusterId) => {
    return `/tke/cluster/sub/list/service/svc?rid=${rid}&clusterId=${clusterId}`;
  },
  Ingress: (rid, clusterId) => {
    return `/tke/cluster/sub/list/service/ingress?rid=${rid}&clusterId=${clusterId}`;
  },
  ConfigMap: (rid, clusterId) => {
    return `/tke/cluster/sub/list/config/configmap?rid=${rid}&clusterId=${clusterId}`;
  },
  Secret: (rid, clusterId) => {
    return `/tke/cluster/sub/list/config/secret?rid=${rid}&clusterId=${clusterId}`;
  },
  PersistentVolume: (rid, clusterId) => {
    return `/tke/cluster/sub/list/storage/pv?rid=${rid}&clusterId=${clusterId}`;
  },
  PersistentVolumeClaim: (rid, clusterId) => {
    return `/tke/cluster/sub/list/storage/pvc?rid=${rid}&clusterId=${clusterId}`;
  },
  Storagress: (rid, clusterId) => {
    return `/tke/cluster/sub/list/storage/sc?rid=${rid}&clusterId=${clusterId}`;
  }
};
