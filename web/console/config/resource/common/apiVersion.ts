import {
  apiServerVersion,
  businessServerVersion,
  notifyServerVersion,
  basicServerVersion,
  authServerVersion,
  monitorServerVersion
} from '../../apiServerVersion';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ConsoleModuleEnum } from '../../../Wrapper';

export type ApiVersionKeyName = keyof ApiVersion;

export interface ApiVersion {
  deployment?: ResourceApiInfo;
  statefulset?: ResourceApiInfo;
  daemonset?: ResourceApiInfo;
  job?: ResourceApiInfo;
  cronjob?: ResourceApiInfo;
  tapp?: ResourceApiInfo;
  pods?: ResourceApiInfo;
  rc?: ResourceApiInfo;
  rs?: ResourceApiInfo;
  svc?: ResourceApiInfo;
  ingress?: ResourceApiInfo;
  np?: ResourceApiInfo;
  configmap?: ResourceApiInfo;
  secret?: ResourceApiInfo;
  pv?: ResourceApiInfo;
  pvc?: ResourceApiInfo;
  sc?: ResourceApiInfo;
  hpa?: ResourceApiInfo;
  cronhpa?: ResourceApiInfo;
  event?: ResourceApiInfo;
  node?: ResourceApiInfo;
  masteretcd?: ResourceApiInfo;
  cluster?: ResourceApiInfo;
  pe?: ResourceApiInfo;
  ns?: ResourceApiInfo;
  localidentity?: ResourceApiInfo;
  policy?: ResourceApiInfo;
  user?: ResourceApiInfo;
  role?: ResourceApiInfo;
  localgroup?: ResourceApiInfo;
  group?: ResourceApiInfo;
  category?: ResourceApiInfo;
  machines?: ResourceApiInfo;
  helm?: ResourceApiInfo;
  logcs?: ResourceApiInfo;
  clustercredential?: ResourceApiInfo;

  lbcf?: ResourceApiInfo;
  lbcf_bg?: ResourceApiInfo;
  lbcf_br?: ResourceApiInfo;
  /** ============= 以下是addon的相关配置 ============== */
  addon?: ResourceApiInfo;
  addon_helm?: ResourceApiInfo;
  addon_gpumanager?: ResourceApiInfo;
  addon_logcollector?: ResourceApiInfo;
  addon_persistentevent?: ResourceApiInfo;
  addon_tappcontroller?: ResourceApiInfo;
  addon_csioperator?: ResourceApiInfo;
  addon_lbcf?: ResourceApiInfo;
  addon_cronhpa?: ResourceApiInfo;
  addon_coredns?: ResourceApiInfo;
  addon_galaxy?: ResourceApiInfo;
  addon_prometheus?: ResourceApiInfo;
  addon_volumedecorator?: ResourceApiInfo;
  addon_ipam?: ResourceApiInfo;
  /** ============= 以上是addon的相关配置 ============== */

  /** ============= 以上是 用户相关的配置信息 的相关配置 ============== */
  module?: ResourceApiInfo;
  info?: ResourceApiInfo;
  redirect?: ResourceApiInfo;
  renew?: ResourceApiInfo;
  logout?: ResourceApiInfo;
  /** ============= 以上是 用户相关的配置信息 的相关配置 ============== */

  /** ============= 以上是 业务配置信息 的相关配置 ============== */
  projects?: ResourceApiInfo;
  platforms?: ResourceApiInfo;
  portal?: ResourceApiInfo;
  namespaces?: ResourceApiInfo;
  /** ============= 以上是 业务配置信息 的相关配置 ============== */

  /** 告警配置 */
  prometheus?: ResourceApiInfo;
  alarmPolicy?: ResourceApiInfo;
  channel?: ResourceApiInfo;
  template?: ResourceApiInfo;
  message?: ResourceApiInfo;
  receiver?: ResourceApiInfo;
  receiverGroup?: ResourceApiInfo;

  /** 服务网格资源 */
  serviceForMesh?: ResourceApiInfo;
  gateway?: ResourceApiInfo;
  virtualservice?: ResourceApiInfo;
  destinationrule?: ResourceApiInfo;
  serviceentry?: ResourceApiInfo;
  controlPlane?: ResourceApiInfo;

  /** 组织资源 */
  apiKey?: ResourceApiInfo;
}

interface ResourceApiInfo {
  group: string;
  version: string;
  basicEntry: string;
  headTitle?: string;
  watchModule?: ConsoleModuleEnum;
}

/** 这里的apiVersion的配置，只列出与1.8不同的，其余都保持和1.8相同的配置 */
const k8sApiVersionFor17: ApiVersion = {
  deployment: {
    group: 'apps',
    version: 'v1beta1',
    basicEntry: 'apis',
    headTitle: 'Deployment'
  },
  statefulset: {
    group: 'apps',
    version: 'v1beta1',
    basicEntry: 'apis',
    headTitle: 'StatefulSet'
  },
  daemonset: {
    group: 'extensions',
    version: 'v1beta1',
    basicEntry: 'apis',
    headTitle: 'DaemonSet'
  },
  rs: {
    group: 'extensions',
    version: 'v1beta1',
    basicEntry: 'apis',
    headTitle: 'ReplicaSet'
  },
  cronjob: {
    group: 'batch',
    version: 'v2alpha1',
    basicEntry: 'apis',
    headTitle: 'CronJob'
  },
  hpa: {
    group: 'autoscaling',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'HorizontalPodAutoscaler'
  },
  gateway: {
    group: 'networking.istio.io',
    version: 'v1alpha3',
    basicEntry: 'apis',
    headTitle: 'Gateway'
  },
  virtualservice: {
    group: 'networking.istio.io',
    version: 'v1alpha3',
    basicEntry: 'apis',
    headTitle: 'Virtual Service'
  },
  destinationrule: {
    group: 'networking.istio.io',
    version: 'v1alpha3',
    basicEntry: 'apis',
    headTitle: 'Destination Rule'
  },
  serviceentry: {
    group: 'networking.istio.io',
    version: 'v1alpha3',
    basicEntry: 'apis',
    headTitle: 'None'
  }
};

/** ================== start 1.14 的apiversion配置 ======================== */
const k8sApiVersionFor114: ApiVersion = {
  deployment: {
    group: 'apps',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'Deployment'
  },
  statefulset: {
    group: 'apps',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'StatefulSet'
  },
  daemonset: {
    group: 'apps',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'DaemonSet'
  },
  rs: {
    group: 'apps',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'ReplicaSet'
  }
};
/** ================== start 1.14 的apiversion配置 ======================== */

/** 以1.8的为基准，后续有新增再继续更改 */
const k8sApiVersionFor18: ApiVersion = {
  deployment: {
    group: 'apps',
    version: 'v1beta2',
    basicEntry: 'apis',
    headTitle: 'Deployment'
  },
  statefulset: {
    group: 'apps',
    version: 'v1beta2',
    basicEntry: 'apis',
    headTitle: 'StatefulSet'
  },
  daemonset: {
    group: 'apps',
    version: 'v1beta2',
    basicEntry: 'apis',
    headTitle: 'DaemonSet'
  },
  job: {
    group: 'batch',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'Job'
  },
  cronjob: {
    group: 'batch',
    version: 'v1beta1',
    basicEntry: 'apis',
    headTitle: 'CronJob'
  },
  tapp: {
    group: 'apps.tkestack.io',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'TApp'
  },
  pods: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'Pod'
  },
  np: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'Namespace'
  },
  rc: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'ReplicationController'
  },
  rs: {
    group: 'apps',
    version: 'v1beta2',
    basicEntry: 'apis',
    headTitle: 'ReplicaSet'
  },
  svc: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'Service'
  },
  ingress: {
    group: 'extensions',
    version: 'v1beta1',
    basicEntry: 'apis',
    headTitle: 'Ingress'
  },
  configmap: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'ConfigMap'
  },
  secret: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'Secret'
  },
  pv: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'PersistentVolume'
  },
  pvc: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'PersistentVolumeClaim'
  },
  sc: {
    group: 'storage.k8s.io',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'StorageClass'
  },
  hpa: {
    group: 'autoscaling',
    version: 'v2beta1',
    basicEntry: 'apis',
    headTitle: 'HorizontalPodAutoscaler'
  },
  cronhpa: {
    group: 'extensions.tkestack.io',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'CronHPA'
  },
  event: {
    group: '',
    version: 'v1',
    basicEntry: 'api'
  },
  node: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'Node'
  },
  masteretcd: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'Node'
  },
  cluster: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Cluster'
  },
  pe: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'PersistentEvent'
  },
  ns: {
    group: '',
    version: 'v1',
    basicEntry: 'api',
    headTitle: 'Namespace'
  },

  localidentity: {
    group: authServerVersion.group,
    version: authServerVersion.version,
    basicEntry: authServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Auth,
    headTitle: 'Localidentities'
  },
  policy: {
    group: authServerVersion.group,
    version: authServerVersion.version,
    basicEntry: authServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Auth,
    headTitle: 'Policies'
  },
  user: {
    group: authServerVersion.group,
    version: authServerVersion.version,
    basicEntry: authServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Auth,
    headTitle: 'Users'
  },
  role: {
    group: authServerVersion.group,
    version: authServerVersion.version,
    basicEntry: authServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Auth,
    headTitle: 'Roles'
  },
  localgroup: {
    group: authServerVersion.group,
    version: authServerVersion.version,
    basicEntry: authServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Auth,
    headTitle: 'Localgroups'
  },
  group: {
    group: authServerVersion.group,
    version: authServerVersion.version,
    basicEntry: authServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Auth,
    headTitle: 'Groups'
  },
  apiKey: {
    group: authServerVersion.group,
    version: authServerVersion.version,
    basicEntry: authServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Auth,
    headTitle: 'Apikeys'
  },
  category: {
    group: authServerVersion.group,
    version: authServerVersion.version,
    basicEntry: authServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Auth,
    headTitle: 'Categories'
  },
  machines: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Machine'
  },
  helm: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Helm'
  },

  logcs: {
    group: 'tke.cloud.tencent.com',
    version: 'v1',
    basicEntry: 'apis',
    headTitle: 'LogCollector'
  },

  clustercredential: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Clustercredential'
  },

  lbcf: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'LBCF'
  },
  lbcf_bg: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'BackendGroup'
  },
  lbcf_br: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'BackendRecord'
  },
  gateway: {
    group: 'networking.istio.io',
    version: 'v1alpha3',
    basicEntry: 'apis',
    headTitle: 'Gateway'
  },
  virtualservice: {
    group: 'networking.istio.io',
    version: 'v1alpha3',
    basicEntry: 'apis',
    headTitle: 'Virtual Service'
  },
  destinationrule: {
    group: 'networking.istio.io',
    version: 'v1alpha3',
    basicEntry: 'apis',
    headTitle: 'Destination Rule'
  },
  serviceentry: {
    group: 'networking.istio.io',
    version: 'v1alpha3',
    basicEntry: 'apis',
    headTitle: 'None'
  }
};

/** addon 的相关配置 */
const k8sApiVersionForAddon: ApiVersion = {
  addon: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Addon'
  },
  addon_helm: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Helm'
  },
  addon_gpumanager: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'GPUManager'
  },
  addon_logcollector: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'LogCollector'
  },
  addon_persistentevent: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'PersistentEvent'
  },
  addon_tappcontroller: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'TappController'
  },
  addon_csioperator: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'CSIOperator'
  },
  addon_lbcf: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'LBCF'
  },
  addon_cronhpa: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'CronHPA'
  },
  addon_coredns: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'CoreDNS'
  },
  addon_galaxy: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Galaxy'
  },
  addon_prometheus: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Prometheus'
  },
  addon_volumedecorator: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'VolumeDecorator'
  },
  addon_ipam: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'IPAM'
  }
};

/** 用户相关的配置信息 */
const k8sConsoleApiVersion: ApiVersion = {
  module: {
    group: basicServerVersion.group,
    version: basicServerVersion.version,
    basicEntry: basicServerVersion.basicUrl,
    headTitle: 'Module'
  },
  info: {
    group: basicServerVersion.group,
    version: basicServerVersion.version,
    basicEntry: basicServerVersion.basicUrl,
    headTitle: 'Info'
  },
  redirect: {
    group: basicServerVersion.group,
    version: basicServerVersion.version,
    basicEntry: basicServerVersion.basicUrl,
    headTitle: 'Redirect'
  },
  renew: {
    group: basicServerVersion.group,
    version: basicServerVersion.version,
    basicEntry: basicServerVersion.basicUrl,
    headTitle: 'Renew'
  },
  logout: {
    group: basicServerVersion.group,
    version: basicServerVersion.version,
    basicEntry: basicServerVersion.basicUrl,
    headTitle: 'Logout'
  }
};

/** 业务相关的配置信息 */
const k8sProjectApiVersion: ApiVersion = {
  projects: {
    group: businessServerVersion.group,
    version: businessServerVersion.version,
    basicEntry: businessServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Business,
    headTitle: 'Project'
  },
  namespaces: {
    group: businessServerVersion.group,
    version: businessServerVersion.version,
    basicEntry: businessServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Business,
    headTitle: 'Namespace'
  },
  portal: {
    group: businessServerVersion.group,
    version: businessServerVersion.version,
    basicEntry: businessServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Business,
    headTitle: 'Portal'
  },
  platforms: {
    group: businessServerVersion.group,
    version: businessServerVersion.version,
    basicEntry: businessServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Business,
    headTitle: 'Platform'
  }
};

const alarmPolicyApiVersion: ApiVersion = {
  alarmPolicy: {
    group: '',
    version: monitorServerVersion.version,
    basicEntry: monitorServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Monitor,
    headTitle: 'AlarmPolicy'
  },
  prometheus: {
    group: apiServerVersion.group,
    version: apiServerVersion.version,
    basicEntry: apiServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.PLATFORM,
    headTitle: 'Prometheus'
  }
};

const notifyApiVersion: ApiVersion = {
  channel: {
    group: notifyServerVersion.group,
    version: notifyServerVersion.version,
    basicEntry: notifyServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Notify,
    headTitle: t('通知渠道')
  },
  template: {
    group: notifyServerVersion.group,
    version: notifyServerVersion.version,
    basicEntry: notifyServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Notify,
    headTitle: t('通知模版')
  },
  message: {
    group: notifyServerVersion.group,
    version: notifyServerVersion.version,
    basicEntry: notifyServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Notify,

    headTitle: 'Message'
  },
  receiver: {
    group: notifyServerVersion.group,
    version: notifyServerVersion.version,
    basicEntry: notifyServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Notify,
    headTitle: t('接收人')
  },
  receiverGroup: {
    group: notifyServerVersion.group,
    version: notifyServerVersion.version,
    basicEntry: notifyServerVersion.basicUrl,
    watchModule: ConsoleModuleEnum.Notify,
    headTitle: t('接收组')
  }
};

interface FinalApiVersion {
  [props: string]: ApiVersion;
}

const basicApiVersion = Object.assign(
  {},
  k8sProjectApiVersion,
  k8sConsoleApiVersion,
  k8sApiVersionForAddon,
  k8sApiVersionFor18,
  alarmPolicyApiVersion,
  notifyApiVersion
);

/**
 * 这里配置的是k8s各版本的的k8s资源的 group 和 version 的最新版本
 * @pre 只需要保证每个k8s大版本当中的 group 和 version使用最新的即可，会向下兼容
 */
export const apiVersion: FinalApiVersion = {
  '1.8': basicApiVersion,
  '1.7': Object.assign({}, basicApiVersion, k8sApiVersionFor17),
  '1.14': Object.assign({}, basicApiVersion, k8sApiVersionFor114)
};
