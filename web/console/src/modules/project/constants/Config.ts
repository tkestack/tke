import { uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { initValidator } from '../../common/models/Validation';
import { NamespaceEdition } from '../models';
import { ProjectEdition } from '../models/Project';

export const canNotOperateCluster = {
  clusterStatus: ['Creating', 'Deleting', 'Upgrading', 'Isolated', 'Abnormal']
};

/** 集群状态 */
export const clusterStatus = {
  normal: {
    text: t('运行中'),
    classname: 'text-success'
  },
  '-': {
    text: '-',
    classname: 'text-restart'
  }
};

/** 节点状态 */
/**
 * namesapce状态
 */
export const NamespaceStatus = {
  Available: 'success',
  Terminating: 'label',
  Pending: 'label',
  Failed: 'danger'
};

export const resourceLimitTypeList = [
  {
    text: t('Pod数目'),
    value: 'pods'
  },
  {
    text: t('Services数目'),
    value: 'services'
  },
  {
    text: t('quotas'),
    value: 'resourcequotas'
  },
  {
    text: t('Secrets数目'),
    value: 'secrets'
  },
  {
    text: t('Configmaps数目'),
    value: 'configmaps'
  },
  {
    text: t('PVC数目'),
    value: 'persistentvolumeclaims'
  },
  {
    text: t('nodeports模式服务数目'),
    value: 'services.nodeports'
  },
  {
    text: t('LB模式服务数目'),
    value: 'services.loadbalancers'
  },
  {
    text: t('CPU Request'),
    value: 'requests.cpu'
  },
  {
    text: t('Mem Request'),
    value: 'requests.memory'
  },
  {
    text: t('Local ephemeral storage request'),
    value: 'requests.ephemeral-storage'
  },
  {
    text: t('CPU Limits'),
    value: 'limits.cpu'
  },
  {
    text: t('Mem Limits'),
    value: 'limits.memory'
  },
  {
    text: t('Local ephemeral storage Limits'),
    value: 'limits.ephemeral-storage'
  }
];

export const resourceTypeToUnit = {
  pods: t('个'),
  services: t('个'),
  resourcequotas: t('个'),
  secrets: t('个'),
  configmaps: t('个'),
  persistentvolumeclaims: t('个'),
  'services.nodeports': t('个'),
  'services.loadbalancers': t('个'),
  'requests.cpu': t('核'),
  'limits.cpu': t('核'),
  'requests.memory': 'MiB',
  'limits.memory': 'MiB',
  'requests.ephemeral-storage': 'MiB',
  'limits.ephemeral-storage': 'MiB'
};
export const resourceLimitTypeToText = {
  pods: t('Pod数目'),
  services: t('Services数目'),
  resourcequotas: t('quotas'),
  secrets: t('Secrets数目'),
  configmaps: t('Configmaps数目'),
  persistentvolumeclaims: t('PVC数目'),
  'services.nodeports': t('nodeports模式服务数目'),
  'services.loadbalancers': t('LB模式服务数目'),
  'requests.cpu': t('CPU Request'),
  'limits.cpu': t('CPU Limits'),
  'requests.memory': t('Mem Request'),
  'limits.memory': t('Mem Limits'),
  'requests.ephemeral-storage': t('Local ephemeral storage Request'),
  'limits.ephemeral-storage': t('Local ephemeral storage Limits')
};

export const initProjectResourceLimit = {
  type: 'requests.cpu',
  value: '',
  v_type: initValidator,
  v_value: initValidator
};

export const initProjectEdition: ProjectEdition = {
  id: '',
  resourceVersion: '',
  displayName: '',
  v_displayName: initValidator,
  members: [],
  clusters: [{ name: '', v_name: initValidator, resourceLimits: [] }],

  parentProject: '',

  status: {}
};

export const projectStatus = {
  Active: 'success',
  Terminating: 'label',
  Failed: 'danger'
};

export const initNamespaceEdition: NamespaceEdition = {
  id: '',
  resourceVersion: '',
  clusterName: '',
  v_clusterName: initValidator,
  namespaceName: '',
  v_namespaceName: initValidator,
  resourceLimits: [],
  status: {}
};
