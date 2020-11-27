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
            let type = extraStore[0] ? extraStore[0].type : '';
            let status = ValidatorStatusEnum.Init,
              message = '';
            if (store.spec.type === 'Imported') {
              if (value !== '') {
                status = ValidatorStatusEnum.Success;
                message = t('');
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
