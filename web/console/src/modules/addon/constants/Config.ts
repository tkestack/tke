import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { ApiVersionKeyName } from '../../../../config/resource/common';
/** ========================= start FFRedux的相关配置 ======================== */
export const FFReduxActionName = {
  REGION: 'region',
  CLUSTER: 'cluster',
  OPENADDON: 'openAddon',
  ADDON: 'addon',
  LOGSET: 'logset',
  TOPIC: 'topic'
};
/** ========================= end FFRedux的相关配置 ======================== */

/** ========================= start Addon的相关配置项 ======================== */
export enum ClusterTypeEnum {
  Imported = 'Imported',
  Baremetal = 'Baremetal'
}

/** Addon组件类型的名称映射 */
export const AddonTypeMap = {
  Enhance: '增强组件',
  Basic: '基础组件'
};

/** Addon组件的类型 */
export enum AddonTypeEnum {
  Enhance = 'Enhance',
  Basic = 'Basic'
}

/** Addon状态的主题映射 */
export const AddonStatusThemeMap = {
  running: 'success',
  failed: 'danger',
  checking: 'text',
  initializing: 'text',
  reinitializing: 'text',
  upgrading: 'text',
  '-': 'text'
};

/** Addon状态的中文名称映射 */
export const AddonStatusNameMap = {
  running: '运行中',
  failed: '异常',
  checking: '检查中',
  initializing: '初始化中',
  reinitializing: '重新初始化中',
  upgrading: '升级中',
  '-': '未知'
};

/** Addon状态的Enum */
export enum AddonStatusEnum {
  Running = 'running',
  Failed = 'failed',
  Checking = 'checking',
  Initializing = 'initializing',
  Reinitializing = 'reinitializing',
  Upgrading = 'upgrading'
}

/** Addon的中文名称映射 */
export const AddonNameMap = {
  VolumeDecorator: 'VolumeDecorator组件',
  TappController: 'TappController组件',
  Prometheus: 'Prometheus组件',
  PersistentEvent: '事件持久化组件',
  LogCollector: '日志采集组件',
  LBCF: 'LBCF组件',
  IPAM: 'IPAM组件',
  Helm: 'Helm应用管理组件',
  Galaxy: 'Galaxy组件',
  GPUManager: 'GPUManager组件',
  CronHPA: 'CronHPA组件',
  CoreDNS: 'CoreDNS组件',
  CSIOperator: 'CSIOperator组件'
};

/** 所有addon的名称 */
export enum AddonNameEnum {
  VolumeDecorator = 'VolumeDecorator',
  TappController = 'TappController',
  Prometheus = 'Prometheus',
  PersistentEvent = 'PersistentEvent',
  LogCollector = 'LogCollector',
  LBCF = 'LBCF',
  IPAM = 'IPAM',
  Helm = 'Helm',
  Galaxy = 'Galaxy',
  GPUManager = 'GPUManager',
  CronHPA = 'CronHPA',
  CoreDNS = 'CoreDNS',
  CSIOperator = 'CSIOperator'
}

/** 创建addon generatedName映射，这个与apiVersion当中的headTitle的配置有关 */
export const AddonNameMapToGenerateName = {
  VolumeDecorator: 'vd',
  TappController: 'tapp',
  Prometheus: 'prometheus',
  PersistentEvent: 'pe',
  LogCollector: 'lc',
  LBCF: 'lbcf',
  IPAM: 'ipam',
  Helm: 'hm',
  Galaxy: 'galaxy',
  GPUManager: 'gm',
  CronHPA: 'cronhpa',
  CoreDNS: 'coredns',
  CSIOperator: 'csio'
};

/** 资源的映射名称 大写 => resourceConfig的名称映射，与resourceConfig字段当中的相同 */
export const ResourceNameMap: { [props: string]: ApiVersionKeyName } = {
  Helm: 'addon_helm',
  PersistentEvent: 'addon_persistentevent',
  GPUManager: 'addon_gpumanager',
  LogCollector: 'addon_logcollector',
  TappController: 'addon_tappcontroller',
  CSIOperator: 'addon_csioperator',
  LBCF: 'addon_lbcf',
  CronHPA: 'addon_cronhpa',
  CoreDNS: 'addon_coredns',
  Galaxy: 'addon_galaxy',
  Prometheus: 'addon_prometheus',
  VolumeDecorator: 'addon_volumedecorator',
  IPAM: 'addon_ipam'
};
/** ========================= end Addon的相关配置项 ======================== */
