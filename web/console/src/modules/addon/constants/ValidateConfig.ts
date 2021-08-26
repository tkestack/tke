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

import { RuleTypeEnum, FieldConfig } from '../../common';
import { AddonEdit } from '../models';
import { AddonNameEnum } from './Config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { AddonValidator } from '../components/EditAddonPanel';

// interface Test {
//   a: number;
//   b: string;
// }
// type T = keyof Test;
// let a: T; // 普通的类型定义
// let e: [T] = ['a']; // Tuple 元组类型
// type TT = { a: number; b: string };
// let c: { [k in T]: string }[T] = 'asdasdasd';
// let d: 'a' extends T ? true : false;
// let numbers: { [K in T]: Test[K] extends number ? Test[K] : never }[T];
// let g: { a: number; b: string }['a'];
// let j: { [k in T]: Test[k] extends number ? Test[k] : never }[T];

// interface Test1 {
//   a: number;
//   b: {
//     c: string;
//     d: boolean;
//   };
// }

// let o: Test1;
// o.b.c;

// let i: {[k in keyof Test1]: Test1[k] extends (number | string | boolean) ? Test1[k] : keyof Test1[k]};

/** 校验规则 */
// export const addonRules: FieldConfig[] = [
//   {
//     label: t('扩展组件'),
//     vKey: 'addonName',
//     rules: [RuleTypeEnum.isRequire]
//   },
//   {
//     label: t('Elasticsearch地址'),
//     vKey: 'peEdit.esAddress',
//     condition: (value, store: AddonEdit) => {
//       return store.addonName === AddonNameEnum.PersistentEvent && store.peEdit.storeType === 'es';
//     },
//     rules: [
//       RuleTypeEnum.isRequire,
//       {
//         type: RuleTypeEnum.regExp,
//         limit: /^((http|https):\/\/)((25[0-5]|2[0-4]\d|1?\d?\d)\.){3}(25[0-5]|2[0-4]\d|1?\d?\d):([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
//         errorTip: `${t('Elasticsearch地址格式不正确')}，{{scheme}}://{{addr}}:{{port}}`
//       }
//     ]
//   },
//   {
//     label: t('索引名'),
//     vKey: 'peEdit.indexName',
//     condition: (value, store: AddonEdit) => {
//       return store.addonName === AddonNameEnum.PersistentEvent && store.peEdit.storeType === 'es';
//     },
//     rules: [
//       RuleTypeEnum.isRequire,
//       {
//         type: RuleTypeEnum.maxLength,
//         limit: 60
//       },
//       {
//         type: RuleTypeEnum.regExp,
//         limit: /^[a-z][0-9a-z_+-]+$/
//       }
//     ]
//   }
// ];

/** 校验规则 */
// export const addonRules: { [key in keyof AddonValidator]: UniqValidateMethodOptions } = {
// export const addonRules: { [key in keyof AddonValidator]: any } = {
//   addonName: {
//     label: '扩展组件',
//     rules: [RuleTypeEnum.isRequire]
//   },
//   esAddress: {
//     label: 'Elasticsearch地址',
//     isNeedValidate: (value, store: AddonEdit) => {
//       return store.addonName === AddonNameEnum.PersistentEvent && store.peEdit.storeType === 'es';
//     },
//     rules: [
//       RuleTypeEnum.isRequire,
//       {
//         type: RuleTypeEnum.regExp,
//         limit: /^((http|https):\/\/)((25[0-5]|2[0-4]\d|1?\d?\d)\.){3}(25[0-5]|2[0-4]\d|1?\d?\d):([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{4}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$/,
//         errorTip: `${t('Elasticsearch地址格式不正确')}，{{scheme}}://{{addr}}:{{port}}`
//       }
//     ]
//   },
//   indexName: {
//     label: '索引名',
//     isNeedValidate: (value, store: AddonEdit) => {
//       return store.addonName === AddonNameEnum.PersistentEvent && store.peEdit.storeType === 'es';
//     },
//     rules: [
//       RuleTypeEnum.isRequire,
//       {
//         type: RuleTypeEnum.maxLength,
//         limit: 60
//       },
//       {
//         type: RuleTypeEnum.regExp,
//         limit: /^[a-z][0-9a-z_+-]+$/
//       }
//     ]
//   },
//   logset: {
//     label: '日志集',
//     isNeedValidate: (value, store: AddonEdit) => {
//       return store.addonName === AddonNameEnum.PersistentEvent && store.peEdit.storeType === 'cls';
//     },
//     rules: [RuleTypeEnum.isRequire]
//   },
//   logsetTopic: {
//     label: '日志主题',
//     isNeedValidate: (value, store: AddonEdit) => {
//       return store.addonName === AddonNameEnum.PersistentEvent && store.peEdit.storeType === 'cls';
//     },
//     rules: [RuleTypeEnum.isRequire]
//   }
// };
