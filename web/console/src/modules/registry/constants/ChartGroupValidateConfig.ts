import { ModelTypeEnum, RuleTypeEnum, ValidateSchema, ValidatorStatusEnum } from '@tencent/ff-validator';
import { Validation } from '../../common/models';
import { UserInfo } from '../models';

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
        },
        {
          type: RuleTypeEnum.custom,
          customFunc: (value, store, extraStore): Validation => {
            let name = extraStore[0] ? extraStore[0].name : '';
            let status = ValidatorStatusEnum.Init,
              message = '';
            if (store.spec.type === 'personal') {
              if (value === name) {
                status = ValidatorStatusEnum.Success;
                message = t('');
              } else {
                status = ValidatorStatusEnum.Failed;
                message = t('仓库类型为个人时，仓库名称需要是登录账号名');
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
    // {
    //   vKey: 'spec.displayName',
    //   label: t('仓库别名'),
    //   rules: [
    //     RuleTypeEnum.isRequire,
    //     {
    //       type: RuleTypeEnum.maxLength,
    //       limit: 60
    //     }
    //   ]
    // },
    {
      vKey: 'spec.description',
      label: t('描述'),
      rules: [
        {
          type: RuleTypeEnum.maxLength,
          limit: 255
        }
      ]
    },
    {
      vKey: 'spec.type',
      label: t('仓库类型'),
      rules: [RuleTypeEnum.isRequire]
    },
    {
      vKey: 'spec.visibility',
      label: t('仓库权限'),
      rules: [RuleTypeEnum.isRequire]
    }
  ]
};
