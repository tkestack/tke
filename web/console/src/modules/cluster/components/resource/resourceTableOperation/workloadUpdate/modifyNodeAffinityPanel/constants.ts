import { t } from '@/tencent/tea-app/lib/i18n';

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
    tip: t('Label的value在列表中')
  },
  {
    value: NodeAffinityOperatorEnum.NotIn,
    tip: t('Label的value不在列表中')
  },
  {
    value: NodeAffinityOperatorEnum.Exists,
    tip: t('Label的key存在')
  },
  {
    value: NodeAffinityOperatorEnum.DoesNotExist,
    tip: t('Labe的key不存在')
  },
  {
    value: NodeAffinityOperatorEnum.Gt,
    tip: t('Label的值大于列表值（字符串匹配）')
  },
  {
    value: NodeAffinityOperatorEnum.Lt,
    tip: t('Label的值小于列表值（字符串匹配）')
  }
];
