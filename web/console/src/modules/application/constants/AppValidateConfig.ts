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

import { t } from '@tencent/tea-app/lib/i18n';

export const AppValidateSchema: ValidateSchema = {
  formKey: 'AppValidator',
  fields: [
    {
      vKey: 'spec.name',
      label: t('应用名称'),
      rules: [
        RuleTypeEnum.isRequire,
        {
          type: RuleTypeEnum.maxLength,
          limit: 60
        },
        {
          type: RuleTypeEnum.regExp,
          limit: /^([a-z]([-a-z0-9]*[a-z0-9])?)*$/
        }
      ]
    },
    {
      vKey: 'spec.targetCluster',
      label: t('运行集群'),
      rules: [RuleTypeEnum.isRequire]
    },
    {
      vKey: 'metadata.namespace',
      label: t('命名空间'),
      rules: [RuleTypeEnum.isRequire]
    },
    {
      vKey: 'spec.type',
      label: t('类型'),
      rules: [RuleTypeEnum.isRequire]
    },
    {
      vKey: 'spec.chart',
      label: t('Chart'),
      rules: [
        RuleTypeEnum.isRequire,
        {
          type: RuleTypeEnum.custom,
          customFunc: (value, store, extraStore): Validation => {
            let status = ValidatorStatusEnum.Init,
              message = '';
            if (!store.spec.chart) {
              status = ValidatorStatusEnum.Failed;
              message = t('Chart不能为空');
            } else if (store.spec.chart.chartGroupName === '') {
              status = ValidatorStatusEnum.Failed;
              message = t('Chart仓库不能为空');
            } else if (store.spec.chart.chartName === '') {
              status = ValidatorStatusEnum.Failed;
              message = t('Chart不能为空');
            } else if (store.spec.chart.chartVersion === '') {
              status = ValidatorStatusEnum.Failed;
              message = t('Chart版本不能为空');
            } else {
              status = ValidatorStatusEnum.Success;
              message = t('');
            }
            return {
              status,
              message
            };
          }
        }
      ]
    }
  ]
};
