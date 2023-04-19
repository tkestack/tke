import { t } from '@/tencent/tea-app/lib/i18n';
import validateJS from 'validator';
import { z } from 'zod';

/** 节点亲和性调度的方式 */
export enum NodeAffinityTypeEnum {
  /** 不使用调度策略 */
  Unset = 'unset',

  /** 指定节点调度 */
  Node = 'node',

  /** 自定义调度规则 */
  Rule = 'rule'
}

export enum TolerationTypeEnum {
  UnSet = 'UnSet',

  Set = 'Set'
}

/** 节点亲和性调度 亲和性调度操作符 */
export enum NodeAffinityOperatorEnum {
  In = 'In',
  NotIn = 'NotIn',
  Exists = 'Exists',
  DoesNotExist = 'DoesNotExist',
  Gt = 'Gt',
  Lt = 'Lt'
}

export const affinityRuleOperatorList = [
  {
    value: NodeAffinityOperatorEnum.In,
    tooltip: t('Label的value在列表中')
  },
  {
    value: NodeAffinityOperatorEnum.NotIn,
    tooltip: t('Label的value不在列表中')
  },
  {
    value: NodeAffinityOperatorEnum.Exists,
    tooltip: t('Label的key存在')
  },
  {
    value: NodeAffinityOperatorEnum.DoesNotExist,
    tooltip: t('Labe的key不存在')
  },
  {
    value: NodeAffinityOperatorEnum.Gt,
    tooltip: t('Label的值大于列表值（字符串匹配）')
  },
  {
    value: NodeAffinityOperatorEnum.Lt,
    tooltip: t('Label的值小于列表值（字符串匹配）')
  }
];

export const ruleSchema = z.object({
  rules: z.array(
    z.object({
      weight: z
        .number()
        .int(t('权重必须为整数'))
        .gte(0, { message: t('权重必须在1～100之前') })
        .lte(100, { message: t('权重必须在1～100之前') }),
      subRules: z.array(
        z
          .object({
            key: z
              .string()
              .min(1, { message: t('标签名不能为空') })
              .max(63, { message: t('标签名长度不能超过63个字符') })
              .regex(/^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/, { message: '标签名格式不正确' }),
            operator: z.nativeEnum(NodeAffinityOperatorEnum),
            value: z.string()
          })
          .superRefine(({ operator, value }, ctx) => {
            console.log('superRefine', value, operator);

            const values = value.split(';');

            let message = '';

            if (operator === NodeAffinityOperatorEnum.Exists || operator === NodeAffinityOperatorEnum.DoesNotExist)
              return;

            const value0 = values?.[0];

            if (!value0) {
              message = t('自定义规则不能为空');
            } else {
              if (operator === NodeAffinityOperatorEnum.Lt || operator === NodeAffinityOperatorEnum.Gt) {
                if (values.length > 1) {
                  message = t('Gt和Lt操作符只支持一个value值');
                } else if (!validateJS.isNumeric(value0)) {
                  message = t('Gt和Lt操作符value值格式必须为数字');
                }
              } else if (values.some(item => !item)) {
                message = t('标签值不能为空');
              } else if (values.some(item => item?.length > 63)) {
                message = t('标签值长度不能超过63个字符');
              } else if (
                values.some(item => !validateJS.matches(item, /^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/))
              ) {
                message = t('标签格式不正确');
              }
            }

            console.log('message--->', message);

            if (message) {
              ctx.addIssue({
                code: z.ZodIssueCode.custom,
                message,
                path: ['value']
              });
            }
          })
      )
    })
  )
});

export const generateDefaultRules = () => [
  {
    weight: 1,
    subRules: [
      {
        key: '',
        operator: NodeAffinityOperatorEnum.In,
        value: ''
      }
    ]
  }
];

export type RuleType = z.infer<typeof ruleSchema>;

export enum TolerationOperatorEnum {
  Equal = 'Equal',
  Exists = 'Exists'
}

export const tolerationOperatorOptions = [
  {
    value: TolerationOperatorEnum.Equal
  },

  {
    value: TolerationOperatorEnum.Exists
  }
];

export enum TolerationEffectEnum {
  All = 'All',
  NoSchedule = 'NoSchedule',
  PreferNoSchedule = 'PreferNoSchedule',
  NoExecute = 'NoExecute'
}

export const tolerationEffectOptions = [
  {
    value: TolerationEffectEnum.All,
    text: '匹配全部'
  },

  {
    value: TolerationEffectEnum.NoSchedule
  },

  {
    value: TolerationEffectEnum.PreferNoSchedule
  },

  {
    value: TolerationEffectEnum.NoExecute
  }
];
