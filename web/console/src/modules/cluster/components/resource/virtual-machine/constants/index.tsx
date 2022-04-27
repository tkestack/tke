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

const mean = str => `mean(${str})`;
const sum = str => `sum(${str})`;

export const vmMonitorFields: Array<{
  expr: string;
  alias: string;
  unit: string;
}> = [
  {
    expr: mean('vm_cpu_core_total'),
    alias: 'CPU 核数',
    unit: '核'
  },

  {
    expr: mean('vm_cpu_usage'),
    alias: 'CPU 使用率',
    unit: '%'
  },

  {
    expr: mean('vm_memory_total'),
    alias: '内存大小',
    unit: 'MB'
  },

  {
    expr: mean('vm_memory_usage'),
    alias: '内存使用率',
    unit: '%'
  },

  {
    expr: mean('vm_network_transmit_bytes_bw'),
    alias: '网络上行带宽 （基于每个网卡）',
    unit: 'Mbps'
  },

  {
    expr: mean('vm_network_receive_bytes_bw'),
    alias: '网络下行带宽 （基于每个网卡）',
    unit: 'Mbps'
  },

  {
    expr: mean('vm_network_transmit_packets_rate'),
    alias: '网络收包速率 （基于每个网卡）',
    unit: 'PPS'
  },

  {
    expr: mean('vm_network_receive_packets_rate'),
    alias: '网络发包速率 （基于每个网卡）',
    unit: 'PPS'
  },

  {
    expr: mean('vm_storage_read_bw'),
    alias: '磁盘读带宽 （基于每块磁盘）',
    unit: 'MB/s'
  },

  {
    expr: mean('vm_storage_write_bw'),
    alias: '磁盘写带宽 （基于每块磁盘）',
    unit: 'MB/s'
  },

  {
    expr: mean('vm_storage_read_iops'),
    alias: '磁盘读IOPS （基于每块磁盘）',
    unit: 'IOPS'
  },

  {
    expr: mean('vm_storage_write_iops'),
    alias: '磁盘写IOPS （基于每块磁盘）',
    unit: 'IOPS'
  }
];
