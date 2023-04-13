export enum WorkloadKindEnum {
  Deployment = 'deployment',

  StatefulSet = 'statefulset',

  DaemonSet = 'daemonset'
}

export enum UpdateTypeEnum {
  ModifyStrategy = 'modifyStrategy',

  ModifyNodeAffinity = 'modifyNodeAffinity'
}

export interface IWrokloadUpdatePanelProps {
  kind: WorkloadKindEnum;

  updateType: UpdateTypeEnum;
}

export const updateType2text = {
  [UpdateTypeEnum.ModifyStrategy]: '设置更新策略',

  [UpdateTypeEnum.ModifyNodeAffinity]: '更新调度策略'
};
