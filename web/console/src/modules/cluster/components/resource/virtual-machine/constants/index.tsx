export enum VolumeModeEnum {
  Filesystem = 'Filesystem',

  Block = 'Block'
}

export const VolumeModeOptions = [
  {
    value: VolumeModeEnum.Filesystem,
    text: '文件系统'
  },

  {
    value: VolumeModeEnum.Block,
    text: '块设备'
  }
];

export enum DiskTypeEnum {
  System = 'system',

  Data = 'data'
}

export interface DiskInterface {
  id: string;
  name: string;
  type: DiskTypeEnum;
  volumeMode: VolumeModeEnum;
  storageClass: string;
  size: number;
}

export enum ActionTypeEnum {
  Delete,

  Add,

  Modify
}

export enum VMDetailTabEnum {
  Info = 'info',

  Yaml = 'yaml',

  Log = 'log'
}

export const VMDetailTabOptions = [
  {
    id: VMDetailTabEnum.Info,
    label: '详情'
  },

  {
    id: VMDetailTabEnum.Yaml,
    label: 'YAML'
  },

  {
    id: VMDetailTabEnum.Log,
    label: '事件'
  }
];
