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

// import { Validation } from '../models';
// import { isEmpty } from '../utils';
// import { ReduxAction } from '@tencent/ff-redux';

// export enum ValidateState {
//   /** 数据为初始状态 */
//   Ready = 'Ready' as any,

//   /** 数据校验中*/
//   Validating = 'Validating' as any,

//   /** 数据校验失败*/
//   Failed = 'Failed' as any
// }

// export enum ValidatorTrigger {
//   /**
//    * trigger a load operation
//    */
//   Start = 'Start' as any,

//   /**
//    * trigger when load for the tolerance duration
//    * */
//   Loading = 'Loading' as any,

//   /** trigger a receive operation */
//   Done = 'Done' as any,

//   /** trigger a failed result */
//   Fail = 'Fail' as any,

//   /** trigger a manual update */
//   Update = 'Update' as any
// }

// export interface Validation {
//   /**验证状态 0: 初始状态；1：校验通过；2：校验不通过；*/
//   status?: number;

//   /**结果描述 */
//   message?: string;
// }

// /**
//  * 校验规则
//  *
//  * @export
//  * @interface Rule
//  */
// export interface ValidatorRule {
//   /**标签名 */
//   label?: string;

//   /**是否为必填 */
//   required?: boolean;

//   /**校验正则 */
//   reg?: RegExp;

//   /**最小长度 */
//   minLength?: number;

//   /**最大长度 */
//   maxLength?: number;

//   /**最小值, 当值为数值时有效 */
//   minValue?: number;

//   /**最大值, 当值为数值时有效 */
//   maxValue?: number;

//   /**自定义校验 */
//   validateFuncs?: Array<Function>;
// }

// /**
//  * 数据项校验Schema
//  */
// export interface ValidateorSchema {
//   /**校验目标 */
//   target?: string;

//   /**前置条件 */
//   condition?: boolean;

//   /**校验规则 */
//   rule?: ValidatorRule;
// }

// interface FormData {
//   [key: string]: {
//     data: any;
//     validator?: {
//       result: Validation;
//       state: ValidateState;
//       schema: ValidatorRule;
//     };
//   };
// }

// /** action payload for trigger */
// export interface ValidatorPayload<TData> {
//   trigger: ValidatorTrigger;
//   data?: TData;
//   error?: Error;
// }

// export type ValidatorAction<TData> = ReduxAction<ValidatorPayload<TData>>;

// export function generateValidator<T>(schemas: Array<Schema<T>>) {
//   /**
//    * 校验是否必填
//    *
//    * @param {any} value
//    * @param {Rule} rule
//    * @returns
//    */
//   function validateRequired(value, rule: Rule<T>) {
//     if (typeof value === 'number') {
//       return { status: 1, message: '' };
//     }

//     let result = !!value;
//     return {
//       status: result ? 1 : 2,
//       message: result ? '' : `${rule.label}不能为空`
//     };
//   }

//   /**
//    * 校验最小长度
//    *
//    * @param {any} value
//    * @param {Rule} rule
//    * @returns
//    */
//   function validateMinLength(value, rule: Rule<T>) {
//     let result = value.length >= rule.minLength;
//     return {
//       status: result ? 1 : 2,
//       message: result ? '' : t('{{label}}长度不能小于{{length}}位', { label: rule.label, length: rule.minLength })
//     };
//   }

//   /**
//    * 校验最大长度
//    *
//    * @param {any} value
//    * @param {Rule} rule
//    * @returns
//    */
//   function validateMaxLength(value, rule: Rule<T>) {
//     let result = value.length <= rule.maxLength;
//     return {
//       status: result ? 1 : 2,
//       message: result ? '' : t('{{label}}长度不能大于{{length}}位', { label: rule.label, length: rule.maxLength })
//     };
//   }

//   /**
//    * 校验最小值
//    *
//    * @param {any} value
//    * @param {Rule} rule
//    * @returns
//    */
//   function validateMinValue(value: number, rule: Rule<T>) {
//     let result = value >= rule.minValue;
//     return {
//       status: result ? 1 : 2,
//       message: result ? '' : t('{{label}}值不能小于{{value}}', { label: rule.label, value: rule.minValue })
//     };
//   }

//   /**
//    *
//    *
//    * @param {number} value
//    * @param {Rule<T>} rule
//    * @returns
//    */
//   function validateMaxValue(value: number, rule: Rule<T>): Validation {
//     let result = value <= rule.maxValue;
//     return {
//       status: result ? 1 : 2,
//       message: result ? '' : t('{{label}}值不能大于{{value}}', { label: rule.label, value: rule.maxValue })
//     };
//   }

//   /**
//    * 按正则表达式校验
//    *
//    * @param {any} value
//    * @param {Rule<T>} rule
//    * @returns
//    */
//   function validateReg(value, rule: Rule<T>) {
//     let result = rule.reg.test(value);
//     return {
//       status: result ? 1 : 2,
//       message: result ? '' : t('{{label}}格式不正确')
//     };
//   }

//   /**
//    * 校验函数
//    *
//    * @param {Object} [obj] 格式为{key: value}
//    */
//   function validate(obj?: Object, model?: T, op?: any) {
//     let re = {},
//       keys = Object.keys(obj);
//     keys.forEach(key => {
//       let schema = schemas.find(s => s.target === key),
//         result = { status: 0, message: '' } as Validation;

//       /**如果未找到相应的校验规则，则视为不校验；否则进行相应的规则校验 */
//       if (!isEmpty(schema)) {
//         let rule = schema.rule,
//           value = obj[key];

//         /**校验是否必须 */
//         if (rule.required) {
//           result = validateRequired(value, rule);
//         }

//         /**校验最小长度 */
//         if (result.status !== 2 && (rule.minLength || rule.minLength === 0)) {
//           result = validateMinLength(value, rule);
//         }

//         /**校验最大长度 */
//         if (result.status !== 2 && (rule.maxLength || rule.minLength === 0)) {
//           result = validateMaxLength(value, rule);
//         }

//         /**校验正则规则 */
//         if (result.status !== 2 && rule.reg) {
//           result = validateReg(value, rule);
//         }

//         /**校验最小值 */
//         if (result.status !== 2 && (rule.minValue || rule.minValue === 0)) {
//           result = validateMinValue(value, rule);
//         }

//         /**校验最大值 */
//         if (result.status !== 2 && (rule.maxValue || rule.maxValue === 0)) {
//           result = validateMaxValue(value, rule);
//         }

//         /**自定义校验 */
//         if (result.status !== 2 && rule.customValidate) {
//           result = rule.customValidate(value, model, op);
//         }

//         /**校验通过 */
//         if (result.status !== 2) {
//           result = { status: 1, message: '' };
//         }
//       }

//       re['v_' + key] = result;
//     });

//     return re;
//   }

//   /**
//    * 判断校验是否通过
//    *
//    * @param {Array<Validation>} validations
//    * @returns {boolean}
//    */
//   function isPassed(result): boolean {
//     let passed = true;
//     for (let key in result) {
//       if (result[key].status === 2) {
//         passed = false;
//       }
//     }

//     return passed;
//   }

//   return { validate, isPassed };
// }
