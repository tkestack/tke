export enum WorkloadKindEnum {
  Deployment = 'deployment',

  StatefulSet = 'statefulset',

  DaemonSet = 'daemonset',

  Cronjob = 'cronjob'
}

export enum UpdateTypeEnum {
  ModifyStrategy = 'modifyStrategy',

  ModifyNodeAffinity = 'modifyNodeAffinity'
}

export interface IWrokloadUpdatePanelProps {
  kind: WorkloadKindEnum;

  updateType: UpdateTypeEnum;

  clusterVersion: string;
}

export const updateType2text = {
  [UpdateTypeEnum.ModifyStrategy]: '设置更新策略',

  [UpdateTypeEnum.ModifyNodeAffinity]: '更新调度策略'
};

export interface IModifyPanelProps {
  kind: WorkloadKindEnum;
  resource: any;
  title: React.ReactNode;
  baseInfo: React.ReactNode;
  onCancel: () => void;
  onUpdate: (data: any) => void;
}
