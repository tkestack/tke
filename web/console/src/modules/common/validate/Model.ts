import { Validation } from '../models';

export interface ValidateSchema {
  /** 表单校验的唯一key，如cluster、service、workload等 */
  formKey: string;

  /** 校验配置项 */
  fields: FieldConfig[];
}

export interface FieldConfig {
  /** 标签名 */
  label: string;

  /** 校验唯一id值 */
  vKey: string;

  /** 是否依赖于其他字段的校验结果 */
  dependentKey?: string;

  /** 校验的方法集合 */
  rules?: (Rule | string)[];
}

export interface Rule {
  /** 规则类型 */
  type: RuleTypeEnum;

  /** 限制条件，非必传 */
  limit?: any;

  /** 自定义校验函数 */
  customFunc?: (value: any, store: any) => Validation;

  /** 描述内容，非必传 */
  errorTip?: React.ReactNode;
}

export enum RuleTypeEnum {
  /** 是否为必须 */
  isRequire = 'isRequire',

  /** 最大长度 */
  maxLength = 'maxLength',

  /** 最小长度 */
  minLength = 'minLength',

  /** 最小值 */
  minValue = 'minValue',

  /** 最大值 */
  maxValue = 'maxValue',

  /** 自定义正则表达式 */
  regExp = 'regExp',

  /** checkBox最小选择数 */
  minCheckBoxCount = 'minCheckBoxCount',

  /** checkBox最大选择数 */
  maxCheckBoxCount = 'maxCheckBoxCount',

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

export interface ValidatorStore {
  [props: string]: Validation;
}

export interface ValidationIns {
  /** 获取校验结果，返回 initValidator的格式 */
  getValue: (vKey?: string) => Validation[] | Validation;

  /** 具体的校验方法，不传 vKey默认校验所有 */
  validate: (vKey?: string | string[]) => void;

  /** 校验结果是否成功 */
  isValid: (vKey?: string) => boolean;
}
