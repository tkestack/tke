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
  scProvisioner: string;
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

export const vmMonitorGroups = [
  {
    by: ['name'],
    fields: [
      {
        expr: mean('vm_cpu_usage_rate'),
        alias: 'CPU 使用率',
        unit: '%'
      },

      {
        expr: mean('vm_memory_usage'),
        alias: '内存使用量',
        unit: 'MB'
      },

      {
        expr: mean('vm_memory_usage_rate'),
        alias: '内存使用率',
        unit: '%'
      },

      {
        expr: mean('vm_network_transmit_bw'),
        alias: '网络出带宽',
        unit: 'Mbps'
      },

      {
        expr: mean('vm_network_receive_bw'),
        alias: '网络入带宽',
        unit: 'Mbps'
      },

      {
        expr: mean('vm_network_transmit_packets_rate'),
        alias: '网络出包量',
        unit: '个/s'
      },

      {
        expr: mean('vm_network_receive_packets_rate'),
        alias: '网络入包量',
        unit: '个/s'
      }
    ]
  },

  {
    by: ['name', 'drive'],
    fields: [
      {
        expr: mean('vm_storage_read_bw'),
        alias: '磁盘读流量',
        unit: 'MB/s'
      },

      {
        expr: mean('vm_storage_write_bw'),
        alias: '磁盘写流量',
        unit: 'MB/s'
      },

      {
        expr: mean('vm_storage_read_iops'),
        alias: '磁盘读IOPS （基于每块磁盘）',
        unit: '次/s'
      },

      {
        expr: mean('vm_storage_write_iops'),
        alias: '磁盘写IOPS （基于每块磁盘）',
        unit: '次/s'
      }
    ]
  }
];
