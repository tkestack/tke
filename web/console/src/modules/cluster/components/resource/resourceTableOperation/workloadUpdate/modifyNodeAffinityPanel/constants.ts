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

export const affinityRuleSchema = z.array(
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
            } else if (values.some(item => !validateJS.matches(item, /^([A-Za-z0-9][-A-Za-z0-9_./]*)?[A-Za-z0-9]$/))) {
              message = t('标签格式不正确');
            }
          }

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
);

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

export enum AffinityTypeEnum {
  Force = 'force',
  Attempt = 'attempt'
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

const tolerationSchema = z.array(
  z
    .object({
      key: z.string(),
      operator: z.nativeEnum(TolerationOperatorEnum),
      value: z.string(),
      effect: z.nativeEnum(TolerationEffectEnum),
      time: z.number().min(0)
    })
    .superRefine(({ key, operator, value, effect, time }, ctx) => {
      if (operator === TolerationOperatorEnum.Equal && !key) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['key'],
          message: 'key 不能为空'
        });
      }

      if (operator === TolerationOperatorEnum.Equal && !value) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ['value'],
          message: 'value 不能为空'
        });
      }
    })
);

export const nodeAffinitySchema = z.object({
  affinityRules: z.object({
    [AffinityTypeEnum.Force]: affinityRuleSchema,
    [AffinityTypeEnum.Attempt]: affinityRuleSchema
  }),
  tolerationRules: tolerationSchema,
  nodeAffinityType: z.nativeEnum(NodeAffinityTypeEnum),
  tolerationType: z.nativeEnum(TolerationTypeEnum)
});

export type NodeAffinityFormType = z.infer<typeof nodeAffinitySchema>;

export const defaultNodeAffinityFormData: NodeAffinityFormType = {
  affinityRules: {
    force: [],
    attempt: []
  },
  tolerationRules: [],
  nodeAffinityType: NodeAffinityTypeEnum.Unset,
  tolerationType: TolerationTypeEnum.UnSet
};

export enum ScheduleTypeEnum {
  ScheduleByNode = 'scheduleByNode',
  ScheduleByLabel = 'scheduleByLabel'
}

export const appendAffinityRuleSchema = z.object({
  rules: affinityRuleSchema
});

export type AppendAffinityRuleFormType = z.infer<typeof appendAffinityRuleSchema>;

export const defaultAppendAffinityRuleFormData: AppendAffinityRuleFormType = {
  rules: generateDefaultRules()
};
