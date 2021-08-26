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

import { ModelTypeEnum, RuleTypeEnum, ValidateSchema, ValidatorStatusEnum } from '@tencent/ff-validator';
import { Validation } from '../../common/models';
import { UserInfo } from '../models';
import validatorjs from 'validator';

import { t } from '@tencent/tea-app/lib/i18n';

export const ChartGroupValidateSchema: ValidateSchema = {
  formKey: 'ChartGroupValidator',
  fields: [
    {
      vKey: 'spec.name',
      label: t('仓库名称'),
      rules: [
        RuleTypeEnum.isRequire,
        {
          type: RuleTypeEnum.maxLength,
          limit: 60
        },
        {
          type: RuleTypeEnum.regExp,
          limit: /^([A-Za-z0-9][-A-Za-z0-9_\.]*)?[A-Za-z0-9]$/
        }
      ]
    },
    {
      vKey: 'spec.visibility',
      label: t('权限范围'),
      rules: [RuleTypeEnum.isRequire]
    },
    {
      vKey: 'spec.importedInfo.addr',
      label: t('第三方仓库地址'),
      rules: [
        {
          type: RuleTypeEnum.custom,
          customFunc: (value, store, extraStore): Validation => {
            let status = ValidatorStatusEnum.Init,
              message = '';

            if (store.spec.type === 'Imported') {
              if (value !== '') {
                if (validatorjs.isURL(value)) {
                  status = ValidatorStatusEnum.Success;
                  message = t('');
                } else {
                  status = ValidatorStatusEnum.Failed;
                  message = t('仓库地址格式不正确');
                }
              } else {
                status = ValidatorStatusEnum.Failed;
                message = t('仓库类型为导入时，仓库地址不能为空');
              }
            }
            return {
              status,
              message
            };
          }
        }
      ]
    },
    {
      vKey: 'spec.description',
      label: t('描述'),
      rules: [
        {
          type: RuleTypeEnum.maxLength,
          limit: 255
        }
      ]
    }
  ]
};
