import { ValidateSchema, RuleTypeEnum } from '../../common/validate/Model';

export const clusterValidateSchema: ValidateSchema = {
  formKey: 'cluster',
  fields: [
    {
      label: '集群名称',
      vKey: 'clusterName',
      rules: [
        RuleTypeEnum.isRequire,
        { type: RuleTypeEnum.minLength, limit: 0 },
        { type: RuleTypeEnum.maxLength, limit: 5 }
      ]
    },
    {
      label: '集群描述',
      vKey: 'description',
      rules: [RuleTypeEnum.isRequire]
    }
  ]
};
