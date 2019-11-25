import { ProjectEdition } from './../models/Project';
import { initValidator } from './../../common/models/Validation';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { uuid } from '@tencent/qcloud-lib';
import { NamespaceEdition } from '../models';
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
export enum K8SUNIT {
  m = 'm',
  unit = 'unit',
  K = 'k',
  M = 'M',
  G = 'G',
  T = 'T',
  P = 'P',
  Ki = 'Ki',
  Mi = 'Mi',
  Gi = 'Gi',
  Ti = 'Ti',
  Pi = 'Pi'
}

export function valueLabels1000(value, targetUnit) {
  return transformField(
    value,
    1000,
    3,
    [K8SUNIT.unit, K8SUNIT.K, K8SUNIT.M, K8SUNIT.G, K8SUNIT.T, K8SUNIT.P],
    targetUnit
  );
}

export function valueLabels1024(value, targetUnit) {
  return transformField(
    value,
    1024,
    3,
    [K8SUNIT.unit, K8SUNIT.Ki, K8SUNIT.Mi, K8SUNIT.Gi, K8SUNIT.Ti, K8SUNIT.Pi],
    targetUnit
  );
}

const UNITS = [K8SUNIT.unit, K8SUNIT.Ki, K8SUNIT.Mi, K8SUNIT.Gi, K8SUNIT.Ti, K8SUNIT.Pi];

/**
 * 进行单位换算
 * 实现k8s数值各单位之间的相互转换
 * @param {string} value
 * @param {number} thousands
 * @param {number} toFixed
 */
export function transformField(_value: string, thousands, toFixed = 3, units = UNITS, targetUnit: K8SUNIT) {
  let reg = /^(\d+(\.\d{1,2})?)([A-Za-z]+)?$/;
  let value;
  let unitBase;
  if (reg.test(_value)) {
    [value, unitBase] = [+RegExp.$1, RegExp.$3];
    if (unitBase === '') {
      unitBase = K8SUNIT.unit;
    }
  } else {
    return '0';
  }

  let i = units.indexOf(unitBase),
    targetI = units.indexOf(targetUnit);
  if (thousands) {
    if (targetI >= i) {
      while (i < targetI) {
        value /= thousands;
        ++i;
      }
    } else {
      while (targetI < i) {
        value *= thousands;
        ++targetI;
      }
    }
  }
  let svalue;
  if (value > 1) {
    svalue = value.toFixed(toFixed);
    svalue = svalue.replace(/0+$/, '');
    svalue = svalue.replace(/\.$/, '');
  } else if (value) {
    // 如果数值很小，保留toFixed位有效数字
    let tens = 0;
    let v = Math.abs(value);
    while (v < 1) {
      v *= 10;
      ++tens;
    }
    svalue = value.toFixed(tens + toFixed - 1);
    svalue = svalue.replace(/0+$/, '');
    svalue = svalue.replace(/\.$/, '');
  } else {
    svalue = value;
  }
  return String(svalue);
}
