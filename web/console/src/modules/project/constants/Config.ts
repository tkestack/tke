/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

export const VALIDATE_PASSWORD_RULE = {
  pattern: /^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9!@#$%^&*-_=+]{10,16}$/,
  message: t('长10~16位，需包括大小写字母及数字')
};

export const VALIDATE_PHONE_RULE = {
  pattern: /^1[3|4|5|7|8][0-9]{9}$/,
  message: t('请输入正确的手机号')
};

export const VALIDATE_EMAIL_RULE = {
  pattern: /^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/,
  message: t('请输入正确的邮箱')
};

export const VALIDATE_NAME_RULE = {
  //   pattern: /^[a-z]([-a-z0-9]{0,18}[a-z0-9])?$/,
  //   message: t('长1~20位，需以小写字母开始，小写字母或数字结尾，包含小写字母、数字、-')
  pattern: /^[a-z0-9][-a-z0-9]{1,30}[a-z0-9]$/,
  message: t('长3~32位，需以小写字母或数字开头结尾，中间包含小写字母、数字、-')
};

export const STRATEGY_TYPE = ['自定义策略', '预设策略'];

export const FFReduxActionName = {
  ProjectUserInfo: 'ProjectUserInfo',
  NamespaceKubectlConfig: 'NamespaceKubectlConfig',
  UserManagedProjects: 'UserManagedProjects',
  UserInfo: 'UserInfo'
};
export enum PlatformTypeEnum {
  /** 平台 */
  Manager = 'manager',

  /** 业务 */
  Business = 'business'
}
