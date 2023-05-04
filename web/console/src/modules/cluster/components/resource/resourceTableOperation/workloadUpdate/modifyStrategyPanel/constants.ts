import { t } from '@/tencent/tea-app/lib/i18n';
import { WorkloadKindEnum } from '../constants';

/** 滚动更新的策略选择 */
export enum RollingUpdateTypeEnum {
  /** 启动新的pod，停止旧的pod */
  CreatePod = 'createPod',

  /** 停止旧的pod，启动新的pod */
  DestroyPod = 'destroyPod',

  /** 用户自定义 */
  UserDefined = 'userDefined'
}

export enum RegistryUpdateTypeEnum {
  /** 滚动更新 */
  RollingUpdate = 'RollingUpdate',

  /** 快速更新 */
  Recreate = 'Recreate',

  /** OnDelete */
  OnDelete = 'OnDelete'
}

export const updateStrategyOptions = [
  {
    text: t('启动新的Pod,停止旧的Pod'),
    value: RollingUpdateTypeEnum.CreatePod
  },

  {
    text: t('停止旧的Pod，启动新的Pod'),
    value: RollingUpdateTypeEnum.DestroyPod
  },

  {
    text: t('自定义'),
    value: RollingUpdateTypeEnum.UserDefined
  }
];

export const getUpdateTypeOptionsForKind = (kind: WorkloadKindEnum) => {
  const fullOptions = [
    {
      text: t('滚动更新（推荐）'),
      value: RegistryUpdateTypeEnum.RollingUpdate
    },

    {
      text: t('快速更新'),
      value: RegistryUpdateTypeEnum.Recreate
    },

    {
      text: 'OnDelete',
      value: RegistryUpdateTypeEnum.OnDelete
    }
  ];

  return fullOptions.filter(({ value }) => {
    if (kind === WorkloadKindEnum.Deployment && value !== RegistryUpdateTypeEnum.OnDelete) return true;

    if (
      (kind === WorkloadKindEnum.StatefulSet || kind === WorkloadKindEnum.DaemonSet) &&
      value !== RegistryUpdateTypeEnum.Recreate
    )
      return true;

    return false;
  });
};
