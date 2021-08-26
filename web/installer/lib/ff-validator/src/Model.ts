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

import { Validation } from './Validation';

export interface ValidatorModel {
  [props: string]: Validation;
}

export interface ValidateSchema {
  /** 表单校验的唯一key，如cluster、service、workload等 */
  formKey: string;

  /** 校验配置项 */
  fields: FieldConfig[];
}

export interface FieldConfig {
  /** 标签名 */
  label: string;

  /** 校验唯一id值，和store当中定义的key相同，即被校验项 */
  vKey: string;

  /** 实际上的值的获取，如果不设置valueField，默认只会取 store[vKey] 作为被校验值 */
  valueField?: string | ((value: any) => any);

  /** 前置检查条件，是否需要进行校验，如果前置条件不通过，则不需要进行rules的检验，并且会初始化校验结果 */
  condition?: (value: any, store: any) => boolean;

  /** 是否依赖于其他字段的校验结果，A --(watch)--> B，当B变化的时候，A也需要校验 */
  watchKey?: string;

  /** 校验的方法集合 */
  rules?: (Rule | string)[];

  /** 当前校验的数据类型，是否为特殊处理，如ff-redux */
  modelType?: ModelTypeEnum;
}

export interface Rule {
  /** 规则类型 */
  type: RuleTypeEnum;

  /** 限制条件，非必传 */
  limit?: any;

  /** 自定义校验函数 */
  customFunc?: (value: any, store: any, extraStore: any) => Validation;

  /** 描述内容，非必传 */
  errorTip?: React.ReactNode;
}

export enum ModelTypeEnum {
  /** ff-redux */
  FFRedux = 'ff-redux',

  /** 正常的数据类型 number、string、undefined、null、boolean等 */
  Normal = 'normal'
}

export enum RuleTypeEnum {
  /** 是否为必须 */
  isRequire = 'isRequire',

  /** 字符串最大长度 */
  maxLength = 'maxLength',

  /** 字符串最小长度 */
  minLength = 'minLength',

  /** 最小值 */
  minValue = 'minValue',

  /** 最大值 */
  maxValue = 'maxValue',

  /** 自定义正则表达式 */
  regExp = 'regExp',

  /** 数组最小选择数 */
  minSelect = 'minSelect',

  /** 数组最大选择数 */
  maxSelect = 'maxSelect',

  /** 用户自定义校验 */
  custom = 'custom'
}

export enum ValidatorStatusEnum {
  /** 初始状态 */
  Init,
  /** 正确状态 */
  Success,
  /** 错误状态 */
  Failed
}

export interface ValidationIns {
  /** 具体的校验方法，不传 vKey默认校验所有 */
  validate: (vKey?: string | string[], callback?: (validateResult: ValidatorModel) => void) => void;
}

export interface ValidateMethodOptions {
  /** 传入的全局的state */
  store: any;

  /** 需被校验的值 */
  value: any;

  /** 校验项的label */
  label: string;

  /** 校验的规则 */
  rule: Rule;

  /** 当前校验的数据类型，是否为特殊处理，如ff-redux */
  modelType: ModelTypeEnum;

  /** 校验所依赖的额外数据项 */
  extraStore: any;
}
