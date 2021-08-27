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
import { Identifiable } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Validation, initValidator } from '@tencent/ff-validator';

/* eslint-disable */
export namespace ContainerEnv {
  export interface Item {
    /** 当前的环境变量的类型 */
    type: EnvTypeEnum;

    /** 环境变量的名称 */
    name: string;
    v_name: Validation;

    /** 自定义环境变量的值 */
    value: string;

    /** configMap的名称 */
    configMapName: string;
    configMapDataKey: string;
    v_configMapName: Validation;
    v_configMapDataKey: Validation;

    /** secret的相关配置 */
    secretName: string;
    secretDataKey: string;
    v_secretName: Validation;
    v_secretDataKey: Validation;

    /** FieldRef的Key */
    fieldName: FieldKeyNameEnum;
    apiVersion: string;

    /** ResourceFieldRef的Key */
    resourceFieldName: ResourceFieldKeyNameEnum;
    divisor: string;
  }

  /** 包含uuid的数据 */
  export interface ItemWithId extends Item, Identifiable {}

  /** 环境变量的类型 */
  export enum EnvTypeEnum {
    /** 用户自定义环境变量 */
    UserDefined = 'UserDefined',

    /** Secret */
    SecretKeyRef = 'SecretKeyRef',

    /** Configmap */
    ConfigMapRef = 'ConfigMapKeyRef',

    /** field */
    FieldRef = 'FieldRef',

    /** resourceFieldRef */
    ResourceFieldRef = 'ResourceFieldRef'
  }

  /** Field的类型 */
  export enum FieldKeyNameEnum {
    /** metadata.name */
    MetadataName = 'metadata.name',

    /** metadata.namespace */
    MetadataNs = 'metadata.namespace',

    /** metadata.labels */
    MetadataLabels = 'metadata.labels',

    /** metadata.annotations */
    MetadataAnnotations = 'metadata.annotations',

    /** spec.nodeName */
    SpecNodeName = 'spec.nodeName',

    /** spec.serviceAccountName */
    SpecServiceAccountName = 'spec.serviceAccountName',

    /** status.hostIP */
    StatusHostIP = 'status.hostIP',

    /** status.podIP */
    StatusPodIP = 'status.podIP',

    /** status.podIPs */
    StatusPodIPs = 'status.podIPs'
  }

  /** ResourceField的类型 */
  export enum ResourceFieldKeyNameEnum {
    /** limits.cpu */
    LimitsCPU = 'limits.cpu',

    /** limits.memory */
    LimitsMem = 'limits.memory',

    /** limits.ephemeral-storage */
    LimitsES = 'limits.ephemeral-storage',

    /** requests.cpu */
    RequestCPU = 'requests.cpu',

    /** requests.memory */
    RequestMem = 'requests.memory',

    /** requests.ephemeral-storage */
    RequestES = 'requests.ephemeral-storage'
  }

  /** envItem的初始值 */
  export const initEnvItem: Item = {
    type: EnvTypeEnum.UserDefined,
    name: '',
    v_name: initValidator,
    value: '',
    configMapName: '',
    configMapDataKey: '',
    v_configMapName: initValidator,
    v_configMapDataKey: initValidator,
    secretName: '',
    secretDataKey: '',
    v_secretName: initValidator,
    v_secretDataKey: initValidator,
    fieldName: FieldKeyNameEnum.MetadataName,
    apiVersion: 'v1',
    resourceFieldName: ResourceFieldKeyNameEnum.LimitsCPU,
    divisor: '1'
  };

  /** FieldOptions列表 */
  export const FieldRefOptions = [
    {
      value: FieldKeyNameEnum.MetadataName,
      text: FieldKeyNameEnum.MetadataName
    },
    {
      value: FieldKeyNameEnum.MetadataNs,
      text: FieldKeyNameEnum.MetadataNs
    },
    {
      value: FieldKeyNameEnum.MetadataLabels,
      text: FieldKeyNameEnum.MetadataLabels
    },
    {
      value: FieldKeyNameEnum.MetadataAnnotations,
      text: FieldKeyNameEnum.MetadataAnnotations
    },
    {
      value: FieldKeyNameEnum.SpecNodeName,
      text: FieldKeyNameEnum.SpecNodeName
    },
    {
      value: FieldKeyNameEnum.SpecServiceAccountName,
      text: FieldKeyNameEnum.SpecServiceAccountName
    },
    {
      value: FieldKeyNameEnum.StatusHostIP,
      text: FieldKeyNameEnum.StatusHostIP
    },
    {
      value: FieldKeyNameEnum.StatusPodIP,
      text: FieldKeyNameEnum.StatusPodIP
    },
    {
      value: FieldKeyNameEnum.StatusPodIPs,
      text: FieldKeyNameEnum.StatusPodIPs
    }
  ];

  /** ResourceFieldOptions列表 */
  export const ResourceFieldRefOptions = [
    {
      value: ResourceFieldKeyNameEnum.LimitsCPU,
      text: ResourceFieldKeyNameEnum.LimitsCPU
    },
    {
      value: ResourceFieldKeyNameEnum.LimitsMem,
      text: ResourceFieldKeyNameEnum.LimitsMem
    },
    {
      value: ResourceFieldKeyNameEnum.LimitsES,
      text: ResourceFieldKeyNameEnum.LimitsES
    },
    {
      value: ResourceFieldKeyNameEnum.RequestCPU,
      text: ResourceFieldKeyNameEnum.RequestCPU
    },
    {
      value: ResourceFieldKeyNameEnum.RequestMem,
      text: ResourceFieldKeyNameEnum.RequestMem
    },
    {
      value: ResourceFieldKeyNameEnum.RequestES,
      text: ResourceFieldKeyNameEnum.RequestES
    }
  ];

  /** EnvTypeOptions列表 */
  export const EnvTypeOptions = [
    {
      value: EnvTypeEnum.UserDefined,
      text: t('自定义')
    },
    {
      value: EnvTypeEnum.ConfigMapRef,
      text: 'ConfigMap'
    },
    {
      value: EnvTypeEnum.SecretKeyRef,
      text: 'Secret'
    },
    {
      value: EnvTypeEnum.FieldRef,
      text: 'Field'
    },
    {
      value: EnvTypeEnum.ResourceFieldRef,
      text: 'ResourceField'
    }
  ];
}
