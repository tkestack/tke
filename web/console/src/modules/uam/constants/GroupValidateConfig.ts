import { ModelTypeEnum, RuleTypeEnum, ValidateSchema, ValidatorStatusEnum } from '@tencent/ff-validator';
import { t } from '@tencent/tea-app/lib/i18n';

export const GroupValidateSchema: ValidateSchema = {
  formKey: 'GroupValidator',
  fields: [
    {
      vKey: 'spec.displayName',
      label: t('名称'),
      rules: [
        RuleTypeEnum.isRequire,
        {
          type: RuleTypeEnum.maxLength,
          limit: 60
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
        }]
    }
  ]
};
